package main

// page_uipack.go finishes the ui pack: motion-prompt #049–#072 also names the
// overlays, the picks and the navigational furniture. page_uikit.go draws
// twelve of them out in full; these are the other twelve, as tiles on the same
// shared driver.
//
// The mocks are laid out on a wide canvas (480 wide) so they fill the card's
// stage instead of floating as a small portrait thumbnail in the middle of a
// much wider well.
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
	uipCard  = core.RGB(0x1b2534)
)

// uipW is the mock canvas width. The stage is ~533 wide; a canvas a little
// narrower reads as a framed window rather than a full-bleed page.
const uipW = 480.0

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

// uipBox is a rounded block (a card or a panel of a mock).
func uipBox(w, h float32, radius core.Points, col core.Color, anim *core.AnimatedStyle, children ...core.Node) core.Node {
	return core.View{Style: core.Style{
		Width: core.Pt(core.Points(w)), Height: core.Pt(core.Points(h)),
		CornerRadius: radius, Background: col, Animated: anim,
	}, Children: children}
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

// uipPage lays out a wide, generic app page: a top nav, a heading and a row of
// cards. It is the backdrop an overlay dims and drops on top of.
func uipPage(w, h float32, accent core.Color) []core.Node {
	cw := (w - 24 - 24) / 3
	return []core.Node{
		At(8, 10, uipChip(30, 11, accent)),
		At(8, 54, row(11,
			uipBar(22, 6, uipLine, nil), uipBar(18, 6, uipLine, nil),
			uipBar(26, 6, uipLine, nil), uipBar(16, 6, uipLine, nil))),
		At(7, w-44, uipChip(34, 12, uipLine)),
		At(24, 10, uipBar(w-20, 1, uipLine, nil)),
		At(36, 12, uipBar(110, 9, uipLine, nil)),
		At(50, 12, uipBar(w-24, 6, uipLine, nil)),
		At(62, 12, uipBox(cw, h-68, 5, uipLine, nil)),
		At(62, 12+cw+12, uipBox(cw, h-68, 5, uipLine, nil)),
		At(62, 12+2*(cw+12), uipBox(cw, h-68, 5, uipLine, nil)),
	}
}

// uipDim is the full-canvas scrim an overlay drops behind its sheet.
func uipDim(a core.AnimatedNode, max float64) core.Node {
	return core.View{Style: core.Style{
		Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
		Background: core.RGB(0x070d15), Animated: &core.AnimatedStyle{Opacity: lerp(a, 0, max)},
	}}
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
			nodes := uipPage(uipW, 96, colBlue)
			nodes = append(nodes,
				uipDim(a, 0.66),
				Centered(uipSheet(250, 86, 10, &core.AnimatedStyle{Opacity: a, Scale: lerp(a, 0.92, 1)},
					row(6,
						uipLabel("Delete project?", 12, uipInk),
						core.View{Style: core.Style{FlexGrow: 1}},
						uipLabel("×", 13, uipGrey),
					),
					uipBar(200, 7, uipLine, nil),
					uipBar(150, 7, uipLine, nil),
					row(8,
						core.View{Style: core.Style{FlexGrow: 1}},
						uipChip(64, 20, uipLine),
						uipChip(72, 20, uipBlue),
					),
				)),
			)
			return uipScreen(uipW, 96, nodes...)
		}})
}

// 050 Drawer: a tall panel slides in from the left edge, the page dims behind it.
func uiDrawer(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "050", Title: "Drawer (off-canvas)", Duration: 2300 * ms, Accent: colViolet, Height: 104,
		Note: "A tall panel slides in from the left; the rest of the page dims, and the panel goes back the same way.",
		Stage: func(p *core.AnimatedValue) core.Node {
			a := openClose(p, 0.52)
			menu := []core.Node{uipBar(56, 8, uipLine, nil)}
			for i := 0; i < 4; i++ {
				col := uipLine
				if i == 0 {
					col = uipBlue
				}
				menu = append(menu, row(6, uipChip(11, 11, col), uipBar(110, 6, col, nil)))
			}
			nodes := uipPage(uipW, 96, colViolet)
			nodes = append(nodes,
				uipDim(a, 0.55),
				core.View{Style: core.Style{
					Position: core.Absolute, Top: core.Pt(0), Bottom: core.Pt(0), Left: core.Pt(0),
					Animated: &core.AnimatedStyle{TranslateX: lerp(a, -200, 0)},
				}, Children: []core.Node{uipSheet(200, 96, 0, nil, menu...)}},
			)
			return uipScreen(uipW, 96, nodes...)
		}})
}

// 051 Bottom sheet: the panel rises from the bottom edge, grab handle first.
func uiSheet(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "051", Title: "Bottom sheet", Duration: 2200 * ms, Accent: colGreen, Height: 104,
		Note: "A sheet rises from the bottom and stops short of the top; the grab handle at its edge is what drags it back down.",
		Stage: func(p *core.AnimatedValue) core.Node {
			a := openClose(p, 0.35)
			nodes := uipPage(uipW, 96, colGreen)
			nodes = append(nodes,
				uipDim(a, 0.5),
				core.View{Style: core.Style{
					Position: core.Absolute, Left: core.Pt(0), Bottom: core.Pt(0),
					Animated: &core.AnimatedStyle{TranslateY: lerp(a, 72, 0)},
				}, Children: []core.Node{
					uipSheet(uipW, 72, 12, nil,
						centredBar(36, 4, uipGrey),
						uipLabel("Details", 11, uipInk),
						row(8, uipChip(14, 14, uipBlue), uipBar(300, 7, uipLine, nil)),
						row(8, uipChip(14, 14, uipPale), uipBar(220, 7, uipLine, nil)),
					),
				}},
			)
			return uipScreen(uipW, 96, nodes...)
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
					Height: core.Pt(11), CornerRadius: 3, Background: bg,
					Padding: core.Edges{Left: 6}, Direction: core.Row, AlignItems: core.AlignCenter,
				}, Children: []core.Node{o}})
			}
			return uipScreen(uipW, 102,
				At(6, 12, uipBar(70, 6, uipLine, nil)),
				core.View{Style: core.Style{
					Position: core.Absolute, Top: core.Pt(14), Left: core.Pt(12),
					Width: core.Pt(210), Height: core.Pt(22), CornerRadius: 5,
					BorderWidth: 2, Background: core.RGB(0x0e1622),
					Animated: &core.AnimatedStyle{BorderColor: core.Interpolate(open, core.Interpolation{
						Input: []float64{0, 1}, Output: []float64{0, 1},
						Colors: []core.Color{uipLine, uipBlue},
					})},
					Direction: core.Row, AlignItems: core.AlignCenter, Padding: core.Edges{Left: 8}, Gap: 4,
				}, Children: []core.Node{
					core.Text{Content: "Choose a city", Style: core.Style{Position: core.Absolute, FontSize: 9, Color: uipGrey, Animated: &core.AnimatedStyle{Opacity: placeholder}}},
					core.Text{Content: "Nagoya", Style: core.Style{Position: core.Absolute, FontSize: 9, FontWeight: core.FontWeightBold, Color: uipInk, Animated: &core.AnimatedStyle{Opacity: picked}}},
					core.View{Style: core.Style{FlexGrow: 1}},
					core.Text{Content: "v", Style: core.Style{FontSize: 9, Color: uipGrey, Animated: &core.AnimatedStyle{Rotate: lerp(open, 0, 180)}}},
					core.View{Style: core.Style{Width: core.Pt(4)}},
				}},
				core.View{Style: core.Style{
					Position: core.Absolute, Top: core.Pt(40), Left: core.Pt(12),
					Width: core.Pt(210), CornerRadius: 6, Background: uipPaper, Padding: core.Edges{Top: 5, Right: 5, Bottom: 5, Left: 5}, Gap: 2,
					Animated: &core.AnimatedStyle{Opacity: open, TranslateY: lerp(open, -6, 0)},
				}, Children: []core.Node{core.View{Style: core.Style{Gap: 2}, Children: list}}},
				core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(14), Left: core.Pt(250)}, Children: []core.Node{
					uipSheet(216, 80, 8, nil,
						uipLabel("Selection", 10, uipInk),
						uipBar(180, 6, uipLine, nil),
						uipBar(140, 6, uipLine, nil),
						row(6, uipChip(66, 16, uipPale), uipChip(48, 16, uipBlue)),
					),
				}},
			)
		}})
}

// 055 Popover: a small panel opens beside the button that asked for it.
func uiPopover(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "055", Title: "Popover", Duration: 2200 * ms, Accent: colAmber, Height: 104,
		Note: "The panel belongs to the button, so it slides the last six points in from that side and points back at it.",
		Stage: func(p *core.AnimatedValue) core.Node {
			a := openClose(p, 0.5)
			return uipScreen(uipW, 96,
				At(8, 12, uipBar(140, 8, uipLine, nil)),
				At(24, 12, uipBar(320, 7, uipLine, nil)),
				At(38, 12, uipBar(260, 6, uipLine, nil)),
				core.View{Style: core.Style{
					Position: core.Absolute, Top: core.Pt(58), Left: core.Pt(12),
					Width: core.Pt(96), Height: core.Pt(22), CornerRadius: 5,
					BorderWidth: 2, BorderColor: uipBlue, Background: core.RGB(0x0e1622),
					AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter,
				}, Children: []core.Node{uipLabel("Details", 10, uipBlue)}},
				core.View{Style: core.Style{
					Position: core.Absolute, Top: core.Pt(18), Left: core.Pt(126),
					Animated: &core.AnimatedStyle{Opacity: a, TranslateX: lerp(a, 6, 0)},
				}, Children: []core.Node{
					row(0,
						core.View{Style: core.Style{Width: core.Pt(9), AlignItems: core.AlignCenter}, Children: []core.Node{
							core.View{Style: core.Style{Width: core.Pt(9), Height: core.Pt(9), Background: uipPaper,
								Animated: &core.AnimatedStyle{Rotate: lerp(a, 45, 45)}}},
						}},
						uipSheet(210, 68, 8, nil,
							uipLabel("Autosave", 10, uipInk),
							uipBar(170, 5, uipLine, nil),
							uipBar(130, 5, uipLine, nil),
							row(6,
								core.View{Style: core.Style{FlexGrow: 1}},
								uipChip(34, 14, uipLine),
								uipChip(38, 14, uipBlue),
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
				var on core.AnimatedNode
				switch i {
				case 1, 2:
					on = newOn
				case 3:
					on = freeOn
				}
				bg := uipLine
				fg := uipInk
				if on != nil {
					bg = uipPale
					fg = uipBlue
				}
				body := []core.Node{
					core.Text{Content: s, Style: core.Style{FontSize: 9, FontWeight: core.FontWeightBold, Color: fg}},
				}
				if on != nil {
					body = append([]core.Node{core.Text{Content: "v", Style: core.Style{FontSize: 9, FontWeight: core.FontWeightBold, Color: uipBlue, Animated: &core.AnimatedStyle{Opacity: on}}}}, body...)
				}
				chipStyle := core.Style{
					Height: core.Pt(18), CornerRadius: 9, Background: bg,
					Padding: core.Edges{Left: 10, Right: 10}, Direction: core.Row, AlignItems: core.AlignCenter, Gap: 5,
				}
				if on != nil {
					chipStyle.Animated = &core.AnimatedStyle{ScaleX: core.Interpolate(on, core.Interpolation{Input: []float64{0, 1}, Output: []float64{1, 1.18}})}
				}
				chips = append(chips, core.View{Style: chipStyle, Children: body})
			}
			dimNew := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.2, 0.8, 1}, Output: []float64{1, 0.2, 0.2, 1}})
			cw := float32((uipW - 24 - 3*10) / 4)
			tile := func(dim *core.AnimatedStyle) core.Node {
				return core.View{Style: core.Style{Width: core.Pt(core.Points(cw)), Height: core.Pt(44), CornerRadius: 6,
					Background: uipCard, Padding: core.Edges{Top: 6, Left: 6, Right: 6, Bottom: 6}, Gap: 4, Animated: dim},
					Children: []core.Node{uipBox(cw-12, 18, 4, uipGrey, nil), uipBar(cw*0.55, 5, uipGrey, nil)}}
			}
			tiles := []core.Node{tile(&core.AnimatedStyle{Opacity: dimNew}), tile(&core.AnimatedStyle{Opacity: dimNew}), tile(nil), tile(nil)}
			return uipScreen(uipW, 96,
				At(8, 12, core.View{Style: core.Style{Direction: core.Row, Gap: 7}, Children: chips}),
				At(42, 12, row(10, tiles...)),
			)
		}})
}

// 062 Range slider: two handles narrow the range; the blue part is the range itself.
func uiRange(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "062", Title: "Range slider", Duration: 2400 * ms, Accent: colCyan, Height: 104,
		Note: "The two handles are dragged inward and the blue part between them is the range; the bubbles carry the numbers.",
		Stage: func(p *core.AnimatedValue) core.Node {
			const track = 400.0
			const left = 40.0
			const trackTop = 46.0
			lo := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.5, 1}, Output: []float64{0.1 * track, 0.32 * track, 0.1 * track}})
			hi := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.45, 0.9, 1}, Output: []float64{1 * track, 0.8 * track, 0.8 * track, 1 * track}})
			span := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.45, 0.5, 0.9, 1}, Output: []float64{0.9, 0.48, 0.48, 0.48, 0.9}})
			shift := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.45, 0.5, 0.9, 1}, Output: []float64{0.05 * track, 0.26 * track, 0.26 * track, 0.26 * track, 0.05 * track}})
			return uipScreen(uipW, 96,
				At(8, 12, uipBar(80, 8, uipLine, nil)),
				At(trackTop, left, uipBar(track, 6, uipCard, nil)),
				At(trackTop, left, core.View{Style: core.Style{
					Width: core.Pt(track), Height: core.Pt(6), CornerRadius: 3, Background: uipBlue,
					Animated: &core.AnimatedStyle{ScaleX: span, TranslateX: shift},
				}}),
				At(trackTop-4, left, core.View{Style: core.Style{
					Position: core.Absolute, Left: core.Pt(0), Animated: &core.AnimatedStyle{TranslateX: lo},
				}, Children: []core.Node{rangeHandle()}}),
				At(trackTop-4, left, core.View{Style: core.Style{
					Position: core.Absolute, Left: core.Pt(0), Animated: &core.AnimatedStyle{TranslateX: hi},
				}, Children: []core.Node{rangeHandle()}}),
				At(trackTop-20, left-8, core.View{Style: core.Style{
					Position: core.Absolute, Left: core.Pt(0), Animated: &core.AnimatedStyle{TranslateX: lo},
				}, Children: []core.Node{uipChip(30, 14, uipBlue)}}),
				At(trackTop-20, left-8, core.View{Style: core.Style{
					Position: core.Absolute, Left: core.Pt(0), Animated: &core.AnimatedStyle{TranslateX: hi},
				}, Children: []core.Node{uipChip(30, 14, uipBlue)}}),
				At(trackTop+9, left, uipBar(16, 5, uipGrey, nil)),
				At(trackTop+9, left+track-18, uipBar(18, 5, uipGrey, nil)),
			)
		}})
}

// rangeHandle is a white knob with a blue rim.
func rangeHandle() core.Node {
	return core.View{Style: core.Style{
		Width: core.Pt(14), Height: core.Pt(14), CornerRadius: 7,
		Background: uipPaper, BorderWidth: 2, BorderColor: uipBlue,
	}}
}

// 064 Pagination: the marker circle travels to the page that was pressed, and
// the list above dips out and back with the new page in it.
func uiPager(cx *core.Scope) core.Node {
	return Fx(cx, FxSpec{ID: "064", Title: "Pagination", Duration: 2100 * ms, Accent: colViolet, Height: 104,
		Note: "One marker circle travels between the page numbers; the list above blinks out for a moment as its contents change.",
		Stage: func(p *core.AnimatedValue) core.Node {
			marker := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.33, 0.66, 1}, Output: []float64{0, 32, 64, 0}})
			dip := core.Interpolate(p, core.Interpolation{
				Input:  []float64{0, 0.3, 0.33, 0.63, 0.66, 0.96, 1},
				Output: []float64{1, 1, 0.4, 1, 0.4, 1, 1},
			})
			var rows []core.Node
			for i := 0; i < 3; i++ {
				rows = append(rows, row(8, uipChip(16, 14, uipLine), core.View{Style: core.Style{Gap: 3}, Children: []core.Node{
					uipBar(300, 6, uipLine, nil), uipBar(180, 5, uipGrey, nil),
				}}))
			}
			var nums []core.Node
			for i := 0; i < 5; i++ {
				nums = append(nums, core.View{Style: core.Style{
					Width: core.Pt(20), Height: core.Pt(20), AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter,
				}, Children: []core.Node{uipLabel(string(rune('1'+i)), 10, uipGrey)}})
			}
			return uipScreen(uipW, 96,
				At(6, 78, core.View{Style: core.Style{Gap: 6, Animated: &core.AnimatedStyle{Opacity: dip}}, Children: rows}),
				At(66, 166, core.View{Style: core.Style{Direction: core.Row, Gap: 12, Position: core.Relative}, Children: []core.Node{
					core.View{Style: core.Style{
						Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0),
						Width: core.Pt(20), Height: core.Pt(20), CornerRadius: 10, Background: uipBlue,
						Animated: &core.AnimatedStyle{TranslateX: marker},
					}},
					core.View{Style: core.Style{Direction: core.Row, Gap: 12}, Children: nums},
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
			return uipScreen(uipW, 96,
				At(8, 12, row(4, uipLabel("Home", 9, uipBlue), uipLabel(">", 9, uipGrey),
					core.Text{Content: "Courses", Style: core.Style{FontSize: 9, FontWeight: core.FontWeightBold, Color: uipBlue,
						Animated: &core.AnimatedStyle{TextColor: core.Interpolate(up, core.Interpolation{Input: []float64{0, 1}, Output: []float64{0, 1}, Colors: []core.Color{uipBlue, uipInk}})}}},
					core.View{Style: core.Style{Animated: &core.AnimatedStyle{Opacity: core.Interpolate(up, core.Interpolation{Input: []float64{0, 1}, Output: []float64{1, 0}})}}, Children: []core.Node{
						row(4, uipLabel(">", 9, uipGrey), uipLabel("Web design", 9, uipInk)),
					}},
				)),
				At(24, 12, uipBar(uipW-24, 1, uipLine, nil)),
				core.View{Style: core.Style{
					Position: core.Absolute, Top: core.Pt(36), Left: core.Pt(12),
					Animated: &core.AnimatedStyle{Opacity: detail},
				}, Children: []core.Node{core.View{Style: core.Style{Gap: 6},
					Children: []core.Node{
						uipBar(140, 10, uipInk, nil),
						uipBar(300, 6, uipLine, nil),
						uipBar(240, 6, uipLine, nil),
						uipBar(200, 26, uipLine, nil),
					}}}},
				core.View{Style: core.Style{
					Position: core.Absolute, Top: core.Pt(36), Left: core.Pt(12),
					Animated: &core.AnimatedStyle{Opacity: up},
				}, Children: []core.Node{row(12,
					core.View{Style: core.Style{Width: core.Pt(148), Height: core.Pt(52), CornerRadius: 6, Background: uipPaper, Padding: core.Edges{Top: 6, Left: 6, Right: 6, Bottom: 6}, Gap: 4}, Children: []core.Node{uipBox(136, 20, 4, uipGrey, nil), uipBar(96, 6, uipInk, nil), uipBar(64, 4, uipGrey, nil)}},
					core.View{Style: core.Style{Width: core.Pt(148), Height: core.Pt(52), CornerRadius: 6, Background: uipPaper, Padding: core.Edges{Top: 6, Left: 6, Right: 6, Bottom: 6}, Gap: 4}, Children: []core.Node{uipBox(136, 20, 4, uipGrey, nil), uipBar(96, 6, uipInk, nil), uipBar(64, 4, uipGrey, nil)}},
					core.View{Style: core.Style{Width: core.Pt(148), Height: core.Pt(52), CornerRadius: 6, Background: uipPaper, Padding: core.Edges{Top: 6, Left: 6, Right: 6, Bottom: 6}, Gap: 4}, Children: []core.Node{uipBox(136, 20, 4, uipGrey, nil), uipBar(96, 6, uipInk, nil), uipBar(64, 4, uipGrey, nil)}},
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
					links = append(links, uipBar(float32(46+k*8), 5, col, nil))
				}
				cols = append(cols, core.View{Style: core.Style{Gap: 5, Width: core.Pt(100)}, Children: links})
			}
			cols = append(cols, uipBox(110, 54, 6, core.RGB(0x24303f), nil))
			return uipScreen(uipW, 102,
				At(8, 12, core.View{Style: core.Style{Direction: core.Row, Gap: 18}, Children: []core.Node{
					uipLabel("Products", 9, uipBlue),
					uipLabel("Pricing", 9, uipGrey),
					uipLabel("Support", 9, uipGrey),
				}}),
				At(24, 12, uipBar(uipW-24, 1, uipLine, nil)),
				At(32, 12, uipBar(110, 8, uipLine, nil)),
				At(46, 12, uipBar(180, 34, uipLine, nil)),
				core.View{Style: core.Style{
					Position: core.Absolute, Top: core.Pt(28), Left: core.Pt(0),
					Animated: &core.AnimatedStyle{Opacity: a, TranslateY: lerp(a, -12, 0)},
				}, Children: []core.Node{
					uipSheet(uipW, 76, 6, nil,
						core.View{Style: core.Style{Direction: core.Row, Gap: 18, AlignItems: core.AlignCenter}, Children: cols},
					),
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
			const step = 130.0
			const start = 52.0
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
					Position: core.Absolute, FontSize: 11, FontWeight: core.FontWeightBold, Color: uipInk,
				}}
				check := core.Text{Content: "v", Style: core.Style{
					Position: core.Absolute, FontSize: 11, FontWeight: core.FontWeightBold, Color: uipBlue,
					Animated: &core.AnimatedStyle{Opacity: done},
				}}
				circles = append(circles, core.View{Style: core.Style{
					Width: core.Pt(20), Height: core.Pt(20), CornerRadius: 10,
					Background: core.RGB(0x0e1622),
					Animated: &core.AnimatedStyle{Background: core.Interpolate(done, core.Interpolation{
						Input: []float64{0, 1}, Output: []float64{0, 1}, Colors: []core.Color{core.RGB(0x0e1622), uipPale},
					})},
					BorderWidth: 2, BorderColor: uipBlue,
					AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter,
				}, Children: []core.Node{mark, check}})
			}
			return uipScreen(uipW, 96,
				At(22, start, core.View{Style: core.Style{Position: core.Relative, Width: core.Pt(3 * step), Height: core.Pt(20)}, Children: []core.Node{
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(9), Left: core.Pt(10), Width: core.Pt(3 * step), Height: core.Pt(4), CornerRadius: 2, Background: uipCard}},
					core.View{Style: core.Style{
						Position: core.Absolute, Top: core.Pt(9), Left: core.Pt(10),
						Width: core.Pt(3 * step), Height: core.Pt(4), CornerRadius: 2, Background: uipBlue,
						Animated: &core.AnimatedStyle{ScaleX: grow, TranslateX: lerp(grow, 0, -1.5*step)},
					}},
					core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(0), Left: core.Pt(0), Direction: core.Row, Gap: core.Points(step - 20)}, Children: circles},
				}}),
				At(52, 12, uipBar(160, 8, uipLine, nil)),
				At(66, 12, uipBar(320, 7, uipLine, nil)),
				At(78, 400, uipChip(64, 18, uipBlue)),
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
			shift := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.18, 0.24, 0.86, 1}, Output: []float64{0, 0, -60, -60, 0}})
			shift2 := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.36, 0.42, 0.86, 1}, Output: []float64{0, 0, -60, -60, 0}})
			lift := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.36, 0.42, 0.86, 1}, Output: []float64{0, 0, -18, -18, 0}})
			empty := core.Interpolate(p, core.Interpolation{Input: []float64{0, 0.45, 0.56, 0.9, 1}, Output: []float64{0, 0, 1, 1, 0}})
			rowCard := func(barW float32) core.Node {
				return core.View{Style: core.Style{
					Width: core.Pt(uipW - 24), Height: core.Pt(26), CornerRadius: 5,
					BorderWidth: 1, BorderColor: uipLine, Padding: core.Edges{Left: 8},
					Direction: core.Row, AlignItems: core.AlignCenter, Gap: 8,
				}, Children: []core.Node{uipChip(14, 14, uipLine), uipBar(barW, 7, uipLine, nil), core.View{Style: core.Style{FlexGrow: 1}}, uipChip(9, 11, uipGrey)}}
			}
			return uipScreen(uipW, 102,
				At(6, 12, uipBar(60, 8, uipLine, nil)),
				core.View{Style: core.Style{
					Position: core.Absolute, Top: core.Pt(22), Left: core.Pt(12),
					Animated: &core.AnimatedStyle{Opacity: first, TranslateX: shift},
				}, Children: []core.Node{rowCard(300)}},
				core.View{Style: core.Style{
					Position: core.Absolute, Top: core.Pt(56), Left: core.Pt(12),
					Animated: &core.AnimatedStyle{Opacity: second, TranslateX: shift2, TranslateY: lift},
				}, Children: []core.Node{rowCard(240)}},
				core.View{Style: core.Style{
					Position: core.Absolute, Top: core.Pt(60), Left: core.Pt(0), Width: core.Pt(uipW), AlignItems: core.AlignCenter,
					Animated: &core.AnimatedStyle{Opacity: empty, TranslateY: lerp(empty, 8, 0)},
				}, Children: []core.Node{
					core.View{Style: core.Style{Gap: 4, AlignItems: core.AlignCenter}, Children: []core.Node{
						core.View{Style: core.Style{Width: core.Pt(26), Height: core.Pt(16), CornerRadius: 4, BorderWidth: 2, BorderColor: uipBlue}},
						uipLabel("Nothing here yet", 10, uipGrey),
						uipChip(56, 15, uipBlue),
					}},
				}},
			)
		}})
}
