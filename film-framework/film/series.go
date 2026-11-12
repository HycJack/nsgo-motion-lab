package film

import (
	"github.com/oneo-me/nsgo/core"
)

// 分镜的装配层：Series 自动排布、Seq 场景局部时间、转场声明。
//
// 设计约束：所有新字段零值 = 现状——Shot{A, B, Fn} 的老工程一行不用改。

// Transition 是幕的进/出场形态（转场）。零值 Cross 就是现状的交叉溶解。
type Transition uint8

const (
	// Cross 是交叉溶解：与相邻幕重叠 fade 秒，透明度互为反相。
	Cross Transition = iota
	// Cut 是硬切：进场的幕在 A 时刻瞬间满alpha（叠在出幕之上，盖住即切）。
	Cut
	// FadeBlack 经底色过渡：前半程出幕淡出+底色浮现，后半程入幕从底色淡入。
	// 底色就是 Film.Bg——把 Bg 设成黑就是黑场过渡。
	FadeBlack
	// FadeWhite 同上，白色闪场。
	FadeWhite
	// Wipe 是擦除：入幕被一个从左向右扩张的矩形裁剪，逐渐盖住出幕。
	Wipe
	// Slide 是推挤：入幕从右侧滑入，出幕同时向左推出（画中露出 Film.Bg）。
	Slide
	// Iris 是圆形揭示：入幕被一个从屏幕中心扩张的圆裁剪，逐渐盖住出幕。
	Iris
)

// Series 把按「时长」声明的幕自动排成时间轴：第一幕从 0 起，
// 相邻幕重叠 fade 秒（前搭后）。时长改了全片自动重排，
// 不再手写 A/B、不用手动和 Film.Fade 保持一致。
//
// 每幕的 In/Out（转场）原样保留；幕里 Fn 之外的字段也是。
func Series(fade float64, shots ...Shot) []Shot {
	out := make([]Shot, 0, len(shots))
	var next float64 // 下一幕的起点（前一份钟表的 B）
	for i, s := range shots {
		if s.D <= 0 {
			continue
		}
		if i == 0 {
			s.A = 0
		} else {
			s.A = next - fade
		}
		s.B = s.A + s.D
		next = s.B
		out = append(out, s)
	}
	return out
}

// Seq 把场景函数放到 start 时刻，场景内部从 0 计时——
// Frame.LocalT 是场景内时间、LocalD 是场景时长。场景因此可复用、
// 可整体移动，不必把所有坐标挂在全片绝对时间上。
//
// 本质就是一枚普通 Shot：LocalT 恒等于 f.T - A，对每一幕都成立。
func Seq(start, dur float64, scene func(d *core.Drawing, f *Frame)) Shot {
	return Shot{A: start, B: start + dur, Fn: scene}
}
