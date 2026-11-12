package film

// 审片模式的控制条：播放/暂停、当前时间、可点可拖的进度条。
//
// 时间显示用裸秒数（12.34s），和分镜表、字幕表里的 A/B 以及
// frames.py 抓帧的时间戳同一个口径，审片时直接对得上。
//
// 结构上把「会动的」和「接手指的」分开：进度条的轨道是静态节点，
// Responder 挂在它身上，拖动全程不被重建打断；填色、圆钮、时间字
// 用 core.Dynamic 按信号重建，只在视觉层换。

import (
	"fmt"
	"time"

	"github.com/oneo-me/nsgo/base"
	"github.com/oneo-me/nsgo/core"

	"nsgofilm/filmkit"
)

const (
	auditTrackW = 560.0 // 轨道宽（pt）；跳转位置按它归一化
	auditHitH   = 44.0  // 命中区高：比轨道本身高，好按
	auditRailH  = 6.0   // 轨道高
)

// auditBar 搭出审片控制条。filmT 是影片时钟，playing 是播放状态，
// media 是音频载体（无声片为 nil，进度条只推画面时钟），wake 报告
// 「有动静」——碰控制条也算，播放中的自动倒计时重新来过。
func auditBar(cx *core.Scope, f Film, filmT *core.Signal[float64], playing *core.Signal[bool], media *core.Ref, wake func()) core.Node {
	seek := func(t float64) {
		if f.Duration <= 0 {
			return
		}
		t = filmkit.Clamp01(t/f.Duration) * f.Duration
		if media != nil {
			media.Seek(time.Duration(t * float64(time.Second)))
		}
		filmT.Set(t)
	}

	// 轨道：按下即认领，起点和拖动都跳。拖动时先暂停——
	// 审片的手势语义是「掐住片子找位置」，松手后自己按播放。
	// 只在按住期间扫动：杂散的 move（无按下的鼠标路过）不许拽动播放头。
	var scrubbing bool
	scrub := func(ev core.ResponderEvent) {
		wake()
		playing.Set(false)
		seek(filmkit.Clamp01(float64(ev.Position.X)/auditTrackW) * f.Duration)
	}
	track := core.View{
		Style: core.Style{
			Width:  core.Pt(auditTrackW),
			Height: core.Pt(auditHitH),
		},
		Responder: &core.Responder{
			Claims:  []core.TouchClaim{{}},
			OnGrant: func(core.ResponderEvent) { scrubbing = true },
			OnStart: func(ev core.ResponderEvent) { scrubbing = true; scrub(ev) },
			OnMove: func(ev core.ResponderEvent) {
				if scrubbing {
					scrub(ev)
				}
			},
			OnRelease:   func(core.ResponderEvent) { scrubbing = false },
			OnTerminate: func(core.ResponderEvent) { scrubbing = false },
		},
		Children: []core.Node{
			// 轨道底
			core.View{Style: core.Style{
				Position: core.Absolute, Left: core.Pt(0), Right: core.Pt(0),
				Top: core.Pt((auditHitH - auditRailH) / 2), Height: core.Pt(auditRailH),
				Background:   core.RGBA(0xffffff, 0.22),
				CornerRadius: auditRailH / 2,
			}},
			// 已播部分 + 圆钮：跟着影片时钟走
			core.Dynamic(func(*core.Scope) core.Node {
				frac := filmkit.Clamp01(filmT.Get() / max(f.Duration, 0.001))
				return core.View{Style: core.Style{
					Position: core.Absolute, Left: core.Pt(0),
					Top:          core.Pt((auditHitH - auditRailH) / 2),
					Width:        core.Pt(core.Points(frac * auditTrackW)),
					Height:       core.Pt(auditRailH),
					Background:   core.RGB(0xf2b544),
					CornerRadius: auditRailH / 2,
				}}
			}),
			core.Dynamic(func(*core.Scope) core.Node {
				frac := filmkit.Clamp01(filmT.Get() / max(f.Duration, 0.001))
				const knob = 16.0
				return core.View{Style: core.Style{
					Position:     core.Absolute,
					Left:         core.Pt(core.Points(frac*auditTrackW - knob/2)),
					Top:          core.Pt((auditHitH - knob) / 2),
					Width:        core.Pt(knob),
					Height:       core.Pt(knob),
					Background:   core.RGB(0xf2b544),
					CornerRadius: knob / 2,
				}}
			}),
		},
	}

	// 播放/暂停：只在播放状态变化时重建。
	playBtn := core.Dynamic(func(*core.Scope) core.Node {
		label := "播放"
		if playing.Get() {
			label = "暂停"
		}
		return base.Button{
			Label: label,
			OnPress: func(core.PressEvent) {
				wake()
				if playing.Get() {
					playing.Set(false)
					return
				}
				// 播到头了再按播放 = 重看一遍。
				t := filmT.GetUntracked()
				if t >= f.Duration-0.01 {
					t = 0
				}
				// 起播前把媒体重新对齐到影片时钟：跳转命令若因媒体
				// 尚未加载而丢失，这里能补上；已在位上则无害。
				if media != nil {
					media.Seek(time.Duration(t * float64(time.Second)))
				}
				filmT.Set(t)
				playing.Set(true)
			},
			Style: core.Style{
				Background:   core.RGB(0xf2b544),
				CornerRadius: 18,
				Padding:      core.Symmetric(6, 22),
			},
			Children: []core.Node{
				core.Text{Content: label, Style: core.Style{
					Color:    core.RGB(0x241605),
					FontSize: core.Points(20),
				}},
			},
		}
	})

	// 时间：裸秒数，两位小数，和分镜表同一口径。
	clock := core.Dynamic(func(*core.Scope) core.Node {
		return core.Text{
			Content: fmt.Sprintf("%.2fs / %.2fs", filmT.Get(), f.Duration),
			Style: core.Style{
				Color:    core.RGBA(0xffffff, 0.85),
				FontSize: core.Points(20),
			},
		}
	})

	// 提示行 + 主行：竖排。
	row := core.View{
		Style: core.Style{
			Direction: core.Row, AlignItems: core.AlignCenter,
			Gap: core.Points(22),
		},
		Children: []core.Node{playBtn, clock, track},
	}
	return core.View{
		Style: core.Style{
			Position: core.Absolute, Left: core.Pt(0), Right: core.Pt(0), Bottom: core.Pt(0),
			Direction: core.Column, AlignItems: core.AlignCenter,
			Gap:          core.Points(8),
			Padding:      core.All(24),
			Background:   core.RGBA(0x0b0f14, 0.62),
			CornerRadius: 16,
		},
		OnHoverMove: func(core.HoverEvent) { wake() },
		Children: []core.Node{
			core.Text{Content: "空格 播放/暂停 · ← → 逐帧 · Shift+← → ±1s · J 0.5× · L 1× · Shift+L 2×",
				Style: core.Style{Color: core.RGBA(0xffffff, 0.55), FontSize: core.Points(18)}},
			row,
		},
	}
}
