const { spawn } = require("node:child_process");
const { setTimeout: sleep } = require("node:timers/promises");
(async () => {
  const p = spawn("C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe", [
    "--headless=new", "--remote-debugging-port=9341",
    "--window-size=1240,900", "--hide-scrollbars",
    "--enable-unsafe-swiftshader", "--use-gl=angle", "--use-angle=swiftshader",
    "--no-first-run", "--disable-extensions",
    "--user-data-dir=C:\\Users\\huangyicao\\Downloads\\LobeHubDemo\\shots\\probe" + Date.now(),
    "about:blank",
  ], { stdio: "inherit" });
  p.on("exit", (c) => console.log("chrome exit:", c));
  await sleep(5000);
  try {
    const l = await (await fetch("http://127.0.0.1:9341/json/list")).json();
    console.log("targets:", l.map(t => t.type).join(","));
  } catch (e) { console.log("cdp:", e.message); }
  p.kill(); process.exit(0);
})();
