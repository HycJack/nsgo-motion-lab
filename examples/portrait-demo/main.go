// 竖屏片（1080×1920，10 秒）：层、窗口、全片层。
//
// 和 take-demo 的差别不是画幅，而是三种层同时存在：
//
//	rail   —— 窗口两端都是 0，铺满全片（一条走到底的进度轨）
//	hero   —— 0.4–4.6 秒的个人秀
//	stats  —— 4.6–7.6 秒，三行数字先后进场
//	end    —— 7.6–10 秒，落版
//
// 画布分镜（warm/cool）则自己交叉溶解一次，于是「画布换幕」和「字卡换层」
// 是两件独立的事，各按各的时间走——这正是 Layers 存在的理由。
package main

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/nsgo-dev/nsgo/core"

	"nsgofilm/film"
	"nsgofilm/keyframe"
)

const (
	dw, dh = 1080.0, 1920.0 // 设计画幅：竖屏
	ms     = time.Millisecond
)

var (
	ink   = core.RGB(0x0e0b09)
	sand  = core.RGB(0xf4e9d8)
	amber = core.RGB(0xe0a13c)
	clay  = core.RGB(0x8f4a2e)
	muted = core.RGB(0x9b8b78)
	line  = core.RGBA(0x4a3a2c, 0.9)
)

func main() {
	rail, hero, stats, end := beats()

	film.Run(film.Film{
		W:        dw,
		H:        dh,
		Duration: 10,
		Fade:     0.5,
		Bg:       ink,
		Shots: film.Series(0.5,
			film.Shot{D: 5.5, Fn: warm}, // 0 – 5.5
			film.Shot{D: 5.0, Fn: cool}, // 5 – 10（和上一幕叠 0.5 秒交叉溶解）
		),
		Layers: []film.Layer{
			film.Over(rail, 0, 0), // 两个 0 = 铺满全片
			film.Over(hero, 0.4, 4.6),
			film.Over(stats, 4.6, 7.6),
			film.Over(end, 7.6, 10),
		},
		Subs: []film.Sub{
			{A: 0.9, B: 4.2, Text: "纸轻了，字还在"},
			{A: 4.9, B: 7.3, Text: "同样的 240 页，重量少了一半"},
			{A: 7.9, B: 9.7, Text: "PAPER 轻装本 · 现已上市"},
		},
		// 竖屏的字幕基线在 1748，不是 1080p 的 952——SubStyle 就是为这个存在的。
		SubStyle: film.SubStyle{
			Font:       "Noto Sans CJK SC",
			Size:       44,
			Y:          1748,
			Narration:  sand,
			Dialogue:   amber,
			Plate:      ink,
			PlateAlpha: 0.55,
		},
	})
}

func beats() (rail, hero, stats, end *keyframe.Take) {
	return beatRail(), beatHero(), beatStats(), beatEnd()
}

// warm 是前半段的底：横格纸的纹路极慢下移，右下角一团暖光在长大。
func warm(d *core.Drawing, f *film.Frame) {
	d.SetFillColor(ink)
	d.FillRect(0, 0, dw, dh)

	shift := float32(math.Mod(f.T*9, 96))
	d.SetStrokeColor(core.RGBA(0x33261d, 0.6))
	d.SetLineWidth(1)
	for y := shift - 96; y < dh+96; y += 96 {
		d.BeginPath()
		d.MoveTo(0, y)
		d.LineTo(dw, y)
		d.Stroke()
	}

	// 竖屏的下半幅要靠光撑住，否则字在上三分之一、下面全空。
	film.Glow(d, 300, 1300, 560+float32(f.T)*30, clay, 0.62)
	film.Glow(d, 820, 1620, 380, clay, 0.3)
	film.Glow(d, 880, 480, 300, amber, 0.16)
	film.Vignette(d, dw, dh, 0.6)
}

// cool 是后半段的底：三道同心圆慢慢呼气（各差 0.9 弧度，所以永不同步）。
func cool(d *core.Drawing, f *film.Frame) {
	d.SetFillColor(ink)
	d.FillRect(0, 0, dw, dh)

	d.SetLineWidth(2)
	for i := 0; i < 3; i++ {
		r := 210 + float32(i)*150 + float32(math.Sin(f.T*0.55+float64(i)*0.9))*26
		d.SetStrokeColor(core.RGBA(0x6a4a2c, 0.4))
		d.BeginPath()
		d.Arc(dw/2, 1180, r, 0, float32(2*math.Pi), false)
		d.Stroke()
	}
	film.Glow(d, dw/2, 1180, 460, clay, 0.32)
	film.Vignette(d, dw, dh, 0.58)
}

// beatRail 是唯一铺满全片的层：一条轨道 + 三个刻度 + 一颗匀速下行的点。
// 它的 0 就是影片的 0（Delay 为 0），所以缩略图上它从头走到尾。
func beatRail() *keyframe.Take {
	take := keyframe.New("rail", 10*time.Second)

	take.Box("track").At(938, 560).Size(2, 800).Fill(line)
	for i := 0; i < 3; i++ {
		y := 560 + float64(i)*400
		take.Box("tick"+strconv.Itoa(i)).At(928, y).Size(22, 2).Fill(line)
		take.Text("no"+strconv.Itoa(i), fmt.Sprintf("%02d", i+1)).Font(24).Color(muted).
			At(966, y-15).
			On("opacity", keyframe.Win(300*ms+time.Duration(i)*220*ms, 9600*ms, 420*ms))
	}
	take.Disc("dot").At(931, 560).Size(16, 16).Fill(amber).
		On("y", keyframe.Keys(0, 560, 10, 1360)).
		On("opacity", keyframe.Win(0*ms, 9600*ms, 300*ms))
	return take
}

// beatHero 是开场三行字：小标、两行大字、一条短线、一行数据。
// 每条轨道都用 Win 收进出场，窗口怎么挪都不用改。
func beatHero() *keyframe.Take {
	take := keyframe.New("hero", 4200*ms)

	take.Text("eyebrow", "PAPER · 轻装本").Font(30).Bold().Color(amber).
		At(96, 520).
		On("opacity", keyframe.Win(100*ms, 3800*ms, 320*ms)).
		On("x", keyframe.KeysAt(100*ms, 72, 700*ms, 96).Ease(keyframe.EaseOut))

	for i, s := range []string{"轻一点", "再轻一点"} {
		name := "l" + strconv.Itoa(i+1)
		start := 180*ms + time.Duration(i)*100*ms
		take.Text(name, s).Font(154).Bold().Color(sand).
			At(96, 610+float64(i)*180).
			On("opacity", keyframe.Win(start, 3700*ms, 380*ms)).
			On("y", keyframe.KeysAt(start, 632+float64(i)*180, start+800*ms, 610+float64(i)*180).
				Ease(keyframe.EaseOut))
	}

	take.Box("rule").At(100, 1010).Size(180, 4).Fill(amber).
		On("scaleX", keyframe.KeysAt(700*ms, 0, 1200*ms, 1).Ease(keyframe.EaseOut)).
		On("opacity", keyframe.Win(680*ms, 3600*ms, 260*ms))

	take.Text("cap", "240 页，580 g → 242 g").Font(36).Color(muted).
		At(96, 1060).
		On("opacity", keyframe.Win(1000*ms, 3600*ms, 400*ms)).
		On("y", keyframe.KeysAt(1000*ms, 1076, 1500*ms, 1060).Ease(keyframe.EaseOut))
	return take
}

// beatStats 是三行数字，一行一组：组管进退（y 与 opacity），行内只管自己的字。
func beatStats() *keyframe.Take {
	take := keyframe.New("stats", 3000*ms)

	take.Text("head", "三件事变了").Font(62).Bold().Color(sand).
		At(96, 470).
		On("opacity", keyframe.Win(0*ms, 2600*ms, 400*ms)).
		On("y", keyframe.KeysAt(0*ms, 494, 500*ms, 470).Ease(keyframe.EaseOut))

	rows := []struct{ num, label string }{
		{"−58%", "整本更轻"},
		{"1 次", "对折不裂"},
		{"120 年", "纸张寿命"},
	}
	for i, r := range rows {
		y := 640 + float64(i)*200
		start := 220*ms + time.Duration(i)*140*ms
		g := take.Group("row" + strconv.Itoa(i))
		g.On("opacity", keyframe.Win(start, 2700*ms, 380*ms)).
			On("y", keyframe.KeysAt(start, 28, start+700*ms, 0).Ease(keyframe.EaseOut))

		// 组本身是一个满幅合成箱：孩子用绝对坐标写，组只负责整体位移与淡出。
		g.Add(keyframe.Box, "num"+strconv.Itoa(i)).At(96, y).Size(4, 96).Fill(amber).
			On("scaleY", keyframe.KeysAt(start, 0, start+520*ms, 1).Ease(keyframe.EaseOut))
		g.Add(keyframe.Text, "val"+strconv.Itoa(i)).Text(r.num).Font(88).Bold().Color(sand).
			At(140, y+2)
		g.Add(keyframe.Text, "lab"+strconv.Itoa(i)).Text(r.label).Font(34).Color(muted).
			At(140, y+112)
		// 横线止于 880：再往右就穿进右侧的进度轨了。
		g.Add(keyframe.Box, "hair"+strconv.Itoa(i)).At(96, y+176).Size(784, 1).Fill(line)
	}
	return take
}

// beatEnd 是落版：一行品牌字、一条短线、一行渠道。
func beatEnd() *keyframe.Take {
	take := keyframe.New("end", 2400*ms)

	take.Text("brand", "PAPER").Font(40).Bold().Color(amber).
		At(96, 700).
		On("opacity", keyframe.Win(0*ms, 2100*ms, 360*ms))
	take.Text("l1", "轻装本").Font(120).Bold().Color(sand).
		At(96, 760).
		On("opacity", keyframe.Win(120*ms, 2100*ms, 400*ms)).
		On("y", keyframe.KeysAt(120*ms, 786, 780*ms, 760).Ease(keyframe.EaseOut))

	take.Box("rule").At(100, 950).Size(180, 4).Fill(amber).
		On("scaleX", keyframe.KeysAt(500*ms, 0, 1000*ms, 1).Ease(keyframe.EaseOut)).
		On("opacity", keyframe.Win(480*ms, 2100*ms, 300*ms))

	take.Text("where", "全渠道上市").Font(40).Color(muted).
		At(96, 1000).
		On("opacity", keyframe.Win(760*ms, 2100*ms, 400*ms)).
		On("y", keyframe.KeysAt(760*ms, 1016, 1260*ms, 1000).Ease(keyframe.EaseOut))
	return take
}
