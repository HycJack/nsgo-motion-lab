package main

// page_lp.go is the landing-page pack (motion-prompt #073–#096): the sections
// a marketing page is built from, each with the entrance that makes it read.
//
// The upstream titles for this range could not be fetched while this was
// written, so the twenty-four names here are our own labels for the same
// section patterns. lp2 (145–168) carries the upstream names verbatim.
//
// Catalog structure and timings: github.com/showheyohtaki/motion-prompt (MIT).

import (
	"github.com/oneo-me/nsgo/core"
)

// lpFrame is a mock browser page: a nav strip and a content area, both empty.
func lpFrame(nav bool, children ...core.Node) core.Node {
	body := core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
		Overflow: core.Hidden, Background: colCard}, Children: children}
	if !nav {
		return body
	}
	return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
		Overflow: core.Hidden, Background: colCard}, Children: []core.Node{
		body,
		core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Height: core.Pt(18),
			Background: colBar, Direction: core.Row, AlignItems: core.AlignCenter, Padding: core.Edges{Left: 10, Right: 10}, Gap: 6},
			Children: []core.Node{
				Disc(5, colCyan, nil),
				core.View{Style: core.Style{Width: core.Pt(22), Height: core.Pt(3), CornerRadius: 2, Background: colLine}},
				core.View{Style: core.Style{FlexGrow: 1}},
				lpNavChip(), lpNavChip(), lpNavChip(),
			}},
	}}
}

// lpNavChip is one nav item.
func lpNavChip() core.Node {
	return core.View{Style: core.Style{Width: core.Pt(14), Height: core.Pt(3), CornerRadius: 2, Background: colLine}}
}

// lpText is a line of body copy of a given width.
func lpText(w float32, col core.Color, anim *core.AnimatedStyle) core.Node {
	return core.View{Style: core.Style{Width: core.Pt(core.Points(w)), Height: core.Pt(4), CornerRadius: 2, Background: col, Animated: anim}}
}

// lpButton is a filled call to action.
func lpButton(w float32, col core.Color, anim *core.AnimatedStyle) core.Node {
	return core.View{Style: core.Style{Width: core.Pt(core.Points(w)), Height: core.Pt(13), CornerRadius: 7, Background: col,
		AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter, Animated: anim},
		Children: []core.Node{core.View{Style: core.Style{Width: core.Pt(core.Points(w * 0.45)), Height: core.Pt(2.5), CornerRadius: 2, Background: core.RGB(0x07101c)}}}}
}

// pageLP lays the pack out.
func pageLP(cx *core.Scope) core.Node {
	return core.View{Children: []core.Node{
		core.View{Style: core.Style{Padding: core.Edges{Top: 14, Right: 18, Left: 18}}, Children: []core.Node{
			note("Landing page pack 073–096: the section entrance is the effect. Patterns follow motion-prompt (github.com/showheyohtaki/motion-prompt); the labels are ours, the code is NSGo's."),
		}},
		grid(
			lpHero(cx),
			lpStickyNav(cx),
			lpFeatureGrid(cx),
			lpSplitHero(cx),
			lpLogoCloud(cx),
			lpTestimonial(cx),
			lpPricingToggle(cx),
			lpFunnel(cx),
			lpProductTour(cx),
			lpScreenshot(cx),
			lpMetrics(cx),
			lpBento(cx),
			lpCTA(cx),
			lpFAQ(cx),
			lpBlogCard(cx),
			lpTeam(cx),
			lpFooter(cx),
			lpCookie(cx),
			lpCookieBar(cx),
			lpSearch(cx),
			lpSignup(cx),
			lpOnboarding(cx),
			lpEmpty(cx),
			lpSkeleton(cx),
		),
	}}
}

// #073 Hero entrance: the words arrive before the picture does.
func lpHero(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "073", Width: 300, Title: "Hero entrance", Duration: 2200 * ms, Accent: colCyan, Height: 110,
		Note: "Headline, subhead, button, picture — in that order, sixty milliseconds apart, each one rising out of nothing.",
		Stage: func(p *core.AnimatedValue) core.Node {
			at := func(delay float64) core.AnimatedNode {
				return core.Interpolate(p, core.Interpolation{Input: []float64{delay, delay + 0.35}, Output: []float64{0, 1}, Easing: bezier(core.EaseOut)})
			}
			rise := func(q core.AnimatedNode) core.AnimatedNode {
				return core.Interpolate(q, core.Interpolation{Input: []float64{0, 1}, Output: []float64{12, 0}})
			}
			return lpFrame(true,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(30), Left: core.Pt(12), AlignItems: core.AlignCenter, Gap: 5,
					Animated: &core.AnimatedStyle{Opacity: at(0.05), TranslateY: rise(at(0.05))}}, Children: []core.Node{
					lpText(84, colText, nil),
					lpText(60, colText, nil),
					lpText(70, colMuted, nil),
					lpButton(46, colCyan, nil),
				}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(34), Left: core.Pt(122), Width: core.Pt(100), Height: core.Pt(62),
					CornerRadius: 6, Background: colWell, Animated: &core.AnimatedStyle{Opacity: at(0.45), Scale: core.Interpolate(at(0.45), core.Interpolation{Input: []float64{0, 1}, Output: []float64{0.9, 1}})}}},
			)
		}})
}

// #074 Sticky nav: the bar detaches, shrinks and stays.
func lpStickyNav(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "074", Width: 300, Title: "Sticky nav", Duration: 2200 * ms, Accent: colGreen, Height: 110,
		Note: "The page scrolls under a bar that has already stopped: height and elevation change once, at the moment it pins.",
		Stage: func(p *core.AnimatedValue) core.Node {
			pin := core.Interpolate(p, core.Interpolation{Input: []float64{0.3, 0.55}, Output: []float64{0, 1}})
			return lpFrame(false,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0)},
					Children: []core.Node{
						core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
							AlignItems: core.AlignCenter, Gap: 8, Animated: &core.AnimatedStyle{TranslateY: lerp(p, 0, -40)}}, Children: []core.Node{
							lpText(90, colLine, nil), lpText(120, colLine, nil), lpText(70, colLine, nil),
						}},
						core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Height: core.Pt(18),
							Background: colBar, Direction: core.Row, AlignItems: core.AlignCenter, Padding: core.Edges{Left: 10, Right: 10}, Gap: 6,
							Animated: &core.AnimatedStyle{ScaleY: core.Interpolate(pin, core.Interpolation{Input: []float64{0, 1}, Output: []float64{1, 0.72}})}},
							Children: []core.Node{Disc(5, colGreen, nil), core.View{Style: core.Style{FlexGrow: 1}}, lpNavChip(), lpNavChip()}},
						core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(18), Left: core.Pt(0), Right: core.Pt(0), Height: core.Pt(1),
							Background: colLine, Animated: &core.AnimatedStyle{Opacity: pin}}},
					}},
			)
		}})
}

// #075 Feature grid: rows of cards, staggered by row and by column.
func lpFeatureGrid(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "075", Width: 300, Title: "Feature grid stagger", Duration: 2400 * ms, Accent: colCyan, Height: 110,
		Note: "Stagger by column and by row at once: the eye reads left to right and then down, which is how the grid fills.",
		Stage: func(p *core.AnimatedValue) core.Node {
			var cards []core.Node
			for i := 0; i < 6; i++ {
				row, col := i/3, i%3
				delay := 0.08*float64(col) + 0.2*float64(row)
				q := core.Interpolate(p, core.Interpolation{Input: []float64{0.05 + delay, 0.35 + delay}, Output: []float64{0, 1}, Easing: bezier(core.EaseOut)})
				cards = append(cards, core.View{Style: core.Style{Position: core.Absolute,
					Top: core.Pt(core.Points(float32(26 + row*38))), Left: core.Pt(core.Points(float32(10 + col*78))),
					Width: core.Pt(72), Height: core.Pt(32), CornerRadius: 5, Background: colWell, Padding: core.Edges{Top: 6, Left: 7},
					Animated: &core.AnimatedStyle{Opacity: q, TranslateY: core.Interpolate(q, core.Interpolation{Input: []float64{0, 1}, Output: []float64{10, 0}})}},
					Children: []core.Node{
						Disc(9, []core.Color{colCyan, colViolet, colGreen, colAmber, colPink, colBlue}[i], nil),
						core.View{Style: core.Style{Margin: core.Edges{Top: 4}}, Children: []core.Node{lpText(40, colLine, nil)}},
					}})
			}
			return lpFrame(true, core.View{Children: cards})
		}})
}

// #076 Split hero: text one side, picture the other, meeting in the middle.
func lpSplitHero(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "076", Width: 300, Title: "Split hero", Duration: 2200 * ms, Accent: colViolet, Height: 110,
		Note: "The two halves come in from opposite edges and land together; anything that arrives from the same side reads as one block.",
		Stage: func(p *core.AnimatedValue) core.Node {
			in := core.Interpolate(p, core.Interpolation{Input: []float64{0.1, 0.7}, Output: []float64{0, 1}, Easing: bezier(core.EaseOut)})
			return lpFrame(true,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(34), Left: core.Pt(12), Width: core.Pt(110), AlignItems: core.AlignCenter, Gap: 5,
					Animated: &core.AnimatedStyle{TranslateX: core.Interpolate(in, core.Interpolation{Input: []float64{0, 1}, Output: []float64{-50, 0}}), Opacity: in}},
					Children: []core.Node{lpText(96, colText, nil), lpText(74, colMuted, nil), lpButton(44, colViolet, nil)}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(26), Left: core.Pt(130), Width: core.Pt(108), Height: core.Pt(70),
					CornerRadius: 8, Background: colWell, Animated: &core.AnimatedStyle{TranslateX: core.Interpolate(in, core.Interpolation{Input: []float64{0, 1}, Output: []float64{50, 0}}), Opacity: in}}},
			)
		}})
}

// #077 Logo cloud: proof by numbers, drifting rather than popping.
func lpLogoCloud(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "077", Width: 300, Title: "Logo cloud", Duration: 2200 * ms, Accent: colMuted, Height: 110,
		Note: "Eight marks, no emphasis: the arrangement is a texture, so they fade in as one block and then drift slowly.",
		Stage: func(p *core.AnimatedValue) core.Node {
			fade := core.Interpolate(p, core.Interpolation{Input: []float64{0.05, 0.5}, Output: []float64{0, 1}})
			drift := core.Interpolate(p, core.Interpolation{Input: []float64{0.3, 1}, Output: []float64{0, -8}})
			var marks []core.Node
			for i := 0; i < 8; i++ {
				marks = append(marks, core.View{Style: core.Style{Position: core.Absolute,
					Top: core.Pt(core.Points(float32(34 + (i/4)*30))), Left: core.Pt(core.Points(float32(12 + (i%4)*58))),
					Width: core.Pt(40), Height: core.Pt(9), CornerRadius: 4, Background: colLine,
					Animated: &core.AnimatedStyle{Opacity: fade, TranslateY: drift}}})
			}
			return lpFrame(false,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(18), Left: core.Pt(0), Right: core.Pt(0),
					AlignItems: core.AlignCenter, Animated: &core.AnimatedStyle{Opacity: fade}}, Children: []core.Node{lpText(60, colMuted, nil)}},
				core.View{Children: marks},
			)
		}})
}

// #078 Testimonial slider: one quote at a time, and the dots agree.
func lpTestimonial(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "078", Width: 300, Title: "Testimonial slider", Duration: 2600 * ms, Accent: colAmber, Height: 110,
		Note: "Each quote is on for a third of the pass and the two that are off slide out sideways, which is what makes it a slider.",
		Stage: func(p *core.AnimatedValue) core.Node {
			var slides []core.Node
			for i := 0; i < 3; i++ {
				from := float64(i) / 3
				to := float64(i+1) / 3
				on := core.Interpolate(p, core.Interpolation{Input: []float64{from, from + 0.12, to - 0.12, to},
					Output: []float64{0, 1, 1, 0}})
				slides = append(slides, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(28), Left: core.Pt(14), Width: core.Pt(200),
					Animated: &core.AnimatedStyle{Opacity: on, TranslateX: core.Interpolate(on, core.Interpolation{Input: []float64{0, 1}, Output: []float64{30, 0}})}},
					Children: []core.Node{
						lpText(180, colText, nil),
						core.View{Style: core.Style{Margin: core.Edges{Top: 5}}, Children: []core.Node{lpText(60, colMuted, nil)}},
					}})
			}
			return lpFrame(true, core.View{Children: append(slides,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(76), Left: core.Pt(14), Direction: core.Row, Gap: 4},
					Children: []core.Node{
						Disc(4, colAmber, nil), Disc(4, colLine, nil), Disc(4, colLine, nil),
					}})})
		}})
}

// #079 Pricing toggle: monthly to yearly, and the numbers move.
func lpPricingToggle(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "079", Width: 300, Title: "Pricing toggle", Duration: 2000 * ms, Accent: colGreen, Height: 110,
		Note: "Only the knob and the two prices change: everything else holds still, so the switch is unmistakably the cause.",
		Stage: func(p *core.AnimatedValue) core.Node {
			on := core.Interpolate(p, core.Interpolation{Input: []float64{0.25, 0.7}, Output: []float64{0, 1}, Easing: bezier(core.EaseInOut)})
			return lpFrame(true,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(24), Left: core.Pt(0), Right: core.Pt(0), AlignItems: core.AlignCenter},
					Children: []core.Node{
						core.View{Style: core.Style{Width: core.Pt(52), Height: core.Pt(12), CornerRadius: 6, Background: colLine,
							Animated: &core.AnimatedStyle{Background: core.Interpolate(on, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1}, Colors: []core.Color{colLine, colGreen}})}}},
						core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(2), Left: core.Pt(2),
							Animated: &core.AnimatedStyle{TranslateX: core.Interpolate(on, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 40}})}},
							Children: []core.Node{Disc(8, core.RGB(0x0a1018), nil)}},
					}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(48), Left: core.Pt(0), Right: core.Pt(0), Direction: core.Row,
					JustifyContent: core.JustifyCenter, Gap: 12}, Children: []core.Node{
					core.View{Style: core.Style{Width: core.Pt(80), Height: core.Pt(44), CornerRadius: 6, Background: colWell,
						Padding: core.Edges{Top: 8, Left: 8}, Gap: 4}, Children: []core.Node{
						lpText(30, colMuted, nil),
						Word("$9", 16, colText, &core.AnimatedStyle{Opacity: core.Interpolate(on, core.Interpolation{Input: []float64{0, 0.4, 0.6, 1}, Output: []float64{1, 0, 0, 1}})}),
					}},
					core.View{Style: core.Style{Width: core.Pt(80), Height: core.Pt(44), CornerRadius: 6, Background: colWell,
						Padding: core.Edges{Top: 8, Left: 8}, Gap: 4}, Children: []core.Node{
						lpText(30, colMuted, nil),
						Word("$7", 16, colGreen, &core.AnimatedStyle{Opacity: on}),
					}},
				}},
			)
		}})
}

// #080 Funnel: four steps, each one narrower and one beat later.
func lpFunnel(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "080", Width: 300, Title: "Funnel steps", Duration: 2200 * ms, Accent: colAmber, Height: 110,
		Note: "Each step is less wide than the one above it; the widths are the data, the stagger is only polish.",
		Stage: func(p *core.AnimatedValue) core.Node {
			var rows []core.Node
			for i := 0; i < 4; i++ {
				delay := 0.1 * float64(i)
				q := core.Interpolate(p, core.Interpolation{Input: []float64{delay, delay + 0.3}, Output: []float64{0, 1}, Easing: bezier(core.EaseOut)})
				rows = append(rows, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(core.Points(float32(22 + i*20))),
					Left: core.Pt(core.Points(float32(40 + i*14))), Width: core.Pt(core.Points(float32(170 - i*28))), Height: core.Pt(14),
					CornerRadius: 4, Background: []core.Color{colBlue, colCyan, colGreen, colAmber}[i],
					Animated: &core.AnimatedStyle{ScaleX: q}}})
			}
			return lpFrame(true, core.View{Children: rows})
		}})
}

// #081 Product tour: a spotlight walks through the interface.
func lpProductTour(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "081", Width: 300, Title: "Product tour", Duration: 2800 * ms, Accent: colCyan, Height: 110,
		Note: "The spotlight moves, the app does not: the tour is a window sliding over a screenshot, with a caption that keeps up.",
		Stage: func(p *core.AnimatedValue) core.Node {
			spot := core.Interpolate(p, core.Interpolation{Input: []float64{0.15, 0.85}, Output: []float64{0, 1}, Easing: bezier(core.EaseInOut)})
			return lpFrame(true,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(26), Left: core.Pt(14), Width: core.Pt(200), Height: core.Pt(70),
					CornerRadius: 6, Background: colWell, Padding: core.Edges{Top: 8, Left: 8}, Gap: 5}, Children: []core.Node{
					lpText(160, colLine, nil), lpText(140, colLine, nil), lpText(170, colLine, nil), lpText(110, colLine, nil),
				}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(28), Left: core.Pt(16), Width: core.Pt(70), Height: core.Pt(18),
					CornerRadius: 4, BorderWidth: 2, BorderColor: colCyan,
					Animated: &core.AnimatedStyle{TranslateY: core.Interpolate(spot, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 44}})}}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(94), Left: core.Pt(16), Width: core.Pt(90), Height: core.Pt(10),
					CornerRadius: 3, Background: colCyan, Animated: &core.AnimatedStyle{Opacity: core.Interpolate(spot, core.Interpolation{Input: []float64{0, 0.3}, Output: []float64{0, 1}})}}},
			)
		}})
}

// #082 Screenshot parallax: the window tilts and the screen inside moves less.
func lpScreenshot(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "082", Width: 300, Title: "Screenshot parallax", Duration: 2400 * ms, Accent: colBlue, Height: 110,
		Note: "Two speeds inside one object: the frame tilts slowly and the content inside it slides, which reads as depth.",
		Stage: func(p *core.AnimatedValue) core.Node {
			tilt := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.5, 1}, Output: []float64{-8, 0, 8}, Easing: bezier(core.EaseInOut)})
			return lpFrame(true,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(28), Left: core.Pt(56), Width: core.Pt(130), Height: core.Pt(70),
					CornerRadius: 7, Background: colWell, Overflow: core.Hidden,
					Animated: &core.AnimatedStyle{Rotate: tilt, TranslateY: core.Interpolate(tilt, core.Interpolation{Input: []float64{-8, 8}, Output: []float64{4, -4}})}},
					Children: []core.Node{
						core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(-10), Left: core.Pt(-16), Width: core.Pt(170), Height: core.Pt(100)},
							Children: []core.Node{camScene(core.RGB(0x1b3a5c), core.RGB(0x152436), camLayers{
								Far:  &core.AnimatedStyle{TranslateX: lerp(p, 30, -30)},
								Near: &core.AnimatedStyle{TranslateX: lerp(p, 60, -60)},
							})}},
					}},
			)
		}})
}

// #083 Metrics band: three numbers, counted at their own speeds.
func lpMetrics(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "083", Width: 300, Title: "Metrics band", Duration: 2400 * ms, Accent: colGreen, Height: 110,
		Note: "Three counters, three speeds, one arrival: they finish together because the biggest number has the furthest to go.",
		Stage: func(p *core.AnimatedValue) core.Node {
			values := []string{"0", "12", "48", "96", "127"}
			count := core.Interpolate(p, core.Interpolation{Input: []float64{0.15, 0.85}, Output: []float64{0, 1}, Easing: bezier(core.EaseOut)})
			_ = count
			make1 := func(left float32, lab float32) core.Node {
				var words []core.Node
				for i, v := range values {
					from := 0.1 + 0.15*float64(i)
					words = append(words, core.View{Style: core.Style{Position: core.Absolute,
						Animated: &core.AnimatedStyle{Opacity: core.Interpolate(p, core.Interpolation{Input: []float64{from, from + 0.16}, Output: []float64{0, 1}})}},
						Children: []core.Node{Word(v, 15, colGreen, nil)}})
				}
				return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(44), Left: core.Pt(core.Points(left)), Width: core.Pt(70),
					Height: core.Pt(40), AlignItems: core.AlignCenter, Gap: 2}, Children: append(words, lpText(lab, colMuted, nil))}
			}
			return lpFrame(true, core.View{Children: []core.Node{make1(16, 24), make1(100, 34), make1(184, 20)}})
		}})
}

// #084 Bento grid: one big tile and four small ones, filling by weight.
func lpBento(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "084", Width: 300, Title: "Bento grid", Duration: 2400 * ms, Accent: colViolet, Height: 110,
		Note: "The big tile lands first: hierarchy is a timing decision as much as a size one.",
		Stage: func(p *core.AnimatedValue) core.Node {
			tiles := []struct{ top, left, w, h, delay float32 }{
				{24, 12, 108, 74, 0}, {24, 126, 52, 34, 0.12}, {62, 126, 52, 36, 0.18},
				{24, 184, 52, 34, 0.24}, {62, 184, 52, 36, 0.3},
			}
			var nodes []core.Node
			for i, t := range tiles {
				q := core.Interpolate(p, core.Interpolation{Input: []float64{0.05 + float64(t.delay), 0.4 + float64(t.delay)},
					Output: []float64{0, 1}, Easing: bezier(core.EaseOut)})
				nodes = append(nodes, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(core.Points(t.top)),
					Left: core.Pt(core.Points(t.left)), Width: core.Pt(core.Points(t.w)), Height: core.Pt(core.Points(t.h)),
					CornerRadius: 6, Background: colWell, Padding: core.Edges{Top: 8, Left: 8},
					Animated: &core.AnimatedStyle{Opacity: q, Scale: core.Interpolate(q, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0.94, 1}})}},
					Children: []core.Node{
						Disc(8, []core.Color{colViolet, colCyan, colGreen, colAmber, colPink}[i], nil),
						core.View{Style: core.Style{Margin: core.Edges{Top: 6}}, Children: []core.Node{lpText(t.w*0.6, colLine, nil)}},
					}})
			}
			return lpFrame(true, core.View{Children: nodes})
		}})
}

// #085 CTA: a button that asks to be pressed without shouting.
func lpCTA(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "085", Width: 300, Title: "CTA attention", Duration: 2000 * ms, Accent: colCyan, Height: 110, Loop: true,
		Note: "One breath, one ring, and nothing else on the page moves: a call to action that repeats its own motion becomes noise.",
		Stage: func(p *core.AnimatedValue) core.Node {
			pulse := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.25, 0.5, 1}, Output: []float64{1, 1.05, 1, 1}})
			ring := core.Interpolate(p, core.Interpolation{Input: []float64{0.2, 0.9}, Output: []float64{0, 1}})
			return lpFrame(true,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(34), Left: core.Pt(0), Right: core.Pt(0),
					AlignItems: core.AlignCenter, Gap: 8}, Children: []core.Node{
					lpText(120, colText, nil),
					lpText(90, colMuted, nil),
					core.View{Style: core.Style{Position: core.Relative, Width: core.Pt(90), Height: core.Pt(22), AlignItems: core.AlignCenter,
						JustifyContent: core.JustifyCenter}, Children: []core.Node{
						core.View{Style: core.Style{Position: core.Absolute, Width: core.Pt(90), Height: core.Pt(20), CornerRadius: 10,
							BorderWidth: 1, BorderColor: colCyan,
							Animated: &core.AnimatedStyle{Opacity: core.Interpolate(ring, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0.7, 0}}),
								Scale: core.Interpolate(ring, core.Interpolation{Input: []float64{0, 1}, Output: []float64{1, 1.35}})}}},
						lpButton(84, colCyan, &core.AnimatedStyle{Scale: pulse}),
					}},
				}},
			)
		}})
}

// #086 FAQ accordion: one opens, the other closes, and the layout moves with them.
func lpFAQ(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "086", Width: 300, Title: "FAQ accordion", Duration: 2400 * ms, Accent: colAmber, Height: 110,
		Note: "The panel does not fade, it grows: height is what tells you something opened, so the height is what animates.",
		Stage: func(p *core.AnimatedValue) core.Node {
			first := core.Interpolate(p, core.Interpolation{Input: []float64{0.1, 0.5}, Output: []float64{0, 1}, Easing: bezier(core.EaseInOut)})
			second := core.Interpolate(p, core.Interpolation{Input: []float64{0.5, 0.9}, Output: []float64{0, 1}, Easing: bezier(core.EaseInOut)})
			item := func(top float32, open core.AnimatedNode, col core.Color) core.Node {
				return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(core.Points(top)), Left: core.Pt(14), Width: core.Pt(200),
					Animated: &core.AnimatedStyle{TranslateY: core.Interpolate(open, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, -14}})}},
					Children: []core.Node{
						core.View{Style: core.Style{Direction: core.Row, AlignItems: core.AlignCenter, Gap: 6}, Children: []core.Node{
							lpText(120, colText, nil),
							core.View{Style: core.Style{FlexGrow: 1}},
							core.View{Style: core.Style{Width: core.Pt(6), Height: core.Pt(1.5), Background: col,
								Animated: &core.AnimatedStyle{Rotate: core.Interpolate(open, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 90}})}}},
						}},
						core.View{Style: core.Style{Height: core.Pt(0), Overflow: core.Hidden, Animated: &core.AnimatedStyle{ScaleY: open}},
							Children: []core.Node{
								core.View{Style: core.Style{Height: core.Pt(26), Padding: core.Edges{Top: 5}, Gap: 3}, Children: []core.Node{
									lpText(170, colMuted, nil), lpText(150, colMuted, nil),
								}},
							}},
					}}
			}
			return lpFrame(true,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(80), Left: core.Pt(14), Right: core.Pt(14), Height: core.Pt(1), Background: colLine}},
				item(28, first, colAmber),
				item(62, second, colAmber),
			)
		}})
}

// #087 Blog card hover: lift, shadow, and the image pushes in.
func lpBlogCard(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "087", Width: 300, Title: "Blog card hover", Duration: 2000 * ms, Accent: colPink, Height: 110,
		Note: "Three small movements on one hover: up, in, and a line that draws under the title — any one alone is not enough.",
		Stage: func(p *core.AnimatedValue) core.Node {
			up := core.Interpolate(p, core.Interpolation{Input: []float64{0.1, 0.6}, Output: []float64{0, 1}, Easing: bezier(core.EaseInOut)})
			return lpFrame(true,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(26), Left: core.Pt(14), Width: core.Pt(96), Height: core.Pt(72),
					CornerRadius: 7, Background: colWell, Overflow: core.Hidden,
					Animated: &core.AnimatedStyle{TranslateY: core.Interpolate(up, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, -5}}),
						Scale: core.Interpolate(up, core.Interpolation{Input: []float64{0, 1}, Output: []float64{1, 1.03}})}},
					Children: []core.Node{
						core.View{Style: core.Style{Height: core.Pt(38), Overflow: core.Hidden}, Children: []core.Node{
							core.View{Style: core.Style{Width: core.Pt(96), Height: core.Pt(38), Background: colBlue,
								Animated: &core.AnimatedStyle{Scale: core.Interpolate(up, core.Interpolation{Input: []float64{0, 1}, Output: []float64{1, 1.12}})}}},
						}},
						core.View{Style: core.Style{Padding: core.Edges{Top: 6, Left: 7}, Gap: 4}, Children: []core.Node{
							lpText(64, colText, nil),
							core.View{Style: core.Style{Height: core.Pt(1.5), Width: core.Pt(24), Background: colPink,
								Animated: &core.AnimatedStyle{ScaleX: up}}},
							lpText(50, colMuted, nil),
						}},
					}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(26), Left: core.Pt(122), Width: core.Pt(96), Height: core.Pt(72),
					CornerRadius: 7, Background: colWell, Padding: core.Edges{Top: 8, Left: 8}, Gap: 4}, Children: []core.Node{
					lpText(64, colLine, nil), lpText(50, colLine, nil),
				}},
			)
		}})
}

// #088 Team grid: eight faces, arriving in a wave.
func lpTeam(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "088", Width: 300, Title: "Team grid", Duration: 2200 * ms, Accent: colCyan, Height: 110,
		Note: "The wave runs along the row and then overflows to the next one, so the grid reads as people, not as tiles.",
		Stage: func(p *core.AnimatedValue) core.Node {
			var people []core.Node
			for i := 0; i < 8; i++ {
				delay := 0.06 * float64(i)
				q := core.Interpolate(p, core.Interpolation{Input: []float64{0.05 + delay, 0.35 + delay}, Output: []float64{0, 1}, Easing: bezier(core.EaseOut)})
				people = append(people, core.View{Style: core.Style{Position: core.Absolute,
					Top: core.Pt(core.Points(float32(26 + (i/4)*42))), Left: core.Pt(core.Points(float32(16 + (i%4)*56))),
					Width: core.Pt(44), AlignItems: core.AlignCenter, Gap: 4,
					Animated: &core.AnimatedStyle{Opacity: q, Scale: core.Interpolate(q, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0.85, 1}})}},
					Children: []core.Node{
						Disc(20, []core.Color{colCyan, colViolet, colGreen, colAmber, colPink, colBlue, colCyan, colViolet}[i], nil),
						lpText(30, colLine, nil),
					}})
			}
			return lpFrame(true, core.View{Children: people})
		}})
}

// #089 Footer reveal: the page ends by uncovering the footer.
func lpFooter(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "089", Width: 300, Title: "Footer reveal", Duration: 2400 * ms, Accent: colMuted, Height: 110,
		Note: "The footer is already there; the page above it slides away, which is cheaper than animating four columns.",
		Stage: func(p *core.AnimatedValue) core.Node {
			lift := core.Interpolate(p, core.Interpolation{Input: []float64{0.15, 0.8}, Output: []float64{0, 1}, Easing: bezier(core.EaseInOut)})
			return core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
				Overflow: core.Hidden, Background: colBar}, Children: []core.Node{
				// The footer, already in place at the bottom.
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(68), Left: core.Pt(10), Right: core.Pt(10), Bottom: core.Pt(8),
					Direction: core.Row, Gap: 8}, Children: []core.Node{
					core.View{Style: core.Style{Width: core.Pt(60), Gap: 3}, Children: []core.Node{lpText(40, colFaint, nil), lpText(52, colLine, nil), lpText(34, colLine, nil)}},
					core.View{Style: core.Style{Width: core.Pt(60), Gap: 3}, Children: []core.Node{lpText(40, colFaint, nil), lpText(46, colLine, nil), lpText(30, colLine, nil)}},
					core.View{Style: core.Style{Width: core.Pt(60), Gap: 3}, Children: []core.Node{lpText(40, colFaint, nil), lpText(50, colLine, nil)}},
				}},
				// The page body, leaving upwards.
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Height: core.Pt(96),
					Background: colCard, Padding: core.Edges{Top: 16, Left: 12}, Gap: 6,
					Animated: &core.AnimatedStyle{TranslateY: core.Interpolate(lift, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, -96}})}},
					Children: []core.Node{lpText(110, colText, nil), lpText(86, colMuted, nil), lpText(120, colLine, nil), lpText(70, colLine, nil)}},
			}}
		}})
}

// #090 Cookie banner: it arrives from the bottom, and it leaves the same way.
func lpCookie(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "090", Width: 300, Title: "Cookie banner", Duration: 2400 * ms, Accent: colAmber, Height: 110,
		Note: "A banner is a question, so it comes up, waits, and goes down — the press is what ends it, not the timer.",
		Stage: func(p *core.AnimatedValue) core.Node {
			up := core.Interpolate(p, core.Interpolation{Input: []float64{0.05, 0.3, 0.6, 0.85}, Output: []float64{0, 1, 1, 0}, Easing: bezier(core.EaseInOut)})
			return lpFrame(true,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(30), Left: core.Pt(14), Right: core.Pt(14), Gap: 5}, Children: []core.Node{
					lpText(150, colLine, nil), lpText(120, colLine, nil), lpText(90, colLine, nil),
				}},
				core.View{Style: core.Style{Position: core.Absolute, Bottom: core.Pt(8), Left: core.Pt(10), Right: core.Pt(10), Height: core.Pt(30),
					CornerRadius: 6, Background: colBar, Padding: core.Edges{Left: 8, Right: 8}, Direction: core.Row, AlignItems: core.AlignCenter, Gap: 8,
					Animated: &core.AnimatedStyle{TranslateY: core.Interpolate(up, core.Interpolation{Input: []float64{0, 1}, Output: []float64{48, 0}})}},
					Children: []core.Node{
						lpText(90, colMuted, nil),
						core.View{Style: core.Style{FlexGrow: 1}},
						lpButton(30, colAmber, nil),
					}},
			)
		}})
}

// #091 Cookie progress: consent has a shelf life, so it has a bar.
func lpCookieBar(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "091", Width: 300, Title: "Consent progress", Duration: 2600 * ms, Accent: colGreen, Height: 110,
		Note: "The bar is the only thing moving and it moves linearly: a countdown that eases is a countdown you misread.",
		Stage: func(p *core.AnimatedValue) core.Node {
			left := core.Interpolate(p, core.Interpolation{Input: []float64{0, 1}, Output: []float64{1, 0}})
			return lpFrame(true,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(20), Left: core.Pt(12), Right: core.Pt(12), Gap: 6}, Children: []core.Node{
					core.View{Style: core.Style{Direction: core.Row, AlignItems: core.AlignCenter, Gap: 6}, Children: []core.Node{
						lpText(60, colText, nil),
						core.View{Style: core.Style{FlexGrow: 1}},
						lpText(24, colMuted, nil),
					}},
					core.View{Style: core.Style{Height: core.Pt(4), CornerRadius: 2, Background: colWell}},
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(31), Left: core.Pt(12), Height: core.Pt(4), Width: core.Pt(200),
						CornerRadius: 2, Background: colGreen, Animated: &core.AnimatedStyle{ScaleX: left}}},
				}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(52), Left: core.Pt(12), Right: core.Pt(12), Gap: 5}, Children: []core.Node{
					lpText(170, colLine, nil), lpText(140, colLine, nil), lpText(100, colLine, nil),
				}},
			)
		}})
}

// #092 Search overlay: the field opens, the results fill in behind it.
func lpSearch(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "092", Width: 300, Title: "Search overlay", Duration: 2400 * ms, Accent: colCyan, Height: 110,
		Note: "The overlay scales up from the field that opened it: the origin is the point, otherwise it is just a fade.",
		Stage: func(p *core.AnimatedValue) core.Node {
			open := core.Interpolate(p, core.Interpolation{Input: []float64{0.1, 0.5}, Output: []float64{0, 1}, Easing: bezier(core.EaseOut)})
			var results []core.Node
			for i := 0; i < 4; i++ {
				delay := 0.4 + 0.07*float64(i)
				results = append(results, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(core.Points(float32(52 + i*14))),
					Left: core.Pt(18), Animated: &core.AnimatedStyle{Opacity: core.Interpolate(p, core.Interpolation{Input: []float64{delay, delay + 0.2}, Output: []float64{0, 1}}),
						TranslateX: core.Interpolate(p, core.Interpolation{Input: []float64{delay, delay + 0.2}, Output: []float64{-8, 0}})}},
					Children: []core.Node{lpText(110-float32(i)*14, colLine, nil)}})
			}
			return lpFrame(false,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					Background: core.RGB(0x080c12), Animated: &core.AnimatedStyle{Opacity: core.Interpolate(open, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 0.92}})}}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(20), Left: core.Pt(14), Right: core.Pt(14), Height: core.Pt(20),
					CornerRadius: 5, Background: colWell, BorderWidth: 1, BorderColor: colCyan, Padding: core.Edges{Left: 8},
					Direction: core.Row, AlignItems: core.AlignCenter, Gap: 5,
					Animated: &core.AnimatedStyle{ScaleX: core.Interpolate(open, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0.6, 1}})}},
					Children: []core.Node{
						Disc(5, colCyan, nil),
						lpText(64, colMuted, &core.AnimatedStyle{ScaleX: open}),
						core.View{Style: core.Style{Width: core.Pt(1), Height: core.Pt(10), Background: colCyan,
							Animated: &core.AnimatedStyle{Opacity: core.Interpolate(p, core.Interpolation{Input: []float64{0.3, 0.5, 0.7, 0.9}, Output: []float64{1, 0, 1, 0}})}}},
					}},
				core.View{Children: results},
			)
		}})
}

// #093 Signup flow: three fields, then a fourth that was not there.
func lpSignup(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "093", Width: 300, Title: "Signup flow", Duration: 2600 * ms, Accent: colCyan, Height: 110,
		Note: "One field at a time, and the button never moves until the last one lands: a form that grows upward keeps its target still.",
		Stage: func(p *core.AnimatedValue) core.Node {
			fields := []core.AnimatedNode{
				core.Interpolate(p, core.Interpolation{Input: []float64{0.1, 0.3}, Output: []float64{0, 1}}),
				core.Interpolate(p, core.Interpolation{Input: []float64{0.3, 0.5}, Output: []float64{0, 1}}),
				core.Interpolate(p, core.Interpolation{Input: []float64{0.5, 0.7}, Output: []float64{0, 1}}),
				core.Interpolate(p, core.Interpolation{Input: []float64{0.7, 0.9}, Output: []float64{0, 1}}),
			}
			var rows []core.Node
			for i, q := range fields {
				rows = append(rows, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(core.Points(float32(30 + i*17))),
					Left: core.Pt(40), Width: core.Pt(140), Height: core.Pt(14), CornerRadius: 4, Background: colWell,
					BorderWidth: 1, BorderColor: colLine, Animated: &core.AnimatedStyle{Opacity: q,
						TranslateY: core.Interpolate(q, core.Interpolation{Input: []float64{0, 1}, Output: []float64{-6, 0}})}}})
			}
			return lpFrame(true,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(20), Left: core.Pt(40)}, Children: []core.Node{lpText(70, colText, nil)}},
				core.View{Children: rows},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(96), Left: core.Pt(40)}, Children: []core.Node{
					lpButton(60, colCyan, nil),
				}},
			)
		}})
}

// #094 Onboarding checklist: each tick lands, the bar moves once.
func lpOnboarding(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "094", Width: 300, Title: "Onboarding checklist", Duration: 2400 * ms, Accent: colGreen, Height: 110,
		Note: "Three ticks and one bar: the bar is the sum of them, so it must move exactly when a tick lands.",
		Stage: func(p *core.AnimatedValue) core.Node {
			var rows []core.Node
			var done int8
			for i := 0; i < 3; i++ {
				at := 0.2 + 0.22*float64(i)
				q := core.Interpolate(p, core.Interpolation{Input: []float64{at, at + 0.15}, Output: []float64{0, 1}})
				_ = done
				rows = append(rows, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(core.Points(float32(38 + i*20))),
					Left: core.Pt(16), Direction: core.Row, AlignItems: core.AlignCenter, Gap: 7}, Children: []core.Node{
					core.View{Style: core.Style{Width: core.Pt(11), Height: core.Pt(11), CornerRadius: 3, BorderWidth: 1, BorderColor: colLine,
						Background: colWell, Animated: &core.AnimatedStyle{Scale: core.Interpolate(q, core.Interpolation{Input: []float64{0, 1}, Output: []float64{1, 1.1}})}},
						Children: []core.Node{
							core.View{Style: core.Style{Width: core.Pt(7), Height: core.Pt(7), CornerRadius: 2, Background: colGreen,
								Animated: &core.AnimatedStyle{Opacity: q, Scale: q}}},
						}},
					core.View{Style: core.Style{Animated: &core.AnimatedStyle{Opacity: core.Interpolate(q, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0.45, 1}})}},
						Children: []core.Node{lpText(90-float32(i)*10, colText, nil)}},
				}})
			}
			progress := core.Interpolate(p, core.Interpolation{Input: []float64{0.2, 0.84}, Output: []float64{0, 1}})
			return lpFrame(true,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(20), Left: core.Pt(16), Right: core.Pt(16), Height: core.Pt(4),
					CornerRadius: 2, Background: colWell}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(20), Left: core.Pt(16), Width: core.Pt(200), Height: core.Pt(4),
					CornerRadius: 2, Background: colGreen, Animated: &core.AnimatedStyle{ScaleX: progress}}},
				core.View{Children: rows},
			)
		}})
}

// #095 Empty dashboard: the blank state teaches the first action.
func lpEmpty(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "095", Width: 300, Title: "Empty state", Duration: 2200 * ms, Accent: colViolet, Height: 110,
		Note: "The illustration is drawn line by line, then the arrow, then the button: an empty state is a first instruction.",
		Stage: func(p *core.AnimatedValue) core.Node {
			draw := core.Interpolate(p, core.Interpolation{Input: []float64{0.1, 0.6}, Output: []float64{0, 1}, Easing: bezier(core.EaseInOut)})
			return lpFrame(true,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(28), Left: core.Pt(0), Right: core.Pt(0),
					AlignItems: core.AlignCenter, Gap: 6}, Children: []core.Node{
					core.View{Style: core.Style{Width: core.Pt(56), Height: core.Pt(40), CornerRadius: 5, BorderWidth: 1, BorderColor: colViolet,
						AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter,
						Animated: &core.AnimatedStyle{ScaleX: draw, Opacity: draw}}, Children: []core.Node{
						core.View{Style: core.Style{Width: core.Pt(6), Height: core.Pt(2), Background: colViolet}},
						core.View{Style: core.Style{Width: core.Pt(2), Height: core.Pt(10), Background: colViolet, Margin: core.Edges{Top: -6}}},
					}},
					core.View{Style: core.Style{Animated: &core.AnimatedStyle{Opacity: core.Interpolate(p, core.Interpolation{Input: []float64{0.6, 0.8}, Output: []float64{0, 1}})}},
						Children: []core.Node{lpText(66, colMuted, nil)}},
					core.View{Style: core.Style{Animated: &core.AnimatedStyle{Opacity: core.Interpolate(p, core.Interpolation{Input: []float64{0.75, 0.95}, Output: []float64{0, 1}})}},
						Children: []core.Node{lpButton(50, colViolet, nil)}},
				}},
			)
		}})
}

// #096 Loading skeleton: the placeholder and the real content swap.
func lpSkeleton(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "096", Width: 300, Title: "Loading skeleton", Duration: 2200 * ms, Accent: colMuted, Height: 110,
		Note: "The skeleton pulses while it waits, and the swap is a cross-fade in place: nothing jumps, because nothing moved.",
		Stage: func(p *core.AnimatedValue) core.Node {
			wait := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.6, 1}, Output: []float64{1, 1, 0}})
			real := core.Interpolate(p, core.Interpolation{Input: []float64{0.6, 0.85}, Output: []float64{0, 1}})
			shimmer := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.3, 0.6}, Output: []float64{0.35, 0.8, 0.35}})
			var bones, realLines []core.Node
			for i := 0; i < 3; i++ {
				bones = append(bones, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(core.Points(float32(26 + i*20))),
					Left: core.Pt(14), Width: core.Pt(core.Points(float32(170 - i*30))), Height: core.Pt(12), CornerRadius: 4, Background: colWell,
					Animated: &core.AnimatedStyle{Opacity: shimmer}}})
				realLines = append(realLines, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(core.Points(float32(26 + i*20))),
					Left: core.Pt(14), Width: core.Pt(core.Points(float32(150 - i*20))), Height: core.Pt(core.Points(float32(i*4 + 3))), CornerRadius: 2,
					Background: []core.Color{colText, colMuted, colLine}[i], Animated: &core.AnimatedStyle{Opacity: real}}})
			}
			return lpFrame(true,
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					Animated: &core.AnimatedStyle{Opacity: wait}}, Children: bones},
				core.View{Children: realLines},
			)
		}})
}
