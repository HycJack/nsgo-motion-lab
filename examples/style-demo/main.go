// 《Style Demo》—— tpl 模版系统的风格画廊。
//
// 六幕六风格（对应 html-ppt-designer STYLES.md），每幕一个预制模版：
// TED 标题卡 → Keynote 要点 → Consulting 数据 → Editorial 引用 →
// Swiss 章节 → Neo-Tokyo 片尾。风格之间用黑场过渡（各风格气质差异大，
// 黑场是唯一不抢戏的缝）。
//
// 一部「PPT 式视频」的完整写法就长这样：选 Style + 按内容排 Series。
package main

import (
	"nsgofilm/film"
	"nsgofilm/tpl"
)

const (
	dw    float32 = 1920.0
	dh    float32 = 1080.0
	fade          = 0.5
)

func main() {
	film.Run(film.Film{
		W:    dw,
		H:    dh,
		Fade: fade,
		Bg:   tpl.TED.Bg, // 幕间底色；各幕自带自己的底色
		Shots: film.Series(fade,
			film.Shot{D: 5, Fn: tpl.Title(tpl.TED,
				"把想法拍成视频", "选一种风格，剩下的交给模版"),
			},
			film.Shot{D: 6.5, Fn: tpl.Bullets(tpl.Keynote,
				"为什么用模版拍视频",
				[]string{
					"配色、字体、动效一次定义，全片一致",
					"改内容不用改视觉：换数据就是新片子",
					"分镜是代码——精确到帧，改完即录",
				}),
			},
			film.Shot{D: 6, Fn: tpl.KPI(tpl.Consulting,
				"三个月的效率提升",
				[]tpl.KPIItem{
					{Value: 3.2, Unit: "×", Label: "产出速度", Decimals: 1},
					{Value: 47, Unit: "%", Label: "成本下降"},
					{Value: 12, Unit: "h", Label: "每周省下"},
				}),
			},
			film.Shot{D: 6, Fn: tpl.Quote(tpl.Editorial,
				"设计的本质不是装饰，而是传达。",
				"Paul Rand"),
			},
			film.Shot{D: 4.5, Fn: tpl.Section(tpl.Swiss, 2, "开始你的第一部片子")},
			film.Shot{D: 5, Fn: tpl.Outro(tpl.NeoTokyo,
				"MAKE IDEAS MOVE", "nsgo film · 声明式视频框架"),
				Out: film.FadeBlack,
			},
		),
	})
}
