// 《Math Demo》—— mathkit + film 引擎表达能力的能力演示片。
//
// 这版演示装配层的三件新武器（老写法 Shot{A,B,Fn} 依然可用）：
//   - Series：分镜只声明时长，A/B 与转场重叠自动排布；Duration 自动取幕尾
//   - Seq + Frame.LocalT：场景内从 0 计时，可整体搬家、可复用
//   - Transition：In=FadeBlack（经底色过渡）/ In=Cut（硬切）/ Out=FadeBlack（片尾淡出）
// 字幕演示：自动换行（长句）、SubTop、SubCenter、逐条 Size/Color 覆写。
//
// 无声片：观看模式按墙钟走，审片用 ?audit=1 跳转。
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
)

var (
	ink   = core.RGB(0xe8eef7) // 主线条
	gold  = core.RGB(0xf2b544) // 强调
	cyan  = core.RGB(0x6fd3e8) // 曲线
	rose  = core.RGB(0xf08fa0) // 面积
	paper = core.RGB(0x0d1520) // 底色（FadeBlack 过渡到的颜色）
)

func main() {
	film.Run(film.Film{
		W: dw,
		H: dh,
		// Duration 不填：自动取最晚的幕尾（15s）
		Fade: fade,
		Bg:   paper,
		Shots: film.Series(fade,
			film.Shot{D: 6.5, Fn: sceneDerivative},
			film.Shot{D: 6, Fn: sceneEquations, In: film.FadeBlack},
			film.Shot{D: 7, Fn: sceneSurface3D, In: film.Cut, Out: film.FadeBlack},
		),
		Subs: []film.Sub{
			// 长句自动换行成两行，底板随行数走。
			{A: 1.8, B: 5.4, Text: "抛物线的切线斜率随 x 增大而增大——这就是导数的几何意义。"},
			// 顶部 + 逐条覆写颜色字号。
			{A: 7.0, B: 8.4, Text: "公式逐笔写（Write）", Pos: film.SubTop, Size: 40, Color: cyan},
			{A: 9.2, B: 11.0, Text: "变形（Transform）", Dialogue: true},
			// 居中卡片。
			{A: 12.4, B: 14.6, Text: "3D：CPU 投影 + 深度排序", Pos: film.SubCenter, Size: 56, Color: gold},
		},
		SubStyle: film.DefaultSubStyle(ink, gold, core.RGB(0x000000)),
	})
}

// mathAxes 是导数场景的 2D 坐标系。
var mathAxes = mathkit.Axes{
	X:    [2]float64{-2.2, 2.2},
	Y:    [2]float64{-0.6, 4.6},
	At:   [2]float32{460, 170},
	Size: [2]float32{1000, 720},
	Nums: true,
	Tick: 8,
}

func parabola(x float64) float64 { return x * x }

// sceneDerivative：建轴 → 描画抛物线 → 面积 → 切线滑动。
// 全部用 f.LocalT——场景从 0 计时，搬到任何起点都不用改一个数。
func sceneDerivative(d *core.Drawing, f *film.Frame) {
	axesA := filmkit.F32(filmkit.Smooth(filmkit.Seg(f.LocalT, 0.3, 1.5)))
	mathAxes.Draw(d, float64(axesA), ink, film.A(0.85), 3)

	plotP := filmkit.Seg(f.LocalT, 1.5, 3.2)
	mathkit.Plot(d, mathAxes, parabola, -2, 2, plotP, 120, cyan, film.A(1), 5)
	if plotP >= 1 {
		mathkit.Label(d, "y = x²", mathAxes.At[0]+mathAxes.Size[0]-40, mathAxes.At[1]+80, 44, cyan, film.A(1), 1)
	}

	areaP := filmkit.Smooth(filmkit.Seg(f.LocalT, 3.8, 5.2))
	if areaP > 0 {
		x1 := -1 + 2.2*float64(areaP)
		mathkit.Area(d, mathAxes, parabola, -1, x1, 60, rose, film.A(0.35))
	}

	tanP := filmkit.Seg(f.LocalT, 5.0, 6.3)
	if tanP > 0 {
		x0 := -1.4 + 2.8*float64(tanP)
		mathkit.Tangent(d, mathAxes, parabola, x0, 130, gold, film.A(0.95), 3)
		px, py := mathAxes.Screen(x0, parabola(x0))
		mathkit.Dot(d, px, py, 10, gold, film.A(1)) // 接触点：跟着切线滑的圆点
		mathkit.Arrow(d, px+66, py-118, px+24, py-34, 14, ink, film.A(0.8), 2)
		mathkit.Label(d, "f′(x)", px+74, py-128, 38, gold, film.A(0.95), 0)
	}
}

// sceneEquations：逐笔写 y=x²，再变形为 dy/dx=2x。
func sceneEquations(d *core.Drawing, f *film.Frame) {
	const scale = 2.4
	left := (dw - max(EqEqParabola.W, EqEqDerivative.W)*scale) / 2
	baseY := dh / 2

	writeP := filmkit.Smooth(filmkit.Seg(f.LocalT, 0.5, 2.2))
	if writeP < 1 {
		// 逐笔写期间只画描边；写完立即换成填充态，描边不残留。
		mathkit.EqWrite(d, EqEqParabola, left, baseY, scale, float64(writeP), ink, film.A(1), 4)
		return
	}

	morphP := filmkit.Smooth(filmkit.Seg(f.LocalT, 2.8, 4.3))
	if morphP <= 0 {
		mathkit.EqDraw(d, EqEqParabola, left, baseY, scale, ink, film.A(1))
		return
	}
	if morphP >= 1 {
		// 变形落定：弹簧给一个 5% 的过冲（spring 的手感），落回 1。
		pop := filmkit.Spring{Stiffness: 130, Damping: 8}.At(f.LocalT - 4.3)
		s := scale * filmkit.F32(1+0.05*(pop-1))
		mathkit.EqMorph(d, EqEqParabola, EqEqDerivative, left, baseY, s, 1, gold, film.A(1), 4)
		return
	}
	fadeOld := 1 - filmkit.F32(filmkit.Seg(f.LocalT, 2.8, 3.4))
	if fadeOld > 0 {
		mathkit.EqDraw(d, EqEqParabola, left, baseY, scale, ink, film.A(fadeOld*0.5))
	}
	mathkit.EqMorph(d, EqEqParabola, EqEqDerivative, left, baseY, scale, float64(morphP), gold, film.A(1), 4)
}

// sceneSurface3D：3D 坐标轴 + 鞍面线框 + 螺旋线，相机环绕。
// 局部时间轴（7s）：0.2-1.0 建轴 → 0.8-2.0 铺鞍面 → 2.0-6.5
// 相机绕 y 轴转一圈（Eye 随 LocalT 走，仍是纯函数），2.6-4.2 螺旋线描画。
func sceneSurface3D(d *core.Drawing, f *film.Frame) {
	// 相机环绕：从右侧 40° 起步，转 300°。
	orbit := func(t float64) [3]float64 {
		theta := 0.7 + 5.2*filmkit.Seg(t, 2.0, 6.5)
		return [3]float64{3.4 * math.Cos(theta), 1.9, 3.4 * math.Sin(theta)}
	}
	a := &mathkit.Axes3D{
		X:    [2]float64{-1.3, 1.3},
		Y:    [2]float64{-0.9, 0.9},
		Z:    [2]float64{-1.3, 1.3},
		At:   [2]float32{360, 240},
		Size: [2]float32{1200, 700},
		Eye:  orbit(f.LocalT),
		Look: [3]float64{0, -0.05, 0},
	}
	a.Draw(d, float64(filmkit.Smooth(filmkit.Seg(f.LocalT, 0.2, 1.0))), ink, film.A(0.8), 2)

	saddle := func(u, v float64) [3]float64 {
		return [3]float64{u, (u*u - v*v) * 0.6, v}
	}
	surfP := float64(filmkit.EaseInOut(filmkit.Seg(f.LocalT, 0.8, 2.0)))
	// 实心面片：法线打光，光从右上来。
	a.SurfaceShaded(d, saddle, -1.1, 1.1, -1.1, 1.1, 18, 18, surfP, cyan, core.RGB(0x1f5f7a), [3]float64{0.5, 0.9, 0.3}, film.A(0.9))
	// 细线框叠上去，网格感。
	if surfP >= 1 {
		a.SurfaceWireframe(d, saddle, -1.1, 1.1, -1.1, 1.1, 18, 18, 1, ink, film.A(0.25), 1)
	}

	// 螺旋线：绕 y 轴两圈，随进度描画。
	helixP := filmkit.EaseInOut(filmkit.Seg(f.LocalT, 2.6, 4.2))
	if helixP > 0 {
		pts := make([][3]float64, 0, 121)
		for i := 0; i <= 120; i++ {
			t := float64(i) / 120
			ang := 2 * math.Pi * 2 * t
			pts = append(pts, [3]float64{0.85 * math.Cos(ang), -0.75 + 1.5*t, 0.85 * math.Sin(ang)})
		}
		a.Poly3(d, pts, float64(helixP), gold, film.A(0.95), 3)
	}

	// 虚线：鞍面角点到地面的辅助垂线。
	if surfP >= 1 {
		a.Poly3Dashed(d, [][3]float64{{1.1, 0.72, 0}, {1.1, -0.5, 0}}, 1, 10, 8, ink, film.A(0.5), 2)
	}

	// 3D 标签：锚在世界坐标上，相机转它跟着走。
	lblA := film.A(filmkit.F32(filmkit.Seg(f.LocalT, 1.6, 2.0)))
	a.Label3D(d, "y", [3]float64{0, 0.95, 0}, 40, ink, lblA, 0.5)
	a.Label3D(d, "x", [3]float64{1.42, 0, 0}, 40, ink, lblA, 0.5)
	a.Label3D(d, "z", [3]float64{0, 0, 1.42}, 40, ink, lblA, 0.5)
	a.Label3D(d, "z = x² − y²", [3]float64{0, 1.18, 0}, 52, cyan, film.A(filmkit.F32(filmkit.Seg(f.LocalT, 1.6, 2.0))), 0.5)
}
