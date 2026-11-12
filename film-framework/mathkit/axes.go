package mathkit

import (
	"math"
	"strconv"

	"github.com/oneo-me/nsgo/core"
)

// Axes 是 2D 坐标系：数据范围映射到画布上的一块矩形。
// 轴穿过数据原点（不在范围内就钳到最近的边），manim 的样子：
// 箭头、轴上刻度、可选的刻度数字。
type Axes struct {
	X, Y [2]float64 // 数据范围
	At   [2]float32 // 画布左上角
	Size [2]float32 // 画布宽高（像素）
	Nums bool       // 画刻度数字
	Tick int        // 每轴的目标刻度数，0 = 4
}

// Screen 把数据点映射到画布坐标。
func (a Axes) Screen(x, y float64) (float32, float32) {
	sx := a.At[0] + float32((x-a.X[0])/(a.X[1]-a.X[0]))*a.Size[0]
	sy := a.At[1] + a.Size[1] - float32((y-a.Y[0])/(a.Y[1]-a.Y[0]))*a.Size[1]
	return sx, sy
}

// Draw 画坐标系。p 是建轴进度 0..1：两条轴从原点向两端长出来，
// 刻度与数字随后浮现。
func (a Axes) Draw(d *core.Drawing, p float64, c core.Color, alpha float32, width float32) {
	if alpha <= 0.001 || p <= 0 {
		return
	}
	d.SetStrokeColor(c)
	d.SetAlpha(alpha)
	d.SetLineWidth(width)

	// 轴的锚点：数据原点钳进范围。
	ox := clampd(0, a.X[0], a.X[1])
	oy := clampd(0, a.Y[0], a.Y[1])
	ocx, ocy := a.Screen(ox, oy)

	// 横轴：从原点向两端长。
	lx, _ := a.Screen(a.X[0], oy)
	rx, _ := a.Screen(a.X[1], oy)
	half1 := p
	d.BeginPath()
	d.MoveTo(mix(ocx, lx, half1), ocy)
	d.LineTo(mix(ocx, rx, half1), ocy)
	d.Stroke()
	// 纵轴。
	_, by := a.Screen(ox, a.Y[0])
	_, ty := a.Screen(ox, a.Y[1])
	d.BeginPath()
	d.MoveTo(ocx, mix(ocy, by, half1))
	d.LineTo(ocx, mix(ocy, ty, half1))
	d.Stroke()

	// 箭头：正方向的两端（轴长出来之后浮现）。
	arrA := alpha * float32(smooth((p-0.7)/0.3))
	if arrA > 0.001 {
		d.SetAlpha(arrA)
		a.arrow(d, rx, ocy, 1, 0, c, arrA)
		a.arrow(d, ocx, ty, 0, -1, c, arrA)
	}

	// 刻度与数字：轴长出来之后浮现。
	tickA := alpha * float32(smooth((p-0.6)/0.4))
	if tickA <= 0.001 {
		return
	}
	d.SetAlpha(tickA)
	target := a.Tick
	if target <= 0 {
		target = 5
	}
	const tick = 6
	step := niceStep(a.X[1]-a.X[0], target)
	for v := math.Ceil(a.X[0]/step) * step; v <= a.X[1]+1e-9; v += step {
		if math.Abs(v) < step/2 {
			continue // 原点不给自己画刻度
		}
		tx, _ := a.Screen(v, oy)
		d.BeginPath()
		d.MoveTo(tx, ocy-tick)
		d.LineTo(tx, ocy+tick)
		d.Stroke()
		if a.Nums {
			mathkitNum(d, formatNum(v), tx, ocy+22, 24, c, tickA, 0.5)
		}
	}
	step = niceStep(a.Y[1]-a.Y[0], target)
	for v := math.Ceil(a.Y[0]/step) * step; v <= a.Y[1]+1e-9; v += step {
		if math.Abs(v) < step/2 {
			continue
		}
		_, yy := a.Screen(ox, v)
		d.BeginPath()
		d.MoveTo(ocx-tick, yy)
		d.LineTo(ocx+tick, yy)
		d.Stroke()
		if a.Nums {
			mathkitNum(d, formatNum(v), ocx-10, yy+8, 24, c, tickA, 1)
		}
	}
	d.SetAlpha(1)
}

// arrow 在 (x,y) 沿屏幕方向 (dx,dy) 画一个小实心三角。
func (a Axes) arrow(d *core.Drawing, x, y, dx, dy float32, c core.Color, alpha float32) {
	const s = 12.0
	// 垂直方向。
	px, py := -dy, dx
	d.SetFillColor(c)
	d.SetAlpha(alpha)
	d.BeginPath()
	d.MoveTo(x+dx*s*1.4, y+dy*s*1.4)
	d.LineTo(x-dx*s+px*s*0.6, y-dy*s+py*s*0.6)
	d.LineTo(x-dx*s-px*s*0.6, y-dy*s-py*s*0.6)
	d.ClosePath()
	d.Fill()
	d.SetAlpha(1)
}

// niceStep 是好看的刻度步长：1/2/5 × 10^k。
func niceStep(span float64, target int) float64 {
	if target <= 0 || span <= 0 {
		return 1
	}
	raw := span / float64(target)
	mag := math.Pow(10, math.Floor(math.Log10(raw)))
	for _, m := range []float64{1, 2, 5, 10} {
		if raw <= m*mag {
			return m * mag
		}
	}
	return 10 * mag
}

func formatNum(v float64) string {
	if v == math.Trunc(v) && math.Abs(v) < 1e15 {
		return strconv.FormatInt(int64(v), 10)
	}
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func clampd(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// mathkitNum 是 Label 的本地版（mathkit 不依赖 film）。
func mathkitNum(d *core.Drawing, s string, x, y, size float32, c core.Color, alpha float32, align float32) {
	Label(d, s, x, y, size, c, alpha, align)
}
