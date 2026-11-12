// 《甲烷的燃烧》—— 球棍模型分子动画教学片。
//
// CH₄ + 2O₂ →(点燃) CO₂ + 2H₂O
//
// 反应的本质是**原子的重新组合**：全片只有 9 个原子球（1C + 4H + 4O），
// 反应幕里每个原子从反应物位置沿弧线飞向生成物位置，化学键随之断开
// 与生成——没有贴图、没有粒子特效，全部是投影 + 渐变球 + 深度排序。
//
// 五幕（Series，一台相机）：
//  1. 标题 + 配平的化学方程式
//  2. 反应物亮相：甲烷四面体旋转、两分子氧气就位
//  3. 点燃：原子飞行重排（键断开/生成）
//  4. 生成物：CO₂ 与两分子水
//  5. 落版：质量守恒
package main

import (
	"math"
	"sort"

	"nsgofilm/film"
	"nsgofilm/filmkit"
	"nsgofilm/mathkit"

	"github.com/oneo-me/nsgo/core"
)

const (
	dw   float32 = 1920.0
	dh   float32 = 1080.0
	fade         = 0.5

	bl = 1.05 // C–H 键长（世界单位）
)

var (
	ink   = core.RGB(0xF1F5F9)
	gold  = core.RGB(0xFFC53D)
	pale  = core.RGB(0x94A3B8)
	paper = core.RGB(0x0B0F16)
)

func main() {
	film.Run(film.Film{
		W:    dw,
		H:    dh,
		Fade: fade,
		Bg:   paper,
		Shots: film.Series(fade,
			film.Shot{D: 5, Fn: sceneTitle},
			film.Shot{D: 7, Fn: sceneReactants},
			film.Shot{D: 8, Fn: sceneReaction},
			film.Shot{D: 6.5, Fn: sceneProducts},
			film.Shot{D: 5, Fn: sceneConservation, Out: film.FadeBlack},
		),
		Subs: []film.Sub{
			{A: 1.0, B: 4.4, Text: "甲烷是最简单的有机物——一个碳，四个氢。"},
			{A: 5.4, B: 11.2, Text: "点燃甲烷和氧气的混合气体，剧烈的氧化反应开始了。"},
			{A: 11.6, B: 14.6, Text: "化学键断裂，原子重新组合。"},
			{A: 15.0, B: 20.2, Text: "碳和氧结合成二氧化碳，氢和氧结合成水。"},
			{A: 20.8, B: 24.2, Text: "反应前后，原子的种类和个数都不变——质量守恒。"},
		},
		SubStyle: film.DefaultSubStyle(ink, gold, core.RGB(0x000000)),
	})
}

// ---------------------------------------------------------------- 原子表

// atomDef 描述一个原子：反应物位置（rx..）、生成物位置（px..）、
// 半径、元素、起飞延迟（反应波）。
type atomDef struct {
	rx, ry, rz float64
	px, py, pz float64
	r          float64
	el         string
	delay      float64
	// methH：甲烷的氢——反应物位置绕碳旋转（旋转由场景传 θ）。
	methH bool
	off   [3]float64
}

// 四面体方向：立方体四角归一化。
var tetra = [][3]float64{
	{1, 1, 1}, {1, -1, -1}, {-1, 1, -1}, {-1, -1, 1},
}

const (
	methX = -4.6 // 甲烷碳的初始位置
	methY = 2.0
)

// atoms 是全片 9 个原子的编排表（反应物左，生成物右）。
func atomTable() []atomDef {
	is := 1.0 / math.Sqrt(3)
	ats := []atomDef{
		{rx: methX, ry: methY, rz: 0, px: 3.6, py: 3.0, pz: 0, r: 0.40, el: "C", delay: 0.0},
	}
	for i, d := range tetra {
		ats = append(ats, atomDef{
			rx: methX + bl*d[0]*is, ry: methY + bl*d[1]*is, rz: bl * d[2] * is,
			r: 0.24, el: "H",
			delay: 0.28 + 0.09*float64(i), methH: true,
			off: [3]float64{bl * d[0] * is, bl * d[1] * is, bl * d[2] * is},
		})
	}
	ats = append(ats,
		atomDef{rx: -1.6, ry: 3.15, rz: 0, px: 2.3, py: 3.0, pz: 0, r: 0.36, el: "O", delay: 0.06},
		atomDef{rx: -0.5, ry: 3.15, rz: 0, px: 4.9, py: 3.0, pz: 0, r: 0.36, el: "O", delay: 0.16},
		atomDef{rx: -1.6, ry: 0.85, rz: 0, px: 2.75, py: 0.95, pz: 0, r: 0.36, el: "O", delay: 0.26},
		atomDef{rx: -0.5, ry: 0.85, rz: 0, px: 5.35, py: 0.85, pz: 0, r: 0.36, el: "O", delay: 0.36},
	)
	// 甲烷氢的生成物位置：两分子水（O 已在表里，H 挂上去）。
	homes := [][3]float64{
		{2.75 - 0.52, 0.95 + 0.68, 0}, {2.75 + 0.52, 0.95 + 0.68, 0},
		{5.35 - 0.52, 0.85 + 0.68, 0}, {5.35 + 0.52, 0.85 + 0.68, 0},
	}
	hi := 0
	for i := range ats {
		if ats[i].el == "H" {
			ats[i].px, ats[i].py, ats[i].pz = homes[hi][0], homes[hi][1], homes[hi][2]
			hi++
		}
	}
	return ats
}

// 化学键（原子下标对）。product 标记生成物键（反应后才出现）。
var bonds = []struct {
	a, b    int
	order   int // 1 单键 2 双键
	product bool
}{
	{0, 1, 1, false}, {0, 2, 1, false}, {0, 3, 1, false}, {0, 4, 1, false}, // C–H ×4
	{5, 6, 2, false}, {7, 8, 2, false}, // O=O ×2
	{0, 5, 2, true}, {0, 6, 2, true}, // C=O ×2（CO₂）
	{7, 1, 1, true}, {7, 2, 1, true}, // H–O–H（水一）
	{8, 3, 1, true}, {8, 4, 1, true}, // H–O–H（水二）
}

// ---------------------------------------------------------------- 绘制

var elColors = map[string]struct{ base, hi core.Color }{
	"C": {core.RGB(0x39404A), core.RGB(0x9AA3AE)},
	"H": {core.RGB(0xE9ECEF), core.RGB(0xFFFFFF)},
	"O": {core.RGB(0xD64545), core.RGB(0xFF9B8A)},
}

// ball 画一个有高光的渐变球。
func ball(d *core.Drawing, sx, sy, r float32, el string, alpha float32) {
	cc := elColors[el]
	g := core.CanvasRadial(sx-r*0.35, sy-r*0.4, r*0.08, sx, sy, r,
		core.GradientStop{Offset: 0, Color: cc.hi.Faded(alpha)},
		core.GradientStop{Offset: 0.55, Color: cc.base.Faded(alpha)},
		core.GradientStop{Offset: 1, Color: cc.base.Faded(alpha * 0.55)},
	)
	d.SetFillGradient(g)
	d.BeginPath()
	d.Arc(sx, sy, r, 0, float32(2*math.Pi), false)
	d.Fill()
}

// stick 画一根化学键（粗线），宽度随深度。
func stick(d *core.Drawing, ax, ay, bx, by, depth, w float32, c core.Color, alpha float32) {
	d.SetStrokeColor(c)
	d.SetAlpha(alpha)
	d.SetLineWidth(w * float32(math.Max(0.55, math.Min(1.6, 10.0/float64(depth)))))
	d.SetLineCap(core.CapRound)
	d.BeginPath()
	d.MoveTo(ax, ay)
	d.LineTo(bx, by)
	d.Stroke()
	d.SetAlpha(1)
}

// chemSeg 是化学式的一段：sub 为下标。
type chemSeg struct {
	t   string
	sub bool
}

// chemText 画带下标的化学式（align: 0 左 0.5 中 1 右）。
func chemText(d *core.Drawing, segs []chemSeg, x, y, size float32, c core.Color, alpha float32, align float32) {
	total := float32(0)
	widths := make([]float32, len(segs))
	for i, s := range segs {
		fs := size
		if s.sub {
			fs = size * 0.6
		}
		widths[i] = d.MeasureText(s.t, core.Font{Size: core.Points(fs)})
		total += widths[i]
	}
	cur := x - align*total
	for i, s := range segs {
		fs := size
		base := y
		if s.sub {
			fs = size * 0.6
			base = y + size*0.22
		}
		d.SetFillColor(c)
		d.SetAlpha(alpha)
		d.FillText(s.t, cur, base, core.Font{Size: core.Points(fs)})
		cur += widths[i]
	}
}

// equation 是配平的化学方程式。
var equation = []chemSeg{
	{t: "CH"}, {t: "4", sub: true},
	{t: " + 2O"}, {t: "2", sub: true},
	{t: "  →  CO"}, {t: "2", sub: true},
	{t: " + 2H"}, {t: "2", sub: true}, {t: "O"},
}

func drawEqCenter(d *core.Drawing, y float32, size float32, c core.Color, alpha float32) {
	chemText(d, equation, dw/2, y, size, c, alpha, 0.5)
}

// ---------------------------------------------------------------- 幂与相机

// camPose：固定机位 + 极缓慢漂移，后段随生成物右移视线。
func camPose(t float64) ([3]float64, [3]float64) {
	drift := 0.25 * math.Sin(t/6)
	look := [3]float64{0.3 + 2.0*filmkit.Seg(t, 14.5, 19.5), 1.7, 0}
	return [3]float64{0.4 + drift, 2.0, 9.8}, look
}

func chemAxes(eye, look [3]float64) *mathkit.Axes3D {
	return &mathkit.Axes3D{
		X: [2]float64{-6.5, 7.5}, Y: [2]float64{-0.4, 4.8}, Z: [2]float64{-2.2, 2.2},
		At: [2]float32{280, 100}, Size: [2]float32{1360, 860},
		Eye:  eye,
		Look: look,
	}
}

// riseT：淡入 + 上浮（模板同款入场）。
func riseT(t, at, dur float64) (float64, float32) {
	p := filmkit.Clamp01((t - at) / dur)
	e := filmkit.EaseOut(p)
	return e, float32((1 - e) * 30)
}

// ---------------------------------------------------------------- 幕

func sceneTitle(d *core.Drawing, f *film.Frame) {
	e, dy := riseT(f.LocalT, 0.2, 0.9)
	coreText(d, "甲烷的燃烧", dw/2, dh/2-120+dy, 120, ink, float32(e), 0.5)
	barP := filmkit.Seg(f.LocalT, 0.9, 1.5)
	if barP > 0 {
		w := 200.0 * float32(barP)
		d.SetFillColor(gold)
		d.SetAlpha(float32(e))
		d.BeginPath()
		d.Rect(dw/2-w/2, dh/2-84, w, 8)
		d.Fill()
		d.SetAlpha(1)
	}
	e2, dy2 := riseT(f.LocalT, 1.2, 0.9)
	drawEqCenter(d, dh/2+40+dy2, 72, pale, float32(e2))
	e3, dy3 := riseT(f.LocalT, 1.6, 0.9)
	coreText(d, "—— 用球棍模型看化学反应 ——", dw/2, dh/2+170+dy3, 38, pale, float32(e3), 0.5)
}

// sceneReactants：甲烷四面体旋转，氧气双分子就位，标签浮现。
func sceneReactants(d *core.Drawing, f *film.Frame) {
	eye, look := camPose(f.T)
	a := chemAxes(eye, look)
	drawMolecules(d, a, f.LocalT, 0, film.A(1))
	labelIn(d, a, f.LocalT, "甲烷 CH₄", methX, 3.4, 2.2)
	labelIn(d, a, f.LocalT, "氧气 O₂", -1.05, 4.2, 2.6)
	labelIn(d, a, f.LocalT, "氧气 O₂", -1.05, 1.7, 2.6)
}

// sceneReaction：原子飞行重排，"点燃" 点题。
func sceneReaction(d *core.Drawing, f *film.Frame) {
	eye, look := camPose(f.T)
	a := chemAxes(eye, look)
	K := filmkit.Seg(f.LocalT, 2.0, 6.4)
	drawMolecules(d, a, f.LocalT, float64(K), film.A(1))
	if K > 0.02 && K < 0.92 {
		coreText(d, "点 燃", dw/2, 170, 56, gold, float32(filmkit.Seg(f.LocalT, 2.0, 2.5)), 0.5)
	}
}

// sceneProducts：生成物列队 + 标注。
func sceneProducts(d *core.Drawing, f *film.Frame) {
	eye, look := camPose(f.T)
	a := chemAxes(eye, look)
	drawMolecules(d, a, f.LocalT, 1, film.A(1))
	labelIn(d, a, f.LocalT, "二氧化碳 CO₂", 3.6, 4.35, 1.0)
	labelIn(d, a, f.LocalT, "水 H₂O ×2", 4.05, 2.3, 1.0)
}

// sceneConservation：方程式大字落版 + 质量守恒。
func sceneConservation(d *core.Drawing, f *film.Frame) {
	e, dy := riseT(f.LocalT, 0.2, 1.0)
	drawEqCenter(d, dh/2-60+dy, 96, ink, float32(e))
	e2, dy2 := riseT(f.LocalT, 1.0, 1.0)
	coreText(d, "原子的种类和个数不变 —— 质量守恒", dw/2, dh/2+80+dy2, 52, gold, float32(e2), 0.5)
}

// labelIn 用本幕的投影在世界上方写字（billboard）。
func labelIn(d *core.Drawing, a *mathkit.Axes3D, t float64, s string, x, y float64, at float64) {
	e, dy := riseT(t, at, 0.7)
	sx, sy, _, ok := a.Project(x, y, 0)
	if !ok {
		return
	}
	coreText(d, s, sx, sy+float32(dy), 40, pale, float32(e), 0.5)
}

// drawMolecules 画出全部原子与键：K=0 反应物，K=1 生成物。
func drawMolecules(d *core.Drawing, a *mathkit.Axes3D, t float64, K float64, alpha float32) {
	ats := atomTable()
	theta := 0.55 * t // 甲烷绕竖轴的旋转角

	// 每个原子的当前位置与飞行进度。
	cur := make([][3]float64, len(ats))
	ks := make([]float64, len(ats))
	for i, at := range ats {
		ki := filmkit.Clamp01((K - at.delay) / 0.5)
		e := filmkit.EaseInOut(ki)
		rx, ry, rz := at.rx, at.ry, at.rz
		if at.methH {
			// 甲烷氢：绕碳的竖轴转。
			co, so := math.Cos(theta), math.Sin(theta)
			ox := at.off[0]*co + at.off[2]*so
			oz := -at.off[0]*so + at.off[2]*co
			rx = methX + ox
			ry = methY + at.off[1]
			rz = oz
		}
		x := rx + (at.px-rx)*e
		y := ry + (at.py-ry)*e
		z := rz + (at.pz-rz)*e
		// 飞行弧线：正弦抬升（起飞的抛物线感）。
		if ki > 0 && ki < 1 {
			y += math.Sin(math.Pi*e) * 0.4
		}
		cur[i] = [3]float64{x, y, z}
		ks[i] = ki
	}

	// 键：反应物键在任一端起飞时淡出；生成物键在两端落位后浮现。
	drawBonds := func() {
		for _, b := range bonds {
			pa, pb := cur[b.a], cur[b.b]
			sa, sya, dza, saok := project(a, pa[0], pa[1], pa[2])
			sb, syb, dzb, sbok := project(a, pb[0], pb[1], pb[2])
			if !saok || !sbok {
				continue
			}
			var alpha float32
			if !b.product {
				kmax := math.Max(ks[b.a], ks[b.b])
				alpha = float32(1 - clamp01((kmax-0.06)/0.3))
			} else {
				kmin := math.Min(ks[b.a], ks[b.b])
				alpha = float32(clamp01((kmin - 0.72) / 0.28))
			}
			if alpha <= 0.01 {
				continue
			}
			depth := (dza + dzb) / 2
			w := float32(0.10)
			if b.order == 2 {
				// 双键：沿屏幕垂线偏移的两根平行棍。
				dx, dy := sb-sa, syb-sya
				l := float32(math.Hypot(float64(dx), float64(dy)))
				if l < 1 {
					continue
				}
				px, py := -dy/l*5, dx/l*5
				stick(d, sa+px, sya+py, sb+px, syb+py, float32(depth), w*0.7, pale, alpha*0.95)
				stick(d, sa-px, sya-py, sb-px, syb-py, float32(depth), w*0.7, pale, alpha*0.95)
				continue
			}
			stick(d, sa, sya, sb, syb, float32(depth), w, pale, alpha)
		}
	}

	// 原子球：深度排序（远→近），元素符号标在球心。
	type rendered struct {
		depth      float64
		sx, sy, sr float32
		el         string
	}
	var balls []rendered
	for i, at := range ats {
		sx, sy, ppu, ok := project(a, cur[i][0], cur[i][1], cur[i][2])
		if !ok {
			continue
		}
		balls = append(balls, rendered{depthAt(a, cur[i]), sx, sy, ballRadius(ppu, at.r), at.el})
	}
	drawBonds()
	sort.Slice(balls, func(i, j int) bool { return balls[i].depth > balls[j].depth })
	for _, b := range balls {
		ball(d, b.sx, b.sy, b.sr, b.el, alpha)
		if b.sr > 20 {
			col := ink
			if b.el == "H" {
				col = core.RGB(0x2B3240)
			}
			coreText(d, b.el, b.sx, b.sy+float32(b.sr)*0.16, float32(b.sr)*0.72, col, alpha*0.95, 0.5)
		}
	}
}

// project 投影到屏幕；返回屏幕坐标、深度（排序用）与每世界单位的像素数。
func project(a *mathkit.Axes3D, x, y, z float64) (float32, float32, float64, bool) {
	sx, sy, cz, ok := a.Project(x, y, z)
	if !ok {
		return 0, 0, 0, false
	}
	focal := (float64(a.Size[0]) + float64(a.Size[1])) / 2
	return sx, sy, focal / cz, true
}

// ballRadius 世界半径 → 屏幕像素半径。
func ballRadius(pxPerUnit float64, r float64) float32 {
	return float32(r * pxPerUnit)
}

// depthAt 返回世界点的相机深度（排序用）。
func depthAt(a *mathkit.Axes3D, p [3]float64) float64 {
	_, _, cz, _ := a.Project(p[0], p[1], p[2])
	return cz
}

// coreText 左/中/右对齐文字（align 0 左 0.5 中 1 右）。
func coreText(d *core.Drawing, s string, x, y, size float32, c core.Color, alpha float32, align float32) {
	if alpha <= 0.001 || s == "" {
		return
	}
	w := d.MeasureText(s, core.Font{Size: core.Points(size)})
	d.SetFillColor(c)
	d.SetAlpha(alpha)
	d.FillText(s, x-align*w, y, core.Font{Size: core.Points(size)})
}

func clamp01(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}
