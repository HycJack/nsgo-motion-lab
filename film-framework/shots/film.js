// film.js - record an NSGo film scene frame by frame into an mp4.
//
//   node film.js <url> <outDir> <durSec> [fps]
//
// The scene must read filmkit.Now(): on web that clock honors the JS
// global __filmT, so the recorder steps time one frame at a time —
// deterministically, not in real time — screenshots each frame, and then
// (ENCODE=1, the default) hands the frame folder to ffmpeg.
//
// Parallel: WORKERS chrome processes each render a contiguous slice of
// the frame range (the scene is a pure function of __filmT, so any
// worker can draw any frame), files land under their global numbers,
// ffmpeg assembles them in order. Resume: frames already on disk are
// skipped, so an interrupted run picks up where it died.
//
//   SIZE=1280x720   viewport (default 1920x1080)
//   CHROME=/path    browser binary (default: first Chrome-like app found)
//   ENCODE=0        keep frames only, skip ffmpeg
//   FPS=60          frames per second (default 30; positional [fps] wins)
//   WORKERS=4       parallel chrome processes (default: min(4, cpus))
//   READY=selector  wait for this before recording (default: canvas)
//   PRESS=text      click the element with that exact text first (a tab)
//   CLIP=selector   record only that element's rect (a card, a stage) instead
//                   of the viewport
//   HIDE=sel,sel    hide page furniture before recording (a tab bar, a chip)
//   SETTLE=ms       extra wait per frame; the default is 60 on a page whose
//                   nodes carry a [data-testid^=take-], else 0 — a keyframe
//                   take is seeked on one frame and written to the tree on the
//                   next, so it needs one more frame than a canvas scene
//
// Frames land in <outDir>/frames/f00000.jpg… and the mp4 in
// <outDir>/film.mp4.

const { spawn, spawnSync } = require("node:child_process");
const { existsSync, mkdirSync, writeFileSync, statSync } = require("node:fs");
const os = require("node:os");
const { setTimeout: sleep } = require("node:timers/promises");
const path = require("node:path");

const CHROME_CANDIDATES = [
  "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
  "/Applications/Chromium.app/Contents/MacOS/Chromium",
  "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome Canary",
  "/Applications/Microsoft Edge.app/Contents/MacOS/Google Edge",
  "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe",
  "C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe",
  "/usr/bin/google-chrome",
  "/usr/bin/chromium-browser",
];

const url = process.argv[2];
const outDir = process.argv[3];
const durSec = Number(process.argv[4] || 16);
const fps = Number(process.argv[5] || process.env.FPS || 30);
const size = (process.env.SIZE || "1920x1080").split("x");
const W = Number(size[0]), H = Number(size[1]);
const N = Math.round(fps * durSec);
const WORKERS = Math.max(1, Number(process.env.WORKERS || Math.min(4, os.cpus().length)));
const port = Number(process.env.CDP_PORT || 9386);

if (!url || !outDir) {
  console.error("usage: node film.js <url> <outDir> <durSec> [fps]");
  process.exit(2);
}
mkdirSync(outDir, { recursive: true });
mkdirSync(path.join(outDir, "frames"), { recursive: true });

const chromePath =
  process.env.CHROME || CHROME_CANDIDATES.find((p) => existsSync(p));
if (!chromePath) throw new Error("no Chrome found; set CHROME=/path/to/chrome");
console.log("browser:", chromePath, "| workers:", WORKERS);

const frameFile = (i) => path.join(outDir, "frames", `f${String(i).padStart(5, "0")}.jpg`);

// 把帧区间切成连续的段：每个工人一段。
function splitFrames(n, w) {
  const chunks = [];
  const per = Math.ceil(n / w);
  for (let k = 0; k < w; k++) {
    const from = k * per;
    if (from >= n) break;
    chunks.push({ from, to: Math.min(from + per, n) });
  }
  return chunks;
}

// 一个工人：自己的 Chrome、自己的调试端口、自己的一段帧区间。
async function worker(k, chunk) {
  const myPort = port + k;
  const profDir = path.resolve(outDir) + `/prof-film-${Date.now()}-${k}`;
  const chrome = spawn(chromePath, [
    "--headless=new",
    `--remote-debugging-port=${myPort}`,
    `--window-size=${W},${H}`,
    "--hide-scrollbars",
    "--no-first-run",
    "--no-sandbox",
    "--disable-dev-shm-usage",
    "--disable-extensions",
    "--user-data-dir=" + profDir,
    "about:blank",
    // No GPU flags: the lab renders canvas-2D, and angle/swiftshader makes
    // Chrome SIGTRAP on macOS sandboxed hosts.
  ], { stdio: "ignore" });
  chrome.on("exit", (code, sig) => {
    if (code !== 0 && code !== null) console.error(`worker ${k}: chrome exited code=${code} sig=${sig}`);
  });

  const jp = async (p) => (await fetch(`http://127.0.0.1:${myPort}${p}`)).json();

  try {
    let target = null;
    for (let i = 0; i < 80; i++) {
      try {
        const list = await jp("/json/list");
        target = list.find((t) => t.type === "page");
        if (target) break;
      } catch { /* browser not up yet */ }
      await sleep(250);
    }
    if (!target) throw new Error(`worker ${k}: no CDP page target`);

    const ws = new WebSocket(target.webSocketDebuggerUrl);
    const pending = new Map();
    // Fail in-flight sends when the browser dies, instead of hanging forever.
    ws.addEventListener("close", () => {
      for (const ok of pending.values()) ok({});
      pending.clear();
    });
    await new Promise((ok, bad) => {
      ws.onopen = ok;
      ws.onerror = () => bad(new Error(`worker ${k}: ws failed`));
    });
    let id = 0;
    ws.onmessage = (e) => {
      const m = JSON.parse(e.data);
      if (m.id && pending.has(m.id)) {
        pending.get(m.id)(m.result);
        pending.delete(m.id);
      }
    };
    const send = (method, params = {}) =>
      new Promise((ok) => {
        const n = ++id;
        pending.set(n, ok);
        ws.send(JSON.stringify({ id: n, method, params }));
      });
    const evaluate = async (expression, awaitPromise = false) => {
      const r = await send("Runtime.evaluate", { expression, returnByValue: true, awaitPromise });
      return r.result ? r.result.value : undefined;
    };

    await send("Page.enable");
    await send("Runtime.enable");
    // --window-size 给的是窗口尺寸，去掉浏览器自身的边框后视口会短一截
    // （1920x1080 的窗口，视口只有 1920x937）。用 CDP 把布局视口钉死。
    await send("Emulation.setDeviceMetricsOverride", {
      width: W, height: H, deviceScaleFactor: 1, mobile: false,
    });
    // 在页面自己的脚本跑之前就把 __filmT 设成 0：App 一启动就知道
    // 「有人在外部推时间」，于是不播声音、不显示「开始」按钮。
    await send("Page.addScriptToEvaluateOnNewDocument", { source: "window.__filmT = 0;" });
    await send("Page.navigate", { url });

    // 应用起来了没有：默认等 canvas（画布分镜），READY 可以换成任何选择器
    // ——节点分镜没有 canvas，等它自己的页签或舞台更实在。
    const readySel = process.env.READY || "canvas";
    let up = false;
    for (let i = 0; i < 120; i++) {
      if (await evaluate(`!!document.querySelector(${JSON.stringify(readySel)})`)) { up = true; break; }
      await sleep(500);
    }
    if (!up) throw new Error(`worker ${k}: ${readySel} never appeared — is the app running?`);
    await sleep(1500);

    const raf2 = () =>
      evaluate(`new Promise(r => requestAnimationFrame(() => requestAnimationFrame(r)))`, true);
    const rafN = (n) =>
      evaluate(
        `new Promise(r => { let i = 0; const f = () => (++i >= ${n} ? r(1) : requestAnimationFrame(f)); requestAnimationFrame(f); })`,
        true,
      );

    // PRESS=文本：录整个页面里的某一块（一个页签、一个按钮）之前先点进去。
    // 点完等页面稳定，再开始推时钟。
    if (process.env.PRESS) {
      const v = await evaluate(`(() => {
        const T = ${JSON.stringify(process.env.PRESS)};
        const el = [...document.querySelectorAll('*')].filter(e => e.textContent && e.textContent.trim() === T).pop();
        if (!el) return "none";
        const b = el.getBoundingClientRect();
        return JSON.stringify({ x: Math.round(b.x + b.width / 2), y: Math.round(b.y + b.height / 2) });
      })()`);
      if (!v || v === "none") throw new Error(`worker ${k}: PRESS target not found: ${process.env.PRESS}`);
      const p = JSON.parse(v);
      await send("Input.dispatchMouseEvent", { type: "mouseMoved", x: p.x, y: p.y, buttons: 0 });
      await send("Input.dispatchMouseEvent", { type: "mousePressed", x: p.x, y: p.y, button: "left", buttons: 1, clickCount: 1 });
      await sleep(60);
      await send("Input.dispatchMouseEvent", { type: "mouseReleased", x: p.x, y: p.y, button: "left", buttons: 0, clickCount: 1 });
      // 点进去之后等它真的画出来：节点分镜等自己的 take 节点，别的页面等
      // 一个固定的稳定期。固定 sleep 会在慢机器上截到还没挂载的帧——那几帧
      // 是空的，而且只有看成品才发现。
      let ready = false;
      for (let i = 0; i < 40; i++) {
        if (await evaluate(`!!document.querySelector('[data-testid^="take-"]')`)) { ready = true; break; }
        await sleep(250);
      }
      await sleep(ready ? 500 : 1500);
    }

    // A 节点分镜（keyframe）多要一帧：seek 先落进时钟，平台的帧时钟再把
    // 动画值写进视图树。画布分镜在 Draw 里读时钟，两帧就够。
    const nodeTakes = await evaluate(`!!document.querySelector('[data-testid^="take-"]')`);
    const rafFrames = nodeTakes ? 3 : 2;
    const settleMS = Number(process.env.SETTLE || (nodeTakes ? 60 : 0));
    if (nodeTakes) console.log(`worker ${k}: node takes present — ${rafFrames} rAF + ${settleMS}ms per frame`);

    // HIDE=sel,sel 把页面家具藏掉（页签栏、重播按钮），CLIP=sel 只录那个元素：
    // 录一支 take，不是录一个网页。
    for (const sel of (process.env.HIDE || "").split(",").filter(Boolean)) {
      await evaluate(`(() => { for (const el of document.querySelectorAll(${JSON.stringify(sel)})) el.style.display = "none"; return true; })()`);
    }
    let clip = null;
    if (process.env.CLIP) {
      const v = await evaluate(`(() => {
        const el = document.querySelector(${JSON.stringify(process.env.CLIP)});
        if (!el) return "none";
        el.scrollIntoView({ block: "center" });
        const b = el.getBoundingClientRect();
        return JSON.stringify({ x: Math.round(b.x), y: Math.round(b.y), width: Math.round(b.width), height: Math.round(b.height), scale: 1 });
      })()`);
      if (!v || v === "none") throw new Error(`worker ${k}: CLIP selector matched nothing: ${process.env.CLIP}`);
      clip = JSON.parse(v);
      // scrollIntoView 可能还在滚：等一帧再量一次，位置稳定了才开始。
      await rafN(2);
      const again = JSON.parse(await evaluate(`(() => {
        const b = document.querySelector(${JSON.stringify(process.env.CLIP)}).getBoundingClientRect();
        return JSON.stringify({ x: Math.round(b.x), y: Math.round(b.y), width: Math.round(b.width), height: Math.round(b.height), scale: 1 });
      })()`));
      if (again.height > 0 && again.width > 0) clip = again;
      if (clip.width < 2 || clip.height < 2) {
        throw new Error(`worker ${k}: CLIP element measured ${clip.width}x${clip.height} — is it visible?`);
      }
      console.log(`worker ${k}: clip ${clip.width}x${clip.height} at ${clip.x},${clip.y}`);
    }

    const t0 = Date.now();
    let done = 0, skipped = 0;
    for (let i = chunk.from; i < chunk.to; i++) {
      const file = frameFile(i);
      // 续录：已有的帧直接跳过。
      if (existsSync(file) && statSync(file).size > 0) { skipped++; continue; }
      await evaluate(`window.__filmT = ${i / fps}`);
      await rafN(rafFrames);
      if (settleMS > 0) await sleep(settleMS);
      const shot = clip
        ? await send("Page.captureScreenshot", { format: "jpeg", quality: 93, clip })
        : await send("Page.captureScreenshot", { format: "jpeg", quality: 93 });
      writeFileSync(file, Buffer.from(shot.data, "base64"));
      done++;
      const idx = i - chunk.from;
      if (idx % 60 === 0) console.log(`worker ${k}: ${idx}/${chunk.to - chunk.from} (${Math.round((Date.now() - t0) / 1000)}s)`);
    }
    console.log(`worker ${k}: done ${done} frames, skipped ${skipped} (${Math.round((Date.now() - t0) / 1000)}s)`);
    ws.close();
  } finally {
    chrome.kill();
  }
}

async function encode() {
  if (process.env.ENCODE === "0") return;
  // 缺帧检查：ffmpeg 的序号输入遇到空洞会静默截断。
  const missing = [];
  for (let i = 0; i < N; i++) {
    if (!existsSync(frameFile(i))) missing.push(i);
  }
  if (missing.length) {
    console.error(`missing ${missing.length} frames, first: ${missing.slice(0, 5).join(", ")}`);
    process.exit(1);
  }
  const mp4 = path.join(outDir, "film.mp4");
  const enc = spawnSync("ffmpeg", [
    "-y", "-framerate", String(fps), "-i", path.join(outDir, "frames", "f%05d.jpg"),
    "-vf", "scale=trunc(iw/2)*2:trunc(ih/2)*2",
    "-c:v", "libx264", "-preset", "medium", "-crf", "18", "-pix_fmt", "yuv420p", mp4,
  ], { stdio: ["ignore", "ignore", "pipe"] });
  if (enc.status !== 0) {
    console.error("ffmpeg failed:", String(enc.stderr).split("\n").slice(-6).join("\n"));
    process.exit(1);
  }
  console.log("wrote", mp4);
}

async function main() {
  const chunks = splitFrames(N, WORKERS);
  console.log(`${N} frames -> ${chunks.map((c) => `[${c.from},${c.to})`).join(" ")}`);
  const t0 = Date.now();
  await Promise.all(chunks.map((c, k) => worker(k, c)));
  console.log(`all workers done in ${Math.round((Date.now() - t0) / 1000)}s`);
  await encode();
  process.exit(0);
}

main().catch((e) => {
  console.error(e);
  process.exit(1);
});
