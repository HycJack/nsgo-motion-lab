# Atlas Intro — NSGo canvas edition

The Atlas intro film, redrawn with **NSGo's own canvas engine** (the same
framework the Motion Lab uses) instead of a browser video framework.

Everything is one `core.Canvas` covering the window. `Root` drives a
`Signal` from the page clock once per animation frame and the canvas is
repainted as a **pure function of time**, so the six-scene film is fully
deterministic:

| scene | seconds | content |
|---|---|---|
| S1 | 0–9.5 | 拆开来，才学得透。 + tilted kinematics card (curve draws, ball rides) |
| S2 | 9.5–16 | flat page → turning 3D card |
| S3 | 16–25 | rotate / explode with 位置 y、速度 v、加速度 a call-outs |
| S4 | 25–33 | AP-style multiple-choice, wrong B ✕ / correct C ✓, mini plot |
| S5 | 33–42 | AI tutor chat, answer reveals line by line |
| S6 | 42–53 | 1,696 / 265 / 51 counters, subject chips, 从一节课开始。 |

## Run it

```bash
# from the repo root (nsgo-motion-lab)
../run-motionlab.sh            # runs the Motion Lab on :8099
# this experiment, web target on :8100:
cd lab && go tool nsgo dev -target web -open=false -port 8100 ./atlasintro
```

Open <http://127.0.0.1:8100/>. It plays in real time.

## Capture to video

`window.__atlasT` overrides the clock: set it to a number of seconds and
the canvas draws that exact frame, so a recorder can step frames without
dropping any. The pipeline used here:

1. Playwright (Chrome) loads the page at 1920×1080.
2. For `i` in `0 … 1589`: set `__atlasT = i/30`, wait two rAF ticks,
   screenshot to `f%05d.jpg`.
3. `ffmpeg -framerate 30 -i frames/f%05d.jpg -i narration.wav …` muxes the
   six Chinese voice-over clips (same script as the Remotion cut) at the
   scene starts 0.5 / 10 / 16.5 / 25.5 / 33.5 / 42.5 s.

Output: `atlas-video/AtlasIntro_nsgo.mp4` (1920×1080, 30 fps, 53 s).

## Files

- `main.go` — palette, math helpers, `Root`, all six scenes, the plot
  drawing helper.
- `clock_js.go` — film time on web (reads `window.__atlasT`, else the page
  clock); build tag `js && wasm`.
- `clock_other.go` — wall-clock time for native builds.
- `nsgo.yaml` — app id and web target.
