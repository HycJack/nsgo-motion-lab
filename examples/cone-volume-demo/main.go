// 《圆锥的体积》—— 切片法积分推导教学片（呼应圆柱侧面展开片）。
//
// CH₄ 那支证明了「原子重新组合」，这一支证明「一个积分把立体变成
// 精确公式」：圆锥体积 = 同底等高圆柱的 1/3。
//
// 四幕（Series，一台连续相机：环绕 → 升到侧面 → 缓慢推近公式）：
//  1. 同底等高：圆柱线框幽灵 + 实心圆锥，V_柱 = πr²h
//  2. 切片：圆锥被切成薄片，逐片浮现——每片半径 ρ(z) = r(1−z/h)
//  3. 相似三角形：母线上的剖面（apex→底缘），ρ(z) 与 z 的标注
//  4. 积分：V = ∫₀ʰ π ρ(z)² dz = ⅓ π r² h，落版「几何与微积分汇合」
//
// 全部画面是 mathkit 的纯函数：圆锥/圆柱曲面、圆盘、剖面线、公式。
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

	r3 = 1.5 // 底面半径（世界单位）
	h3 = 3.2 // 高
)

var (
	ink   = core.RGB(0xE8EEF7)
	gold  = core.RGB(0xFFC53D)
	cyan  = core.RGB(0x53C8E8)
	cBack = core.RGB(0x1B6E8C) // 锥面反面
	coneF = core.RGB(0x53C8E8)
	ghost = core.RGB(0x5A6B85)
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
			film.Shot{D: 6.5, Fn: sceneCompare},
			film.Shot{D: 8, Fn: sceneSlice},
			film.Shot{D: 6.5, Fn: sceneTriangle},
			film.Shot{D: 6.5, Fn: sceneIntegral, Out: film.FadeBlack},
		),
		Subs: []film.Sub{
			{A: 1.0, B: 5.6, Text: "同底等高的圆柱和圆锥——圆锥的体积恰好是它的三分之一。"},
			{A: 6.2, B: 9.2, Text: "把圆锥切成薄片：每一片都是一个薄薄的圆盘。"},
			{A: 9.6, B: 13.2, Text: "相似三角形给出每片的半径：ρ(z) = r(1 − z/h)。"},
			{A: 13.6, B: 17.0, Text: "把所有圆盘的体积加起来——一个积分就够了。"},
			{A: 17.4, B: 19.6, Text: "几何与微积分，在同一个 1/3 上汇合。", Dialogue: true},
		},
		SubStyle: film.DefaultSubStyle(ink, gold, core.RGB(0x000000)),
	})
}

// ---------------------------------------------------------------- 几何

// conePt 是圆锥/圆柱上的一点：u ∈ [0, 2π]（角度）、v ∈ [0, 1]
// （0 = 底面、1 = 锥顶）。rhoFactor 是半径系数：圆锥 1−v，圆柱 1。
func conePt(u, v, rhoFactor float64) [3]float64 {
	rho := r3 * rhoFactor * (1 - v)
	return [3]float64{rho * math.Cos(u), v * h3, rho * math.Sin(u)}
}

// cylinderPt 是圆柱上的一点（半径不随高度收缩）。
func cylinderPt(u, v float64) [3]float64 {
	return [3]float64{r3 * math.Cos(u), v * h3, r3 * math.Sin(u)}
}

// coneAxes：全片一台相机——环绕亮相后缓推。
func coneAxes(t float64) *mathkit.Axes3D {
	th := 0.55 + 0.18*filmkit.Seg(t, 0.5, 6.5)
	push := filmkit.Seg(t, 14.0, 18.0) // 末段缓慢推近
	eye := [3]float64{
		6.4*math.Cos(th) + 1.2*push,
		3.3 + 1.1*filmkit.Seg(t, 8.0, 14.0),
		6.4*math.Sin(th),
	}
	look := [3]float64{0, 1.4 + 0.25*push, 0}
	return &mathkit.Axes3D{
		X: [2]float64{-2.6, 2.6}, Y: [2]float64{-0.6, 4.0}, Z: [2]float64{-2.6, 2.6},
		At: [2]float32{300, 110}, Size: [2]float32{1320, 850},
		Eye:  eye,
		Look: look,
	}
}

// disk 在高度 y 画一个薄圆盘（投影后填充 + 描边）。
func disk(d *core.Drawing, a *mathkit.Axes3D, y, rho float64, fill core.Color, alpha float32) {
	n := 48
	pts := make([][2]float32, 0, n)
	for k := 0; k < n; k++ {
		th := 2 * math.Pi * float64(k) / float64(n)
		sx, sy, ok := project(a, rho*math.Cos(th), y, rho*math.Sin(th))
		if !ok {
			continue
		}
		pts = append(pts, [2]float32{sx, sy})
	}
	if len(pts) < 3 {
		return
	}
	d.SetFillColor(fill)
	d.SetAlpha(alpha)
	d.BeginPath()
	d.MoveTo(pts[0][0], pts[0][1])
	for _, p := range pts[1:] {
		d.LineTo(p[0], p[1])
	}
	d.ClosePath()
	d.Fill()
	d.SetStrokeColor(core.RGB(0xFFFFFF))
	d.SetAlpha(alpha * 0.35)
	d.SetLineWidth(1.5)
	d.BeginPath()
	d.MoveTo(pts[0][0], pts[0][1])
	for _, p := range pts[1:] {
		d.LineTo(p[0], p[1])
	}
	d.ClosePath()
	d.Stroke()
	d.SetAlpha(1)
}

// ring 描一条圆边（Poly3 圆）。
func ring(a *mathkit.Axes3D, y, rho float64, n int) [][3]float64 {
	pts := make([][3]float64, 0, n+1)
	for k := 0; k <= n; k++ {
		th := 2 * math.Pi * float64(k) / float64(n)
		pts = append(pts, [3]float64{rho * math.Cos(th), y, rho * math.Sin(th)})
	}
	return pts
}

func project(a *mathkit.Axes3D, x, y, z float64) (float32, float32, bool) {
	sx, sy, _, ok := a.Project(x, y, z)
	return sx, sy, ok
}

// ---------------------------------------------------------------- 幂

// sceneCompare：圆柱线框幽灵 + 实心圆锥——同底等高，1/3 关系亮相。
func sceneCompare(d *core.Drawing, f *film.Frame) {
	a := coneAxes(f.T)
	a.Draw(d, float64(filmkit.Smooth(filmkit.Seg(f.LocalT, 0.1, 0.9))), ink, film.A(0.3), 2)

	// 圆柱线框幽灵：上下圆边 + 四根竖线。
	fadeIn := film.A(filmkit.F32(filmkit.Seg(f.LocalT, 0.4, 1.2)))
	a.Poly3(d, ring(a, 0, r3, 72), 1, ghost, fadeIn, 2)
	a.Poly3(d, ring(a, h3, r3, 72), 1, ghost, fadeIn, 2)
	for k := 0; k < 4; k++ {
		th := 2 * math.Pi * float64(k) / 4
		a.Poly3(d, [][3]float64{
			cylinderPt(th, 0), cylinderPt(th, 1),
		}, 1, ghost, fadeIn*0.7, 2)
	}

	// 实心圆锥：底圆 → 侧面。
	a.Poly3(d, ring(a, 0, r3, 72), float64(filmkit.EaseInOut(filmkit.Seg(f.LocalT, 1.4, 2.4))), cyan, film.A(1), 4)
	a.SurfaceShaded(d, func(u, v float64) [3]float64 { return conePt(u, v, 1) },
		0, 2*math.Pi, 0, 1, 48, 8,
		float64(filmkit.EaseInOut(filmkit.Seg(f.LocalT, 2.4, 3.8))), coneF,
		core.RGB(0x1B6E8C), [3]float64{0.5, 0.9, 0.3}, film.A(0.9))

	// 文字：同底等高 + 1/3。
	if e := film.A(filmkit.F32(filmkit.Seg(f.LocalT, 3.9, 4.5))); e > 0 {
		coreTextL(d, "同底等高", dw*0.06, dh*0.16, 46, ink, e, 0)
	}
	if e := film.A(filmkit.F32(filmkit.Seg(f.LocalT, 4.6, 5.2))); e > 0 {
		coreTextL(d, "V(锥) = ⅓ · V(柱)", dw*0.06, dh*0.16+64, 46, gold, e, 0)
	}
}

// sceneSlice：圆锥切成薄片，逐片浮现（自顶向下），同时柱面幽灵衬托。
func sceneSlice(d *core.Drawing, f *film.Frame) {
	a := coneAxes(f.T)
	a.Draw(d, 1, ink, film.A(0.3), 2)

	const nDisks = 12
	start := 0.6
	per := 0.42

	// 每片的浮现进度（自顶向下：薄片一片片摞上来）。
	type dslice struct {
		y, rho float64
		p      float64
	}
	var slices []dslice
	for i := 0; i < nDisks; i++ {
		y := (float64(i) + 0.5) / nDisks * h3
		rho := r3 * (1 - y/h3)
		at := start + float64(nDisks-1-i)*per
		p := filmkit.Clamp01((f.LocalT - at) / 0.4)
		if p > 0 {
			slices = append(slices, dslice{y, rho, filmkit.Smooth(p)})
		}
	}

	// 远→近画（深度排序：相机在 +z 前方，z 小的远）。
	sort.Slice(slices, func(i, j int) bool { return slices[i].y < slices[j].y })

	// 柱面幽灵衬托（切片的"容器"）。
	a.Poly3(d, ring(a, 0, r3, 72), 1, ghost, film.A(0.5), 2)
	a.Poly3(d, ring(a, h3, r3, 72), 1, ghost, film.A(0.35), 2)

	for _, s := range slices {
		rr := s.rho * s.p
		disk(d, a, s.y, rr, core.RGBA(0x53C8E8, 0.5), film.A(0.85))
	}

	// 每片的半径标注：ρ(z) 随高度递减。
	if e := film.A(filmkit.F32(filmkit.Seg(f.LocalT, start+float64(nDisks)*per-0.6, start+float64(nDisks)*per))); e > 0 {
		coreTextL(d, "每片半径 ρ(z) = r(1 − z/h)", dw*0.06, dh*0.16, 44, ink, e, 0)
	}
}

// sceneTriangle：相似三角形——母线、ρ(z) 水平线、z 竖直线。
func sceneTriangle(d *core.Drawing, f *film.Frame) {
	a := coneAxes(f.T)
	a.Draw(d, 1, ink, film.A(0.35), 2)

	// 剖面：apex → 底缘的金色母线（相似三角形的斜边）。
	z0 := h3 * 0.62
	rho0 := r3 * (1 - z0/h3)
	a.Poly3(d, [][3]float64{{0, h3, 0}, {r3, 0, 0}}, float64(filmkit.Seg(f.LocalT, 0.4, 1.2)), gold, film.A(1), 4)

	// ρ(z) 水平线 + z 竖直线（虚线），构成相似三角形。
	if p := filmkit.Seg(f.LocalT, 1.4, 2.2); p > 0 {
		a.Poly3Dashed(d, [][3]float64{{0, z0, 0}, {rho0 * float64(p), z0, 0}}, 1, 10, 8, cyan, film.A(1), 3)
		a.Label3D(d, "ρ(z)", [3]float64{rho0 / 2, z0 + 0.18, 0}, 40, cyan, film.A(float32(p)), 0.5)
	}
	if p := filmkit.Seg(f.LocalT, 2.0, 2.7); p > 0 {
		a.Poly3Dashed(d, [][3]float64{{rho0, 0, 0}, {rho0, z0 * float64(p), 0}}, 1, 10, 8, ink, film.A(0.7), 3)
		a.Label3D(d, "z", [3]float64{rho0 + 0.22, z0 * float64(p) / 2, 0}, 40, ink, film.A(float32(p)), 0)
	}

	// 相似三角形比：ρ(z)/r = (h−z)/h。
	if e := film.A(filmkit.F32(filmkit.Seg(f.LocalT, 3.0, 3.7))); e > 0 {
		coreTextL(d, "ρ(z) / r = (h − z) / h", dw*0.06, dh*0.16, 46, gold, e, 0)
	}
}

// sceneIntegral：积分式逐段写 + 结果 1/3 落版。
func sceneIntegral(d *core.Drawing, f *film.Frame) {
	a := coneAxes(f.T)
	a.Draw(d, 1, ink, film.A(0.3), 2)
	// 背景薄切片淡淡衬托。
	for i := 0; i < 10; i++ {
		y := (float64(i) + 0.5) / 10 * h3
		disk(d, a, y, r3*(1-y/h3), core.RGBA(0x53C8E8, 0.10), film.A(0.5))
	}

	// 公式（带下标/上标的分段文字，逐段浮现）。
	segs := []formulaSeg{
		{t: "V  =  ∫"}, {t: "0", sub: true}, {t: "h", sup: true},
		{t: "  π r"}, {t: "2", sup: true},
		{t: " (1 − z/h)"}, {t: "2", sup: true},
		{t: "  dz"},
	}
	fP := filmkit.Seg(f.LocalT, 1.6, 3.4)
	drawFormula(d, segs, dw/2, dh/2-40, 84, ink, float32(fP), 0.5)

	// 结果。
	if e := film.A(filmkit.F32(filmkit.Seg(f.LocalT, 3.8, 4.6))); e > 0 {
		result := []formulaSeg{
			{t: "V  =  "}, {t: "1", sup: false}, {t: "/", sub: false}, {t: "3"}, {t: "  π r"}, {t: "2", sup: true}, {t: "h"},
		}
		drawFormula(d, result, dw/2, dh/2+90, 96, gold, float32(e), 0.5)
	}
	if e := film.A(filmkit.F32(filmkit.Seg(f.LocalT, 4.8, 5.4))); e > 0 {
		coreTextL(d, "几何与微积分，在同一个 1/3 上汇合", dw/2, dh*0.85, 46, pale, e, 0.5)
	}
}

// ---------------------------------------------------------------- 公式文字

// formulaSeg 是公式的一段：sub 下标 / sup 上标。
type formulaSeg struct {
	t   string
	sub bool
	sup bool
}

// drawFormula 画带上下标的公式文字（align: 0 左 0.5 中 1 右）。
func drawFormula(d *core.Drawing, segs []formulaSeg, x, y, size float32, c core.Color, alpha float32, align float32) {
	total := float32(0)
	widths := make([]float32, len(segs))
	for i, s := range segs {
		fs := size
		if s.sub || s.sup {
			fs = size * 0.62
		}
		widths[i] = d.MeasureText(s.t, core.Font{Size: core.Points(fs)})
		total += widths[i]
	}
	cur := x - align*total
	for i, s := range segs {
		fs := size
		base := y
		if s.sub {
			fs = size * 0.62
			base = y + size*0.24
		}
		if s.sup {
			fs = size * 0.62
			base = y - size*0.42
		}
		d.SetFillColor(c)
		d.SetAlpha(alpha)
		d.FillText(s.t, cur, base, core.Font{Size: core.Points(fs)})
		cur += widths[i]
	}
}

// coreTextL 左/中/右对齐文字。
func coreTextL(d *core.Drawing, s string, x, y, size float32, c core.Color, alpha float32, align float32) {
	if alpha <= 0.001 || s == "" {
		return
	}
	w := d.MeasureText(s, core.Font{Size: core.Points(size)})
	d.SetFillColor(c)
	d.SetAlpha(alpha)
	d.FillText(s, x-align*w, y, core.Font{Size: core.Points(size)})
}
