package main

// page_motion2.go is the motion2 pack (motion-prompt #121–#144): the classic
// techniques, then text and shapes, then four harder transitions.
//
// Catalog names and timings: github.com/showheyohtaki/motion-prompt (MIT).

import (
	"math"
	"time"

	"github.com/oneo-me/nsgo/core"
)

// pageMotion2 lays the pack out.
func pageMotion2(cx *core.Scope) core.Node {
	return core.View{Children: []core.Node{
		core.View{Style: core.Style{Padding: core.Edges{Top: 14, Right: 18, Left: 18}}, Children: []core.Node{
			note("Motion2 pack 121–144 of motion-prompt (github.com/showheyohtaki/motion-prompt): anticipation, follow-through, loops, text and four transitions."),
		}},
		grid(
			m2Anticipation(cx),
			m2FollowThrough(cx),
			m2Bounce(cx),
			m2Wiggle(cx),
			m2Pulse(cx),
			m2Path(cx),
			m2PingPong(cx),
			m2StopMotion(cx),
			m2Scramble(cx),
			m2WaveText(cx),
			m2TextOnPath(cx),
			m2VariableFont(cx),
			m2TextMask(cx),
			m2Repeater(cx),
			m2Kaleidoscope(cx),
			m2LineBoil(cx),
			m2Push(cx),
			m2Blinds(cx),
			m2Cube(cx),
			m2Curl(cx),
			m2Pixelate(cx),
			m2Shatter(cx),
			m2Exploded(cx),
			m2Split(cx),
		),
	}}
}

// constant is a fixed value where an animated node is expected: the packs
// build a rotation or an offset from the frame counter only sometimes.
func constant(cx *core.Scope, v float64) core.AnimatedNode {
	return core.NewAnimatedValue(cx, v)
}

// m2Char is the little square character most of the technique tiles use.
func m2Char(size float32, col core.Color, anim *core.AnimatedStyle) core.Node {
	return core.View{Style: core.Style{Width: core.Pt(core.Points(size)), Height: core.Pt(core.Points(size)), CornerRadius: core.Points(size) / 5,
		Background: col, Animated: anim}}
}

// #121 Anticipation: back, coil, then forward — the wind-up is the move.
func m2Anticipation(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "121", Title: "Anticipation", Duration: 1800 * ms, Accent: colCyan,
		Note: "The direction is set by going the wrong way first: back, squash, then across the whole stage.",
		Stage: func(p *core.AnimatedValue) core.Node {
			x := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.3, 0.45, 0.8, 1}, Output: []float64{0, -22, -26, 110, 110}})
			squash := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.3, 0.45, 0.55, 1}, Output: []float64{1, 1.15, 0.8, 1, 1}})
			return uipScreen(170, 96,
				At(86, 10, uipBar(150, 2, core.RGB(0x1b2534), nil)),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(52), Left: core.Pt(18),
					Animated: &core.AnimatedStyle{TranslateX: x}},
					Children: []core.Node{m2Char(26, colCyan, &core.AnimatedStyle{ScaleX: squash})}},
			)
		}})
}

// #122 Follow-through: the body stops dead, the flag keeps going.
func m2FollowThrough(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "122", Title: "Follow-through", Duration: 2000 * ms, Accent: colAmber,
		Note: "The body arrives and stops; the tail and the flag on its head arrive late and settle after it.",
		Stage: func(p *core.AnimatedValue) core.Node {
			body := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.55, 1}, Output: []float64{-40, 60, 60}})
			swing := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.55, 0.68, 0.8, 0.92, 1}, Output: []float64{0, 34, -22, 14, -6, 0}})
			return uipScreen(150, 96,
				At(84, 6, uipBar(140, 2, core.RGB(0x1b2534), nil)),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(44), Left: core.Pt(26),
					Animated: &core.AnimatedStyle{TranslateX: body}}, Children: []core.Node{
					core.View{Style: core.Style{Position: core.Relative}, Children: []core.Node{
						m2Char(26, colAmber, nil),
						core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(-14), Left: core.Pt(22),
							Animated: &core.AnimatedStyle{Rotate: swing}}, Children: []core.Node{
							core.View{Style: core.Style{Width: core.Pt(16), Height: core.Pt(9), CornerRadius: 2, Background: colPink}},
						}},
					}}},
				},
			)
		}})
}

// #123 Bounce: each landing gives back less than the last.
func m2Bounce(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "123", Title: "Bounce", Duration: 2400 * ms, Accent: colGreen,
		Note: "Four arcs, each lower than the one before; the ground is where the energy goes.",
		Stage: func(p *core.AnimatedValue) core.Node {
			y := core.Interpolate(p, core.Interpolation{
				Input:  []float64{0, 0.22, 0.4, 0.55, 0.68, 0.8, 0.9, 1},
				Output: []float64{-64, 0, -34, 0, -16, 0, -6, 0},
			})
			squash := core.Interpolate(p, core.Interpolation{
				Input:  []float64{0, 0.2, 0.24, 0.38, 0.42, 0.53, 0.57},
				Output: []float64{1, 1, 0.72, 1, 0.8, 1, 0.88},
			})
			return uipScreen(150, 96,
				At(80, 8, uipBar(130, 2, core.RGB(0x1b2534), nil)),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(54), Left: core.Pt(56),
					Animated: &core.AnimatedStyle{TranslateY: y, ScaleY: squash, ScaleX: core.Interpolate(squash, core.Interpolation{Input: []float64{0.7, 1}, Output: []float64{1.3, 1}})}},
					Children: []core.Node{m2Char(26, colGreen, nil)}},
			)
		}})
}

// #124 Wiggle: never a single hit, always a small irregular shake.
func m2Wiggle(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "124", Title: "Wiggle", Duration: 1200 * ms, Accent: colAmber, Loop: true,
		Note: "A bell that rings itself: the shake is irregular on purpose, so it reads as ringing, not as a loop.",
		Stage: func(p *core.AnimatedValue) core.Node {
			jit := core.Interpolate(p, core.Interpolation{
				Input:  []float64{0, 0.12, 0.24, 0.33, 0.41, 0.55, 0.64, 0.73, 0.86, 1},
				Output: []float64{0, -4, 3.5, -3, 4, -2.5, 3, -1.5, 1, 0},
			})
			tilt := core.Interpolate(p, core.Interpolation{
				Input:  []float64{0, 0.15, 0.3, 0.45, 0.6, 0.8, 1},
				Output: []float64{0, -7, 6, -5, 4, -2, 0},
			})
			return uipScreen(140, 96,
				Centered(core.View{Style: core.Style{Position: core.Relative, Width: core.Pt(90), Height: core.Pt(90)}, Children: []core.Node{
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(20), Left: core.Pt(22),
						Animated: &core.AnimatedStyle{TranslateX: jit, Rotate: tilt}}, Children: []core.Node{
						core.View{Style: core.Style{Width: core.Pt(46), Height: core.Pt(34), CornerRadius: 20, Background: colAmber}},
						core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(34), Left: core.Pt(19)},
							Children: []core.Node{Disc(9, core.RGB(0x2a2013), nil)}},
					}},
				}}))
		}})
}

// #125 Pulse: a double beat, with rings that leave on each one.
func m2Pulse(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "125", Title: "Pulse", Duration: 1400 * ms, Accent: colPink, Loop: true,
		Note: "Two beats and a pause: the second is shorter, which is what makes it a heartbeat instead of a throb.",
		Stage: func(p *core.AnimatedValue) core.Node {
			beat := core.Interpolate(p, core.Interpolation{
				Input:  []float64{0, 0.08, 0.16, 0.26, 0.34, 0.6, 1},
				Output: []float64{1, 1.14, 1, 1.09, 1, 1, 1},
			})
			ring := func(delay float64) core.Node {
				q := core.Interpolate(p, core.Interpolation{Input: []float64{delay, 0.7 + delay}, Output: []float64{0, 1}})
				return core.View{Style: core.Style{Animated: &core.AnimatedStyle{
					Scale:   core.Interpolate(q, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0.9, 2.4}}),
					Opacity: core.Interpolate(q, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0.8, 0}}),
				}}, Children: []core.Node{Ring(34, 2, colPink, nil)}}
			}
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
				AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter}, Children: []core.Node{
				core.View{Style: core.Style{Position: core.Relative, Width: core.Pt(40), Height: core.Pt(40),
					AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter}, Children: []core.Node{ring(0), ring(0.18)}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter, Animated: &core.AnimatedStyle{Scale: beat}},
					Children: []core.Node{Disc(30, colPink, nil)}},
			}}
		}})
}

// #126 Motion path: the plane follows the dotted line and turns with it.
func m2Path(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "126", Title: "Motion path", Duration: 2600 * ms, Accent: colCyan,
		Note: "Position and heading come from the same curve, which is why the plane never slides sideways.",
		Stage: func(p *core.AnimatedValue) core.Node {
			x := core.Interpolate(p, core.Interpolation{
				Input:  []float64{0, 0.25, 0.5, 0.75, 1},
				Output: []float64{6, 60, 116, 60, 6},
			})
			y := core.Interpolate(p, core.Interpolation{
				Input:  []float64{0, 0.25, 0.5, 0.75, 1},
				Output: []float64{52, 8, 46, 80, 52},
			})
			head := core.Interpolate(p, core.Interpolation{
				Input:  []float64{0, 0.25, 0.5, 0.75, 1},
				Output: []float64{-30, -90, 60, 140, 330},
			})
			var dots []core.Node
			for i := 0; i < 7; i++ {
				t := float64(i) / 6
				dots = append(dots, core.View{Style: core.Style{Position: core.Absolute,
					Top:  core.Pt(core.Points(float32(80 - 34*math.Sin(t*math.Pi*2)))),
					Left: core.Pt(core.Points(float32(8 + 110*t)))}, Children: []core.Node{Disc(2, core.RGB(0x2f4256), nil)}})
			}
			return uipScreen(150, 96, append([]core.Node{
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0),
					Animated: &core.AnimatedStyle{TranslateX: x, TranslateY: y, Rotate: head}}, Children: []core.Node{
					core.View{Style: core.Style{Width: core.Pt(0), Height: core.Pt(0), BorderWidth: 0}},
					Panel(14, 8, colCyan, nil),
				}},
			}, dots...)...)
		}})
}

// #127 Ping-pong loop: out and back along exactly the same path.
func m2PingPong(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "127", Title: "Ping-pong loop", Duration: 1600 * ms, Accent: colGreen, Loop: true,
		Note: "One keyframe pair, played forwards and then backwards: at the ends the motion has to stop dead.",
		Stage: func(p *core.AnimatedValue) core.Node {
			x := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.5, 1}, Output: []float64{0, 116, 0}, Easing: bezier(core.EaseInOut)})
			return uipScreen(150, 96,
				At(70, 8, uipBar(130, 2, core.RGB(0x1b2534), nil)),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(42), Left: core.Pt(10),
					Animated: &core.AnimatedStyle{TranslateX: x}}, Children: []core.Node{m2Char(24, colGreen, nil)}},
			)
		}})
}

// #128 Stop-motion: eight frames a second, every one of them a little off.
func m2StopMotion(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "128", Title: "Stop-motion style", Duration: 2000 * ms, Accent: colAmber,
		Note: "Sampled, not smoothed: seven hops of a twelfth each, and each frame is a pixel or two from where it belongs.",
		Stage: func(p *core.AnimatedValue) core.Node {
			var star []core.Node
			for i := 0; i < 7; i++ {
				t0 := float64(i) / 7
				jitter := []float64{0.5, -0.7, 0.9, -0.4, 0.6, -0.8, 0.3}[i]
				step := core.Interpolate(p, core.Interpolation{
					Input:  []float64{t0, t0 + 0.06},
					Output: []float64{0, 1},
				})
				star = append(star, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(core.Points(float32(30 - 14*jitter))),
					Left:     core.Pt(core.Points(float32(6 + 19*float64(i)))),
					Animated: &core.AnimatedStyle{Scale: step, Rotate: constant(cx, 12*jitter)}}, Children: []core.Node{Disc(11, colAmber, nil)}})
			}
			return uipScreen(150, 96, append([]core.Node{
				At(70, 8, uipBar(130, 2, core.RGB(0x1b2534), nil)),
			}, star...)...)
		}})
}

// #129 Text scramble: noise, then the right letters lock in from the left.
func m2Scramble(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "129", Title: "Text scramble", Duration: 2000 * ms, Accent: colGreen,
		Note: "Each column settles on its own schedule; the left one locks first, so the word resolves like a decode.",
		Stage: func(p *core.AnimatedValue) core.Node {
			word := "DECODED"
			n := len(word)
			shown := core.NewSignal(cx, "")
			cx.AfterFunc(30*ms, func() {
				for i := 1; i <= n; i++ {
					at := i
					cx.AfterFunc(time.Duration(at)*180*ms, func() {
						s := word[:at]
						for k := 0; k < 3 && at+k < n; k++ {
							s += string(byte("XKQZ#%&"[int(math.Mod(float64(at+k)*7, 7))]))
						}
						shown.Set(s)
					})
				}
				cx.AfterFunc(time.Duration(n+1)*180*ms, func() { shown.Set(word) })
			})
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
				AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter}, Children: []core.Node{
				core.Text{ContentFunc: func() string { return shown.Get() },
					Style: core.Style{FontSize: 22, FontWeight: core.FontWeightBold, Color: colGreen, LetterSpacing: 2}},
			}}
		}})
}

// #130 Wave text: one letter after another, the whole word surfing.
func m2WaveText(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "130", Title: "Wave text", Duration: 1800 * ms, Accent: colCyan, Loop: true,
		Note: "The letters do not move together: each one is a tenth of a cycle behind its neighbour.",
		Stage: func(p *core.AnimatedValue) core.Node {
			word := "WAVE TEXT"
			var letters []core.Node
			for i := 0; i < len(word); i++ {
				if word[i] == ' ' {
					letters = append(letters, core.View{Style: core.Style{Width: core.Pt(6)}})
					continue
				}
				off := core.Interpolate(p, core.Interpolation{
					Input:  []float64{0, 0.2, 0.4, 0.6, 0.8, 1},
					Output: waveOffset(i),
				})
				letters = append(letters, core.View{Style: core.Style{Animated: &core.AnimatedStyle{TranslateY: off, Rotate: core.Interpolate(off, core.Interpolation{Input: []float64{-12, 12}, Output: []float64{-6, 6}})}},
					Children: []core.Node{Word(string(word[i]), 20, colCyan, nil)}})
			}
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
				AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter}, Children: []core.Node{row(1, letters...)}}
		}})
}

// waveOffset gives the six knots of a sine for letter i, so neighbours differ.
func waveOffset(i int) []float64 {
	phase := 0.7 * float64(i)
	out := make([]float64, 6)
	for k := 0; k < 6; k++ {
		out[k] = -10 * math.Sin(phase+float64(k)*1.05)
	}
	return out
}

// #131 Text on a path: every letter sits on the curve and turns with it.
func m2TextOnPath(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "131", Title: "Text on a path", Duration: 2400 * ms, Accent: colViolet,
		Note: "Each letter is placed and turned from the same curve, so the line of text stays welded to the wave as it travels.",
		Stage: func(p *core.AnimatedValue) core.Node {
			word := "follow the line"
			travel := core.Interpolate(p, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 40}})
			var letters []core.Node
			for i := 0; i < len(word); i++ {
				t := float64(i)/float64(len(word)) + 0.02
				if word[i] == ' ' {
					continue
				}
				letters = append(letters, core.View{Style: core.Style{Position: core.Absolute,
					Top:  core.Pt(core.Points(float32(46 - 26*math.Sin(t*math.Pi*2)))),
					Left: core.Pt(core.Points(float32(4 + 160*t))),
					Animated: &core.AnimatedStyle{TranslateX: travel,
						Rotate: core.Interpolate(travel, core.Interpolation{Input: []float64{0, 40}, Output: []float64{-8 * math.Cos(t*math.Pi*2), -8 * math.Cos(t*math.Pi*2)}})}},
					Children: []core.Node{Word(string(word[i]), 11, colViolet, nil)}})
			}
			return uipScreen(180, 96, letters...)
		}})
}

// #132 Variable font: thin to black and back, with the axis shown.
func m2VariableFont(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "132", Title: "Variable font", Duration: 2000 * ms, Accent: colText,
		Note: "A real variable font would slide one axis; here two weights cross-fade, which is the same thing at this size.",
		Stage: func(p *core.AnimatedValue) core.Node {
			bold := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.45, 0.55, 1}, Output: []float64{0, 1, 1, 0}})
			knob := core.Interpolate(p, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 120}, Easing: bezier(core.EaseInOut)})
			return uipScreen(160, 96,
				At(26, 10, core.View{Style: core.Style{Position: core.Relative}, Children: []core.Node{
					core.Text{Content: "Weight", Style: core.Style{FontSize: 24, FontWeight: core.FontWeightLight, Color: colText}},
					core.Text{Content: "Weight", Style: core.Style{Position: core.Absolute, FontSize: 24, FontWeight: core.FontWeightBold, Color: colText,
						Animated: &core.AnimatedStyle{Opacity: bold}}},
				}}),
				At(72, 10, uipBar(120, 2, core.RGB(0x1b2534), nil)),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(66), Left: core.Pt(10),
					Animated: &core.AnimatedStyle{TranslateX: knob}}, Children: []core.Node{Disc(8, colText, nil)}},
			)
		}})
}

// #133 Text mask: the picture only shows inside the letters.
func m2TextMask(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "133", Title: "Text mask", Duration: 2400 * ms, Accent: colBlue,
		Note: "There is no clipping to a glyph here: the scene is a window, and the word is a stencil drawn over it.",
		Stage: func(p *core.AnimatedValue) core.Node {
			reveal := core.Interpolate(p, core.Interpolation{Input: []float64{0.75, 1}, Output: []float64{0, 1}})
			window := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.7}, Output: []float64{0.85, 1}})
			return core.View{Style: stageFill, Children: []core.Node{
				camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{
					Far:  &core.AnimatedStyle{TranslateX: lerp(p, 40, -40)},
					Near: &core.AnimatedStyle{TranslateX: lerp(p, 80, -80)},
				}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter,
					Animated: &core.AnimatedStyle{Opacity: core.Interpolate(reveal, core.Interpolation{Input: []float64{0, 1}, Output: []float64{1, 0}}),
						Scale: window}},
					Children: []core.Node{
						core.View{Style: core.Style{Background: core.RGB(0x080b11), Padding: core.Edges{Top: 6, Right: 8, Bottom: 6, Left: 8}},
							Children: []core.Node{Word("SEA", 30, core.RGB(0x080b11), nil)}},
					}},
			}}
		}})
}

// #134 Repeater: one shape, eight copies, each a beat behind.
func m2Repeater(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "134", Title: "Repeater", Duration: 1800 * ms, Accent: colCyan, Loop: true,
		Note: "The copies are identical; the stagger is the only thing that turns a ring into something moving.",
		Stage: func(p *core.AnimatedValue) core.Node {
			var ring []core.Node
			for i := 0; i < 8; i++ {
				delay := 0.06 * float64(i)
				q := core.Interpolate(p, core.Interpolation{Input: []float64{delay, 0.5 + delay, 1}, Output: []float64{0.7, 1.15, 0.7}})
				ang := 45.0 * float64(i)
				ring = append(ring, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(core.Points(float32(40 - 30*math.Sin(ang*math.Pi/180)))),
					Left:     core.Pt(core.Points(float32(40 + 30*math.Cos(ang*math.Pi/180)))),
					Animated: &core.AnimatedStyle{Scale: q, Rotate: constant(cx, ang)}},
					Children: []core.Node{Panel(9, 9, colCyan, nil)}})
			}
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
				AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter}, Children: []core.Node{
				core.View{Style: core.Style{Position: core.Relative, Width: core.Pt(80), Height: core.Pt(80)}, Children: append(ring,
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(34), Left: core.Pt(34)},
						Children: []core.Node{Disc(12, core.RGB(0xeaf1fa), nil)}})},
			}}
		}})
}

// #135 Kaleidoscope: eight wedges, one turn, symmetry doing the work.
func m2Kaleidoscope(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "135", Title: "Kaleidoscope", Duration: 2000 * ms, Accent: colViolet,
		Note: "Eight copies of the same wedge at eight angles: the pattern comes from the symmetry, not from drawing eight shapes.",
		Stage: func(p *core.AnimatedValue) core.Node {
			turn := core.Interpolate(p, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 45}})
			grow := core.Interpolate(p, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0.6, 1.25}, Easing: bezier(core.EaseInOut)})
			var wedges []core.Node
			for i := 0; i < 8; i++ {
				wedges = append(wedges, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter,
					Animated: &core.AnimatedStyle{Rotate: constant(cx, 45.0*float64(i))}},
					Children: []core.Node{
						core.View{Style: core.Style{Width: core.Pt(6), Height: core.Pt(34), CornerRadius: 3,
							Background: []core.Color{colViolet, colPink, colCyan, colBlue}[i%4],
							Animated:   &core.AnimatedStyle{Scale: grow, TranslateY: constant(cx, -22)}}},
					}})
			}
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
				AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter,
				Animated: &core.AnimatedStyle{Rotate: turn}}, Children: wedges}
		}})
}

// #136 Line boil: a still drawing that wobbles on twos.
func m2LineBoil(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "136", Title: "Line boil", Duration: 1400 * ms, Accent: colText, Loop: true,
		Note: "The drawing never moves; the line does, on a six-frame cycle of small offsets — the hand-drawn illusion.",
		Stage: func(p *core.AnimatedValue) core.Node {

			wob := func(seed float64) core.AnimatedNode {
				return core.Interpolate(p, core.Interpolation{
					Input:  []float64{0, 0.16, 0.33, 0.5, 0.66, 0.83, 1},
					Output: []float64{0, seed, -seed, seed * 0.6, -seed * 0.8, seed * 0.5, 0},
				})
			}
			return uipScreen(140, 96,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(24), Left: core.Pt(14),
					Width: core.Pt(112), Height: core.Pt(50), BorderWidth: 2, BorderColor: colText, CornerRadius: 6,
					Animated: &core.AnimatedStyle{TranslateX: wob(1.2), TranslateY: wob(0.8)}}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(44), Left: core.Pt(34),
					Width: core.Pt(40), Height: core.Pt(30), BorderWidth: 2, BorderColor: colText, CornerRadius: 15,
					Animated: &core.AnimatedStyle{TranslateX: wob(0.9)}}},
			)
		}})
}

// #137 Push: the next shot arrives and takes the frame with it.
func m2Push(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "137", Title: "Push transition", Duration: 2000 * ms, Accent: colGreen, Height: 110,
		Note: "Both scenes travel the same distance at the same speed; the seam between them is the whole transition.",
		Stage: func(p *core.AnimatedValue) core.Node {
			slide := core.Interpolate(p, core.Interpolation{Input: []float64{0.15, 0.85}, Output: []float64{0, 260}, Easing: bezier(core.EaseInOut)})
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0), Overflow: core.Hidden},
				Children: []core.Node{
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
						Animated: &core.AnimatedStyle{TranslateX: slide}}, Children: []core.Node{camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{})}},
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(-260), Bottom: core.Pt(0), Width: core.Pt(260),
						Animated: &core.AnimatedStyle{TranslateX: slide}}, Children: []core.Node{camScene(core.RGB(0x5a2b3c), core.RGB(0x2a1420), camLayers{})}},
				}}
		}})
}

// #138 Blinds: strips that narrow until the next scene is fully open.
func m2Blinds(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "138", Title: "Blinds transition", Duration: 2000 * ms, Accent: colAmber, Height: 110,
		Note: "Six strips, one scale each: they shrink on their own axes and the gaps between them are the next picture.",
		Stage: func(p *core.AnimatedValue) core.Node {
			open := core.Interpolate(p, core.Interpolation{Input: []float64{0.2, 0.8, 1}, Output: []float64{1, 0.04, 0.04}, Easing: bezier(core.EaseInOut)})
			var strips []core.Node
			for i := 0; i < 6; i++ {
				strips = append(strips, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(core.Points(float32(i) * 19)),
					Left: core.Pt(0), Right: core.Pt(0), Height: core.Pt(17), Background: core.RGB(0x1b3a5c),
					Animated: &core.AnimatedStyle{ScaleY: open}}})
			}
			return core.View{Style: stageFill, Children: []core.Node{
				camScene(core.RGB(0x5a2b3c), core.RGB(0x2a1420), camLayers{}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0)}, Children: strips},
			}}
		}})
}

// #139 Cube: the front face swings away and the next face comes round.
func m2Cube(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "139", Title: "Cube transition", Duration: 2200 * ms, Accent: colBlue, Height: 110,
		Note: "A cube turn is two faces and one rotation: the old face going away is already ninety degrees round.",
		Stage: func(p *core.AnimatedValue) core.Node {
			turn := core.Interpolate(p, core.Interpolation{Input: []float64{0.1, 0.9}, Output: []float64{0, -90}, Easing: bezier(core.EaseInOut)})
			back := core.Interpolate(p, core.Interpolation{Input: []float64{0.1, 0.9}, Output: []float64{90, 0}, Easing: bezier(core.EaseInOut)})
			return core.View{Style: stageFill, Children: []core.Node{
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					Animated: &core.AnimatedStyle{RotateY: back}}, Children: []core.Node{camScene(core.RGB(0x5a2b3c), core.RGB(0x2a1420), camLayers{})}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					Animated: &core.AnimatedStyle{RotateY: turn}}, Children: []core.Node{camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{})}},
			}}
		}})
}

// #140 Page curl: a corner lifts, folds over, and the next page is under it.
func m2Curl(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "140", Title: "Page curl", Duration: 2400 * ms, Accent: colAmber, Height: 110,
		Note: "The sheet is one rectangle with a corner of shadow; the fold is a rotation and a shear we cannot do, so it lifts as a whole.",
		Stage: func(p *core.AnimatedValue) core.Node {
			lift := core.Interpolate(p, core.Interpolation{Input: []float64{0.15, 0.75}, Output: []float64{0, 1}, Easing: bezier(core.EaseInOut)})
			return core.View{Style: stageFill, Children: []core.Node{
				camScene(core.RGB(0x5a2b3c), core.RGB(0x2a1420), camLayers{}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Right: core.Pt(0), Width: core.Pt(260), Bottom: core.Pt(0),
					Overflow: core.Hidden}, Children: []core.Node{
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
						Background: core.RGB(0xf2e6cf),
						Animated: &core.AnimatedStyle{TranslateX: core.Interpolate(lift, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 240}}),
							Rotate:  core.Interpolate(lift, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, -8}}),
							ScaleY:  core.Interpolate(lift, core.Interpolation{Input: []float64{0, 1}, Output: []float64{1, 0.86}}),
							Opacity: core.Interpolate(lift, core.Interpolation{Input: []float64{0.7, 1}, Output: []float64{1, 0.2}})}},
						Children: []core.Node{
							core.View{Style: core.Style{Padding: core.Edges{Top: 16, Left: 14}}, Children: []core.Node{
								Word("page one", 12, core.RGB(0x3a3324), nil),
							}},
						}},
				}},
			}}
		}})
}

// #141 Pixelate: the picture coarsens, switches at the worst of it, then sharpens.
func m2Pixelate(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "141", Title: "Pixelate", Duration: 2200 * ms, Accent: colViolet, Height: 110,
		Note: "Blocks are how a mosaic reads: big squares, one colour each, growing until the picture is only a suggestion.",
		Stage: func(p *core.AnimatedValue) core.Node {
			coarse := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.45}, Output: []float64{0.05, 1}})
			swap := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.48, 0.52, 0.8}, Output: []float64{0, 0, 1, 0}})
			tints := []core.Color{core.RGB(0x1b3a5c), core.RGB(0x24486b), core.RGB(0x2f5a7a), core.RGB(0x3a2b3c), core.RGB(0x5a2b3c), core.RGB(0x243040)}
			var blocks []core.Node
			for i := 0; i < 24; i++ {
				col := tints[(i*7)%len(tints)]
				if i%3 == 0 {
					col = core.RGB(0x14202e)
				}
				blocks = append(blocks, core.View{Style: core.Style{Position: core.Absolute,
					Top: core.Pt(core.Points(float32(i/6) * 22)), Left: core.Pt(core.Points(float32(i%6) * 44)),
					Width: core.Pt(44), Height: core.Pt(22), Background: col,
					Animated: &core.AnimatedStyle{Scale: core.Interpolate(coarse, core.Interpolation{Input: []float64{0.05, 1}, Output: []float64{0.05, 1}}),
						Opacity: coarse, Background: core.Interpolate(swap, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1},
							Colors: []core.Color{col, core.RGB(0x5a2b3c)}})}}})
			}
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0), Overflow: core.Hidden},
				Children: blocks}
		}})
}

// #142 Shatter: the frame cracks and the shards leave together.
func m2Shatter(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "142", Title: "Shatter", Duration: 2000 * ms, Accent: colCyan, Height: 110,
		Note: "Six shards of the same picture, each flung along its own radius: the cut is hidden at the moment they separate.",
		Stage: func(p *core.AnimatedValue) core.Node {
			burst := core.Interpolate(p, core.Interpolation{Input: []float64{0.35, 0.9}, Output: []float64{0, 1}, Easing: bezier(core.EaseIn)})
			var shards []core.Node
			offsets := [][2]float64{{-1, -1}, {1, -1}, {-1, 1}, {1, 1}, {0, -1.4}, {0, 1.4}}
			for i, o := range offsets {
				tint := []core.Color{core.RGB(0x1b3a5c), core.RGB(0x24486b), core.RGB(0x14202e), core.RGB(0x2f5a7a), core.RGB(0x1b3a5c), core.RGB(0x24486b)}[i]
				shards = append(shards, core.View{Style: core.Style{Position: core.Absolute,
					Top: core.Pt(core.Points(float32(26 + 30*o[1]))), Left: core.Pt(core.Points(float32(60 + 52*o[0]))),
					Width: core.Pt(70), Height: core.Pt(50), Background: tint,
					Animated: &core.AnimatedStyle{
						TranslateX: core.Interpolate(burst, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 90 * o[0]}}),
						TranslateY: core.Interpolate(burst, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 80 * o[1]}}),
						Rotate:     core.Interpolate(burst, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 22 * o[0]}}),
						Opacity:    core.Interpolate(p, core.Interpolation{Input: []float64{0.7, 1}, Output: []float64{1, 0}}),
					}}})
			}
			return core.View{Style: stageFill, Children: []core.Node{
				camScene(core.RGB(0x5a2b3c), core.RGB(0x2a1420), camLayers{}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					Overflow: core.Hidden}, Children: shards},
			}}
		}})
}

// #143 Exploded view: the layers of the object separate and come back.
func m2Exploded(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "143", Title: "Exploded view", Duration: 2600 * ms, Accent: colGreen, Height: 110,
		Note: "Four layers, one axis: what makes an exploded view readable is that nothing moves in any other direction.",
		Stage: func(p *core.AnimatedValue) core.Node {
			spread := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.5}, Output: []float64{0, 1}, Easing: bezier(core.EaseInOut)})
			layers := []core.Color{core.RGB(0x2f6bf0), core.RGB(0x35d6de), core.RGB(0x39d98a), core.RGB(0xffb347)}
			var parts []core.Node
			for i, col := range layers {
				gap := float64(i) - 1.5
				parts = append(parts, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(40), Left: core.Pt(70),
					Animated: &core.AnimatedStyle{TranslateY: core.Interpolate(spread, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 12 * gap}}),
						TranslateX: core.Interpolate(spread, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 10 * gap}})}},
					Children: []core.Node{core.View{Style: core.Style{Width: core.Pt(56), Height: core.Pt(12), CornerRadius: 3,
						Background: col, BorderWidth: 1, BorderColor: core.RGB(0x0c1119)}}}})
			}
			return uipScreen(170, 96, parts...)
		}})
}

// #144 Split screen: one picture becomes two, then three, each alive.
func m2Split(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "144", Title: "Split screen", Duration: 2400 * ms, Accent: colBlue, Height: 110,
		Note: "The frame divides rather than the scene: each panel keeps its own moving picture, which is what makes it read as three shots.",
		Stage: func(p *core.AnimatedValue) core.Node {

			tints := []core.Color{core.RGB(0x1b3a5c), core.RGB(0x5a2b3c), core.RGB(0x1d4a3f)}
			var panels []core.Node
			for i, col := range tints {
				// The panel keeps its width and grows across it; splitting the
				// frame is the same gesture three times, a beat apart.
				grow := core.Interpolate(p, core.Interpolation{Input: []float64{0.1 + 0.2*float64(i), 0.5 + 0.2*float64(i)}, Output: []float64{0.05, 1}})
				panels = append(panels, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0),
					Left: core.Pt(core.Points(float32(i * 176))), Width: core.Pt(176), Height: core.Pt(110),
					Background: col, Overflow: core.Hidden,
					Animated: &core.AnimatedStyle{ScaleX: grow}},
					Children: []core.Node{
						core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Width: core.Pt(176), Height: core.Pt(110),
							Animated: &core.AnimatedStyle{TranslateX: lerp(p, float64(-20+i*10), float64(20-i*10))}},
							Children: []core.Node{camScene(col, core.RGB(0x0e1420), camLayers{})}},
					}})
			}
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0), Overflow: core.Hidden},
				Children: panels}
		}})
}
