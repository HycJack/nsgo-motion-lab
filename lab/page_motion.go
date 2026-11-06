package main

// page_motion.go is the motion pack (motion-prompt #001–#024): entrances,
// motion polish, transitions and highlights, on a bare stage.
//
// Catalog names and timings: github.com/showheyohtaki/motion-prompt (MIT).
// The code is NSGo's own.

import (
	"time"

	"github.com/oneo-me/nsgo/core"
)

// pageMotion lays the pack out.
func pageMotion(cx *core.Scope) core.Node {
	return core.View{Children: []core.Node{
		core.View{Style: core.Style{Padding: core.Edges{Top: 14, Right: 18, Left: 18}}, Children: []core.Node{
			note("Motion pack 001–024 of motion-prompt (github.com/showheyohtaki/motion-prompt): entrances, polish, transitions. Code is NSGo's own."),
		}},
		grid(
			moPopIn(cx),
			moStagger(cx),
			moBurst(cx),
			moLineDraw(cx),
			moMaskReveal(cx),
			moTypewriter(cx),
			moCountUp(cx),
			moKinetic(cx),
			moEasing(cx),
			moOvershoot(cx),
			moSquash(cx),
			moMorph(cx),
			moLiquid(cx),
			moTrail(cx),
			moParticles(cx),
			moParallax(cx),
			moMaskTransition(cx),
			moShapeTransition(cx),
			moLiquidTransition(cx),
			moZoomThrough(cx),
			moFlip(cx),
			moLightSweep(cx),
			moGlitch(cx),
			moImpact(cx),
		),
	}}
}

// #001 Pop-in: the icon overshoots, settles, then the badge pops.
func moPopIn(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "001", Title: "Pop-in", Duration: 1600 * ms, Accent: colCyan,
		Note: "The icon arrives with an overshoot and a settle; the badge is a separate beat, a moment later.",
		Stage: func(p *core.AnimatedValue) core.Node {
			icon := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.35, 0.55, 0.75}, Output: []float64{0.4, 1.16, 0.94, 1}})
			badge := core.Interpolate(p, core.Interpolation{Input: []float64{0.55, 0.72, 0.85, 1}, Output: []float64{0, 1.3, 0.9, 1}})
			return uipScreen(120, 96,
				Centered(core.View{Style: core.Style{Position: core.Relative}, Children: []core.Node{
					core.View{Style: core.Style{Width: core.Pt(44), Height: core.Pt(44), CornerRadius: 12, Background: core.RGB(0x2f6bf0),
						Animated: &core.AnimatedStyle{Scale: icon}}},
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(-6), Left: core.Pt(32),
						Animated: &core.AnimatedStyle{Scale: badge}},
						Children: []core.Node{Disc(16, core.RGB(0xff4d6d), nil)}},
				}}))
		}})
}

// #002 Stagger: bars rise one after another, each a tenth of the pass later.
func moStagger(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "002", Title: "Stagger", Duration: 1400 * ms, Accent: colCyan,
		Note: "Five bars, five starts: the delay is what makes the group read as a row instead of a wall.",
		Stage: func(p *core.AnimatedValue) core.Node {
			heights := []float32{28, 42, 34, 54, 46}
			var bars []core.Node
			for i, h := range heights {
				from := float64(i) * 0.1
				rise := core.Interpolate(p, core.Interpolation{
					Input: []float64{from, from + 0.35, 1}, Output: []float64{0, 1, 1},
				})
				bars = append(bars, core.View{Style: core.Style{Width: core.Pt(18), Height: core.Pt(core.Points(h)),
					CornerRadius: 4, Background: core.RGB(0x1d3a52),
					Animated: &core.AnimatedStyle{
						Background: core.Interpolate(rise, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1},
							Colors: []core.Color{core.RGB(0x1d3a52), core.RGB(0x35d6de)}}),
						ScaleY:     rise,
						TranslateY: core.Interpolate(rise, core.Interpolation{Input: []float64{0, 1}, Output: []float64{float64(h) / 2, 0}}),
					}}})
			}
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
				AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter}, Children: []core.Node{row(8, bars...)}}
		}})
}

// #003 Burst: the heart fills, and a ring and dots leave it at once.
func moBurst(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "003", Title: "Burst", Duration: 1500 * ms, Accent: colPink,
		Note: "The fill lands first; the ring, the eight lines and the dots all leave on that same beat.",
		Stage: func(p *core.AnimatedValue) core.Node {
			out := core.Interpolate(p, core.Interpolation{Input: []float64{0.3, 1}, Output: []float64{0, 1}})
			fill := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.3}, Output: []float64{0.6, 1}})
			var dots []core.Node
			for i := 0; i < 8; i++ {
				dots = append(dots, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(40), Left: core.Pt(52),
					Animated: &core.AnimatedStyle{
						TranslateX: core.Interpolate(out, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 46 * circleX(i)}}),
						TranslateY: core.Interpolate(out, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 46 * circleY(i)}}),
						Opacity:    core.Interpolate(out, core.Interpolation{Input: []float64{0, 0.7, 1}, Output: []float64{1, 0.8, 0}}),
					}}, Children: []core.Node{Disc(5, colPink, nil)}})
			}
			return core.View{Style: stageFill, Children: append([]core.Node{
				Centered(core.View{Style: core.Style{Position: core.Relative, Width: core.Pt(120), Height: core.Pt(96)}, Children: []core.Node{
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(36), Left: core.Pt(46),
						Width: core.Pt(28), Height: core.Pt(26), CornerRadius: 13, Background: colPink,
						Animated: &core.AnimatedStyle{Scale: fill}}},
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(30), Left: core.Pt(40),
						Animated: &core.AnimatedStyle{Scale: core.Interpolate(out, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0.6, 2.2}}),
							Opacity: core.Interpolate(out, core.Interpolation{Input: []float64{0, 1}, Output: []float64{1, 0}})}},
						Children: []core.Node{Ring(40, 2, colPink, nil)}},
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0)}, Children: dots},
				}}),
			}, nil...)}
		}})
}

// circleX and circleY place the eight burst dots on a circle.
func circleX(i int) float64 { return []float64{1, 0.7, 0, -0.7, -1, -0.7, 0, 0.7}[i] }
func circleY(i int) float64 { return []float64{0, 0.7, 1, 0.7, 0, -0.7, -1, -0.7}[i] }

// #004 Line draw: the line draws on, then is erased from where it started.
func moLineDraw(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "004", Title: "Line draw", Duration: 2400 * ms, Accent: colGreen,
		Note: "A line can be drawn from one end and erased back to the same end; the dot rides the tip either way.",
		Stage: func(p *core.AnimatedValue) core.Node {
			grow := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.6, 1}, Output: []float64{0, 1, 0}})
			return uipScreen(180, 96,
				At(20, 10, uipBar(180, 1, core.RGB(0x1b2534), nil)),
				At(58, 10, uipBar(180, 1, core.RGB(0x1b2534), nil)),
				At(10, 10, core.View{Style: core.Style{Width: core.Pt(180), Height: core.Pt(3), CornerRadius: 2, Background: colGreen,
					Animated: &core.AnimatedStyle{ScaleX: grow, TranslateX: lerp(grow, 0, 90)}}}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(54), Left: core.Pt(10),
					Animated: &core.AnimatedStyle{TranslateX: lerp(grow, 0, 180)}},
					Children: []core.Node{Disc(7, core.RGB(0xeaf1fa), nil)}},
			)
		}})
}

// #005 Mask reveal: the headline slides out from behind a hard edge.
func moMaskReveal(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "005", Title: "Mask reveal", Duration: 1600 * ms, Accent: colViolet,
		Note: "Nothing fades: the words are simply below a hard edge and slide up through it.",
		Stage: func(p *core.AnimatedValue) core.Node {
			rise := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.7, 1}, Output: []float64{30, 0, 0}})
			return uipScreen(180, 96,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(30), Left: core.Pt(10), Width: core.Pt(170), Height: core.Pt(28),
					Overflow: core.Hidden}, Children: []core.Node{
					core.View{Style: core.Style{Animated: &core.AnimatedStyle{TranslateY: rise}}, Children: []core.Node{
						Word("Big news", 22, colText, nil),
					}},
				}},
				At(58, 10, uipBar(170, 1, core.RGBA(0xa2b3c9, 0.35), nil)),
				At(70, 10, uipBar(120, 6, core.RGB(0x1d2836), nil)),
			)
		}})
}

// #006 Typewriter: an input fills a character at a time, then the send button lights.
func moTypewriter(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "006", Title: "Typewriter", Duration: 2400 * ms, Accent: colGreen,
		Note: "A caret, a character every sixth of the pass, and the send button only lights once the line is complete.",
		Stage: func(p *core.AnimatedValue) core.Node {
			msg := core.NewSignal(cx, "")
			full := "Ship it today"
			cx.AfterFunc(60*ms, func() {
				for i := 1; i <= len(full); i++ {
					at := i
					cx.AfterFunc(time.Duration(at)*90*ms, func() { msg.Set(full[:at]) })
				}
			})
			send := core.Interpolate(p, core.Interpolation{Input: []float64{0.8, 0.95}, Output: []float64{0, 1}})
			return uipScreen(190, 96,
				At(26, 8, core.View{Style: core.Style{Width: core.Pt(174), Height: core.Pt(26), CornerRadius: 7,
					BorderWidth: 1, BorderColor: core.RGB(0x2a3644), Padding: core.Edges{Left: 8}, Direction: core.Row, AlignItems: core.AlignCenter}, Children: []core.Node{
					core.Text{ContentFunc: func() string { return msg.Get() }, Style: core.Style{FontSize: 10, Color: colText}},
					core.View{Style: core.Style{Width: core.Pt(4), Height: core.Pt(12), Background: colGreen,
						Animated: &core.AnimatedStyle{Opacity: core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.5, 0.51, 1}, Output: []float64{1, 1, 0, 1}})}}},
				}}),
				At(58, 130, core.View{Style: core.Style{Width: core.Pt(52), Height: core.Pt(20), CornerRadius: 6, Background: core.RGB(0x1d2836),
					AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter,
					Animated: &core.AnimatedStyle{Background: core.Interpolate(send, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1},
						Colors: []core.Color{core.RGB(0x1d2836), core.RGB(0x39d98a)}})}}, Children: []core.Node{
					uipLabel("send", 9, core.RGB(0x0a0e14)),
				}}),
			)
		}})
}

// #007 Count-up: the number races and locks onto its value.
func moCountUp(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "007", Title: "Count-up", Duration: 1800 * ms, Accent: colAmber,
		Note: "The value is a signal, not a re-render: one text node, its content read from the number as it runs.",
		Stage: func(p *core.AnimatedValue) core.Node {
			n := core.NewSignal(cx, 0.0)
			cx.AfterFunc(30*ms, func() {
				cx.Animate(1500*ms, core.EaseOut, func(t float64) { n.Set(4820 * t) })
			})
			bar := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.9, 1}, Output: []float64{0, 1, 1}})
			return uipScreen(180, 96,
				At(24, 10, core.Text{ContentFunc: func() string { return formatCount(n.Get()) },
					Style: core.Style{FontSize: 28, FontWeight: core.FontWeightBold, Color: colText}}),
				At(64, 10, uipBar(150, 3, core.RGB(0x1d2836), nil)),
				At(63, 10, core.View{Style: core.Style{Width: core.Pt(150), Height: core.Pt(3), CornerRadius: 2, Background: colAmber,
					Animated: &core.AnimatedStyle{ScaleX: bar, TranslateX: lerp(bar, 0, 75)}}}),
			)
		}})
}

// formatCount prints an integer with thousands separators.
func formatCount(v float64) string {
	n := int(v + 0.5)
	if n < 1000 {
		return itoa(n)
	}
	return itoa(n/1000) + "," + pad3(n%1000)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b []byte
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}

func pad3(n int) string {
	s := itoa(n)
	for len(s) < 3 {
		s = "0" + s
	}
	return s
}

// #008 Kinetic typography: four words, four ways of meaning it.
func moKinetic(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "008", Title: "Kinetic typography", Duration: 2200 * ms, Accent: colViolet, Height: 110,
		Note: "Stretch, bounce, spin, slam: the word does what it says, one beat each.",
		Stage: func(p *core.AnimatedValue) core.Node {
			stretch := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.25, 0.4}, Output: []float64{1, 1.6, 1}})
			bounce := core.Interpolate(p, core.Interpolation{Input: []float64{0.25, 0.4, 0.5, 0.6}, Output: []float64{0, -14, 0, -5}})
			spin := core.Interpolate(p, core.Interpolation{Input: []float64{0.5, 0.85}, Output: []float64{0, 360}})
			slam := core.Interpolate(p, core.Interpolation{Input: []float64{0.85, 0.95, 1}, Output: []float64{0, 1.15, 1}})
			words := []core.Node{
				Word("STRETCH", 13, colCyan, &core.AnimatedStyle{ScaleX: stretch}),
				Word("bounce", 13, colGreen, &core.AnimatedStyle{TranslateY: bounce}),
				Word("spin", 13, colAmber, &core.AnimatedStyle{Rotate: spin}),
				Word("SLAM", 13, colPink, &core.AnimatedStyle{Scale: slam}),
			}
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
				AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter}, Children: []core.Node{row(10, words...)}}
		}})
}

// #009 Easing: same distance, same time, one linear and one eased.
func moEasing(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "009", Title: "Easing", Duration: 2000 * ms, Accent: colCyan,
		Note: "Both squares arrive together; only the one on the bottom spends its time differently.",
		Stage: func(p *core.AnimatedValue) core.Node {
			lin := lerp(p, 0, 140)
			eased := core.Interpolate(p, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 140},
				Easing: bezier(core.EaseInOut)})
			return uipScreen(170, 96,
				At(26, 6, uipBar(150, 1, core.RGB(0x1b2534), nil)),
				At(66, 6, uipBar(150, 1, core.RGB(0x1b2534), nil)),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(20), Left: core.Pt(6),
					Animated: &core.AnimatedStyle{TranslateX: lin}}, Children: []core.Node{Panel(16, 14, colCyan, nil)}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(60), Left: core.Pt(6),
					Animated: &core.AnimatedStyle{TranslateX: eased}}, Children: []core.Node{Panel(16, 14, colViolet, nil)}},
			)
		}})
}

// #010 Overshoot: the needle goes past the mark and comes back.
func moOvershoot(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "010", Title: "Overshoot", Duration: 1800 * ms, Accent: colAmber,
		Note: "The needle swings past the target it is aiming at, comes back and settles — the spring you can see.",
		Stage: func(p *core.AnimatedValue) core.Node {
			ang := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.42, 0.62, 0.8, 1}, Output: []float64{-70, 68, 34, 42, 40}})
			return uipScreen(120, 96,
				Centered(core.View{Style: core.Style{Position: core.Relative, Width: core.Pt(96), Height: core.Pt(90)}, Children: []core.Node{
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(6), Left: core.Pt(6),
						BorderWidth: 2, BorderColor: core.RGB(0x233040), Width: core.Pt(84), Height: core.Pt(84), CornerRadius: 42}},
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(31), Left: core.Pt(31), Width: core.Pt(34), Height: core.Pt(34),
						Animated: &core.AnimatedStyle{Rotate: ang}}, Children: []core.Node{
						core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(16),
							Width: core.Pt(2), Height: core.Pt(17), Background: colAmber}},
					}},
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(44), Left: core.Pt(44)},
						Children: []core.Node{Disc(9, colText, nil)}},
				}}))
		}})
}

// #011 Squash & stretch: the ball deforms on the two frames that matter.
func moSquash(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "011", Title: "Squash & stretch", Duration: 1600 * ms, Accent: colGreen,
		Note: "Down fast, squashed flat for an instant, up stretched tall: the two frames that make a ball feel heavy.",
		Stage: func(p *core.AnimatedValue) core.Node {
			y := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.35, 0.42, 0.7, 0.8, 1}, Output: []float64{-40, 0, 0, 0, 0, 0}})
			sy := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.34, 0.44, 0.6, 0.72, 1}, Output: []float64{1.1, 1, 0.66, 1.2, 1, 1}})
			sx := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.34, 0.44, 0.6, 0.72, 1}, Output: []float64{0.9, 1, 1.4, 0.86, 1, 1}})
			return uipScreen(140, 96,
				At(86, 10, uipBar(120, 2, core.RGB(0x1b2534), nil)),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(40), Left: core.Pt(52),
					Animated: &core.AnimatedStyle{TranslateY: y}}, Children: []core.Node{
					core.View{Style: core.Style{Width: core.Pt(26), Height: core.Pt(26), CornerRadius: 13, Background: colGreen,
						Animated: &core.AnimatedStyle{ScaleX: sx, ScaleY: sy}}},
				}},
			)
		}})
}

// #012 Shape morph: circle, square, diamond, and the colour comes along.
func moMorph(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "012", Title: "Shape morph", Duration: 2200 * ms, Accent: colViolet,
		Note: "There is no path here: the corner radius is fixed and the shape turns, but the colour carries the change.",
		Stage: func(p *core.AnimatedValue) core.Node {
			rot := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.5, 1}, Output: []float64{0, 45, 135}})
			col := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.5, 1}, Output: []float64{0, 1, 2},
				Colors: []core.Color{colViolet, colCyan, colPink}})
			round := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.5, 1}, Output: []float64{1, 0, 0.5}})
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
				AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter}, Children: []core.Node{
				core.View{Style: core.Style{Width: core.Pt(40), Height: core.Pt(40), Background: colViolet,
					Animated: &core.AnimatedStyle{Background: col, Rotate: rot, Scale: round}}},
			}}
		}})
}

// #013 Liquid motion: a blob stretches, splits into drops, and comes back.
func moLiquid(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "013", Title: "Liquid motion", Duration: 2200 * ms, Accent: colCyan,
		Note: "One body becomes three and then one again; the drops are only convincing while they stay in line.",
		Stage: func(p *core.AnimatedValue) core.Node {
			split := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.45, 0.75, 1}, Output: []float64{0, 1, 1, 0}})
			stretch := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.35, 0.6}, Output: []float64{1, 1.7, 1}})
			var drops []core.Node
			for i, d := range []float64{-1, 1} {
				_ = i
				drops = append(drops, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(40), Left: core.Pt(52),
					Animated: &core.AnimatedStyle{TranslateX: core.Interpolate(split, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 26 * d}}),
						Opacity: split, Scale: core.Interpolate(split, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0.4, 1}})}},
					Children: []core.Node{Disc(14, colCyan, nil)}})
			}
			return core.View{Style: stageFill, Children: append([]core.Node{
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(36), Left: core.Pt(48),
					Animated: &core.AnimatedStyle{ScaleX: stretch}}, Children: []core.Node{Disc(22, colCyan, nil)}},
			}, drops...)}
		}})
}

// #014 Motion trail: the ball leads, the copies follow behind it.
func moTrail(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "014", Title: "Motion trail", Duration: 2400 * ms, Accent: colViolet,
		Note: "An afterimage is the same motion, sampled later: each copy is the lead's position a few frames ago.",
		Stage: func(p *core.AnimatedValue) core.Node {
			var ghosts []core.Node
			for i := 5; i >= 1; i-- {
				lag := core.Interpolate(p, core.Interpolation{
					Input:  []float64{0, 1},
					Output: []float64{0, 1},
				})
				_ = lag
				ghosts = append(ghosts, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(40), Left: core.Pt(10),
					Animated: &core.AnimatedStyle{
						TranslateX: core.Interpolate(p, core.Interpolation{
							Input:  []float64{0, 0.25, 0.5, 0.75, 1},
							Output: []float64{0, 120, 0, 120, 0},
							Easing: bezier(core.EaseInOut),
						}),
						TranslateY: core.Interpolate(p, core.Interpolation{
							Input:  []float64{0, 0.25, 0.5, 0.75, 1},
							Output: []float64{0, 30 * float64(i-3), 0, 30 * float64(i-3), 0},
						}),
						Opacity: core.Interpolate(p, core.Interpolation{Input: []float64{0, 1}, Output: []float64{1, 0.08 * float64(6-i)}}),
					}}, Children: []core.Node{Disc(11, colViolet, nil)}})
			}
			return uipScreen(170, 96, append([]core.Node{
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(40), Left: core.Pt(10),
					Animated: &core.AnimatedStyle{
						TranslateX: core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.25, 0.5, 0.75, 1}, Output: []float64{0, 120, 0, 120, 0}, Easing: bezier(core.EaseInOut)}),
						TranslateY: core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.25, 0.5, 0.75, 1}, Output: []float64{0, 0, 0, 0, 0}}),
					}}, Children: []core.Node{Disc(13, colText, nil)}},
			}, ghosts...)...)
		}})
}

// #015 Particles: an emitter, gravity, and a spread that keeps going.
func moParticles(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "015", Title: "Particles", Duration: 1800 * ms, Accent: colAmber, Loop: true,
		Note: "Six emitters, six delays, one gravity: the loop is what turns a single burst into a fountain.",
		Stage: func(p *core.AnimatedValue) core.Node {
			var parts []core.Node
			spread := []float64{-1, -0.6, -0.2, 0.2, 0.6, 1}
			for i, s := range spread {
				delay := 0.07 * float64(i%3)
				up := core.Interpolate(p, core.Interpolation{Input: []float64{delay, 0.55 + delay, 1}, Output: []float64{0, 1, 0}})
				parts = append(parts, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(74), Left: core.Pt(82),
					Animated: &core.AnimatedStyle{
						TranslateX: core.Interpolate(up, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 34 * s}}),
						TranslateY: core.Interpolate(up, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, -52}}),
						Opacity:    core.Interpolate(p, core.Interpolation{Input: []float64{delay, 0.8, 1}, Output: []float64{1, 1, 0}}),
					}}, Children: []core.Node{Disc(5, colAmber, nil)}})
			}
			return uipScreen(170, 96, append([]core.Node{
				At(78, 80, Disc(8, core.RGB(0x3a2d18), nil)),
			}, parts...)...)
		}})
}

// #016 Parallax: one sideways pan, four speeds.
func moParallax(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "016", Title: "Parallax", Duration: 2600 * ms, Accent: colGreen, Height: 110,
		Note: "Nothing about the layers is 3D; the illusion is that they slide by four different amounts.",
		Stage: func(p *core.AnimatedValue) core.Node {
			return camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{
				Sun:    &core.AnimatedStyle{TranslateX: lerp(p, 26, -26)},
				Far:    &core.AnimatedStyle{TranslateX: lerp(p, 60, -60)},
				Near:   &core.AnimatedStyle{TranslateX: lerp(p, 130, -130)},
				Figure: &core.AnimatedStyle{TranslateX: lerp(p, 200, -200)},
			})
		}})
}

// #017 Mask transition: a window opens on the next scene, growing from a point.
func moMaskTransition(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "017", Title: "Mask transition", Duration: 2200 * ms, Accent: colBlue, Height: 110,
		Note: "The new scene is already there, hidden; a circular window opens on it until only a corner of the old one is left.",
		Stage: func(p *core.AnimatedValue) core.Node {
			open := core.Interpolate(p, core.Interpolation{Input: []float64{0.1, 0.85}, Output: []float64{0.04, 1.6}})
			return core.View{Style: stageFill, Children: []core.Node{
				camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter}, Children: []core.Node{
					core.View{Style: core.Style{Width: core.Pt(150), Height: core.Pt(150), Overflow: core.Hidden,
						Animated: &core.AnimatedStyle{Scale: open}},
						Children: []core.Node{core.View{Style: core.Style{Width: core.Pt(150), Height: core.Pt(150)}, Children: []core.Node{
							camScene(core.RGB(0x5a2b3c), core.RGB(0x2a1420), camLayers{}),
						}}}},
				}},
			}}
		}})
}

// #018 Shape transition: bands sweep across and the next scene is behind them.
func moShapeTransition(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "018", Title: "Shape transition", Duration: 2000 * ms, Accent: colAmber, Height: 110,
		Note: "Five bands, five small delays: the stagger is the whole effect, and the cut happens under the cover.",
		Stage: func(p *core.AnimatedValue) core.Node {
			var bands []core.Node
			tints := []core.Color{colAmber, colPink, colViolet, colBlue, colCyan}
			for i := 0; i < 5; i++ {
				from := 0.06 * float64(i)
				x := core.Interpolate(p, core.Interpolation{Input: []float64{from, 0.42 + from, 0.58 + from, 1}, Output: []float64{-60, 40, 40, 200}})
				bands = append(bands, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(core.Points(float32(i) * 22)), Left: core.Pt(0),
					Width: core.Pt(58), Height: core.Pt(24), Background: tints[i],
					Animated: &core.AnimatedStyle{TranslateX: x}}})
			}
			for i := range bands {
				_ = i
			}
			return core.View{Style: stageFill, Children: []core.Node{
				camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					Overflow: core.Hidden}, Children: bands},
			}}
		}})
}

// #019 Liquid transition: drips come down and fill the frame with the next scene.
func moLiquidTransition(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "019", Title: "Liquid transition", Duration: 2400 * ms, Accent: colPink, Height: 110,
		Note: "Four drips at four lengths; they meet at the bottom and the flood takes the rest of the frame.",
		Stage: func(p *core.AnimatedValue) core.Node {
			flood := core.Interpolate(p, core.Interpolation{Input: []float64{0.55, 1}, Output: []float64{0, 110}})
			var drips []core.Node
			for i, left := range []float32{10, 58, 106, 168} {
				delay := 0.07 * float64(i)
				k := core.Interpolate(p, core.Interpolation{Input: []float64{delay, 0.6 + delay}, Output: []float64{0.02, 1}})
				drips = append(drips, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(core.Points(left)),
					Width: core.Pt(26), Height: core.Pt(96), CornerRadius: 13, Background: colPink,
					// A drip grows downward from a fixed top: scale about the
					// centre, then push back by half what was lost.
					Animated: &core.AnimatedStyle{ScaleY: k,
						TranslateY: core.Interpolate(k, core.Interpolation{Input: []float64{0.02, 1}, Output: []float64{-47, 0}})}}})
			}
			return core.View{Style: stageFill, Children: []core.Node{
				camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					Overflow: core.Hidden}, Children: append(drips, core.View{Style: core.Style{
					Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					Animated: &core.AnimatedStyle{TranslateY: core.Interpolate(flood, core.Interpolation{Input: []float64{0, 110}, Output: []float64{-110, 0}})},
				}, Children: []core.Node{camScene(core.RGB(0x5a2b3c), core.RGB(0x2a1420), camLayers{})}})},
			}}
		}})
}

// #020 Zoom through: the hole of the O opens and becomes the next scene.
func moZoomThrough(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "020", Title: "Zoom through", Duration: 2200 * ms, Accent: colCyan, Height: 110,
		Note: "The ring rushes past the camera; what was a hole in a letter is a doorway by the time we are through it.",
		Stage: func(p *core.AnimatedValue) core.Node {
			ring := core.Interpolate(p, core.Interpolation{Input: []float64{0.15, 0.9}, Output: []float64{1, 9}})
			inner := core.Interpolate(p, core.Interpolation{Input: []float64{0.45, 1}, Output: []float64{0.2, 1.4}})
			return core.View{Style: stageFill, Children: []core.Node{
				camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter}, Children: []core.Node{
					core.View{Style: core.Style{Width: core.Pt(60), Height: core.Pt(60), Overflow: core.Hidden, CornerRadius: 30,
						Animated: &core.AnimatedStyle{Scale: inner}}, Children: []core.Node{
						core.View{Style: core.Style{Width: core.Pt(60), Height: core.Pt(60)}, Children: []core.Node{
							camScene(core.RGB(0x5a2b3c), core.RGB(0x2a1420), camLayers{}),
						}}}},
				}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter}, Children: []core.Node{
					Ring(52, 16, core.RGB(0x0a0e14), &core.AnimatedStyle{Scale: ring}),
				}},
			}}
		}})
}

// #021 Flip: the card turns over around its vertical axis.
func moFlip(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "021", Title: "Flip", Duration: 1800 * ms, Accent: colViolet,
		Note: "One rotation, two faces: the back is only visible after the halfway point.",
		Stage: func(p *core.AnimatedValue) core.Node {
			turn := core.Interpolate(p, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 180}})
			back := core.Interpolate(p, core.Interpolation{Input: []float64{0.45, 0.56}, Output: []float64{0, 1}})
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
				AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter}, Children: []core.Node{
				core.View{Style: core.Style{Width: core.Pt(64), Height: core.Pt(44), CornerRadius: 8, Background: colViolet,
					AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter,
					Animated: &core.AnimatedStyle{RotateY: turn, Background: core.Interpolate(back, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1}, Colors: []core.Color{colViolet, colCyan}})}},
					Children: []core.Node{core.Text{Content: "A", Style: core.Style{FontSize: 18, FontWeight: core.FontWeightBold, Color: core.RGB(0x0a0e14),
						Animated: &core.AnimatedStyle{Opacity: core.Interpolate(back, core.Interpolation{Input: []float64{0, 0.4}, Output: []float64{1, 0}})}}},
						core.Text{Content: "B", Style: core.Style{Position: core.Absolute, FontSize: 18, FontWeight: core.FontWeightBold, Color: core.RGB(0x0a0e14),
							Animated: &core.AnimatedStyle{Opacity: back}}},
					}},
			}}
		}})
}

// #022 Light sweep: a band of light crosses the badge at an angle.
func moLightSweep(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "022", Title: "Light sweep", Duration: 1600 * ms, Accent: colAmber,
		Note: "A highlight that never moves the object: only the band of light travels, and only while it is over the badge.",
		Stage: func(p *core.AnimatedValue) core.Node {
			x := core.Interpolate(p, core.Interpolation{Input: []float64{0.15, 0.85}, Output: []float64{-70, 140}, Easing: bezier(core.EaseInOut)})
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
				AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter}, Children: []core.Node{
				core.View{Style: core.Style{Width: core.Pt(70), Height: core.Pt(48), CornerRadius: 10, Background: core.RGB(0x24303f),
					Overflow: core.Hidden}, Children: []core.Node{
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(-20), Width: core.Pt(18), Height: core.Pt(90),
						Background: core.RGB(0xffffff),
						Animated: &core.AnimatedStyle{TranslateX: x, Rotate: lerp(p, 20, 20),
							Opacity: core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.3, 0.7, 1}, Output: []float64{0, 0.55, 0.55, 0}}),
							ScaleX:  core.Interpolate(p, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0.4, 1}})}}},
				}},
			}}
		}})
}

// #023 Glitch: slices shift, colours split, bars appear.
func moGlitch(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "023", Title: "Glitch", Duration: 1400 * ms, Accent: colBlue,
		Note: "The word stays put and its slices do not: two copies split off at the colour fringes, and bars of noise cut through.",
		Stage: func(p *core.AnimatedValue) core.Node {
			jit := core.Interpolate(p, core.Interpolation{
				Input:  []float64{0, 0.2, 0.32, 0.44, 0.56, 0.68, 0.8, 1},
				Output: []float64{0, -6, 5, -3, 7, -4, 2, 0},
			})
			split := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.3, 0.7, 1}, Output: []float64{0, 3, 3, 0}})
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
				AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter}, Children: []core.Node{
				core.View{Style: core.Style{Position: core.Absolute, Width: core.Pt(170), Height: core.Pt(50)}, Children: []core.Node{
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(14), Left: core.Pt(30),
						Animated: &core.AnimatedStyle{TranslateX: jit}}, Children: []core.Node{Word("SIGNAL", 20, core.RGB(0x35d6de), nil)}},
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(14), Left: core.Pt(30),
						Animated: &core.AnimatedStyle{TranslateX: split, Opacity: core.Interpolate(split, core.Interpolation{Input: []float64{0, 3}, Output: []float64{0, 0.7}})}},
						Children: []core.Node{Word("SIGNAL", 20, core.RGB(0xff4d6d), nil)}},
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(14), Left: core.Pt(30),
						Animated: &core.AnimatedStyle{TranslateX: core.Interpolate(split, core.Interpolation{Input: []float64{0, 3}, Output: []float64{0, -3}}),
							Opacity: core.Interpolate(split, core.Interpolation{Input: []float64{0, 3}, Output: []float64{0, 0.7}})}},
						Children: []core.Node{Word("SIGNAL", 20, core.RGB(0x39d98a), nil)}},
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(6), Left: core.Pt(20), Width: core.Pt(130), Height: core.Pt(3),
						Background: core.RGB(0xeaf1fa), Animated: &core.AnimatedStyle{TranslateX: jit, Opacity: split}}},
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(38), Left: core.Pt(50), Width: core.Pt(110), Height: core.Pt(2),
						Background: core.RGB(0xeaf1fa), Animated: &core.AnimatedStyle{TranslateX: core.Interpolate(jit, core.Interpolation{Input: []float64{-6, 7}, Output: []float64{-9, 11}}), Opacity: split}}},
				}}},
			}
		}})
}

// #024 Impact: the slab lands and the whole frame takes the hit.
func moImpact(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "024", Title: "Impact shake", Duration: 1800 * ms, Accent: colPink, Height: 110,
		Note: "The drop is the easy half; the shake that follows is what makes it land, and it decays.",
		Stage: func(p *core.AnimatedValue) core.Node {
			drop := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.42, 0.5, 1}, Output: []float64{-96, -96, 0, 0}})
			shake := core.Interpolate(p, core.Interpolation{
				Input:  []float64{0, 0.5, 0.56, 0.62, 0.68, 0.74, 0.82, 1},
				Output: []float64{0, 0, 5, -4, 3, -2, 1, 0},
			})
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
				Overflow: core.Hidden, Animated: &core.AnimatedStyle{TranslateY: shake}}, Children: []core.Node{
				camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{}),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(36), Left: core.Pt(24),
					Animated: &core.AnimatedStyle{TranslateY: drop}}, Children: []core.Node{
					core.View{Style: core.Style{Width: core.Pt(150), Height: core.Pt(26), CornerRadius: 4, Background: core.RGB(0xeaf1fa)}},
				}},
			}}
		}})
}
