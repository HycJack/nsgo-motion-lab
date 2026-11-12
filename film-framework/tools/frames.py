#!/usr/bin/env python3
"""按时间点给 NSGo 影片抓帧 —— 用来审片。

片子的时钟是 filmkit.Now()，web 端认 JS 全局 __filmT。
所以：设 __filmT → 等两帧 → 截图，得到的是那一瞬间该有的画面。
"""
import base64, os, sys, pathlib
from playwright.sync_api import sync_playwright

URL = sys.argv[1] if len(sys.argv) > 1 else "http://127.0.0.1:8099/"
OUT = pathlib.Path(sys.argv[2] if len(sys.argv) > 2 else "./frames")
TIMES = [float(x) for x in (sys.argv[3].split(",") if len(sys.argv) > 3
                            else "3,8,11,15,19,22,26,30,35,42,46,51,56,60,64".split(","))]
W = int(os.environ.get("W", 1920))
H = int(os.environ.get("H", 1080))

OUT.mkdir(parents=True, exist_ok=True)
with sync_playwright() as p:
    b = p.chromium.launch(executable_path=os.environ["CHROME_BIN"],
                          args=["--no-sandbox", "--disable-gpu", "--disable-dev-shm-usage"])
    ctx = b.new_context(viewport={"width": W, "height": H})
    page = ctx.new_page()
    page.route("**/agent.minimax.cn/**", lambda r: r.abort())
    logs = []
    page.on("console", lambda m: logs.append(f"[{m.type}] {m.text}"))
    page.on("pageerror", lambda e: logs.append(f"[pageerror] {e}"))
    page.goto(URL, wait_until="domcontentloaded", timeout=60000)
    for _ in range(40):
        page.wait_for_timeout(500)
        if page.evaluate("() => !!document.querySelector('canvas')"):
            break
    page.wait_for_timeout(4000)

    cdp = ctx.new_cdp_session(page)
    for t in TIMES:
        page.evaluate(f"() => {{ window.__filmT = {t}; }}")
        # 两帧 rAF：让 app 的循环读到新时钟并重画
        page.evaluate("() => new Promise(r => requestAnimationFrame(() => requestAnimationFrame(r)))", )
        page.wait_for_timeout(350)
        page.evaluate("() => new Promise(r => requestAnimationFrame(() => requestAnimationFrame(r)))")
        page.wait_for_timeout(350)
        shot = cdp.send("Page.captureScreenshot", {"format": "jpeg", "quality": 88})
        f = OUT / f"t{t:06.2f}.jpg"
        f.write_bytes(base64.b64decode(shot["data"]))
        print("shot", f.name)
    if logs:
        print("CONSOLE:", *logs[:20], sep="\n  ")
    b.close()
