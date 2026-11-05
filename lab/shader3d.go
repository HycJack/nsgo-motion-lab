package main

import "github.com/oneo-me/nsgo/core"

// rayProgram is a ray marcher: a fragment shader that walks a ray through a
// signed distance field until it meets a surface, then shades that surface
// from its normal. There is no geometry anywhere — three spheres and a box
// exist only as formulas — so the camera is the whole scene's state, and it
// arrives as three uniforms.
//
// Eye and Look are the camera's position and target, Fov its vertical field
// of view in radians, Tint a 0..1 mix for the fill light. Because all four
// can be animated lanes (core.Uniforms takes AnimatedNodes), a camera move
// is stepped by the platform, not by the engine.
var rayProgram = core.MustCompileShader(raySource)

const raySource = `package main

var Eye vec3
var Look vec3
var Fov float
var Tint float

func sdSphere(p vec3, r float) float {
	return length(p) - r
}

func sdBox(p vec3, b vec3) float {
	d := abs(p) - b
	return length(max(d, 0.0)) + min(max(d.x, max(d.y, d.z)), 0.0)
}

// scene is the whole world: a floor, two spheres and a cube.
func scene(p vec3) float {
	ground := p.y
	s1 := sdSphere(p - vec3(-0.98, 0.60, 0.05), 0.60)
	s2 := sdSphere(p - vec3(0.88, 0.42, -0.42), 0.42)
	b1 := sdBox(p - vec3(0.10, 0.30, 0.92), vec3(0.30, 0.30, 0.30))
	return min(ground, min(s1, min(s2, b1)))
}

// normalAt is the field's gradient: six samples, the usual trick.
func normalAt(p vec3) vec3 {
	h := 0.0018
	n := vec3(
		scene(p + vec3(h, 0.0, 0.0)) - scene(p - vec3(h, 0.0, 0.0)),
		scene(p + vec3(0.0, h, 0.0)) - scene(p - vec3(0.0, h, 0.0)),
		scene(p + vec3(0.0, 0.0, h)) - scene(p - vec3(0.0, 0.0, h)))
	return normalize(n)
}

// shadowAt walks toward the key light and keeps the tightest gap: a soft
// shadow without a shadow map, and without a second pass.
func shadowAt(p vec3) float {
	light := normalize(vec3(-0.55, 0.80, 0.34))
	pen := 1.0
	t := 0.12
	for i := 0; i < 16; i++ {
		h := scene(p + light*t)
		pen = min(pen, 20.0*h/t)
		t += 0.30
	}
	return clamp(pen, 0.18, 1.0)
}

func Fragment(dst vec4, pos vec2) vec4 {
	size := rectSize()
	uv := (pos*2.0 - size) / size.y

	// A camera basis from two uniforms: the direction it looks, and up.
	fwd := normalize(Look - Eye)
	right := normalize(cross(fwd, vec3(0.0, 1.0, 0.0)))
	up := cross(right, fwd)
	// A ray at the edge of the box is tan(fov/2) off the axis: that is what
	// makes Fov a field of view rather than a zoom factor.
	lens := tan(Fov*0.5)
	dir := normalize(right*(uv.x*lens) - up*(uv.y*lens) + fwd)

	// The march. There is no break here on purpose: a hit freezes t by
	// multiplying the step, which keeps the loop's trip count fixed — the
	// shape every GPU translation stays happy with.
	t := 0.0
	hit := 0.0
	for i := 0; i < 96; i++ {
		d := scene(Eye + dir*t)
		if d < 0.002 {
			hit = 1.0
		}
		t += max(d, 0.002) * (1.0 - hit)
	}

	p := Eye + dir*t
	n := vec3(0.0, 1.0, 0.0)
	if hit > 0.5 {
		n = normalAt(p)
	}

	// Three materials, told apart by where the surface is: warm spheres, a
	// teal cube, a floor with a grid on it.
	albedo := vec3(0.80, 0.45, 0.28)
	grid := abs(fract(p.x*0.45) - 0.5) * abs(fract(p.z*0.45) - 0.5)
	if p.y < 0.05 {
		// Grid lines, faded with distance: at a grazing angle the pattern
		// would otherwise alias into rings.
		line := (1.0 - smoothstep(0.0, 0.30, grid)) * exp(-0.10*t)
		albedo = mix(vec3(0.100, 0.125, 0.160), vec3(0.260, 0.300, 0.360), line)
	}
	if abs(p.x - 0.10) < 0.34 && abs(p.y - 0.30) < 0.34 && abs(p.z - 0.92) < 0.34 {
		albedo = vec3(0.22, 0.58, 0.64)
	}

	key := normalize(vec3(-0.55, 0.80, 0.34))
	fillDir := normalize(vec3(0.72, 0.24, 0.52))
	diff := max(dot(n, key), 0.0)
	sh := 1.0
	if hit > 0.5 {
		sh = shadowAt(p + n*0.03)
	}

	keyCol := vec3(1.35, 1.16, 0.86)
	fillCol := mix(vec3(0.24, 0.52, 0.72), vec3(0.62, 0.36, 0.92), Tint)
	amb := vec3(0.17, 0.21, 0.28) * (0.55 + 0.45*max(n.y, 0.0))

	col := albedo * (keyCol*(diff*sh) + fillCol*(0.34*max(dot(n, fillDir), 0.0)) + amb)
	spec := pow(max(dot(reflect(dir, n), key), 0.0), 28.0) * diff * sh
	col = col + keyCol*spec*0.32

	sky := mix(vec3(0.040, 0.060, 0.100), vec3(0.090, 0.135, 0.205), clamp(dir.y*0.9 + 0.35, 0.0, 1.0))
	sky = sky + vec3(0.10, 0.075, 0.055)*exp(-max(dir.y, 0.0)*9.0)
	if hit < 0.5 {
		col = sky
	} else {
		fog := 1.0 - exp(-t*0.085)
		col = mix(col, sky, clamp(fog, 0.0, 1.0))
	}
	col = pow(col, vec3(0.92, 0.95, 1.0))
	return vec4(col, 1.0)
}
`
