package film

import (
	"math"

	"github.com/oneo-me/nsgo/core"

	"nsgofilm/filmkit"
)

// Shot 是一幕。两种写法：
//   - 手写时间轴：A、B 是全片时间（秒）；
//   - Series 排布：只给 D（时长），A/B 由 Series 算。
//
// In/Out 是进/出场转场，零值 Cross = 现状的交叉溶解。
// Fn 画的是世界坐标里的内容；字幕由引擎最后补画（屏幕空间）。
type Shot struct {
	A, B float64
	D    float64 // Series 用：时长（A/B 为零时由 Series 填）
	Fn   func(d *core.Drawing, f *Frame)
	In   Transition
	Out  Transition
}

// Frame 是一帧的全部上下文，传给每个分镜函数。
type Frame struct {
	T    float64 // 影片时间（秒）
	P    float64 // 本幕进度 0..1
	Ease float64 // 本幕的缓动进度 0..1（EaseOut）
	// LocalT/LocalD 是本幕的局部时间与时长（T-A、B-A）——
	// 用它们写场景，场景就能整体搬家、复用（Seq 的全部秘密）。
	LocalT, LocalD float64
	W, H           float32 // 设计画幅
}

// Cam 把取景（平移 + 缩放）装到画布上，之后所有坐标都是世界坐标。
// 调用方自己配 Save/Restore。
//
// 约束：zoom ≥ 1。背景板只有一张画幅大小，zoom < 1 会露黑边；
// 推近靠放大角色（绘本就是这么干的），不靠拉镜头。
// 平移范围 |x| ≤ W/2·(zoom-1)、|y| ≤ H/2·(zoom-1)。
func (f *Frame) Cam(d *core.Drawing, x, y, zoom float32) {
	d.Translate(f.W/2+x, f.H/2+y)
	d.Scale(zoom, zoom)
	d.Translate(-f.W/2, -f.H/2)
}

// resolve 补全 Film 的派生量：Duration 未填时取最晚的幕尾。
func (f *Film) resolve() {
	if f.Duration <= 0 {
		for i := range f.Shots {
			if f.Shots[i].B > f.Duration {
				f.Duration = f.Shots[i].B
			}
		}
	}
}

// render 把这一帧画出来：底色 → 依次画各幕（转场包络）→ 转场扫瞙 → 字幕。
func render(d *core.Drawing, size core.Size, t float64, f Film) {
	f.resolve()
	s := float32(min(float64(size.Width)/float64(f.W), float64(size.Height)/float64(f.H)))
	d.Save()
	d.Translate((float32(size.Width)-f.W*s)/2, (float32(size.Height)-f.H*s)/2)
	d.Scale(s, s)

	// 底色：幕间露出的那层。
	d.SetAlpha(1)
	d.SetFillColor(f.Bg)
	d.BeginPath()
	d.Rect(0, 0, f.W, f.H)
	d.Fill()

	last := f.Duration
	for i := range f.Shots {
		sc := &f.Shots[i]
		if t < sc.A-0.001 || t > sc.B+0.001 {
			continue
		}
		fr := &Frame{
			T: t, W: f.W, H: f.H,
			LocalT: t - sc.A, LocalD: sc.B - sc.A,
			P: filmkit.Clamp01((t - sc.A) / (sc.B - sc.A)),
		}
		fr.Ease = filmkit.EaseOut(fr.P)

		// 几何转场参数（Wipe/Slide/Iris 的入幕包络恒 1，靠几何揭示）。
		inP := filmkit.Seg(t, sc.A, sc.A+f.Fade)
		offX := float32(0)
		clipW := float32(-1) // -1 = 不裁
		irisR := float32(-1)
		switch sc.In {
		case Wipe:
			clipW = float32(inP * float64(f.W))
		case Iris:
			irisR = float32(inP*math.Hypot(float64(f.W)/2, float64(f.H)/2)) + 1
		case Slide:
			offX = float32((1 - inP) * float64(f.W))
		}

		// 进出包络：Cross/Cut 是两端各自独立；经底色的转场把
		// 入场压到重叠区的后半程，底色扫瞙盖住接缝（见下）。
		var aIn float64
		switch sc.In {
		case Cut, Wipe, Slide, Iris:
			aIn = 1
		case FadeBlack, FadeWhite:
			aIn = filmkit.Seg(t, sc.A+f.Fade/2, sc.A+f.Fade)
		default:
			aIn = filmkit.Seg(t, sc.A, sc.A+f.Fade)
		}
		var aOut float64
		switch {
		case i == len(f.Shots)-1 && sc.B >= last-0.001 && sc.Out != FadeBlack && sc.Out != FadeWhite:
			// 最后一幕不淡出——除非声明了淡出到底色。
			aOut = 1
		case sc.Out == Cut:
			aOut = 1
		default:
			aOut = 1 - filmkit.Seg(t, sc.B-f.Fade, sc.B)
		}
		if sc.Out == FadeBlack || sc.Out == FadeWhite {
			// 片尾淡出：在最后一幕的尾部让位给底色。
			if sc.B >= last-0.001 {
				aOut = 1 - filmkit.Seg(t, sc.B-f.Fade/2, sc.B)
			}
		}
		// 下一幕是几何转场（Wipe/Iris 盖住、Slide 推走）时，
		// 出幕保持满 alpha——遮盖/推移本身就是出场方式。
		if i+1 < len(f.Shots) {
			switch f.Shots[i+1].In {
			case Wipe, Iris, Slide:
				if t >= f.Shots[i+1].A-0.001 {
					aOut = 1
				}
			}
		}

		a := aIn * aOut
		if a <= 0.0001 {
			continue
		}
		curAlpha = float32(a)
		if offX != 0 || clipW >= 0 || irisR >= 0 {
			d.Save()
			if offX != 0 {
				d.Translate(offX, 0)
			}
			if clipW >= 0 {
				d.BeginPath()
				d.Rect(0, 0, clipW, f.H)
				d.Clip()
			}
			if irisR >= 0 {
				d.BeginPath()
				d.Arc(f.W/2, f.H/2, irisR, 0, float32(2*math.Pi), false)
				d.Clip()
			}
			sc.Fn(d, fr)
			d.Restore()
		} else {
			sc.Fn(d, fr)
		}
	}
	curAlpha = 1

	// 转场扫瞙：经底色的转场在重叠区的前半程把底色盖上来、
	// 后半程揭下去；片尾淡出则只升不降。
	fadeScrims(d, t, f)

	drawSubs(d, t, f)

	d.Restore()
}

// fadeScrims 画经底色转场的底色层。
func fadeScrims(d *core.Drawing, t float64, f Film) {
	for i := range f.Shots {
		sc := &f.Shots[i]
		var white bool
		var a float64
		switch {
		case sc.In == FadeWhite:
			white = true
			a = filmkit.Seg(t, sc.A, sc.A+f.Fade/2) * (1 - filmkit.Seg(t, sc.A+f.Fade/2, sc.A+f.Fade))
		case sc.In == FadeBlack:
			a = filmkit.Seg(t, sc.A, sc.A+f.Fade/2) * (1 - filmkit.Seg(t, sc.A+f.Fade/2, sc.A+f.Fade))
		case sc.Out == FadeWhite && sc.B >= f.Duration-0.001:
			white = true
			a = filmkit.Seg(t, sc.B-f.Fade/2, sc.B)
		case sc.Out == FadeBlack && sc.B >= f.Duration-0.001:
			a = filmkit.Seg(t, sc.B-f.Fade/2, sc.B)
		default:
			continue
		}
		if a <= 0.0001 {
			continue
		}
		c := f.Bg
		if white {
			c = core.RGB(0xffffff)
		}
		d.SetFillColor(c)
		d.SetAlpha(float32(a))
		d.BeginPath()
		d.Rect(0, 0, f.W, f.H)
		d.Fill()
		d.SetAlpha(1)
	}
}
