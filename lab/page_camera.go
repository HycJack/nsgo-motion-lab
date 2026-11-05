package main

// page_camera.go is the camera pack (motion-prompt #025–#048): a camera in
// front of one small landscape, and the edits that join shots.
//
// The landscape is built once by camScene and every tile moves its layers by
// different amounts — that difference is what makes a dolly look like a dolly
// and a pan look like a pan:
//
//	pan    a rotation, so every layer slides the same distance
//	dolly  the camera moves, so the near layer grows and passes faster
//	zoom   no movement at all, so everything grows by the same ratio
//
// That model is honest for anything the depth already gives the picture, and
// wrong for the two shots whose whole point is depth: an orbit has to show the
// other side of the set, and a crane has to look down on it. Those two are
// rendered on the ray-marched set from page_3d.go instead, with the camera as
// uniforms — the same shader, one tile shorter.
//
// Catalog names and timings: github.com/showheyohtaki/motion-prompt (MIT).

import (
	"math"

	"github.com/oneo-me/nsgo/core"
)

// arcLane is cos or sin of an arc, sampled into an interpolation: an animated
// lane can only interpolate, so a curve becomes a short polyline.
func arcLane(p *core.AnimatedValue, radius, from, to float64, cosine bool) core.AnimatedNode {
	const n = 4
	in := make([]float64, n+1)
	out := make([]float64, n+1)
	for i := 0; i <= n; i++ {
		t := float64(i) / n
		in[i] = t
		a := from + (to-from)*t
		if cosine {
			out[i] = radius * math.Cos(a)
		} else {
			out[i] = radius * math.Sin(a)
		}
	}
	return core.Interpolate(p, core.Interpolation{Input: in, Output: out})
}

// camLayers are the moving parts of the scene. A tile sets only the ones it
// needs; a nil style leaves that layer still.
type camLayers struct {
	Sky    *core.AnimatedStyle
	Sun    *core.AnimatedStyle
	Far    *core.AnimatedStyle
	Near   *core.AnimatedStyle
	Figure *core.AnimatedStyle
}

// camScene is a horizon, a sun, two hills and a person standing between them.
// Everything is absolutely placed so a layer can be scaled or slid without
// disturbing its neighbours.
func camScene(sky, ground core.Color, l camLayers) core.Node {
	figure := core.View{Style: core.Style{
		Position: core.Absolute, Top: core.Pt(58), Left: core.Pt(92),
		AlignItems: core.AlignCenter, Animated: l.Figure,
	}, Children: []core.Node{
		Disc(10, core.RGB(0xf2dfc8), nil),
		Panel(13, 22, core.RGB(0x2a3646), nil),
	}}
	return core.View{Style: core.Style{
		Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0),
		Right: core.Pt(0), Bottom: core.Pt(0), Overflow: core.Hidden,
	}, Children: []core.Node{
		core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Height: core.Pt(86),
			Background: sky, Animated: l.Sky}},
		core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(86), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
			Background: ground, Animated: l.Sky}},
		core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(14), Left: core.Pt(38), Animated: l.Sun},
			Children: []core.Node{Disc(16, core.RGB(0xffd27a), nil)}},
		core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(58), Left: core.Pt(40), Animated: l.Far},
			Children: []core.Node{Disc(96, core.RGB(0x27405a), nil)}},
		core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(76), Left: core.Pt(120), Animated: l.Near},
			Children: []core.Node{Disc(124, core.RGB(0x1b2c40), nil)}},
		figure,
	}}
}

// camShots is the two-shot pair the edit tiles cut between.
func camShots(a, b *core.AnimatedStyle) core.Node {
	wide := camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{})
	return core.View{Style: core.Style{
		Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
	}, Children: []core.Node{
		core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0), Animated: a},
			Children: []core.Node{wide}},
		core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0), Animated: b},
			Children: []core.Node{camScene(core.RGB(0x5a2b3c), core.RGB(0x2a1420), camLayers{})}},
	}}
}

// pageCamera lays the pack out.
func pageCamera(cx *core.Scope) core.Node {
	return core.View{Children: []core.Node{
		core.View{Style: core.Style{Padding: core.Edges{Top: 14, Right: 18, Left: 18}}, Children: []core.Node{
			note("Camera pack 025–048 of motion-prompt (github.com/showheyohtaki/motion-prompt). One landscape, seen from 24 places; the code is NSGo's own. Layers carry the depth the picture already has (dolly 027 vs zoom 028 is the proof), and the two shots that need a viewpoint — orbit 030, crane 031 — run on the ray-marched set instead."),
		}},
		grid(
			camPan(cx),
			camTilt(cx),
			camDolly(cx),
			camZoom(cx),
			camTracking(cx),
			camOrbit(cx),
			camCrane(cx),
			camDollyZoom(cx),
			camRackFocus(cx),
			camHandheld(cx),
			camFPV(cx),
			camDutch(cx),
			camSlowMotion(cx),
			camSpeedRamp(cx),
			camTimelapse(cx),
			camFreeze(cx),
			camJumpCut(cx),
			camMatchCut(cx),
			camWhipPan(cx),
			camDissolve(cx),
			camWipe(cx),
			camSmashCut(cx),
			camJCut(cx),
			camLCut(cx),
		),
	}}
}

// #025 Pan: a rotation, so near and far slide by the same amount.
func camPan(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "025", Title: "Pan", Duration: 2400 * ms, Accent: colCyan, Height: 110,
		Note: "The camera turns instead of moving, so the hills and the sun all slide by the same 90 points.",
		Stage: func(p *core.AnimatedValue) core.Node {
			move := &core.AnimatedStyle{TranslateX: lerp(p, 60, -60)}
			return camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{Sun: &core.AnimatedStyle{TranslateX: lerp(p, 40, -40)}, Far: move, Near: move, Figure: move})
		}})
}

// #026 Tilt: the camera looks up, so the horizon sinks and the sky takes over.
func camTilt(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "026", Title: "Tilt", Duration: 2200 * ms, Accent: colCyan, Height: 110,
		Note: "Looking up from the person to the sky: the ground and the near hill drop out of frame first.",
		Stage: func(p *core.AnimatedValue) core.Node {
			return camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{
				Sun:    &core.AnimatedStyle{TranslateY: lerp(p, 0, 26)},
				Far:    &core.AnimatedStyle{TranslateY: lerp(p, 0, 46)},
				Near:   &core.AnimatedStyle{TranslateY: lerp(p, 0, 74)},
				Figure: &core.AnimatedStyle{TranslateY: lerp(p, 0, 96)},
			})
		}})
}

// #027 Dolly in: near layers grow and pass; the far hill barely changes.
func camDolly(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "027", Title: "Dolly in", Duration: 2200 * ms, Accent: colCyan, Height: 110,
		Note: "The camera moves in: the near hill swells by half again while the far hill grows by a twentieth.",
		Stage: func(p *core.AnimatedValue) core.Node {
			return camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{
				Far:    &core.AnimatedStyle{Scale: lerp(p, 1, 1.06)},
				Near:   &core.AnimatedStyle{Scale: lerp(p, 1, 1.5), TranslateY: lerp(p, 0, 18)},
				Sun:    &core.AnimatedStyle{Scale: lerp(p, 1, 1.08)},
				Figure: &core.AnimatedStyle{Scale: lerp(p, 1, 1.22), TranslateY: lerp(p, 0, -6)},
			})
		}})
}

// #028 Zoom in: no movement, so every layer grows by the same ratio.
func camZoom(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "028", Title: "Zoom in", Duration: 2200 * ms, Accent: colCyan, Height: 110,
		Note: "The camera stands still and only the focal length lengthens, so the person and the hills grow by the same 30%.",
		Stage: func(p *core.AnimatedValue) core.Node {
			s := &core.AnimatedStyle{Scale: lerp(p, 1, 1.3)}
			return camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{Sun: s, Far: s, Near: s, Figure: s})
		}})
}

// #029 Tracking: the person stays put while the ground sweeps past faster than
// the hills behind it.
func camTracking(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "029", Title: "Tracking shot", Duration: 2400 * ms, Accent: colGreen, Height: 110,
		Note: "The camera walks with the person, so he never moves on screen and everything else slides left at two speeds.",
		Stage: func(p *core.AnimatedValue) core.Node {
			return camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{
				Far:  &core.AnimatedStyle{TranslateX: lerp(p, 30, -30)},
				Near: &core.AnimatedStyle{TranslateX: lerp(p, 80, -80)},
				Sun:  &core.AnimatedStyle{TranslateX: lerp(p, 18, -18)},
			})
		}})
}

// #030 Orbit: the camera travels around the subject, so the background
// sweeps sideways and the person stays centred.
//
// A layer slide cannot do this shot: an orbit changes which side of the set
// is lit and which sphere hides which, and no picture of the front can show
// that. So this one runs on the ray-marched set, the camera on a 90-degree
// arc — x swings ±1.7 while z bulges from 1.7 to 2.4 and back.
func camOrbit(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "030", Title: "Orbit (arc shot)", Duration: 3200 * ms, Accent: colGreen, Height: 110,
		Note: "A true arc on a ray-marched set: the subject holds still while the cube slides across the far sphere and back — the parallax a flat picture cannot fake.",
		Stage: func(p *core.AnimatedValue) core.Node {
			return core.Shader{
				Program: rayProgram,
				Label:   "Orbit (arc shot)",
				Style:   stageFill,
				Uniforms: core.Uniforms{
					"Eye":  []any{arcLane(p, 2.4, -0.79, 0.79, false), constant(cx, 0.88), arcLane(p, 2.4, -0.79, 0.79, true)},
					"Look": []any{constant(cx, 0.0), constant(cx, 0.55), constant(cx, 0.0)},
					"Fov":  constant(cx, 0.85),
					"Tint": lerp(p, 0.0, 0.35),
				},
			}
		}})
}

// #031 Crane up: the camera rises and tilts down, so the person and the near
// hill sink out of the bottom of the frame.
//
// A crane changes the viewpoint's height, which is the one thing a stack of
// pictures cannot do: the floor has to open up. Knee height to twice the set's
// height, look target held on the subject, and the ray marcher does the rest.
func camCrane(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "031", Title: "Crane up", Duration: 3000 * ms, Accent: colGreen, Height: 110,
		Note: "Rising from knee height to twice the set's: the floor opens under the shot and the horizon drops, with the look target held on the subject.",
		Stage: func(p *core.AnimatedValue) core.Node {
			return core.Shader{
				Program: rayProgram,
				Label:   "Crane up",
				Style:   stageFill,
				Uniforms: core.Uniforms{
					"Eye":  []any{constant(cx, 0.62), lerp(p, 0.42, 2.55), lerp(p, 3.00, 1.80)},
					"Look": []any{constant(cx, 0.0), constant(cx, 0.55), constant(cx, 0.0)},
					"Fov":  constant(cx, 0.88),
					"Tint": lerp(p, 0.0, 0.2),
				},
			}
		}})
}

// #032 Dolly zoom: pull back while zooming in, so the person keeps his size
// and only the background closes in.
func camDollyZoom(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "032", Title: "Dolly zoom (Vertigo)", Duration: 2600 * ms, Accent: colViolet, Height: 110,
		Note: "The signature of the shot: the person does not change size at all, and everything behind him rushes outward.",
		Stage: func(p *core.AnimatedValue) core.Node {
			return camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{
				Far:    &core.AnimatedStyle{Scale: lerp(p, 1, 1.7), TranslateY: lerp(p, 0, 30)},
				Near:   &core.AnimatedStyle{Scale: lerp(p, 1, 2.1), TranslateY: lerp(p, 0, 46)},
				Sun:    &core.AnimatedStyle{Scale: lerp(p, 1, 1.5), TranslateY: lerp(p, 0, -14)},
				Sky:    &core.AnimatedStyle{Background: core.Interpolate(p, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1}, Colors: []core.Color{core.RGB(0x1b3a5c), core.RGB(0x2a2350)}})},
				Figure: nil,
			})
		}})
}

// #033 Rack focus: the foreground flower goes soft while the person behind it
// comes in; softness is shown by dimming, since there is no blur.
func camRackFocus(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "033", Title: "Rack focus", Duration: 2400 * ms, Accent: colAmber, Height: 110,
		Note: "Focus travels from the flower in front to the person behind: one dialogue dims out, the other sharpens, and never both at once.",
		Stage: func(p *core.AnimatedValue) core.Node {
			flower := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.45, 1}, Output: []float64{1, 0.5, 0.25}})
			person := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.55, 1}, Output: []float64{0.3, 0.45, 1}})
			return core.View{Style: stageFill, Children: []core.Node{
				camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{Figure: &core.AnimatedStyle{Opacity: person}}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(52), Left: core.Pt(30),
					Animated: &core.AnimatedStyle{Opacity: flower, Scale: lerp(p, 1.1, 0.9)}},
					Children: []core.Node{Disc(30, core.RGB(0xff8fb1), nil)}},
			}}
		}})
}

// #034 Handheld: irregular shake, never periodic.
func camHandheld(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "034", Title: "Handheld", Duration: 1600 * ms, Accent: colAmber, Height: 110,
		Note: "A shake that never repeats: the frame drifts a few points in four uneven beats.",
		Stage: func(p *core.AnimatedValue) core.Node {
			shake := core.Interpolate(p, core.Interpolation{
				Input:  []float64{0, 0.18, 0.31, 0.44, 0.62, 0.77, 0.9, 1},
				Output: []float64{0, 3.5, -2.5, 3, -3.5, 2, -1.5, 0},
			})
			bob := core.Interpolate(p, core.Interpolation{
				Input:  []float64{0, 0.22, 0.4, 0.58, 0.72, 0.88, 1},
				Output: []float64{0, -2.5, 2, -3, 1.5, -1, 0},
			})
			whole := &core.AnimatedStyle{TranslateX: shake, TranslateY: bob}
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0), Animated: whole},
				Children: []core.Node{camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{})}}
		}})
}

// #035 FPV: low, fast, between two rows of trees, banking.
func camFPV(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "035", Title: "FPV drone shot", Duration: 2200 * ms, Accent: colBlue, Height: 110,
		Note: "First person, low and fast: the trees rush past on both sides and the whole frame banks into the turn.",
		Stage: func(p *core.AnimatedValue) core.Node {
			tree := func(left float32, delay float64) core.Node {
				q := core.Interpolate(p, core.Interpolation{Input: []float64{0, 1 - delay, 1}, Output: []float64{1, 1.6, 2.6}})
				return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(34), Left: core.Pt(core.Points(left)),
					Animated: &core.AnimatedStyle{Scale: q, TranslateY: lerp(p, 0, 40)}},
					Children: []core.Node{Panel(8, 46, core.RGB(0x14261f), nil)}}
			}
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
				Animated: &core.AnimatedStyle{Rotate: lerp(p, 0, -7), Scale: lerp(p, 1, 1.15)}}, Children: []core.Node{
				camScene(core.RGB(0x1d5a4a), core.RGB(0x10281f), camLayers{}),
				tree(10, 0), tree(200, 0.25), tree(48, 0.5), tree(228, 0.7),
			}}
		}})
}

// #036 Dutch angle: the camera rolls and holds with the horizon on a diagonal.
func camDutch(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "036", Title: "Dutch angle", Duration: 2200 * ms, Accent: colViolet, Height: 110,
		Note: "The whole frame rolls fourteen degrees and stays there: a level camera held wrong on purpose.",
		Stage: func(p *core.AnimatedValue) core.Node {
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(-14), Left: core.Pt(-14), Right: core.Pt(-14), Bottom: core.Pt(-14),
				Animated: &core.AnimatedStyle{Rotate: lerp(p, 0, -14), Scale: lerp(p, 1, 1.16)}},
				Children: []core.Node{camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{})}}
		}})
}

// #037 Slow motion: a jump, stretched out, with leaves drifting down.
func camSlowMotion(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "037", Title: "Slow motion", Duration: 2600 * ms, Accent: colCyan, Height: 110,
		Note: "The camera is bolted down; the person rises and falls over two and a half seconds and the leaves fall slowly past him.",
		Stage: func(p *core.AnimatedValue) core.Node {
			jump := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.45, 1}, Output: []float64{0, -34, 0}})
			leaf := func(left float32, delay float64) core.Node {
				q := core.Interpolate(p, core.Interpolation{Input: []float64{delay, 1}, Output: []float64{0, 1}})
				return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(6), Left: core.Pt(core.Points(left)),
					Animated: &core.AnimatedStyle{TranslateY: lerp(q, 0, 92), Rotate: lerp(q, 0, 160), Opacity: core.Interpolate(q, core.Interpolation{Input: []float64{0, 0.8, 1}, Output: []float64{1, 1, 0}})}},
					Children: []core.Node{Panel(5, 3, core.RGB(0xc9d97a), nil)}}
			}
			return core.View{Style: stageFill, Children: []core.Node{
				camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{Figure: &core.AnimatedStyle{TranslateY: jump}}),
				leaf(30, 0.1), leaf(160, 0.3), leaf(190, 0.5),
			}}
		}})
}

// #038 Speed ramp: the same distance, first fast then slow, with the speed drawn.
func camSpeedRamp(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "038", Title: "Speed ramp", Duration: 2400 * ms, Accent: colGreen, Height: 110,
		Note: "One move, two speeds: he crosses most of the frame early and crawls through the last part; the line below is the speed.",
		Stage: func(p *core.AnimatedValue) core.Node {
			run := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.35, 1}, Output: []float64{0, 150, 178}})
			speed := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.35, 1}, Output: []float64{46, 6, 2}})
			return core.View{Style: stageFill, Children: []core.Node{
				camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{Figure: &core.AnimatedStyle{TranslateX: run}}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(98), Left: core.Pt(12),
					Animated: &core.AnimatedStyle{}},
					Children: []core.Node{core.View{Style: core.Style{Width: core.Pt(150), Height: core.Pt(2), Background: core.RGB(0x39d98a),
						Animated: &core.AnimatedStyle{ScaleX: core.Interpolate(speed, core.Interpolation{Input: []float64{2, 46}, Output: []float64{0.06, 1}})}}}}},
			}}
		}})
}

// #039 Time-lapse: the sun crosses the sky and the light changes with it.
func camTimelapse(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "039", Title: "Time-lapse", Duration: 2600 * ms, Accent: colAmber, Height: 110,
		Note: "A fixed camera and a whole day in two seconds: the sun crosses overhead and the sky goes from morning grey to sunset orange.",
		Stage: func(p *core.AnimatedValue) core.Node {
			return camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{
				Sky: &core.AnimatedStyle{Background: core.Interpolate(p, core.Interpolation{
					Input: []float64{0, 0.4, 0.75, 1}, Output: []float64{0, 1, 2, 3},
					Colors: []core.Color{core.RGB(0x2b3a52), core.RGB(0x3d6fa8), core.RGB(0xa8622f), core.RGB(0x35203a)},
				})},
				Sun: &core.AnimatedStyle{TranslateX: lerp(p, -26, 190), TranslateY: core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.5, 1}, Output: []float64{52, -8, 52}})},
			})
		}})
}

// #040 Freeze frame: the picture stops at the top of the jump, then resumes.
func camFreeze(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "040", Title: "Freeze frame", Duration: 2600 * ms, Accent: colCyan, Height: 110,
		Note: "The jump stops dead at its highest point and holds; the dark bar at the bottom says the clock has stopped too.",
		Stage: func(p *core.AnimatedValue) core.Node {
			hold := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.4, 0.75, 1}, Output: []float64{0, -30, -30, 0}})
			tint := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.45, 0.7, 1}, Output: []float64{0, 0.35, 0.35, 0}})
			return core.View{Style: stageFill, Children: []core.Node{
				camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{Figure: &core.AnimatedStyle{TranslateY: hold}}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					Background: core.RGB(0x000000), Animated: &core.AnimatedStyle{Opacity: tint}}},
			}}
		}})
}

// #041 Jump cut: same framing, the person is somewhere else.
func camJumpCut(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "041", Title: "Jump cut", Duration: 2000 * ms, Accent: colPink, Height: 110,
		Note: "No cut-away: the same shot, twice, with time removed in between. He is suddenly closer and the second hand has moved.",
		Stage: func(p *core.AnimatedValue) core.Node {
			jump := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.49, 0.5, 1}, Output: []float64{0, 0, 64, 64}})
			zoom := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.49, 0.5, 1}, Output: []float64{1, 1, 1.5, 1.5}})
			sun := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.49, 0.5, 1}, Output: []float64{0, 0, 18, 18}})
			return camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{
				Sun:    &core.AnimatedStyle{TranslateX: sun},
				Figure: &core.AnimatedStyle{TranslateX: jump, Scale: zoom},
			})
		}})
}

// #042 Match cut: a thrown ball becomes the sun at the same place and size.
func camMatchCut(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "042", Title: "Match cut", Duration: 2200 * ms, Accent: colAmber, Height: 110,
		Note: "Two shots joined by one shape: the ball leaves the hand, and where it was is the sun in the next shot.",
		Stage: func(p *core.AnimatedValue) core.Node {
			ball := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.48, 0.5, 1}, Output: []float64{0, 1, 1, 0}})
			sun := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.5, 0.56, 1}, Output: []float64{0, 0, 1, 1}})
			return core.View{Style: stageFill, Children: []core.Node{
				camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(26), Left: core.Pt(120),
					Animated: &core.AnimatedStyle{Opacity: ball, Scale: lerp(p, 0.6, 1.4)}},
					Children: []core.Node{Disc(16, core.RGB(0x9fe8ff), nil)}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(24), Left: core.Pt(118),
					Animated: &core.AnimatedStyle{Opacity: sun}},
					Children: []core.Node{Disc(20, core.RGB(0xffd27a), nil)}},
			}}
		}})
}

// #043 Whip pan: too fast to see, so the two shots are joined inside the streaked blur.
func camWhipPan(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "043", Title: "Whip pan", Duration: 2000 * ms, Accent: colBlue, Height: 110,
		Note: "The camera swings so fast that the first shot is gone before you can read it; the second one arrives out of the same blur.",
		Stage: func(p *core.AnimatedValue) core.Node {
			swing := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.4, 0.55, 1}, Output: []float64{0, -180, 180, 0}})
			second := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.42, 0.5, 1}, Output: []float64{0, 0, 1, 1}})
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
				Overflow: core.Hidden}, Children: []core.Node{
				camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					Animated: &core.AnimatedStyle{TranslateX: swing}}, Children: []core.Node{camShots(
					&core.AnimatedStyle{},
					&core.AnimatedStyle{Opacity: second},
				)}},
			}}
		}})
}

// #044 Cross dissolve: the next picture rises over this one.
func camDissolve(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "044", Title: "Cross dissolve", Duration: 2200 * ms, Accent: colViolet, Height: 110,
		Note: "Both pictures are on screen at once; the second one simply climbs over the first and the first never cuts.",
		Stage: func(p *core.AnimatedValue) core.Node {
			next := core.Interpolate(p, core.Interpolation{Input: []float64{0.15, 0.8}, Output: []float64{0, 1}})
			return core.View{Style: stageFill, Children: []core.Node{camShots(&core.AnimatedStyle{}, &core.AnimatedStyle{Opacity: next})}}
		}})
}

// #045 Wipe: an edge travels across and the next picture is behind it.
func camWipe(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "045", Title: "Wipe", Duration: 2200 * ms, Accent: colViolet, Height: 110,
		Note: "A hard edge, not a blend: the first shot slides off to the right and the second is already underneath.",
		Stage: func(p *core.AnimatedValue) core.Node {
			slide := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.75, 1}, Output: []float64{0, 240, 260}})
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0), Overflow: core.Hidden},
				Children: []core.Node{
					camScene(core.RGB(0x5a2b3c), core.RGB(0x2a1420), camLayers{}),
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Bottom: core.Pt(0), Width: core.Pt(260),
						Animated: &core.AnimatedStyle{TranslateX: slide}}, Children: []core.Node{
						camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{}),
					}},
				}}
		}})
}

// #046 Smash cut: quiet to loud with no warning.
func camSmashCut(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "046", Title: "Smash cut", Duration: 2000 * ms, Accent: colPink, Height: 110,
		Note: "No dissolve, no ramp: the cut lands on a frame boundary and the level meter jumps with it.",
		Stage: func(p *core.AnimatedValue) core.Node {
			loud := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.55, 0.56, 1}, Output: []float64{0, 0, 1, 1}})
			level := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.55, 0.56, 1}, Output: []float64{0.16, 0.16, 1, 1}})
			return core.View{Style: stageFill, Children: []core.Node{
				camShots(&core.AnimatedStyle{}, &core.AnimatedStyle{Opacity: loud}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(88), Left: core.Pt(12),
					Animated: &core.AnimatedStyle{}},
					Children: []core.Node{core.View{Style: core.Style{Width: core.Pt(140), Height: core.Pt(4), CornerRadius: 2,
						Background: core.RGB(0xff6b9d), Animated: &core.AnimatedStyle{ScaleX: level}}}}},
			}}
		}})
}

// #047 J-cut: the next scene's sound arrives before its picture.
func camJCut(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "047", Title: "J-cut", Duration: 2400 * ms, Accent: colBlue, Height: 110,
		Note: "Two timelines: the second sound track starts under the first picture, and only then does the picture follow it.",
		Stage: func(p *core.AnimatedValue) core.Node {
			audio := core.Interpolate(p, core.Interpolation{Input: []float64{0.2, 0.5}, Output: []float64{0, 1}})
			picture := core.Interpolate(p, core.Interpolation{Input: []float64{0.5, 0.85}, Output: []float64{0, 1}})
			return core.View{Style: stageFill, Children: []core.Node{
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Height: core.Pt(86)},
					Children: []core.Node{camShots(&core.AnimatedStyle{}, &core.AnimatedStyle{Opacity: picture})}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(90), Left: core.Pt(12), Right: core.Pt(12), Height: core.Pt(14),
					Direction: core.Row, Gap: 4, AlignItems: core.AlignCenter}, Children: []core.Node{
					uipLabel("A", 8, colFaint),
					core.View{Style: core.Style{Width: core.Pt(70), Height: core.Pt(4), CornerRadius: 2, Background: core.RGB(0x35d6de)}},
					uipLabel("B", 8, colFaint),
					core.View{Style: core.Style{Width: core.Pt(70), Height: core.Pt(4), CornerRadius: 2, Background: core.RGB(0xa78bfa),
						Animated: &core.AnimatedStyle{Opacity: audio}}},
				}},
			}}
		}})
}

// #048 L-cut: the old sound carries on past the new picture.
func camLCut(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "048", Title: "L-cut", Duration: 2400 * ms, Accent: colBlue, Height: 110,
		Note: "The picture changes first; the first scene's sound runs on underneath for another half second, then hands over.",
		Stage: func(p *core.AnimatedValue) core.Node {
			picture := core.Interpolate(p, core.Interpolation{Input: []float64{0.15, 0.3}, Output: []float64{0, 1}})
			old := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.62, 0.95}, Output: []float64{1, 1, 0}})
			return core.View{Style: stageFill, Children: []core.Node{
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Height: core.Pt(86)},
					Children: []core.Node{camShots(&core.AnimatedStyle{}, &core.AnimatedStyle{Opacity: picture})}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(90), Left: core.Pt(12), Right: core.Pt(12), Height: core.Pt(14),
					Direction: core.Row, Gap: 4, AlignItems: core.AlignCenter}, Children: []core.Node{
					uipLabel("A", 8, colFaint),
					core.View{Style: core.Style{Width: core.Pt(70), Height: core.Pt(4), CornerRadius: 2, Background: core.RGB(0x35d6de),
						Animated: &core.AnimatedStyle{Opacity: old}}},
					uipLabel("B", 8, colFaint),
					core.View{Style: core.Style{Width: core.Pt(70), Height: core.Pt(4), CornerRadius: 2, Background: core.RGB(0xa78bfa),
						Animated: &core.AnimatedStyle{Opacity: core.Interpolate(p, core.Interpolation{Input: []float64{0.62, 0.95}, Output: []float64{0, 1}})}}},
				}},
			}}
		}})
}
