#!/usr/bin/env python3
"""确定性回归：同一部片录两遍，逐帧哈希比对——全等才算过。

    CHROME=/path/to/chrome python3 tools/diff.py <url> <tmpdir> <durSec> [fps]

把 README 里「录制的片子、观看的画面、审片的单帧，三者逐像素一致」
从承诺变成测试：任何渲染层改动之后跑一遍，退出码 0 = 承诺还在。

实现：直接调 shots/film.js 录两遍（ENCODE=0 只出帧），然后按文件名
配对比对。同一浏览器同一版本，JPEG 编码是确定的；某帧不一致时
打印帧号，方便直接开 frames.py 抓出来看。
"""
import hashlib
import os
import pathlib
import subprocess
import sys

HERE = os.path.dirname(os.path.abspath(__file__))


def record(url, out, dur, fps):
    env = dict(os.environ, ENCODE="0")
    r = subprocess.run(["node", os.path.join(HERE, "..", "shots", "film.js"),
                        url, out, str(dur), str(fps)],
                       env=env)
    if r.returncode != 0:
        sys.exit(f"recording failed: {out}")


def frames(d):
    fd = pathlib.Path(d) / "frames"
    return sorted(fd.glob("f*.jpg"))


def main():
    if len(sys.argv) < 4:
        print(__doc__)
        sys.exit(2)
    url, tmp, dur = sys.argv[1], sys.argv[2], sys.argv[3]
    fps = sys.argv[4] if len(sys.argv) > 4 else "30"

    a, b = os.path.join(tmp, "run-a"), os.path.join(tmp, "run-b")
    os.makedirs(a, exist_ok=True)
    os.makedirs(b, exist_ok=True)
    print("recording run A ...")
    record(url, a, dur, fps)
    print("recording run B ...")
    record(url, b, dur, fps)

    fa, fb = frames(a), frames(b)
    if len(fa) != len(fb):
        sys.exit(f"frame count mismatch: {len(fa)} vs {len(fb)}")
    if not fa:
        sys.exit("no frames recorded")

    bad = []
    for x, y in zip(fa, fb):
        ha = hashlib.sha256(x.read_bytes()).hexdigest()
        hb = hashlib.sha256(y.read_bytes()).hexdigest()
        if ha != hb:
            bad.append(x.name)
    if bad:
        n = len(bad)
        print(f"NON-DETERMINISTIC: {n}/{len(fa)} frames differ")
        print("first differing:", *bad[:10], sep="\n  ")
        sys.exit(1)
    print(f"OK: {len(fa)} frames identical across two runs")


if __name__ == "__main__":
    main()
