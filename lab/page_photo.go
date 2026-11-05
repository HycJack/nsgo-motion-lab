package main

// page_photo.go is the photo pack (motion-prompt #097–#120): one portrait,
// lit and graded twenty-four ways.
//
// Everything here is honest about what NSGo can do: there is no blur and no
// blend mode, so a soft light is a scrim whose opacity moves, and a blurred
// background is a few low-opacity offset copies.
//
// Catalog names and timings: github.com/showheyohtaki/motion-prompt (MIT).

import (
	"github.com/oneo-me/nsgo/core"
)

// phScene is a horizon, a sun and a bust. The head and the body are separate
// so a scrim can darken one side and not the other.
func phScene(sky, ground core.Color, sun core.AnimatedStyle, l phLayers) core.Node {
	return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
		Overflow: core.Hidden}, Children: []core.Node{
		core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Height: core.Pt(84),
			Background: sky, Animated: l.Sky}},
		core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(84), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
			Background: ground, Animated: l.Ground}},
		core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(16), Left: core.Pt(30), Animated: &sun},
			Children: []core.Node{Disc(20, core.RGB(0xffd27a), nil)}},
		core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(30), Left: core.Pt(core.Points(52)),
			AlignItems: core.AlignCenter, Animated: l.Bust}, Children: []core.Node{
			core.View{Style: core.Style{Width: core.Pt(30), Height: core.Pt(32), CornerRadius: 15, Background: core.RGB(0xd9b79a),
				Animated: l.Head}},
			core.View{Style: core.Style{Width: core.Pt(64), Height: core.Pt(52), CornerRadius: 22, Background: core.RGB(0x3a4a5e),
				Animated: l.Body, Margin: core.Edges{Top: 2}}},
		}},
	}}
}

// solid wraps a plain colour where a colour node is expected: an animated
// style may override a background with a colour, but only as a node.
func solid(cx *core.Scope, c core.Color) core.AnimatedNode {
	return core.Interpolate(constant(cx, 0), core.Interpolation{
		Input: []float64{0, 1}, Output: []float64{0, 1}, Colors: []core.Color{c, c},
	})
}

// phLayers are the parts a tile may grade separately.
type phLayers struct {
	Sky, Ground, Head, Body, Bust *core.AnimatedStyle
}

// phPage renders the page with a credit and the grid.
func pagePhoto(cx *core.Scope) core.Node {
	return core.View{Children: []core.Node{
		core.View{Style: core.Style{Padding: core.Edges{Top: 14, Right: 18, Left: 18}}, Children: []core.Node{
			note("Photo pack 097–120 of motion-prompt (github.com/showheyohtaki/motion-prompt): one portrait, lit and graded. Code is NSGo's own; no blur, so softness is a scrim."),
		}},
		grid(
			phBacklight(cx),
			phRim(cx),
			phSilhouette(cx),
			phSplit(cx),
			phRembrandt(cx),
			phButterfly(cx),
			phGolden(cx),
			phBlueHour(cx),
			phBokeh(cx),
			phTiltShift(cx),
			phLongExposure(cx),
			phPanning(cx),
			phFlare(cx),
			phFisheye(cx),
			phAberration(cx),
			phVignette(cx),
			phHighKey(cx),
			phLowKey(cx),
			phMono(cx),
			phDuotone(cx),
			phTealOrange(cx),
			phGrain(cx),
			phLightLeak(cx),
			phDouble(cx),
		),
	}}
}

// phScrim is a full-frame wash: opacity and colour, over everything.
func phScrim(col core.Color, opacity core.AnimatedNode) core.Node {
	return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
		Background: col, Animated: &core.AnimatedStyle{Opacity: opacity}}}
}

// #097 Backlighting: the sun goes behind the head and the front goes dark.
func phBacklight(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "097", Title: "Backlighting", Duration: 2400 * ms, Accent: colAmber, Height: 110,
		Note: "The sun swings behind the head: an edge appears all the way round the figure, a warm haze rises, the front turns to shadow.",
		Stage: func(p *core.AnimatedValue) core.Node {
			behind := core.Interpolate(p, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1}, Easing: bezier(core.EaseInOut)})
			return core.View{Style: stageFill, Children: []core.Node{
				phScene(core.RGB(0x2a3f5c), core.RGB(0x1b2a3c), core.AnimatedStyle{TranslateX: lerp(behind, 60, 74), TranslateY: lerp(behind, 20, 26)},
					phLayers{Head: &core.AnimatedStyle{Background: core.Interpolate(behind, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1},
						Colors: []core.Color{core.RGB(0xd9b79a), core.RGB(0x6a5546)}}),
						}}),
				phScrim(core.RGB(0xffb070), core.Interpolate(behind, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 0.28}})),
			}}
		}})
}

// #098 Rim light: two lights behind, two thin edges, everything else dark.
func phRim(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "098", Title: "Rim light", Duration: 2200 * ms, Accent: colCyan, Height: 110,
		Note: "Two lights, both behind: the edges of the head and shoulders pick up a line and nothing else is lit.",
		Stage: func(p *core.AnimatedValue) core.Node {
			on := core.Interpolate(p, core.Interpolation{Input: []float64{0.15, 0.85}, Output: []float64{0, 1}})
			edge := func(left float32, col core.Color) core.Node {
				return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(26), Left: core.Pt(core.Points(left)),
					Width: core.Pt(3), Height: core.Pt(74), CornerRadius: 2, Background: col,
					Animated: &core.AnimatedStyle{Opacity: on, ScaleY: core.Interpolate(on, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0.4, 1}})}}}
			}
			return core.View{Style: stageFill, Children: []core.Node{
				phScene(core.RGB(0x0c1219), core.RGB(0x070b10), core.AnimatedStyle{Opacity: constant(cx, 0)},
					phLayers{Head: &core.AnimatedStyle{Background: solid(cx, core.RGB(0x2a2118))}, Body: &core.AnimatedStyle{Background: solid(cx, core.RGB(0x1a2230))}}),
				edge(50, core.RGB(0xd8f2ff)),
				edge(80, core.RGB(0xd8f2ff)),
			}}
		}})
}

// #099 Silhouette: exposure set for the sky, so the figure is a hole.
func phSilhouette(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "099", Title: "Silhouette", Duration: 2200 * ms, Accent: colAmber, Height: 110,
		Note: "Nothing about the figure changes: the exposure does, and the person becomes the black shape the sky leaves.",
		Stage: func(p *core.AnimatedValue) core.Node {
			black := core.Interpolate(p, core.Interpolation{Input: []float64{0.1, 0.9}, Output: []float64{0, 1}})
			return phScene(core.RGB(0xffb066), core.RGB(0xd98a4a), core.AnimatedStyle{},
				phLayers{Head: &core.AnimatedStyle{Background: core.Interpolate(black, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1},
					Colors: []core.Color{core.RGB(0xd9b79a), core.RGB(0x0a0c10)}})},
					Body: &core.AnimatedStyle{Background: core.Interpolate(black, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1},
						Colors: []core.Color{core.RGB(0x3a4a5e), core.RGB(0x0a0c10)}})}})
		}})
}

// #100 Split lighting: the light swings round until the shadow halves the face.
func phSplit(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "100", Title: "Split lighting", Duration: 2200 * ms, Accent: colAmber, Height: 110,
		Note: "Twenty degrees to ninety: the shadow on the head grows from a crescent until it is exactly half the face.",
		Stage: func(p *core.AnimatedValue) core.Node {
			shadow := core.Interpolate(p, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0.08, 0.62}, Easing: bezier(core.EaseInOut)})
			return core.View{Style: stageFill, Children: []core.Node{
				phScene(core.RGB(0x1b2534), core.RGB(0x141c27), core.AnimatedStyle{TranslateX: lerp(p, 44, 8)},
					phLayers{Head: &core.AnimatedStyle{Background: solid(cx, core.RGB(0xe0bd9d))}}),
				// A shadow that grows across the face is a panel whose width
				// grows from the shadow side.
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(30), Left: core.Pt(52), Height: core.Pt(32),
					Width: core.Pt(30), Background: core.RGB(0x241a14),
					Animated: &core.AnimatedStyle{ScaleX: shadow}}},
			}}
		}})
}

// #101 Rembrandt: 45 up and 45 round leaves one triangle on the far cheek.
func phRembrandt(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "101", Title: "Rembrandt lighting", Duration: 2200 * ms, Accent: colAmber, Height: 110,
		Note: "The whole effect is that small inverted triangle; without it a Rembrandt is just a dark portrait.",
		Stage: func(p *core.AnimatedValue) core.Node {
			on := core.Interpolate(p, core.Interpolation{Input: []float64{0.2, 0.8}, Output: []float64{0, 1}})
			return core.View{Style: stageFill, Children: []core.Node{
				phScene(core.RGB(0x161d28), core.RGB(0x101620), core.AnimatedStyle{TranslateX: lerp(p, 24, 18), TranslateY: lerp(p, 30, 10)},
					phLayers{Head: &core.AnimatedStyle{Background: solid(cx, core.RGB(0x3a2c22))}}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(52), Left: core.Pt(74),
					Width: core.Pt(9), Height: core.Pt(7), Background: core.RGB(0xe8c49a),
					Animated: &core.AnimatedStyle{Opacity: on}}},
			}}
		}})
}

// #102 Butterfly: light high and in front, so only the chin has a shadow.
func phButterfly(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "102", Title: "Butterfly lighting", Duration: 2000 * ms, Accent: colAmber, Height: 110,
		Note: "Both sides of the face light evenly and one small shadow sits under the chin — the least dramatic of the three.",
		Stage: func(p *core.AnimatedValue) core.Node {
			on := core.Interpolate(p, core.Interpolation{Input: []float64{0.15, 0.7}, Output: []float64{0, 1}})
			return core.View{Style: stageFill, Children: []core.Node{
				phScene(core.RGB(0x26313f), core.RGB(0x1a222c), core.AnimatedStyle{TranslateY: lerp(p, 6, 2)},
					phLayers{Head: &core.AnimatedStyle{Background: solid(cx, core.RGB(0xe6c4a4))}}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(62), Left: core.Pt(68),
					Width: core.Pt(22), Height: core.Pt(5), CornerRadius: 3, Background: core.RGB(0x3a2a1e),
					Animated: &core.AnimatedStyle{Opacity: on}}},
			}}
		}})
}

// #103 Golden hour: the sun reaches the horizon and the shadows stretch.
func phGolden(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "103", Title: "Golden hour", Duration: 2600 * ms, Accent: colAmber, Height: 110,
		Note: "The sun drops to the horizon, the sky and the ground turn orange and the shadows lie down to the right.",
		Stage: func(p *core.AnimatedValue) core.Node {
			low := core.Interpolate(p, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1}, Easing: bezier(core.EaseInOut)})
			return core.View{Style: stageFill, Children: []core.Node{
				phScene(core.RGB(0x2a3f5c), core.RGB(0x1b2a3c), core.AnimatedStyle{TranslateY: lerp(low, 10, 46), TranslateX: lerp(low, 20, 6)},
					phLayers{
						Sky: &core.AnimatedStyle{Background: core.Interpolate(low, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1},
							Colors: []core.Color{core.RGB(0x2a3f5c), core.RGB(0xe8873c)}})},
						Ground: &core.AnimatedStyle{Background: core.Interpolate(low, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1},
							Colors: []core.Color{core.RGB(0x1b2a3c), core.RGB(0xb05f2a)}})},
						Head: &core.AnimatedStyle{Background: core.Interpolate(low, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1},
							Colors: []core.Color{core.RGB(0xd9b79a), core.RGB(0xf0b077)}})},
					}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(96), Left: core.Pt(72), Height: core.Pt(6),
					Width: core.Pt(90), Background: core.RGB(0x6d3a18),
					Animated: &core.AnimatedStyle{ScaleX: core.Interpolate(low, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0.15, 1}})}}},
			}}
		}})
}

// #104 Blue hour: two bands of deep blue and the windows come on amber.
func phBlueHour(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "104", Title: "Blue hour", Duration: 2400 * ms, Accent: colBlue, Height: 110,
		Note: "After sunset the sky splits into two blues and the lighting turns from the sun to the windows.",
		Stage: func(p *core.AnimatedValue) core.Node {
			dusk := core.Interpolate(p, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1}, Easing: bezier(core.EaseInOut)})
			win := core.Interpolate(p, core.Interpolation{Input: []float64{0.45, 1}, Output: []float64{0, 1}})
			var windows []core.Node
			for i := 0; i < 5; i++ {
				windows = append(windows, core.View{Style: core.Style{Position: core.Absolute,
					Top: core.Pt(core.Points(float32(74 + (i%2)*12))), Left: core.Pt(core.Points(float32(8 + i*34))),
					Width: core.Pt(9), Height: core.Pt(7), CornerRadius: 1, Background: core.RGB(0xffb347),
					Animated: &core.AnimatedStyle{Opacity: win}}})
			}
			return core.View{Style: stageFill, Children: append([]core.Node{
				phScene(core.RGB(0x16243a), core.RGB(0x0e1826), core.AnimatedStyle{Opacity: core.Interpolate(dusk, core.Interpolation{Input: []float64{0, 1}, Output: []float64{1, 0.1}})},
					phLayers{
						Sky: &core.AnimatedStyle{Background: core.Interpolate(dusk, core.Interpolation{Input: []float64{0, 0.5, 1}, Output: []float64{0, 1, 2},
							Colors: []core.Color{core.RGB(0x2a3f5c), core.RGB(0x1b3355), core.RGB(0x101f38)}})},
						Head: &core.AnimatedStyle{Background: core.Interpolate(dusk, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1},
							Colors: []core.Color{core.RGB(0xd9b79a), core.RGB(0x2a3a55)}})},
					}),
			}, windows...)}
		}})
}

// #105 Bokeh: the background lights open into discs while the figure holds.
func phBokeh(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "105", Title: "Bokeh", Duration: 2200 * ms, Accent: colAmber, Height: 110,
		Note: "The aperture opens: the lights behind spread into big pale discs and everything back there loses its edge.",
		Stage: func(p *core.AnimatedValue) core.Node {
			var discs []core.Node
			for i := 0; i < 6; i++ {
				delay := 0.08 * float64(i)
				q := core.Interpolate(p, core.Interpolation{Input: []float64{delay, 0.85}, Output: []float64{0, 1}, Easing: bezier(core.EaseInOut)})
				discs = append(discs, core.View{Style: core.Style{Position: core.Absolute,
					Top: core.Pt(core.Points(float32(18 + (i%3)*22))), Left: core.Pt(core.Points(float32(10 + i*38))),
					Animated: &core.AnimatedStyle{Scale: core.Interpolate(q, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0.3, 2.2}}),
						Opacity: core.Interpolate(q, core.Interpolation{Input: []float64{0, 0.6, 1}, Output: []float64{0, 0.55, 0.35}})}},
					Children: []core.Node{Disc(16, core.RGB(0xffe3ad), nil)}})
			}
			return core.View{Style: stageFill, Children: append([]core.Node{
				phScene(core.RGB(0x1d2a3a), core.RGB(0x141c27), core.AnimatedStyle{},
					phLayers{Head: &core.AnimatedStyle{Background: solid(cx, core.RGB(0xe2c1a2))}}),
			}, discs...)}
		}})
}

// #106 Tilt-shift: two bands of softness, one band of focus.
func phTiltShift(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "106", Title: "Tilt-shift", Duration: 2200 * ms, Accent: colCyan, Height: 110,
		Note: "Only the middle band stays sharp; the top and the bottom soften, and that is what makes a street look like a model.",
		Stage: func(p *core.AnimatedValue) core.Node {
			soft := core.Interpolate(p, core.Interpolation{Input: []float64{0.1, 0.9}, Output: []float64{0, 0.55}})
			band := func(top float32, h float32) core.Node {
				return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(core.Points(top)), Left: core.Pt(0), Right: core.Pt(0),
					Height: core.Pt(core.Points(h)), Background: core.RGB(0x0e141d),
					Animated: &core.AnimatedStyle{Opacity: soft}}}
			}
			return core.View{Style: stageFill, Children: []core.Node{
				phScene(core.RGB(0x2a3f5c), core.RGB(0x1b2a3c), core.AnimatedStyle{}, phLayers{}),
				band(0, 38), band(78, 32),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(37), Left: core.Pt(0), Right: core.Pt(0), Height: core.Pt(1),
					Background: core.RGB(0x9fd8ff), Animated: &core.AnimatedStyle{Opacity: core.Interpolate(soft, core.Interpolation{Input: []float64{0, 0.55}, Output: []float64{0, 0.5}})}}},
			}}
		}})
}

// #107 Long exposure: the car goes, the lights stay.
func phLongExposure(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "107", Title: "Long exposure", Duration: 2400 * ms, Accent: colPink, Height: 110,
		Note: "The shutter stays open: the car dissolves while the two lines of light it drew stay behind it.",
		Stage: func(p *core.AnimatedValue) core.Node {
			night := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.5}, Output: []float64{0, 1}, Easing: bezier(core.EaseInOut)})
			trail := core.Interpolate(p, core.Interpolation{Input: []float64{0.25, 0.9}, Output: []float64{0, 1}})
			return core.View{Style: stageFill, Children: []core.Node{
				phScene(core.RGB(0x1b2a3c), core.RGB(0x121a24), core.AnimatedStyle{},
					phLayers{Sky: &core.AnimatedStyle{Background: core.Interpolate(night, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1},
						Colors: []core.Color{core.RGB(0x2a3f5c), core.RGB(0x0a1424)}})}}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(86), Left: core.Pt(30),
					Animated: &core.AnimatedStyle{TranslateX: lerp(p, 0, 120), Opacity: core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.55}, Output: []float64{1, 0.25}})}},
					Children: []core.Node{Panel(26, 12, core.RGB(0x2f3a48), nil)}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(84), Left: core.Pt(30), Width: core.Pt(150), Height: core.Pt(2),
					Background: core.RGB(0xff5a4d), Animated: &core.AnimatedStyle{ScaleX: trail}}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(94), Left: core.Pt(30), Width: core.Pt(150), Height: core.Pt(2),
					Background: core.RGB(0xfff0c0), Animated: &core.AnimatedStyle{ScaleX: trail}}},
			}}
		}})
}

// #108 Panning: the camera goes with the car, so only the background streaks.
func phPanning(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "108", Title: "Panning shot", Duration: 2400 * ms, Accent: colAmber, Height: 110,
		Note: "The car does not move on screen at all; the only thing that does is the world behind it, sideways.",
		Stage: func(p *core.AnimatedValue) core.Node {
			streak := core.Interpolate(p, core.Interpolation{Input: []float64{0, 1}, Output: []float64{60, -60}})
			var lines []core.Node
			for i := 0; i < 5; i++ {
				lines = append(lines, core.View{Style: core.Style{Position: core.Absolute,
					Top: core.Pt(core.Points(float32(22 + i*14))), Left: core.Pt(0), Width: core.Pt(90), Height: core.Pt(2),
					Background: core.RGB(0x24303f),
					Animated: &core.AnimatedStyle{TranslateX: core.Interpolate(streak, core.Interpolation{Input: []float64{60, -60}, Output: []float64{60 * float64(i+1), -60 * float64(i+1)}})}}})
			}
			return core.View{Style: stageFill, Children: append([]core.Node{
				phScene(core.RGB(0x2a3f5c), core.RGB(0x1b2a3c), core.AnimatedStyle{TranslateX: streak}, phLayers{}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(74), Left: core.Pt(96)},
					Children: []core.Node{Panel(40, 16, core.RGB(0xd94b3f), nil)}},
			}, lines...)}
		}})
}

// #109 Lens flare: rays around the sun and a line of ghosts across the frame.
func phFlare(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "109", Title: "Lens flare", Duration: 2200 * ms, Accent: colAmber, Height: 110,
		Note: "The flare is not light on the scene: it is the lens answering the sun, so it sits in front of everything.",
		Stage: func(p *core.AnimatedValue) core.Node {
			on := core.Interpolate(p, core.Interpolation{Input: []float64{0.15, 0.6, 1}, Output: []float64{0, 1, 0.7}})
			var rays, ghosts []core.Node
			for i := 0; i < 6; i++ {
				rays = append(rays, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(28), Left: core.Pt(46),
					Width: core.Pt(90), Height: core.Pt(1), Background: core.RGB(0xffe9b8),
					Animated: &core.AnimatedStyle{Opacity: on, Rotate: constant(cx, 30.0*float64(i)),
						ScaleX: core.Interpolate(on, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0.2, 1}})}}})
			}
			for i := 0; i < 5; i++ {
				ghosts = append(ghosts, core.View{Style: core.Style{Position: core.Absolute,
					Top: core.Pt(core.Points(float32(24 + 12*i))), Left: core.Pt(core.Points(float32(52 + 26*i))),
					Animated: &core.AnimatedStyle{Opacity: core.Interpolate(on, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 0.5}})}},
					Children: []core.Node{Disc(6+float32(i)*3, core.RGB(0x8fd6ff), nil)}})
			}
			return core.View{Style: stageFill, Children: append(append([]core.Node{
				phScene(core.RGB(0x2a3f5c), core.RGB(0x1b2a3c), core.AnimatedStyle{TranslateX: lerp(p, 40, 26)},
					phLayers{Head: &core.AnimatedStyle{Background: solid(cx, core.RGB(0x5c4a3a))}}),
			}, rays...), ghosts...)}
		}})
}

// #110 Fisheye: the middle bulges and the edges bend outward.
func phFisheye(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "110", Title: "Fisheye lens", Duration: 2200 * ms, Accent: colCyan, Height: 110,
		Note: "An 8 mm lens is a bulge we cannot do in geometry, so the scene overfills a rounded frame and the corners fall away.",
		Stage: func(p *core.AnimatedValue) core.Node {
			bulge := core.Interpolate(p, core.Interpolation{Input: []float64{0, 1}, Output: []float64{1, 1.34}, Easing: bezier(core.EaseInOut)})
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
				AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter}, Children: []core.Node{
				core.View{Style: core.Style{Width: core.Pt(250), Height: core.Pt(130), CornerRadius: 62, Overflow: core.Hidden,
					Animated: &core.AnimatedStyle{Scale: bulge}}, Children: []core.Node{
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(-10), Left: core.Pt(-6), Width: core.Pt(262), Height: core.Pt(130)},
						Children: []core.Node{phScene(core.RGB(0x2a3f5c), core.RGB(0x1b2a3c), core.AnimatedStyle{}, phLayers{})}},
				}},
			}}
		}})
}

// #111 Chromatic aberration: three copies, three sizes, fringes at the edges.
func phAberration(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "111", Title: "Chromatic aberration", Duration: 2000 * ms, Accent: colViolet, Height: 110,
		Note: "Red, green and blue are laid back a little out of register: the fringe is biggest at the edges, so the centre stays clean.",
		Stage: func(p *core.AnimatedValue) core.Node {
			off := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.6, 1}, Output: []float64{0, 1, 0.6}})
			layer := func(tint core.Color, dx float64, dy float64) core.Node {
				return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					Background: tint, Animated: &core.AnimatedStyle{
						TranslateX: core.Interpolate(off, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, dx}}),
						TranslateY: core.Interpolate(off, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, dy}}),
						Opacity:    core.Interpolate(off, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0.0, 0.34}})}},
					Children: []core.Node{phScene(tint, tint, core.AnimatedStyle{}, phLayers{})}}
			}
			return core.View{Style: stageFill, Children: []core.Node{
				phScene(core.RGB(0x2a3f5c), core.RGB(0x1b2a3c), core.AnimatedStyle{}, phLayers{}),
				layer(core.RGB(0xff0000), -4, 0),
				layer(core.RGB(0x00ff00), 4, 0),
				layer(core.RGB(0x0000ff), 0, 4),
			}}
		}})
}

// #112 Vignette: the corners fall away and the middle stays.
func phVignette(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "112", Title: "Vignette", Duration: 2000 * ms, Accent: colFaint, Height: 110,
		Note: "Four corners, one wash: the falloff is what pulls the eye to the face, not the darkness itself.",
		Stage: func(p *core.AnimatedValue) core.Node {
			dark := core.Interpolate(p, core.Interpolation{Input: []float64{0.1, 0.95}, Output: []float64{0, 0.85}})
			corner := func(top float32, left float32) core.Node {
				return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(core.Points(top)), Left: core.Pt(core.Points(left)),
					Width: core.Pt(110), Height: core.Pt(70), CornerRadius: 30, Background: core.RGB(0x05070a),
					Animated: &core.AnimatedStyle{Opacity: dark}}}
			}
			return core.View{Style: stageFill, Children: []core.Node{
				phScene(core.RGB(0x2a3f5c), core.RGB(0x1b2a3c), core.AnimatedStyle{},
					phLayers{Head: &core.AnimatedStyle{Background: solid(cx, core.RGB(0xf0d0b0))}}),
				corner(-18, -22), corner(-18, 120), corner(70, -22), corner(70, 120),
			}}
		}})
}

// #113 High-key: everything goes light, the shadows go last.
func phHighKey(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "113", Title: "High-key", Duration: 2000 * ms, Accent: colText, Height: 110,
		Note: "Not an exposure to white: the tones move up together, so the face is still readable when the frame is nearly white.",
		Stage: func(p *core.AnimatedValue) core.Node {
			up := core.Interpolate(p, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1}, Easing: bezier(core.EaseInOut)})
			return core.View{Style: stageFill, Children: []core.Node{
				phScene(core.RGB(0x2a3f5c), core.RGB(0x1b2a3c), core.AnimatedStyle{},
					phLayers{
						Sky: &core.AnimatedStyle{Background: core.Interpolate(up, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1},
							Colors: []core.Color{core.RGB(0x2a3f5c), core.RGB(0xdce9f4)}})},
						Ground: &core.AnimatedStyle{Background: core.Interpolate(up, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1},
							Colors: []core.Color{core.RGB(0x1b2a3c), core.RGB(0xc8d8e6)}})},
						Head: &core.AnimatedStyle{Background: core.Interpolate(up, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1},
							Colors: []core.Color{core.RGB(0xd9b79a), core.RGB(0xf7e7d8)}})},
						Body: &core.AnimatedStyle{Background: core.Interpolate(up, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1},
							Colors: []core.Color{core.RGB(0x3a4a5e), core.RGB(0xbccedd)}})},
					}),
				phScrim(core.RGB(0xffffff), core.Interpolate(up, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 0.22}})),
			}}
		}})
}

// #114 Low-key: only the lit edge survives.
func phLowKey(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "114", Title: "Low-key", Duration: 2200 * ms, Accent: colAmber, Height: 110,
		Note: "Ninety percent of the frame goes, and the ten percent left is a single lit edge on the right of the bust.",
		Stage: func(p *core.AnimatedValue) core.Node {
			down := core.Interpolate(p, core.Interpolation{Input: []float64{0.1, 0.95}, Output: []float64{0, 1}, Easing: bezier(core.EaseInOut)})
			return core.View{Style: stageFill, Children: []core.Node{
				phScene(core.RGB(0x1b2634), core.RGB(0x111823), core.AnimatedStyle{Opacity: core.Interpolate(down, core.Interpolation{Input: []float64{0, 1}, Output: []float64{1, 0.25}})},
					phLayers{Head: &core.AnimatedStyle{Background: solid(cx, core.RGB(0x241c16))}}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(30), Left: core.Pt(82), Width: core.Pt(3), Height: core.Pt(78),
					CornerRadius: 2, Background: core.RGB(0xffc98a),
					Animated: &core.AnimatedStyle{Opacity: down, ScaleY: core.Interpolate(down, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0.2, 1}})}}},
				phScrim(core.RGB(0x05070a), core.Interpolate(down, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 0.66}})),
			}}
		}})
}

// #115 Black and white: the colour drains, the shapes stay.
func phMono(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "115", Title: "Black and white", Duration: 2000 * ms, Accent: colText, Height: 110,
		Note: "The sun and the skin lose their colour and only brightness is left to tell them apart — which is the whole lesson of the effect.",
		Stage: func(p *core.AnimatedValue) core.Node {
			gray := core.Interpolate(p, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1}, Easing: bezier(core.EaseInOut)})
			mix := func(a, b core.Color) core.AnimatedNode {
				return core.Interpolate(gray, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1}, Colors: []core.Color{a, b}})
			}
			return phScene(core.RGB(0x2a3f5c), core.RGB(0x1b2a3c), core.AnimatedStyle{},
				phLayers{
					Sky:    &core.AnimatedStyle{Background: mix(core.RGB(0x2a3f5c), core.RGB(0x6d7a88))},
					Ground: &core.AnimatedStyle{Background: mix(core.RGB(0x1b2a3c), core.RGB(0x4a545f))},
					Head:   &core.AnimatedStyle{Background: mix(core.RGB(0xd9b79a), core.RGB(0xb9b9b9))},
					Body:   &core.AnimatedStyle{Background: mix(core.RGB(0x3a4a5e), core.RGB(0x555555))},
				})
		}})
}

// #116 Duotone: two inks, one for the darks and one for the lights.
func phDuotone(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "116", Title: "Duotone", Duration: 2200 * ms, Accent: colViolet, Height: 110,
		Note: "Not a tint: every tone is remapped into a range between two colours, which is why the picture reads as a poster.",
		Stage: func(p *core.AnimatedValue) core.Node {
			mix := core.Interpolate(p, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1}, Easing: bezier(core.EaseInOut)})
			tone := func(a, b core.Color) core.AnimatedNode {
				return core.Interpolate(mix, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1}, Colors: []core.Color{a, b}})
			}
			return phScene(core.RGB(0x2a3f5c), core.RGB(0x1b2a3c), core.AnimatedStyle{},
				phLayers{
					Sky:    &core.AnimatedStyle{Background: tone(core.RGB(0x2a3f5c), core.RGB(0x2b2a6b))},
					Ground: &core.AnimatedStyle{Background: tone(core.RGB(0x1b2a3c), core.RGB(0x1a1948))},
					Head:   &core.AnimatedStyle{Background: tone(core.RGB(0xd9b79a), core.RGB(0xff9ec4))},
					Body:   &core.AnimatedStyle{Background: tone(core.RGB(0x3a4a5e), core.RGB(0x4e3a8a))},
				})
		}})
}

// #117 Teal and orange: the complement is the grade.
func phTealOrange(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "117", Title: "Teal and orange", Duration: 2400 * ms, Accent: colCyan, Height: 110,
		Note: "Two colours pulling opposite ways: the sky and the shadows go teal, the sun and the lit edge go orange.",
		Stage: func(p *core.AnimatedValue) core.Node {
			grade := core.Interpolate(p, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1}, Easing: bezier(core.EaseInOut)})
			tone := func(a, b core.Color) core.AnimatedNode {
				return core.Interpolate(grade, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1}, Colors: []core.Color{a, b}})
			}
			return core.View{Style: stageFill, Children: []core.Node{
				phScene(core.RGB(0x2a3f5c), core.RGB(0x1b2a3c), core.AnimatedStyle{TranslateX: lerp(p, 40, 16)},
					phLayers{
						Sky:    &core.AnimatedStyle{Background: tone(core.RGB(0x2a3f5c), core.RGB(0x146b70))},
						Ground: &core.AnimatedStyle{Background: tone(core.RGB(0x1b2a3c), core.RGB(0x0c4a50))},
						Head:   &core.AnimatedStyle{Background: tone(core.RGB(0xd9b79a), core.RGB(0xe0854a))},
						Body:   &core.AnimatedStyle{Background: tone(core.RGB(0x3a4a5e), core.RGB(0x1d3a44))},
					}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(30), Left: core.Pt(82), Width: core.Pt(3), Height: core.Pt(40),
					CornerRadius: 2, Background: core.RGB(0xffa759), Animated: &core.AnimatedStyle{Opacity: grade}}},
			}}
		}})
}

// #118 Film grain: fine noise, warm, on top of everything.
func phGrain(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "118", Title: "Film grain", Duration: 1600 * ms, Accent: colAmber, Height: 110, Loop: true,
		Note: "Grain is a hundred small decisions a second: the dots move on every sample, and the whole picture leans warm.",
		Stage: func(p *core.AnimatedValue) core.Node {
			var specks []core.Node
			for i := 0; i < 26; i++ {
				top := (i * 37) % 100
				left := (i * 53) % 240
				flick := core.Interpolate(p, core.Interpolation{
					Input:  []float64{0, 0.2, 0.35, 0.5, 0.65, 0.8, 1},
					Output: []float64{0.15, 0.5, 0.1, 0.45, 0.2, 0.55, 0.15},
				})
				specks = append(specks, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(core.Points(float32(top))),
					Left: core.Pt(core.Points(float32(left))), Width: core.Pt(1), Height: core.Pt(1), Background: core.RGB(0xfff2d0),
					Animated: &core.AnimatedStyle{Opacity: flick}}})
			}
			return core.View{Style: stageFill, Children: append([]core.Node{
				phScene(core.RGB(0x2a3f5c), core.RGB(0x1b2a3c), core.AnimatedStyle{}, phLayers{}),
				phScrim(core.RGB(0xffd9a0), constant(cx, 0.1)),
			}, specks...)}
		}})
}

// #119 Light leak: a band of stray light enters from the right edge.
func phLightLeak(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "119", Title: "Light leak", Duration: 2400 * ms, Accent: colPink, Height: 110,
		Note: "A leak is not in the scene at all: it is light that got past the seal, so it lies over the picture and moves slowly.",
		Stage: func(p *core.AnimatedValue) core.Node {
			sway := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.5, 1}, Output: []float64{30, -10, 30}})
			strength := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.35, 0.7, 1}, Output: []float64{0.1, 0.7, 0.55, 0.2}})
			return core.View{Style: stageFill, Children: []core.Node{
				phScene(core.RGB(0x2a3f5c), core.RGB(0x1b2a3c), core.AnimatedStyle{}, phLayers{}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(-30), Left: core.Pt(200), Width: core.Pt(46), Height: core.Pt(180),
					Background: core.RGB(0xffe0a0), Animated: &core.AnimatedStyle{TranslateX: sway, Rotate: constant(cx, 12), Opacity: strength}}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(-30), Left: core.Pt(238), Width: core.Pt(8), Height: core.Pt(180),
					Background: core.RGB(0xff6a5a), Animated: &core.AnimatedStyle{TranslateX: sway, Rotate: constant(cx, 12), Opacity: strength}}},
			}}
		}})
}

// #120 Double exposure: a landscape, only where the bust is.
func phDouble(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "120", Title: "Double exposure", Duration: 2600 * ms, Accent: colCyan, Height: 110,
		Note: "The second frame is not over the first, it is inside it: the bust is the window, and the hills show through it.",
		Stage: func(p *core.AnimatedValue) core.Node {
			inside := core.Interpolate(p, core.Interpolation{Input: []float64{0.15, 0.9}, Output: []float64{0, 1}, Easing: bezier(core.EaseInOut)})
			return core.View{Style: stageFill, Children: []core.Node{
				phScene(core.RGB(0xf4f1ea), core.RGB(0xe8e2d6), core.AnimatedStyle{Opacity: constant(cx, 0)},
					phLayers{Head: &core.AnimatedStyle{Background: solid(cx, core.RGB(0x1a2230))}, Body: &core.AnimatedStyle{Background: solid(cx, core.RGB(0x1a2230))}}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(30), Left: core.Pt(core.Points(52)), Width: core.Pt(64), Height: core.Pt(80),
					Overflow: core.Hidden, AlignItems: core.AlignCenter, Animated: &core.AnimatedStyle{Opacity: inside}},
					Children: []core.Node{
						core.View{Style: core.Style{Width: core.Pt(30), Height: core.Pt(32), CornerRadius: 15, Overflow: core.Hidden}, Children: []core.Node{
							core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(-6), Left: core.Pt(-10), Width: core.Pt(50), Height: core.Pt(50)},
								Children: []core.Node{camScene(core.RGB(0x2a3f5c), core.RGB(0x1b2a3c), camLayers{})}},
						}},
						core.View{Style: core.Style{Width: core.Pt(64), Height: core.Pt(52), CornerRadius: 22, Overflow: core.Hidden}, Children: []core.Node{
							core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(-20), Left: core.Pt(-30), Width: core.Pt(130), Height: core.Pt(90)},
								Children: []core.Node{camScene(core.RGB(0x2a3f5c), core.RGB(0x1b2a3c), camLayers{})}},
						}},
					}},
			}}
		}})
}
