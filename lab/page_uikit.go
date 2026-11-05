package main

// The ui kit page rebuilds effects from the "ui" pack of motion-prompt
// (github.com/showheyohtaki/motion-prompt): 24 named UI component
// motions, each with the timing its catalog gives. Only the names and
// timings come from there; every effect is written from scratch with
// NSGo's own primitives, so each card shows the API that does the work.
//
// The rule every card follows: state signals are read inside closures
// (Show.When, StyleFunc, ContentFunc), never at render time. Reading one
// during render would rebuild this subtree on every change, and the
// animated values created here would restart with it.

import (
	"strconv"
	"time"

	"github.com/oneo-me/nsgo/base"
	"github.com/oneo-me/nsgo/core"
)

func uikitPage(cx *core.Scope) core.Node {
	return core.View{Style: stageFill, Children: []core.Node{
		core.View{Style: core.Style{Padding: core.Edges{Top: 14, Right: 18, Left: 18}}, Children: []core.Node{
			note("Names and timings from the ui pack of motion-prompt (github.com/showheyohtaki/motion-prompt). The code is NSGo's own."),
		}},
		grid(
			accordionCard(cx),
			toastCard(cx),
			toggleCard(cx),
			segmentedCard(cx),
			stepperCard(cx),
			badgeCard(cx),
			skeletonCard(cx),
			spinnerCard(cx),
			tooltipCard(cx),
			hamburgerCard(cx),
			slidingTabsCard(cx),
			carouselCard(cx),
		),
		core.View{Style: core.Style{Padding: core.Edges{Top: 4, Right: 18, Left: 18}}, Children: []core.Node{
			note("The rest of the pack — 049 to 072 — as tiles: the same names and timings, one effect each."),
		}},
		grid(uiPackTiles(cx)...),
	}}
}

// pill is a small label on a demo surface.
func pill(text string, col core.Color) core.Node {
	return core.View{
		Style: core.Style{Padding: core.Symmetric(4, 9), CornerRadius: 7, Background: core.RGBA(0xffffff, 0.06)},
		Children: []core.Node{core.Text{
			Content: text,
			Style:   core.Style{FontSize: 11, FontWeight: core.FontWeightMedium, Color: col},
		}},
	}
}

// #052 Accordion: the clicked item opens over 0.35s and its arrow turns
// 180 degrees. AnimateLayout animates the frame change; the arrow is a
// timing.
func accordionCard(cx *core.Scope) core.Node {
	open := core.NewSignal(cx, true)
	turn := core.NewAnimatedValue(cx, 1)

	toggle := func() {
		next := !open.Get()
		// AnimateLayout applies to the next commit that has any, so it
		// goes before the signal that causes it.
		cx.AnimateLayout(core.LayoutAnimation{
			Update: &core.LayoutTransition{Curve: core.LayoutEaseInEaseOut, Duration: 350 * time.Millisecond},
			Create: &core.LayoutTransition{Curve: core.LayoutEaseOut, Duration: 200 * time.Millisecond, Property: core.LayoutOpacity},
			Delete: &core.LayoutTransition{Curve: core.LayoutEaseIn, Duration: 150 * time.Millisecond, Property: core.LayoutOpacity},
		})
		open.Set(next)
		target := 0.0
		if next {
			target = 1
		}
		cx.StartAnimation(core.Timing(turn, target, core.TimingConfig{
			Duration: 350 * time.Millisecond, Easing: bezier(core.EaseOut),
		}), nil)
	}

	head := base.Button{
		Label:   "Toggle the item",
		TestID:  "accordion",
		OnPress: func(core.PressEvent) { toggle() },
		Style:   core.Style{Direction: core.Row, AlignItems: core.AlignCenter, Gap: 8, Padding: core.Symmetric(8, 10), CornerRadius: 8, Background: core.RGBA(0xffffff, 0.05)},
		Children: []core.Node{
			core.Text{Content: "What animates a layout change?", Style: core.Style{FlexGrow: 1, FontSize: 12, Color: colText}},
			core.Text{Content: "›", Style: core.Style{FontSize: 13, Color: colCyan, Animated: &core.AnimatedStyle{Rotate: lerp(turn, 0, 90)}}},
		},
	}
	body := core.Text{
		Content: "Scope.AnimateLayout takes a LayoutAnimation: a curve and a duration per kind of change. Create and Delete animate an appearing or vanishing view's opacity; Update eases the frames in between.",
		Style:   core.Style{FontSize: 12, Color: colMuted, LineHeight: 18},
	}

	return card("Accordion · #052", "Opens over 0.35s, and the chevron turns 90 degrees.", colCyan,
		core.View{Style: core.Style{Gap: 8}, Children: []core.Node{
			head,
			core.Show{When: open.Get, Then: func(*core.Scope) core.Node { return body }},
		}},
	)
}

// #056 Toast notification: slides in from the right in 0.35s and
// dismisses itself after 2s.
func toastCard(cx *core.Scope) core.Node {
	shown := core.NewSignal(cx, false)
	enter := core.NewAnimatedValue(cx, 0)
	var pending func() bool

	show := func() {
		if pending != nil {
			pending()
			pending = nil
		}
		shown.Set(true)
		enter.Set(0)
		cx.StartAnimation(core.Timing(enter, 1, core.TimingConfig{
			Duration: 350 * time.Millisecond, Easing: bezier(core.EaseOut),
		}), nil)
		pending = cx.AfterFunc(2*time.Second, func() {
			pending = nil
			cx.StartAnimation(core.Timing(enter, 0, core.TimingConfig{
				Duration: 250 * time.Millisecond, Easing: bezier(core.EaseIn),
			}), func(bool) { shown.Set(false) })
		})
	}

	toast := core.View{
		Style: core.Style{
			Direction: core.Row, AlignItems: core.AlignCenter, Gap: 8, Padding: core.Symmetric(8, 12),
			CornerRadius: 9, Background: core.RGB(0x14301f), BorderWidth: 1, BorderColor: core.RGBA(0x39d98a, 0.5),
			Animated: &core.AnimatedStyle{TranslateX: lerp(enter, 150, 0), Opacity: enter},
		},
		Children: []core.Node{
			core.Text{Content: "done", Style: core.Style{FontSize: 11, Color: colGreen}},
			core.Text{Content: "Saved", Style: core.Style{FontSize: 12, Color: colText}},
		},
	}

	return card("Toast notification · #056", "Slides in from the right in 0.35s, then goes away on its own.", colGreen,
		core.View{Style: core.Style{Gap: 8}, Children: []core.Node{
			core.View{Style: core.Style{Direction: core.Row, Height: core.Pt(42), AlignItems: core.AlignCenter, Overflow: core.Hidden}, Children: []core.Node{
				core.Show{When: shown.Get, Then: func(*core.Scope) core.Node { return toast }},
			}},
			chip("Show a toast", colGreen, show),
			note("cx.AfterFunc arms the dismissal; its stop function cancels a toast pressed twice."),
		}},
	)
}

// #059 Toggle switch: the knob slides in 0.2s and the track turns green
// when on. One 0..1 value drives the slide and the colour.
func toggleCard(cx *core.Scope) core.Node {
	on := core.NewSignal(cx, false)
	p := core.NewAnimatedValue(cx, 0)

	flip := func() {
		next := !on.Get()
		on.Set(next)
		target := 0.0
		if next {
			target = 1
		}
		cx.StartAnimation(core.Timing(p, target, core.TimingConfig{
			Duration: 200 * time.Millisecond, Easing: bezier(core.EaseOut),
		}), nil)
	}

	// Opaque colours on purpose: an interpolated background is a paint, and
	// 12% white over a dark card reads as a light grey track rather than a
	// switched-off one.
	track := core.Interpolate(p, core.Interpolation{
		Input:  []float64{0, 1},
		Colors: []core.Color{core.RGB(0x2b3644), core.RGB(0x2a9d63)},
	})
	knob := base.Button{
		Label:   "Toggle the setting",
		TestID:  "toggle",
		OnPress: func(core.PressEvent) { flip() },
		Style: core.Style{
			Width: core.Pt(46), Height: core.Pt(26), CornerRadius: 13, Padding: core.All(3),
			JustifyContent: core.JustifyCenter,
			Animated:       &core.AnimatedStyle{Background: track},
		},
		Children: []core.Node{core.View{
			Style: core.Style{
				Width: core.Pt(20), Height: core.Pt(20), CornerRadius: 10, Background: core.RGB(0xf7fbff),
				Animated: &core.AnimatedStyle{TranslateX: lerp(p, 0, 20)},
			},
		}},
	}

	return card("Toggle switch · #059", "The knob slides in 0.2s and the track turns green.", colGreen,
		row(12,
			knob,
			core.Text{ContentFunc: func() string {
				if on.Get() {
					return "notifications on"
				}
				return "notifications off"
			}, Style: core.Style{FontSize: 12, Color: colMuted}},
		),
		note("Interpolate maps the same 0..1 node to a colour, so the fill and the knob cannot drift apart."),
	)
}

// #060 Segmented control: the selected background slides over 0.25s and
// the data changes with it.
func segmentedCard(cx *core.Scope) core.Node {
	const segW = 78.0
	periods := []string{"Day", "Week", "Month"}
	// Seven bars per period, in eighths of the chart height.
	bars := [3][7]float32{
		{4, 7, 5, 9, 6, 8, 5},
		{6, 8, 7, 10, 9, 9, 7},
		{5, 9, 8, 11, 10, 12, 9},
	}

	sel := core.NewSignal(cx, 0)
	pos := core.NewAnimatedValue(cx, 0)
	fade := core.NewAnimatedValue(cx, 1)

	pick := func(i int) {
		sel.Set(i)
		cx.StartAnimation(core.Timing(pos, float64(i)*segW, core.TimingConfig{
			Duration: 250 * time.Millisecond, Easing: bezier(core.EaseInOut),
		}), nil)
		fade.Set(0.3)
		cx.StartAnimation(core.Timing(fade, 1, core.TimingConfig{
			Duration: 200 * time.Millisecond, Easing: bezier(core.Linear),
		}), nil)
	}

	segs := make([]core.Node, len(periods))
	for i, name := range periods {
		i, name := i, name
		segs[i] = base.Button{
			Label:   name,
			TestID:  "seg-" + name,
			OnPress: func(core.PressEvent) { pick(i) },
			Style:   core.Style{Width: core.Pt(segW), Height: core.Pt(28), AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter},
			Children: []core.Node{core.Text{
				Content: name,
				StyleFunc: func() core.Style {
					colour := colMuted
					if sel.Get() == i {
						colour = colText
					}
					return core.Style{FontSize: 12, FontWeight: core.FontWeightMedium, Color: colour}
				},
			}},
		}
	}

	// Each bar reads the period in its own StyleFunc, so the change
	// costs seven restyles rather than a rebuild of the card.
	chart := make([]core.Node, 7)
	for i := range chart {
		i := i
		chart[i] = core.View{
			Style: core.Style{Width: core.Pt(14), CornerRadius: 4, Background: colBlue},
			StyleFunc: func() core.Style {
				return core.Style{Height: core.Pt(core.Points(8 * bars[sel.Get()][i]))}
			},
		}
	}

	return card("Segmented control · #060", "The selected background slides over 0.25s; the data cross-fades.", colBlue,
		core.View{Style: core.Style{Direction: core.Row, Padding: core.All(4), CornerRadius: 10, Background: core.RGBA(0xffffff, 0.05), Position: core.Relative}, Children: append([]core.Node{
			core.View{Style: core.Style{
				Position: core.Absolute, Top: core.Pt(4), Left: core.Pt(4), Width: core.Pt(segW), Height: core.Pt(28),
				CornerRadius: 7, Background: core.RGBA(0x4d8dff, 0.32),
				Animated: &core.AnimatedStyle{TranslateX: pos},
			}},
		}, segs...)},
		core.View{Style: core.Style{Direction: core.Row, AlignItems: core.AlignEnd, Gap: 5, Height: core.Pt(96), Animated: &core.AnimatedStyle{Opacity: fade}}, Children: chart},
	)
}

// #063 Stepper: plus and minus change the number by one, and it pops for
// 0.25s on each press.
func stepperCard(cx *core.Scope) core.Node {
	count := core.NewSignal(cx, 1)
	pop := core.NewAnimatedValue(cx, 1)

	bump := func(d int) {
		count.Update(func(n int) int {
			n += d
			if n < 0 {
				return 0
			}
			if n > 99 {
				return 99
			}
			return n
		})
		pop.Set(1)
		cx.StartAnimation(core.Sequence(
			core.Timing(pop, 1.3, core.TimingConfig{Duration: 100 * time.Millisecond, Easing: bezier(core.EaseOut)}),
			core.Timing(pop, 1, core.TimingConfig{Duration: 150 * time.Millisecond, Easing: bezier(core.EaseIn)}),
		), nil)
	}

	button := func(label, id string, d int) core.Node {
		return base.Button{
			Label:   label,
			TestID:  id,
			OnPress: func(core.PressEvent) { bump(d) },
			Style:   core.Style{Width: core.Pt(28), Height: core.Pt(28), CornerRadius: 8, Background: core.RGBA(0xffffff, 0.08), AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter},
			Children: []core.Node{core.Text{Content: label, Style: core.Style{FontSize: 14, Color: colText}}},
		}
	}

	return card("Stepper · #063", "Each press pops the number for 0.25s.", colAmber,
		row(10,
			button("-", "step-minus", -1),
			core.View{Style: core.Style{Width: core.Pt(46), AlignItems: core.AlignCenter, Animated: &core.AnimatedStyle{Scale: pop}}, Children: []core.Node{
				core.Text{TestID: "step-value", ContentFunc: func() string { return strconv.Itoa(count.Get()) }, Style: core.Style{FontSize: 18, FontWeight: core.FontWeightBold, Color: colAmber}},
			}},
			button("+", "step-plus", 1),
		),
		note("Sequence chains the pop and the settle; Scale is one AnimatedStyle field."),
	)
}

// #071 Badge: the count pops for 0.3s each time it goes up, and the
// badge goes away when it is read.
func badgeCard(cx *core.Scope) core.Node {
	count := core.NewSignal(cx, 0)
	pop := core.NewAnimatedValue(cx, 1)

	add := func() {
		count.Update(func(n int) int { return n + 1 })
		pop.Set(1)
		cx.StartAnimation(core.Sequence(
			core.Timing(pop, 1.45, core.TimingConfig{Duration: 120 * time.Millisecond, Easing: bezier(core.EaseOut)}),
			core.Timing(pop, 1, core.TimingConfig{Duration: 180 * time.Millisecond, Easing: bezier(core.EaseIn)}),
		), nil)
	}
	markRead := func() { count.Set(0) }

	bell := core.View{Style: core.Style{
		Position: core.Relative, Width: core.Pt(34), Height: core.Pt(34), CornerRadius: 9,
		Background: core.RGBA(0xffffff, 0.06), AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter,
	}, Children: []core.Node{
		core.Text{Content: "bell", Style: core.Style{FontSize: 10, Color: colMuted}},
		core.Show{When: func() bool { return count.Get() > 0 }, Then: func(*core.Scope) core.Node {
			return core.View{
				Style: core.Style{
					Position: core.Absolute, Top: core.Pt(-5), Right: core.Pt(-5), MinWidth: core.Pt(17), Height: core.Pt(17),
					Padding: core.Symmetric(0, 4), CornerRadius: 9, Background: colPink,
					AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter,
					Animated: &core.AnimatedStyle{Scale: pop},
				},
				Children: []core.Node{core.Text{
					TestID:      "badge-count",
					ContentFunc: func() string { return strconv.Itoa(count.Get()) },
					Style:       core.Style{FontSize: 10, FontWeight: core.FontWeightBold, Color: core.RGB(0x24101a)},
				}},
			}
		}},
	}}

	return card("Badge · #071", "Pops for 0.3s each time the count goes up, and hides when read.", colPink,
		row(12, bell, pill("unread", colFaint)),
		core.View{Style: core.Style{Direction: core.Row, Gap: 8}, Children: []core.Node{
			chip("Add one", colPink, add),
			chip("Mark read", colFaint, markRead),
		}},
	)
}

// #069 Skeleton screen: grey blocks in the shape of the content, with a
// shine sweeping left to right every 0.9s.
func skeletonCard(cx *core.Scope) core.Node {
	sheen := core.NewAnimatedValue(cx, 0)
	cx.StartAnimation(core.Loop(core.Timing(sheen, 1, core.TimingConfig{
		Duration: 900 * time.Millisecond, Easing: bezier(core.Linear),
	}), 0), nil)

	block := func(w, h, alpha float32) core.Node {
		return core.View{Style: core.Style{
			Width: core.Pt(core.Points(w)), Height: core.Pt(core.Points(h)),
			CornerRadius: 5, Background: core.RGBA(0xffffff, alpha),
		}}
	}
	// Three stacked bands of falling alpha read as a soft swipe without
	// a gradient; the well clips them.
	band := func(alpha float32, offset float64) core.Node {
		return core.View{Style: core.Style{
			Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Width: core.Pt(26), Height: core.Pct(100),
			Background: core.RGBA(0xffffff, alpha),
			Animated:   &core.AnimatedStyle{TranslateX: lerp(sheen, -120+offset, 300+offset)},
		}}
	}

	return card("Skeleton screen · #069", "A shine sweeps left to right every 0.9s while content loads.", colFaint,
		core.View{
			Style: well(96),
			Children: []core.Node{
				core.View{Style: core.Style{Padding: core.All(12), Gap: 8}, Children: []core.Node{
					row(10,
						core.View{Style: core.Style{Width: core.Pt(28), Height: core.Pt(28), CornerRadius: 14, Background: core.RGBA(0xffffff, 0.10)}},
						core.View{Style: core.Style{Gap: 6, FlexGrow: 1}, Children: []core.Node{block(120, 8, 0.10), block(78, 8, 0.06)}},
					),
					block(210, 8, 0.10),
					block(160, 8, 0.06),
				}},
				band(0.02, -14),
				band(0.05, 0),
				band(0.02, 14),
			},
		},
		note("Loop(anim, 0) repeats forever; the platform steps it, so the engine commits nothing per frame."),
	)
}

// #070 Spinner: turns a little over once per second, then becomes a
// check when the work is done.
func spinnerCard(cx *core.Scope) core.Node {
	spin := core.NewAnimatedValue(cx, 0)
	cx.StartAnimation(core.Loop(core.Timing(spin, 1, core.TimingConfig{
		Duration: time.Second, Easing: bezier(core.Linear),
	}), 0), nil)

	done := core.NewSignal(cx, false)
	check := core.NewAnimatedValue(cx, 0)

	finish := func() {
		next := !done.Get()
		done.Set(next)
		target := 0.0
		if next {
			target = 1
		}
		cx.StartAnimation(core.Spring(check, target, core.SpringConfig{Stiffness: 240, Damping: 14, Mass: 1}), nil)
	}

	ring := core.View{Style: core.Style{
		Width: core.Pt(34), Height: core.Pt(34), CornerRadius: 17,
		Borders:      core.Edges{Top: 3, Right: 3},
		BorderColors: core.EdgeColors{Top: colCyan, Right: core.RGBA(0x35d6de, 0.15)},
		Animated:     &core.AnimatedStyle{Rotate: lerp(spin, 0, 360)},
	}}
	tick := core.View{Style: core.Style{
		Width: core.Pt(34), Height: core.Pt(34), CornerRadius: 17, Background: core.RGBA(0x39d98a, 0.18),
		AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter,
		Animated:   &core.AnimatedStyle{Scale: check},
	}, Children: []core.Node{core.Text{Content: "ok", Style: core.Style{FontSize: 12, Color: colGreen}}}}

	return card("Spinner · #070", "A little over one turn per second; press Finish to swap it for a check.", colCyan,
		row(12,
			core.View{Style: core.Style{Width: core.Pt(34), Height: core.Pt(34), AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter}, Children: []core.Node{
				core.Show{When: func() bool { return !done.Get() }, Then: func(*core.Scope) core.Node { return ring }, Else: func(*core.Scope) core.Node { return tick }},
			}},
			core.Text{ContentFunc: func() string {
				if done.Get() {
					return "sent"
				}
				return "sending"
			}, Style: core.Style{FontSize: 12, Color: colMuted}},
		),
		chip("Finish", colGreen, finish),
		note("Only two of the ring's sides have a border, so a turning view reads as an arc."),
	)
}

// #054 Tooltip: shows the note beside the icon 0.2s after hover, and
// hides it when the cursor leaves.
func tooltipCard(cx *core.Scope) core.Node {
	shown := core.NewSignal(cx, false)
	enter := core.NewAnimatedValue(cx, 0)
	var pending func() bool

	hoverIn := func() {
		if pending != nil {
			pending()
		}
		pending = cx.AfterFunc(200*time.Millisecond, func() {
			pending = nil
			shown.Set(true)
			cx.StartAnimation(core.Timing(enter, 1, core.TimingConfig{
				Duration: 160 * time.Millisecond, Easing: bezier(core.EaseOut),
			}), nil)
		})
	}
	hoverOut := func() {
		if pending != nil {
			pending()
			pending = nil
		}
		cx.StartAnimation(core.Timing(enter, 0, core.TimingConfig{
			Duration: 120 * time.Millisecond, Easing: bezier(core.EaseIn),
		}), func(bool) { shown.Set(false) })
	}

	anchor := core.View{
		Style: core.Style{
			Position: core.Relative, Width: core.Pt(26), Height: core.Pt(26), CornerRadius: 13,
			Background: core.RGBA(0xffffff, 0.08), AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter,
		},
		OnHoverIn:  hoverIn,
		OnHoverOut: hoverOut,
		Children: []core.Node{
			core.Text{Content: "?", Style: core.Style{FontSize: 12, FontWeight: core.FontWeightBold, Color: colMuted}},
			core.Show{When: shown.Get, Then: func(*core.Scope) core.Node {
				// The width is explicit: an absolutely positioned view with
				// only a Left edge takes the space left over in its
				// containing block, which here is a 26pt circle, and the
				// label then wraps one character per line.
				return core.View{
					Style: core.Style{
						Position: core.Absolute, Left: core.Pt(34), Top: core.Pt(-4), Width: core.Pt(124),
						Padding: core.Symmetric(5, 9), CornerRadius: 7, Background: colBar,
						BorderWidth: 1, BorderColor: colLine,
						Animated: &core.AnimatedStyle{Opacity: enter, TranslateX: lerp(enter, -6, 0)},
					},
					Children: []core.Node{core.Text{Content: "Hover, not press", Style: core.Style{FontSize: 11, Color: colText}}},
				}
			}},
		},
	}

	return card("Tooltip · #054", "Hover the icon: the note comes 0.2s later, and leaves with the cursor.", colAmber,
		core.View{Style: core.Style{Height: core.Pt(56), JustifyContent: core.JustifyCenter, AlignItems: core.AlignStart}, Children: []core.Node{anchor}},
		note("OnHoverIn and OnHoverOut are View handlers; the 0.2s is cx.AfterFunc, cancelled if you leave first."),
	)
}

// #065 Hamburger menu: the three lines become a cross in 0.3s and the
// menu drops down from the top.
func hamburgerCard(cx *core.Scope) core.Node {
	open := core.NewSignal(cx, false)
	p := core.NewAnimatedValue(cx, 0)
	mounted := core.NewSignal(cx, false)

	flip := func() {
		next := !open.Get()
		open.Set(next)
		if next {
			mounted.Set(true)
		}
		target := 0.0
		if next {
			target = 1
		}
		cx.StartAnimation(core.Timing(p, target, core.TimingConfig{
			Duration: 300 * time.Millisecond, Easing: bezier(core.EaseInOut),
		}), func(bool) {
			// The panel stays mounted until the closing timing ends, so
			// it fades out instead of vanishing.
			if !open.Get() {
				mounted.Set(false)
			}
		})
	}

	line := func(top float64, anim *core.AnimatedStyle) core.Node {
		return core.View{Style: core.Style{
			Position: core.Absolute, Left: core.Pt(8), Top: core.Pt(core.Points(top)),
			Width: core.Pt(18), Height: core.Pt(2), CornerRadius: 1, Background: colText,
			Animated: anim,
		}}
	}

	icon := base.Button{
		Label:   "Menu",
		TestID:  "hamburger",
		OnPress: func(core.PressEvent) { flip() },
		Style:   core.Style{Position: core.Relative, Width: core.Pt(34), Height: core.Pt(26), CornerRadius: 8, Background: core.RGBA(0xffffff, 0.06)},
		Children: []core.Node{
			line(7, &core.AnimatedStyle{TranslateY: lerp(p, 0, 6), Rotate: lerp(p, 0, 45)}),
			line(11, &core.AnimatedStyle{Opacity: lerp(p, 1, 0)}),
			line(15, &core.AnimatedStyle{TranslateY: lerp(p, 0, -6), Rotate: lerp(p, 0, -45)}),
		},
	}

	panel := core.View{
		Style: core.Style{
			Width: core.Pt(96), Padding: core.All(8), Gap: 6, CornerRadius: 9,
			Background: colBar, BorderWidth: 1, BorderColor: colLine,
			Animated: &core.AnimatedStyle{Opacity: p, TranslateY: lerp(p, -6, 0)},
		},
		Children: []core.Node{
			core.Text{Content: "Home", Style: core.Style{FontSize: 11, Color: colMuted}},
			core.Text{Content: "Projects", Style: core.Style{FontSize: 11, Color: colMuted}},
			core.Text{Content: "Settings", Style: core.Style{FontSize: 11, Color: colMuted}},
		},
	}

	return card("Hamburger menu · #065", "Three lines cross in 0.3s; the menu drops from the top.", colViolet,
		core.View{Style: core.Style{Direction: core.Row, Gap: 14, Height: core.Pt(78), AlignItems: core.AlignStart}, Children: []core.Node{
			icon,
			core.Show{When: mounted.Get, Then: func(*core.Scope) core.Node { return panel }},
		}},
		note("Two edges of one 0..1 value: TranslateY folds the lines, Rotate crosses them."),
	)
}

// #057 Tabs: the underline slides over 0.25s and the content
// cross-fades in 0.2s.
func slidingTabsCard(cx *core.Scope) core.Node {
	const tabW = 68.0
	names := []string{"Overview", "Activity", "Billing"}
	texts := []string{"Three tabs, one underline that moves.", "Nothing here yet today.", "One invoice, paid."}

	sel := core.NewSignal(cx, 0)
	pos := core.NewAnimatedValue(cx, 0)
	fade := core.NewAnimatedValue(cx, 1)

	pick := func(i int) {
		sel.Set(i)
		cx.StartAnimation(core.Timing(pos, float64(i)*tabW, core.TimingConfig{
			Duration: 250 * time.Millisecond, Easing: bezier(core.EaseOut),
		}), nil)
		fade.Set(0.35)
		cx.StartAnimation(core.Timing(fade, 1, core.TimingConfig{
			Duration: 200 * time.Millisecond, Easing: bezier(core.Linear),
		}), nil)
	}

	tabs := make([]core.Node, len(names))
	for i, name := range names {
		i, name := i, name
		tabs[i] = base.Button{
			Label:   name,
			TestID:  "utab-" + name,
			OnPress: func(core.PressEvent) { pick(i) },
			Style:   core.Style{Width: core.Pt(tabW), Height: core.Pt(26), AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter},
			Children: []core.Node{core.Text{Content: name, StyleFunc: func() core.Style {
				colour := colMuted
				if sel.Get() == i {
					colour = colText
				}
				return core.Style{FontSize: 12, Color: colour}
			}}},
		}
	}

	return card("Tabs · #057", "The underline slides over 0.25s, the content cross-fades in 0.2s.", colViolet,
		core.View{Style: core.Style{Position: core.Relative, Direction: core.Row}, Children: append([]core.Node{
			core.View{Style: core.Style{
				Position: core.Absolute, Top: core.Pt(24), Left: core.Pt(0), Width: core.Pt(tabW), Height: core.Pt(2),
				CornerRadius: 1, Background: colViolet,
				Animated: &core.AnimatedStyle{TranslateX: pos},
			}},
		}, tabs...)},
		core.View{Style: core.Style{Height: core.Pt(20), Animated: &core.AnimatedStyle{Opacity: fade}}, Children: []core.Node{
			core.Text{ContentFunc: func() string { return texts[sel.Get()] }, Style: core.Style{FontSize: 12, Color: colMuted}},
		}},
	)
}

// #058 Carousel: the arrow moves one slide in 0.4s, with position dots
// below.
func carouselCard(cx *core.Scope) core.Node {
	const slideW = 196.0
	slides := []struct {
		label string
		col   core.Color
	}{
		{"01", colCyan}, {"02", colViolet}, {"03", colGreen},
	}
	idx := core.NewSignal(cx, 0)
	x := core.NewAnimatedValue(cx, 0)

	goTo := func(d int) {
		n := len(slides)
		next := ((idx.Get()+d)%n + n) % n
		idx.Set(next)
		cx.StartAnimation(core.Timing(x, -float64(next)*slideW, core.TimingConfig{
			Duration: 400 * time.Millisecond, Easing: bezier(core.EaseInOut),
		}), nil)
	}

	panels := make([]core.Node, len(slides))
	for i, s := range slides {
		panels[i] = core.View{
			Style: core.Style{
				Width: core.Pt(slideW), Height: core.Pt(88), CornerRadius: 10,
				Background: core.RGBA(0xffffff, 0.04), AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter,
			},
			Children: []core.Node{
				core.Text{Content: s.label, Style: core.Style{FontSize: 32, FontWeight: core.FontWeightBold, Color: s.col}},
			},
		}
	}

	dots := make([]core.Node, len(slides))
	for i := range slides {
		i := i
		dots[i] = core.View{
			Style: core.Style{Width: core.Pt(6), Height: core.Pt(6), CornerRadius: 3},
			StyleFunc: func() core.Style {
				colour := core.RGBA(0xffffff, 0.18)
				if idx.Get() == i {
					colour = colCyan
				}
				return core.Style{Background: colour}
			},
		}
	}

	return card("Carousel · #058", "One slide per press, 0.4s each way, with dots below.", colCyan,
		core.View{Style: core.Style{Width: core.Pt(slideW), Height: core.Pt(88), CornerRadius: 10, Overflow: core.Hidden}, Children: []core.Node{
			core.View{Style: core.Style{Direction: core.Row, Animated: &core.AnimatedStyle{TranslateX: x}}, Children: panels},
		}},
		core.View{Style: core.Style{Direction: core.Row, AlignItems: core.AlignCenter, Gap: 10}, Children: append(append([]core.Node{
			chip("Prev", colFaint, func() { goTo(-1) }),
		}, dots...), chip("Next", colCyan, func() { goTo(1) }))},
	)
}
