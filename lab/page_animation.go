package main

import (
	"time"

	"github.com/oneo-me/nsgo/base"
	"github.com/oneo-me/nsgo/core"
)

// animationPage collects the animation graph's examples. Everything here
// is stepped by the platform's frame clock: the engine does no work per
// frame while these run.
func animationPage(cx *core.Scope) core.Node {
	return grid(
		easingRace{}.Render(cx),
		springCompare(cx),
		keyframesDemo(cx),
		compositionDemo(cx),
		presetGrid{}.Render(cx),
	)
}

// travel is how far a dot moves across its track, in points.
const travel = 200.0

// lane is one labelled track with a dot on it.
type lane struct {
	name  string
	value *core.AnimatedValue
	col   core.Color
}

// tracks lays lanes out one per line.
func tracks(lanes []lane) core.Node {
	rows := make([]core.Node, len(lanes))
	for i, l := range lanes {
		rows[i] = row(10,
			core.Text{Content: l.name, Style: core.Style{FontSize: 11, Color: colMuted, Width: core.Pt(72)}},
			core.View{
				Style: core.Style{Width: core.Pt(212), Height: core.Pt(8), CornerRadius: 4, Background: core.RGBA(0xffffff, 0.06)},
				Children: []core.Node{
					core.View{Style: core.Style{
						Width: core.Pt(12), Height: core.Pt(12), CornerRadius: 6, Background: l.col,
						Animated: &core.AnimatedStyle{TranslateX: lerp(l.value, 0, travel)},
					}},
				},
			},
		)
	}
	return core.View{Style: core.Style{Gap: 9}, Children: rows}
}

// easingRace runs the same timing under four of CSS's curves.
type easingRace struct{}

func (easingRace) Render(cx *core.Scope) core.Node {
	defs := []struct {
		name  string
		ease  core.Easing
		col   core.Color
	}{
		{"linear", bezier(core.Linear), colFaint},
		{"ease", bezier(core.Ease), colCyan},
		{"ease-out", bezier(core.EaseOut), colGreen},
		{"ease-in-out", bezier(core.EaseInOut), colViolet},
	}
	lanes := make([]lane, len(defs))
	for i, d := range defs {
		lanes[i] = lane{name: d.name, value: core.NewAnimatedValue(cx, 0), col: d.col}
	}
	replay := func() {
		for i, d := range defs {
			lanes[i].value.Set(0)
			cx.StartAnimation(core.Timing(lanes[i].value, 1, core.TimingConfig{
				Duration: 1400 * time.Millisecond,
				Easing:   d.ease,
			}), nil)
		}
	}
	replay()
	return card("Easing curves", "A timing's Easing takes CSS's curves as cubic Béziers, or React Native's named ones.", colViolet,
		tracks(lanes),
		chip("Replay curves", colViolet, replay),
	)
}

// springCompare races a timing against a damped spring over the same
// distance, so the spring's overshoot shows.
func springCompare(cx *core.Scope) core.Node {
	lanes := []lane{
		{name: "timing", value: core.NewAnimatedValue(cx, 0), col: colBlue},
		{name: "spring", value: core.NewAnimatedValue(cx, 0), col: colAmber},
	}
	replay := func() {
		for _, l := range lanes {
			l.value.Set(0)
		}
		cx.StartAnimation(core.Timing(lanes[0].value, 1, core.TimingConfig{
			Duration: 900 * time.Millisecond,
			Easing:   bezier(core.EaseOut),
		}), nil)
		cx.StartAnimation(core.Spring(lanes[1].value, 1, core.SpringConfig{
			Stiffness: 230, Damping: 9, Mass: 1,
		}), nil)
	}
	replay()
	return card("Timing vs spring", "Spring is a damped spring from React Native; zero Stiffness, Damping and Mass take 100, 10 and 1.", colAmber,
		tracks(lanes),
		chip("Replay spring", colAmber, replay),
	)
}

// keyframesDemo plays one value through stops, the way Web Animations
// does it.
func keyframesDemo(cx *core.Scope) core.Node {
	bounce := core.NewAnimatedValue(cx, 0)
	replay := func() {
		bounce.Set(0)
		cx.StartAnimation(core.Loop(core.ValueKeyframes(bounce,
			core.KeyframeTiming{Duration: 1600 * time.Millisecond, Iterations: 1},
			core.StopAt(0, 0),
			core.StopAt(0.3, 1),
			core.StopAt(0.55, 0.32),
			core.StopAt(0.8, 0.88),
			core.StopAt(1, 0),
		), 0), nil)
	}
	replay()
	return card("Keyframes", "ValueKeyframes takes offsets and values as Web Animations spreads them.", colGreen,
		core.View{
			Style: core.Style{Height: core.Pt(110), Background: colWell, CornerRadius: 10, AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter},
			Children: []core.Node{
				core.View{Style: core.Style{
					Width: core.Pt(30), Height: core.Pt(30), CornerRadius: 8, Background: colGreen,
					Animated: &core.AnimatedStyle{TranslateY: lerp(bounce, 0, -34)},
				}},
			},
		},
		chip("Replay bounce", colGreen, replay),
	)
}

// compositionDemo runs six dots under one Stagger, then loops a
// Sequence that breathes.
func compositionDemo(cx *core.Scope) core.Node {
	const dots = 6
	values := make([]*core.AnimatedValue, dots)
	cells := make([]core.Node, dots)
	for i := range values {
		values[i] = core.NewAnimatedValue(cx, 0)
	}
	pulse := core.NewAnimatedValue(cx, 0)

	wave := func() {
		steps := make([]*core.Animation, dots)
		for i := range values {
			values[i].Set(0)
			steps[i] = core.Timing(values[i], 1, core.TimingConfig{
				Duration: 520 * time.Millisecond,
				Easing:   bezier(core.EaseInOut),
			})
		}
		// Stagger starts each step after the one before, so the dots
		// ripple instead of moving together.
		cx.StartAnimation(core.Stagger(70*time.Millisecond, steps...), nil)
	}
	wave()

	cx.StartAnimation(core.Loop(core.Sequence(
		core.Timing(pulse, 1, core.TimingConfig{Duration: 800 * time.Millisecond, Easing: bezier(core.EaseInOut)}),
		core.Timing(pulse, 0, core.TimingConfig{Duration: 800 * time.Millisecond, Easing: bezier(core.EaseInOut)}),
	), 0), nil)

	for i := range cells {
		cells[i] = core.View{Style: core.Style{
			Width: core.Pt(16), Height: core.Pt(16), CornerRadius: 8, Background: colBlue,
			Animated: &core.AnimatedStyle{
				TranslateY: lerp(values[i], 0, -22),
				Opacity:    lerp(values[i], 0.25, 1),
			},
		}}
	}
	return card("Composition", "Stagger, Sequence and Loop compose animations; each step starts where the last ended.", colBlue,
		core.View{
			Style: core.Style{Height: core.Pt(96), Background: colWell, CornerRadius: 10, Direction: core.Row, AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter, Gap: 10},
			Children: cells,
		},
		row(8,
			chip("Ripple", colBlue, wave),
			core.View{Style: core.Style{
				Width: core.Pt(16), Height: core.Pt(16), CornerRadius: 8, Background: colViolet,
				Animated: &core.AnimatedStyle{
					Scale:   lerp(pulse, 0.6, 1.5),
					Opacity: lerp(pulse, 0.45, 1),
				},
			}},
			note("looping breath"),
		),
	)
}

// presets are the MotionPresets shown on the grid: entrances, attention
// seekers and loops, all built from NSH's presets.
var presets = []core.MotionPreset{
	core.FadeIn, core.FadeInUp, core.FadeInDown, core.SlideInUp,
	core.ZoomIn, core.PopIn, core.BounceIn, core.FlipInX,
	core.Shake, core.Pulse, core.Heartbeat, core.Wiggle,
	core.Spin, core.Breathe, core.Float, core.Tada,
}

// presetGrid is a MotionPreset per cell, each replaying on press.
type presetGrid struct{}

func (presetGrid) Render(cx *core.Scope) core.Node {
	cells := make([]core.Node, len(presets))
	for i, p := range presets {
		cells[i] = presetCell{preset: p}.Render(cx)
	}
	return card("Motion presets", "A Motion moves opacity, translation, scale and rotation together; press one to play it.", colPink,
		core.View{Style: core.Style{Direction: core.Row, Wrap: core.Wrap, Gap: 8}, Children: cells},
	)
}

// presetCell is one square that plays its preset when pressed.
type presetCell struct {
	preset core.MotionPreset
}

func (c presetCell) Render(cx *core.Scope) core.Node {
	motion := core.NewMotion(cx, nil)
	return base.Button{
		Label:   c.preset.Name(),
		TestID:  "preset-" + c.preset.Name(),
		OnPress: func(core.PressEvent) { motion.Play(c.preset.Spec(), nil) },
		Style: core.Style{
			Width: core.Pt(96), Height: core.Pt(62), CornerRadius: 9, Gap: 5,
			Background: core.RGBA(0xffffff, 0.04), BorderWidth: 1, BorderColor: colLine,
			AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter,
		},
		Children: []core.Node{
			core.View{Style: core.Style{
				Width: core.Pt(18), Height: core.Pt(18), CornerRadius: 5, Background: colPink,
				Animated: motion.Style(),
			}},
			core.Text{Content: c.preset.Name(), Style: core.Style{FontSize: 10, Color: colFaint}},
		},
	}
}
