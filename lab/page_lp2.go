package main

// page_lp2.go is the landing-page interactions pack (motion-prompt #145–#168):
// the twenty-four things a marketing page does when you touch it or scroll it.
//
// Upstream they are driven by scroll position; here each one is driven by the
// stage clock, so the same move can be watched without a mouse.
//
// Catalog names and timings: github.com/showheyohtaki/motion-prompt (MIT).

import (
	"github.com/oneo-me/nsgo/core"
)

// pageLP2 lays the pack out.
func pageLP2(cx *core.Scope) core.Node {
	return core.View{Children: []core.Node{
		core.View{Style: core.Style{Padding: core.Edges{Top: 14, Right: 18, Left: 18}}, Children: []core.Node{
			note("Landing page interactions 145–168 of motion-prompt (github.com/showheyohtaki/motion-prompt): scroll, hover and input, driven by the stage clock here. Code is NSGo's own."),
		}},
		grid(
			lp2Timeline(cx),
			lp2TextReveal(cx),
			lp2Stacking(cx),
			lp2BgChange(cx),
			lp2HideHeader(cx),
			lp2Scrollspy(cx),
			lp2BackToTop(cx),
			lp2ZoomScroll(cx),
			lp2Tilt(cx),
			lp2Spotlight(cx),
			lp2ImageSwap(cx),
			lp2TextRoll(cx),
			lp2HoverReveal(cx),
			lp2FloatingLabel(cx),
			lp2InlineValidation(cx),
			lp2LoadingButton(cx),
			lp2ViewTransition(cx),
			lp2Lightbox(cx),
			lp2Hotspots(cx),
			lp2Countdown(cx),
			lp2ShapeDivider(cx),
			lp2Like(cx),
			lp2ScrollCue(cx),
			lp2BgVideo(cx),
		),
	}}
}

// #145 Timeline on scroll: the line draws and the dots light up behind it.
func lp2Timeline(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "145", Width: 300, Title: "Timeline on scroll", Duration: 2400 * ms, Accent: colCyan, Height: 110,
		Note: "The line grows first and the dots come after it: a dot that lights before the line reaches it breaks the illusion of scrolling.",
		Stage: func(p *core.AnimatedValue) core.Node {
			draw := core.Interpolate(p, core.Interpolation{Input: []float64{0.05, 0.85}, Output: []float64{0, 1}, Easing: bezier(core.EaseInOut)})
			var dots []core.Node
			for i := 0; i < 4; i++ {
				at := 0.1 + 0.2*float64(i)
				on := core.Interpolate(p, core.Interpolation{Input: []float64{at, at + 0.12}, Output: []float64{0, 1}})
				dots = append(dots, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(core.Points(float32(24+i*18) - 3)),
					Left: core.Pt(20), Animated: &core.AnimatedStyle{Scale: core.Interpolate(on, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0.7, 1.15}})}},
					Children: []core.Node{Disc(8, colCyan, &core.AnimatedStyle{Opacity: core.Interpolate(on, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0.25, 1}})})}})
			}
			return lpFrame(false,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(24), Left: core.Pt(23), Width: core.Pt(2), Height: core.Pt(66),
					Background: colWell}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(24), Left: core.Pt(23), Width: core.Pt(2), Height: core.Pt(66),
					Background: colCyan, Animated: &core.AnimatedStyle{ScaleY: draw}}},
				core.View{Children: dots},
			)
		}})
}

// #146 Scroll text reveal: words un-dim one line at a time.
func lp2TextReveal(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "146", Width: 300, Title: "Scroll text reveal", Duration: 2400 * ms, Accent: colText, Height: 110,
		Note: "The text is always there, only its brightness changes: that is why a reveal like this never reflows the paragraph.",
		Stage: func(p *core.AnimatedValue) core.Node {
			var lines []core.Node
			for i := 0; i < 4; i++ {
				at := 0.1 + 0.18*float64(i)
				lit := core.Interpolate(p, core.Interpolation{Input: []float64{at, at + 0.2}, Output: []float64{0, 1}})
				lines = append(lines, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(core.Points(float32(24 + i*16))), Left: core.Pt(18),
					Width: core.Pt(core.Points(float32(200 - (i%2)*30))), Height: core.Pt(6), CornerRadius: 3, Background: colLine,
					Animated: &core.AnimatedStyle{Background: core.Interpolate(lit, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1},
						Colors: []core.Color{colLine, colText}})}}})
			}
			return lpFrame(false, core.View{Children: lines})
		}})
}

// #147 Stacking cards: each card pins and the next one lands on top of it.
func lp2Stacking(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "147", Width: 300, Title: "Stacking cards", Duration: 2600 * ms, Accent: colViolet, Height: 110,
		Note: "Every card stops at the same height and shrinks a little: the stack reads because the ones underneath get smaller and darker.",
		Stage: func(p *core.AnimatedValue) core.Node {
			var cards []core.Node
			for i := 0; i < 3; i++ {
				arrive := 0.25 * float64(i)
				land := core.Interpolate(p, core.Interpolation{Input: []float64{arrive, arrive + 0.25}, Output: []float64{0, 1}, Easing: bezier(core.EaseOut)})
				push := core.Interpolate(p, core.Interpolation{Input: []float64{arrive + 0.2, arrive + 0.45}, Output: []float64{0, 1}})
				cards = append(cards, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(core.Points(float32(30 + i*20))),
					Left: core.Pt(core.Points(float32(30 + i*6))), Width: core.Pt(core.Points(float32(150 - i*12))), Height: core.Pt(60),
					CornerRadius: 7, Background: []core.Color{colWell, colBar, colCard}[i], BorderWidth: 1, BorderColor: colLine,
					Animated: &core.AnimatedStyle{
						TranslateY: core.Interpolate(land, core.Interpolation{Input: []float64{0, 1}, Output: []float64{80, 0}}),
						Scale:      core.Interpolate(push, core.Interpolation{Input: []float64{0, 1}, Output: []float64{1, 0.94}}),
						Opacity:    core.Interpolate(push, core.Interpolation{Input: []float64{0, 1}, Output: []float64{1, 0.7}})}},
					Children: []core.Node{
						core.View{Style: core.Style{Padding: core.Edges{Top: 10, Left: 10}, Gap: 5}, Children: []core.Node{
							lpText(60, colText, nil), lpText(90, colMuted, nil),
						}},
					}})
			}
			return lpFrame(true, core.View{Children: cards})
		}})
}

// #148 Background colour change: the page tint follows the section.
func lp2BgChange(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "148", Width: 300, Title: "Background colour change", Duration: 2600 * ms, Accent: colViolet, Height: 110,
		Note: "Three sections, three tints, no cut: the colour is interpolated across the scroll rather than switched between sections.",
		Stage: func(p *core.AnimatedValue) core.Node {
			tint := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.5, 1}, Output: []float64{0, 1, 2},
				Colors: []core.Color{core.RGB(0x121a26), core.RGB(0x1b1430), core.RGB(0x0f241f)}})
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
				Overflow: core.Hidden, Background: colCard, Animated: &core.AnimatedStyle{Background: tint}}, Children: []core.Node{
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(30), Left: core.Pt(16), Gap: 6}, Children: []core.Node{
					lpText(120, colText, nil),
					lpText(160, colMuted, nil),
					lpButton(48, colViolet, nil),
				}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(34), Left: core.Pt(0), Right: core.Pt(0),
					AlignItems: core.AlignCenter, Animated: &core.AnimatedStyle{TranslateY: lerp(p, 0, -36)}}, Children: []core.Node{
					core.View{Style: core.Style{Width: core.Pt(120), Height: core.Pt(40), CornerRadius: 6, Background: colWell}},
				}},
			}}
		}})
}

// #149 Hide-on-scroll header: it leaves going up and returns going down.
func lp2HideHeader(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "149", Width: 300, Title: "Hide-on-scroll header", Duration: 2600 * ms, Accent: colCyan, Height: 110,
		Note: "Down hides it, up brings it back: the header follows the direction of the scroll, not the scroll distance.",
		Stage: func(p *core.AnimatedValue) core.Node {
			hide := core.Interpolate(p, core.Interpolation{Input: []float64{0.1, 0.35, 0.6, 0.85}, Output: []float64{0, -26, -26, 0}, Easing: bezier(core.EaseInOut)})
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
				Overflow: core.Hidden, Background: colCard}, Children: []core.Node{
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(20), Left: core.Pt(16), Right: core.Pt(16), Gap: 6,
					Animated: &core.AnimatedStyle{TranslateY: lerp(p, 0, -30)}}, Children: []core.Node{
					lpText(100, colLine, nil), lpText(140, colLine, nil), lpText(80, colLine, nil), lpText(120, colLine, nil),
				}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Height: core.Pt(20),
					Background: colBar, Direction: core.Row, AlignItems: core.AlignCenter, Padding: core.Edges{Left: 10, Right: 10}, Gap: 6,
					Animated: &core.AnimatedStyle{TranslateY: hide}}, Children: []core.Node{
					Disc(5, colCyan, nil), core.View{Style: core.Style{FlexGrow: 1}}, lpNavChip(), lpNavChip(), lpNavChip(),
				}},
			}}
		}})
}

// #150 Scrollspy: the nav underline moves to whichever section you are in.
func lp2Scrollspy(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "150", Width: 300, Title: "Scrollspy", Duration: 2600 * ms, Accent: colAmber, Height: 110,
		Note: "One underline for four items: it slides instead of switching, so you can see it travel to the new section.",
		Stage: func(p *core.AnimatedValue) core.Node {
			stop := core.Interpolate(p, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 3}})
			x := core.Interpolate(stop, core.Interpolation{Input: []float64{0, 1, 2, 3}, Output: []float64{0, 42, 84, 126}, Easing: bezier(core.EaseInOut)})
			var tabs []core.Node
			for i := 0; i < 4; i++ {
				tabs = append(tabs, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(core.Points(float32(i * 42))),
					Width: core.Pt(42), AlignItems: core.AlignCenter}, Children: []core.Node{lpText(24, colMuted, nil)}})
			}
			return lpFrame(true,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(24), Left: core.Pt(16), Width: core.Pt(170), Height: core.Pt(6)},
					Children: tabs},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(34), Left: core.Pt(16), Height: core.Pt(2), Width: core.Pt(28),
					Background: colAmber, Animated: &core.AnimatedStyle{TranslateX: x}}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(48), Left: core.Pt(16), Right: core.Pt(16), Gap: 6,
					Animated: &core.AnimatedStyle{TranslateY: lerp(p, 0, -20)}}, Children: []core.Node{
					lpText(140, colLine, nil), lpText(100, colLine, nil), lpText(120, colLine, nil),
				}},
			)
		}})
}

// #151 Back-to-top: the button fades in late and the page snaps up.
func lp2BackToTop(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "151", Width: 300, Title: "Back-to-top", Duration: 2400 * ms, Accent: colGreen, Height: 110,
		Note: "It appears only after the page is long enough to need it, and the jump is quick: nobody wants to ride a page back up.",
		Stage: func(p *core.AnimatedValue) core.Node {
			show := core.Interpolate(p, core.Interpolation{Input: []float64{0.35, 0.5}, Output: []float64{0, 1}})
			jump := core.Interpolate(p, core.Interpolation{Input: []float64{0.7, 0.95}, Output: []float64{0, 1}, Easing: bezier(core.EaseInOut)})
			return lpFrame(true,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(24), Left: core.Pt(16), Right: core.Pt(16), Gap: 6,
					Animated: &core.AnimatedStyle{TranslateY: core.Interpolate(jump, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 60}})}},
					Children: []core.Node{
						lpText(150, colLine, nil), lpText(120, colLine, nil), lpText(170, colLine, nil), lpText(90, colLine, nil),
					}},
				core.View{Style: core.Style{Position: core.Absolute, Bottom: core.Pt(10), Right: core.Pt(12), Width: core.Pt(24), Height: core.Pt(24),
					CornerRadius: 12, Background: colGreen, AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter,
					Animated: &core.AnimatedStyle{Opacity: show, Scale: core.Interpolate(show, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0.6, 1}})}},
					Children: []core.Node{
						core.View{Style: core.Style{Width: core.Pt(2), Height: core.Pt(8), Background: core.RGB(0x07101c),
							Animated: &core.AnimatedStyle{TranslateY: constant(cx, -2)}}},
					}},
			)
		}})
}

// #152 Zoom on scroll: the picture grows as it passes the middle of the screen.
func lp2ZoomScroll(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "152", Width: 300, Title: "Zoom on scroll", Duration: 2400 * ms, Accent: colBlue, Height: 110,
		Note: "The zoom peaks as the picture crosses the centre and is back to normal at the edges: it reads as a lens, not a loop.",
		Stage: func(p *core.AnimatedValue) core.Node {
			peak := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.5, 1}, Output: []float64{1, 1.22, 1}, Easing: bezier(core.EaseInOut)})
			return lpFrame(true,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(26), Left: core.Pt(16), Right: core.Pt(16), Height: core.Pt(70),
					CornerRadius: 7, Overflow: core.Hidden,
					Animated: &core.AnimatedStyle{TranslateY: core.Interpolate(p, core.Interpolation{Input: []float64{0, 1}, Output: []float64{14, -14}})}},
					Children: []core.Node{
						core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
							Animated: &core.AnimatedStyle{Scale: peak}}, Children: []core.Node{
							camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{}),
						}},
					}},
			)
		}})
}

// #153 3D tilt: the card answers the cursor with a rotation on both axes.
func lp2Tilt(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "153", Width: 300, Title: "3D tilt card", Duration: 2400 * ms, Accent: colViolet, Height: 110,
		Note: "Two axes, mirrored: the card turns towards the pointer and the sheen slides the other way, which is what sells the depth.",
		Stage: func(p *core.AnimatedValue) core.Node {
			point := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.35, 0.7, 1}, Output: []float64{0, 1, -1, 0}, Easing: bezier(core.EaseInOut)})
			return lpFrame(true,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(26), Left: core.Pt(66), Width: core.Pt(110), Height: core.Pt(66),
					CornerRadius: 8, Background: colWell, BorderWidth: 1, BorderColor: colLine, Overflow: core.Hidden,
					Animated: &core.AnimatedStyle{RotateY: core.Interpolate(point, core.Interpolation{Input: []float64{-1, 1}, Output: []float64{-14, 14}}),
						RotateX: core.Interpolate(point, core.Interpolation{Input: []float64{-1, 1}, Output: []float64{10, -10}})}},
					Children: []core.Node{
						core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(-30), Left: core.Pt(-30), Width: core.Pt(40), Height: core.Pt(140),
							Background: core.RGB(0x2a3a52), Animated: &core.AnimatedStyle{
								TranslateX: core.Interpolate(point, core.Interpolation{Input: []float64{-1, 1}, Output: []float64{150, 10}}),
								Rotate:     constant(cx, 18), Opacity: constant(cx, 0.5)}}},
						core.View{Style: core.Style{Padding: core.Edges{Top: 12, Left: 12}, Gap: 5}, Children: []core.Node{
							Disc(14, colViolet, nil), lpText(60, colText, nil), lpText(40, colMuted, nil),
						}},
					}},
			)
		}})
}

// #154 Spotlight: a light follows the cursor and lights only what it covers.
func lp2Spotlight(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "154", Width: 300, Title: "Spotlight", Duration: 2400 * ms, Accent: colAmber, Height: 110,
		Note: "A dark layer with a soft hole in it, and the hole moves: the picture underneath is never touched.",
		Stage: func(p *core.AnimatedValue) core.Node {
			path := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.5, 1}, Output: []float64{0, 1, 0}, Easing: bezier(core.EaseInOut)})
			return lpFrame(true,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(26), Left: core.Pt(16), Right: core.Pt(16), Height: core.Pt(70),
					CornerRadius: 7, Overflow: core.Hidden, Background: colWell, Padding: core.Edges{Top: 8, Left: 8}, Gap: 5},
					Children: []core.Node{
						lpText(140, colLine, nil), lpText(110, colLine, nil), lpText(160, colLine, nil), lpText(90, colLine, nil),
					}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					Background: core.RGB(0x05080d), Opacity: core.Alpha(0.72)}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(2), Left: core.Pt(0), Width: core.Pt(120), Height: core.Pt(120),
					Animated: &core.AnimatedStyle{TranslateX: core.Interpolate(path, core.Interpolation{Input: []float64{0, 1}, Output: []float64{20, 150}})}}, Children: []core.Node{
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Width: core.Pt(120), Height: core.Pt(120), CornerRadius: 60, Background: core.RGB(0xffd27a), Opacity: core.Alpha(0.05)}},
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(12), Left: core.Pt(12), Width: core.Pt(96), Height: core.Pt(96), CornerRadius: 48, Background: core.RGB(0xffd98a), Opacity: core.Alpha(0.06)}},
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(24), Left: core.Pt(24), Width: core.Pt(72), Height: core.Pt(72), CornerRadius: 36, Background: core.RGB(0xffe2a4), Opacity: core.Alpha(0.08)}},
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(36), Left: core.Pt(36), Width: core.Pt(48), Height: core.Pt(48), CornerRadius: 24, Background: core.RGB(0xffeec8), Opacity: core.Alpha(0.11)}},
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Width: core.Pt(120), Height: core.Pt(120), CornerRadius: 60, BorderWidth: 1, BorderColor: colAmber, Opacity: core.Alpha(0.30)}},
				}},
			)
		}})
}

// #155 Image swap: one picture replaces another on hover, in place.
func lp2ImageSwap(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "155", Width: 300, Title: "Image swap", Duration: 2000 * ms, Accent: colPink, Height: 110,
		Note: "No fade: the second picture wipes across the first from the left, so you see it be replaced rather than appear.",
		Stage: func(p *core.AnimatedValue) core.Node {
			wipe := core.Interpolate(p, core.Interpolation{Input: []float64{0.15, 0.7}, Output: []float64{0, 1}, Easing: bezier(core.EaseInOut)})
			return lpFrame(true,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(26), Left: core.Pt(16), Width: core.Pt(120), Height: core.Pt(70),
					CornerRadius: 7, Overflow: core.Hidden, Background: core.RGB(0x1b3a5c)}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(26), Left: core.Pt(16), Width: core.Pt(120), Height: core.Pt(70),
					CornerRadius: 7, Overflow: core.Hidden, Animated: &core.AnimatedStyle{ScaleX: wipe}},
					Children: []core.Node{
						core.View{Style: core.Style{Width: core.Pt(120), Height: core.Pt(70)}, Children: []core.Node{
							camScene(core.RGB(0x5a2b3c), core.RGB(0x2a1420), camLayers{}),
						}},
					}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(34), Left: core.Pt(150), Gap: 5}, Children: []core.Node{
					lpText(70, colText, nil), lpText(50, colMuted, nil),
				}},
			)
		}})
}

// #156 Text roll: the label rolls up as the next one arrives.
func lp2TextRoll(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "156", Width: 300, Title: "Text roll", Duration: 2000 * ms, Accent: colCyan, Height: 110,
		Note: "Both words move the same distance in the same time, inside a window: a roll is a mask, not a fade.",
		Stage: func(p *core.AnimatedValue) core.Node {
			roll := core.Interpolate(p, core.Interpolation{Input: []float64{0.2, 0.65}, Output: []float64{0, 1}, Easing: bezier(core.EaseInOut)})
			return lpFrame(true,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(40), Left: core.Pt(0), Right: core.Pt(0),
					AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter}, Children: []core.Node{
					core.View{Style: core.Style{Width: core.Pt(160), Height: core.Pt(20), Overflow: core.Hidden,
						AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter}, Children: []core.Node{
						core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Height: core.Pt(20),
							AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter,
							Animated: &core.AnimatedStyle{TranslateY: core.Interpolate(roll, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, -20}})}},
							Children: []core.Node{Word("Deploy", 15, colText, nil)}},
						core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(20), Left: core.Pt(0), Right: core.Pt(0), Height: core.Pt(20),
							AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter,
							Animated: &core.AnimatedStyle{TranslateY: core.Interpolate(roll, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, -20}})}},
							Children: []core.Node{Word("Deployed", 15, colGreen, nil)}},
					}},
					core.View{Style: core.Style{Margin: core.Edges{Top: 6}}, Children: []core.Node{lpText(70, colMuted, nil)}},
				}},
			)
		}})
}

// #157 Hover image reveal: the picture appears where the pointer already is.
func lp2HoverReveal(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "157", Width: 300, Title: "Hover image reveal", Duration: 2200 * ms, Accent: colAmber, Height: 110,
		Note: "The link is a list item until you hover it, and then it is a picture: the reveal is the only motion in the row.",
		Stage: func(p *core.AnimatedValue) core.Node {
			reveal := core.Interpolate(p, core.Interpolation{Input: []float64{0.2, 0.6, 0.9, 1}, Output: []float64{0, 1, 1, 0}, Easing: bezier(core.EaseInOut)})
			var rows []core.Node
			for i := 0; i < 3; i++ {
				rows = append(rows, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(core.Points(float32(30 + i*18))), Left: core.Pt(16),
					Animated: &core.AnimatedStyle{Opacity: core.Interpolate(reveal, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0.5, 1}})}},
					Children: []core.Node{lpText(90-float32(i)*14, colText, nil)}})
			}
			return lpFrame(true,
				core.View{Children: rows},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(32), Left: core.Pt(130), Width: core.Pt(80), Height: core.Pt(52),
					CornerRadius: 6, Overflow: core.Hidden,
					Animated: &core.AnimatedStyle{Opacity: reveal, Scale: core.Interpolate(reveal, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0.85, 1}}),
						Rotate: core.Interpolate(reveal, core.Interpolation{Input: []float64{0, 1}, Output: []float64{-4, 0}})}},
					Children: []core.Node{
						core.View{Style: core.Style{Width: core.Pt(80), Height: core.Pt(52)}, Children: []core.Node{
							camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{}),
						}},
					}},
			)
		}})
}

// #158 Floating label: the placeholder becomes the label above the field.
func lp2FloatingLabel(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "158", Width: 300, Title: "Floating label", Duration: 2200 * ms, Accent: colCyan, Height: 110,
		Note: "The label never changes size: it travels up and out of the field, and the field keeps focus so the caret stays.",
		Stage: func(p *core.AnimatedValue) core.Node {
			float := core.Interpolate(p, core.Interpolation{Input: []float64{0.25, 0.65}, Output: []float64{0, 1}, Easing: bezier(core.EaseInOut)})
			return lpFrame(true,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(40), Left: core.Pt(30), Width: core.Pt(150), Height: core.Pt(26),
					CornerRadius: 5, Background: colWell, BorderWidth: 1, BorderColor: colCyan}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(50), Left: core.Pt(40),
					Animated: &core.AnimatedStyle{TranslateY: core.Interpolate(float, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, -18}}),
						TranslateX: core.Interpolate(float, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, -4}}),
						Scale:      core.Interpolate(float, core.Interpolation{Input: []float64{0, 1}, Output: []float64{1, 0.72}})}},
					Children: []core.Node{core.Text{Content: "Email address", Style: core.Style{FontSize: 9, Color: colMuted}}}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(50), Left: core.Pt(40), Width: core.Pt(60), Height: core.Pt(8),
					Animated: &core.AnimatedStyle{Opacity: float, ScaleX: float}}, Children: []core.Node{
					core.View{Style: core.Style{Width: core.Pt(2), Height: core.Pt(8), Background: colCyan}},
				}},
			)
		}})
}

// #159 Inline validation: the border turns, the message slides out under it.
func lp2InlineValidation(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "159", Width: 300, Title: "Inline validation", Duration: 2400 * ms, Accent: colPink, Height: 110,
		Note: "Bad, then good: the field stays put and only the border and the note under it change, so the eye is not sent hunting.",
		Stage: func(p *core.AnimatedValue) core.Node {
			bad := core.Interpolate(p, core.Interpolation{Input: []float64{0.2, 0.3, 0.6, 0.7}, Output: []float64{0, 1, 1, 0}})
			good := core.Interpolate(p, core.Interpolation{Input: []float64{0.6, 0.75}, Output: []float64{0, 1}})
			shake := core.Interpolate(p, core.Interpolation{Input: []float64{0.3, 0.36, 0.42, 0.48, 0.54}, Output: []float64{0, -4, 4, -3, 0}})
			return lpFrame(true,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(34), Left: core.Pt(30), Width: core.Pt(150), Height: core.Pt(24),
					CornerRadius: 5, Background: colWell, BorderWidth: 1, BorderColor: colLine,
					Animated: &core.AnimatedStyle{TranslateX: shake,
						BorderColor: core.Interpolate(good, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1},
							Colors: []core.Color{core.RGB(0xd94b6a), core.RGB(0x39d98a)}})}}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(62), Left: core.Pt(30), Gap: 3}, Children: []core.Node{
					core.View{Style: core.Style{Animated: &core.AnimatedStyle{Opacity: bad, TranslateY: core.Interpolate(bad, core.Interpolation{Input: []float64{0, 1}, Output: []float64{-4, 0}})}},
						Children: []core.Node{core.Text{Content: "That domain does not exist", Style: core.Style{FontSize: 8, Color: core.RGB(0xd94b6a)}}}},
					core.View{Style: core.Style{Animated: &core.AnimatedStyle{Opacity: good, TranslateY: core.Interpolate(good, core.Interpolation{Input: []float64{0, 1}, Output: []float64{-4, 0}})}},
						Children: []core.Node{core.Text{Content: "Looks good", Style: core.Style{FontSize: 8, Color: colGreen}}}},
				}},
			)
		}})
}

// #160 Loading button: the label goes, the spinner takes its place, the tick ends it.
func lp2LoadingButton(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "160", Width: 300, Title: "Loading button", Duration: 2400 * ms, Accent: colCyan, Height: 110,
		Note: "The button keeps its width through all three states: a button that resizes under the pointer is a button you miss.",
		Stage: func(p *core.AnimatedValue) core.Node {
			busy := core.Interpolate(p, core.Interpolation{Input: []float64{0.2, 0.3, 0.62, 0.72}, Output: []float64{0, 1, 1, 0}})
			done := core.Interpolate(p, core.Interpolation{Input: []float64{0.72, 0.85}, Output: []float64{0, 1}})
			spin := core.Interpolate(p, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 720}})
			// The label leaves with the spinner and stays gone: "Done" ends it,
			// so both texts are never on screen at once.
			label := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.2, 0.3, 1}, Output: []float64{1, 1, 0, 0}})
			return lpFrame(true,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(44), Left: core.Pt(0), Right: core.Pt(0),
					AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter}, Children: []core.Node{
					core.View{Style: core.Style{Width: core.Pt(96), Height: core.Pt(22), CornerRadius: 11, Background: colCyan,
						AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter}, Children: []core.Node{
						core.View{Style: core.Style{Position: core.Absolute, Animated: &core.AnimatedStyle{Opacity: label}},
							Children: []core.Node{Word("Checkout", 10, core.RGB(0x07101c), nil)}},
						core.View{Style: core.Style{Position: core.Absolute, Animated: &core.AnimatedStyle{Opacity: busy, Rotate: spin}},
							Children: []core.Node{Ring(11, 2, core.RGB(0x07101c), nil)}},
						core.View{Style: core.Style{Position: core.Absolute, Animated: &core.AnimatedStyle{Opacity: done}},
							Children: []core.Node{Word("Done", 10, core.RGB(0x07101c), nil)}},
					}},
				}},
			)
		}})
}

// #161 View transition: a card grows into the page it links to.
func lp2ViewTransition(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "161", Width: 300, Title: "View transition", Duration: 2600 * ms, Accent: colViolet, Height: 110,
		Note: "The card and the page share one element: it moves and resizes between them, and both frames are faded around it.",
		Stage: func(p *core.AnimatedValue) core.Node {
			move := core.Interpolate(p, core.Interpolation{Input: []float64{0.15, 0.75}, Output: []float64{0, 1}, Easing: bezier(core.EaseInOut)})
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
				Overflow: core.Hidden, Background: colCard}, Children: []core.Node{
				lpFrame(true,
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(30), Left: core.Pt(16), Right: core.Pt(16), Gap: 6},
						Children: []core.Node{lpText(130, colLine, nil), lpText(90, colLine, nil)}},
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(56), Left: core.Pt(16), Width: core.Pt(60), Height: core.Pt(36),
						CornerRadius: 5, Background: colViolet, Animated: &core.AnimatedStyle{Opacity: core.Interpolate(move, core.Interpolation{Input: []float64{0, 0.4}, Output: []float64{1, 0}})}}},
				),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(26), Left: core.Pt(16), Right: core.Pt(16), Bottom: core.Pt(10),
					CornerRadius: 5, Background: colWell,
					Animated: &core.AnimatedStyle{Scale: core.Interpolate(move, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0.3, 1}}),
						Opacity: core.Interpolate(move, core.Interpolation{Input: []float64{0, 0.6, 1}, Output: []float64{0, 0.9, 1}})}},
					Children: []core.Node{
						core.View{Style: core.Style{Padding: core.Edges{Top: 10, Left: 10}, Gap: 5}, Children: []core.Node{
							lpText(120, colText, nil), lpText(150, colMuted, nil), lpText(100, colLine, nil),
						}},
					}},
			}}
		}})
}

// #162 Lightbox: the thumbnail lifts out of the grid into a full frame.
func lp2Lightbox(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "162", Width: 300, Title: "Lightbox", Duration: 2400 * ms, Accent: colBlue, Height: 110,
		Note: "The bigger frame starts where the thumbnail is: the shared origin is what makes it a lift rather than a new window.",
		Stage: func(p *core.AnimatedValue) core.Node {
			open := core.Interpolate(p, core.Interpolation{Input: []float64{0.15, 0.65}, Output: []float64{0, 1}, Easing: bezier(core.EaseOut)})
			return core.View{Style: stageFill, Children: []core.Node{
				lpFrame(true,
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(30), Left: core.Pt(16), Width: core.Pt(50), Height: core.Pt(34),
						CornerRadius: 4, Background: colBlue}},
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(30), Left: core.Pt(76), Width: core.Pt(50), Height: core.Pt(34),
						CornerRadius: 4, Background: colWell}},
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(70), Left: core.Pt(16), Width: core.Pt(50), Height: core.Pt(30),
						CornerRadius: 4, Background: colWell}},
				),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					Background: core.RGB(0x05080d), Animated: &core.AnimatedStyle{Opacity: core.Interpolate(open, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 0.86}})}}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(22), Left: core.Pt(30), Right: core.Pt(30), Height: core.Pt(66),
					CornerRadius: 6, Overflow: core.Hidden,
					Animated: &core.AnimatedStyle{Opacity: open, Scale: core.Interpolate(open, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0.45, 1}})}},
					Children: []core.Node{
						core.View{Style: core.Style{Width: core.Pt(200), Height: core.Pt(66)}, Children: []core.Node{
							camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{}),
						}},
					}},
			}}
		}})
}

// #163 Hotspots: the dots pulse and the card opens on the one you touch.
func lp2Hotspots(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "163", Width: 300, Title: "Hotspots", Duration: 2600 * ms, Accent: colPink, Height: 110,
		Note: "Three dots, one open card: the dots keep pulsing so they still read as targets after the card has opened.",
		Stage: func(p *core.AnimatedValue) core.Node {
			active := core.Interpolate(p, core.Interpolation{Input: []float64{0.2, 0.5}, Output: []float64{0, 1}})
			pulse := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.5, 1}, Output: []float64{1, 1.6, 1}})
			spot := func(left float32, top float32) core.Node {
				return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(core.Points(top)), Left: core.Pt(core.Points(left))},
					Children: []core.Node{
						core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(-4), Left: core.Pt(-4), Animated: &core.AnimatedStyle{Scale: pulse, Opacity: constant(cx, 0.35)}},
							Children: []core.Node{Disc(16, colPink, nil)}},
						Disc(8, colPink, nil),
					}}
			}
			return lpFrame(true,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(26), Left: core.Pt(14), Right: core.Pt(14), Height: core.Pt(70),
					CornerRadius: 6, Background: colWell, Overflow: core.Hidden}, Children: []core.Node{
					core.View{Style: core.Style{Width: core.Pt(200), Height: core.Pt(70)}, Children: []core.Node{
						camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{}),
					}},
				}},
				spot(40, 60),
				spot(120, 44),
				spot(170, 66),
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(30), Left: core.Pt(150), Width: core.Pt(70), Height: core.Pt(30),
					CornerRadius: 5, Background: colBar, Padding: core.Edges{Top: 5, Left: 6}, Gap: 3,
					Animated: &core.AnimatedStyle{Opacity: active, TranslateY: core.Interpolate(active, core.Interpolation{Input: []float64{0, 1}, Output: []float64{6, 0}})}},
					Children: []core.Node{lpText(34, colText, nil), lpText(50, colLine, nil)}},
			)
		}})
}

// #164 Countdown: the numbers roll over, the last second is the loud one.
func lp2Countdown(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "164", Width: 300, Title: "Countdown", Duration: 2600 * ms, Accent: colAmber, Height: 110,
		Note: "Each digit rolls the same distance; the last one changes colour instead of moving faster, because time does not speed up.",
		Stage: func(p *core.AnimatedValue) core.Node {
			roll := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.28, 0.36, 0.62, 0.7, 1},
				Output: []float64{0, 0, -40, -40, -80, -80}, Easing: bezier(core.EaseInOut)})
			var cells []core.Node
			for i := 0; i < 3; i++ {
				cells = append(cells, core.View{Style: core.Style{Width: core.Pt(30), Height: core.Pt(40), CornerRadius: 5, Background: colWell,
					Overflow: core.Hidden, AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter}, Children: []core.Node{
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0),
						AlignItems: core.AlignCenter, Animated: &core.AnimatedStyle{TranslateY: roll}},
						Children: []core.Node{
							core.View{Style: core.Style{Height: core.Pt(40), AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter},
								Children: []core.Node{Word("3", 15, colAmber, nil)}},
							core.View{Style: core.Style{Height: core.Pt(40), AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter},
								Children: []core.Node{Word("2", 15, colAmber, nil)}},
							core.View{Style: core.Style{Height: core.Pt(40), AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter},
								Children: []core.Node{Word("1", 15, colPink, nil)}},
						}},
				}})
			}
			return lpFrame(true,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(28), Left: core.Pt(0), Right: core.Pt(0), Direction: core.Row,
					JustifyContent: core.JustifyCenter, Gap: 5}, Children: cells},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(76), Left: core.Pt(0), Right: core.Pt(0),
					AlignItems: core.AlignCenter}, Children: []core.Node{lpText(60, colMuted, nil)}},
			)
		}})
}

// #165 Shape divider: the section edge is a shape that moves with the page.
func lp2ShapeDivider(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "165", Width: 300, Title: "Shape divider", Duration: 2400 * ms, Accent: colGreen, Height: 110,
		Note: "The divider is made of a few angled panels, not a path: at this size five steps read as a curve, and the drift gives it life.",
		Stage: func(p *core.AnimatedValue) core.Node {
			drift := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.5, 1}, Output: []float64{-6, 6, -6}, Easing: bezier(core.EaseInOut)})
			var steps []core.Node
			for i := 0; i < 7; i++ {
				h := 18 - i*2 + (i%2)*2
				steps = append(steps, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(core.Points(float32(h))),
					Left: core.Pt(core.Points(float32(i) * 34)), Width: core.Pt(38), Height: core.Pt(60), Background: colGreen,
					Animated: &core.AnimatedStyle{TranslateY: core.Interpolate(drift, core.Interpolation{Input: []float64{-6, 6}, Output: []float64{-3, 3}}),
						Rotate: constant(cx, 8)}}})
			}
			return lpFrame(true,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(24), Left: core.Pt(16), Gap: 5}, Children: []core.Node{
					lpText(90, colText, nil), lpText(60, colMuted, nil),
				}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(50), Left: core.Pt(-6), Right: core.Pt(-6), Bottom: core.Pt(0),
					Overflow: core.Hidden}, Children: steps},
			)
		}})
}

// #166 Like button: the heart pops, the count follows.
func lp2Like(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "166", Width: 300, Title: "Like button", Duration: 2000 * ms, Accent: colPink, Height: 110,
		Note: "The heart overshoots and settles, the ring leaves once, and the number changes while the heart is still moving.",
		Stage: func(p *core.AnimatedValue) core.Node {
			pop := core.Interpolate(p, core.Interpolation{Input: []float64{0.15, 0.35, 0.5, 0.65}, Output: []float64{1, 1.5, 0.92, 1}})
			ring := core.Interpolate(p, core.Interpolation{Input: []float64{0.35, 0.9}, Output: []float64{0, 1}})
			return lpFrame(true,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(38), Left: core.Pt(0), Right: core.Pt(0),
					AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter, Direction: core.Row, Gap: 10}, Children: []core.Node{
					core.View{Style: core.Style{Position: core.Relative, Width: core.Pt(44), Height: core.Pt(44), AlignItems: core.AlignCenter,
						JustifyContent: core.JustifyCenter}, Children: []core.Node{
						core.View{Style: core.Style{Position: core.Absolute, Width: core.Pt(34), Height: core.Pt(34), CornerRadius: 17,
							BorderWidth: 2, BorderColor: colPink,
							Animated: &core.AnimatedStyle{Opacity: core.Interpolate(ring, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0.9, 0}}),
								Scale: core.Interpolate(ring, core.Interpolation{Input: []float64{0, 1}, Output: []float64{1, 2.1}})}}},
						// A heart from two lobes and a turned square.
						core.View{Style: core.Style{Direction: core.Row, Gap: 0, Animated: &core.AnimatedStyle{Scale: pop}}, Children: []core.Node{
							Disc(9, colPink, nil),
							Disc(9, colPink, nil),
						}},
					}},
					core.View{Style: core.Style{Animated: &core.AnimatedStyle{Opacity: core.Interpolate(p, core.Interpolation{Input: []float64{0.35, 0.45}, Output: []float64{0, 1}})}},
						Children: []core.Node{Word("1,284", 13, colText, nil)}},
				}},
			)
		}})
}

// #167 Scroll cue: the arrow that says there is more below.
func lp2ScrollCue(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "167", Width: 300, Title: "Scroll cue", Duration: 2000 * ms, Accent: colMuted, Height: 110, Loop: true,
		Note: "One element, one direction: the cue falls, fades and repeats, and it stops the moment the page really does scroll.",
		Stage: func(p *core.AnimatedValue) core.Node {
			fall := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.7, 1}, Output: []float64{0, 16, 16}})
			byebye := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.55, 1}, Output: []float64{1, 1, 0}})
			return lpFrame(true,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(26), Left: core.Pt(0), Right: core.Pt(0),
					AlignItems: core.AlignCenter, Gap: 6}, Children: []core.Node{
					lpText(90, colText, nil),
					lpText(60, colMuted, nil),
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(64), Left: core.Pt(0), Right: core.Pt(0),
						AlignItems: core.AlignCenter, Animated: &core.AnimatedStyle{Opacity: byebye}}, Children: []core.Node{
						core.View{Style: core.Style{Width: core.Pt(20), Height: core.Pt(30), CornerRadius: 10, BorderWidth: 2, BorderColor: colText}, Children: []core.Node{
							core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(5), Left: core.Pt(0), Right: core.Pt(0),
								AlignItems: core.AlignCenter, Animated: &core.AnimatedStyle{TranslateY: fall}}, Children: []core.Node{Disc(4, colText, nil)}},
						}},
					}},
				}},
			)
		}})
}

// #168 Background video: the clip runs and a scrim keeps the words readable.
func lp2BgVideo(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "168", Width: 300, Title: "Background video", Duration: 2600 * ms, Accent: colCyan, Height: 110,
		Note: "The scrim is not decoration: the same clip has to carry white text in every frame, so the darkening is the design.",
		Stage: func(p *core.AnimatedValue) core.Node {
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
				Overflow: core.Hidden, Background: colBar}, Children: []core.Node{
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					Animated: &core.AnimatedStyle{Scale: core.Interpolate(p, core.Interpolation{Input: []float64{0, 1}, Output: []float64{1, 1.14}})}},
					Children: []core.Node{camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{
						Far:    &core.AnimatedStyle{TranslateX: lerp(p, 20, -20)},
						Near:   &core.AnimatedStyle{TranslateX: lerp(p, 50, -50)},
						Figure: &core.AnimatedStyle{TranslateY: lerp(p, 0, -6)},
					})}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					Background: core.RGB(0x05080d), Opacity: core.Alpha(0.55)}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(34), Left: core.Pt(16), Gap: 6,
					Animated: &core.AnimatedStyle{Opacity: core.Interpolate(p, core.Interpolation{Input: []float64{0.1, 0.5}, Output: []float64{0, 1}})}},
					Children: []core.Node{
						Word("Ship it", 16, colText, nil),
						lpText(110, colMuted, nil),
						lpButton(52, colCyan, nil),
					}},
			}}
		}})
}
