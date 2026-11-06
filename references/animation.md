# Motion in NSGo

Everything below was checked against the module (`go doc`), then built into a
231-card app and run on web and Windows. Where a claim comes from a bug we
actually hit, the card count is in the sentence.

## The two clocks

| Clock | Who steps it | Cost | What lives there |
|---|---|---|---|
| **Commit** | Go, on a signal or a `Show`/`Each`/`Dynamic` rebuild | full re-render of that span | which nodes exist, their layout, static style |
| **Frame** | the platform's animator, at display rate | no commit | `AnimatedStyle` fields, shader uniforms bound to `AnimatedNode` |
| **GPU** | the fragment shader | no commit, per pixel | anything a Kage program computes from uniforms |

**The point of the whole API is to keep motion on the frame clock.** A card that
re-renders 60 times a second works, and is wrong: it rebuilds its subtree, throws
away its `NewAnimatedValue`s and restarts every animation from zero — the
symptom is "it jumps back to the start when I touch it".

## Five ways to move something

| Need | Use | Lands on |
|---|---|---|
| A state change (hover, press, selection, an element that changes size) | `Style.Transition` | frame, no graph at all |
| A standard entrance/exit/attention/loop | `core.Motion` + `core.MotionPreset` | frame |
| A custom choreography (overlap, holds, swaps) | `AnimatedStyle` + `AnimatedNode` graph | frame |
| Per-pixel or infinite-field motion | `core.Shader` uniforms | GPU |
| Procedural drawing (charts, projections) | `core.Canvas` `Draw` + `Bind` | frame (re-record) |

## The five rules (each cost us a broken card)

1. **Never read a signal or an animated value during render.** Put it in a
   closure — `StyleFunc`, `ContentFunc`, `Show.When`, `Drawing.Draw` — or in
   `AnimatedStyle`. Reading in render re-renders the subtree on every change.
2. **Prefer restyling over rebuilding.** One `StyleFunc` that reads a signal is
   a style-only commit; a `Dynamic` around the same thing rebuilds and restarts
   the animations inside.
3. **The last frame must be the readable state.** A finite `Timing` *holds* its
   final value, so a card whose loop returns to the beginning looks blank once
   it settles: three of ours did (`pixelate`, `exploded`, `montage`). Rewrite the
   stops so `p == 1` is the interesting frame, or loop for real.
4. **An animated element needs a measured parent.** A `View` whose children are
   all absolute (and which has no size of its own) measures 0×0 and clips them
   — everything "does not animate", nothing is visible. Give the stage
   `Position: Absolute, Top/Left/Right/Bottom: 0` (or a real height). 57 of 192
   cards were blank from this.
5. **`AnimatedStyle` fields are nodes, not numbers.** `Scale: 1` does not
   compile; `core.Constant(0.8)` / `lerp(p, 24, 0)` do. Static colour →
   `core.Interpolate(core.Constant(0), core.Interpolation{Colors: []core.Color{c, c}})`
   or just `Style.Background`.

### Easing menu that reads as motion

| Situation | Easing |
|---|---|
| entrance | `core.BezierEasing(core.EaseOut)` (fast out, settle) — or steeper: `core.BezierEasing(core.Curve{X1: 0.16, Y1: 1, X2: 0.3, Y2: 1})` |
| exit | `core.BezierEasing(core.EaseIn)` — leaving should accelerate away |
| overshoot ("pop") | a curve whose y passes 1: `core.BezierEasing(core.Curve{X1: 0.34, Y1: 1.56, X2: 0.64, Y2: 1})` — or `core.Easing{Kind: core.EasingBack, Mode: core.EaseModeOut}` |
| continuous loop | `core.BezierEasing(core.EaseInOut)`, or `core.BezierEasing(core.Linear)` for spins |
| camera / mechanical | `core.BezierEasing(core.EaseInOut)` on *every* lane, identical |
| stepped / typed | `core.StepsEasing(n, core.JumpEnd)` |

`core.EaseOut`, `core.EaseIn`, `core.EaseInOut` and `core.Linear` are `core.Curve`s
(CSS bezier control points); `TimingConfig.Easing` and `Interpolation.Easing` want
a `core.Easing`, so wrap them — `core.BezierEasing(curve)` — while
`Style.Transition.Easing` takes the bare `Curve`.

Distances and durations that worked for UI-size cards: **8–32 pt** for controls,
**24–64 pt** for cards and panels; **150–300 ms** for micro-interactions,
**300–600 ms** for entrances, **600–900 ms** for attention and loops. Anything
past 900 ms in a live UI reads as lag unless it is a deliberate loop.

## Recipe 1 — a state change animates itself

```go
Style: core.Style{
    Background: core.RGB(0x16202c),
    Transition: core.Transition{
        Duration:   180 * time.Millisecond,
        Easing:     core.EaseOut,                   // a Curve, not an Easing
        Properties: core.TransitionBackground | core.TransitionTransform,
    },
}
```

- No graph and no `Motion`: when the style value changes (a signal flips the
  colour, a `StyleFunc` returns a new one), the platform interpolates it.
- `Properties` is a bitmask — `TransitionOpacity`, `TransitionBackground`,
  `TransitionTransform`, `TransitionBorder`, `TransitionShadow`,
  `TransitionLayout` — and **0 means every property except `TransitionLayout`**,
  which must be asked for by name.
- `TransitionLayout` is the FLIP trick: the view moves and resizes from where it
  showed, while its children take their new frames inside it at once.
- Tailwind does the same in a class: `"transition-colors duration-200 ease-out"`
  (`transition` = all five visual properties, `transition-all` adds layout,
  plus `transition-opacity`, `transition-transform`, `transition-shadow`,
  `delay-*`).
- Transitions are skipped under reduce-motion. `Easing` is a `core.Curve`
  (`core.EaseOut`), while `TimingConfig` and `Interpolation` take a `core.Easing`:
  `core.BezierEasing(curve)` converts.
- `View.OnTransitionEnd func()` runs once the restyle's `Delay + Duration` is
  over (`element.go:88`) — the hook for "chained" motion without timers, and it
  fires at once when there is no duration.

## Recipe 2 — a preset, six lines

```go
var m *core.Motion
spec := core.FadeInUp.Spec()          // or your own core.MotionSpec
m = core.NewMotion(cx, &spec)         // builds values and plays at once

return core.View{
    Style:    core.Style{Animated: m.Style()}, // ← the only binding
    Children: []core.Node{card()},
}
```

- `NewMotion(cx, nil)` gives every channel at rest (opacity 1, transforms 0) —
  use it when the motion is triggered later: `m.Play(spec, nil)`.
- `m.Play` and the presets **honour the system reduce-motion setting**
  (`cx.Env().ReduceMotion()`): each channel jumps to where it ends. Hand-built
  graphs do not, so check `cx.Env().ReduceMotion()` yourself if you care.
- One `Motion` per view. Its channels are the 9 in `core.MotionChannel`
  (opacity, translate x/y, rotate, rotate x/y, scale, scale x/y).

### What the 38 presets already are

| Group | Presets (duration in ms) |
|---|---|
| Entrances | `FadeIn` 300 · `FadeInUp/Down/Left/Right` 450 (+24 pt translate) · `ZoomIn` 350 (scale .9) · `PopIn` 450 (scale .5, back-eased) · `BounceIn` 650 (5-stop scale) · `SlideInUp/Down` 400 (64 pt) · `FlipInX/Y` 600 (90→−20→10→−5→0°) |
| Exits | `FadeOut` 250 · `FadeOutUp/Down/Left/Right` 300 · `ZoomOut` 250 · `PopOut` 300 · `SlideOutDown/Up` 300 |
| Attention | `Shake` 500 · `HeadShake` 800 · `Wiggle` 600 · `Swing` 900 · `Pulse` 500 · `Heartbeat` 900 · `RubberBand` 800 · `Jello` 800 · `Tada` 900 · `Jump` 700 · `Flash` 800 |
| Loops (`Forever`) | `Spin` 1000 linear · `Float` 1500 alternate · `Breathe` 1800 alternate · `Sway` 1600 alternate · `Bob` 1000 · `Glow` 2000 |

`p.Name()` is camelCase (`fadeIn`), `p.Kind()` is `MotionEnter/Exit/Attention/Loop`,
`p.Spec()` is the data — copy it and edit stops rather than writing one by hand.

## Recipe 3 — your own choreography as data

```go
spec := core.MotionSpec{
    Timing: core.KeyframeTiming{
        Duration: 700 * time.Millisecond,
        Easing:   core.BezierEasing(core.EaseOut),
    },
    Tracks: []core.MotionTrack{
        {Channel: core.MotionOpacity, Stops: core.Stops(0, 1)},
        // 18 pt down → overshoot 3 pt up → rest: no spring needed
        {Channel: core.MotionTranslateY, Stops: []core.KeyframeStop{
            core.StopAt(0, 18), core.StopAt(0.6, -3), core.StopAt(1, 0),
        }},
    },
}
```

- `core.Stops(...)` spreads evenly (first at 0, last at 1); `core.StopAt(offset, v)`
  places one; `stop.Eased(e)` curves only the segment that starts at it.
- `KeyframeTiming.Direction`: `PlayAlternate` turns a one-way track into a
  there-and-back loop (that is how `Breathe`/`Float`/`Sway` are built);
  `Iterations: core.Forever` repeats. `KeyframeTiming.Duration` defaults to 300 ms.
- Drive a plain value the same way: `core.ValueKeyframes(v, timing, stops...)`.

## Recipe 4 — the value graph

```go
ms := time.Millisecond
p := core.NewAnimatedValue(cx, 0)                       // the driver, 0 → 1
cx.StartAnimation(core.Timing(p, 1, core.TimingConfig{  // nil done = fire and forget
    Duration: 600 * ms, Easing: core.BezierEasing(core.EaseOut),
}), nil)

core.View{Style: core.Style{Animated: &core.AnimatedStyle{
    Opacity:    core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.2}, Output: []float64{0, 1}}),
    TranslateY: lerp(p, 24, 0),
    Scale:      core.Constant(1),                       // a static number is a node
}}}
```

`lerp(p, a, b)` above is the app-side one-liner this whole reference assumes:
`core.Interpolate(p, core.Interpolation{Input: []float64{0, 1}, Output: []float64{a, b}})`.
Keep one such helper next to the palette; every card reads better for it.

| Combinator | Shape |
|---|---|
| `core.Timing(v, to, TimingConfig{Duration, Easing, Delay})` | one value → one target (500 ms when 0) |
| `core.Spring(v, to, SpringConfig{Stiffness, Damping, Mass, Velocity})` | 100 / 10 / 1 when 0 |
| `core.Decay(v, DecayConfig{Velocity, Deceleration, Clamp, Min, Max, Bounce})` | a fling |
| `core.Sequence` / `core.Parallel` / `core.Stagger(interval, ...)` / `core.Delay(d)` | compose |
| `core.Loop(a, n)` | reset the values each iteration; **`n = 0` is forever** |
| `core.LoopOn(a, n)` | keep going from where the last iteration ended |
| `core.Interpolate(in, r)` | piecewise map: `Input`, `Output`, `Colors`, `Easing`, `Extrapolate` |
| `core.Add/Sub/Mul/Div/Mod` | arithmetic on nodes |
| `core.DiffClamp(in, min, max)` | follow a scroll offset inside a range (header hide/show) |
| `core.SnappingDiffClamp` | same, settling at the nearer bound in 200 ms |

`Interpolate` is the workhorse — most of the 231 cards are stops, not code:

| Want | Stops |
|---|---|
| fade in over the first 20 % | `In{0,0.2} Out{0,1}` |
| hold, then leave | `In{0,0.6,0.8} Out{1,1,0}` |
| swap two states at a moment | `In{0,0.48,0.52,1} Out{0,0,1,1}` |
| stagger a second element | same curve, `In{0.15,0.45,0.65,1}` |
| undo back to the start | `In{0,0.5,1} Out{0,1,0}` |

## Recipe 5 — shader uniforms (zero commits)

```go
var plasmaProgram = core.MustCompileShader(plasmaSource) // panics at start-up if wrong

core.Shader{
    Program: plasmaProgram,
    Style:   fill,
    Uniforms: core.Uniforms{
        "Time":  clock,                       // *core.AnimatedValue
        "Press": []any{x, y, age, weight},    // []any fills the vec4's lanes, in order
    },
}
```

- `Uniforms` values may be a number, an `AnimatedNode`, a `Color`, or `[]any` of
  those. Numbers and `AnimatedNode`s both work; nodes are stepped by the
  platform, so **animating a uniform costs no commit**.
- `var Time float` at package level of the shader source *is* the uniform, in
  declaration order.
- Available in Kage: `rectSize()`, `pixelScale()`, `image0At`…, the GLSL builtins,
  `for`/`if`, own functions, `mat2..mat4`, arrays. Entry point:
  `func Fragment(dst vec4, pos vec2) vec4`, **premultiplied alpha** in the result.
- A shader has `OnPress` but **no** `OnMouseMove` — wrap it in a `core.View` for hover.
- **Compile every source in a test**: `core.CompileShader(src)` returns the
  error, `MustCompileShader` panics on start-up.

**Ray marching is the honest 3D path** (no meshes, no z-buffer, no
`preserve-3d` in NSGo). Camera = three uniforms:

```go
eye := []any{lerp(p, -1.7, 1.7), core.Constant(0.95), core.Constant(2.75)}
// in the shader:
fwd   := normalize(Look - Eye)
right := normalize(cross(fwd, vec3(0, 1, 0)))
lens  := tan(Fov * 0.5)          // NOT 1/tan: that turns wide-angle into a telephoto
dir   := normalize(right*(uv.x*lens) - up*(uv.y*lens) + fwd)
```

Three findings worth keeping: `Fov` is `tan(fov/2)`; `BackfaceHidden` hidden on
web hides *both* faces (cross-fade opacity instead); a rotating node turns around
its **own centre**, so a full-bleed child swings across the screen.

## Recipe 6 — canvas motion

`Canvas.Draw` re-runs **whenever a signal it read changed, or the size changed** —
that is a re-record, not a commit, but it is still work per frame.

```go
core.Canvas{Style: fill, Draw: func(d *core.Drawing, size core.Size) {
    d.SetFillColor(core.RGB(0x080d14))
    d.FillRect(0, 0, float32(size.Width), float32(size.Height))
    a := angle.Get()          // reading it here is what schedules the redraw
    ...
}}
```

- `d.Bind(core.ArgA1, node)` binds one argument of the last op to a node: the
  platform replays the recorded ops with the node's value — and again whenever it
  changes — **with no work on the engine**, so no re-record. A gauge sweep or a
  progress arc needs no `Draw` closure per frame.
- Procedural tick: keep the state in a `core.Signal` and re-arm
  `cx.AfterFunc(33*ms, tick)` from inside the tick.
- Projection must be normalised to the size (`unit := h * 0.4`), or the drawing
  collapses to one pixel — a real bug: `x/(eye-z)*f` stays in units.
- Painter's sort is manual: `sort.Slice` the quads by depth, then fill **opaque**
  faces or the back ones show through.

## Recipe 7 — one take: many beats, no cuts

A card answers "how does one move look good"; a one take answers "how do fifteen
beats share one clock". The trick is to stop thinking in sequences and think in
**one value plus windows**.

```go
// p: 0→1 over the whole take, LINEAR in time.
p := core.AnimatedValue{}                                     // started by core.Timing(p, 1, cfg)
cam := core.Interpolate(p, core.Interpolation{                // a polyline, not an orbit equation
    Input:  []float64{0, 0.14, 0.30, 0.45, 0.60, 0.75, 0.88, 1},
    Output: []float64{-2.60, -2.15, -1.55, -1.05, -0.45, 0.40, 1.55, 1.70},
    Easing: core.Easing{Bezier: core.Linear},
})
title := core.Interpolate(p, core.Interpolation{              // a window: 3%→28%, ease inside
    Input:  []float64{0, 0.03, 0.20, 0.28},
    Output: []float64{0, 1, 1, 0},
    Easing: core.Easing{Bezier: core.EaseInOut},
})
```

- **`p` must be linear in time.** Windows are written as "second 4 to second 9",
  so easing belongs inside each window; ease the master lane and every window
  boundary drifts, which shows up as captions firing early.
- **Nothing mounts, nothing unmounts.** Hidden means `Opacity 0`. An `if` that
  removes a node re-mounts it later, and a re-mounted node starts from zero —
  that discontinuity *is* a cut.
- **Every caption brings its own plate.** A fragment shader cannot know what is
  behind it; white text over an orange sphere is unreadable without a
  `rgba(0,0,0,0.5)` plate.
- **A camera path is a sampled polyline.** `AnimatedStyle` interpolates, it does
  not evaluate expressions, so `Eye`/`Look`/`Fov`/`Tint` are each a lane with the
  same stop times, and the camera eases through them.
- **Same lane, same ground.** Anything that must not drift against the scenery
  (a walker on a panning ridge) rides the *same* lane as that scenery.

```go
// camera as uniforms: the whole take costs zero commits
core.Shader{Program: rayProgram, Style: fill, Uniforms: core.Uniforms{
    "Eye":  []any{eyeX, eyeY, eyeZ},
    "Look": []any{lookX, lookY, lookZ},
    "Fov":  fov, "Tint": tint,
}}
```

Draw the take's schedule once as a chart (one row per lane, windows as bars) — if
the rows line up, there is no cut to make; if a row has to *end* before the next
begins, that gap is where you would have cut.

## Testing motion

```go
app := nsgotest.Start(t, Root, nsgotest.Size(1180, 860))
app.Press("id:tab-Animation")
app.Press("id:chip-Replay fadeIn")
app.Settle()                                  // run the frame clock to rest
v := app.Find("id:stage-fadeIn")
if v.Animated.Opacity != 1 { t.Fatal("did not settle") }
```

| Call | What it does |
|---|---|
| `Settle()` | waits until no commit came and the platform animated nothing for 100 ms (5 s cap) — the resting state |
| `Advance(d)` | moves the test clock (timers, `AfterFunc`); the platform's own animations still run in real time |
| `nsgotest.ManualFrames()` | animations move *only* as `Advance` moves the clock, so a mid-flight frame is reproducible |
| `WaitFor(sel, d)` / `WaitForGone` | for state that arrives asynchronously |
| `View.Animated` | the values the animator last sent (`Opacity`, `TranslateX/Y`, `Scale…`, `Rotate`, colours); `Set` is the bitmask of what is bound |
| `app.AnimatedValue(nodeID)` | one node's raw value; `app.Bound(sel)` gives the gesture/scroll nodes a view binds |
| `Screenshot(name)` | pixels only under `nsgo test -target <platform>`; headless it logs and does nothing |

Assert the **resting** state by default: frame timing differs per target, so a
"half open" assertion is flaky unless `ManualFrames` makes it deterministic. Give
every animation a stable `TestID` (`chip-Replay <id>`, `preset-<name>`) so a test
can restart it.

## Cost, in one table

| You write | Commits | Frame work |
|---|---|---|
| `Style.Transition` | 0 (the style value still changes) | the listed properties only; `TransitionLayout` moves the frame |
| `Motion.Play` / keyframes / `AnimatedStyle` | 0 after mount | transform + opacity only (composited natively) |
| Shader uniform on a node | 0 | GPU draws it |
| `Canvas` `Draw` reading a signal | 0 | re-record + replay every change |
| `Canvas` `d.Bind` | 0 | replay only |
| Reading a signal inside `Render` | one per change, for the whole subtree | everything restarts |

## Mistakes that look like animation bugs

| Symptom | Cause |
|---|---|
| "It jumps back to the start" | a signal read during render, or `NewAnimatedValue` recreated inside the render of a `Dynamic` |
| Nothing visible, no error | absolute-only parent measured 0×0 (see rule 4) |
| Page will not scroll | the column above `core.Scroll` needs `FlexGrow: 1, FlexShrink: 1` |
| Card is blank once it settles | the loop's last frame equals its first (rule 3) |
| Interpolated fill is solid white | alpha in `Interpolation.Colors`; use opaque colours, keep alpha in `Style` |
| `Style` has no `Rotate` / no `Children` / no `BorderWidth` | those live on `AnimatedStyle`, in `Children:` of the element, or on the border fields |
| A transition animates more than you asked | `Transition.Properties` of 0 means every visual property; set the bitmask (or use a Tailwind `transition-*` class) |
| `transition-all` drags the layout along | its `Properties` include `TransitionLayout`; use `transition-colors` / `transition-transform` |
| `AnimatedStyle` has no `Width`/`Height` | animate `ScaleX` on a fixed-width node instead |
| `AnimatedStyle` has no `Perspective` | put `Perspective` on the parent's `Style` |
| `Scale: 1` will not compile | it is an `AnimatedNode`: `core.Constant(1)` |
| `tree: unknown field Children in struct literal of type core.Style` | close `Style` before `Children:` — it is a field of the element, not the style |
| `core.Alpha(0.6)` everywhere | it is only for `Style.Opacity`; colours take `core.RGBA(hex, a)` |
| `JustifyContent: core.AlignCenter` | `AlignItems` and `JustifyContent` are different enums (`core.JustifyCenter`) |
| `MotionPreset.Name()` shows `fadeIn` | it is NSH's JS name, camelCase |
| Shader crashes on start-up, build was fine | Kage compiles at run time — test with `core.CompileShader` |
| Shader loops forever / GPU hangs on some device | no `break` in Kage loops: freeze with `if d < eps { hit = 1.0 }` and `t += max(d, eps)*(1.0-hit)` |
