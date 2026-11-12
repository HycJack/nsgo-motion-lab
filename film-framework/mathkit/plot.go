package mathkit

import (
	"math"

	"github.com/oneo-me/nsgo/core"
)

// Plot 画 y=fn(x)：从 x0 描到 x0+(x1-x0)*p——p 是描画进度 0..1，
// 就是 manim 的 Create。折线 steps 段近似（0 取 96）。
func Plot(d *core.Drawing, a Axes, fn func(float64) float64, x0, x1, p float64, steps int, c core.Color, alpha float32, width float32) {
	if alpha <= 0.001 || p <= 0 || steps <= 0 {
		return
	}
	xEnd := x0 + (x1-x0)*p
	d.SetStrokeColor(c)
	d.SetAlpha(alpha)
	d.SetLineWidth(width)
	d.BeginPath()
	for i := 0; i <= steps; i++ {
		x := x0 + (xEnd-x0)*float64(i)/float64(steps)
		sx, sy := a.Screen(x, fn(x))
		if i == 0 {
			d.MoveTo(sx, sy)
		} else {
			d.LineTo(sx, sy)
		}
	}
	d.Stroke()
	d.SetAlpha(1)
}

// Area 填充 fn 下方 [x0,x1] 与 y=0 之间的区域（积分图示）。
func Area(d *core.Drawing, a Axes, fn func(float64) float64, x0, x1 float64, steps int, c core.Color, alpha float32) {
	if alpha <= 0.001 || steps <= 0 {
		return
	}
	_, zy := a.Screen(0, 0) // y=0 的屏幕高度
	d.SetFillColor(c)
	d.SetAlpha(alpha)
	d.BeginPath()
	base, _ := a.Screen(x0, 0)
	d.MoveTo(base, zy)
	for i := 0; i <= steps; i++ {
		x := x0 + (x1-x0)*float64(i)/float64(steps)
		sx, sy := a.Screen(x, fn(x))
		d.LineTo(sx, sy)
	}
	ex, _ := a.Screen(x1, 0)
	d.LineTo(ex, zy)
	d.ClosePath()
	d.Fill()
	d.SetAlpha(1)
}

// Tangent 在 x0 处画切线（数值导数），halfLen 是屏幕像素的半长。
// 导数与切线动画：x0 随影片时钟滑动，切线就贴着曲线扫过去。
func Tangent(d *core.Drawing, a Axes, fn func(float64) float64, x0 float64, halfLen float32, c core.Color, alpha float32, width float32) {
	if alpha <= 0.001 {
		return
	}
	const eps = 1e-4
	slope := (fn(x0+eps) - fn(x0-eps)) / (2 * eps)
	// 屏幕方向：数据里 x 走一步，屏幕走 (dx, -slope*dy)。
	dx := a.Size[0] / float32(a.X[1]-a.X[0])
	dy := a.Size[1] / float32(a.Y[1]-a.Y[0])
	vx, vy := dx, -float32(slope)*dy
	l := float32(math.Sqrt(float64(vx*vx + vy*vy)))
	if l <= 0 {
		return
	}
	vx, vy = vx/l*halfLen, vy/l*halfLen
	cx, cy := a.Screen(x0, fn(x0))
	d.SetStrokeColor(c)
	d.SetAlpha(alpha)
	d.SetLineWidth(width)
	d.BeginPath()
	d.MoveTo(cx-vx, cy-vy)
	d.LineTo(cx+vx, cy+vy)
	d.Stroke()
	d.SetAlpha(1)
}

// Label 写一行字（屏幕坐标）：align 0 左对齐 / 0.5 居中 / 1 右对齐。
func Label(d *core.Drawing, s string, x, y, size float32, c core.Color, alpha float32, align float32) {
	if alpha <= 0.001 || s == "" {
		return
	}
	w := d.MeasureText(s, core.Font{Size: core.Points(size)})
	d.SetFillColor(c)
	d.SetAlpha(alpha)
	d.FillText(s, x-float32(align)*w, y, core.Font{Size: core.Points(size)})
	d.SetAlpha(1)
}
