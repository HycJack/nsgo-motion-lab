// 《圆柱的侧面展开》—— manim 风格教学片。
//
// 全片是一个连续的 3D 世界 + 一条全局相机路径（camPose，影片时间的
// 纯函数）：环绕亮相 → 滑翔到侧面看滚动展开 → 升高俯视量尺寸。
// 幕与幕之间没有剪切，视角的切换就是相机本身的运动。
//
// 三幕（Series 排布，幕间 Cross 重叠衔接，Duration 自动）：
//  1. 圆柱：圆边、打光侧面，标出 r（截面半径）和 h（母线长）
//  2. 展开：放倒的圆柱沿地面无滑动滚动一圈，侧面沿母线展开铺平
//  3. 量一量：长方形的 C、h 尺寸线 + 公式 S = C·h = 2πrh
//
// 无声 + 字幕，审片 ?audit=1（空格/←→/J L），时间轴 ?timeline=1。
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

	r3 = 1.0 // 底面半径（世界单位）
	h3 = 2.0 // 高 = 母线长（轴长）
	z0 = 0.5 // 放倒后侧面沿 z 轴的起点
)

var (
	ink   = core.RGB(0xe8eef7)
	gold  = core.RGB(0xf2b544)
	cyan  = core.RGB(0x6fd3e8)
	cBack = core.RGB(0x2f6d80) // 侧面反面（内侧）
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
			film.Shot{D: 8, Fn: sceneUnroll},
			film.Shot{D: 6, Fn: sceneMeasure, Out: film.FadeBlack},
		),
		Subs: []film.Sub{
			{A: 1.0, B: 6.2, Text: "一个圆柱：底面是半径为 r 的圆，高 h 是母线长。"},
			{A: 7.4, B: 11.0, Text: "放倒，沿地面滚动一圈——侧面就沿着母线展开。"},
			{A: 11.4, B: 13.8, Text: "铺平后是一个长方形。"},
			{A: 14.8, B: 17.6, Text: "长恰好是底面圆的周长 C = 2πr，宽就是高 h。"},
			{A: 18.6, B: 20.2, Text: "所以，圆柱的侧面积 S = C · h = 2πrh。", Dialogue: true},
		},
		SubStyle: film.DefaultSubStyle(ink, gold, core.RGB(0x000000)),
	})
}

// ---------------------------------------------------------------- 相机

// pose 是一个机位。
type pose struct {
	eye, look [3]float64
}

func lerpPose(a, b pose, p float64) pose {
	q := pose{}
	for i := 0; i < 3; i++ {
		q.eye[i] = a.eye[i] + (b.eye[i]-a.eye[i])*p
		q.look[i] = a.look[i] + (b.look[i]-a.look[i])*p
	}
	return q
}

// camPose 是全片唯一的相机路径（f.T 的纯函数）：
// 环绕亮相（缓慢漂移，6.5s 冻结）→ EaseInOut 滑翔到侧视图（看滚动）
// → 升高到俯视（量尺寸）+ 缓慢推近公式区。机位之间从不剪切。
func camPose(t float64) pose {
	side := pose{eye: [3]float64{3.14, 2.9, 8.6}, look: [3]float64{3.14, 0.5, 1.5}}
	top := pose{eye: [3]float64{math.Pi + 2.4, 7.0, 6.7}, look: [3]float64{math.Pi, 0, 1.5}}

	orbit := func(t float64) pose {
		th := 0.9 + 0.15*filmkit.Seg(t, 0.5, 6.5)
		return pose{
			eye:  [3]float64{5.2 * math.Cos(th), 2.6, 1.5 + 5.2*math.Sin(th)},
			look: [3]float64{0, 0.9, 1.5},
		}
	}

	switch {
	case t < 6.8:
		return orbit(t)
	case t < 8.6:
		return lerpPose(orbit(6.8), side, filmkit.EaseInOut(filmkit.Seg(t, 6.8, 8.6)))
	case t < 14.8:
		return side
	case t < 16.6:
		return lerpPose(side, top, filmkit.EaseInOut(filmkit.Seg(t, 14.8, 16.6)))
	default:
		// 高位缓慢推近公式区。
		return lerpPose(top, pose{eye: top.look, look: top.look}, 0.35*filmkit.Seg(t, 16.6, 20.5))
	}
}

// ---------------------------------------------------------------- 几何

// sheetPos 是侧面上一点：u ∈ [0, 2π] 是沿母线剪开后从剪缘量起的材料角，
// v ∈ [0, 1] 沿轴。K 是滚动进度：圆柱向 +x 无滑动地滚过 r·φ，圆心在
// (rφ, r)——
//   - 已离角的材料（u ≤ φ）平躺地面：x = r·u（拖在圆柱身后）；
//   - 未离角的材料包在柱面上，在地面接触点与平铺部分切线连续，
//     沿前进方向的前面包上去：圆心 + r·(sin β, −cos β)。
//
// 方向别写反：取 −sin 会让 sheet 在接触点掉头向后翻（变成「倒着卷」）。
func sheetPos(u, v, K float64) [3]float64 {
	phi := K * 2 * math.Pi
	z := z0 + v*h3
	if u <= phi {
		return [3]float64{r3 * u, 0, z}
	}
	beta := u - phi
	return [3]float64{
		r3*phi + r3*math.Sin(beta),
		r3 - r3*math.Cos(beta),
		z,
	}
}

// rimPoly 采样一条圆边（v=0 近端 / v=1 远端，随展开一起变形）。
func rimPoly(v float64, k float64, n int) [][3]float64 {
	pts := make([][3]float64, 0, n+1)
	for i := 0; i <= n; i++ {
		pts = append(pts, sheetPos(2*math.Pi*float64(i)/float64(n), v, k))
	}
	return pts
}

func cylAxes(eye, look [3]float64) *mathkit.Axes3D {
	return &mathkit.Axes3D{
		X:    [2]float64{-1.6, 9.6},
		Y:    [2]float64{-0.6, 3.2},
		Z:    [2]float64{-1.2, 3.6},
		At:   [2]float32{260, 140},
		Size: [2]float32{1400, 800},
		Eye:  eye,
		Look: look,
	}
}

// ---------------------------------------------------------------- 三幕

// sceneCylinder：轴系 → 近/远圆边 → 打光侧面 → 标 r（截面半径，
// 圆心到触地点）和 h（母线长，悬在轴线上方）。相机在环绕。
func sceneCylinder(d *core.Drawing, f *film.Frame) {
	c := camPose(f.T)
	a := cylAxes(c.eye, c.look)
	a.Draw(d, float64(filmkit.Smooth(filmkit.Seg(f.LocalT, 0.1, 0.9))), ink, film.A(0.35), 2)

	a.Poly3(d, rimPoly(0, 0, 72), float64(filmkit.EaseInOut(filmkit.Seg(f.LocalT, 0.3, 1.2))), cyan, film.A(1), 4)
	a.Poly3(d, rimPoly(1, 0, 72), float64(filmkit.EaseInOut(filmkit.Seg(f.LocalT, 1.2, 2.0))), cyan, film.A(1), 4)
	a.SurfaceShaded(d, func(u, v float64) [3]float64 { return sheetPos(u, v, 0) },
		0, 2*math.Pi, 0, 1, 36, 6,
		float64(filmkit.EaseInOut(filmkit.Seg(f.LocalT, 2.0, 3.4))), cyan, cBack, [3]float64{0.5, 0.9, 0.3}, film.A(0.85))

	// r：近端截面里从圆心竖直向下到触地点的半径。
	if p := filmkit.Seg(f.LocalT, 3.6, 4.2); p > 0 {
		a.Poly3(d, [][3]float64{{0, r3, z0}, {0, 0, z0}}, float64(p), gold, film.A(1), 4)
		if p >= 1 {
			a.Label3D(d, "r", [3]float64{0, r3 / 2, z0 - 0.24}, 44, gold, film.A(1), 0.5)
		}
	}
	// h：悬在圆柱正上方、沿轴线方向的线段（母线长）。
	if p := filmkit.Seg(f.LocalT, 4.5, 5.2); p > 0 {
		a.Poly3(d, [][3]float64{{0, 2*r3 + 0.16, z0}, {0, 2*r3 + 0.16, z0 + h3}}, float64(p), gold, film.A(1), 3)
		if p >= 1 {
			a.Label3D(d, "h", [3]float64{0, 2*r3 + 0.52, z0 + h3 / 2}, 44, gold, film.A(1), 0.5)
		}
	}
}

// sceneUnroll：侧面无滑动滚动展开（K 由局部时间驱动），相机已滑翔到侧
// 视图就位，滚动开始时相机刚好停稳。
func sceneUnroll(d *core.Drawing, f *film.Frame) {
	c := camPose(f.T)
	a := cylAxes(c.eye, c.look)
	K := float64(filmkit.Seg(f.LocalT, 2.1, 6.4))

	a.SurfaceShaded(d, func(u, v float64) [3]float64 { return sheetPos(u, v, K) },
		0, 2*math.Pi, 0, 1, 90, 6, 1, cyan, cBack, [3]float64{0.5, 0.9, 0.3}, film.A(0.85))
	a.Poly3(d, rimPoly(0, K, 90), 1, cyan, film.A(1), 3)
	a.Poly3(d, rimPoly(1, K, 90), 1, cyan, film.A(1), 3)
	edge := func(u float64) [][3]float64 {
		pts := make([][3]float64, 0, 9)
		for i := 0; i <= 8; i++ {
			pts = append(pts, sheetPos(u, float64(i)/8, K))
		}
		return pts
	}
	a.Poly3(d, edge(0), 1, ink, film.A(0.6), 2)
	a.Poly3(d, edge(2*math.Pi), 1, ink, film.A(0.6), 2)
}

// sceneMeasure：铺平的长方形上标 C（前缘尺寸线）与 h（侧缘尺寸线），
// 公式以屏幕空间字幕落版；相机升到俯视并缓慢推近。
func sceneMeasure(d *core.Drawing, f *film.Frame) {
	c := camPose(f.T)
	a := cylAxes(c.eye, c.look)

	a.SurfaceShaded(d, func(u, v float64) [3]float64 { return sheetPos(u, v, 1) },
		0, 2*math.Pi, 0, 1, 90, 6, 1, cyan, cBack, [3]float64{0.5, 0.9, 0.3}, film.A(0.85))
	a.Poly3(d, rimPoly(0, 1, 90), 1, cyan, film.A(1), 3)
	a.Poly3(d, rimPoly(1, 1, 90), 1, cyan, film.A(1), 3)

	// C：前缘外侧的尺寸线 + 两端刻度。
	if p := filmkit.Seg(f.LocalT, 2.2, 2.8); p > 0 {
		a.Poly3(d, [][3]float64{{0, 0.03, z0 - 0.24}, {2 * math.Pi * r3 * float64(p), 0.03, z0 - 0.24}}, 1, gold, film.A(1), 3)
		if p >= 1 {
			a.Poly3(d, [][3]float64{{0, 0.03, z0 - 0.10}, {0, 0.03, z0 - 0.38}}, 1, gold, film.A(1), 2)
			a.Poly3(d, [][3]float64{{2 * math.Pi * r3, 0.03, z0 - 0.10}, {2 * math.Pi * r3, 0.03, z0 - 0.38}}, 1, gold, film.A(1), 2)
			if fade := film.A(filmkit.F32(filmkit.Seg(f.LocalT, 2.9, 3.3))); fade > 0 {
				a.Label3D(d, "C = 2πr", [3]float64{math.Pi * r3, 0.30, z0 - 0.24}, 52, gold, fade, 0.5)
			}
		}
	}
	// h：右端外侧的尺寸线。
	if p := filmkit.Seg(f.LocalT, 3.4, 4.0); p > 0 {
		a.Poly3(d, [][3]float64{{2*math.Pi*r3 + 0.26, 0.03, z0}, {2*math.Pi*r3 + 0.26, 0.03, z0 + h3*float64(p)}}, 1, gold, film.A(1), 3)
		if p >= 1 {
			if fade := film.A(filmkit.F32(filmkit.Seg(f.LocalT, 4.0, 4.4))); fade > 0 {
				a.Label3D(d, "h", [3]float64{2*math.Pi*r3 + 0.55, 0.06, z0 + h3 / 2}, 48, gold, fade, 0.5)
			}
		}
	}

	// 公式：屏幕空间底部落版（相机推近时它在画面里是稳定的）。
	if p := filmkit.Seg(f.LocalT, 4.4, 5.6); p > 0 {
		mathkit.EqWriteGradient(d, EqEqS, dw/2-EqEqS.W*1.5/2, dh-170, 1.5, float64(p), gold, ink, film.A(1), 5)
	}
}
