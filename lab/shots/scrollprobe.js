// scrollprobe.js — is the app scrollable, and does an indicator show?
const { spawn } = require("node:child_process");
const { writeFileSync, mkdirSync } = require("node:fs");
const { setTimeout: sleep } = require("node:timers/promises");

const CHROME = "C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe";
const url = process.argv[2] || "http://localhost:8099/";
const outDir = process.argv[3] || ".";
mkdirSync(outDir, { recursive: true });

const port = 9334;
const chrome = spawn(CHROME, [
  "--headless=new",
  `--remote-debugging-port=${port}`,
  "--window-size=1240,900",
  // No --hide-scrollbars this time: the browser's own bar must be visible.
  "--enable-unsafe-swiftshader",
  "--use-gl=angle",
  "--use-angle=swiftshader",
  "--no-first-run",
  "--disable-extensions",
  "--user-data-dir=" + outDir + "\\profile2",
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
  await new Promise((ok, bad) => { ws.onopen = ok; ws.onerror = () => bad(new Error("ws")); });
  let id = 0;
  const pending = new Map();
  ws.onmessage = (e) => {
    const m = JSON.parse(e.data);
    if (m.id && pending.has(m.id)) { pending.get(m.id)(m.result); pending.delete(m.id); }
  };
  const send = (method, params = {}) => new Promise((ok) => {
    const n = ++id; pending.set(n, ok);
    ws.send(JSON.stringify({ id: n, method, params }));
  });
  const evalJS = async (expr) => {
    const r = await send("Runtime.evaluate", { expression: expr, returnByValue: true });
    return r.result && r.result.value;
  };
  const shot = async (name) => {
    const { data } = await send("Page.captureScreenshot", { format: "png" });
    writeFileSync(`${outDir}\\${name}.png`, Buffer.from(data, "base64"));
    console.log(`wrote ${name}.png`);
  };

  await send("Page.enable");
  await send("Runtime.enable");
  await send("Page.navigate", { url });
  await sleep(9000);

  // Which elements actually overflow, and how are their bars configured?
  const report = await evalJS(`(() => {
    const out = [];
    for (const el of document.querySelectorAll('*')) {
      if (el.scrollHeight > el.clientHeight + 4 || el.scrollWidth > el.clientWidth + 4) {
        const cs = getComputedStyle(el);
        out.push({
          tag: el.tagName,
          cls: (el.className || '').toString().slice(0, 40),
          overflowY: cs.overflowY,
          scrollbarWidth: cs.scrollbarWidth,
          client: el.clientHeight,
          scrollHeight: el.scrollHeight,
          scrollTop: el.scrollTop,
        });
      }
    }
    return JSON.stringify({ found: out.length, inner: innerHeight, items: out.slice(0, 6) }, null, 1);
  })()`);
  console.log("overflow report:", report);

  // Scroll the overflow container and keep the frame right after.
  const moved = await evalJS(`(() => {
    const el = [...document.querySelectorAll('*')].find(e => e.scrollHeight > e.clientHeight + 4);
    if (!el) return 'no scroller';
    el.scrollTop = 420;
    return 'scrollTop=' + el.scrollTop + ' of ' + el.scrollHeight;
  })()`);
  console.log("wheel:", moved);
  await sleep(150);
  await shot("07-scrolled");

  ws.close();
  chrome.kill();
}

main().catch((e) => { console.error("failed:", e.message); chrome.kill(); process.exit(1); });
