// uipack.js 鈥?capture the twelve tile effects at the end of the UI Kit page.
//
// usage: node uipack.js <url> <outDir>

const { spawn } = require("node:child_process");
const { writeFileSync, mkdirSync } = require("node:fs");
const { setTimeout: sleep } = require("node:timers/promises");

const CHROME = "C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe";
const url = process.argv[2] || "http://localhost:8099/";
const outDir = process.argv[3] || ".";

mkdirSync(outDir, { recursive: true });

const port = 9351;
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
  "--user-data-dir=" + outDir + "\\prof-" + Date.now(),
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
      const hit = [...document.querySelectorAll('*')].filter(e => e.textContent && e.textContent.trim() === ${JSON.stringify(label)});
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

  async function wheel(dy) {
    await send("Input.dispatchMouseEvent", { type: "mouseWheel", x: 620, y: 500, deltaX: 0, deltaY: dy });
    await sleep(900);
  }

  console.log(await press("UI Kit"));
  await sleep(2500);
  await wheel(1400);
  await wheel(1400);
  // The tiles replay when they mount, so a second look a moment later catches
  // them mid-pass.
  await shot("ui-tiles-a");
  await wheel(700);
  await shot("ui-tiles-b");
  await wheel(700);
  await shot("ui-tiles-c");
  await wheel(700);
  await shot("ui-tiles-d");

  ws.close();
  chrome.kill();
  process.exit(0);
}

main().catch((e) => {
  console.error(e);
  chrome.kill();
  process.exit(1);
});




