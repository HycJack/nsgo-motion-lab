package main

import (
	"fmt"
	"math"
	"time"

	"github.com/oneo-me/nsgo/core"
)

// canvasPage collects the core.Canvas examples. Canvas.Draw records 2D
// operations on the engine and the platform plays them; the recording is
// made again when a signal Draw read, or the canvas's size, changes.
func canvasPage(cx *core.Scope) core.Node {
	return grid(
		canvasWave(cx),
		canvasGauge(cx),
		canvasBars(cx),
		canvasStar(),
		canvasClip(),
		canvasHitTest(cx),
	)
}

// canvasWave redraws every frame from a signal the frame loop sets: the
// plain way to animate a canvas.
func canvasWave(cx *core.Scope) core.Node {
	phase := core.NewSignal(cx, 0.0)
	ticker(cx, func(t float64) { phase.Set(t) })
	return card("Sine wave", "Draw runs again whenever a signal it read changes — here a frame loop.", colCyan,
		core.Canvas{
			Label:  "Sine waves",
			TestID: "canvas-wave",
			Style:  well(150),
			Draw: func(d *core.Drawing, size core.Size) {
				w, h := float32(size.Width), float32(size.Height)
				mid := h / 2
				ph := float32(phase.Get())

				d.SetLineWidth(1)
				d.SetStrokeColor(core.RGBA(0x93a7c4, 0.13))
				for i := 1; i < 4; i++ {
					y := h * float32(i) / 4
					d.BeginPath()
					d.MoveTo(0, y)
					d.LineTo(w, y)
					d.Stroke()
				}

				colours := []core.Color{colCyan, colViolet}
				for k := range colours {
					d.SetStrokeColor(colours[k])
					d.SetLineWidth(2.5)
					d.BeginPath()
					for x := float32(0); x <= w; x += 2 {
						u := x / w
						wave := math.Sin(float64(u*2*math.Pi*2) + float64(ph*1.6) + float64(k)*1.1)
						y := mid - float32(wave)*mid*0.62*float32(1-float64(u)*0.35)
						if x == 0 {
							d.MoveTo(x, y)
						} else {
							d.LineTo(x, y)
						}
					}
					d.Stroke()
				}
			},
		},
		note("Two commits a frame: a signal drives the recording."),
	)
}

// canvasGauge binds an arc's end angle to an animated value, so the
// platform redraws it on its own frame clock with no commit per frame.
func canvasGauge(cx *core.Scope) core.Node {
	progress := core.NewAnimatedValue(cx, 0)
	cx.StartAnimation(core.Loop(core.Timing(progress, 1, core.TimingConfig{Duration: 2400 * time.Millisecond}), 0), nil)
	sweep := core.Interpolate(progress, core.Interpolation{
		Input:  []float64{0, 1},
		Output: []float64{0, 2 * math.Pi},
	})
	return card("Gauge", "Drawing.Bind ties an arc's angle to an animated value: no engine work per frame.", colGreen,
		core.Canvas{
			Label:  "Gauge, bound to an animated value",
			TestID: "canvas-gauge",
			Style:  well(150),
			Draw: func(d *core.Drawing, size core.Size) {
				cxp, cyp := float32(size.Width)/2, float32(size.Height)/2
				r := float32(math.Min(float64(cxp), float64(cyp))) - 22

				d.SetLineWidth(11)
				d.SetLineCap(core.CapRound)
				d.SetStrokeColor(core.RGBA(0xffffff, 0.07))
				d.BeginPath()
				d.Arc(cxp, cyp, r, 0, 2*math.Pi, false)
				d.Stroke()

				// The arc records a1 as 0; Bind replaces it as the
				// platform draws, from the animated value's sweep.
				d.SetStrokeColor(colGreen)
				d.BeginPath()
				d.Arc(cxp, cyp, r, 0, 0, false)
				d.Bind(core.ArgA1, sweep)
				d.Stroke()

				d.SetFillColor(colText)
				d.SetTextAlign(core.TextAlignCenter)
				d.SetTextBaseline(core.BaselineMiddle)
				d.FillText("bind", cxp, cyp, core.Font{Size: 15})
			},
		},
		note("core.ArgA1 bound to Interpolate(progress) — redrawn by the platform."),
	)
}

// canvasBars fills rounded bars with a gradient and writes their labels
// with measured text.
func canvasBars(cx *core.Scope) core.Node {
	seed := core.NewSignal(cx, 1)
	values := func() []float64 {
		s := seed.Get()
		out := make([]float64, 7)
		for i := range out {
			// A cheap deterministic jitter so Shuffle is reproducible.
			out[i] = 0.18 + 0.82*math.Abs(math.Sin(float64(s*7+i*13)*12.9898))
		}
		return out
	}
	return card("Bars", "A linear gradient, rounded rectangles, and text measured by the platform.", colAmber,
		core.Canvas{
			Label:  "Bar chart",
			TestID: "canvas-bars",
			Style:  well(160),
			Draw: func(d *core.Drawing, size core.Size) {
				w, h := float32(size.Width), float32(size.Height)
				v := values()
				const gap = 10
				bw := (w - gap*(float32(len(v))+1)) / float32(len(v))

				d.SetFillGradient(core.CanvasLinear(0, 0, 0, h,
					core.GradientStop{Offset: 0, Color: colAmber},
					core.GradientStop{Offset: 1, Color: core.RGBA(0xffb454, 0.18)}))
				for i, value := range v {
					bh := float32(value) * (h - 34)
					d.RoundRect(gap+float32(i)*(bw+gap), h-18-bh, bw, bh, 5)
					d.Fill()
				}

				d.SetFillColor(colFaint)
				d.SetTextAlign(core.TextAlignCenter)
				d.SetTextBaseline(core.BaselineTop)
				for i, value := range v {
					l := fmt.Sprintf("%.0f", value*100)
					d.FillText(l, gap+float32(i)*(bw+gap)+bw/2, h-15, core.Font{Size: 10})
				}
			},
		},
		chip("Shuffle", colAmber, func() { seed.Update(func(n int) int { return n + 1 }) }),
	)
}

// starPath is SVG path data read once by the core's path parser.
var starPath = core.MustParsePath("M 50 2 L 61 38 L 99 38 L 68 60 L 79 96 L 50 74 L 21 96 L 32 60 L 1 38 L 39 38 Z")

// canvasStar fills an SVG path with the even-odd rule.
func canvasStar() core.Node {
	return card("SVG paths", "ParsePath reads SVG path data; FillPath fills it under the even-odd rule.", colPink,
		core.Canvas{
			Label:  "Star",
			TestID: "canvas-star",
			Style:  well(150),
			Draw: func(d *core.Drawing, size core.Size) {
				scale := float32(math.Min(float64(size.Width), float64(size.Height))) / 110
				d.Save()
				d.Translate(float32(size.Width)/2-50*scale, float32(size.Height)/2-50*scale)
				d.Scale(scale, scale)
				d.SetFillGradient(core.CanvasLinear(0, 0, 100, 100,
					core.GradientStop{Offset: 0, Color: colPink},
					core.GradientStop{Offset: 1, Color: colViolet}))
				d.FillPath(starPath, core.EvenOdd)
				d.Restore()
			},
		},
		note("Path data is in its own coordinates and takes the drawing's transform."),
	)
}

// canvasClip clips a stripe field to a circle, then restores.
func canvasClip() core.Node {
	return card("Clipping", "ClipPath restricts everything drawn after it until Restore.", colBlue,
		core.Canvas{
			Label:  "Stripes inside a circle",
			TestID: "canvas-clip",
			Style:  well(150),
			Draw: func(d *core.Drawing, size core.Size) {
				w, h := float32(size.Width), float32(size.Height)
				cxp, cyp := w/2, h/2
				r := float32(math.Min(float64(cxp), float64(cyp))) - 14

				d.Save()
				d.BeginPath()
				d.Arc(cxp, cyp, r, 0, 2*math.Pi, false)
				d.Clip()

				d.SetFillColor(colWell)
				d.FillRect(0, 0, w, h)
				for i := 0; i < 26; i++ {
					x := -h + float32(i)*22
					d.SetFillGradient(core.CanvasLinear(x, 0, x+h, h,
						core.GradientStop{Offset: 0, Color: colBlue},
						core.GradientStop{Offset: 1, Color: colCyan}))
					d.BeginPath()
					d.MoveTo(x, 0)
					d.LineTo(x+11, 0)
					d.LineTo(x+11+h, h)
					d.LineTo(x+h, h)
					d.ClosePath()
					d.Fill()
				}
				d.Restore()

				d.SetLineWidth(1.5)
				d.SetStrokeColor(core.RGBA(0x93a7c4, 0.35))
				d.BeginPath()
				d.Arc(cxp, cyp, r, 0, 2*math.Pi, false)
				d.Stroke()
			},
		},
		note("Save / Clip / draw / Restore, as the web's canvas has them."),
	)
}

// canvasHitTest answers presses from the shapes of the last recording:
// PathContains replays a path's geometry without the platform.
func canvasHitTest(cx *core.Scope) core.Node {
	picked := core.NewSignal(cx, -1)
	var last *core.Drawing
	// The circles are recorded each draw, so their paths follow the size.
	var circles []*core.Path
	return card("Hit testing", "Tap a circle: PathContains tests the recorded shapes.", colCyan,
		core.Canvas{
			Label:  "Circles to tap",
			TestID: "canvas-hit",
			Style:  well(150),
			OnPress: func(e core.PressEvent) {
				x, y := float32(e.Position.X), float32(e.Position.Y)
				for i, p := range circles {
					if p != nil && last != nil && last.PathContains(p, x, y, core.NonZero) {
						picked.Set(i)
						return
					}
				}
				picked.Set(-1)
			},
			Draw: func(d *core.Drawing, size core.Size) {
				last = d
				w, h := float32(size.Width), float32(size.Height)
				chosen := picked.Get()
				colours := []core.Color{colCyan, colViolet, colAmber, colPink}
				r := float32(math.Min(float64(w)/8, float64(h)/3.4))
				circles = circles[:0]
				for i := range colours {
					cxp := w * float32(i+1) / 5
					cyp := h / 2
					p := core.MustParsePath(fmt.Sprintf(
						"M %.2f %.2f a %.2f %.2f 0 1 0 %.2f 0 a %.2f %.2f 0 1 0 %.2f 0 Z",
						cxp-r, cyp, r, r, 2*r, r, r, -2*r))
					circles = append(circles, p)
					if i == chosen {
						d.SetFillColor(colours[i])
					} else {
						d.SetFillColor(core.RGBA(0xffffff, 0.08))
					}
					d.FillPath(p, core.NonZero)
					if i == chosen {
						d.SetLineWidth(2)
						d.SetStrokeColor(colText)
						d.StrokePath(p)
					}
				}
			},
		},
		row(6,
			core.Text{Content: "picked:", Style: core.Style{FontSize: 11, Color: colFaint}},
			noteFunc(func() string {
				if picked.Get() < 0 {
					return "none"
				}
				return fmt.Sprintf("%d", picked.Get()+1)
			}),
		),
	)
}
