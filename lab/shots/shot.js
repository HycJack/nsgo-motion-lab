// shot.js 鈥?capture the NSGo web app with CDP, waiting for it to paint.
//
// Chrome's --screenshot fires at load, which for a wasm app is before the
// first frame; this drives CDP instead so each page is captured after it
// has rendered.
//
// usage: node shot.js <url> <outDir>

const { spawn } = require("node:child_process");
const { writeFileSync, mkdirSync } = require("node:fs");
const { setTimeout: sleep } = require("node:timers/promises");

const CHROME = "C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe";
const url = process.argv[2] || "http://localhost:8099/";
const outDir = process.argv[3] || ".";

mkdirSync(outDir, { recursive: true });

const port = 9333;
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
  "--user-data-dir=" + outDir + "\\profile",
  "about:blank",
], { stdio: "ignore" });

async function json(path) {
  const res = await fetch(`http://127.0.0.1:${port}${path}`);
  return res.json();
}

async function main() {
  // Wait for the debugging endpoint.
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
  // The wasm app needs to compile and mount; give it room.
  await sleep(9000);

  // Capture a named page and click a tab by its visible text.
  async function shot(name) {
    const { data } = await send("Page.captureScreenshot", { format: "png" });
    writeFileSync(`${outDir}\\${name}.png`, Buffer.from(data, "base64"));
    console.log(`wrote ${name}.png`);
  }

  async function clickTab(label) {
    const expr = `(() => {
      const els = [...document.querySelectorAll('*')];
      const hit = els.filter(e => e.textContent && e.textContent.trim() === ${JSON.stringify(label)});
      const el = hit[hit.length - 1];
      if (!el) return 'none';
      el.dispatchEvent(new MouseEvent('mousedown', {bubbles: true}));
      el.dispatchEvent(new MouseEvent('mouseup', {bubbles: true}));
      el.dispatchEvent(new MouseEvent('click', {bubbles: true}));
      el.click && el.click();
      return el.tagName + ':' + el.textContent.trim();
    })()`;
    const r = await send("Runtime.evaluate", { expression: expr, returnByValue: true });
    return r.result && r.result.value;
  }

  await shot("01-canvas");

  for (const [label, file] of [["Animation", "02-animation"], ["Shader", "03-shader"]]) {
    const r = await clickTab(label);
    console.log(`click ${label}: ${r}`);
    await sleep(2500);
    await shot(file);
  }

  // Drive the shader page: hover the dot grid, then press the ripple box.
  const rects = await send("Runtime.evaluate", {
    expression: `JSON.stringify([...document.querySelectorAll('canvas')].map(c => {
      const r = c.getBoundingClientRect();
      return {x: Math.round(r.x), y: Math.round(r.y), w: Math.round(r.width), h: Math.round(r.height)};
    }))`,
    returnByValue: true,
  });
  const boxes = JSON.parse(rects.result.value);
  console.log("shader boxes:", JSON.stringify(boxes));

  if (boxes.length >= 2) {
    const grid = boxes[1];
    const gx = grid.x + grid.w * 0.5;
    const gy = grid.y + grid.h * 0.5;
    for (let i = 0; i < 4; i++) {
      await send("Input.dispatchMouseEvent", { type: "mouseMoved", x: gx - 40 + i * 20, y: gy, buttons: 0 });
      await sleep(120);
    }
    await sleep(500);
    await shot("04-shader-hover");
  }

  if (boxes.length >= 4) {
    const ripple = boxes[3];
    const rx = ripple.x + ripple.w * 0.45;
    const ry = ripple.y + ripple.h * 0.5;
    await send("Input.dispatchMouseEvent", { type: "mouseMoved", x: rx, y: ry, buttons: 0 });
    await send("Input.dispatchMouseEvent", { type: "mousePressed", x: rx, y: ry, button: "left", buttons: 1, clickCount: 1 });
    await sleep(60);
    await send("Input.dispatchMouseEvent", { type: "mouseReleased", x: rx, y: ry, button: "left", buttons: 0, clickCount: 1 });
    await sleep(90);
    await shot("05-shader-ripple");
    await sleep(300);
    await shot("06-shader-ripple-late");
  }

  // Report whether the shader canvases really are WebGL surfaces.
  const info = await send("Runtime.evaluate", {
    expression: `JSON.stringify({
      canvases: [...document.querySelectorAll('canvas')].map(c => ({w: c.width, h: c.height, gl: !!(c.getContext('webgl2') || c.getContext('webgl'))})),
      size: [innerWidth, innerHeight]
    })`,
    returnByValue: true,
  });
  console.log("dom:", info.result && info.result.value);

  ws.close();
  chrome.kill();
}

main().catch((e) => {
  console.error("failed:", e.message);
  chrome.kill();
  process.exit(1);
});
