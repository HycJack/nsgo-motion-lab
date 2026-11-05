// uikit.js — capture the UI Kit page and a few of its effects with CDP.
//
// usage: node uikit.js <url> <outDir>

const { spawn } = require("node:child_process");
const { writeFileSync, mkdirSync } = require("node:fs");
const { setTimeout: sleep } = require("node:timers/promises");

const CHROME = "C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe";
const url = process.argv[2] || "http://localhost:8099/";
const outDir = process.argv[3] || ".";

mkdirSync(outDir, { recursive: true });

const port = 9335;
const chrome = spawn(CHROME, [
  "--headless=new",
  `--remote-debugging-port=${port}`,
  "--window-size=1240,900",
  "--hide-scrollbars",
  "--enable-unsafe-swiftshader",
  "--use-gl=angle",
  "--use-angle=swiftshader",
  "--no-first-run",
  "--disable-extensions",
  "--user-data-dir=" + outDir + "\\profile3",
  "about:blank",
], { stdio: "ignore" });

async function json(path) {
  const res = await fetch(`http://127.0.0.1:${port}${path}`);
  return res.json();
}

async function main() {
  let target = null;
  for (let i = 0; i < 60; i++) {
    try {
      const list = await json("/json/list");
      target = list.find((t) => t.type === "page");
      if (target) break;
    } catch {}
    await sleep(250);
  }
  if (!target) throw new Error("no CDP page target");

  const ws = new WebSocket(target.webSocketDebuggerUrl);
  await new Promise((ok, bad) => {
    ws.onopen = ok;
    ws.onerror = () => bad(new Error("ws failed"));
  });

  let id = 0;
  const pending = new Map();
  ws.onmessage = (e) => {
    const msg = JSON.parse(e.data);
    if (msg.id && pending.has(msg.id)) {
      pending.get(msg.id)(msg.result);
      pending.delete(msg.id);
    }
  };
  const send = (method, params = {}) =>
    new Promise((ok) => {
      const n = ++id;
      pending.set(n, ok);
      ws.send(JSON.stringify({ id: n, method, params }));
    });

  await send("Page.enable");
  await send("Runtime.enable");
  await send("Page.navigate", { url });
  await sleep(9000);

  async function shot(name) {
    const { data } = await send("Page.captureScreenshot", { format: "png" });
    writeFileSync(`${outDir}\\${name}.png`, Buffer.from(data, "base64"));
    console.log(`wrote ${name}.png`);
  }

  async function press(label) {
    const expr = `(() => {
      const els = [...document.querySelectorAll('*')];
      const hit = els.filter(e => e.textContent && e.textContent.trim() === ${JSON.stringify(label)});
      const el = hit[hit.length - 1];
      if (!el) return 'none';
      const box = el.getBoundingClientRect();
      return JSON.stringify({x: Math.round(box.x + box.width / 2), y: Math.round(box.y + box.height / 2)});
    })()`;
    const r = await send("Runtime.evaluate", { expression: expr, returnByValue: true });
    const v = r.result && r.result.value;
    if (!v || v === "none") return v;
    const { x, y } = JSON.parse(v);
    await send("Input.dispatchMouseEvent", { type: "mouseMoved", x, y, buttons: 0 });
    await send("Input.dispatchMouseEvent", { type: "mousePressed", x, y, button: "left", buttons: 1, clickCount: 1 });
    await sleep(50);
    await send("Input.dispatchMouseEvent", { type: "mouseReleased", x, y, button: "left", buttons: 0, clickCount: 1 });
    return `${label} at ${x},${y}`;
  }

  async function hover(label) {
    const expr = `(() => {
      const els = [...document.querySelectorAll('*')];
      const hit = els.filter(e => e.textContent && e.textContent.trim() === ${JSON.stringify(label)});
      const el = hit[hit.length - 1];
      if (!el) return 'none';
      const box = el.getBoundingClientRect();
      return JSON.stringify({x: Math.round(box.x + box.width / 2), y: Math.round(box.y + box.height / 2)});
    })()`;
    const r = await send("Runtime.evaluate", { expression: expr, returnByValue: true });
    const v = r.result && r.result.value;
    if (!v || v === "none") return v;
    const { x, y } = JSON.parse(v);
    await send("Input.dispatchMouseEvent", { type: "mouseMoved", x, y, buttons: 0 });
    return `${label} at ${x},${y}`;
  }

  async function wheel(dy) {
    await send("Input.dispatchMouseEvent", { type: "mouseWheel", x: 620, y: 500, deltaX: 0, deltaY: dy });
    await sleep(700);
  }

  // Controls without text are found by geometry: the first role=button
  // under a card's heading.
  async function pressUnderHeading(heading) {
    const expr = `(() => {
      const els = [...document.querySelectorAll('*')];
      const head = els.find(e => e.textContent && e.textContent.trim() === ${JSON.stringify(heading)});
      if (!head) return 'no-heading';
      const hb = head.getBoundingClientRect();
      const cands = [...document.querySelectorAll('[role=button]')]
        .map(e => ({ e, r: e.getBoundingClientRect() }))
        .filter(o => o.r.y >= hb.bottom - 2 && o.r.y <= hb.bottom + 140 && o.r.width < 200)
        // Nearest to the heading's own column, not merely the topmost:
        // both cards in a row start at the same y.
        .sort((a, b) => (Math.abs(a.r.x - hb.x) - Math.abs(b.r.x - hb.x)) || (a.r.y - b.r.y));
      if (!cands.length) return 'no-button';
      const r = cands[0].r;
      return JSON.stringify({ x: Math.round(r.x + r.width / 2), y: Math.round(r.y + r.height / 2) });
    })()`;
    const res = await send("Runtime.evaluate", { expression: expr, returnByValue: true });
    const v = res.result && res.result.value;
    if (!v || v.startsWith("no-")) return v;
    const { x, y } = JSON.parse(v);
    await send("Input.dispatchMouseEvent", { type: "mouseMoved", x, y, buttons: 0 });
    await send("Input.dispatchMouseEvent", { type: "mousePressed", x, y, button: "left", buttons: 1, clickCount: 1 });
    await sleep(50);
    await send("Input.dispatchMouseEvent", { type: "mouseReleased", x, y, button: "left", buttons: 0, clickCount: 1 });
    return `${heading} -> ${x},${y}`;
  }

  console.log(await press("UI Kit"));
  await sleep(1800);
  await shot("10-uikit-top");

  // Things reachable without scrolling: the toggle, the toast, the badge.
  console.log(await press("Show a toast"));
  await sleep(400);
  await shot("11-uikit-toast");

  console.log(await press("Add one"));
  await sleep(500);
  console.log(await pressUnderHeading("Toggle switch \u00b7 #059"));
  await sleep(600);
  await shot("12-uikit-toggled");

  await wheel(420);
  await shot("14-uikit-middle");

  console.log(await press("Finish"));
  await sleep(500);
  console.log(await pressUnderHeading("Hamburger menu \u00b7 #065"));
  await sleep(700);
  await shot("15-uikit-hamburger");

  // The tooltip's anchor is a View, not a button, so it has no role to
  // find: it is reached by the "?" it shows. The cursor staying there is
  // what keeps the note up.
  console.log(await hover("?"));
  await sleep(800);
  await shot("16-uikit-tooltip");

  await wheel(600);
  console.log(await press("Next"));
  await sleep(700);
  await shot("17-uikit-carousel");

  ws.close();
  chrome.kill();
}

main().catch((e) => {
  console.error("failed:", e.message);
  chrome.kill();
  process.exit(1);
});
