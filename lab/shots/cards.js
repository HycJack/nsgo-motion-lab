// cards.js — replay one pack card and capture it mid-flight and at rest.
//
// usage: node cards.js <url> <outDir>

const { spawn } = require("node:child_process");
const { writeFileSync, mkdirSync } = require("node:fs");
const { setTimeout: sleep } = require("node:timers/promises");
const path = require("node:path");

const CHROME = "C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe";
const url = process.argv[2] || "http://localhost:8099/";
const outDir = process.argv[3] || ".";
const profDir = path.resolve(outDir) + "\\prof-" + Date.now();
mkdirSync(outDir, { recursive: true });

// tab, card title, replay chip label, one pass in ms
const CARDS = [
  ["Camera", "Orbit (arc shot)", "Replay 030", 3200],
  ["Camera", "Crane up", "Replay 031", 3000],
  ["Camera2", "Fly-through", "Replay 176", 3000],
  ["Camera2", "POV shot", "Replay 179", 2600],
  ["Camera2", "Over-the-shoulder", "Replay 180", 2400],
];

const port = 9381;
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
  "--user-data-dir=" + profDir,
  "about:blank",
], { stdio: "ignore" });
chrome.on("error", (e) => console.log("chrome spawn error:", e.message));

async function json(p) {
  const res = await fetch(`http://127.0.0.1:${port}${p}`);
  return res.json();
}

async function main() {
  let target = null;
  for (let i = 0; i < 60; i++) {
    try {
      const list = await json("/json/list");
      target = list.find((t) => t.type === "page");
      if (target) break;
    } catch (e) {
      if (i % 8 === 0) console.log("poll", i, e.message);
    }
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

  const evaluate = async (expression) => {
    const r = await send("Runtime.evaluate", { expression, returnByValue: true });
    return r.result && r.result.value;
  };

  // press hits the element whose trimmed text equals label, topmost match.
  async function press(label) {
    const v = await evaluate(`(() => {
      const el = [...document.querySelectorAll('*')].filter(e => e.textContent && e.textContent.trim() === ${JSON.stringify(label)}).pop();
      if (!el) return 'none';
      const b = el.getBoundingClientRect();
      window.scrollTo({ top: window.scrollY + b.y - 220 });
      const c = el.getBoundingClientRect();
      return JSON.stringify({ x: Math.round(c.x + c.width / 2), y: Math.round(c.y + c.height / 2) });
    })()`);
    if (!v || v === "none") return v;
    const { x, y } = JSON.parse(v);
    await send("Input.dispatchMouseEvent", { type: "mouseMoved", x, y, buttons: 0 });
    await send("Input.dispatchMouseEvent", { type: "mousePressed", x, y, button: "left", buttons: 1, clickCount: 1 });
    await sleep(40);
    await send("Input.dispatchMouseEvent", { type: "mouseReleased", x, y, button: "left", buttons: 0, clickCount: 1 });
    return `${label} @ ${x},${y}`;
  }

  // cardBox scrolls the title into view and returns its card's rectangle.
  // A card's title reads "030 · Orbit (arc shot)", so match by substring and
  // take the innermost element that carries it.
  async function cardBox(title) {
    const v = await evaluate(`(() => {
      const T = ${JSON.stringify(title)};
      const cands = [...document.querySelectorAll('*')].filter(e => e.textContent && e.textContent.includes(T));
      if (!cands.length) return 'none';
      const t = cands.find(e => ![...e.children].some(c => c.textContent && c.textContent.includes(T))) || cands[cands.length - 1];
      let el = t;
      while (el.parentElement && (el.getBoundingClientRect().height < 170 || el.getBoundingClientRect().width < 260)) el = el.parentElement;
      if (!el) el = t;
      el.scrollIntoView({ block: 'center' });
      const b = el.getBoundingClientRect();
      return JSON.stringify({ x: Math.round(b.x), y: Math.round(b.y), w: Math.round(b.width), h: Math.round(b.height) });
    })()`);
    if (!v || v === "none") return null;
    return JSON.parse(v);
  }

  const shot = async (name) => {
    const { data } = await send("Page.captureScreenshot", { format: "png" });
    writeFileSync(`${outDir}\\${name}.png`, Buffer.from(data, "base64"));
    console.log(`wrote ${name}.png`);
  };

  const clipShot = async (name, box) => {
    const { data } = await send("Page.captureScreenshot", {
      format: "png",
      clip: { x: Math.max(0, box.x - 6), y: Math.max(0, box.y - 6), width: box.w + 12, height: box.h + 12, scale: 2 },
    });
    writeFileSync(`${outDir}\\${name}.png`, Buffer.from(data, "base64"));
    console.log(`wrote ${name}.png`);
  };

  for (const [tab, title, chip, dur] of CARDS) {
    console.log(await press(tab));
    await sleep(1800);
    const box = await cardBox(title);
    if (!box) {
      console.log("no card", title);
      continue;
    }
    console.log(await press(chip));
    await sleep(Math.round(dur * 0.55));
    await clipShot(`${title.toLowerCase().replace(/[^a-z0-9]+/g, "-")}-mid`, box);
    await sleep(Math.round(dur * 0.45) + 250);
    await clipShot(`${title.toLowerCase().replace(/[^a-z0-9]+/g, "-")}-end`, box);
    await shot(`tab-${tab.toLowerCase()}`);
  }

  ws.close();
  chrome.kill();
  process.exit(0);
}

main().catch((e) => {
  console.error(e);
  chrome.kill();
  process.exit(1);
});
