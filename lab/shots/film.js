// film.js - record an NSGo film scene frame by frame into an mp4.
//
//   node film.js <url> <outDir> <durSec> [fps]
//
// The scene must read filmkit.Now(): on web that clock honors the JS
// global __filmT, so the recorder steps time one frame at a time —
// deterministically, not in real time — screenshots each frame, and then
// (ENCODE=1, the default) hands the frame folder to ffmpeg.
//
//   SIZE=1280x720   viewport (default 1920x1080)
//   CHROME=/path    browser binary (default: first Chrome-like app found)
//   ENCODE=0        keep frames only, skip ffmpeg
//   FPS=60          frames per second (default 30; positional [fps] wins)
//
// Frames land in <outDir>/frames/f00000.jpg… and the mp4 in
// <outDir>/film.mp4.

const { spawn, spawnSync } = require("node:child_process");
const { existsSync, mkdirSync, writeFileSync } = require("node:fs");
const { setTimeout: sleep } = require("node:timers/promises");
const path = require("node:path");

const CHROME_CANDIDATES = [
  "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
  "/Applications/Chromium.app/Contents/MacOS/Chromium",
  "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome Canary",
  "/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
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
console.log("browser:", chromePath);
const profDir = path.resolve(outDir) + "/prof-film-" + Date.now();

const chrome = spawn(chromePath, [
  "--headless=new",
  `--remote-debugging-port=${port}`,
  `--window-size=${W},${H}`,
  "--hide-scrollbars",
  "--no-first-run",
  "--disable-extensions",
  "--user-data-dir=" + profDir,
  "about:blank",
  // No GPU flags: the lab renders canvas-2D, and angle/swiftshader makes
  // Chrome SIGTRAP on macOS sandboxed hosts.
], { stdio: "ignore" });
chrome.on("error", (e) => console.error("spawn error:", e.message));
chrome.on("exit", (code, sig) => console.log(`chrome exited: code=${code} sig=${sig}`));
// Pin the event loop: an in-flight await alone can let Node drain and exit 0.
const heartbeat = setInterval(() => {}, 1 << 30);

async function jp(p) {
  const res = await fetch(`http://127.0.0.1:${port}${p}`);
  return res.json();
}

async function main() {
  console.log("cdp: polling", port);
  let target = null;
  for (let i = 0; i < 80; i++) {
    try {
      const list = await jp("/json/list");
      target = list.find((t) => t.type === "page");
      if (target) break;
    } catch { /* browser not up yet */ }
    await sleep(250);
  }
  if (!target) throw new Error("no CDP page target");
  console.log("cdp: target ok");

  const ws = new WebSocket(target.webSocketDebuggerUrl);
  // Fail in-flight sends when the browser dies, instead of hanging forever.
  ws.addEventListener("close", () => {
    for (const ok of pending.values()) ok({});
    pending.clear();
  });
  await new Promise((ok, bad) => {
    ws.onopen = ok;
    ws.onerror = () => bad(new Error("ws failed"));
  });

  let id = 0;
  const pending = new Map();
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
  await send("Page.navigate", { url });

  // Wait for the wasm app's canvas, then settle.
  let up = false;
  for (let i = 0; i < 120; i++) {
    if (await evaluate(`!!document.querySelector('canvas')`)) { up = true; break; }
    await sleep(500);
  }
  if (!up) throw new Error("canvas never appeared — is the app running?");
  await sleep(1500);

  const raf2 = () =>
    evaluate(`new Promise(r => requestAnimationFrame(() => requestAnimationFrame(r)))`, true);

  const t0 = Date.now();
  for (let i = 0; i < N; i++) {
    await evaluate(`window.__filmT = ${i / fps}`);
    await raf2();
    const { data } = await send("Page.captureScreenshot", { format: "jpeg", quality: 93 });
    writeFileSync(path.join(outDir, "frames", `f${String(i).padStart(5, "0")}.jpg`), Buffer.from(data, "base64"));
    if (i % 60 === 0) console.log(i, "/", N, Math.round((Date.now() - t0) / 1000) + "s");
  }
  console.log("DONE", N, "frames in", Math.round((Date.now() - t0) / 1000) + "s");

  ws.close();
  chrome.kill();

  if (process.env.ENCODE !== "0") {
    const mp4 = path.join(outDir, "film.mp4");
    const enc = spawnSync("ffmpeg", [
      "-y", "-framerate", String(fps), "-i", path.join(outDir, "frames", "f%05d.jpg"),
      "-c:v", "libx264", "-preset", "medium", "-crf", "18", "-pix_fmt", "yuv420p", mp4,
    ], { stdio: ["ignore", "ignore", "pipe"] });
    if (enc.status !== 0) {
      console.error("ffmpeg failed:", String(enc.stderr).split("\n").slice(-6).join("\n"));
      process.exit(1);
    }
    console.log("wrote", mp4);
  }
  clearInterval(heartbeat);
  process.exit(0);
}

main().catch((e) => {
  console.error(e);
  chrome.kill();
  process.exit(1);
});
