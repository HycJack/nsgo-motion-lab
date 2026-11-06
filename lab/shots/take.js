// take.js - capture a one-take card at several moments of its pass.
//
//   node take.js <url> <outDir> <tabLabel> <cardTitle> <durationMs> <atSec,..>
//
// The point is continuity: frames at 10%, 30%, 55%, 80%, 100% of the pass have
// to read as one move, and the last frame has to be a finished picture.

const { spawn } = require("node:child_process");
const { writeFileSync } = require("node:fs");
const { setTimeout: sleep } = require("node:timers/promises");
const path = require("node:path");

const CHROME = "C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe";
const url = process.argv[2] || "http://localhost:8099/";
const outDir = process.argv[3] || ".";
const tabLabel = process.argv[4] || "One take";
const cardTitle = process.argv[5] || "One take";
const atSec = (process.argv[6] || "1,4,8,12,17.5").split(",").map(Number);
// FULL=1 captures the whole viewport instead of the card's stage, which is what
// a page of view-built demos needs.
const full = process.env.FULL === "1";
const tag = (process.argv[7] || cardTitle.split(" · ")[0]).replace(/[^\w.-]+/g, "-");
const port = Number(process.env.CDP_PORT || 9385);
const profDir = path.resolve(outDir) + "\\prof-take-" + Date.now();

const chrome = spawn(CHROME, [
  "--headless=new",
  `--remote-debugging-port=${port}`,
  "--window-size=1240,1000",
  "--hide-scrollbars",
  "--enable-unsafe-swiftshader",
  "--use-gl=angle",
  "--use-angle=swiftshader",
  "--no-first-run",
  "--disable-extensions",
  "--user-data-dir=" + profDir,
  "about:blank",
], { stdio: "ignore" });

async function jp(p) {
  const res = await fetch(`http://127.0.0.1:${port}${p}`);
  return res.json();
}

async function main() {
  let target = null;
  for (let i = 0; i < 80; i++) {
    try {
      const list = await jp("/json/list");
      target = list.find((t) => t.type === "page");
      if (target) break;
    } catch (e) {
      if (i === 0) console.log("(browser not up yet)");
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
  const evaluate = async (expression) => {
    const r = await send("Runtime.evaluate", { expression, returnByValue: true });
    return r.result ? r.result.value : undefined;
  };

  await send("Page.enable");
  await send("Runtime.enable");
  await send("Page.navigate", { url });
  await sleep(9000);

  async function press(exact) {
    const v = await evaluate(`(() => {
      const el = [...document.querySelectorAll('*')].filter(e => e.textContent && e.textContent.trim() === ${JSON.stringify(exact)}).pop();
      if (!el) return "none";
      const b = el.getBoundingClientRect();
      return JSON.stringify({ x: Math.round(b.x + b.width / 2), y: Math.round(b.y + b.height / 2) });
    })()`);
    if (!v || v === "none") return v;
    const { x, y } = JSON.parse(v);
    await send("Input.dispatchMouseEvent", { type: "mouseMoved", x, y, buttons: 0 });
    await send("Input.dispatchMouseEvent", { type: "mousePressed", x, y, button: "left", buttons: 1, clickCount: 1 });
    await sleep(40);
    await send("Input.dispatchMouseEvent", { type: "mouseReleased", x, y, button: "left", buttons: 0, clickCount: 1 });
    return exact;
  }

  // The stage's own box: the shader canvas inside the card, in viewport
  // coordinates. Capturing only this keeps a frame small enough that the
  // screenshot cost cannot eat the schedule.
  const stageBox = async () => {
    const v = await evaluate(`(() => {
      const T = ${JSON.stringify(cardTitle)};
      const cands = [...document.querySelectorAll('*')].filter(e => e.textContent && e.textContent.includes(T));
      if (!cands.length) return "none";
      const t = cands.find(e => ![...e.children].some(c => c.textContent && c.textContent.includes(T))) || cands[cands.length - 1];
      let el = t;
      while (el.parentElement && (el.getBoundingClientRect().height < 200 || el.getBoundingClientRect().width < 300)) el = el.parentElement;
      el.scrollIntoView({ block: "start" });
      const canvas = el.querySelector("canvas");
      const b = (canvas || el).getBoundingClientRect();
      return JSON.stringify({ x: Math.round(b.x), y: Math.round(b.y), w: Math.round(b.width), h: Math.round(b.height) });
    })()`);
    if (!v || v === "none") return null;
    return JSON.parse(v);
  };

  console.log(await press(tabLabel));
  await sleep(2500);
  const box = await stageBox();
  if (!box) throw new Error(`no card ${cardTitle}`);
  await sleep(600);
  const box2 = await stageBox();
  console.log(`stage at ${box2.x},${box2.y} ${box2.w}x${box2.h}`);

  console.log(await press(`Replay ${cardTitle.split(" · ")[0]}`));
  const t0 = Date.now();
  for (const at of atSec) {
    const wait = Math.max(0, at * 1000 - (Date.now() - t0));
    await sleep(wait);
    const clip = full
      ? { x: 0, y: 0, width: 1240, height: 900, scale: 1 }
      : { x: Math.max(0, box2.x), y: Math.max(0, box2.y), width: box2.w, height: box2.h, scale: 1 };
    const { data } = await send("Page.captureScreenshot", { format: "png", clip });
    const name = `take-${tag}-${String(at).replace(".", "_")}s.png`;
    writeFileSync(path.join(outDir, name), Buffer.from(data, "base64"));
    console.log(`wrote ${name}  (+${((Date.now() - t0) / 1000).toFixed(1)}s)`);
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
