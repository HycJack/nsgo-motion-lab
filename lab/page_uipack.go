package main

// page_uipack.go finishes the ui pack: motion-prompt #049–#072 also names the
// overlays, the picks and the navigational furniture. page_uikit.go draws
// twelve of them out in full; these are the other twelve, as tiles on the same
// shared driver.
//
// Catalog and timings: github.com/showheyohtaki/motion-prompt (MIT). The code
// is NSGo's own.

import (
	"github.com/oneo-me/nsgo/core"
)

// The ui pack draws white sheets on the dark stage, so it needs its own ink.
var (
	uipPaper = core.RGB(0xf3f6fa)
	uipInk   = core.RGB(0x1f2a37)
	uipGrey  = core.RGB(0xc9d3df)
	uipLine  = core.RGB(0xdfe5ee)
	uipBlue  = core.RGB(0x2f6bf0)
	uipPale  = core.RGB(0xd9e4fb)
)

// uipScreen is the fixed-size box a mock UI is laid out in. A fixed width
// matters: positions inside it have to be animated in points, and a
// percentage cannot be animated.
func uipScreen(w, h float32, children ...core.Node) core.Node {
	return Centered(core.View{Style: core.Style{
		Width:     core.Pt(core.Points(w)),
		Height:    core.Pt(core.Points(h)),
		Position:  core.Relative,
		Direction: core.Column,
	}, Children: children})
}

// uipBar is a line of fake text.
func uipBar(w, h float32, col core.Color, anim *core.AnimatedStyle) core.Node {
	return core.View{Style: core.Style{
		Width: core.Pt(core.Points(w)), Height: core.Pt(core.Points(h)),
		CornerRadius: core.Points(h / 2), Background: col, Animated: anim,
	}}
}

// uipSheet is the white surface of a dialog, drawer or menu.
func uipSheet(w, h float32, radius core.Points, anim *core.AnimatedStyle, children ...core.Node) core.Node {
	return core.View{Style: core.Style{
		Width: core.Pt(core.Points(w)), Height: core.Pt(core.Points(h)),
		CornerRadius: radius, Background: uipPaper, Animated: anim,
		Padding: core.Edges{Top: 9, Right: 10, Bottom: 9, Left: 10}, Gap: 5,
	}, Children: children}
}

// uipChip is a small rounded label, used for buttons in the mocks.
func uipChip(w, h float32, col core.Color, children ...core.Node) core.Node {
	return core.View{Style: core.Style{
		Width: core.Pt(core.Points(w)), Height: core.Pt(core.Points(h)),
		CornerRadius: core.Points(h / 2), Background: col,
		AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter,
	}, Children: children}
}

// uipLabel is text inside a mock, on the paper.
func uipLabel(s string, size float32, col core.Color) core.Node {
	return core.Text{Content: s, Style: core.Style{FontSize: core.Points(size), FontWeight: core.FontWeightBold, Color: col}}
}

// openClose maps one 0..1 pass to "open, hold, close" so a single replay shows
// both edges of an overlay, the way the catalog's cycle does.
func openClose(p *core.AnimatedValue, shutAt float64) core.AnimatedNode {
	return core.Interpolate(p, core.Interpolation{
		Input:  []float64{0, shutAt, 0.86, 1},
		Output: []float64{0, 1, 1, 0},
	})
}

// uiPackTiles are the twelve effects page_uikit.go does not draw in full.
func uiPackTiles(cx *core.Scope) []core.Node {
	return []core.Node{
		uiModal(cx),
		uiDrawer(cx),
		uiSheet(cx),
		uiDropdown(cx),
		uiPopover(cx),
		uiChips(cx),
		uiRange(cx),
		uiPager(cx),
		uiCrumbs(cx),
		uiMega(cx),
		uiSteps(cx),
		uiEmpty(cx),
	}
}

// 049 Modal dialog: the page dims, a window scales up into the middle, × closes it.
func uiModal(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "049", Title: "Modal dialog", Duration: 2400 * ms, Accent: colBlue, Height: 104,
		Note: "The page dims and a white window scales up into the centre; the × at its corner closes it again.",
		Stage: func(p *core.AnimatedValue) core.Node {
			a := openClose(p, 0.5)
			return uipScreen(210, 96,
				At(10, 10, uipBar(70, 9, uipLine, nil)),
				At(22, 10, uipBar(180, 9, uipLine, nil)),
				At(34, 10, uipBar(150, 9, uipLine, nil)),
				At(46, 10, uipChip(44, 14, uipBlue)),
				At(64, 110, uipBar(90, 9, uipLine, nil)),
				core.View{Style: core.Style{
					Position: core.Absolute,
					Top:      core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					Background: core.RGB(0x070d15), Animated: &core.AnimatedStyle{Opacity: lerp(a, 0, 0.66)},
				}},
				Centered(uipSheet(146, 76, 10, &core.AnimatedStyle{Opacity: a, Scale: lerp(a, 0.9, 1)},
					row(6,
						uipLabel("Delete project?", 11, uipInk),
						core.View{Style: core.Style{FlexGrow: 1}},
						uipLabel("×", 12, uipGrey),
					),
					uipBar(120, 6, uipLine, nil),
					uipBar(96, 6, uipLine, nil),
					row(6,
						core.View{Style: core.Style{FlexGrow: 1}},
						uipChip(34, 14, uipLine),
						uipChip(40, 14, uipBlue),
					),
				)),
			)
		}})
}

// 050 Drawer: a tall panel slides in from the left edge, the page dims behind it.
func uiDrawer(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "050", Title: "Drawer (off-canvas)", Duration: 2300 * ms, Accent: colViolet, Height: 104,
		Note: "A tall panel slides in from the left; the rest of the page dims, and the panel goes back the same way.",
		Stage: func(p *core.AnimatedValue) core.Node {
			a := openClose(p, 0.52)
			menu := []core.Node{uipBar(56, 7, uipGrey, nil)}
			for i := 0; i < 4; i++ {
				col := uipLine
				if i == 0 {
					col = uipBlue
				}
				menu = append(menu, row(6, uipChip(10, 10, col), uipBar(48, 6, col, nil)))
			}
			return uipScreen(210, 96,
				At(8, 8, uipChip(46, 16, uipLine)),
				At(30, 8, uipBar(96, 6, uipLine, nil)),
				At(42, 8, uipBar(76, 6, uipLine, nil)),
				At(56, 8, uipBar(160, 6, uipLine, nil)),
				At(68, 8, uipBar(120, 6, uipLine, nil)),
				core.View{Style: core.Style{
					Position: core.Absolute,
					Top:      core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					Background: core.RGB(0x070d15), Animated: &core.AnimatedStyle{Opacity: lerp(a, 0, 0.55)},
				}},
				core.View{Style: core.Style{
					Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0),
					Animated: &core.AnimatedStyle{TranslateX: lerp(a, -86, 0)},
				}, Children: []core.Node{
					uipSheet(78, 96, 0, nil, menu...),
				}},
			)
		}})
}

// 051 Bottom sheet: the panel rises from the bottom edge, grab handle first.
func uiSheet(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "051", Title: "Bottom sheet", Duration: 2200 * ms, Accent: colGreen, Height: 104,
		Note: "A sheet rises from the bottom and stops short of the top; the grab handle at its edge is what drags it back down.",
		Stage: func(p *core.AnimatedValue) core.Node {
			a := openClose(p, 0.35)
			return uipScreen(210, 96,
				At(10, 10, uipBar(80, 8, uipLine, nil)),
				At(24, 10, uipBar(46, 34, uipLine, nil)),
				At(24, 118, uipBar(46, 34, uipLine, nil)),
				At(64, 10, uipBar(150, 6, uipLine, nil)),
				At(76, 10, uipBar(110, 6, uipLine, nil)),
				At(84, 76, uipChip(58, 14, uipBlue)),
				core.View{Style: core.Style{
					Position: core.Absolute,
					Top:      core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
					Background: core.RGB(0x070d15), Animated: &core.AnimatedStyle{Opacity: lerp(a, 0, 0.5)},
				}},
				core.View{Style: core.Style{
					Position: core.Absolute, Left: core.Pt(0), Bottom: core.Pt(0),
					Animated: &core.AnimatedStyle{TranslateY: lerp(a, 58, 0)},
				}, Children: []core.Node{
					uipSheet(210, 58, 12, nil,
						centredBar(26, 4, uipGrey),
						uipLabel("Details", 10, uipInk),
						row(6, uipChip(12, 12, uipBlue), uipBar(110, 6, uipLine, nil)),
						row(6, uipChip(12, 12, uipPale), uipBar(84, 6, uipLine, nil)),
					),
				}},
			)
		}})
}

// centredBar is a small handle, centred in its row.
func centredBar(w, h float32, col core.Color) core.Node {
	return core.View{Style: core.Style{Direction: core.Row, JustifyContent: core.JustifyCenter},
		Children: []core.Node{uipBar(w, h, col, nil)}}
}

// 053 Dropdown: the list opens under the field, the chip turns over, the choice fills in.
func uiDropdown(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "053", Title: "Dropdown (select)", Duration: 2400 * ms, Accent: colCyan, Height: 110,
		Note: "The list drops the last six points into place while it fades in, the chevron turns over, and the chosen city replaces the placeholder.",
		Stage: func(p *core.AnimatedValue) core.Node {
			open := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.18, 0.72, 0.86}, Output: []float64{0, 1, 1, 0}})
			picked := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.5, 0.7, 1}, Output: []float64{0, 0, 1, 1}})
			placeholder := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.44, 0.5, 1}, Output: []float64{1, 1, 0, 0}})
			options := []core.Node{
				uipBar(52, 7, uipInk, nil),
				uipBar(44, 7, uipInk, nil),
				uipBar(48, 7, uipBlue, nil),
				uipBar(40, 7, uipInk, nil),
			}
			var list []core.Node
			for i, o := range options {
				bg := core.Color(0)
				if i == 2 {
					bg = uipPale
				}
				list = append(list, core.View{Style: core.Style{
					Height: core.Pt(14), CornerRadius: 4, Background: bg,
					Padding: core.Edges{Left: 6}, Direction: core.Row, AlignItems: core.AlignCenter,
				}, Children: []core.Node{o}})
			}
			return uipScreen(210, 102,
				At(4, 4, uipBar(60, 6, uipLine, nil)),
				core.View{Style: core.Style{
					Position: core.Absolute, Top: core.Pt(16), Left: core.Pt(4),
					Width: core.Pt(120), Height: core.Pt(20), CornerRadius: 5,
					BorderWidth: 2, Background: core.RGB(0x0e1622),
					Animated: &core.AnimatedStyle{BorderColor: core.Interpolate(open, core.Interpolation{
						Input:  []float64{0, 1},
						Output: []float64{0, 1},
						Colors: []core.Color{uipLine, uipBlue},
					})},
					Direction: core.Row, AlignItems: core.AlignCenter, Padding: core.Edges{Left: 6}, Gap: 4,
				}, Children: []core.Node{
					core.Text{Content: "Choose a city", Style: core.Style{Position: core.Absolute, FontSize: 8, Color: uipGrey, Animated: &core.AnimatedStyle{Opacity: placeholder}}},
					core.Text{Content: "Nagoya", Style: core.Style{Position: core.Absolute, FontSize: 8, FontWeight: core.FontWeightBold, Color: uipInk, Animated: &core.AnimatedStyle{Opacity: picked}}},
					core.View{Style: core.Style{FlexGrow: 1}},
					core.Text{Content: "v", Style: core.Style{FontSize: 8, Color: uipGrey, Animated: &core.AnimatedStyle{Rotate: lerp(open, 0, 180)}}},
					core.View{Style: core.Style{Width: core.Pt(4)}},
				}},
				core.View{Style: core.Style{
					Position: core.Absolute, Top: core.Pt(40), Left: core.Pt(4),
					Width: core.Pt(120), CornerRadius: 5, Background: uipPaper, Padding: core.Edges{Top: 4, Right: 4, Bottom: 4, Left: 4}, Gap: 2,
					Animated: &core.AnimatedStyle{Opacity: open, TranslateY: lerp(open, -6, 0)},
				}, Children: []core.Node{core.View{Style: core.Style{Gap: 2}, Children: list}}},
				At(16, 136, uipBar(66, 34, uipLine, nil)),
				At(56, 136, uipBar(60, 6, uipLine, nil)),
				At(66, 136, uipBar(44, 6, uipLine, nil)),
			)
		}})
}

// 055 Popover: a small panel opens beside the button that asked for it.
func uiPopover(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "055", Title: "Popover", Duration: 2200 * ms, Accent: colAmber, Height: 104,
		Note: "The panel belongs to the button, so it slides the last six points in from that side and points back at it.",
		Stage: func(p *core.AnimatedValue) core.Node {
			a := openClose(p, 0.5)
			return uipScreen(210, 96,
				At(6, 6, uipBar(70, 7, uipLine, nil)),
				At(18, 6, uipBar(120, 7, uipLine, nil)),
				core.View{Style: core.Style{
					Position: core.Absolute, Top: core.Pt(56), Left: core.Pt(6),
					Width: core.Pt(62), Height: core.Pt(18), CornerRadius: 5,
					BorderWidth: 2, BorderColor: uipBlue, Background: core.RGB(0x0e1622),
					AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter,
				}, Children: []core.Node{uipLabel("Details", 9, uipBlue)}},
				At(70, 6, uipBar(150, 7, uipLine, nil)),
				core.View{Style: core.Style{
					Position: core.Absolute, Top: core.Pt(30), Left: core.Pt(74),
					Animated: &core.AnimatedStyle{Opacity: a, TranslateX: lerp(a, 6, 0)},
				}, Children: []core.Node{
					row(0,
						core.View{Style: core.Style{
							Width: core.Pt(8), AlignItems: core.AlignCenter,
						}, Children: []core.Node{
							core.View{Style: core.Style{Width: core.Pt(8), Height: core.Pt(8), Background: uipPaper,
								Animated: &core.AnimatedStyle{Rotate: lerp(a, 45, 45)}}},
						}},
						uipSheet(118, 54, 8, nil,
							uipLabel("Autosave", 9, uipInk),
							uipBar(96, 5, uipLine, nil),
							uipBar(74, 5, uipLine, nil),
							row(6,
								core.View{Style: core.Style{FlexGrow: 1}},
								uipChip(26, 12, uipLine),
								uipChip(28, 12, uipBlue),
							),
						),
					),
				}},
			)
		}})
}

// 061 Filter chips: the chips a list is filtered by widen for a check mark, and
// the rows that no longer match drop to a twelfth of their opacity.
func uiChips(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "061", Title: "Filter chips", Duration: 2100 * ms, Accent: colPink, Height: 104,
		Note: "Choosing a chip widens it to make room for a check; rows that do not match stay on screen at 12%.",
		Stage: func(p *core.AnimatedValue) core.Node {
			newOn := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.2, 0.8, 1}, Output: []float64{0, 1, 1, 0}})
			freeOn := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.34, 0.72, 0.86}, Output: []float64{0, 1, 1, 0}})
			var chips []core.Node
			for i, s := range []string{"All", "New", "Popular", "Free"} {
				var on core.AnimatedNode = newOn
				if i == 3 {
					on = freeOn
				} else if i < 2 {
					on = nil
				}
				bg := uipLine
				fg := uipInk
				if on != nil {
					bg = uipPale
					fg = uipBlue
				}
				body := []core.Node{
					core.Text{Content: s, Style: core.Style{FontSize: 8, FontWeight: core.FontWeightBold, Color: fg}},
				}
				if on != nil {
					body = append([]core.Node{core.Text{Content: "v", Style: core.Style{FontSize: 8, FontWeight: core.FontWeightBold, Color: uipBlue, Animated: &core.AnimatedStyle{Opacity: on}}}}, body...)
				}
				chips = append(chips, core.View{Style: core.Style{
					Height: core.Pt(16), CornerRadius: 8, Background: bg,
					Padding: core.Edges{Left: 8, Right: 8}, Direction: core.Row, AlignItems: core.AlignCenter, Gap: 4,
				}, Children: body})
			}
			dimNew := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.2, 0.8, 1}, Output: []float64{1, 0.12, 0.12, 1}})
			return uipScreen(210, 96,
				At(6, 6, core.View{Style: core.Style{Direction: core.Row, Gap: 5}, Children: chips}),
				At(32, 6, core.View{Style: core.Style{Direction: core.Row, Gap: 6}, Children: []core.Node{
					core.View{Style: core.Style{Gap: 4, Animated: &core.AnimatedStyle{Opacity: dimNew}}, Children: []core.Node{
						uipChip(96, 24, core.RGB(0x131c28)),
						uipChip(96, 24, core.RGB(0x131c28)),
					}},
					core.View{Style: core.Style{Gap: 4}, Children: []core.Node{
						uipChip(96, 24, core.RGB(0x131c28)),
						uipChip(96, 24, core.RGB(0x131c28)),
					}},
				}}),
			)
		}})
}

// 062 Range slider: two handles narrow the range; the blue part is the range itself.
func uiRange(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "062", Title: "Range slider", Duration: 2400 * ms, Accent: colCyan, Height: 104,
		Note: "The two handles are dragged inward and the blue part between them is the range; the bubbles carry the numbers.",
		Stage: func(p *core.AnimatedValue) core.Node {
			const track = 180.0
			// lo and hi are the two handle positions as fractions, so the
			// blue segment's scale and offset follow from the same breakpoints.
			lo := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.5, 1}, Output: []float64{0.1 * track, 0.32 * track, 0.1 * track}})
			hi := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.45, 0.9, 1}, Output: []float64{1 * track, 0.8 * track, 0.8 * track, 1 * track}})
			span := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.45, 0.5, 0.9, 1}, Output: []float64{0.9, 0.48, 0.48, 0.48, 0.9}})
			shift := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.45, 0.5, 0.9, 1}, Output: []float64{9, 9.4, 9.4, 9.4, 9}})
			return uipScreen(200, 96,
				At(10, 10, uipBar(30, 7, uipLine, nil)),
				At(10, 128, uipBar(52, 7, uipBlue, nil)),
				At(54, 10, uipBar(track, 5, core.RGB(0x1b2534), nil)),
				At(54, 10, core.View{Style: core.Style{
					Width: core.Pt(track), Height: core.Pt(5), CornerRadius: 3, Background: uipBlue,
					Animated: &core.AnimatedStyle{ScaleX: span, TranslateX: shift},
				}}),
				At(47, 0, core.View{Style: core.Style{
					Position: core.Absolute, Left: core.Pt(10), Animated: &core.AnimatedStyle{TranslateX: lo},
				}, Children: []core.Node{rangeHandle()}}),
				At(47, 0, core.View{Style: core.Style{
					Position: core.Absolute, Left: core.Pt(10), Animated: &core.AnimatedStyle{TranslateX: hi},
				}, Children: []core.Node{rangeHandle()}}),
				At(78, 10, uipBar(14, 5, uipGrey, nil)),
				At(78, 176, uipBar(18, 5, uipGrey, nil)),
			)
		}})
}

// rangeHandle is a white knob with a blue rim.
func rangeHandle() core.Node {
	return core.View{Style: core.Style{
		Width: core.Pt(12), Height: core.Pt(12), CornerRadius: 6,
		Background: uipPaper, BorderWidth: 2, BorderColor: uipBlue,
	}}
}

// 064 Pagination: the marker circle travels to the page that was pressed, and
// the list above dips out and back with the new page in it.
func uiPager(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "064", Title: "Pagination", Duration: 2100 * ms, Accent: colViolet, Height: 104,
		Note: "One marker circle travels between the page numbers; the list above blinks out for a moment as its contents change.",
		Stage: func(p *core.AnimatedValue) core.Node {
			marker := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.33, 0.66, 1}, Output: []float64{0, 26, 52, 0}})
			dip := core.Interpolate(p, core.Interpolation{
				Input:  []float64{0, 0.3, 0.33, 0.63, 0.66, 0.96, 1},
				Output: []float64{1, 1, 0.15, 1, 0.15, 1, 1},
			})
			var rows []core.Node
			for i := 0; i < 3; i++ {
				rows = append(rows, row(6, uipChip(20, 16, uipLine), core.View{Style: core.Style{Gap: 3}, Children: []core.Node{
					uipBar(100, 7, uipLine, nil), uipBar(70, 5, uipGrey, nil),
				}}))
			}
			var nums []core.Node
			for i := 0; i < 5; i++ {
				nums = append(nums, core.View{Style: core.Style{
					Width: core.Pt(18), Height: core.Pt(18), AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter,
				}, Children: []core.Node{uipLabel(string(rune('1'+i)), 9, uipGrey)}})
			}
			return uipScreen(210, 96,
				At(6, 6, core.View{Style: core.Style{Gap: 5, Animated: &core.AnimatedStyle{Opacity: dip}}, Children: rows}),
				At(70, 4, core.View{Style: core.Style{Direction: core.Row, Gap: 8, Position: core.Relative}, Children: []core.Node{
					core.View{Style: core.Style{
						Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0),
						Width: core.Pt(18), Height: core.Pt(18), CornerRadius: 9, Background: uipBlue,
						Animated: &core.AnimatedStyle{TranslateX: marker},
					}},
					core.View{Style: core.Style{Direction: core.Row, Gap: 8}, Children: nums},
				}}),
			)
		}})
}

// 066 Breadcrumbs: going up a level rewrites the trail and cross-fades the page.
func uiCrumbs(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "066", Title: "Breadcrumbs", Duration: 2200 * ms, Accent: colGreen, Height: 104,
		Note: "Up a level: the last crumb drops out, the one before it becomes the current page, and the content cross-fades.",
		Stage: func(p *core.AnimatedValue) core.Node {
			up := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.35, 0.72, 0.95}, Output: []float64{0, 1, 1, 0}})
			detail := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.35, 0.72, 0.95}, Output: []float64{1, 0, 0, 1}})
			return uipScreen(210, 96,
				At(4, 6, row(4, uipLabel("Home", 8, uipBlue), uipLabel(">", 8, uipGrey),
					core.Text{Content: "Courses", Style: core.Style{FontSize: 8, FontWeight: core.FontWeightBold, Color: uipBlue,
						Animated: &core.AnimatedStyle{TextColor: core.Interpolate(up, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1}, Colors: []core.Color{uipBlue, uipInk}})}}},
					core.View{Style: core.Style{Animated: &core.AnimatedStyle{Opacity: core.Interpolate(up, core.Interpolation{Input: []float64{0, 1}, Output: []float64{1, 0}})}}, Children: []core.Node{
						row(4, uipLabel(">", 8, uipGrey), uipLabel("Web design", 8, uipInk)),
					}},
				)),
				At(20, 4, uipBar(200, 1, uipLine, nil)),
				core.View{Style: core.Style{
					Position: core.Absolute, Top: core.Pt(28), Left: core.Pt(6),
					Animated: &core.AnimatedStyle{Opacity: detail},
				}, Children: []core.Node{core.View{Style: core.Style{Gap: 5},
					Children: []core.Node{
						uipBar(90, 9, uipInk, nil),
						uipBar(180, 6, uipLine, nil),
						uipBar(150, 6, uipLine, nil),
						uipBar(120, 26, uipLine, nil),
					}}}},
				core.View{Style: core.Style{
					Position: core.Absolute, Top: core.Pt(28), Left: core.Pt(6),
					Animated: &core.AnimatedStyle{Opacity: up},
				}, Children: []core.Node{row(6,
					core.View{Style: core.Style{Gap: 3, Width: core.Pt(62)}, Children: []core.Node{uipBar(62, 40, uipLine, nil), uipBar(40, 6, uipInk, nil), uipBar(50, 5, uipGrey, nil)}},
					core.View{Style: core.Style{Gap: 3, Width: core.Pt(62)}, Children: []core.Node{uipBar(62, 40, uipLine, nil), uipBar(40, 6, uipInk, nil), uipBar(50, 5, uipGrey, nil)}},
					core.View{Style: core.Style{Gap: 3, Width: core.Pt(62)}, Children: []core.Node{uipBar(62, 40, uipLine, nil), uipBar(40, 6, uipInk, nil), uipBar(50, 5, uipGrey, nil)}},
				)}},
			)
		}})
}

// 067 Mega menu: hovering the header drops a full-width panel of columns under it.
func uiMega(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "067", Title: "Mega menu", Duration: 2200 * ms, Accent: colBlue, Height: 110,
		Note: "A hover, not a click: the full-width panel drops twelve points into place and the hovered link underlines.",
		Stage: func(p *core.AnimatedValue) core.Node {
			a := openClose(p, 0.55)
			var cols []core.Node
			for c := 0; c < 3; c++ {
				links := []core.Node{uipBar(34, 6, uipBlue, nil)}
				for k := 0; k < 3; k++ {
					col := uipGrey
					if c == 1 && k == 1 {
						col = uipBlue
					}
					links = append(links, uipBar(float32(46+k*6), 5, col, nil))
				}
				cols = append(cols, core.View{Style: core.Style{Gap: 4, Width: core.Pt(58)}, Children: links})
			}
			cols = append(cols, core.View{Style: core.Style{Width: core.Pt(52), Height: core.Pt(42), CornerRadius: 4, Background: core.RGB(0x24303f)}})
			return uipScreen(210, 102,
				At(4, 6, core.View{Style: core.Style{Direction: core.Row, Gap: 14}, Children: []core.Node{
					uipLabel("Products", 8, uipBlue),
					uipLabel("Pricing", 8, uipGrey),
					uipLabel("Support", 8, uipGrey),
				}}),
				At(22, 4, uipBar(202, 1, uipLine, nil)),
				At(30, 6, uipBar(80, 7, uipLine, nil)),
				At(44, 6, uipBar(140, 26, uipLine, nil)),
				core.View{Style: core.Style{
					Position: core.Absolute, Top: core.Pt(26), Left: core.Pt(0),
					Animated: &core.AnimatedStyle{Opacity: a, TranslateY: lerp(a, -12, 0)},
				}, Children: []core.Node{
					core.View{Style: core.Style{
						Width: core.Pt(210), Background: uipPaper, CornerRadius: 6,
						Padding:   core.Edges{Top: 8, Right: 10, Bottom: 8, Left: 10},
						Direction: core.Row, Gap: 8, AlignItems: core.AlignStart,
					}, Children: cols},
				}},
			)
		}})
}

// 068 Step indicator: each press turns the current circle into a check and
// stretches the blue line to the next one.
func uiSteps(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "068", Title: "Step indicator", Duration: 2200 * ms, Accent: colGreen, Height: 104,
		Note: "Three presses: the circles behind you become checks and the line to the next step grows; done is a full row of checks.",
		Stage: func(p *core.AnimatedValue) core.Node {
			const step = 48.0
			grow := core.Interpolate(p, core.Interpolation{
				Input:  []float64{0, 0.3, 0.34, 0.6, 0.64, 0.9, 0.94, 1},
				Output: []float64{0, 1.0 / 3, 1.0 / 3, 2.0 / 3, 2.0 / 3, 1, 1, 1},
			})
			var circles []core.Node
			for i := 0; i < 4; i++ {
				at := float64(i) / 3
				done := core.Interpolate(p, core.Interpolation{
					Input:  []float64{0, at * 0.9, at*0.9 + 0.06, 1},
					Output: []float64{0, 0, 1, 1},
				})
				mark := core.Text{Content: string(rune('1' + i)), Style: core.Style{
					Position: core.Absolute, FontSize: 9, FontWeight: core.FontWeightBold, Color: uipInk,
				}}
				check := core.Text{Content: "v", Style: core.Style{
					Position: core.Absolute, FontSize: 9, FontWeight: core.FontWeightBold, Color: uipBlue,
					Animated: &core.AnimatedStyle{Opacity: done},
				}}
				circles = append(circles, core.View{Style: core.Style{
					Width: core.Pt(18), Height: core.Pt(18), CornerRadius: 9,
					Background: core.RGB(0x0e1622),
					Animated: &core.AnimatedStyle{Background: core.Interpolate(done, core.Interpolation{
						Input: []float64{0, 1}, Output: []float64{0, 1}, Colors: []core.Color{core.RGB(0x0e1622), uipPale},
					})},
					BorderWidth: 2, BorderColor: uipBlue,
					AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter,
				}, Children: []core.Node{mark, check}})
			}
			return uipScreen(210, 96,
				At(22, 8, core.View{Style: core.Style{Position: core.Relative, Width: core.Pt(3 * step), Height: core.Pt(18)}, Children: []core.Node{
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(8), Left: core.Pt(9), Width: core.Pt(3 * step), Height: core.Pt(3), CornerRadius: 2, Background: core.RGB(0x1b2534)}},
					core.View{Style: core.Style{
						Position: core.Absolute, Top: core.Pt(8), Left: core.Pt(9),
						Width: core.Pt(3 * step), Height: core.Pt(3), CornerRadius: 2, Background: uipBlue,
						Animated: &core.AnimatedStyle{ScaleX: grow, TranslateX: lerp(grow, 0, -1.5*step)},
					}},
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Direction: core.Row, Gap: core.Points(step - 18)}, Children: circles},
				}}),
				At(50, 6, uipBar(120, 7, uipLine, nil)),
				At(64, 6, uipBar(170, 6, uipLine, nil)),
				At(76, 150, uipChip(46, 16, uipBlue)),
			)
		}})
}

// 072 Empty state: the rows are deleted to the left, and what is left is an
// empty tray, one line of copy and the way to add something again.
func uiEmpty(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "072", Title: "Empty state", Duration: 2300 * ms, Accent: colAmber, Height: 110,
		Note: "Both rows leave to the left, the second slot closes up, and the tray rises eight points into the space they left.",
		Stage: func(p *core.AnimatedValue) core.Node {
			first := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.18, 0.24, 0.86, 1}, Output: []float64{1, 1, 0, 0, 1}})
			second := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.36, 0.42, 0.86, 1}, Output: []float64{1, 1, 0, 0, 1}})
			shift := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.18, 0.24, 0.86, 1}, Output: []float64{0, 0, -40, -40, 0}})
			shift2 := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.36, 0.42, 0.86, 1}, Output: []float64{0, 0, -40, -40, 0}})
			lift := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.36, 0.42, 0.86, 1}, Output: []float64{0, 0, -22, -22, 0}})
			empty := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.45, 0.56, 0.9, 1}, Output: []float64{0, 0, 1, 1, 0}})
			return uipScreen(210, 102,
				At(6, 6, uipBar(50, 7, uipLine, nil)),
				core.View{Style: core.Style{
					Position: core.Absolute, Top: core.Pt(20), Left: core.Pt(6),
					Animated: &core.AnimatedStyle{Opacity: first, TranslateX: shift},
				}, Children: []core.Node{
					core.View{Style: core.Style{
						Width: core.Pt(198), Height: core.Pt(24), CornerRadius: 4,
						BorderWidth: 1, BorderColor: uipLine, Padding: core.Edges{Left: 6},
						Direction: core.Row, AlignItems: core.AlignCenter, Gap: 6,
					}, Children: []core.Node{uipChip(12, 12, uipLine), uipBar(90, 6, uipLine, nil), core.View{Style: core.Style{FlexGrow: 1}}, uipChip(8, 10, uipGrey)}},
				}},
				core.View{Style: core.Style{
					Position: core.Absolute, Top: core.Pt(50), Left: core.Pt(6),
					Animated: &core.AnimatedStyle{Opacity: second, TranslateX: shift2, TranslateY: lift},
				}, Children: []core.Node{
					core.View{Style: core.Style{
						Width: core.Pt(198), Height: core.Pt(24), CornerRadius: 4,
						BorderWidth: 1, BorderColor: uipLine, Padding: core.Edges{Left: 6},
						Direction: core.Row, AlignItems: core.AlignCenter, Gap: 6,
					}, Children: []core.Node{uipChip(12, 12, uipLine), uipBar(70, 6, uipLine, nil), core.View{Style: core.Style{FlexGrow: 1}}, uipChip(8, 10, uipGrey)}},
				}},
				core.View{Style: core.Style{
					Position: core.Absolute, Top: core.Pt(58), Left: core.Pt(0), Width: core.Pt(210), AlignItems: core.AlignCenter,
					Animated: &core.AnimatedStyle{Opacity: empty, TranslateY: lerp(empty, 8, 0)},
				}, Children: []core.Node{
					core.View{Style: core.Style{Gap: 4, AlignItems: core.AlignCenter}, Children: []core.Node{
						core.View{Style: core.Style{Width: core.Pt(26), Height: core.Pt(16), CornerRadius: 4, BorderWidth: 2, BorderColor: uipBlue}},
						uipLabel("Nothing here yet", 9, uipGrey),
						uipChip(44, 14, uipBlue),
					}},
				}},
			)
		}})
}
