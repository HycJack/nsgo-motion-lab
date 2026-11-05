// Command Motion Lab is a tour of what NSGo draws with: core.Canvas
// records 2D drawing, the animation graph moves values on the platform's
// frame clock, core.AnimatedStyle carries those values into ordinary
// component state, and core.Shader runs Kage fragment shaders on the GPU.
//
// Each page holds a handful of small, self-contained examples.
package main

import (
	"github.com/oneo-me/nsgo"
	"github.com/oneo-me/nsgo/base"
	"github.com/oneo-me/nsgo/core"
)

// pageDef is one tab.
type pageDef struct {
	name string
	view func(cx *core.Scope) core.Node
}

var pages = []pageDef{
	{name: "Canvas", view: canvasPage},
	{name: "Animation", view: animationPage},
	{name: "UI Kit", view: uikitPage},
	{name: "Shader", view: shaderPage},
	// The eight motion-prompt packs. The ui pack (049–072) lives at the
	// bottom of the UI Kit page, next to the components it belongs to.
	{name: "Motion", view: pageMotion},
	{name: "Camera", view: pageCamera},
	{name: "Landing", view: pageLP},
	{name: "Photo", view: pagePhoto},
	{name: "Motion2", view: pageMotion2},
	{name: "Landing2", view: pageLP2},
	{name: "Camera2", view: pageCamera2},
	// The camera packs move a picture in two dimensions; this page is the
	// three-dimensional version of the same moves.
	{name: "3D", view: page3D},
}

func main() {
	nsgo.Run(Root, nsgo.WindowSize(1180, 860))
}

// Root is the app's first screen: a title bar with the three tabs, and
// the page they choose.
func Root(cx *core.Scope) core.Node {
	current := core.NewSignal(cx, 0)
	return core.View{
		// FlexShrink matters as much as FlexGrow here: without it the
		// column grows to its content, the scroll view inherits that
		// height, and nothing ever overflows to scroll.
		Style: core.Style{FlexGrow: 1, FlexShrink: 1, Background: colBg, Direction: core.Column},
		Children: []core.Node{
			titleBar(cx, current),
			core.Scroll{
				TestID: "page-scroll",
				Style:  core.Style{FlexGrow: 1},
				Children: []core.Node{
					// Dynamic rebuilds only this span, so switching tabs
					// replaces the page and every demo on it restarts.
					core.Dynamic(func(cx *core.Scope) core.Node {
						return pages[current.Get()].view(cx)
					}),
				},
			},
		},
	}
}

// titleBar is the app's header: its name, and a tab per page.
func titleBar(cx *core.Scope, current *core.Signal[int]) core.Node {
	tabs := make([]core.Node, len(pages))
	for i := range pages {
		i := i
		tabs[i] = base.Button{
			Label:   pages[i].name,
			TestID:  "tab-" + pages[i].name,
			OnPress: func(core.PressEvent) { current.Set(i) },
			Selected: func() bool {
				return current.Get() == i
			},
			StyleFunc: func() core.Style {
				background := core.RGBA(0xffffff, 0.04)
				if current.Get() == i {
					background = core.RGBA(0x4d8dff, 0.22)
				}
				return core.Style{Padding: core.Symmetric(7, 14), CornerRadius: 9, Background: background}
			},
			Children: []core.Node{core.Text{
				Content: pages[i].name,
				StyleFunc: func() core.Style {
					colour := colMuted
					if current.Get() == i {
						colour = colText
					}
					return core.Style{FontSize: 13, FontWeight: core.FontWeightMedium, Color: colour}
				},
			}},
		}
	}
	header := []core.Node{
		core.Text{Content: "NSGo Motion Lab", Style: core.Style{FontSize: 16, FontWeight: core.FontWeightBold, Color: colText}},
		core.Text{Content: "canvas · animation · ui · shader · 192 motion examples · 3d", Style: core.Style{FontSize: 12, Color: colFaint}},
	}
	// The tabs wrap under the title on a narrow window, so they stay
	// reachable there.
	return core.View{
		Style: core.Style{
			Padding: core.Symmetric(12, 18), Background: colBar, Gap: 8,
			BorderWidth: 1, BorderColor: colLine,
		},
		Children: []core.Node{
			row(12, header...),
			core.View{Style: core.Style{Direction: core.Row, Wrap: core.Wrap, Gap: 8}, Children: tabs},
		},
	}
}
