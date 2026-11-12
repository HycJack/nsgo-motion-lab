#!/usr/bin/env python3
"""成片合成：录制的无声视频 + 主音轨 + SFX 表 → 最终成片。

    python3 tools/mux.py <video.mp4> <audio.m4a> <out.mp4> --sfx sfx.json

sfx.json 和 Film.Sfx 同一份事实源（film.SfxJSON 生成或照抄格式）：

    [
      {"t": 3.2, "file": "asset:ding.m4a", "gain": 0.8},
      ...
    ]

file 里的 "asset:" 前缀会被剥掉，按 resources/ 下的相对路径找。
主音轨时长就是成片时长；SFX 用 adelay 推到各自的 t、amix 混入。
"""
import argparse
import json
import os
import subprocess
import sys


def amix_cmd(video, audio, sfx, out):
    """拼 ffmpeg 命令：视频直拷，主音轨 + 各 SFX 混音。"""
    inputs = ["-i", video, "-i", audio]
    parts = []
    mixins = "[1:a]"
    for i, s in enumerate(sfx):
        path = s["file"]
        if path.startswith("asset:"):
            path = path[len("asset:"):]
        if not os.path.exists(path):
            print(f"warn: sfx file missing: {path}", file=sys.stderr)
            continue
        inputs += ["-i", path]
        delay = int(round(s["t"] * 1000))
        gain = s.get("gain", 1)
        parts.append(
            f"[{i + 2}:a]adelay={delay}:{delay},volume={gain}[s{i}]"
        )
        mixins += f"[s{i}]"
    if parts:
        n = len(parts) + 1
        fc = ";".join(parts) + f";{mixins}amix=inputs={n}:normalize=0[aout]"
        acodec = ["-filter_complex", fc, "-map", "0:v", "-map", "[aout]"]
    else:
        acodec = ["-map", "0:v", "-map", "1:a"]
    return ["ffmpeg", "-y", *inputs, *acodec, "-c:v", "copy",
            "-c:a", "aac", "-shortest", out]


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("video", help="录制的无声视频（shots/film.js 的产物）")
    ap.add_argument("audio", help="主音轨（TTS+配乐的 final_audio.m4a）")
    ap.add_argument("out", help="成片路径")
    ap.add_argument("--sfx", help="SFX 表 JSON（与 Film.Sfx 同源）")
    args = ap.parse_args()

    sfx = []
    if args.sfx:
        with open(args.sfx) as f:
            sfx = json.load(f)

    cmd = amix_cmd(args.video, args.audio, sfx, args.out)
    r = subprocess.run(cmd, capture_output=True, text=True)
    if r.returncode != 0:
        print(r.stderr[-2000:], file=sys.stderr)
        sys.exit(1)
    print(f"wrote {args.out} ({len(sfx)} sfx)")


if __name__ == "__main__":
    main()
