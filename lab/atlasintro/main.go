// Command atlasintro draws the Atlas intro film with NSGo's canvas:
// one core.Canvas covering the window, redrawn every frame from a clock,
// so the whole six-scene film is a pure function of time.
package main

import (
	"math"
	"time"

	"github.com/oneo-me/nsgo"
	"github.com/oneo-me/nsgo/core"
)

// ---------------------------------------------------------------- palette

var (
	cream = core.RGB(0xf4f2eb)
	card_ = core.RGB(0xfbfaf6)
	ink   = core.RGB(0x1b1a17)
	muted = core.RGB(0x6f6b61)
	faint = core.RGB(0x9a958a)
	hair  = core.RGB(0xd8d3c6)
	red   = core.RGB(0x9c3b2e)
	redBg = core.RGB(0xf6e4e0)
	grn   = core.RGB(0x1f7a4d)
	grnBg = core.RGB(0xe4efe7)
	blue  = core.RGB(0x2f5fd0)
	blueD = core.RGB(0x3d63b8)
	gold  = core.RGB(0xc8952f)
)

const (
	serif = "Songti SC"
	sans  = "PingFang SC"
)

const (
	total = 53.0
	dw    = 1920.0
	dh    = 1080.0
)

// curAlpha multiplies every alpha() call, so a scene can be faded as a whole
// during transitions even though its parts set their own opacities.
var curAlpha float32 = 1

func main() {
	nsgo.Run(Root, nsgo.WindowSize(dw, dh))
}

// Root redraws a full-window canvas every frame from the wall clock (or,
// during capture, from window.__atlasT set by the recorder).
func Root(cx *core.Scope) core.Node {
	clock := core.NewSignal(cx, 0.0)
	var loop func(time.Time)
	loop = func(at time.Time) {
		clock.Set(atlasTime())
		cx.RequestAnimationFrame(loop)
	}
	cx.RequestAnimationFrame(loop)
	return core.Canvas{
		Style: core.Style{Width: core.Pct(100), Height: core.Pct(100), Background: cream},
		Draw: func(d *core.Drawing, size core.Size) {
			render(d, float32(size.Width), float32(size.Height), clock.Get())
		},
	}
}

// ---------------------------------------------------------------- math

func clamp01(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}
func seg(t, a, b float64) float64 { return clamp01((t - a) / (b - a)) }
func lerp(a, b, p float64) float64 { return a + (b-a)*p }
func lerpf(a, b, p float32) float32 { return a + (b-a)*p }
func smooth(x float64) float64     { x = clamp01(x); return x * x * (3 - 2*x) }
func easeOut(x float64) float64 {
	x = clamp01(x)
	return 1 - math.Pow(1-x, 3)
}
func easeInOut(x float64) float64 {
	x = clamp01(x)
	if x < 0.5 {
		return 4 * x * x * x
	}
	return 1 - math.Pow(-2*x+2, 3)/2
}
func f32(x float64) float32 { return float32(x) }

// ---------------------------------------------------------------- render

func render(d *core.Drawing, w, h float32, t float64) {
	s := float32(math.Min(float64(w)/dw, float64(h)/dh))
	ox := (w - float32(dw)*s) / 2
	oy := (h - float32(dh)*s) / 2
	d.Save()
	d.Translate(ox, oy)
	d.Scale(s, s)

	// base wash
	d.SetAlpha(1)
	d.SetFillColor(cream)
	d.BeginPath()
	d.Rect(0, 0, float32(dw), float32(dh))
	d.Fill()

	// scene table: start, end, draw
	type sc struct {
		a, b float64
		fn   func(*core.Drawing, float64)
	}
	scenes := []sc{
		{0, 9.5, scene1},
		{9.5, 16, scene2},
		{16, 25, scene3},
		{25, 33, scene4},
		{33, 42, scene5},
		{42, total, scene6},
	}

	// find current scene
	idx := 0
	for i, scn := range scenes {
		if t >= scn.a {
			idx = i
		}
	}
	cur := scenes[idx]
	// dip-to-paper transition: the outgoing scene fades out, the incoming
	// fades in, with a clean cream beat between — no double-exposed titles.
	fadeIn := 0.6
	if idx > 0 && t < cur.a+fadeIn {
		prev := scenes[idx-1]
		e := (t - cur.a) / fadeIn
		if e < 0.5 {
			curAlpha = float32(1 - smooth(e/0.5))
		} else {
			curAlpha = float32(smooth((e - 0.5) / 0.5))
		}
		d.SetAlpha(curAlpha)
		if e < 0.5 {
			prev.fn(d, t)
		} else {
			cur.fn(d, t)
		}
		curAlpha = 1
		d.SetAlpha(1)
	} else {
		curAlpha = 1
		d.SetAlpha(1)
		cur.fn(d, t)
	}

	d.Restore()
}

// ---------------------------------------------------------------- primitives

func bg(d *core.Drawing) {
	d.SetFillColor(cream)
	d.BeginPath()
	d.Rect(0, 0, float32(dw), float32(dh))
	d.Fill()
}

func rect(d *core.Drawing, x, y, w, h float32) {
	d.BeginPath()
	d.Rect(x, y, w, h)
	d.Fill()
}

func rrect(d *core.Drawing, x, y, w, h, r float32) {
	d.BeginPath()
	d.RoundRect(x, y, w, h, r)
	d.Fill()
}

func strokeRRect(d *core.Drawing, x, y, w, h, r float32) {
	d.BeginPath()
	d.RoundRect(x, y, w, h, r)
	d.Stroke()
}

func txt(d *core.Drawing, s string, x, y, size float32, fam string, wt core.FontWeight, c core.Color) {
	d.SetFillColor(c)
	d.FillText(s, x, y, core.Font{Family: fam, Size: core.Points(size), Weight: wt})
}

func measure(d *core.Drawing, s string, size float32, fam string, wt core.FontWeight) float32 {
	return d.MeasureText(s, core.Font{Family: fam, Size: core.Points(size), Weight: wt})
}

func alpha(d *core.Drawing, a float32, fn func()) {
	f := a * curAlpha
	if f <= 0.002 {
		return
	}
	d.SetAlpha(f)
	fn()
	d.SetAlpha(curAlpha)
}

func rotateAbout(d *core.Drawing, cx, cy, ang float32, fn func()) {
	d.Save()
	d.Translate(cx, cy)
	d.Rotate(ang)
	d.Translate(-cx, -cy)
	fn()
	d.Restore()
}

func roundRectPath(d *core.Drawing, x, y, w, h, r float32) {
	d.BeginPath()
	d.RoundRect(x, y, w, h, r)
}

func circle(d *core.Drawing, cx, cy, r float32) {
	d.BeginPath()
	d.Arc(cx, cy, r, 0, float32(math.Pi*2), false)
	d.Fill()
}

func pill(d *core.Drawing, x, y, w, h float32, fill, stroke core.Color, sw float32, r float32) {
	d.SetFillColor(fill)
	rrect(d, x, y, w, h, r)
	if sw > 0 {
		d.SetStrokeColor(stroke)
		d.SetLineWidth(sw)
		strokeRRect(d, x, y, w, h, r)
	}
}

func arrow(d *core.Drawing, x1, y1, x2, y2, width float32, c core.Color) {
	d.SetStrokeColor(c)
	d.SetFillColor(c)
	d.SetLineWidth(width)
	d.SetLineCap(core.CapRound)
	d.BeginPath()
	d.MoveTo(x1, y1)
	d.LineTo(x2, y2)
	d.Stroke()
	ang := math.Atan2(float64(y2-y1), float64(x2-x1))
	hl := float64(width) * 3.4
	a1 := ang + math.Pi - 0.42
	a2 := ang + math.Pi + 0.42
	d.BeginPath()
	d.MoveTo(x2, y2)
	d.LineTo(float32(float64(x2)+hl*math.Cos(a1)), float32(float64(y2)+hl*math.Sin(a1)))
	d.LineTo(float32(float64(x2)+hl*math.Cos(a2)), float32(float64(y2)+hl*math.Sin(a2)))
	d.ClosePath()
	d.Fill()
}

func shadowOn(d *core.Drawing, a float32) {
	d.SetShadow(core.RGBA(0x1b1a17, a), 0, 16, 44)
}
func shadowOff(d *core.Drawing) {
	d.SetShadow(core.RGBA(0x000000, 0), 0, 0, 0)
}

// ---------------------------------------------------------------- chrome

func nav(d *core.Drawing, a float64) {
	if a <= 0.001 {
		return
	}
	al := float32(smooth(a))
	alpha(d, al, func() {
		// hairline
		d.SetStrokeColor(hair)
		d.SetLineWidth(1)
		d.BeginPath()
		d.MoveTo(0, 64)
		d.LineTo(float32(dw), 64)
		d.Stroke()
		// wordmark
		txt(d, "Atlas", 120, 45, 27, serif, core.FontWeightBold, ink)
		// right nav
		items := []string{"学科", "练习", "Premium", "教师"}
		x := float32(1200)
		for _, it := range items {
			txt(d, it, x, 41, 19, sans, core.FontWeightMedium, muted)
			x += measure(d, it, 19, sans, core.FontWeightMedium) + 46
		}
		// login pill
		pill(d, 1560, 16, 118, 34, cream, hair, 1, 17)
		d.SetTextAlign(core.TextAlignCenter)
		txt(d, "登录", 1619, 39, 18, sans, core.FontWeightMedium, ink)
		d.SetTextAlign(core.TextAlignStart)
	})
}

// ================================================================ S1 hero

func scene1(d *core.Drawing, t float64) {
	bg(d)
	nav(d, seg(t, 0, 0.8))
	p := seg(t, 0.15, 9.5)

	// eyebrow
	alpha(d, float32(easeOut(seg(t, 0.35, 0.95))), func() {
		txt(d, "ATLAS  ·  交互式 3D 教科书", 122, 300, 20, sans, core.FontWeightMedium, faint)
	})

	// headline, two lines rising out
	drawLine := func(s string, x, y float32, from, to float64) {
		e := easeOut(seg(t, from, to))
		alpha(d, float32(e), func() {
			dy := (1 - float32(e)) * 42
			txt(d, s, x, y+dy, 96, serif, core.FontWeightBold, ink)
		})
	}
	drawLine("拆开来，", 118, 460, 0.45, 1.15)
	drawLine("才学得透。", 118, 578, 0.62, 1.35)

	alpha(d, float32(easeOut(seg(t, 0.95, 1.6))), func() {
		txt(d, "互动   ·   学习   ·   练习", 124, 676, 30, sans, core.FontWeightMedium, muted)
	})

	// CTA pill
	alpha(d, float32(easeOut(seg(t, 1.2, 1.9))), func() {
		pill(d, 120, 736, 250, 66, ink, ink, 0, 33)
		d.SetTextAlign(core.TextAlignCenter)
		txt(d, "看一节课  →", 245, 778, 24, sans, core.FontWeightMedium, cream)
		d.SetTextAlign(core.TextAlignStart)
	})

	// right plot card, floating, tilted
	floatY := float32(math.Sin(t*0.55)*7) * 1.0
	cardX, cardY := float32(1210), float32(250)
	cw, chh := float32(600), float32(470)
	rotateAbout(d, cardX+cw/2, cardY+chh/2+floatY, float32(-6*math.Pi/180), func() {
		shadowOn(d, 0.16)
		d.SetFillColor(card_)
		rrect(d, cardX, cardY+floatY, cw, chh, 8)
		shadowOff(d)
		d.SetStrokeColor(hair)
		d.SetLineWidth(1)
		strokeRRect(d, cardX, cardY+floatY, cw, chh, 8)
		// curve draws on, ball rides
		cp := easeInOut(seg(t, 0.5, 3.0))
		ballU := easeInOut(seg(t, 0.7, 4.2))
		plot(d, cardX, cardY+floatY, cw, chh, float32(cp), float32(ballU), true, true)
	})

	alpha(d, float32(easeOut(seg(t, 1.6, 2.2))), func() {
		d.SetTextAlign(core.TextAlignCenter)
		txt(d, "物理 · 运动学 — 竖直上抛的小球", cardX+cw/2, cardY+chh+floatY+44, 18, sans, core.FontWeightMedium, faint)
		d.SetTextAlign(core.TextAlignStart)
	})
	_ = p
}

// plot draws a kinematics parabola card interior.
// cp: curve draw progress 0..1, u: ball position 0..1 along the curve.
func plot(d *core.Drawing, x, y, w, h float32, cp, u float32, vectors, labels bool) {
	ox := x + 78
	oy := y + h - 62
	W := w - 150
	H := h - 150
	if H < 40 {
		H = 40
	}
	// axes
	d.SetStrokeColor(core.RGB(0x9a958a))
	d.SetLineWidth(1.4)
	d.BeginPath()
	d.MoveTo(ox, oy+6)
	d.LineTo(ox, oy-H-26)
	d.MoveTo(ox-6, oy)
	d.LineTo(ox+W+18, oy)
	d.Stroke()
	// ticks + labels
	d.SetFillColor(faint)
	d.SetTextAlign(core.TextAlignEnd)
	d.SetTextBaseline(core.BaselineMiddle)
	for i := 0; i <= 5; i++ {
		yy := oy - H*float32(i)/5
		d.SetStrokeColor(core.RGB(0xe4dfd2))
		d.SetLineWidth(1)
		d.BeginPath()
		d.MoveTo(ox-5, yy)
		d.LineTo(ox+W+18, yy)
		d.Stroke()
		txt(d, itoa(i)+" m", ox-16, yy, 15, sans, core.FontWeightMedium, faint)
	}
	d.SetTextBaseline(core.BaselineAlphabetic)
	d.SetTextAlign(core.TextAlignStart)
	txt(d, "t (s)", ox+W+18, oy+30, 15, sans, core.FontWeightMedium, faint)
	txt(d, "位置 y (m)", ox+W-70, oy-H-38, 16, sans, core.FontWeightMedium, muted)

	// parabola, progressive
	parab := func(uu float32) (float32, float32) {
		return ox + uu*W, oy - 4*H*uu*(1-uu)
	}
	d.SetStrokeColor(red)
	d.SetLineWidth(3)
	d.SetLineCap(core.CapRound)
	d.BeginPath()
	steps := 90
	n := int(float32(steps) * cp)
	for i := 0; i <= n; i++ {
		uu := float32(i) / float32(steps)
		px, py := parab(uu)
		if i == 0 {
			d.MoveTo(px, py)
		} else {
			d.LineTo(px, py)
		}
	}
	d.Stroke()

	// ball + vectors
	if cp > 0.02 {
		bx, by := parab(u)
		// acceleration: straight down
		if vectors {
			arrow(d, bx, by+10, bx, by+72, 3, blueD)
			// velocity: tangent
			dx := W
			dy := -4 * H * (1 - 2*u)
			l := float32(math.Hypot(float64(dx), float64(dy)))
			vx := bx + dx/l*80
			vy := by + dy/l*80
			arrow(d, bx, by, vx, vy, 3, red)
		}
		shadowOn(d, 0.22)
		d.SetFillColor(card_)
		circle(d, bx, by, 12)
		shadowOff(d)
		d.SetStrokeColor(hair)
		d.SetLineWidth(1)
		d.BeginPath()
		d.Arc(bx, by, 12, 0, float32(math.Pi*2), false)
		d.Stroke()

		if labels {
			// small pills near vectors
			labelPill := func(s string, px, py float32, c core.Color) {
				wl := measure(d, s, 15, sans, core.FontWeightBold) + 22
				d.SetFillColor(c)
				rrect(d, px, py, wl, 26, 13)
				d.SetTextAlign(core.TextAlignCenter)
				d.SetTextBaseline(core.BaselineMiddle)
				txt(d, s, px+wl/2, py+13, 15, sans, core.FontWeightBold, core.RGB(0xffffff))
				d.SetTextAlign(core.TextAlignStart)
				d.SetTextBaseline(core.BaselineAlphabetic)
			}
			labelPill("v", bx+58, by-46, red)
			labelPill("a", bx+16, by+60, blueD)
		}
	}
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	neg := i < 0
	if neg {
		i = -i
	}
	var b [8]byte
	p := len(b)
	for i > 0 {
		p--
		b[p] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		p--
		b[p] = '-'
	}
	return string(b[p:])
}

// ================================================================ S2 morph

func scene2(d *core.Drawing, t float64) {
	bg(d)
	nav(d, 1)
	tl := t - 9.5
	p := seg(tl, 0.0, 6.5)
	_ = p

	// left copy
	alpha(d, float32(easeOut(seg(tl, 0.2, 0.9))), func() {
		txt(d, "抽象的概念，", 120, 470, 74, serif, core.FontWeightBold, ink)
	})
	alpha(d, float32(easeOut(seg(tl, 0.4, 1.1))), func() {
		txt(d, "变成能转动的", 120, 566, 74, serif, core.FontWeightBold, ink)
		txt(d, "立体图像。", 120, 662, 74, serif, core.FontWeightBold, ink)
	})
	alpha(d, float32(easeOut(seg(tl, 0.9, 1.6))), func() {
		txt(d, "拖动从任意角度看；点按某个部分，", 124, 742, 24, sans, core.FontWeightMedium, muted)
		txt(d, "查看它的拆解与实时数值。", 124, 782, 24, sans, core.FontWeightMedium, muted)
	})

	// right: flat page turning into a 3D card
	cx, cy := float32(1390), float32(500)
	cw, chh := float32(560), float32(400)
	turn := easeInOut(seg(tl, 0.6, 4.2)) // 0 flat -> 1 card
	scaleX := float32(1) - float32(math.Sin(float64(turn)*math.Pi))*0.55
	angle := float32(-24*math.Pi/180) * float32(turn)

	drawFlat := func(a float32) {
		shadowOn(d, 0.14)
		d.SetFillColor(card_)
		rrect(d, cx-cw/2, cy-chh/2, cw, chh, 8)
		shadowOff(d)
		d.SetStrokeColor(hair)
		d.SetLineWidth(1)
		strokeRRect(d, cx-cw/2, cy-chh/2, cw, chh, 8)
		// text lines
		d.SetFillColor(core.RGB(0xded8c9))
		for i := 0; i < 7; i++ {
			lw := float32(0.62+0.22*math.Sin(float64(i)*1.7))
			rect(d, cx-cw/2+48, cy-chh/2+70+float32(i)*38, cw*lw-96+120, 12)
		}
		d.SetFillColor(core.RGB(0xcfc8b6))
		rect(d, cx-cw/2+48, cy-chh/2+40, 150, 18)
	}

	d.Save()
	d.Translate(cx, cy)
	d.Scale(scaleX, 1)
	d.Rotate(angle)
	d.Translate(-cx, -cy)
	// flat page fades out
	alpha(d, 1-float32(turn), func() { drawFlat(1) })
	// 3D card fades in
	alpha(d, float32(turn), func() {
		shadowOn(d, 0.16)
		d.SetFillColor(card_)
		rrect(d, cx-cw/2, cy-chh/2, cw, chh, 10)
		shadowOff(d)
		d.SetStrokeColor(hair)
		d.SetLineWidth(1)
		strokeRRect(d, cx-cw/2, cy-chh/2, cw, chh, 10)
		plot(d, cx-cw/2, cy-chh/2, cw, chh, 1, 0.5, true, true)
	})
	d.Restore()
}

// ================================================================ S3 explode

func scene3(d *core.Drawing, t float64) {
	bg(d)
	nav(d, 1)
	tl := t - 16

	alpha(d, float32(easeOut(seg(tl, 0.2, 0.9))), func() {
		txt(d, "每一张图都能旋转，", 120, 470, 66, serif, core.FontWeightBold, ink)
		txt(d, "点开每个部分都能拆解。", 120, 560, 66, serif, core.FontWeightBold, ink)
	})
	alpha(d, float32(easeOut(seg(tl, 0.8, 1.5))), func() {
		txt(d, "拖动从任意角度看；点按某个部分，查看它的拆解与实时数值。", 124, 640, 24, sans, core.FontWeightMedium, muted)
		// chips
		pill(d, 124, 690, 150, 46, card_, hair, 1, 23)
		d.SetTextAlign(core.TextAlignCenter)
		txt(d, "⟳ 拖动旋转", 199, 719, 18, sans, core.FontWeightMedium, ink)
		pill(d, 290, 690, 150, 46, card_, hair, 1, 23)
		txt(d, "◉ 点按拆解", 365, 719, 18, sans, core.FontWeightMedium, ink)
		d.SetTextAlign(core.TextAlignStart)
	})

	cx, cy := float32(1330), float32(520)
	cw, chh := float32(620), float32(500)
	ang := float32(math.Sin(tl*0.5)*4) * math.Pi / 180
	rotateAbout(d, cx, cy, ang, func() {
		shadowOn(d, 0.16)
		d.SetFillColor(card_)
		rrect(d, cx-cw/2, cy-chh/2, cw, chh, 10)
		shadowOff(d)
		d.SetStrokeColor(hair)
		d.SetLineWidth(1)
		strokeRRect(d, cx-cw/2, cy-chh/2, cw, chh, 10)
		plot(d, cx-cw/2, cy-chh/2, cw, chh, 1, 0.5, true, false)
	})

	// callout pills appear sequentially
	callout := func(at float64, s string, px, py float32, tx, ty float32, c core.Color) {
		e := easeOut(seg(tl, at, at+0.6))
		if e <= 0.001 {
			return
		}
		alpha(d, float32(e), func() {
			// leader
			d.SetStrokeColor(c)
			d.SetLineWidth(2)
			d.SetLineDash()
			d.BeginPath()
			d.MoveTo(px, py)
			d.LineTo(tx, ty)
			d.Stroke()
			wl := measure(d, s, 19, sans, core.FontWeightBold) + 30
			pill(d, tx, ty-19, wl, 40, card_, hair, 1, 10)
			d.SetFillColor(c)
			circle(d, tx+16, ty+1, 5)
			txt(d, s, tx+30, ty+7, 19, sans, core.FontWeightBold, ink)
		})
	}
	// 位置 y (left of apex)
	callout(1.4, "位置 y  4.2 m", cx-80, cy-120, 800, 260, red)
	callout(2.0, "速度 v  12.6 m/s", cx+30, cy-200, 1370, 190, red)
	callout(2.6, "加速度 a  −9.8 m/s²", cx+70, cy+150, 1380, 830, blue)
}

// ================================================================ S4 practice

func scene4(d *core.Drawing, t float64) {
	bg(d)
	nav(d, 1)
	tl := t - 25

	alpha(d, float32(easeOut(seg(tl, 0.2, 0.9))), func() {
		txt(d, "每一页都有", 120, 450, 72, serif, core.FontWeightBold, ink)
		txt(d, "AP 风格的题目。", 120, 540, 72, serif, core.FontWeightBold, ink)
		txt(d, "选择题与数值题，边做边批改；", 124, 630, 24, sans, core.FontWeightMedium, muted)
		txt(d, "注册后可免费使用任意 3 套。", 124, 668, 24, sans, core.FontWeightMedium, muted)
	})

	// question card
	cx, cy := float32(1330), float32(510)
	cw, chh := float32(640), float32(560)
	e0 := float32(easeOut(seg(tl, 0.5, 1.4)))
	alpha(d, e0, func() {
		shadowOn(d, 0.16)
		d.SetFillColor(card_)
		rrect(d, cx-cw/2, cy-chh/2, cw, chh, 10)
		shadowOff(d)
		d.SetStrokeColor(hair)
		d.SetLineWidth(1)
		strokeRRect(d, cx-cw/2, cy-chh/2, cw, chh, 10)
		x0 := cx - cw/2 + 36
		y0 := cy - chh/2 + 40
		txt(d, "第 3 章 · 运动学  ·  题 4 / 8", x0, y0, 16, sans, core.FontWeightMedium, faint)
		pill(d, cx+cw/2-118, y0-24, 82, 28, grnBg, grn, 1, 14)
		d.SetTextAlign(core.TextAlignCenter)
		txt(d, "易错题", cx+cw/2-77, y0-3, 15, sans, core.FontWeightBold, grn)
		d.SetTextAlign(core.TextAlignStart)
		txt(d, "一小球以 10 m/s 竖直向上抛出，", x0, y0+64, 22, sans, core.FontWeightMedium, ink)
		txt(d, "忽略空气阻力，2 秒后它的速度约为？", x0, y0+100, 22, sans, core.FontWeightMedium, ink)

		// options
		opts := []struct {
			k, v string
			good bool
		}{{"A", "0 m/s", false}, {"B", "10 m/s", false}, {"C", "−10 m/s", true}, {"D", "20 m/s", false}}
		wrongAt, rightAt := 0.5, 0.72
		oy0 := y0 + 140
		for i, o := range opts {
			row := float32(i / 2)
			colX := x0 + float32(i%2)*270
			yy := oy0 + row*64
			ap := float32(easeOut(seg(tl, 1.0+float64(i)*0.1, 1.5+float64(i)*0.1)))
			alpha(d, ap, func() {
				// state highlight
				fill := card_
				stroke := hair
				if o.k == "B" && tl > wrongAt {
					fill = redBg
					stroke = red
				}
				if o.k == "C" && tl > rightAt {
					fill = grnBg
					stroke = grn
				}
				pill(d, colX, yy, 250, 50, fill, stroke, 1.4, 9)
				txt(d, o.k, colX+18, yy+33, 20, sans, core.FontWeightBold, muted)
				txt(d, o.v, colX+52, yy+33, 21, sans, core.FontWeightMedium, ink)
				if o.k == "B" && tl > wrongAt {
					d.SetStrokeColor(red)
					d.SetLineWidth(2.6)
					d.BeginPath()
					d.MoveTo(colX+216, yy+17)
					d.LineTo(colX+234, yy+35)
					d.MoveTo(colX+234, yy+17)
					d.LineTo(colX+216, yy+35)
					d.Stroke()
				}
				if o.k == "C" && tl > rightAt {
					d.SetStrokeColor(grn)
					d.SetLineWidth(3)
					d.SetLineCap(core.CapRound)
					d.BeginPath()
					d.MoveTo(colX+214, yy+26)
					d.LineTo(colX+222, yy+35)
					d.LineTo(colX+238, yy+15)
					d.Stroke()
				}
			})
		}
		// mini plot + legend
		miniY := oy0 + 2*64 + 20
		alpha(d, float32(easeOut(seg(tl, 1.9, 2.6))), func() {
			plotMini(d, x0, miniY, 330, 120)
			legendX := x0 + 380
			d.SetFillColor(red)
			circle(d, legendX, miniY+26, 5)
			txt(d, "你的答案", legendX+14, miniY+32, 15, sans, core.FontWeightMedium, muted)
			d.SetFillColor(grn)
			circle(d, legendX, miniY+56, 5)
			txt(d, "正确答案", legendX+14, miniY+62, 15, sans, core.FontWeightMedium, muted)
		})
	})
}

func plotMini(d *core.Drawing, x, y, w, h float32) {
	ox, oy := x+10, y+h-12
	W, H := w-30, h-30
	d.SetStrokeColor(core.RGB(0xb8b2a4))
	d.SetLineWidth(1.2)
	d.BeginPath()
	d.MoveTo(ox, oy)
	d.LineTo(ox+W, oy)
	d.Stroke()
	d.SetStrokeColor(red)
	d.SetLineWidth(2)
	d.BeginPath()
	for i := 0; i <= 60; i++ {
		u := float32(i) / 60
		px := ox + u*W
		py := oy - 4*H*u*(1-u)
		if i == 0 {
			d.MoveTo(px, py)
		} else {
			d.LineTo(px, py)
		}
	}
	d.Stroke()
	// your answer (wrong) at u=.3 low, correct at apex
	d.SetFillColor(red)
	circle(d, ox+W*0.28, oy-4*H*0.28*0.72-4, 6)
	d.SetFillColor(grn)
	circle(d, ox+W*0.5, oy-H*0.98, 6)
}

// ================================================================ S5 AI tutor

func scene5(d *core.Drawing, t float64) {
	bg(d)
	nav(d, 1)
	tl := t - 33

	alpha(d, float32(easeOut(seg(tl, 0.2, 0.9))), func() {
		txt(d, "AI 导师，针对那道题", 120, 450, 66, serif, core.FontWeightBold, ink)
		txt(d, "问一问。", 120, 540, 66, serif, core.FontWeightBold, ink)
		txt(d, "它能同时看到题目、你的答案和解题过程，", 124, 626, 24, sans, core.FontWeightMedium, muted)
		txt(d, "告诉你这一类题错在哪里。", 124, 664, 24, sans, core.FontWeightMedium, muted)
	})

	cx, cy := float32(1330), float32(500)
	cw, chh := float32(660), float32(600)
	alpha(d, float32(easeOut(seg(tl, 0.5, 1.3))), func() {
		shadowOn(d, 0.16)
		d.SetFillColor(card_)
		rrect(d, cx-cw/2, cy-chh/2, cw, chh, 14)
		shadowOff(d)
		d.SetStrokeColor(hair)
		d.SetLineWidth(1)
		strokeRRect(d, cx-cw/2, cy-chh/2, cw, chh, 14)
		x0 := cx - cw/2 + 34
		y0 := cy - chh/2 + 40
		txt(d, "第 3 章 · 运动学  ·  题 4 / 8", x0, y0, 16, sans, core.FontWeightMedium, faint)
		pill(d, cx+cw/2-130, y0-24, 96, 28, core.RGB(0xf1ece0), gold, 1, 14)
		d.SetTextAlign(core.TextAlignCenter)
		txt(d, "PREMIUM", cx+cw/2-82, y0-3, 15, sans, core.FontWeightBold, gold)

		// user bubble right
		alpha(d, float32(easeOut(seg(tl, 1.1, 1.6))), func() {
			bw := float32(430)
			bx := cx + cw/2 - 34 - bw
			pill(d, bx, y0+40, bw, 96, ink, ink, 0, 14)
			d.SetTextAlign(core.TextAlignStart)
			txt(d, "为什么 2 秒后速度是 −10 m/s？", bx+20, y0+84, 20, sans, core.FontWeightMedium, cream)
			txt(d, "我算出来是 0 m/s。", bx+20, y0+114, 20, sans, core.FontWeightMedium, cream)
		})

		// AI bubble left
		alpha(d, float32(easeOut(seg(tl, 1.5, 2.1))), func() {
			bx := x0
			by := y0 + 170
			pill(d, bx, by, 470, 210, core.RGB(0xf2efe6), hair, 1, 14)
			reveal := func(line string, idx int) {
				e := seg(tl, 2.0+float64(idx)*0.55, 2.5+float64(idx)*0.55)
				if e <= 0.001 {
					return
				}
				alpha(d, float32(e), func() {
					txt(d, line, bx+22, by+48+float32(idx)*52, 20, sans, core.FontWeightMedium, ink)
				})
			}
			reveal("因为重力加速度 g ≈ −9.8 m/s²，", 0)
			reveal("v = v₀ + gt = 10 − 19.6 ≈ −9.6 m/s，", 1)
			reveal("负号表示速度方向向下。", 2)
		})
	})
}

// ================================================================ S6 stats

func scene6(d *core.Drawing, t float64) {
	bg(d)
	nav(d, 1)
	tl := t - 42

	alpha(d, float32(easeOut(seg(tl, 0.2, 0.8))), func() {
		d.SetTextAlign(core.TextAlignCenter)
		txt(d, "覆盖范围", 960, 250, 20, sans, core.FontWeightMedium, faint)
	})
	// three numbers
	cnt := easeOut(seg(tl, 0.4, 2.2))
	stats := []struct {
		n     int
		label string
		cx    float32
	}{{1696, "道练习题", 560}, {265, "套练习", 960}, {51, "个章节 · 5 个学科", 1360}}
	for i, st := range stats {
		v := int(float64(st.n) * cnt)
		alpha(d, float32(easeOut(seg(tl, 0.4+float64(i)*0.12, 1.0+float64(i)*0.12))), func() {
			d.SetTextAlign(core.TextAlignCenter)
			txt(d, comma(v), st.cx, 400, 96, serif, core.FontWeightBold, ink)
			txt(d, st.label, st.cx, 460, 20, sans, core.FontWeightMedium, faint)
		})
	}
	// subject chips
	subs := []string{"AP Physics 1 & C", "Calculus AB & BC", "Biology", "Chemistry", "Statistics"}
	alpha(d, float32(easeOut(seg(tl, 1.6, 2.4))), func() {
		totalW := float32(0)
		widths := make([]float32, len(subs))
		for i, s := range subs {
			widths[i] = measure(d, s, 18, sans, core.FontWeightMedium) + 40
			totalW += widths[i]
		}
		totalW += float32(len(subs)-1) * 16
		x := 960 - totalW/2
		for i, s := range subs {
			pill(d, x, 520, widths[i], 44, card_, hair, 1, 22)
			d.SetTextAlign(core.TextAlignCenter)
			txt(d, s, x+widths[i]/2, 548, 18, sans, core.FontWeightMedium, muted)
			x += widths[i] + 16
		}
	})
	// headline
	alpha(d, float32(easeOut(seg(tl, 2.4, 3.2))), func() {
		d.SetTextAlign(core.TextAlignCenter)
		txt(d, "从一节课开始。", 960, 760, 88, serif, core.FontWeightBold, ink)
	})
	// CTA
	alpha(d, float32(easeOut(seg(tl, 3.2, 3.9))), func() {
		pill(d, 960-125, 820, 250, 64, ink, ink, 0, 32)
		d.SetTextAlign(core.TextAlignCenter)
		txt(d, "看一节课  →", 960, 861, 24, sans, core.FontWeightMedium, cream)
		txt(d, "atlasclass.app", 960, 950, 20, sans, core.FontWeightMedium, faint)
	})
}

func comma(n int) string {
	if n < 1000 {
		return itoa(n)
	}
	// one comma is enough here
	hi := n / 1000
	lo := n % 1000
	ls := itoa(lo)
	for len(ls) < 3 {
		ls = "0" + ls
	}
	return itoa(hi) + "," + ls
}
