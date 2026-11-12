package film

// 时间轴总览（?timeline=1 / NSGO_FILM_TIMELINE=1）：
// 每幕一帧的 contact sheet——render 是影片时钟的纯函数，
// 所以每张卡片就是一个固定 t 的 Canvas，不引入任何新状态。
// 点卡片：时钟与媒体 seek 到幕首，盖层收起，进入审片。

import (
	"fmt"

	"github.com/oneo-me/nsgo/core"
)

// timelineGrid 搭全屏的卡片网格。
func timelineGrid(cx *core.Scope, f Film, filmT *core.Signal[float64], inFilm *core.Signal[bool], wake func()) core.Node {
	f.resolve()
	const cardW = 560.0
	cardH := core.Points(cardW * f.H / f.W)

	var cards []core.Node
	for i := range f.Shots {
		sc := f.Shots[i]
		mid := (sc.A + sc.B) / 2
		card := core.View{
			Style: core.Style{
				Direction: core.Column, AlignItems: core.AlignStart,
				Gap: core.Points(10),
			},
			OnPress: func(core.PressEvent) {
				filmT.Set(sc.A)
				inFilm.Set(true)
				wake()
			},
			Children: []core.Node{
				core.Canvas{
					Style: core.Style{
						Width:        core.Pt(cardW),
						Height:       core.Pt(cardH),
						Background:   f.Bg,
						CornerRadius: 10,
					},
					Draw: func(d *core.Drawing, size core.Size) {
						render(d, size, mid, f)
					},
				},
				core.Text{
					Content: fmt.Sprintf("幕 %d　%.1fs – %.1fs", i+1, sc.A, sc.B),
					Style: core.Style{
						Color:    core.RGBA(0xffffff, 0.75),
						FontSize: core.Points(24),
					},
				},
			},
		}
		cards = append(cards, card)
	}

	grid := core.View{
		Style: core.Style{
			Direction: core.Row, Wrap: core.Wrap, AlignItems: core.AlignStart,
			Gap:          core.Points(28),
			AlignContent: core.AlignStart,
		},
		Children: cards,
	}

	return core.View{
		Style: core.Style{
			Position: core.Absolute, Left: core.Pt(0), Right: core.Pt(0),
			Top: core.Pt(0), Bottom: core.Pt(0),
			Direction: core.Column, AlignItems: core.AlignStart,
			Gap:        core.Points(24),
			Padding:    core.All(40),
			Background: f.Bg,
		},
		Children: []core.Node{
			core.Text{
				Content: fmt.Sprintf("时间轴 · %d 幕 · %.1fs —— 点卡片进入该幕审片", len(f.Shots), f.Duration),
				Style: core.Style{
					Color:      core.RGBA(0xffffff, 0.9),
					FontSize:   core.Points(34),
					FontWeight: core.FontWeightMedium,
				},
			},
			core.Scroll{
				Style:    core.Style{Width: core.Pct(100), Height: core.Pct(100)},
				Children: []core.Node{grid},
			},
		},
	}
}
