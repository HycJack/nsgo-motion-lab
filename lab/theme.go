package main

import (
	"time"

	"github.com/oneo-me/nsgo/base"
	"github.com/oneo-me/nsgo/core"
)

// The palette: one dark surface, and an accent per page — cyan for
// canvas, violet for animation, green for shaders.
var (
	colBg     = core.RGB(0x0a0e14)
	colBar    = core.RGB(0x111823)
	colCard   = core.RGB(0x151d29)
	colWell   = core.RGB(0x0b1119)
	colLine   = core.RGBA(0x93a7c4, 0.16)
	colText   = core.RGB(0xeaf1fa)
	colMuted  = core.RGBA(0xa2b3c9, 0.78)
	colFaint  = core.RGBA(0xa2b3c9, 0.45)
	colCyan   = core.RGB(0x35d6de)
	colViolet = core.RGB(0xa78bfa)
	colGreen  = core.RGB(0x39d98a)
	colAmber  = core.RGB(0xffb454)
	colPink   = core.RGB(0xff6b9d)
	colBlue   = core.RGB(0x4d8dff)
)

// ticker calls fn once per display frame with the seconds since it
// started, registering the next frame from inside the callback — the
// loop that Scope.RequestAnimationFrame documents. It costs the engine
// one commit a frame, so only demos that need per-frame drawing use it;
// everything else runs on animated values the platform steps itself.
func ticker(cx *core.Scope, fn func(t float64)) {
	start := time.Now()
	var loop func(at time.Time)
	loop = func(at time.Time) {
		fn(at.Sub(start).Seconds())
		cx.RequestAnimationFrame(loop)
	}
	cx.RequestAnimationFrame(loop)
}

// card is one example: a title, a line saying what to notice, and the
// demo. Cards wrap two to a row on a wide window.
func card(title, note string, accent core.Color, children ...core.Node) core.Node {
	head := []core.Node{
		core.View{Style: core.Style{Direction: core.Row, AlignItems: core.AlignCenter, Gap: 8}, Children: []core.Node{
			core.View{Style: core.Style{Width: core.Pt(7), Height: core.Pt(7), CornerRadius: 4, Background: accent}},
			core.Text{Content: title, Style: core.Style{FontSize: 14, FontWeight: core.FontWeightBold, Color: colText}},
		}},
	}
	if note != "" {
		head = append(head, core.Text{Content: note, Style: core.Style{FontSize: 12, Color: colMuted, LineHeight: 18}})
	}
	body := append(head, children...)
	return core.View{
		Style: core.Style{
			FlexGrow: 1, FlexBasis: core.Pct(46), Background: colCard, CornerRadius: 14,
			Padding: core.All(16), Gap: 10, BorderWidth: 1, BorderColor: colLine,
		},
		Children: body,
	}
}

// grid lays cards out two to a row, wrapping as the window narrows.
func grid(children ...core.Node) core.Node {
	return core.View{
		Style:    core.Style{Direction: core.Row, Wrap: core.Wrap, Gap: 14, Padding: core.All(16), AlignItems: core.AlignStart},
		Children: children,
	}
}

// chip is a small button that replays a demo.
func chip(text string, accent core.Color, onPress func()) core.Node {
	return base.Button{
		Label:   text,
		TestID:  "chip-" + text,
		OnPress: func(core.PressEvent) { onPress() },
		Style:   core.Style{Padding: core.Symmetric(6, 12), CornerRadius: 8, Background: core.RGBA(0xffffff, 0.06)},
		Children: []core.Node{core.Text{
			Content: text,
			Style:   core.Style{FontSize: 12, FontWeight: core.FontWeightMedium, Color: accent},
		}},
	}
}

// note is a line of small grey text.
func note(text string) core.Node {
	return core.Text{Content: text, Style: core.Style{FontSize: 11, Color: colFaint, LineHeight: 16}}
}

// noteFunc is note for a string that changes.
func noteFunc(fn func() string) core.Node {
	return core.Text{ContentFunc: fn, Style: core.Style{FontSize: 11, Color: colFaint, LineHeight: 16}}
}

// well is the recessed box a demo draws into.
func well(height float32) core.Style {
	return core.Style{Height: core.Pt(core.Points(height)), Background: colWell, CornerRadius: 10, Overflow: core.Hidden}
}

// row is a horizontal line of items.
func row(gap float32, children ...core.Node) core.Node {
	return core.View{
		Style:    core.Style{Direction: core.Row, AlignItems: core.AlignCenter, Gap: core.Points(gap)},
		Children: children,
	}
}

// bezier wraps a CSS curve as an animation easing.
func bezier(c core.Curve) core.Easing {
	return core.Easing{Kind: core.EasingBezier, Bezier: c}
}

// lerp maps 0..1 to a..b as an animated node, for Style.Animated.
func lerp(n core.AnimatedNode, a, b float64) core.AnimatedNode {
	return core.Interpolate(n, core.Interpolation{
		Input:  []float64{0, 1},
		Output: []float64{a, b},
	})
}
