package main

// page_camera2.go is the camera2 pack (motion-prompt #169–#192): the moves a
// second camera makes, the angles, and the last six edit-room tricks.
//
// It reuses the landscape and the two-shot pair from page_camera.go, except
// for the three tiles that need a viewpoint rather than a picture: fly-through
// 176, POV 179 and over-the-shoulder 180 run on the ray-marched set from
// page_3d.go, camera-as-uniforms, at this pack's 110-point stage.
//
// Catalog names and timings: github.com/showheyohtaki/motion-prompt (MIT).

import (
	"github.com/oneo-me/nsgo/core"
)

// pageCamera2 lays the pack out.
func pageCamera2(cx *core.Scope) core.Node {
	return core.View{Children: []core.Node{
		core.View{Style: core.Style{Padding: core.Edges{Top: 14, Right: 18, Left: 18}}, Children: []core.Node{
			note("Camera2 pack 169–192 of motion-prompt (github.com/showheyohtaki/motion-prompt): pedestals, angles, and six more cuts. Code is NSGo's own. Three of these need a viewpoint rather than a picture — fly-through 176, POV 179, over-the-shoulder 180 — so they fly the ray-marched set from the 3D tab."),
		}},
		grid(
			c2Pedestal(cx),
			c2BarrelRoll(cx),
			c2Hyperlapse(cx),
			c2LongTake(cx),
			c2Steadicam(cx),
			c2Follow(cx),
			c2LockedOff(cx),
			c2FlyThrough(cx),
			c2LowAngle(cx),
			c2HighAngle(cx),
			c2POV(cx),
			c2OverShoulder(cx),
			c2BirdEye(cx),
			c2DeepFocus(cx),
			c2Establishing(cx),
			c2ExtremeClose(cx),
			c2Insert(cx),
			c2MorphCut(cx),
			c2FadeBlack(cx),
			c2IrisOut(cx),
			c2Cutaway(cx),
			c2CrossCutting(cx),
			c2BulletTime(cx),
			c2Montage(cx),
		),
	}}
}

// #169 Pedestal: the camera rises without tilting, so the horizon holds.
func c2Pedestal(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "169", Title: "Pedestal", Duration: 2400 * ms, Accent: colCyan, Height: 110,
		Note: "The camera goes straight up and keeps looking level: the horizon stays put and the figure slides down with the ground.",
		Stage: func(p *core.AnimatedValue) core.Node {
			return camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{
				Sun:    &core.AnimatedStyle{TranslateY: lerp(p, 0, 26)},
				Far:    &core.AnimatedStyle{TranslateY: lerp(p, 0, 34)},
				Near:   &core.AnimatedStyle{TranslateY: lerp(p, 0, 44)},
				Figure: &core.AnimatedStyle{TranslateY: lerp(p, 0, 52)},
			})
		}})
}

// #170 Barrel roll: a full turn around the lens axis, back to level.
func c2BarrelRoll(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "170", Title: "Barrel roll", Duration: 2400 * ms, Accent: colViolet, Height: 110,
		Note: "Three hundred and sixty degrees of roll, and the horizon is level at both ends: the frame is scaled up enough to hide the corners on the way round.",
		Stage: func(p *core.AnimatedValue) core.Node {
			roll := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.5, 1}, Output: []float64{0, 180, 360}, Easing: bezier(core.EaseInOut)})
			// Inset well past the frame and scale up hard: at the 90° point of the
			// roll the scene is only 110 tall, so the scale has to be large enough
			// that the short side still covers the card's width (no black corners).
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(-60), Left: core.Pt(-60), Right: core.Pt(-60), Bottom: core.Pt(-60),
				Animated: &core.AnimatedStyle{Rotate: roll, Scale: constant(cx, 2.5)}},
				Children: []core.Node{camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{})}}
		}})
}

// #171 Hyperlapse: a time-lapse that also travels.
func c2Hyperlapse(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "171", Title: "Hyperlapse", Duration: 2600 * ms, Accent: colAmber, Height: 110,
		Note: "Two speeds at once: the world races past in fast forward while the camera itself walks forward through it.",
		Stage: func(p *core.AnimatedValue) core.Node {
			sun := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.5, 1}, Output: []float64{46, -10, 46}})
			return camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{
				Sun:    &core.AnimatedStyle{TranslateX: lerp(p, 20, -60), TranslateY: sun},
				Far:    &core.AnimatedStyle{TranslateX: lerp(p, 60, -160)},
				Near:   &core.AnimatedStyle{TranslateX: lerp(p, 120, -320), Scale: lerp(p, 1, 1.3), TranslateY: lerp(p, 0, 10)},
				Figure: &core.AnimatedStyle{TranslateX: lerp(p, 40, 30), TranslateY: lerp(p, 0, 14), Scale: lerp(p, 1, 1.2)},
			})
		}})
}

// #172 Long take: one unbroken shot, and the timeline says so.
func c2LongTake(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "172", Title: "Long take", Duration: 2600 * ms, Accent: colGreen, Height: 110,
		Note: "The picture never cuts, so the bar underneath is one strip; a cut would put a gap in it.",
		Stage: func(p *core.AnimatedValue) core.Node {
			return core.View{Style: stageFill, Children: []core.Node{
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Height: core.Pt(86)},
					Children: []core.Node{camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{
						Figure: &core.AnimatedStyle{TranslateX: lerp(p, -30, 150)},
						Far:    &core.AnimatedStyle{TranslateX: lerp(p, 40, -40)},
						Near:   &core.AnimatedStyle{TranslateX: lerp(p, 80, -80)},
					})}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(94), Left: core.Pt(10), Right: core.Pt(10), Height: core.Pt(3),
					CornerRadius: 2, Background: core.RGB(0x39d98a)}},
			}}
		}})
}

// #173 Steadicam: a smooth glide, with one gentle sway — never a jitter.
func c2Steadicam(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "173", Title: "Steadicam", Duration: 2400 * ms, Accent: colGreen, Height: 110,
		Note: "The difference from handheld is the spectrum: one slow sway here, four uneven shakes there.",
		Stage: func(p *core.AnimatedValue) core.Node {
			sway := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.5, 1}, Output: []float64{-2, 2, -2}, Easing: bezier(core.EaseInOut)})
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
				Animated: &core.AnimatedStyle{TranslateY: sway, Rotate: core.Interpolate(sway, core.Interpolation{Input: []float64{-2, 2}, Output: []float64{-0.8, 0.8}})}},
				Children: []core.Node{camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{
					Far:    &core.AnimatedStyle{TranslateX: lerp(p, 30, -30)},
					Near:   &core.AnimatedStyle{TranslateX: lerp(p, 90, -90)},
					Figure: &core.AnimatedStyle{TranslateX: lerp(p, 10, 20)},
				})}}
		}})
}

// #174 Follow shot: behind the walker, on the same path.
func c2Follow(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "174", Title: "Follow shot", Duration: 2400 * ms, Accent: colGreen, Height: 110,
		Note: "Third person: the walker weaves and the camera weaves with him, so he stays near the centre and the ground does the moving.",
		Stage: func(p *core.AnimatedValue) core.Node {
			weave := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.25, 0.5, 0.75, 1}, Output: []float64{0, 18, 0, -18, 0}})
			return core.View{Style: stageFill, Children: []core.Node{
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					Animated: &core.AnimatedStyle{TranslateX: core.Interpolate(weave, core.Interpolation{Input: []float64{-18, 18}, Output: []float64{30, -30}})}},
					Children: []core.Node{camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{
						Near: &core.AnimatedStyle{TranslateY: lerp(p, 0, 40), Scale: lerp(p, 1, 1.4)},
					})}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(30), Left: core.Pt(120),
					Animated: &core.AnimatedStyle{TranslateX: weave}}, Children: []core.Node{
					Disc(22, core.RGB(0x2a3646), nil),
					core.View{Style: core.Style{Width: core.Pt(30), Height: core.Pt(56), CornerRadius: 14, Background: core.RGB(0x2a3646)}},
				}},
			}}
		}})
}

// #175 Locked-off shot: the camera is on sticks and does not move a pixel.
func c2LockedOff(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "175", Title: "Locked-off shot", Duration: 2200 * ms, Accent: colCyan, Height: 110,
		Note: "Everything that moves, moves inside the frame: the tripod does not follow, so the walker leaves the shot.",
		Stage: func(p *core.AnimatedValue) core.Node {
			return core.View{Style: stageFill, Children: []core.Node{
				camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(34), Left: core.Pt(20),
					Animated: &core.AnimatedStyle{TranslateX: lerp(p, 0, 230)}}, Children: []core.Node{
					Disc(12, core.RGB(0xf2dfc8), nil),
					Panel(16, 30, core.RGB(0x2a3646), nil),
				}},
			}}
		}})
}

// #176 Fly-through: through the window and out the other side.
//
// A gate that grows is a zoom wearing a costume: nothing passes the lens. The
// real shot travels the whole set, so geometry rushes past — but it has to end
// on something, so the path is a diagonal across the set with the look target
// pinned to its centre: the frame never goes blank, and by the end it is
// looking back at the set from the far side.
func c2FlyThrough(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "176", Title: "Fly-through", Duration: 3000 * ms, Accent: colBlue, Height: 110,
		Note: "A diagonal pass across the set with the look pinned to its centre: the cube goes by the lens at half way, and the shot ends looking back at the set from the far side.",
		Stage: func(p *core.AnimatedValue) core.Node {
			// The push accelerates: a fly-through that starts at speed reads
			// as a pan, so the whole travel is eased in from a standstill.
			surge := core.Interpolate(p, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1}, Easing: bezier(core.EaseIn)})
			lane := func(from, to float64) core.AnimatedNode {
				return core.Interpolate(surge, core.Interpolation{Input: []float64{0, 1}, Output: []float64{from, to}})
			}
			return core.Shader{
				Program: rayProgram,
				Label:   "Fly-through",
				Style:   stageFill,
				Uniforms: core.Uniforms{
					"Eye":  []any{lane(3.20, -2.40), lane(1.10, 0.95), lane(2.40, -0.55)},
					"Look": []any{constant(cx, 0.0), constant(cx, 0.45), constant(cx, 0.35)},
					"Fov":  lane(0.95, 1.00),
					"Tint": constant(cx, 0.0),
				},
			}
		}})
}

// #177 Low angle: from near the ground, the figure becomes a monument.
func c2LowAngle(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "177", Title: "Low angle", Duration: 2200 * ms, Accent: colAmber, Height: 110,
		Note: "The camera drops and tilts up: the figure grows, the horizon goes below the frame, the sky takes the rest.",
		Stage: func(p *core.AnimatedValue) core.Node {
			return camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{
				Sky:    &core.AnimatedStyle{TranslateY: lerp(p, 0, 30)},
				Sun:    &core.AnimatedStyle{TranslateY: lerp(p, 0, 40), TranslateX: lerp(p, 0, 20)},
				Far:    &core.AnimatedStyle{TranslateY: lerp(p, 0, 60)},
				Near:   &core.AnimatedStyle{TranslateY: lerp(p, 0, 76)},
				Figure: &core.AnimatedStyle{Scale: lerp(p, 1, 1.5), TranslateY: lerp(p, 0, -14)},
			})
		}})
}

// #178 High angle: from above, the same figure is small and exposed.
func c2HighAngle(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "178", Title: "High angle", Duration: 2200 * ms, Accent: colBlue, Height: 110,
		Note: "The mirror of the low angle: the camera climbs, the figure shrinks and the ground comes up to fill the frame.",
		Stage: func(p *core.AnimatedValue) core.Node {
			return camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{
				Sun:    &core.AnimatedStyle{TranslateY: lerp(p, 0, -14)},
				Far:    &core.AnimatedStyle{TranslateY: lerp(p, 0, -30)},
				Near:   &core.AnimatedStyle{TranslateY: lerp(p, 0, -14), Scale: lerp(p, 1, 1.2)},
				Figure: &core.AnimatedStyle{Scale: lerp(p, 1, 0.62), TranslateY: lerp(p, 0, 16)},
			})
		}})
}

// #179 POV: you are the one walking, and someone ahead waves at you.
//
// The sway alone is not first person — the world has to come at you. So the
// ray marcher walks the camera forward, its eye bobbing on the step, and the
// two hands ride the same bob: the props and the geometry agree, which is what
// sells it.
func c2POV(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "179", Title: "POV shot", Duration: 2600 * ms, Accent: colCyan, Height: 110,
		Note: "First person is a low wide camera walking in: the floor streams under the frame, the spheres grow, and the hands bob with the same step.",
		Stage: func(p *core.AnimatedValue) core.Node {
			// One step is one rise and fall of the whole body: the eye and the
			// hands take the same lane, so they never part company.
			bob := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.25, 0.5, 0.75, 1}, Output: []float64{0.46, 0.495, 0.46, 0.495, 0.46}})
			handBob := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.25, 0.5, 0.75, 1}, Output: []float64{0, -2.5, 0, -2.5, 0}})
			return core.View{Style: stageFill, Children: []core.Node{
				core.Shader{
					Program: rayProgram,
					Label:   "POV shot",
					Style:   stageFill,
					Uniforms: core.Uniforms{
						"Eye":  []any{constant(cx, -0.88), bob, lerp(p, 2.60, 1.85)},
						"Look": []any{constant(cx, 0.05), constant(cx, 0.54), constant(cx, -0.70)},
						"Fov":  constant(cx, 1.00),
						"Tint": constant(cx, 0.0),
					},
				},
				// The hands come in from the bottom corners, not the middle: put
				// them side by side in a row and they read as two pebbles.
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(84), Left: core.Pt(4),
					Animated: &core.AnimatedStyle{TranslateY: handBob}}, Children: []core.Node{
					Panel(48, 26, core.RGB(0x4a2c1c), nil),
				}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(84), Right: core.Pt(4),
					Animated: &core.AnimatedStyle{TranslateY: handBob}}, Children: []core.Node{
					Panel(48, 26, core.RGB(0x4a2c1c), nil),
				}},
			}}
		}})
}

// #180 Over-the-shoulder: the foreground is a head, and it is soft.
//
// The shoulder is a real object here: the camera stands a hand's width off the
// big sphere, which is cropped by the left edge of the frame as a near-field
// mass while the subject sits beyond it. As the camera drifts, the shoulder
// slides against the background — two depths, two speeds, no painting required.
func c2OverShoulder(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "180", Title: "Over-the-shoulder", Duration: 2400 * ms, Accent: colViolet, Height: 110,
		Note: "The mass cropped by the left edge is the near sphere, a hand's width from the lens: it slides against the subject as the camera drifts, and that difference is the whole shot.",
		Stage: func(p *core.AnimatedValue) core.Node {
			return core.Shader{
				Program: rayProgram,
				Label:   "Over-the-shoulder",
				Style:   stageFill,
				Uniforms: core.Uniforms{
					"Eye":  []any{lerp(p, -1.88, -1.34), lerp(p, 1.16, 1.06), lerp(p, 1.38, 1.14)},
					"Look": []any{constant(cx, 1.10), constant(cx, 0.42), constant(cx, -0.10)},
					"Fov":  constant(cx, 0.95),
					"Tint": constant(cx, 0.15),
				},
			}
		}})
}

// #181 Bird's-eye: straight down, and more comes into view as we rise.
func c2BirdEye(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "181", Title: "Bird's-eye view", Duration: 2400 * ms, Accent: colGreen, Height: 110,
		Note: "Straight down there is no horizon left: the road is a band, the cars are rectangles and the trees are dots.",
		Stage: func(p *core.AnimatedValue) core.Node {
			rise := core.Interpolate(p, core.Interpolation{Input: []float64{0, 1}, Output: []float64{1.5, 1}, Easing: bezier(core.EaseInOut)})
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
				Background: core.RGB(0x1d3a2a), Overflow: core.Hidden,
				Animated: &core.AnimatedStyle{Scale: rise}}, Children: []core.Node{
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(52), Left: core.Pt(0), Right: core.Pt(0), Height: core.Pt(14),
					Background: core.RGB(0x3a3a42)}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(52), Left: core.Pt(0), Right: core.Pt(0), Height: core.Pt(1),
					Background: core.RGB(0xd8d8c8)}},
				Panel(14, 6, core.RGB(0xd94b3f), &core.AnimatedStyle{TranslateX: lerp(p, 0, 90), TranslateY: constant(cx, 56)}),
				Panel(14, 6, core.RGB(0x9fd8ff), &core.AnimatedStyle{TranslateX: lerp(p, 220, 120), TranslateY: constant(cx, 58)}),
			}}
		}})
}

// #182 Deep focus: the flower, the person and the hills are all sharp.
func c2DeepFocus(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "182", Title: "Deep focus", Duration: 2200 * ms, Accent: colGreen, Height: 110,
		Note: "Here the effect is the absence of one: the flower in front sharpens instead of going soft, and nothing else changes.",
		Stage: func(p *core.AnimatedValue) core.Node {
			sharp := core.Interpolate(p, core.Interpolation{Input: []float64{0.1, 0.8}, Output: []float64{0.45, 1}})
			return core.View{Style: stageFill, Children: []core.Node{
				camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(62), Left: core.Pt(18),
					Animated: &core.AnimatedStyle{Opacity: sharp, Scale: sharp}}, Children: []core.Node{Disc(34, core.RGB(0xff8fb1), nil)}},
			}}
		}})
}

// #183 Establishing shot: wide, high, and slow, to say where we are.
func c2Establishing(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "183", Title: "Establishing shot", Duration: 2600 * ms, Accent: colBlue, Height: 110,
		Note: "Wide and slow, with the town in the middle: the push-in is small on purpose and the caption does the naming.",
		Stage: func(p *core.AnimatedValue) core.Node {
			return core.View{Style: stageFill, Children: []core.Node{
				camScene(core.RGB(0x24405e), core.RGB(0x1a2a3a), camLayers{
					Sun:    &core.AnimatedStyle{Scale: lerp(p, 1, 1.1)},
					Far:    &core.AnimatedStyle{Scale: lerp(p, 1, 1.16), TranslateY: lerp(p, 0, -6)},
					Near:   &core.AnimatedStyle{Scale: lerp(p, 1, 1.26), TranslateY: lerp(p, 0, -10)},
					Figure: &core.AnimatedStyle{Scale: lerp(p, 1, 1.1)},
				}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(86), Left: core.Pt(8), Padding: core.Edges{Top: 3, Right: 6, Bottom: 3, Left: 6},
					Background: core.RGB(0x0a1018), CornerRadius: 3,
					Animated: &core.AnimatedStyle{Opacity: core.Interpolate(p, core.Interpolation{Input: []float64{0.35, 0.75}, Output: []float64{0, 1}})}},
					Children: []core.Node{uipLabel("NAGOYA  ·  06:12", 7, core.RGB(0xd8e4f0))}},
			}}
		}})
}

// #184 Extreme close-up: nothing but the eyes.
func c2ExtremeClose(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "184", Title: "Extreme close-up", Duration: 2400 * ms, Accent: colAmber, Height: 110,
		Note: "From a face to two eyes: the framing is the effect, and the push never stops being about scale.",
		Stage: func(p *core.AnimatedValue) core.Node {
			zoom := core.Interpolate(p, core.Interpolation{Input: []float64{0, 1}, Output: []float64{1, 6.4}, Easing: bezier(core.EaseInOut)})
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0), Overflow: core.Hidden},
				Children: []core.Node{
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
						Background: core.RGB(0x241a14)}},
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
						AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter,
						Animated: &core.AnimatedStyle{Scale: zoom}}, Children: []core.Node{
						core.View{Style: core.Style{Width: core.Pt(64), Height: core.Pt(64), CornerRadius: 32, Background: core.RGB(0xd9b79a),
							AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter}, Children: []core.Node{
							core.View{Style: core.Style{Direction: core.Row, Gap: 12}, Children: []core.Node{
								Disc(9, core.RGB(0x2a2018), nil),
								Disc(9, core.RGB(0x2a2018), nil),
							}},
						}},
					}},
				}}
		}})
}

// #185 Insert shot: a detail cut in and out of the scene.
func c2Insert(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "185", Title: "Insert shot", Duration: 2400 * ms, Accent: colCyan, Height: 110,
		Note: "The wide is untouched; a close-up of the watch is laid over it for one and a third seconds and then taken away.",
		Stage: func(p *core.AnimatedValue) core.Node {
			cut := core.Interpolate(p, core.Interpolation{Input: []float64{0.25, 0.27, 0.7, 0.72}, Output: []float64{0, 1, 1, 0}})
			return core.View{Style: stageFill, Children: []core.Node{
				camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{
					Figure: &core.AnimatedStyle{TranslateX: lerp(p, 20, 30)},
				}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					Background: core.RGB(0x1a1410),
					Animated:   &core.AnimatedStyle{Opacity: cut}}, Children: []core.Node{
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(30), Left: core.Pt(84),
						Animated: &core.AnimatedStyle{Scale: lerp(p, 1, 1.1)}}, Children: []core.Node{
						core.View{Style: core.Style{Width: core.Pt(40), Height: core.Pt(40), CornerRadius: 20, Background: core.RGB(0x8a8f96),
							BorderWidth: 3, BorderColor: core.RGB(0xd8d8cf)}},
						core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(16), Left: core.Pt(16),
							Width: core.Pt(2), Height: core.Pt(12), Background: core.RGB(0x22262c)}},
					}},
				}},
			}}
		}})
}

// #186 Morph cut: two takes, joined by moving rather than by cutting.
func c2MorphCut(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "186", Title: "Morph cut", Duration: 2400 * ms, Accent: colViolet, Height: 110,
		Note: "The head slides and settles instead of jumping: the same subject, two takes, no visible seam.",
		Stage: func(p *core.AnimatedValue) core.Node {
			drift := core.Interpolate(p, core.Interpolation{Input: []float64{0.2, 0.8}, Output: []float64{0, 14}, Easing: bezier(core.EaseInOut)})
			return core.View{Style: stageFill, Children: []core.Node{
				camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{
					Figure: &core.AnimatedStyle{TranslateX: drift, Rotate: core.Interpolate(drift, core.Interpolation{Input: []float64{0, 14}, Output: []float64{-2, 2}})},
					Near:   &core.AnimatedStyle{TranslateX: core.Interpolate(drift, core.Interpolation{Input: []float64{0, 14}, Output: []float64{0, -6}})},
				}),
			}}
		}})
}

// #187 Fade to black: out, then the next scene fades up out of the same black.
func c2FadeBlack(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "187", Title: "Fade to black", Duration: 2600 * ms, Accent: colFaint, Height: 110,
		Note: "Two pictures and one black between them: the black is what does the cutting, so neither shot has to move.",
		Stage: func(p *core.AnimatedValue) core.Node {
			black := core.Interpolate(p, core.Interpolation{Input: []float64{0.2, 0.48, 0.52, 0.8}, Output: []float64{0, 1, 1, 0}})
			next := core.Interpolate(p, core.Interpolation{Input: []float64{0.5, 0.501}, Output: []float64{0, 1}})
			return core.View{Style: stageFill, Children: []core.Node{
				camShots(&core.AnimatedStyle{}, &core.AnimatedStyle{Opacity: next}),
				phScrim(core.RGB(0x000000), black),
			}}
		}})
}

// #188 Iris out: a circle closes on the person, then opens on the next scene.
func c2IrisOut(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "188", Title: "Iris out", Duration: 2400 * ms, Accent: colAmber, Height: 110,
		Note: "The old silent-film close: a round hole shrinks until the frame is black, and the next shot opens out of it.",
		Stage: func(p *core.AnimatedValue) core.Node {
			close := core.Interpolate(p, core.Interpolation{Input: []float64{0.1, 0.45, 0.55, 0.9}, Output: []float64{1.6, 0.02, 0.02, 1.6}})
			next := core.Interpolate(p, core.Interpolation{Input: []float64{0.5, 0.501}, Output: []float64{0, 1}})
			return core.View{Style: stageFill, Children: []core.Node{
				camShots(&core.AnimatedStyle{}, &core.AnimatedStyle{Opacity: next}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter}, Children: []core.Node{
					// Four black slabs leaving a hole the ring closes over.
					core.View{Style: core.Style{Width: core.Pt(150), Height: core.Pt(150), CornerRadius: 75, BorderWidth: 76, BorderColor: core.RGB(0x000000),
						Animated: &core.AnimatedStyle{Scale: close}}},
				}},
			}}
		}})
}

// #189 Cutaway: here is what the listener is doing.
func c2Cutaway(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "189", Title: "Cutaway", Duration: 2400 * ms, Accent: colGreen, Height: 110,
		Note: "A different subject, cut in for a moment: the speaker, then the room, then the speaker — nothing in the scene changes.",
		Stage: func(p *core.AnimatedValue) core.Node {
			room := core.Interpolate(p, core.Interpolation{Input: []float64{0.3, 0.32, 0.66, 0.68}, Output: []float64{0, 1, 1, 0}})
			var audience []core.Node
			for i := 0; i < 4; i++ {
				audience = append(audience, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(core.Points(float32(58 + (i%2)*18))),
					Left:     core.Pt(core.Points(float32(24 + i*52))),
					Animated: &core.AnimatedStyle{Scale: lerp(p, 1, 1.04)}}, Children: []core.Node{
					Disc(16, core.RGB(0x2a3646), nil),
					Panel(22, 18, core.RGB(0x1f2a38), nil),
				}})
			}
			return core.View{Style: stageFill, Children: []core.Node{
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0)},
					Children: []core.Node{camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{})}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					Background: core.RGB(0x141d28), Animated: &core.AnimatedStyle{Opacity: room}}, Children: audience},
			}}
		}})
}

// #190 Cross-cutting: two places at once, the cuts getting shorter.
func c2CrossCutting(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "190", Title: "Cross-cutting", Duration: 2600 * ms, Accent: colPink, Height: 110,
		Note: "Nobody moves, but the picture will not sit still: the cuts come faster and faster until the two shots are a flicker.",
		Stage: func(p *core.AnimatedValue) core.Node {
			// A run of alternations with a shortening period.
			var cuts []core.Node
			at := 0.1
			span := 0.22
			side := 0
			for i := 0; i < 5; i++ {
				side = i % 2
				_ = side
				cuts = append(cuts, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					Background: core.RGB(0x05080d),
					Animated:   &core.AnimatedStyle{Opacity: core.Interpolate(p, core.Interpolation{Input: []float64{at, at + 0.005, at + span - 0.005, at + span}, Output: []float64{0, 1, 1, 0}})}}})
				at += span
				span *= 0.62
			}
			return core.View{Style: stageFill, Children: []core.Node{
				camShots(&core.AnimatedStyle{}, &core.AnimatedStyle{}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(88), Left: core.Pt(10), Right: core.Pt(10), Height: core.Pt(14),
					Direction: core.Row, Gap: 3}, Children: []core.Node{
					Panel(52, 12, core.RGB(0x39d98a), &core.AnimatedStyle{TranslateX: lerp(p, 0, 4)}),
					Panel(30, 12, core.RGB(0x35d6de), &core.AnimatedStyle{TranslateX: lerp(p, 0, 140)}),
				}},
			}}
		}})
}

// #191 Bullet time: time nearly stops and the camera goes round.
func c2BulletTime(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "191", Title: "Bullet time", Duration: 2600 * ms, Accent: colViolet, Height: 110,
		Note: "The jump stops at the top and the camera keeps going: the figure barely moves while the world slides behind him.",
		Stage: func(p *core.AnimatedValue) core.Node {
			frozen := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.3, 0.4, 0.85, 1}, Output: []float64{0, -30, -32, -32, -30}})
			orbit := core.Interpolate(p, core.Interpolation{Input: []float64{0.3, 0.9}, Output: []float64{0, 120}, Easing: bezier(core.EaseInOut)})
			return camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{
				Far:    &core.AnimatedStyle{TranslateX: orbit},
				Near:   &core.AnimatedStyle{TranslateX: core.Interpolate(orbit, core.Interpolation{Input: []float64{0, 120}, Output: []float64{0, 190}})},
				Sun:    &core.AnimatedStyle{TranslateX: core.Interpolate(orbit, core.Interpolation{Input: []float64{0, 120}, Output: []float64{0, 80}})},
				Figure: &core.AnimatedStyle{TranslateY: frozen, Rotate: lerp(p, -6, 6)},
			})
		}})
}

// #192 Montage: seven shots, thirty days, one straight cut after another.
func c2Montage(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "192", Title: "Montage", Duration: 2600 * ms, Accent: colAmber, Height: 110,
		Note: "Seven shots at seven steps each: what makes a montage is that the cuts are straight and the date never stops moving.",
		Stage: func(p *core.AnimatedValue) core.Node {
			tints := []core.Color{core.RGB(0x1b3a5c), core.RGB(0x24405e), core.RGB(0x3a2b3c), core.RGB(0x1d4a3f),
				core.RGB(0x5a2b3c), core.RGB(0x2a2350), core.RGB(0x1b3a5c)}
			var shots []core.Node
			for i, col := range tints {
				from := float64(i) / 7
				to := float64(i+1) / 7
				// The last shot stays: a tile that ends on an empty frame
				// looks broken the moment the pass is over.
				fade := []float64{0, 1, 1, 0}
				if i == len(tints)-1 {
					fade = []float64{0, 1, 1, 1}
				}
				shots = append(shots, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					Animated: &core.AnimatedStyle{Opacity: core.Interpolate(p, core.Interpolation{Input: []float64{from, from + 0.004, to - 0.004, to}, Output: fade})}},
					Children: []core.Node{camScene(col, core.RGB(0x0e1420), camLayers{
						Figure: &core.AnimatedStyle{TranslateX: lerp(p, float64(i*6), float64(i*6+14))},
					})}})
			}
			const days = 30
			var dayLayers []core.Node
			for d := 1; d <= days; d++ {
				label := "DAY "
				if d < 10 {
					label += "0"
				}
				label += itoa(d)
				dayLayers = append(dayLayers, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter,
					Animated: &core.AnimatedStyle{Opacity: seqFade(p, d-1, days-1, 1.0, 0.006)}},
					Children: []core.Node{core.Text{Content: label, Style: core.Style{FontSize: 8, Color: core.RGB(0xffd27a)}}}})
			}
			stamp := core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(6), Right: core.Pt(8),
				Padding: core.Edges{Top: 2, Right: 5, Bottom: 2, Left: 5}, Background: core.RGB(0x0a1018), CornerRadius: 3},
				Children: []core.Node{
					core.View{Style: core.Style{Position: core.Relative, Width: core.Pt(40), Height: core.Pt(10)}, Children: dayLayers},
				}}
			return core.View{Style: stageFill, Children: append(append([]core.Node{}, shots...), stamp)}
		}})
}
