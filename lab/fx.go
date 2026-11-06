package main

// fx.go is the driver every named effect in the eight pack pages shares.
//
// motion-prompt names 192 effects; the pages built on Fx (motion, camera,
// ui, lp, photo, motion2, lp2, camera2) show all of them the same way: one
// animated 0..1 value, a stage that draws that value, the name from the
// catalog, and a chip that replays it. A tile is therefore a few lines of
// composition rather than a hundred lines of plumbing.
//
// The catalog and its timings come from github.com/showheyohtaki/motion-prompt
// (MIT). Only names, categories and timings are reused — every effect here
// is written from scratch with NSGo's own primitives.

import (
	"time"

	"github.com/oneo-me/nsgo/core"
)

// ms is the unit the catalog's timings are written in.
const ms = time.Millisecond

// stageFill stretches a stage's wrapper to the whole stage.
//
// A stage is a box of a fixed height, and the shapes inside it are absolute,
// positioned from its top-left corner. Those shapes need a box to resolve
// against: a wrapper with no style is only as tall as its content, and
// absolute children never count towards their parent's height, so an unstyled
// wrapper measures zero and everything inside it is drawn into nothing. Every
// stage therefore returns either one absolute element or a wrapper carrying
// stageFill.
var stageFill = core.Style{
	Position: core.Absolute,
	Top:      core.Pt(0), Left: core.Pt(0),
	Right: core.Pt(0), Bottom: core.Pt(0),
}

// FxSpec describes one catalog effect.
type FxSpec struct {
	// ID is the three-digit catalog number, e.g. "001".
	ID string
	// Title is the catalog title, English half only, e.g. "Pop-in".
	Title string
	// Note is one line saying what to watch.
	Note string
	// Duration is one pass of the effect; 900ms when zero.
	Duration time.Duration
	// Loop replays the pass forever.
	Loop bool
	// PingPong runs 0->1 then 1->0, forever. It wins over Loop.
	PingPong bool
	// Linear drives the pass at a constant rate instead of easing in and out.
	// A take wants it: its progress has to be time, so the windows written in
	// it can be read off a clock.
	Linear bool
	// Height is the stage's height in points; 96 when zero.
	Height float32
	// Width pins the stage to a fixed width and centres it, for effects
	// composed as a page or a screen rather than as a loose shape. Zero
	// lets the stage fill the card.
	Width float32
	// Accent tints the title dot and the replay chip; cyan when zero.
	Accent core.Color
	// Stage draws the effect. It is already inside a Relative box that
	// fills the stage, so absolute children are positioned from its
	// top-left corner.
	Stage func(p *core.AnimatedValue) core.Node
}

// Fx builds one effect tile: the stage playing it, its catalog name, the
// line saying what to watch, and a Replay chip.
func Fx(cx *core.Scope, s FxSpec) core.Node {
	duration := s.Duration
	if duration == 0 {
		duration = 900 * ms
	}
	height := s.Height
	if height == 0 {
		height = 96
	}
	accent := s.Accent
	if accent == 0 {
		accent = colCyan
	}

	p := core.NewAnimatedValue(cx, 0)
	play := func() {
		p.Set(0)
		cfg := core.TimingConfig{Duration: duration, Easing: bezier(core.EaseInOut)}
		if s.Linear {
			// The zero curve is the identity, so this is time itself.
			cfg.Easing = core.Easing{Bezier: core.Linear}
		}
		var anim *core.Animation
		switch {
		case s.PingPong:
			anim = core.Loop(core.Sequence(
				core.Timing(p, 1, cfg),
				core.Timing(p, 0, cfg),
			), 0)
		case s.Loop:
			anim = core.Loop(core.Timing(p, 1, cfg), 0)
		default:
			anim = core.Timing(p, 1, cfg)
		}
		cx.StartAnimation(anim, nil)
	}
	// Starting the pass here is what makes a tile autoplay when the page
	// mounts; the chip does the same thing again on demand.
	play()

	stage := well(height)
	if s.Width > 0 {
		stage.Width = core.Pt(core.Points(s.Width))
		stage.AlignSelf = core.AlignCenter
	}

	return card(s.ID+" · "+s.Title, s.Note, accent,
		core.View{Style: stage, Children: []core.Node{
			core.View{Style: core.Style{Position: core.Relative, FlexGrow: 1}, Children: []core.Node{s.Stage(p)}},
		}},
		core.View{Style: core.Style{Direction: core.Row, Gap: 8}, Children: []core.Node{
			chip("Replay "+s.ID, accent, play),
		}},
	)
}

// Disc is a round block.
func Disc(size float32, col core.Color, anim *core.AnimatedStyle) core.Node {
	return core.View{Style: core.Style{
		Width: core.Pt(core.Points(size)), Height: core.Pt(core.Points(size)),
		CornerRadius: core.Points(size / 2), Background: col, Animated: anim,
	}}
}

// Panel is a rectangle with soft corners.
func Panel(w, h float32, col core.Color, anim *core.AnimatedStyle) core.Node {
	return core.View{Style: core.Style{
		Width: core.Pt(core.Points(w)), Height: core.Pt(core.Points(h)),
		CornerRadius: 6, Background: col, Animated: anim,
	}}
}

// Ring is an outline; give it the whole border and let a rotation make
// the gap read as an arc if you want one.
func Ring(size float32, thick core.Points, col core.Color, anim *core.AnimatedStyle) core.Node {
	return core.View{Style: core.Style{
		Width: core.Pt(core.Points(size)), Height: core.Pt(core.Points(size)),
		CornerRadius: core.Points(size / 2),
		BorderWidth:  thick, BorderColor: col, Animated: anim,
	}}
}

// Word is a text subject for the typography effects.
func Word(text string, size float32, col core.Color, anim *core.AnimatedStyle) core.Node {
	return core.Text{Content: text, Style: core.Style{
		FontSize: core.Points(size), FontWeight: core.FontWeightBold, Color: col, Animated: anim,
	}}
}

// At pins a fixed-size child at top/left inside a Relative parent. It is
// for shapes that carry their own width: an absolute box given only a
// Left edge stretches to the parent's right edge instead of hugging its
// content, so a text child would wrap.
func At(top, left float32, child core.Node) core.Node {
	return core.View{Style: core.Style{
		Position: core.Absolute, Top: core.Pt(core.Points(top)), Left: core.Pt(core.Points(left)),
	}, Children: []core.Node{child}}
}

// Fill covers the whole stage with a solid colour, for skies and scrims.
func Fill(col core.Color) core.Node {
	return core.View{Style: core.Style{
		Position: core.Absolute,
		Top:      core.Pt(core.Points(0)), Left: core.Pt(core.Points(0)),
		Right: core.Pt(core.Points(0)), Bottom: core.Pt(core.Points(0)),
		Background: col,
	}}
}

// Centered puts a child in the middle of the stage.
func Centered(child core.Node) core.Node {
	return core.View{Style: core.Style{
		Position: core.Absolute,
		Top:      core.Pt(core.Points(0)), Left: core.Pt(core.Points(0)),
		Right: core.Pt(core.Points(0)), Bottom: core.Pt(core.Points(0)),
		AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter,
	}, Children: []core.Node{child}}
}

// StageRow is a horizontal line of subjects centred in the stage.
func StageRow(gap float32, children ...core.Node) core.Node {
	return core.View{Style: core.Style{
		Position: core.Absolute,
		Top:      core.Pt(core.Points(0)), Left: core.Pt(core.Points(0)),
		Right: core.Pt(core.Points(0)), Bottom: core.Pt(core.Points(0)),
		Direction: core.Row, AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter,
		Gap: core.Points(gap),
	}, Children: children}
}
