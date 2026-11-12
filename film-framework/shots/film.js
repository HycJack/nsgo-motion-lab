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

    let up = false;
    for (let i = 0; i < 120; i++) {
      if (await evaluate(`!!document.querySelector('canvas')`)) { up = true; break; }
      await sleep(500);
    }
    if (!up) throw new Error(`worker ${k}: canvas never appeared — is the app running?`);
    await sleep(1500);

    const raf2 = () =>
      evaluate(`new Promise(r => requestAnimationFrame(() => requestAnimationFrame(r)))`, true);

    const t0 = Date.now();
    let done = 0, skipped = 0;
    for (let i = chunk.from; i < chunk.to; i++) {
      const file = frameFile(i);
      // 续录：已有的帧直接跳过。
      if (existsSync(file) && statSync(file).size > 0) { skipped++; continue; }
      await evaluate(`window.__filmT = ${i / fps}`);
      await raf2();
      const { data } = await send("Page.captureScreenshot", { format: "jpeg", quality: 93 });
      writeFileSync(file, Buffer.from(data, "base64"));
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
