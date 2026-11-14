package flowkit

// 绘制：边走引擎路由（圆角拐弯 + 箭头 + 虚线 + 药丸标签），节点是
// 状态卡片（圆角矩形 + 左侧 Kind 色条 + 状态角标）。reveal：节点
// stagger 浮现，边按权重从起点生长，标签在边描完时浮现。

import (

	"github.com/oneo-me/nsgo/core"
)

// Draw 画出时刻 t 的整张图（t 超过 Total 时定格在完成态）。
func (c Chart) Draw(d *core.Drawing, t float64) {
	pal := c.pal
	// 边在下层。
	for i := range c.Edges {
		e := &c.Edges[i]
		w := clamp01((t - e.RevealAt) / 0.5)
		if w <= 0.001 {
			continue
		}
		col := pal.Accent
		e.route.draw(d, defaultEdgeWidth, col, e.Dashed, w)
		if w >= 1 && e.Label != "" {
			e.route.label(d, e.Label, pal.Text, pal.LabelBg, pal.LabelLine,
				filmkitA(clamp01((t-e.RevealAt-0.5)/0.4)))
		}
	}
	// 节点在上层。
	for i := range c.Nodes {
		nd := &c.Nodes[i]
		a := clamp01((t - nd.RevealAt) / 0.45)
		if a <= 0.001 {
			continue
		}
		drawCard(d, nd, pal, a)
		if a >= 1 && nd.Note != "" {
			na := clamp01((t - nd.RevealAt - 0.45) / 0.4)
			if na > 0 {
				coreText(d, nd.Note, float32(nd.X), float32(nd.Y+nd.H+30), 26, pal.Muted, float32(na), 0)
			}
		}
	}
}

// drawCard 画一个状态卡片：圆角矩形 + 左侧 Kind 色条 + 状态角标。
var violetFlow = core.RGB(0xA78BFA)

func drawCard(d *core.Drawing, nd *Node, pal Palette, alpha float64) {
	a := float32(alpha)
	x, y := float32(nd.X), float32(nd.Y)
	w, h := float32(nd.W), float32(nd.H)

	fill, stroke := pal.NodeFill, pal.NodeStroke
	switch nd.Status {
	case "done":
		stroke, fill = pal.Done, dimFill(pal.Done)
	case "doing":
		stroke, fill = pal.Doing, dimFill(pal.Doing)
	case "blocked":
		stroke, fill = pal.Blocked, dimFill(pal.Blocked)
	}
	// Kind 色条：start 绿、end 绿、decision 金、data 紫。
	bar := pal.Accent
	switch nd.Kind {
	case "start", "end":
		bar = pal.Done
	case "decision":
		bar = pal.Doing
	case "data":
		bar = violetFlow
	}

	d.SetFillColor(fill)
	d.SetAlpha(a)
	d.BeginPath()
	d.RoundRect(x, y, w, h, 12)
	d.Fill()
	d.SetStrokeColor(stroke)
	d.SetLineWidth(2.5)
	d.BeginPath()
	d.RoundRect(x, y, w, h, 12)
	d.Stroke()
	// 左侧色条。
	d.SetFillColor(bar)
	d.BeginPath()
	d.Rect(x, y+8, 6, h-16)
	d.Fill()

	// 标签（居中，过长缩小）。
	size := float32(30)
	font := core.Font{Family: pal.Font, Size: core.Points(size)}
	for d.MeasureText(nd.Label, font) > w-40 && size > 20 {
		size -= 2
		font = core.Font{Family: pal.Font, Size: core.Points(size)}
	}
	lw := d.MeasureText(nd.Label, font)
	d.SetFillColor(pal.Text)
	// 显式设基线：引擎路由的 label() 把基线设成 Top 且状态泄漏，
	// 不重设的话节点文字会挂到卡片下缘外。
	d.SetTextBaseline(core.BaselineMiddle)
	d.FillText(nd.Label, x+w/2-lw/2, y+h/2, font)

	// 状态角标。
	mark, mc := "", pal.Muted
	switch nd.Status {
	case "done":
		mark, mc = "✓", pal.Done
	case "doing":
		mark, mc = "●", pal.Doing
	case "blocked":
		mark, mc = "✕", pal.Blocked
	}
	if mark != "" {
		d.SetFillColor(mc)
		d.FillText(mark, x+w-34, y+30, core.Font{Family: pal.Font, Size: core.Points(26)})
	}
	d.SetAlpha(1)
}

func coreText(d *core.Drawing, s string, x, y, size float32, c core.Color, alpha float32, align float32) {
	if alpha <= 0.001 || s == "" {
		return
	}
	w := d.MeasureText(s, core.Font{Size: core.Points(size)})
	d.SetFillColor(c)
	d.SetAlpha(alpha)
	d.SetTextBaseline(core.BaselineMiddle)
	d.FillText(s, x-align*w, y, core.Font{Size: core.Points(size)})
}

