// 《Take Demo》—— 一部由节点分镜搭起来的影片。
//
// 底图是一幕 canvas 分镜，三层字卡是 keyframe.Take：两种分镜读同一个时钟
// （film.Run 挂载时把它发布出去），所以暂停、拖进度、审片逐帧一起生效。
// 录成 mp4 的还是同一个录制器，不必知道画面是怎么来的。
//
// 三层字卡各占一小段时间窗，窗口之间留 0.3 秒重叠 —— 那是「切」和「溶」的
// 位置，写在各层自己的 opacity 轨道里，不必改动分镜表。
package main

import (
	"math"
	"strconv"
	"time"

	"github.com/oneo-me/nsgo/core"

	"nsgofilm/film"
	"nsgofilm/keyframe"
)

const (
	dw float32 = 1920
	dh float32 = 1080
	ms         = time.Millisecond
)

var (
	ink   = core.RGB(0x070a10)
	grid  = core.RGBA(0x18222f, 0.85)
	amber = core.RGB(0xf0b429)
	text  = core.RGB(0xeef3fa)
	muted = core.RGB(0x8ea2b8)
	panel = core.RGBA(0x0d1521, 0.9)
	line  = core.RGBA(0x33465c, 0.9)
)

func main() {
	title, cards, outro := beats()

	film.Run(film.Film{
		W:        dw,
		H:        dh,
		Duration: 21,
		Fade:     0.6,
		Bg:       ink,
		Shots:    film.Series(0.6, film.Shot{D: 21, Fn: backdrop}),
		Layers: []film.Layer{
			film.Over(title, 0.4, 7.0),
			film.Over(cards, 7.0, 14.6),
			film.Over(outro, 14.6, 21),
		},
		Subs: []film.Sub{
			{A: 1.2, B: 6.4, Text: "分镜写在 Go 里，时间用秒"},
			{A: 7.6, B: 14.0, Text: "每一帧都是一次求值，不是一段正在跑的动画"},
			{A: 15.0, B: 20.4, Text: "同一个时钟：暂停、拖进度、逐帧，画面跟着走"},
		},
		SubStyle: film.DefaultSubStyle(text, amber, core.RGB(0x04060a)),
	})
}

// backdrop 是全片的底图：一片极淡的竖网格横移，两道缓慢漂移的光。
// 它是普通 canvas 分镜 —— 和字卡一样只读 f.T。
func backdrop(d *core.Drawing, f *film.Frame) {
	d.SetFillColor(ink)
	d.FillRect(0, 0, dw, dh)

	shift := float32(math.Mod(f.T*14, 96))
	d.SetStrokeColor(grid)
	d.SetLineWidth(1)
	for x := shift - 96; x < dw+96; x += 96 {
		d.BeginPath()
		d.MoveTo(x, 0)
		d.LineTo(x, dh)
		d.Stroke()
	}

	// 光斑在缓慢呼吸：两个不同周期，所以不会同时回到原位。
	film.Glow(d, 1460+float32(math.Sin(f.T*0.30))*160, 300, 640, core.RGB(0x1d4a76), 0.55)
	film.Glow(d, 250+float32(math.Cos(f.T*0.19))*90, 880, 560, core.RGB(0x2a6b8f), 0.26)
	film.Vignette(d, dw, dh, 0.6)
}

func beats() (title, cards, outro *keyframe.Take) {
	return beatTitle(), beatCards(), beatOutro()
}

// beatTitle 是开场：一行小标题、两行大字、一条短线。
// 每条轨道都用 Win 收进出场，所以把窗口挪到哪里都不用重写。
func beatTitle() *keyframe.Take {
	take := keyframe.New("title", 6600*ms)

	take.Text("eyebrow", "NSGO · FILM-FRAMEWORK").Font(26).Bold().Color(amber).
		At(240, 358).
		On("opacity", keyframe.Win(150*ms, 6200*ms, 400*ms)).
		On("x", keyframe.KeysAt(150*ms, 216, 900*ms, 240).Ease(keyframe.EaseOut))

	// 两行大字错开 90 毫秒落下：一起到像是整块贴上来的，错开会读成「写上去」。
	take.Text("l1", "一张时间表").Font(104).Bold().Color(text).
		At(240, 452).
		On("opacity", keyframe.Win(300*ms, 6200*ms, 450*ms)).
		On("y", keyframe.KeysAt(300*ms, 470, 1100*ms, 452).Ease(keyframe.EaseOut))
	take.Text("l2", "就是一支片子的剧本").Font(104).Bold().Color(text).
		At(240, 596).
		On("opacity", keyframe.Win(390*ms, 6200*ms, 450*ms)).
		On("y", keyframe.KeysAt(390*ms, 614, 1190*ms, 596).Ease(keyframe.EaseOut))

	take.Box("rule").At(244, 736).Size(200, 3).Fill(amber).
		On("scaleX", keyframe.KeysAt(900*ms, 0, 1500*ms, 1).Ease(keyframe.EaseOut)).
		On("opacity", keyframe.Win(900*ms, 6200*ms, 300*ms))

	take.Text("note", "keyframe.Take · 一秒一个数，没有 0–1 的进度").Font(30).Color(muted).
		At(240, 782).
		On("opacity", keyframe.Win(1500*ms, 6200*ms, 500*ms))

	return take
}

// beatCards 是三张卡：分镜、时钟、导出 —— 一部片子要交代的三件事。
// 卡是 Group，错开 160 毫秒从下面升上来；进场之后卡内部不再动。
//
// 两处尺寸上有讲究：
//   - 组上的 y 是「整组位移」，从 0 涨到 24 才是往上飘，写死成目标的 396
//     会把整张卡推到画面底部；
//   - 字节点量自己的宽度，不会换行，所以正文是两行、两行各一个 Text。
func beatCards() *keyframe.Take {
	take := keyframe.New("cards", 7600*ms)

	take.Text("head", "三件事，一层就够").Font(56).Bold().Color(text).
		At(240, 226).
		On("opacity", keyframe.Win(0, 7200*ms, 400*ms)).
		On("y", keyframe.KeysAt(0*ms, 244, 700*ms, 226).Ease(keyframe.EaseOut))

	cards := []struct {
		num, name string
		lines     []string
	}{
		{"01", "分镜", []string{"一幕一个函数", "几何是算出来的"}},
		{"02", "时钟", []string{"全片一个播放头", "暂停、拖动、逐帧都读它"}},
		{"03", "导出", []string{"逐帧设时间再截图", "帧率无关，重跑也一样"}},
	}

	for i, c := range cards {
		start := 350*ms + time.Duration(i)*160*ms
		x := 240 + float64(i)*500
		g := take.Group("card" + c.num)
		g.On("opacity", keyframe.Win(start, 7300*ms, 420*ms))
		g.On("y", keyframe.KeysAt(start, 26, start+800*ms, 0).Ease(keyframe.EaseOut))

		g.Add(keyframe.Box, "plate").At(x, 330).Size(440, 300).
			Radius(20).Fill(panel)
		g.Add(keyframe.Box, "edge").At(x, 330).Size(440, 3).
			Radius(2).Fill(amber)
		g.Add(keyframe.Text, "num").Text(c.num).Font(26).Bold().Color(amber).
			At(x+34, 362)
		g.Add(keyframe.Text, "name").Text(c.name).Font(46).Bold().Color(text).
			At(x+34, 400)
		for j, ln := range c.lines {
			g.Add(keyframe.Text, "line"+strconv.Itoa(j)).Text(ln).Font(24).Color(muted).
				At(x+34, 500+float64(j)*38)
		}
		g.Add(keyframe.Box, "hair").At(x+34, 580).Size(372, 1).Fill(line)
	}

	return take
}

// beatOutro 收尾：一行结论、一条路径。整层在最后 0.8 秒淡出，
// 让片尾的黑场是黑的（录制器拍到的是最后一帧，不是最后一秒）。
func beatOutro() *keyframe.Take {
	take := keyframe.New("outro", 6400*ms)

	take.Text("l1", "改一行，重录一遍").Font(88).Bold().Color(text).
		At(240, 452).
		On("opacity", keyframe.Win(0, 6000*ms, 450*ms)).
		On("y", keyframe.KeysAt(0*ms, 470, 800*ms, 452).Ease(keyframe.EaseOut))

	take.Box("rule").At(244, 610).Size(200, 3).Fill(amber).
		On("scaleX", keyframe.KeysAt(400*ms, 0, 1100*ms, 1).Ease(keyframe.EaseOut)).
		On("opacity", keyframe.Win(400*ms, 6000*ms, 300*ms))

	take.Text("path", "film-framework/keyframe  ·  shots/film.js").Font(28).Color(muted).
		At(240, 660).
		On("opacity", keyframe.Win(900*ms, 6000*ms, 500*ms))

	return take
}
