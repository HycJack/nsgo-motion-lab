package main

import (
	"fmt"
	"time"

	"github.com/oneo-me/nsgo/core"
)

// shaderPage collects the core.Shader examples. A shader is written in
// Kage — Go's syntax with GLSL's types — compiled once by
// core.MustCompileShader, and run by the GPU for every pixel of its box.
// Values that move are animated lanes, so the platform redraws them
// without the engine's help.
func shaderPage(cx *core.Scope) core.Node {
	return grid(
		shaderPlasma(cx),
		shaderPointer(cx),
		shaderSDF(cx),
		shaderRipple(cx),
	)
}

// shaderSDF draws three shapes from signed distance fields, antialiased
// by the pixel's size — the usual way a shader draws crisp geometry.
func shaderSDF(cx *core.Scope) core.Node {
	clock := core.NewAnimatedValue(cx, 0)
	cx.StartAnimation(core.Loop(core.Timing(clock, 6.283, core.TimingConfig{
		Duration: 6 * time.Second,
	}), 0), nil)
	return card("Distance fields", "Signed distance fields antialiased by 1/min(rectSize()) — shapes without meshes.", colGreen,
		core.Shader{
			Program:  sdfProgram,
			Label:    "Circle, square and ring from distance fields",
			TestID:   "shader-sdf",
			Style:    well(190),
			Uniforms: core.Uniforms{"Time": clock},
		},
		note("min, max, abs, length and smoothstep are Kage's GLSL builtins."),
	)
}

// shaderPlasma drifts two octaves of noise through a scrolling sine
// ramp, its clock an animated lane.
func shaderPlasma(cx *core.Scope) core.Node {
	clock := core.NewAnimatedValue(cx, 0)
	cx.StartAnimation(core.Loop(core.Timing(clock, 60, core.TimingConfig{
		Duration: 60 * time.Second,
	}), 0), nil)
	return card("Plasma", "One animated lane: the platform steps the clock and the GPU draws, with no commit per frame.", colGreen,
		core.Shader{
			Program:  plasmaProgram,
			Label:    "Drifting plasma",
			TestID:   "shader-plasma",
			Style:    well(190),
			Uniforms: core.Uniforms{"Time": clock},
		},
		note("vec4 Fragment(dst vec4, pos vec2) — pos in points, rectSize() the box."),
	)
}

// shaderPointer lights a dot grid around the pointer. The glow is set
// where the pointer moves, as a signal, and the grid's spacing cycles
// through a function, which costs a commit — the distinction between
// Uniforms and UniformsFunc.
func shaderPointer(cx *core.Scope) core.Node {
	x, y := core.NewAnimatedValue(cx, 0), core.NewAnimatedValue(cx, 0)
	presence := core.NewAnimatedValue(cx, 0)
	spacing := core.NewSignal(cx, 1.0)
	inside := false

	follow := func(p core.Point) {
		x.Set(float64(p.X))
		y.Set(float64(p.Y))
		if !inside {
			inside = true
			presence.Set(1)
		}
	}
	leave := func() {
		inside = false
		cx.StartAnimation(core.Timing(presence, 0, core.TimingConfig{Duration: 500 * time.Millisecond}), nil)
	}
	cycle := func() {
		spacing.Update(func(v float64) float64 {
			switch {
			case v < 1.4:
				return 1.8
			case v < 2.2:
				return 0.7
			default:
				return 1.0
			}
		})
	}
	return card("Pointer grid", "A signal sets the glow on each move; UniformsFunc resends the batch when spacing changes.", colCyan,
		core.View{
			Style:       core.Style{Height: core.Pt(190)},
			OnMouseMove: follow,
			OnHoverOut:  leave,
			Children: []core.Node{
				core.Shader{
					Program: pointerProgram,
					Label:   "Dot grid following the pointer",
					TestID:  "shader-pointer",
					Style:   well(190),
					Uniforms: core.Uniforms{
						"Pointer": []core.AnimatedNode{x, y, presence},
					},
					UniformsFunc: func() core.Uniforms {
						return core.Uniforms{"Scale": spacing.Get()}
					},
				},
			},
		},
		row(8, chip("Spacing", colCyan, cycle), noteFunc(func() string {
			return fmt.Sprintf("uniform Scale = %.2g", spacing.Get())
		})),
	)
}

// shaderRipple sends a ring out where the box is pressed: the press only
// sets animated values, which the backend reads as it draws the batch.
func shaderRipple(cx *core.Scope) core.Node {
	x, y := core.NewAnimatedValue(cx, 0), core.NewAnimatedValue(cx, 0)
	age, weight := core.NewAnimatedValue(cx, 1), core.NewAnimatedValue(cx, 0)
	press := func(e core.PressEvent) {
		x.Set(float64(e.Position.X))
		y.Set(float64(e.Position.Y))
		// Another animation of the same value ends the one before it.
		age.Set(0)
		weight.Set(1)
		cx.StartAnimation(core.Timing(age, 1, core.TimingConfig{Duration: 1200 * time.Millisecond}), func(finished bool) {
			if finished {
				weight.Set(0)
			}
		})
	}
	return card("Press ripple", "A press sets four animated lanes; the GPU draws the ring as they move.", colPink,
		core.Shader{
			Program:  rippleProgram,
			Label:    "Ripple on press",
			TestID:   "shader-ripple",
			Style:    well(190),
			OnPress:  press,
			Uniforms: core.Uniforms{"Press": []core.AnimatedNode{x, y, age, weight}},
		},
		note("Tap anywhere in the box."),
	)
}

// plasmaProgram is the plasma shader.
var plasmaProgram = core.MustCompileShader(plasmaSource)

const plasmaSource = `package main

var Time float

func hash(p vec2) float {
	p3 := fract(vec3(p.x, p.y, p.x) * 0.1031)
	p3 += dot(p3, p3.yzx+33.33)
	return fract((p3.x + p3.y) * p3.z)
}

func noise(p vec2) float {
	i := floor(p)
	f := fract(p)
	u := f * f * (3 - 2*f)
	a := hash(i)
	b := hash(i + vec2(1, 0))
	c := hash(i + vec2(0, 1))
	d := hash(i + vec2(1, 1))
	return mix(mix(a, b, u.x), mix(c, d, u.x), u.y)
}

func Fragment(dst vec4, pos vec2) vec4 {
	size := rectSize()
	uv := pos / size
	p := pos / min(size.x, size.y)
	t := Time

	n := noise(p*3.2 + vec2(t*0.30, t*0.21))
	n += 0.5 * noise(p*6.4 - vec2(t*0.19, -t*0.27))

	v := 0.5 + 0.5*sin(6.2831*(uv.x+uv.y) + 3.0*n + t*0.9)
	col := mix(vec3(0.04, 0.07, 0.13), vec3(0.20, 0.82, 0.86), v)
	col = mix(col, vec3(0.64, 0.53, 0.98), 0.42*n)

	return vec4(col, 1.0)
}
`

// pointerProgram is the dot grid.
var pointerProgram = core.MustCompileShader(pointerSource)

const pointerSource = `package main

var Pointer vec3
var Scale float

func hash(p vec2) float {
	p3 := fract(vec3(p.x, p.y, p.x) * 0.1031)
	p3 += dot(p3, p3.yzx+33.33)
	return fract((p3.x + p3.y) * p3.z)
}

func Fragment(dst vec4, pos vec2) vec4 {
	spacing := 20 * Scale
	center := (floor(pos/spacing) + 0.5) * spacing

	d := center - Pointer.xy
	glow := Pointer.z * exp(-dot(d, d)/(2*95*95))

	n := hash(center)
	radius := 0.55 + 0.75*n + 2.2*glow
	cover := 1 - smoothstep(radius-0.8, radius+0.8, distance(pos, center))

	a := clamp((0.10 + 0.10*n) + 0.95*glow, 0.0, 1.0) * cover
	col := mix(vec3(0.42, 0.56, 0.72), vec3(0.22, 0.86, 0.88), clamp(glow, 0.0, 1.0))

	return vec4(col*a, a)
}
`

// rippleProgram sends a ring out from Press.xy as Press.z (the age, 0 to
// 1) grows.
var rippleProgram = core.MustCompileShader(rippleSource)
const rippleSource = `package main

var Press vec4

func Fragment(dst vec4, pos vec2) vec4 {
	size := rectSize()
	uv := pos / size

	bg := mix(vec3(0.04, 0.06, 0.10), vec3(0.08, 0.13, 0.20), uv.y)

	age := clamp(Press.z, 0.0, 1.0)
	fade := 1.0 - age
	radius := age * max(size.x, size.y) * 0.55
	width := 16.0 + 44.0*age
	dist := length(pos - Press.xy)

	k := (dist - radius) / width
	ring := exp(-k*k) * pow(fade, 0.7)
	core := exp(-dist*dist/1400.0) * fade * fade

	light := clamp((ring + 0.7*core) * Press.w, 0.0, 1.0)
	col := bg + vec3(0.20, 0.66, 0.90) * light
	col = col + vec3(0.26, 0.86, 0.70) * core * Press.w * 0.6

	return vec4(col, 1.0)
}
`

// sdfProgram is the distance-field shader.
var sdfProgram = core.MustCompileShader(sdfSource)

const sdfSource = `package main

var Time float

func sdCircle(p vec2, r float) float {
	return length(p) - r
}

func sdBox(p vec2, b vec2) float {
	d := abs(p) - b
	return length(max(d, 0)) + min(max(d.x, d.y), 0)
}

func sdRing(p vec2, r float, w float) float {
	return abs(length(p) - r) - w
}

func Fragment(dst vec4, pos vec2) vec4 {
	size := rectSize()
	uv := pos / size
	p := (pos - size*0.5) / min(size.x, size.y) * 2.0

	// A ring that breathes, a disc and a square.
	ring := sdRing(p - vec2(0.62, 0.02), 0.24 + 0.025*sin(Time), 0.045)
	disc := sdCircle(p - vec2(-0.6, 0.0), 0.28)
	box := sdBox(p, vec2(0.2, 0.2))

	d := min(min(ring, disc), box)
	px := 2.0 / min(size.x, size.y)
	cover := 1.0 - smoothstep(-px, px, d)

	bg := mix(vec3(0.04, 0.06, 0.10), vec3(0.07, 0.11, 0.16), uv.y)
	shade := mix(vec3(0.24, 0.85, 0.55), vec3(0.16, 0.62, 0.92), smoothstep(-0.6, 0.6, p.x))
	glowv := exp(-max(d, 0.0)*7.0) * 0.22

	col := mix(bg, shade, cover)
	col = col + shade * glowv * (1.0 - cover)
	return vec4(col, 1.0)
}
`