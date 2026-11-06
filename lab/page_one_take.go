package main

// page_one_take answers a question the other pages only circle: can one
// continuous move be written as code, with no cut anywhere?
//
// A cut is a node that is swapped, mounted or reset. A take is one value that
// runs 0 -> 1 once while everything on screen reads it. That is exactly what
// an animated track is, so a take is:
//
//   - one master progress, p, played once, linear in time, over the duration;
//   - every participant bound to p through its own window, so a title, a
//     camera and a panel share a clock but not a schedule;
//   - nothing added or removed while p moves: an element that is not on
//     screen yet is at opacity 0, not absent.
//
// Linear matters: the take's p has to be time, or the lane map below would be
// describing a schedule nobody can read off the clock. Easing belongs inside a
// window, which is where window() puts it.
//
// The 3D card is a camera pass over the ray-marched scene; the 2D card is the
// same idea without a shader; the map card draws the channels themselves, which
// is the honest way to show that the take has no seam.

import (
	"math"

	"github.com/oneo-me/nsgo/core"
)

func pageOneTake(cx *core.Scope) core.Node {
	return grid(
		oneTake3D(cx),
		oneTakeMap(cx),
		oneTake2D(cx),
	)
}

// channel is one value over the master progress.
//
// Stops are (p, value) pairs and hit the value exactly at that p. A segment is
// straight unless an easing is given, which is why a camera path is written as
// several stops: a curve has to be sampled into a polyline before a channel can
// carry it. Outside the stops the value holds, and holding is what turns a stop
// list into a window: 0,0 ... 1,1 ... 0,0 is in, hold, out.
func channel(p *core.AnimatedValue, ease core.Easing, stops ...[2]float64) core.AnimatedNode {
	in := make([]float64, len(stops))
	out := make([]float64, len(stops))
	for i, s := range stops {
		in[i], out[i] = s[0], s[1]
	}
	return core.Interpolate(p, core.Interpolation{Input: in, Output: out, Easing: ease})
}

// track is a channel that moves at a constant rate between its stops: for the
// camera, which should not come to rest at every keyframe.
func track(p *core.AnimatedValue, stops ...[2]float64) core.AnimatedNode {
	return channel(p, core.Easing{Bezier: core.Linear}, stops...)
}

// window is a channel that eases in and out of every segment: for opacity and
// the small offsets that go with it.
func window(p *core.AnimatedValue, stops ...[2]float64) core.AnimatedNode {
	return channel(p, bezier(core.EaseInOut), stops...)
}

// oneTake3D is the take itself: eighteen seconds, one camera, one clock.
func oneTake3D(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "T1", Title: "One take · 3D", Accent: colAmber,
		Duration: 18000 * ms, Height: 240, Linear: true,
		Note: "十八秒一个镜头：机位从高处降下来、贴着地面穿过场景、再从右侧撤出。标题、字幕、面板、落版全挂在同一条 0→1 通道上，中间没有一个节点被换掉。",
		Stage: func(p *core.AnimatedValue) core.Node {
			// The camera: eight stops on one path. Eye falls, runs past the set
			// at head height, and steps back out; Look stays on the subject, so
			// the move reads as a single intent rather than four shots. It never
			// comes closer than about two units, so a whole subject stays in
			// frame at every stop — a take whose middle is a frame full of
			// surface is a take nobody finishes.
			eyeX := track(p, [2]float64{0, -2.60}, [2]float64{0.14, -2.15}, [2]float64{0.30, -1.55}, [2]float64{0.45, -1.05},
				[2]float64{0.60, -0.45}, [2]float64{0.75, 0.40}, [2]float64{0.88, 1.55}, [2]float64{1, 1.70})
			eyeY := track(p, [2]float64{0, 1.85}, [2]float64{0.14, 1.35}, [2]float64{0.30, 0.95}, [2]float64{0.45, 0.90},
				[2]float64{0.60, 0.82}, [2]float64{0.75, 0.86}, [2]float64{0.88, 0.95}, [2]float64{1, 1.00})
			eyeZ := track(p, [2]float64{0, 4.30}, [2]float64{0.14, 3.90}, [2]float64{0.30, 3.35}, [2]float64{0.45, 2.95},
				[2]float64{0.60, 2.45}, [2]float64{0.75, 2.30}, [2]float64{0.88, 2.10}, [2]float64{1, 2.55})
			lookX := track(p, [2]float64{0, -0.05}, [2]float64{0.30, 0.02}, [2]float64{0.60, 0.12},
				[2]float64{0.88, 0.15}, [2]float64{1, 0.05})
			lookY := track(p, [2]float64{0, 0.55}, [2]float64{0.45, 0.50}, [2]float64{0.75, 0.46}, [2]float64{1, 0.46})
			lookZ := track(p, [2]float64{0, 0.10}, [2]float64{0.30, -0.02}, [2]float64{0.60, -0.05},
				[2]float64{0.75, -0.10}, [2]float64{1, 0.18})
			// The field of view breathes: it narrows as the camera comes down,
			// then opens wide for the pass through the set, where a narrow lens
			// would fill the frame with surface instead of subjects.
			fov := track(p, [2]float64{0, 0.72}, [2]float64{0.30, 0.66}, [2]float64{0.60, 0.78},
				[2]float64{0.75, 0.70}, [2]float64{0.88, 0.66}, [2]float64{1, 0.60})
			tint := track(p, [2]float64{0, 0.05}, [2]float64{0.30, 0.30}, [2]float64{0.55, 0.70},
				[2]float64{0.70, 0.90}, [2]float64{0.85, 0.35}, [2]float64{1, 0.12})

			// The overlays. Each one is a window on the same clock: in, hold,
			// out — never a mount and never a removal.
			titleIn := window(p, [2]float64{0, 0}, [2]float64{0.03, 0}, [2]float64{0.10, 1},
				[2]float64{0.20, 1}, [2]float64{0.28, 0})
			titleY := window(p, [2]float64{0, 14}, [2]float64{0.10, 0}, [2]float64{0.20, 0}, [2]float64{0.28, -12})
			subIn := window(p, [2]float64{0, 0}, [2]float64{0.16, 0}, [2]float64{0.24, 1},
				[2]float64{0.40, 1}, [2]float64{0.48, 0})
			subY := window(p, [2]float64{0, 10}, [2]float64{0.24, 0}, [2]float64{0.40, 0}, [2]float64{0.48, -8})
			panelIn := window(p, [2]float64{0, 0}, [2]float64{0.46, 0}, [2]float64{0.56, 1},
				[2]float64{0.74, 1}, [2]float64{0.84, 0})
			panelX := window(p, [2]float64{0, 40}, [2]float64{0.56, 0}, [2]float64{0.74, 0}, [2]float64{0.84, 36})
			endIn := window(p, [2]float64{0, 0}, [2]float64{0.86, 0}, [2]float64{0.95, 1}, [2]float64{1, 1})
			endY := window(p, [2]float64{0, 14}, [2]float64{0.95, 0}, [2]float64{1, 0})

			// A plate under each caption: a fragment shader has no idea what is
			// behind it, so the text brings its own contrast.
			plate := func(col core.Color, radius float32) func() core.Style {
				return func() core.Style {
					return core.Style{Background: col, CornerRadius: core.Points(radius), Padding: core.Symmetric(5, 12)}
				}
			}

			return core.View{Style: core.Style{Position: core.Relative, FlexGrow: 1}, Children: []core.Node{
				core.Shader{
					Program: rayProgram,
					Label:   "One-take camera pass over the ray-marched scene",
					Style:   stageFill,
					Uniforms: core.Uniforms{
						"Eye":  []any{eyeX, eyeY, eyeZ},
						"Look": []any{lookX, lookY, lookZ},
						"Fov":  fov,
						"Tint": tint,
					},
				},
				// The slate: top centre, and the only thing on screen for the
				// first fifth of the take.
				core.View{Style: core.Style{
					Position: core.Absolute, Top: core.Pt(16), Left: core.Pt(0), Right: core.Pt(0),
					AlignItems: core.AlignCenter,
					Animated:   &core.AnimatedStyle{Opacity: titleIn, TranslateY: titleY},
				}, Children: []core.Node{
					core.View{StyleFunc: plate(core.RGBA(0x000000, 0.5), 9), Children: []core.Node{
						core.Text{Content: "ONE TAKE", Style: core.Style{
							FontSize: 22, FontWeight: core.FontWeightBold, Color: colText,
						}},
					}},
				}},
				core.View{Style: core.Style{
					Position: core.Absolute, Top: core.Pt(62), Left: core.Pt(0), Right: core.Pt(0),
					AlignItems: core.AlignCenter,
					Animated:   &core.AnimatedStyle{Opacity: subIn, TranslateY: subY},
				}, Children: []core.Node{
					core.View{StyleFunc: plate(core.RGBA(0x000000, 0.38), 7), Children: []core.Node{
						core.Text{Content: "十八秒，一个镜头，没有剪切", Style: core.Style{FontSize: 12, Color: colMuted}},
					}},
				}},
				core.View{Style: core.Style{
					Position: core.Absolute, Right: core.Pt(14), Bottom: core.Pt(14), Width: core.Pt(158),
					Background: core.RGB(0x0d131c), CornerRadius: 10, Padding: core.All(10), Gap: 4,
					BorderWidth: 1, BorderColor: core.RGBA(0x35d6de, 0.35),
					Animated: &core.AnimatedStyle{Opacity: panelIn, TranslateX: panelX},
				}, Children: []core.Node{
					core.Text{Content: "REPLAY", Style: core.Style{FontSize: 10, FontWeight: core.FontWeightBold, Color: colCyan}},
					core.Text{Content: "面板也是通道的一段：滑进来、钉住、滑出去。", Style: core.Style{FontSize: 10, Color: colMuted, LineHeight: 14}},
				}},
				// The landing: bottom right, on the empty floor the exit move
				// leaves there. It takes the panel's corner over, which is why
				// the panel leaves a beat earlier.
				core.View{Style: core.Style{
					Position: core.Absolute, Right: core.Pt(16), Bottom: core.Pt(14),
					Padding: core.Symmetric(10, 14), Gap: 3, Background: core.RGBA(0x000000, 0.55), CornerRadius: 11,
					Animated: &core.AnimatedStyle{Opacity: endIn, TranslateY: endY},
				}, Children: []core.Node{
					core.Text{Content: "MOTION LAB", Style: core.Style{
						FontSize: 20, FontWeight: core.FontWeightBold, Color: colText,
					}},
					core.Text{Content: "NSGo · 一条通道，零剪切", Style: core.Style{FontSize: 10, Color: colFaint}},
				}},
			}}
		}})
}

// oneTakeMap draws the take's channels on one axis. Everything that moves in
// the 3D card is a row here; the rows that overlap are the handoffs that would
// otherwise have been cuts.
func oneTakeMap(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "T2", Title: "Lane map", Accent: colViolet, Height: 240,
		Note: "上面那个镜头的时间表：横轴是同一条 p，每一行都是一个窗口。窗口能对上，就不需要剪切——换镜头才需要。",
		Stage: func(p *core.AnimatedValue) core.Node {
			return core.Canvas{
				Label: "Lane map of the one-take",
				Style: stageFill,
				Draw: func(d *core.Drawing, size core.Size) {
					w, h := float32(size.Width), float32(size.Height)
					d.SetFillColor(core.RGB(0x0b1119))
					d.FillRect(0, 0, w, h)

					left, right := float32(98), w-14
					span := right - left

					// The axis: p, which is the take's own clock.
					for i := 0; i <= 4; i++ {
						x := left + span*float32(i)/4
						d.SetStrokeColor(core.RGBA(0x93a7c4, 0.16))
						d.SetLineWidth(1)
						d.BeginPath()
						d.MoveTo(x, 22)
						d.LineTo(x, h-16)
						d.Stroke()
						d.SetFillColor(core.RGBA(0x93a7c4, 0.45))
						d.FillText([]string{"0s", "4.5s", "9s", "13.5s", "18s"}[i], x-7, 14, core.Font{Size: 8})
					}

					type row struct {
						label  string
						colour core.Color
						// A camera row carries its stops; a window its ends.
						values [][2]float64
						from   float64
						to     float64
					}
					rows := []row{
						{label: "eye.x", colour: colAmber, values: [][2]float64{{0, -2.60}, {0.14, -2.15}, {0.30, -1.55}, {0.45, -1.05}, {0.60, -0.45}, {0.75, 0.40}, {0.88, 1.55}, {1, 1.70}}},
						{label: "eye.y", colour: colAmber, values: [][2]float64{{0, 1.85}, {0.14, 1.35}, {0.30, 0.95}, {0.45, 0.90}, {0.60, 0.82}, {0.75, 0.86}, {0.88, 0.95}, {1, 1.00}}},
						{label: "eye.z", colour: colAmber, values: [][2]float64{{0, 4.30}, {0.14, 3.90}, {0.30, 3.35}, {0.45, 2.95}, {0.60, 2.45}, {0.75, 2.30}, {0.88, 2.10}, {1, 2.55}}},
						{label: "look.x", colour: colPink, values: [][2]float64{{0, -0.05}, {0.30, 0.02}, {0.60, 0.12}, {0.88, 0.15}, {1, 0.05}}},
						{label: "fov", colour: colPink, values: [][2]float64{{0, 0.72}, {0.30, 0.66}, {0.60, 0.78}, {0.75, 0.70}, {0.88, 0.66}, {1, 0.60}}},
						{label: "tint", colour: colViolet, values: [][2]float64{{0, 0.05}, {0.30, 0.30}, {0.55, 0.70}, {0.70, 0.90}, {0.85, 0.35}, {1, 0.12}}},
						{label: "title", colour: colCyan, from: 0.03, to: 0.28},
						{label: "subtitle", colour: colCyan, from: 0.16, to: 0.48},
						{label: "ui panel", colour: colGreen, from: 0.46, to: 0.84},
						{label: "landing", colour: colBlue, from: 0.86, to: 1.00},
					}

					top, bottom := float32(30), h-18
					step := (bottom - top) / float32(len(rows))
					for i, r := range rows {
						y := top + step*float32(i) + step/2
						d.SetFillColor(core.RGBA(0xa2b3c9, 0.72))
						d.FillText(r.label, 8, y+3, core.Font{Size: 9})

						if r.values != nil {
							lo, hi := r.values[0][1], r.values[0][1]
							for _, v := range r.values {
								lo, hi = math.Min(lo, v[1]), math.Max(hi, v[1])
							}
							if hi-lo < 1e-6 {
								hi = lo + 1
							}
							d.SetStrokeColor(r.colour)
							d.SetLineWidth(1.5)
							d.BeginPath()
							for k, v := range r.values {
								x := left + span*float32(v[0])
								yy := y - 8 + 16*float32((v[1]-lo)/(hi-lo))
								if k == 0 {
									d.MoveTo(x, yy)
								} else {
									d.LineTo(x, yy)
								}
							}
							d.Stroke()
							continue
						}

						// A window: the bar is the time on screen, the end caps
						// the fade in and out.
						x0 := left + span*float32(r.from)
						x1 := left + span*float32(r.to)
						d.SetFillColor(core.RGBA(0xffffff, 0.05))
						d.FillRect(left, y-6, span, 12)
						d.SetStrokeColor(r.colour)
						d.SetLineWidth(1.5)
						d.BeginPath()
						d.MoveTo(x0, y)
						d.LineTo(x1, y)
						d.Stroke()
						d.SetFillColor(r.colour)
						d.FillRect(x0, y-5, 1.6, 10)
						d.FillRect(x1-1.6, y-5, 1.6, 10)
					}

					d.SetFillColor(core.RGBA(0x93a7c4, 0.42))
					d.FillText("every row reads the same p · nothing is mounted or removed", 8, h-5, core.Font{Size: 8})
				},
			}
		}})
}

// oneTake2D is the same construction without a shader: a landscape in bands,
// each band panning at its own rate under one clock, and the title landing on
// the last third. Nothing is cut; the bands only differ in rate.
func oneTake2D(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "T3", Title: "One take · 2D", Accent: colCyan,
		Duration: 12000 * ms, Height: 200, Linear: true,
		Note: "不用着色器也是同一个思路：三层风景按不同速度横移，一条通道走完，落版在最后三分之一。切的是各层速度，不是镜头。",
		Stage: func(p *core.AnimatedValue) core.Node {
			// Three rates on one clock. The far band crawls, the near band runs,
			// and the walker rides the near band: same channel, so the figure
			// never drifts against the ground it stands on.
			far := track(p, [2]float64{0, 34}, [2]float64{1, -34})
			mid := track(p, [2]float64{0, 120}, [2]float64{1, -120})
			lift := window(p, [2]float64{0, 0.94}, [2]float64{0.20, 0})
			endIn := window(p, [2]float64{0, 0}, [2]float64{0.62, 0}, [2]float64{0.74, 1}, [2]float64{1, 1})
			endY := window(p, [2]float64{0, 14}, [2]float64{0.74, 0}, [2]float64{1, 0})

			// A band overhangs the stage by more than any rate in the take can
			// travel, so its edge can never come into frame. Its crest is a row
			// of discs in the same colour, which is enough to read as a ridge
			// instead of a rectangle — and they ride the band, so the silhouette
			// moves with the ground it belongs to.
			band := func(h, bottom float32, col core.Color, x core.AnimatedNode, humps ...[2]float32) core.Node {
				kids := make([]core.Node, 0, len(humps))
				for _, hp := range humps {
					kids = append(kids, core.View{Style: core.Style{
						Position: core.Absolute, Left: core.Pt(core.Points(hp[0])), Top: core.Pt(core.Points(-hp[1] / 2)),
						Width: core.Pt(core.Points(hp[1])), Height: core.Pt(core.Points(hp[1])),
						CornerRadius: core.Points(hp[1] / 2), Background: col,
					}})
				}
				return core.View{Style: core.Style{
					Position: core.Absolute, Left: core.Pt(-160), Right: core.Pt(-160), Bottom: core.Pt(core.Points(bottom)),
					Height:     core.Pt(core.Points(h)),
					Background: col, Animated: &core.AnimatedStyle{TranslateX: x},
				}, Children: kids}
			}

			return core.View{Style: core.Style{Position: core.Relative, FlexGrow: 1}, Children: []core.Node{
				Fill(core.RGB(0x0d1726)),
				// A moon that does not move: the one still thing, so the pan has
				// something to be measured against.
				core.View{Style: core.Style{
					Position: core.Absolute, Left: core.Pt(178), Top: core.Pt(16),
					Width: core.Pt(22), Height: core.Pt(22), CornerRadius: 11,
					Background: core.RGBA(0x93a7c4, 0.16), BorderWidth: 1, BorderColor: core.RGBA(0x93a7c4, 0.32),
				}},
				// Far ridge: lit by the sky it sits under, so lighter than it.
				band(46, 88, core.RGB(0x1c2e44), far,
					[2]float32{40, 70}, [2]float32{210, 84}, [2]float32{400, 60}, [2]float32{560, 80}, [2]float32{760, 72}),
				// Near ridge: four times the rate, and darker than the sky, which
				// is what makes it read as the nearest thing in the frame.
				band(66, 0, core.RGB(0x0f1a2b), mid,
					[2]float32{10, 104}, [2]float32{330, 96}, [2]float32{540, 112}, [2]float32{780, 100}),
				// The walker, standing on the near ridge's crest between humps.
				core.View{Style: core.Style{
					Position: core.Absolute, Left: core.Pt(246), Bottom: core.Pt(66),
					AlignItems: core.AlignCenter,
					Animated:   &core.AnimatedStyle{TranslateX: mid},
				}, Children: []core.Node{
					Disc(10, colAmber, nil),
					core.View{Style: core.Style{Width: core.Pt(3), Height: core.Pt(15), Background: colAmber, CornerRadius: 2}},
				}},
				// Fade up from black, then a scrim so the landing is readable.
				core.View{Style: core.Style{
					Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					Background: core.RGB(0x000000), Animated: &core.AnimatedStyle{Opacity: lift},
				}},
				core.View{Style: core.Style{
					Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter,
				}, Children: []core.Node{
					core.View{Style: core.Style{
						Padding: core.Symmetric(14, 20), Gap: 4, CornerRadius: 14,
						Background: core.RGBA(0x060a10, 0.62),
						Animated:   &core.AnimatedStyle{Opacity: endIn, TranslateY: endY},
					}, Children: []core.Node{
						core.Text{Content: "一镜到底", Style: core.Style{
							FontSize: 22, FontWeight: core.FontWeightBold, Color: colText,
						}},
						core.Text{Content: "没有剪切，只有速度差", Style: core.Style{FontSize: 11, Color: colMuted}},
					}},
				}},
			}}
		}})
}
