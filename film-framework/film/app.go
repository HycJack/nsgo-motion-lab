// Package film 把一部 NSGo 影片的「骨架」打包成声明式配置。
//
// 一部片子 = 分镜表 + 字幕表 + 一条音频轨。引擎负责其余的一切：
// 时钟（录制器驱动 or 声音主时钟）、幕与幕之间的交叉溶解、
// 藏在画布后面的音频载体、「开始播放」按钮。
//
// 分镜函数只做一件事：给定时间，画出那一帧。
package film

import (
	"time"

	"github.com/oneo-me/nsgo"
	"github.com/oneo-me/nsgo/base"
	"github.com/oneo-me/nsgo/core"

	"nsgofilm/filmkit"
)

// Film 是一部片子的全部声明。
type Film struct {
	W, H     float32    // 设计画幅（如 1920x1080），实际窗口不同时整体缩放
	Duration float64    // 全片时长（秒）
	Fade     float64    // 幕与幕之间交叉溶解的时长
	Bg       core.Color // 画布底色（幕间露出的那层）
	Audio    string     // 音频载体的 asset: 路径（只带声音的 mp4）；空 = 无声片
	Button   string     // 「开始播放」按钮文案；空 = 「开始播放」
	Shots    []Shot     // 分镜表：相邻两幕重叠 Fade 秒即交叉溶解
	Subs     []Sub      // 字幕表（屏幕空间，不受镜头影响）
	SubStyle SubStyle   // 字幕样式；用 DefaultSubStyle 生成
	Sfx      []Sfx      // 音效表：影片时钟走到 T 时响起（录制模式不出声）
	FPS      float64    // 审片逐帧步进的帧率，0 = 30
}

// Run 用默认窗口大小（即设计画幅）运行这部片子。
func Run(f Film) {
	nsgo.Run(root(f), nsgo.WindowSize(core.Points(f.W), core.Points(f.H)))
}

// New 返回一个 NSGo Root 函数；想自己决定窗口大小时用这个。
func New(f Film) func(*core.Scope) core.Node { return root(f) }

// root 搭出整个画面：画布 + 藏起来的音频 + 「开始播放」按钮。
func root(f Film) func(*core.Scope) core.Node {
	return func(cx *core.Scope) core.Node {
		f.resolve() // Duration 未填时取最晚的幕尾

		filmT := core.NewSignal(cx, 0.0)                   // 影片时间（秒），画布按它画
		mediaT := core.NewSignal(cx, 0.0)                  // 声音的播放位置（秒）
		started := core.NewSignal(cx, false)               // 观众点过「开始」了吗
		playing := core.NewSignal(cx, false)               // 审片模式：正在播吗
		rate := core.NewSignal(cx, 1.0)                    // 审片变速：0.5 / 1 / 2
		recording := recordingMode()                       // 录制器在驱动吗
		audit := auditMode() && !recording                 // 审片模式（URL ?audit=1）
		timeline := timelineMode() && !recording && !audit // 时间轴总览（?timeline=1）
		review := audit || timeline                        // 审片行为：两种入口共用
		inFilm := core.NewSignal(cx, !timeline)            // 时间轴总览 → 影片视图

		// 审片控制条的显隐，照常用播放器的规矩：暂停时常驻；
		// 播放中 2.5 秒没动静就藏起来（把字幕让出来），点画面、
		// 动鼠标或碰控制条任意一处都会唤出。
		var wakeAt time.Time
		wake := func() { wakeAt = time.Now() }
		barVisible := core.NewSignal(cx, true)

		// 音效槽：每条 Sfx 一个藏起来的媒体元素，PausedFunc 由信号
		// 驱动。帧循环检测影片时钟「跨过 T」的沿触发；时钟倒回
		// （拖进度条、重播）就停掉归零，允许再次触发。
		type sfxSlot struct {
			s    Sfx
			play *core.Signal[bool]
			ref  *core.Ref
		}
		var slots []*sfxSlot
		for _, sx := range f.Sfx {
			slots = append(slots, &sfxSlot{s: sx, play: core.NewSignal(cx, false), ref: &core.Ref{}})
		}

		var loop func(time.Time)
		var prevWall float64
		var prevT float64
		loop = func(time.Time) {
			t := filmkit.Now()
			switch {
			case recording:
				// 时钟归录制器管。
			case review && f.Audio != "" && playing.Get():
				// 声音当主时钟：画面跟着声音走。
				t = mediaT.Get()
			case review && f.Audio != "":
				// 审片暂停中：跳转直接写 filmT，时钟停在原地。
				t = filmT.Get()
			case review && playing.Get():
				// 无声片播放：从当前位置按墙钟推进（乘审片变速），播完停住。
				now := filmkit.Now()
				t = filmT.Get() + min(max(now-prevWall, 0), 0.1)*rate.Get()
				if t >= f.Duration {
					t = f.Duration
					playing.Set(false)
				}
			case review:
				// 无声片暂停中。
				t = filmT.Get()
			case f.Audio != "" && started.Get():
				// 声音当主时钟：画面跟着声音走。
				t = mediaT.Get()
			case f.Audio != "":
				// 有声片还没开始，停在片头。
				t = 0
			default:
				// 无声片按墙钟走。
			}
			prevWall = filmkit.Now()
			filmT.Set(t)
			// 音效：检测影片时钟跨过 T 的沿（录制模式不出声）。
			if !recording {
				for _, sl := range slots {
					switch {
					case !sl.play.Get() && prevT < sl.s.T && t >= sl.s.T:
						sl.play.Set(true)
					case sl.play.Get() && t < sl.s.T:
						// 时钟倒回去了：停掉归零，允许再次触发。
						sl.play.Set(false)
						sl.ref.Seek(0)
					}
				}
			}
			prevT = t
			if review {
				barVisible.Set(!playing.Get() || time.Since(wakeAt) < 2500*time.Millisecond)
			}
			cx.RequestAnimationFrame(loop)
		}
		cx.RequestAnimationFrame(loop)

		canvas := core.Canvas{
			Style: core.Style{Width: core.Pct(100), Height: core.Pct(100), Background: f.Bg},
			Draw: func(d *core.Drawing, size core.Size) {
				render(d, size, filmT.Get(), f)
			},
		}

		children := []core.Node{canvas}

		// 音效的媒体元素：信号已在上面建好，元素在 children 就位。
		for _, sl := range slots {
			gain := sl.s.Gain
			if gain <= 0 {
				gain = 1
			}
			children = append(children, core.Video{
				Ref:        sl.ref,
				Source:     sl.s.Asset,
				PausedFunc: func() bool { return !sl.play.Get() },
				Volume:     gain,
				Style: core.Style{
					Position: core.Absolute, Left: core.Pt(0), Top: core.Pt(0),
					Width: core.Pt(2), Height: core.Pt(2),
					Opacity: core.Alpha(0),
				},
			})
		}

		// 审片模式的唤出层：铺满画布的透明视图，点一下、动一下
		// 鼠标、摸一下屏幕都会把控制条唤出来。垫在控制条底下，
		// 控制条自己的交互不会被它截走。
		if review {
			children = append(children, core.View{
				Style: core.Style{
					Position: core.Absolute, Left: core.Pt(0), Right: core.Pt(0),
					Top: core.Pt(0), Bottom: core.Pt(0),
				},
				OnPress:     func(core.PressEvent) { wake() },
				OnHoverMove: func(core.HoverEvent) { wake() },
				OnTouch:     func(core.TouchEvent) { wake() },
			})
		}

		// 音频载体：一条只带声音的 mp4，画面轨是 64×64 的黑场。
		// 藏在画布后面（2×2、全透明），唯一的用途是把播放位置喂给画面时钟。
		var mediaRef *core.Ref
		if f.Audio != "" {
			mediaRef = &core.Ref{}
			paused := func() bool { return !started.Get() }
			if review {
				paused = func() bool { return !playing.Get() }
			}
			children = append(children, core.Video{
				Ref:              mediaRef,
				Source:           f.Audio,
				PausedFunc:       paused,
				RateFunc:         func() float32 { return float32(rate.Get()) },
				Volume:           1,
				ProgressInterval: 16 * time.Millisecond,
				OnProgress:       func(pos, dur time.Duration) { mediaT.Set(pos.Seconds()) },
				Style: core.Style{
					Position: core.Absolute, Left: core.Pt(0), Top: core.Pt(0),
					Width: core.Pt(2), Height: core.Pt(2),
					Opacity: core.Alpha(0),
				},
			})

			// 「开始播放」按钮只在普通观看模式出现；点了就从头播。
			// 审片模式用底部的控制条代替。
			if !recording && !audit {
				label := f.Button
				if label == "" {
					label = "开始播放"
				}
				children = append(children, core.Show{
					When: func() bool { return !started.Get() },
					Then: func(cx *core.Scope) core.Node {
						return core.View{
							Style: core.Style{
								Position: core.Absolute, Left: core.Pt(0), Right: core.Pt(0),
								Top: core.Pt(0), Bottom: core.Pt(0),
								AlignItems: core.AlignCenter, JustifyContent: core.JustifyEnd,
								Padding: core.All(90),
							},
							Children: []core.Node{
								base.Button{
									Label: label,
									OnPress: func(core.PressEvent) {
										mediaRef.Seek(0)
										started.Set(true)
									},
									Style: core.Style{
										Background:   core.RGB(0xf2b544),
										CornerRadius: 30,
										Padding:      core.Symmetric(16, 44),
									},
									Children: []core.Node{
										core.Text{
											Content: label,
											Style: core.Style{
												Color:    core.RGB(0x241605),
												FontSize: core.Points(26),
											},
										},
									},
								},
							},
						}
					},
				})
			}
		}

		// 审片控制条：播放/暂停、时间、带声音跳转的进度条。
		// 显隐由 barVisible 驱动；出现期间会盖住字幕（常用播放器也这样），
		// 播放中静置 2.5 秒就自动让位。
		if review {
			children = append(children, core.Show{
				When: barVisible.Get,
				Then: func(cx *core.Scope) core.Node {
					return auditBar(cx, f, filmT, playing, mediaRef, wake)
				},
			})
		}

		// 审片键盘：空格播放/暂停，←→ 逐帧，Shift+←→ ±1s，
		// J 0.5× / L 1× / Shift+L 2×（K = 暂停）。
		if review {
			fps := f.FPS
			if fps <= 0 {
				fps = 30
			}
			step := func(dt float64) {
				playing.Set(false)
				wake()
				nt := filmkit.Clamp01((filmT.Get()+dt)/max(f.Duration, 0.001)) * max(f.Duration, 0.001)
				filmT.Set(nt)
				if mediaRef != nil {
					mediaRef.Seek(time.Duration(nt * float64(time.Second)))
				}
			}
			toggle := func() {
				wake()
				if playing.Get() {
					playing.Set(false)
					return
				}
				if filmT.GetUntracked() >= f.Duration-0.01 {
					filmT.Set(0)
					if mediaRef != nil {
						mediaRef.Seek(0)
					}
				}
				playing.Set(true)
			}
			setRate := func(r float64) {
				wake()
				rate.Set(r)
				playing.Set(true)
			}
			// 方向键的名字两端不同：web 归一化为 left/right，native 用
			// arrowleft/arrowright——两套都注册，重复注册无害。
			cx.OnKey("space", toggle)
			cx.OnKey("right", func() { step(1 / fps) })
			cx.OnKey("left", func() { step(-1 / fps) })
			cx.OnKey("shift+right", func() { step(1) })
			cx.OnKey("shift+left", func() { step(-1) })
			cx.OnKey("arrowright", func() { step(1 / fps) })
			cx.OnKey("arrowleft", func() { step(-1 / fps) })
			cx.OnKey("shift+arrowright", func() { step(1) })
			cx.OnKey("shift+arrowleft", func() { step(-1) })
			cx.OnKey("l", func() { setRate(1) })
			cx.OnKey("shift+l", func() { setRate(2) })
			cx.OnKey("j", func() { setRate(0.5) })
			cx.OnKey("k", func() { playing.Set(false) })
		}

		// 时间轴总览：全屏盖一层每幕一帧的卡片网格，
		// 点卡片跳进审片模式的那一幕（时钟与媒体都 seek 到幕首）。
		if timeline {
			children = append(children, core.Show{
				When: func() bool { return !inFilm.Get() },
				Then: func(cx *core.Scope) core.Node {
					return timelineGrid(cx, f, filmT, inFilm, wake)
				},
			})
		}

		return core.View{
			Style:    core.Style{Width: core.Pct(100), Height: core.Pct(100)},
			Children: children,
		}
	}
}
