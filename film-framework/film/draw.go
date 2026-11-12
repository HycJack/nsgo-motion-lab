package film

import (
	"math"

	"github.com/oneo-me/nsgo/core"
)

// ---------------------------------------------------------------- 层透明度

// curAlpha 是当前这一幕的整体透明度（交叉溶解用）。
// 每个透明度都经过 A() 走一遍，于是「幕透明度 × 元素透明度」一次算清。
var curAlpha float32 = 1

// A 把元素自己的透明度乘上当前幕的透明度。
// 分镜里自己动手 SetAlpha 的地方都要走这一遍。
func A(a float32) float32 { return curAlpha * a }

// ---------------------------------------------------------------- 图

// Image 是一张带像素尺寸的贴图（宽高比用 AR()，画时按比例算宽）。
type Image interface {
	Source() core.ImageSource
	AR() float32
}

// DrawPic 画一张图，铺进 (x, y, w, h)，a 是它自己的透明度。
func DrawPic(d *core.Drawing, p Image, x, y, w, h, a float32) {
	if a <= 0.001 {
		return
	}
	d.SetAlpha(A(a))
	d.DrawImageIn(p.Source(), x, y, w, h)
}

// PicAt 画一张图，给定水平中心、底部（如地面）位置和高度，
// 按原图比例算宽度。
func PicAt(d *core.Drawing, p Image, cx, bottomY, h, a float32) {
	w := h * p.AR()
	DrawPic(d, p, cx-w/2, bottomY-h, w, h, a)
}

// ---------------------------------------------------------------- 光

// Glow 在 (x, y) 画一团径向渐变的光晕。
func Glow(d *core.Drawing, x, y, r float32, c core.Color, a float32) {
	if r <= 0 || a <= 0.001 {
		return
	}
	d.SetFillGradient(core.CanvasRadial(x, y, 0, x, y, r,
		core.GradientStop{Offset: 0, Color: c.Faded(a)},
		core.GradientStop{Offset: 0.3, Color: c.Faded(a * 0.42)},
		core.GradientStop{Offset: 1, Color: c.Faded(0)},
	))
	d.BeginPath()
	d.Arc(x, y, r, 0, 2*math.Pi, false)
	d.Fill()
}

// GroundPool 是落在地上的一片光（被压扁的椭圆）。
func GroundPool(d *core.Drawing, x, y, rx, ry float32, c core.Color, a float32) {
	if rx <= 0 || a <= 0.001 {
		return
	}
	d.SetFillGradient(core.CanvasRadial(x, y, 0, x, y, rx,
		core.GradientStop{Offset: 0, Color: c.Faded(a)},
		core.GradientStop{Offset: 0.55, Color: c.Faded(a * 0.35)},
		core.GradientStop{Offset: 1, Color: c.Faded(0)},
	))
	d.BeginPath()
	d.Ellipse(x, y, rx, ry)
	d.Fill()
}

// Rays 是光芒：绕 (x, y) 均匀伸出 n 根线，phase 让它们能转起来。
func Rays(d *core.Drawing, x, y, r0, r1 float32, c core.Color, a float32, n int, phase float64) {
	d.SetStrokeColor(c.Faded(a))
	d.SetAlpha(A(1))
	d.SetLineWidth(2)
	for i := 0; i < n; i++ {
		ang := phase + float64(i)*2*math.Pi/float64(n)
		d.BeginPath()
		d.MoveTo(x+float32(math.Cos(ang))*r0, y+float32(math.Sin(ang))*r0)
		d.LineTo(x+float32(math.Cos(ang))*r1, y+float32(math.Sin(ang))*r1)
		d.Stroke()
	}
}

// ---------------------------------------------------------------- 字

// Txt 在 (x, y) 画一行字。
func Txt(d *core.Drawing, s string, x, y, size float32, fam string, wt core.FontWeight, c core.Color, a float32) {
	d.SetFillColor(c)
	d.SetAlpha(A(a))
	d.FillText(s, x, y, core.Font{Family: fam, Size: core.Points(size), Weight: wt})
}

// Center 画一行水平居中的字。
func Center(d *core.Drawing, s string, cx, y, size float32, fam string, wt core.FontWeight, c core.Color, a float32) {
	w := d.MeasureText(s, core.Font{Family: fam, Size: core.Points(size), Weight: wt})
	Txt(d, s, cx-w/2, y, size, fam, wt, c, a)
}

// ---------------------------------------------------------------- 氛围

// Hash 是一个确定性的伪随机：同一个 i 永远给同一个 0..1 的数。
// 星星、灯火的位置必须每次都一样，否则拍两遍不一样。
func Hash(i int) float64 {
	x := uint64(i)*6364136223846793005 + 1442695040888963407
	x ^= x >> 33
	x *= 0xff51afd7ed558ccd
	x ^= x >> 33
	return float64(x%100000) / 100000.0
}

// Scrim 压一层纯色，把素材统一到同一个世界里。
func Scrim(d *core.Drawing, w, h float32, a float32, c core.Color) {
	d.SetFillColor(c)
	d.SetAlpha(A(a))
	d.BeginPath()
	d.Rect(0, 0, w, h)
	d.Fill()
}

// Vignette 四角压暗。
func Vignette(d *core.Drawing, w, h float32, a float32) {
	d.SetFillGradient(core.CanvasRadial(w/2, h/2, h*0.35, w/2, h/2, h*0.95,
		core.GradientStop{Offset: 0, Color: core.RGB(0x000000).Faded(0)},
		core.GradientStop{Offset: 1, Color: core.RGB(0x000000).Faded(a)},
	))
	d.SetAlpha(A(1))
	d.BeginPath()
	d.Rect(0, 0, w, h)
	d.Fill()
}
