// 《登录流程》—— 流程图教学片：.flow 文件 → flowkit → 视频。
//
// 流程图数据来自一个真正的 .flow 文本文件（nsgo-flow 的格式，可手写
// 也可由 Flow Editor 导出），flowkit 解析 + 分层布局 + 逐节点/逐边动画。
//
// 三幕（Series，DarkMode 风格）：
//  1. 标题卡
//  2. 流程图动画：节点逐个浮现（start→input→check→lock/home），
//     边在两端就位后沿折线生长、箭头落位、标签浮现
//  3. 状态语义：done=绿、blocked=红——一条边走向两个结局
package main

import (
	_ "embed"

	"nsgofilm/film"
	"nsgofilm/filmkit"
	"nsgofilm/flowkit"
	"nsgofilm/tpl"

	"github.com/oneo-me/nsgo/core"

	"github.com/oneo-me/nsgo-flow/flowfile"
)

//go:embed login.flow
var loginFlow string

const (
	dw   float32 = 1920.0
	dh   float32 = 1080.0
	fade         = 0.5
)

func main() {
	film.Run(film.Film{
		W:    dw,
		H:    dh,
		Fade: fade,
		Bg:   tpl.DarkMode.Bg,
		Shots: film.Series(fade,
			film.Shot{D: 5, Fn: tpl.Title(tpl.DarkMode,
				"登录流程", "一张 .flow 文件，一部教学视频")},
			film.Shot{D: 11, Fn: sceneChart},
			film.Shot{D: 5.5, Fn: tpl.Outro(tpl.DarkMode,
				"流程图 = 代码", ".flow 文本 → flowkit → 动画"),
				Out: film.FadeBlack,
			},
		),
		Subs: []film.Sub{
			{A: 1.0, B: 4.4, Text: "把登录的每一步画成流程图——节点是动作，箭头是走向。"},
			{A: 5.2, B: 10.8, Text: "校验是判断：成功进入首页，连续失败五次就锁定账号。"},
			{A: 11.6, B: 13.8, Text: "绿色是终点，红色是异常出口。"},
		},
		SubStyle: film.DefaultSubStyle(tpl.DarkMode.Heading, tpl.DarkMode.Accent, core.RGB(0x000000)),
	})
}

// sceneChart：流程图逐节点/逐边动画——flowkit 的时序表被拉伸到
// 恰好铺满本幕（reveal 0.4s 开始，11s 放完），全程纯函数。
func sceneChart(d *core.Drawing, f *film.Frame) {
	doc, err := flowfile.Parse(loginFlow)
	if err != nil {
		return
	}
	c := flowkit.Layout(doc, flowkit.TopToBottom, float64(dw)*0.16, float64(dh)*0.14, float64(dw)*0.68, float64(dh)*0.74, d)
	const start, span = 0.5, 8.5
	drawn := filmkit.Clamp01((f.LocalT-start)/span) * c.Duration()
	c.Draw(d, drawn)
}

// maxF/minF 是 float64 的 min/max（本地版）。
func maxF(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func minF(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
