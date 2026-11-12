// 《圆柱的侧面展开》—— 一支 manim 风格的教学短片。
//
// 三幕（Series 排布，Duration 自动；无声 + 字幕，审片 ?audit=1）：
//  1. 圆柱：3D 画底圆、顶圆、侧面（法线打光），标注 r 和 h
//  2. 展开：侧面沿左缘剥离、逐条铺平到地面（参数化剥离动画），
//     相机随之升到俯视；铺平后正是「长 = 底面周长、宽 = 高」的长方形
//  3. 长方形：回到屏幕空间画长方形，标注 C = 2πr、h，写出 S = C·h = 2πrh
//
// 整支片子零贴图：面片、曲线、公式全是 mathkit 的纯函数。
package main

import (
	"math"

	"nsgofilm/film"
	"nsgofilm/filmkit"
	"nsgofilm/mathkit"

	"github.com/oneo-me/nsgo/core"
)

const (
	dw    float32 = 1920.0
	dh    float32 = 1080.0
	fade          = 0.5

	r3 = 1.0 // 圆柱底面半径（世界单位）
	h3 = 2.0 // 圆柱高
	z0 = 0.5 // 放倒后侧面在轴向上的起点
)

var (
	ink   = core.RGB(0xe8eef7)
	gold  = core.RGB(0xf2b544)
	cyan  = core.RGB(0x6fd3e8)
	paper = core.RGB(0x0d1520)
)

func main() {
	film.Run(film.Film{
		W:    dw,
		H:    dh,
		Fade: fade,
		Bg:   paper,
		Shots: film.Series(fade,
			film.Shot{D: 7, Fn: sceneCylinder},
			film.Shot{D: 8, Fn: sceneUnroll, In: film.FadeBlack},
			film.Shot{D: 7, Fn: sceneRectangle, In: film.FadeBlack, Out: film.FadeBlack},
		),
		Subs: []film.Sub{
			{A: 1.0, B: 6.2, Text: "一个圆柱：底面是半径为 r 的圆，高是 h。"},
			{A: 7.2, B: 11.2, Text: "把圆柱放倒，让它沿地面滚动一圈——侧面就沿着母线展开。"},
			{A: 11.6, B: 14.2, Text: "铺平后是一个长方形。"},
			{A: 14.6, B: 18.4, Text: "长方形的长恰好是底面圆的周长 C = 2πr，宽就是圆柱的高 h。"},
			{A: 18.8, B: 20.6, Text: "所以，圆柱的侧面积 S = C · h = 2πrh。", Dialogue: true},
		},
		SubStyle: film.DefaultSubStyle(ink, gold, core.RGB(0x000000)),
	})
}

// cylAxes 是 3D 幕共用的坐标系：长方形铺开后 x 方向要容得下 ±πr。
func cylAxes(eye, look [3]float64) *mathkit.Axes3D {
	return &mathkit.Axes3D{
		X:    [2]float64{-2.0, 8.4},
		Y:    [2]float64{-0.5, 3.0},
		Z:    [2]float64{-1.5, 3.4},
		At:   [2]float32{260, 140},
		Size: [2]float32{1400, 800},
		Eye:  eye,
		Look: look,
	}
}

// sheetPos 是侧面上一点的位置：u ∈ [0, 2π] 是沿母线剪开后从剪缘量起
// 的材料角，v ∈ [0, 1] 沿轴（圆柱放倒后是水平方向）。
// K 是滚动进度：圆柱沿地面无滑动地滚过 K·2π 的弧——
//   - 已离角的材料（u ≤ φ）平躺在地面：x = r·u；
//   - 未离角的材料仍包在柱面上，在顶部切线处与平铺部分相切连续。
// 这就是「沿母线展开」的真实运动学：展开 = 无滑动滚动。
func sheetPos(u, v, K float64) [3]float64 {
	phi := K * 2 * math.Pi
	z := z0 + v*h3
	if u <= phi {
		return [3]float64{r3 * u, 0, z} // 已铺平：弧长展开在身后
	}
	beta := u - phi // 距离触地还差的角度
	return [3]float64{
		r3*phi - r3*math.Sin(beta),
		r3 - r3*math.Cos(beta),
		z,
	}
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

// rimPoly 采样顶/底圆（随展开一起变形）。
func rimPoly(v float64, k float64, n int) [][3]float64 {
	pts := make([][3]float64, 0, n+1)
	for i := 0; i <= n; i++ {
		pts = append(pts, sheetPos(2*math.Pi*float64(i)/float64(n), v, k))
	}
	return pts
}

// sceneCylinder：底圆 → 顶圆 → 侧面 → r/h 标注。相机缓慢侧移。
func sceneCylinder(d *core.Drawing, f *film.Frame) {
	t := f.LocalT
	theta := 0.65 + 0.25*filmkit.Seg(t, 0, 7)
	a := cylAxes([3]float64{5.0 * math.Cos(theta), 3.1, 5.0 * math.Sin(theta)}, [3]float64{0, 0.9, 1.0})

	// 展开进度恒 0：这一幕圆柱还是完整的。
	const k0 = 0.0
	a.Draw(d, float64(filmkit.Smooth(filmkit.Seg(t, 0.1, 0.9))), ink, film.A(0.5), 2)

	// 底圆 → 顶圆。
	a.Poly3(d, rimPoly(0, k0, 72), float64(filmkit.EaseInOut(filmkit.Seg(t, 0.3, 1.3))), cyan, film.A(1), 4)
	a.Poly3(d, rimPoly(1, k0, 72), float64(filmkit.EaseInOut(filmkit.Seg(t, 1.3, 2.2))), cyan, film.A(1), 4)

	// 侧面（法线打光的面片）。
	a.SurfaceShaded(d, func(u, v float64) [3]float64 { return sheetPos(u, v, k0) },
		0, 2*math.Pi, 0, 1, 36, 6,
		float64(filmkit.EaseInOut(filmkit.Seg(t, 2.2, 3.8))), cyan, core.RGB(0x2f6d80), [3]float64{0.5, 0.9, 0.3}, film.A(0.85))

	// 半径 r：底面圆心到圆缘的一条线 + 标注。
	if p := filmkit.Seg(t, 3.9, 4.5); p > 0 {
		a.Poly3(d, [][3]float64{{0, 0, 0}, {0, 0, r3 * float64(p)}}, 1, gold, film.A(1), 4)
		a.Label3D(d, "r", [3]float64{0, -0.12, r3 * 0.5 * float64(p)}, 44, gold, film.A(1), 0.5)
	}
	// 高 h：柱后一条竖线 + 标注。
	if p := filmkit.Seg(t, 4.7, 5.4); p > 0 {
		a.Poly3(d, [][3]float64{{-1.55 * r3, 0, -1.55 * r3}, {-1.55 * r3, h3 * float64(p), -1.55 * r3}}, 1, gold, film.A(1), 4)
		a.Label3D(d, "h", [3]float64{-1.55 * r3, h3 * float64(p) + 0.15, -1.55 * r3}, 44, gold, film.A(1), 0.5)
	}
}

// sceneUnroll：圆柱放倒沿地面滚一圈，侧面沿母线展开铺在身后。
// 局部时间轴（8s）：0-1.2 静止亮相（放倒的圆柱）→ 1.2-6.2 匀速滚动
// → 6.2-8 铺平的长方形停在地面（长 = 2πr，宽 = h）。
func sceneUnroll(d *core.Drawing, f *film.Frame) {
	t := f.LocalT
	// 固定的侧偏俯视机位，框住整个滚动区间。
	a := cylAxes([3]float64{3.1, 2.9, 8.4}, [3]float64{3.1, 0.8, 1.4})

	K := filmkit.Seg(t, 1.2, 6.2) // 滚动进度：一圈 = 2π
	a.Draw(d, 1, ink, film.A(0.5), 2)

	// 侧面（滚动中的柱面 + 已铺平部分一体）。
	a.SurfaceShaded(d, func(u, v float64) [3]float64 { return sheetPos(u, v, float64(K)) },
		0, 2*math.Pi, 0, 1, 90, 6, 1, cyan, core.RGB(0x2f6d80), [3]float64{0.5, 0.9, 0.3}, film.A(0.85))

	// 边界：两条圆边（v=0/v=1，滚动中变形）+ 两条剪开缘（u=0/u=2π）。
	a.Poly3(d, rimPoly(0, float64(K), 90), 1, cyan, film.A(1), 3)
	a.Poly3(d, rimPoly(1, float64(K), 90), 1, cyan, film.A(1), 3)
	edge := func(u float64) [][3]float64 {
		pts := make([][3]float64, 0, 9)
		for i := 0; i <= 8; i++ {
			pts = append(pts, sheetPos(u, float64(i)/8, float64(K)))
		}
		return pts
	}
	a.Poly3(d, edge(0), 1, ink, film.A(0.7), 2)
	a.Poly3(d, edge(2*math.Pi), 1, ink, film.A(0.7), 2)

	// 铺平后：标出长（= 底面周长）与宽（= 高）。
	if K >= 1 {
		fade := film.A(filmkit.F32(filmkit.Seg(t, 6.4, 7.0)))
		a.Label3D(d, "C = 2πr", [3]float64{math.Pi * r3, 0.05, z0 - 0.28}, 56, gold, fade, 0.5)
		a.Label3D(d, "h", [3]float64{2*math.Pi*r3 + 0.22, 0.02, z0 + h3/2}, 48, gold, fade, 0)
	}
}

// sceneRectangle：屏幕空间的长方形 + 尺寸标注 + 公式。
func sceneRectangle(d *core.Drawing, f *film.Frame) {
	t := f.LocalT
	cx, cy := dw/2, float32(440)
	rw, rh := float32(860), float32(860/(2*math.Pi)*2.0) // 长宽比 = 2πr : h = 2π : 2

	// 长方形从中间向两边展开（承接上一幕的「展开」）。
	rev := filmkit.Smooth(filmkit.Seg(t, 0.2, 1.0))
	halfW := rw / 2 * filmkit.F32(rev)
	d.SetFillColor(core.RGBA(0x6fd3e8, 0.16))
	d.SetAlpha(film.A(1))
	d.BeginPath()
	d.Rect(cx-halfW, cy-rh/2, halfW*2, rh)
	d.Fill()
	d.SetStrokeColor(ink)
	d.SetLineWidth(3)
	d.BeginPath()
	d.Rect(cx-halfW, cy-rh/2, halfW*2, rh)
	d.Stroke()

	if rev < 1 {
		return
	}

	// 上方尺寸线：C = 2πr（底面周长）。
	dimP := filmkit.Seg(t, 1.2, 1.9)
	if dimP > 0 {
		y := cy - rh/2 - 70
		x1 := cx - rw/2 + rw*filmkit.F32(dimP)
		d.SetStrokeColor(gold)
		d.SetLineWidth(2)
		d.BeginPath()
		d.MoveTo(cx-rw/2, cy-rh/2-12)
		d.LineTo(cx-rw/2, y+14)
		d.LineTo(x1, y+14)
		d.LineTo(cx+rw/2, cy-rh/2-12)
		d.Stroke()
		mathkit.EqWriteGradient(d, EqEqC, cx-EqEqC.W*1.1/2, y-26, 1.1, float64(filmkit.Seg(t, 1.9, 2.8)), gold, ink, film.A(1), 4)
	}

	// 右侧尺寸线：h。
	if p := filmkit.Seg(t, 2.6, 3.2); p > 0 {
		x := cx + rw/2 + 70
		y1 := cy - rh/2 + rh*filmkit.F32(p)
		d.SetStrokeColor(gold)
		d.SetLineWidth(2)
		d.BeginPath()
		d.MoveTo(cx+rw/2+12, cy-rh/2)
		d.LineTo(x-14, cy-rh/2)
		d.LineTo(x-14, y1)
		d.LineTo(cx+rw/2+12, cy+rh/2)
		d.Stroke()
		mathkit.Label(d, "h", x+16, cy, 52, gold, film.A(1), 0.5)
	}

	// 公式：S = C·h = 2πrh（逐笔写）。
	if p := filmkit.Seg(t, 3.8, 5.4); p > 0 {
		mathkit.EqWriteGradient(d, EqEqS, cx-EqEqS.W*1.6/2, cy+rh/2+170, 1.6, float64(p), gold, ink, film.A(1), 5)
	}
}
