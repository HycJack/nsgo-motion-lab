// tabshots.js — refresh the README's camera screenshots.
//
// usage: node tabshots.js <url> <docsDir> <tmpDir>

const { spawn } = require("node:child_process");
const { writeFileSync, mkdirSync } = require("node:fs");
const { setTimeout: sleep } = require("node:timers/promises");
const path = require("node:path");

const CHROME = "C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe";
const url = process.argv[2] || "http://localhost:8099/";
const docsDir = process.argv[3] || ".";
const outDir = process.argv[4] || ".";
const profDir = path.resolve(outDir) + "\\prof-" + Date.now();

const port = 9383;
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

  async function press(label) {
    const v = await evaluate(`(() => {
      const el = [...document.querySelectorAll('*')].filter(e => e.textContent && e.textContent.trim() === ${JSON.stringify(label)}).pop();
      if (!el) return 'none';
      const b = el.getBoundingClientRect();
      return JSON.stringify({ x: Math.round(b.x + b.width / 2), y: Math.round(b.y + b.height / 2) });
    })()`);
    if (!v || v === "none") return v;
    const { x, y } = JSON.parse(v);
    await send("Input.dispatchMouseEvent", { type: "mouseMoved", x, y, buttons: 0 });
    await send("Input.dispatchMouseEvent", { type: "mousePressed", x, y, button: "left", buttons: 1, clickCount: 1 });
    await sleep(40);
    await send("Input.dispatchMouseEvent", { type: "mouseReleased", x, y, button: "left", buttons: 0, clickCount: 1 });
    return `${label}`;
  }

  // union scrolls the cards carrying these titles into view and returns the
  // viewport-relative box that covers all of them (this Chrome clips in
  // viewport space, which is also what cards.js relies on).
  async function union(titles) {
    const script = `(() => {
      const titles = ${JSON.stringify(titles)};
      const boxes = [];
      for (const T of titles) {
        const cands = [...document.querySelectorAll('*')].filter(e => e.textContent && e.textContent.includes(T));
        if (!cands.length) continue;
        const t = cands.find(e => ![...e.children].some(c => c.textContent && c.textContent.includes(T))) || cands[cands.length - 1];
        let el = t;
        while (el.parentElement && (el.getBoundingClientRect().height < 170 || el.getBoundingClientRect().width < 260)) el = el.parentElement;
        boxes.push(el);
      }
      if (!boxes.length) return 'none';
      const tops = boxes.map(e => e.getBoundingClientRect().top);
      window.scrollBy(0, Math.min(...tops) - 130);
      return String(boxes.length);
    })()`;
    const n = await evaluate(script);
    if (!n || n === "none") return null;
    await sleep(400);
    const v = await evaluate(`(() => {
      const titles = ${JSON.stringify(titles)};
      const boxes = [];
      for (const T of titles) {
        const cands = [...document.querySelectorAll('*')].filter(e => e.textContent && e.textContent.includes(T));
        if (!cands.length) continue;
        const t = cands.find(e => ![...e.children].some(c => c.textContent && c.textContent.includes(T))) || cands[cands.length - 1];
        let el = t;
        while (el.parentElement && (el.getBoundingClientRect().height < 170 || el.getBoundingClientRect().width < 260)) el = el.parentElement;
        const b = el.getBoundingClientRect();
        boxes.push({ x: b.x, y: b.y, w: b.width, h: b.height });
      }
      if (!boxes.length) return 'none';
      const x = Math.min(...boxes.map(b => b.x)), y = Math.min(...boxes.map(b => b.y));
      const r = Math.max(...boxes.map(b => b.x + b.w)), d = Math.max(...boxes.map(b => b.y + b.h));
      return JSON.stringify({ x: Math.round(x), y: Math.round(y), w: Math.round(r - x), h: Math.round(d - y), n: boxes.length });
    })()`);
    if (!v || v === "none") return null;
    return JSON.parse(v);
  }

  const pageShot = async (name) => {
    const { data } = await send("Page.captureScreenshot", { format: "png" });
    writeFileSync(path.join(docsDir, name), Buffer.from(data, "base64"));
    console.log(`wrote ${name}`);
  };

  const clipShot = async (name, box) => {
    const { data } = await send("Page.captureScreenshot", {
      format: "png",
      clip: { x: Math.max(0, box.x), y: Math.max(0, box.y), width: box.w, height: box.h, scale: 1 },
    });
    writeFileSync(path.join(docsDir, name), Buffer.from(data, "base64"));
    console.log(`wrote ${name} (${box.w}x${box.h}, ${box.n} cards)`);
  };

  await evaluate("window.scrollTo(0,0)");
  console.log(await press("Camera"));
  await sleep(2500);
  await evaluate("window.scrollTo(0,0)");
  await sleep(600);
  await pageShot("pack-camera.png");

  // The card strips are taller than the 809-point viewport, so give the page a
  // tall viewport before clipping: a clip that reaches past the viewport comes
  // back blank.
  await send("Emulation.setDeviceMetricsOverride", { width: 1240, height: 1700, deviceScaleFactor: 1, mobile: false });
  await sleep(1200);
  await evaluate("window.scrollTo(0,0)");
  let box = await union(["Orbit (arc shot)", "Crane up"]);
  await sleep(1800);
  if (box) await clipShot("camera-3d-two.png", box);

  console.log(await press("Camera2"));
  await sleep(2500);
  await evaluate("window.scrollTo(0,0)");
  box = await union(["Fly-through", "POV shot", "Over-the-shoulder"]);
  await sleep(1800);
  if (box) await clipShot("camera-3d-three.png", box);

  // Back to the normal window before the page screenshots.
  await send("Emulation.setDeviceMetricsOverride", { width: 1240, height: 900, deviceScaleFactor: 1, mobile: false });
  await evaluate("window.scrollTo(0,0)");
  await sleep(1200);
  await pageShot("pack-camera2.png");

  ws.close();
  chrome.kill();
  process.exit(0);
}

main().catch((e) => {
  console.error(e);
  chrome.kill();
  process.exit(1);
});
