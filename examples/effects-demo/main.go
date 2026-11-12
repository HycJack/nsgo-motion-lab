// 《组合技画廊》—— 现有原语的组合效果演示 + 转场三件套实战。
//
// 七个组合技，每幕一个（幕间转场就是 Wipe/Slide/Iris/Cut/FadeWhite 实战）：
//  1. Dolly Zoom：eye 推近 + Lens 反向补偿，中环大小恒定、外环拉伸
//  2. 签名描画：EqWrite 沿弧长写五角星 + Dot 笔尖跟随
//  3. 打字机代码：mono 逐字符 reveal + 光标闪烁
//  4. 星空 + 流星：Hash 确定性撒点 + KF 流星轨迹
//  5. 进度环：Circle 弧进度 + 数字计数
//  6. 画中画：三个 Canvas 各画各的（timeline 卡片同款模式）
//
// 全部由现有原语组合，零新代码。
package main

import (
	"fmt"
	"math"

	"nsgofilm/film"
	"nsgofilm/filmkit"
	"nsgofilm/mathkit"
	"nsgofilm/tpl"

	"github.com/oneo-me/nsgo/core"
)

const (
	dw    float32 = 1920.0
	dh    float32 = 1080.0
	fade          = 0.4
)

var (
	ink   = core.RGB(0xF1F5F9)
	cyan  = core.RGB(0x06B6D4)
	blue  = core.RGB(0x3B82F6)
	violet = core.RGB(0xA78BFA)
	green = core.RGB(0x34D399)
	pink  = core.RGB(0xF472B6)
)

func main() {
	film.Run(film.Film{
		W:    dw,
		H:    dh,
		Fade: fade,
		Bg:   tpl.DarkMode.Bg,
		Shots: film.Series(fade,
			film.Shot{D: 5, Fn: tpl.Title(tpl.DarkMode,
				"组合技画廊", "现有原语能出的效果 · 每幕之间就是新转场")},
			film.Shot{D: 6.5, Fn: sceneDollyZoom, In: film.Slide},
			film.Shot{D: 6, Fn: sceneSignature, In: film.Wipe},
			film.Shot{D: 6, Fn: sceneTypewriter, In: film.Cut},
			film.Shot{D: 6, Fn: sceneStarfield, In: film.Iris},
			film.Shot{D: 5.5, Fn: sceneRing, In: film.Wipe},
			film.Shot{D: 6, Fn: scenePiP, In: film.Slide},
			film.Shot{D: 4.5, Fn: tpl.Outro(tpl.DarkMode,
				"原语 × 组合 = 视频", "film-framework · 全部零新代码"),
				In: film.FadeWhite, Out: film.FadeBlack,
			},
		),
	})
}

// ---------------------------------------------------------- 1. Dolly Zoom

// 透视走廊：四圈嵌套，相机推近 + Lens 反向拉宽——中环大小恒定。
func sceneDollyZoom(d *core.Drawing, f *film.Frame) {
	p := filmkit.Seg(f.LocalT, 0.8, 5.8) // 推进进度
	eyeX := 12.0 - 6.5*p
	distMid := eyeX - 4.25 // 到中环平面的距离
	lens := distMid / 6.75 // 起始时 = 1：中环大小恒定，外环被拉扁

	a := &mathkit.Axes3D{
		X: [2]float64{-1, 10}, Y: [2]float64{-1, 4.5}, Z: [2]float64{-2, 5},
		At: [2]float32{360, 140}, Size: [2]float32{1200, 800},
		Eye:  [3]float64{eyeX, 1.8, 1.5},
		Look: [3]float64{4.25, 1.8, 1.5},
		Lens: lens,
	}
	cols := []core.Color{cyan, violet, green, pink}
	for i, x := range []float64{1.5, 3.5, 5.5, 7.5} {
		ring := make([][3]float64, 0, 49)
		for k := 0; k <= 48; k++ {
			th := 2 * math.Pi * float64(k) / 48
			ring = append(ring, [3]float64{x, 1.8 + 1.2*math.Cos(th), 1.5 + 1.2*math.Sin(th)})
		}
		a.Poly3(d, ring, 1, cols[i], film.A(0.9), 4)
	}
	core2Text(d, "Dolly Zoom · 推近 + 变焦补偿", dw/2, dh-140, 44,
		tpl.DarkMode.Heading, film.A(filmkit.F32(filmkit.Seg(f.LocalT, 0.3, 0.8))), 0.5)
}

// ---------------------------------------------------------- 2. 签名描画

// starEq 是手搓的五角星轮廓（演示：任意矢量都能走 EqWrite 管线）。
var starEq = func() mathkit.Eq {
	outer, inner := 170.0, 66.0
	var pts []float32
	for k := 0; k < 10; k++ {
		ang := -math.Pi/2 + float64(k)*math.Pi/5
		r := outer
		if k%2 == 1 {
			r = inner
		}
		pts = append(pts, float32(190+r*math.Cos(ang)), float32(190+r*math.Sin(ang)))
	}
	pts = append(pts, pts[0], pts[1]) // 闭合
	return mathkit.Eq{W: 380, H: 380, Lines: []mathkit.Poly{{Pts: pts, Closed: true}}}
}()

// polyPointAt 返回折线上弧长 l 处的点（笔尖位置）。
func polyPointAt(p mathkit.Poly, l float64) (float32, float32) {
	n := len(p.Pts) / 2
	if n == 0 {
		return 0, 0
	}
	var acc float64
	for i := 1; i < n; i++ {
		dx := float64(p.X(i) - p.X(i-1))
		dy := float64(p.Y(i) - p.Y(i-1))
		seg := math.Sqrt(dx*dx + dy*dy)
		if acc+seg >= l {
			t := 0.0
			if seg > 0 {
				t = (l - acc) / seg
			}
			return p.X(i-1) + float32(dx*t), p.Y(i-1) + float32(dy*t)
		}
		acc += seg
	}
	return p.X(n-1), p.Y(n-1)
}

// sceneSignature：五角星逐笔描 + 笔尖圆点跟着走（写完弹出 "✓"）。
func sceneSignature(d *core.Drawing, f *film.Frame) {
	t := f.LocalT
	cx, cy := dw/2-190, dh/2-190
	p := filmkit.Clamp01(filmkit.Seg(t, 0.8, 3.4))

	if p > 0 {
		mathkit.EqWrite(d, starEq, cx, cy, 1, p, cyan, film.A(1), 7)
		// 笔尖：已写弧长在折线上的位置。
		ln := starEq.Lines[0]
		px, py := polyPointAt(ln, float64(p)*ln.Length())
		mathkit.Dot(d, cx+px, cy+py, 12, goldDot(), film.A(1))
	}
	if p >= 1 {
		if fade := film.A(filmkit.F32(filmkit.Seg(t, 3.6, 4.2))); fade > 0 {
			core2Text(d, "✓ 任意 SVG 轮廓都能这样描", dw/2, cy+440, 46,
				tpl.DarkMode.Body, fade, 0.5)
		}
	}
	core2Text(d, "签名描画 · EqWrite + 笔尖", dw/2, 120, 44,
		tpl.DarkMode.Heading, film.A(filmkit.F32(filmkit.Seg(t, 0.2, 0.7))), 0.5)
}

func goldDot() core.Color { return core.RGB(0xFFD166) }

// ---------------------------------------------------------- 3. 打字机

var codeLines = []string{
	"func main() {",
	`  film.Run(film.Film{`,
	`    Shots: film.Series(0.5,`,
	`      film.Shot{D: 5, Fn: title}),`,
	`  })`,
	"}",
}

// sceneTypewriter：逐字符 reveal + 光标闪烁。
func sceneTypewriter(d *core.Drawing, f *film.Frame) {
	t := f.LocalT
	x0, y0 := dw*0.22, dh*0.3
	lh := float32(72)

	shown := filmkit.Clamp01(filmkit.Seg(t, 0.6, 4.6))
	totalChars := 0
	for _, l := range codeLines {
		totalChars += len(l)
	}
	nChars := int(float64(totalChars) * shown)

	used := 0
	for i, l := range codeLines {
		if nChars <= used {
			break
		}
		show := l
		if nChars-used < len(l) {
			show = l[:nChars-used]
		}
		core2Text(d, show, x0, y0+float32(i)*lh, 44, green, film.A(1), 0)
		used += len(l)
	}

	// 行号。
	for i := range codeLines {
		core2Text(d, fmt.Sprint(i+1), x0-60, y0+float32(i)*lh, 32,
			tpl.DarkMode.Muted, film.A(0.6), 1)
	}

	// 光标：定位到当前行末，闪烁。
	blink := 0.5 + 0.5*math.Sin(float64(t)*6)
	if shown < 1 && blink > 0.5 {
		line := 0
		left := nChars
		for i, l := range codeLines {
			if left <= len(l) {
				line = i
				break
			}
			left -= len(l)
		}
		w := float32(left) * 24
		d.SetFillColor(ink)
		d.SetAlpha(film.A(0.8))
		d.BeginPath()
		d.Rect(x0+w+6, y0+float32(line)*lh-36, 20, 44)
		d.Fill()
	}
	core2Text(d, "打字机 · KF 逐字符 + 光标", dw/2, 120, 44,
		tpl.DarkMode.Heading, film.A(filmkit.F32(filmkit.Seg(t, 0.2, 0.7))), 0.5)
}

// ---------------------------------------------------------- 4. 星空流星

// sceneStarfield：Hash 确定性撒点闪烁 + 一颗 KF 流星拖尾。
func sceneStarfield(d *core.Drawing, f *film.Frame) {
	t := f.LocalT
	for i := 0; i < 90; i++ {
		x := float32(film.Hash(i)) * dw
		y := float32(film.Hash(i + 100)) * dh * 0.92
		r := float32(1 + 2.4*film.Hash(i+200))
		a := 0.25 + 0.75*math.Abs(math.Sin(2*math.Pi*(t/3+film.Hash(i+300))))
		mathkit.Dot(d, x, y, r, ink, film.A(float32(a)))
	}
	// 流星：2.2-3.6s 从左上划到右下，尾部 10 个渐隐点。
	mp := filmkit.Seg(t, 2.2, 3.6)
	if mp > 0 {
		x0, y0 := dw*0.12, dh*0.14
		x1, y1 := dw*0.78, dh*0.55
		for k := 0; k < 10; k++ {
			tt := clamp01f(float64(mp) - float64(k)*0.03)
			if tt <= 0 {
				continue
			}
			x := filmkit.Lerp32(x0, x1, float32(tt))
			y := filmkit.Lerp32(y0, y1, float32(tt))
			mathkit.Dot(d, x, y, float32(9) - float32(k)*0.8, core.RGB(0xFFD166), film.A(float32(1-tt*0.7)))
		}
	}
	core2Text(d, "星空 · Hash 撒点 + KF 流星", dw/2, 120, 44,
		tpl.DarkMode.Heading, film.A(filmkit.F32(filmkit.Seg(t, 0.2, 0.7))), 0.5)
}

func clamp01f(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}

// ---------------------------------------------------------- 5. 进度环

// sceneRing：Circle 弧进度 + 百分比计数。
func sceneRing(d *core.Drawing, f *film.Frame) {
	t := f.LocalT
	p := filmkit.Seg(t, 0.8, 3.6)
	cx, cy, r := dw/2, dh/2, float32(190)

	// 底环（弱）+ 进度弧（粗）。
	mathkit.Circle(d, cx, cy, r, 1, tpl.DarkMode.Card, film.A(1), 34)
	if p > 0 {
		mathkit.Circle(d, cx, cy, r, float64(p), cyan, film.A(1), 34)
	}
	pct := int(math.Round(87 * p))
	core2Text(d, fmt.Sprintf("%d%%", pct), cx, cy+22, 130,
		tpl.DarkMode.Heading, film.A(1), 0.5)
	core2Text(d, "构建中 · Circle 弧进度", dw/2, 120, 44,
		tpl.DarkMode.Heading, film.A(filmkit.F32(filmkit.Seg(t, 0.2, 0.7))), 0.5)
}

// ---------------------------------------------------------- 6. 画中画

// scenePiP：三个 Canvas 各画各的——同一页面的多个独立时钟画面。
func scenePiP(d *core.Drawing, f *film.Frame) {
	t := f.LocalT
	// 背景：三块卡位。
	cards := [][2]float32{{190, 200}, {760, 200}, {1330, 200}}
	for _, c := range cards {
		d.SetFillColor(tpl.DarkMode.Card)
		d.SetAlpha(film.A(1))
		d.BeginPath()
		d.Rect(c[0]-16, c[1]-16, 460+32, 460+32)
		d.Fill()
	}
	// 直接在主画布上画三块（Canvas 嵌套需要 Scope，演示用矩形分区代替）。
	panes := []struct {
		x, y float32
		draw func(d *core.Drawing, x, y float32, t float64)
		ttl  string
	}{
		{cards[0][0], cards[0][1], paneWave, "正弦"},
		{cards[1][0], cards[1][1], paneBars, "脉冲"},
		{cards[2][0], cards[2][1], paneOrbitDot, "轨道点"},
	}
	for _, pn := range panes {
		pn.draw(d, pn.x, pn.y, t)
		core2Text(d, pn.ttl, pn.x+230, pn.y+510, 36, tpl.DarkMode.Muted, film.A(0.8), 0.5)
	}
	core2Text(d, "画中画 · 多画面各画各的", dw/2, 120, 44,
		tpl.DarkMode.Heading, film.A(filmkit.F32(filmkit.Seg(t, 0.2, 0.7))), 0.5)
}

func paneWave(d *core.Drawing, x, y float32, t float64) {
	d.SetStrokeColor(cyan)
	d.SetAlpha(film.A(0.9))
	d.SetLineWidth(4)
	d.BeginPath()
	for i := 0; i <= 60; i++ {
		fx := float32(i) / 60
		fy := float32(0.5 + 0.32*math.Sin(float64(fx)*6 + t))
		px, py := x+fx*440, y+fy*440
		if i == 0 {
			d.MoveTo(px, py)
		} else {
			d.LineTo(px, py)
		}
	}
	d.Stroke()
	d.SetAlpha(1)
}

func paneBars(d *core.Drawing, x, y float32, t float64) {
	for i := 0; i < 5; i++ {
		h := float32(60 + 200*math.Abs(math.Sin(t*1.3+float64(i))))
		d.SetFillColor(violet)
		d.SetAlpha(film.A(0.85))
		d.BeginPath()
		d.Rect(x+50+float32(i)*80, y+420-h, 52, h)
		d.Fill()
	}
	d.SetAlpha(1)
}

func paneOrbitDot(d *core.Drawing, x, y float32, t float64) {
	cx, cy := x+230, y+230
	mathkit.Circle(d, cx, cy, 150, 1, tpl.DarkMode.Muted, film.A(0.5), 2)
	ang := t * 1.8
	mathkit.Dot(d, cx+float32(150*math.Cos(ang)), cy+float32(150*math.Sin(ang)), 14, green, film.A(1))
}

// core2Text 是左/中/右对齐文字（align 0 左 0.5 中 1 右）。
func core2Text(d *core.Drawing, s string, x, y, size float32, c core.Color, alpha float32, align float32) {
	if alpha <= 0.001 || s == "" {
		return
	}
	w := d.MeasureText(s, core.Font{Size: core.Points(size)})
	d.SetFillColor(c)
	d.SetAlpha(alpha)
	d.FillText(s, x-align*w, y, core.Font{Size: core.Points(size)})
}
