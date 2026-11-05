package main

import (
	"math"
	"sort"
	"time"

	"github.com/oneo-me/nsgo/core"
)

// page3D answers a question the camera pack raises. Those 24 cards move a
// picture in two dimensions, because that is all a view transform can do. A
// camera move needs three. NSGo has no 3D engine — no meshes, no vertex
// stage, no depth buffer — but there are three honest ways to put a 3D scene
// on screen, and this page shows each of them:
//
//  1. core.Shader with a ray marcher: the camera is three uniforms and the
//     GPU solves the picture per pixel. The platform steps animated uniforms
//     on its own clock, so a camera move here costs zero commits.
//  2. core.Canvas with the projection done in Go: full control, hit testing
//     included, at the price of re-recording the drawing every frame.
//  3. Style.Perspective with RotateY: CSS 3D, right for cards and panels and
//     only for them — there is no preserve-3d, so children stay flat and a
//     cube would have to be sorted by hand.
func page3D(cx *core.Scope) core.Node {
	return grid(
		shotScene(cx),
		shotDolly(cx),
		shotTruck(cx),
		shotCrane(cx),
		shotOrbit(cx),
		shotDollyZoom(cx),
		shotPan(cx),
		cubeCanvas(cx),
		cubeCSS(cx),
	)
}

// shot3D is one camera move over the ray-marched scene: where the camera
// starts and ends, where it looks, its field of view, and a fill tint.
type shot3D struct {
	ID, Title, Note string
	Accent          core.Color
	Duration        time.Duration
	Loop            bool
	// Height is the stage's height in points; 150 when zero. The camera pack
	// borrows these tiles at its own 110, so both pages keep their rhythm.
	Height    float32
	Eye, Look [2][3]float64
	Fov       [2]float64
	Tint      [2]float64
}

// cam3D builds a tile: the same shader every time, with the camera's numbers
// as animated uniforms. Nothing renders on the engine once the tile is
// mounted — the platform moves the uniforms and the GPU follows.
func cam3D(cx *core.Scope, s shot3D) core.Node {
	fov := s.Fov
	if fov[1] == 0 {
		fov = [2]float64{0.85, 0.85}
	}
	height := s.Height
	if height == 0 {
		height = 150
	}
	return Fx(cx, FxSpec{ID: s.ID, Title: s.Title, Note: s.Note, Accent: s.Accent,
		Duration: s.Duration, Loop: s.Loop, Height: height,
		Stage: func(p *core.AnimatedValue) core.Node {
			axe := func(i int) core.AnimatedNode { return lerp(p, s.Eye[0][i], s.Eye[1][i]) }
			alk := func(i int) core.AnimatedNode { return lerp(p, s.Look[0][i], s.Look[1][i]) }
			return core.Shader{
				Program: rayProgram,
				Label:   s.Title,
				Style:   stageFill,
				Uniforms: core.Uniforms{
					"Eye":  []any{axe(0), axe(1), axe(2)},
					"Look": []any{alk(0), alk(1), alk(2)},
					"Fov":  lerp(p, fov[0], fov[1]),
					"Tint": lerp(p, s.Tint[0], s.Tint[1]),
				},
			}
		}})
}

// shotScene holds the camera still: three bodies on a floor, lit once. This
// is the frame the other six moves are made of.
func shotScene(cx *core.Scope) core.Node {
	return cam3D(cx, shot3D{
		ID: "3D 01", Title: "Ray-marched scene", Accent: colCyan, Duration: 1200 * ms,
		Note: "三个球加一个立方体落在地面上：光线步进在片元着色器里解出每个像素，没有网格，也没有顶点缓冲。",
		Eye:  [2][3]float64{{0.55, 0.95, 2.75}, {0.55, 0.95, 2.75}},
		Look: [2][3]float64{{0.0, 0.50, 0.15}, {0.0, 0.50, 0.15}},
	})
}

// shotDolly pushes the camera along its own axis: the subject grows, the
// background barely moves, and the parallax between them is the tell.
func shotDolly(cx *core.Scope) core.Node {
	return cam3D(cx, shot3D{
		ID: "3D 02", Title: "Dolly in / out", Accent: colBlue, Duration: 3600 * ms, Loop: true,
		Note: "相机沿视线推近再拉远：主体变大、背景几乎不动，两者的视差就是推轨的证据。",
		Eye:  [2][3]float64{{0.5, 1.00, 3.4}, {0.5, 0.90, 1.7}},
		Look: [2][3]float64{{0.0, 0.50, 0.15}, {0.0, 0.50, 0.15}},
	})
}

// shotTruck slides the camera sideways: near objects sweep, far ones crawl.
func shotTruck(cx *core.Scope) core.Node {
	return cam3D(cx, shot3D{
		ID: "3D 03", Title: "Truck (lateral)", Accent: colGreen, Duration: 3800 * ms, Loop: true,
		Note: "机身横移、镜头不转：近处的球掠过画面，远处的地平线几乎不动——横移和摇镜的区别就在这里。",
		Eye:  [2][3]float64{{-1.7, 0.92, 2.45}, {1.7, 0.92, 2.45}},
		Look: [2][3]float64{{0.0, 0.50, 0.15}, {0.0, 0.50, 0.15}},
	})
}

// shotCrane lifts the camera and tilts down with it: the floor opens up.
func shotCrane(cx *core.Scope) core.Node {
	return cam3D(cx, shot3D{
		ID: "3D 04", Title: "Crane up", Accent: colAmber, Duration: 3600 * ms, Loop: true,
		Note: "机位升高、镜头同时下俯：地面在画面里展开，这是升降，不是缩放。",
		Eye:  [2][3]float64{{0.3, 0.50, 2.85}, {0.3, 2.60, 3.4}},
		Look: [2][3]float64{{0.0, 0.50, 0.15}, {0.0, 0.05, -0.6}},
	})
}

// shotOrbit walks the camera round the subject on a circle. A circle needs
// sin and cos, and an animated lane can only interpolate, so the path is a
// twelve-stop polyline — still uniforms, still zero commits.
func shotOrbit(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "3D 05", Title: "Orbit", Accent: colViolet, Duration: 5200 * ms, Loop: true, Height: 150,
		Note: "绕着主体转一圈：x 和 z 各是一条十二段折线（动画通道只会插值，画不出正弦），视差来自真实的机位变化。",
		Stage: func(p *core.AnimatedValue) core.Node {
			return core.Shader{
				Program: rayProgram,
				Label:   "Orbit",
				Style:   stageFill,
				Uniforms: core.Uniforms{
					"Eye":  []any{circleLane(p, 0, 2.4, true), constant(cx, 0.92), circleLane(p, 0, 2.4, false)},
					"Look": []any{constant(cx, 0.0), constant(cx, 0.55), constant(cx, 0.0)},
					"Fov":  constant(cx, 0.85),
					"Tint": circleLane(p, 0, 0.35, true),
				},
			}
		}})
}

// circleLane is cos or sin of one turn, sampled into an interpolation.
func circleLane(p *core.AnimatedValue, phase, radius float64, cosine bool) core.AnimatedNode {
	const n = 12
	in := make([]float64, n+1)
	out := make([]float64, n+1)
	for i := 0; i <= n; i++ {
		t := float64(i) / n
		in[i] = t
		a := phase + 2*math.Pi*t
		if cosine {
			out[i] = radius * math.Cos(a)
		} else {
			out[i] = radius * math.Sin(a)
		}
	}
	return core.Interpolate(p, core.Interpolation{Input: in, Output: out})
}

// shotDollyZoom is Vertigo: the camera pushes in while the lens widens, so
// the subject keeps its size and the background stretches instead.
func shotDollyZoom(cx *core.Scope) core.Node {
	return cam3D(cx, shot3D{
		ID: "3D 06", Title: "Dolly zoom", Accent: colPink, Duration: 4200 * ms, Loop: true,
		Note: "推近的同时把视角拉大：主体尺寸不变、背景被抽走。希区柯克变焦要两个量同时动才算得出来。",
		Eye:  [2][3]float64{{0.0, 0.92, 1.75}, {0.0, 0.92, 3.9}},
		Look: [2][3]float64{{0.0, 0.58, 0.0}, {0.0, 0.58, 0.0}},
		Fov:  [2]float64{0.60, 1.24},
	})
}

// shotPan turns the camera in place: only the look target moves.
func shotPan(cx *core.Scope) core.Node {
	return cam3D(cx, shot3D{
		ID: "3D 07", Title: "Pan (yaw)", Accent: colGreen, Duration: 3400 * ms, Loop: true,
		Note: "机位不动、只转头：画面同样在平移，透视却在变——这才是摇镜，和横移不是一回事。",
		Eye:  [2][3]float64{{0.0, 0.95, 2.65}, {0.0, 0.95, 2.65}},
		Look: [2][3]float64{{-1.3, 0.50, 0.15}, {1.3, 0.50, 0.15}},
	})
}

// cubeCanvas projects a cube in Go and draws it with core.Canvas: eight
// vertices, six quads, painter's sort. This is the path to take when the
// scene must be hit tested or built from data, and its cost is a re-record
// per frame.
func cubeCanvas(cx *core.Scope) core.Node {
	angle := core.NewSignal(cx, 0.0)
	var tick func()
	tick = func() {
		angle.Update(func(a float64) float64 {
			a += 0.03
			if a > 2*math.Pi {
				a -= 2 * math.Pi
			}
			return a
		})
		cx.AfterFunc(33*ms, tick)
	}
	cx.AfterFunc(33*ms, tick)

	return Fx(cx, FxSpec{ID: "3D 08", Title: "CPU projection", Accent: colBlue, Height: 150,
		Note: "同一件事交给 CPU：Go 里算 8 个顶点、投影、按深度排序，再让 Canvas 画出来。慢，但每个顶点都在你手里。",
		Stage: func(p *core.AnimatedValue) core.Node {
			return core.Canvas{
				Label: "Rotating cube projected on the CPU",
				Style: stageFill,
				Draw: func(d *core.Drawing, size core.Size) {
					w, h := float32(size.Width), float32(size.Height)
					d.SetFillColor(core.RGB(0x080d14))
					d.FillRect(0, 0, w, h)

					a := float32(angle.Get())
					spin, tip := a, float32(0.5)
					cx0, cy0 := w/2, h/2+4
					// unit is pixels per world unit at the lens: without it the
					// projection stays in units and the cube lands on one pixel.
					eye, lens := float32(4.9), float32(2.7)
					unit := h * 0.40

					project := func(x, y, z float32) (float32, float32, float32) {
						xr := x*float32(math.Cos(float64(spin))) + z*float32(math.Sin(float64(spin)))
						zr := -x*float32(math.Sin(float64(spin))) + z*float32(math.Cos(float64(spin)))
						yr := y*float32(math.Cos(float64(tip))) - zr*float32(math.Sin(float64(tip)))
						zr = y*float32(math.Sin(float64(tip))) + zr*float32(math.Cos(float64(tip)))
						k := lens / (eye - zr)
						return cx0 + xr*k*unit, cy0 - yr*k*unit, zr
					}

					corner := func(i int) (float32, float32, float32) {
						x, y, z := float32(1), float32(1), float32(1)
						if i&1 == 1 {
							x = -1
						}
						if i&2 == 2 {
							y = -1
						}
						if i&4 == 4 {
							z = -1
						}
						return x, y, z
					}

					faces := [6][4]int{{0, 1, 3, 2}, {4, 6, 7, 5}, {0, 2, 6, 4}, {1, 5, 7, 3}, {0, 4, 5, 1}, {2, 3, 7, 6}}
					light := [6]float32{1.0, 0.34, 0.72, 0.5, 0.86, 0.6}

					type quad struct {
						pts   [4][2]float32
						depth float32
						tone  float32
					}
					quads := make([]quad, 0, 6)
					for i, f := range faces {
						q := quad{tone: light[i]}
						for k, vi := range f {
							x, y, z := project(corner(vi))
							q.pts[k] = [2]float32{x, y}
							q.depth += z
						}
						q.depth /= 4
						quads = append(quads, q)
					}
					sort.Slice(quads, func(i, j int) bool { return quads[i].depth < quads[j].depth })

					for _, q := range quads {
						g := float32(0.55) + 0.42*q.tone
						d.SetFillColor(core.RGBA(0x22405f, g))
						d.SetStrokeColor(core.RGBA(0x8fb4dc, 0.55))
						d.SetLineWidth(1)
						d.BeginPath()
						d.MoveTo(q.pts[0][0], q.pts[0][1])
						d.LineTo(q.pts[1][0], q.pts[1][1])
						d.LineTo(q.pts[2][0], q.pts[2][1])
						d.LineTo(q.pts[3][0], q.pts[3][1])
						d.Fill()
						d.Stroke()
					}

					d.SetFillColor(core.RGBA(0x9fb6d4, 0.5))
					d.FillText("8 vertices · 6 quads · painter sort", 10, h-8, core.Font{Size: 9})
				},
			}
		}})
}

// cubeCSS is the third way: one panel turned out of the plane by a CSS
// transform. Correct for cards, and only for them.
func cubeCSS(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "3D 09", Title: "CSS 3D", Accent: colViolet, Height: 150,
		Note: "Style.Perspective 加 RotateY：UI 卡片该用的就是这一档。没有 preserve-3d，子节点仍是平的，所以它做不了场景。",
		Stage: func(p *core.AnimatedValue) core.Node {
			// BackfaceHidden is set, but backface culling is not the same on
			// the three backends, so the faces also cross-fade by hand: the
			// front side is opaque until the turn passes 90 degrees.
			front := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.24, 0.26, 1}, Output: []float64{1, 1, 0, 0}})
			back := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.24, 0.26, 1}, Output: []float64{0, 0, 1, 1}})
			face := func(label string, accent core.Color, flip bool, show core.AnimatedNode) core.Node {
				style := core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					CornerRadius: 12, Background: colCard, BorderWidth: 1, BorderColor: accent,
					AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter,
					Animated: &core.AnimatedStyle{Opacity: show}}
				if flip {
					style.Transform = core.Transform{RotateY: 180}
				}
				return core.View{Style: style, Children: []core.Node{
					core.Text{Content: label, Style: core.Style{FontSize: 12, FontWeight: core.FontWeightMedium, Color: accent}},
				}}
			}
			// The turning box has a size of its own: rotation happens around
			// the node's centre, so a full-bleed child would swing its face
			// across the whole card.
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0),
				Right: core.Pt(0), Bottom: core.Pt(0), AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter,
				Perspective: 620},
				Children: []core.Node{
					core.View{Style: core.Style{Width: core.Pt(168), Height: core.Pt(100),
						Animated: &core.AnimatedStyle{RotateY: lerp(p, -24, 204)}},
						Children: []core.Node{
							face("FRONT", colCyan, false, front),
							face("BACK", colViolet, true, back),
						}},
				}}
		}})
}
