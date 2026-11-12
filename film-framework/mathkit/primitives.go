package mathkit

import (
	"math"

	"github.com/oneo-me/nsgo/core"
)

// Arrow 画一支从 (x0,y0) 指向 (x1,y1) 的箭头：直线 + 实心三角头。
// head 是箭头头的长度（像素）。教学片里「指一下」的原语。
func Arrow(d *core.Drawing, x0, y0, x1, y1, head float32, c core.Color, alpha float32, width float32) {
	if alpha <= 0.001 {
		return
	}
	dx, dy := x1-x0, y1-y0
	l := float32(math.Hypot(float64(dx), float64(dy)))
	if l < 1 {
		return
	}
	ux, uy := dx/l, dy/l
	px, py := -uy, ux
	if head <= 0 {
		head = 14
	}
	// 主线画到箭头底部，别戳进三角里。
	bx, by := x1-ux*head, y1-uy*head
	d.SetStrokeColor(c)
	d.SetAlpha(alpha)
	d.SetLineWidth(width)
	d.BeginPath()
	d.MoveTo(x0, y0)
	d.LineTo(bx, by)
	d.Stroke()
	// 实心三角头。
	d.SetFillColor(c)
	d.BeginPath()
	d.MoveTo(x1, y1)
	d.LineTo(bx+px*head*0.5, by+py*head*0.5)
	d.LineTo(bx-px*head*0.5, by-py*head*0.5)
	d.ClosePath()
	d.Fill()
	d.SetAlpha(1)
}

// Circle 画圆：p 是描画进度 0..1（从正右方起，画满一整圈）。
// manim 的 Create 对圆的处理。
func Circle(d *core.Drawing, cx, cy, r float32, p float64, c core.Color, alpha float32, width float32) {
	if alpha <= 0.001 || p <= 0 || r <= 0 {
		return
	}
	d.SetStrokeColor(c)
	d.SetAlpha(alpha)
	d.SetLineWidth(width)
	d.BeginPath()
	d.Arc(cx, cy, r, 0, float32(p*2*math.Pi), false)
	d.Stroke()
	d.SetAlpha(1)
}

// Dot 画一个实心圆点（manim 的 Dot）：深芯 + 亮环，比纯色圆点醒目。
func Dot(d *core.Drawing, cx, cy, r float32, c core.Color, alpha float32) {
	if alpha <= 0.001 || r <= 0 {
		return
	}
	d.SetFillColor(c)
	d.SetAlpha(alpha)
	d.BeginPath()
	d.Arc(cx, cy, r, 0, float32(2*math.Pi), false)
	d.Fill()
	// 芯：底色调暗一半，形成内外两圈。
	dr, dg, db, da := c.Components()
	dim := func(v float64) uint32 { return uint32(v*0.45*255 + 0.5) }
	d.SetFillColor(core.Color(dim(dr)<<24 | dim(dg)<<16 | dim(db)<<8 | uint32(da*255+0.5)))
	d.BeginPath()
	d.Arc(cx, cy, r*0.55, 0, float32(2*math.Pi), false)
	d.Fill()
	d.SetAlpha(1)
}
