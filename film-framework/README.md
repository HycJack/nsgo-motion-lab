# nsgo-film · 用 NSGo 拍片子的通用框架

把一部影片拆成三层：**引擎**（这个包）、**素材**（贴图/音频）、**分镜**
（一部片子的全部内容声明）。分镜函数是影片时钟的纯函数——同一个时刻画
出来的东西永远一样，所以录制的片子、观看的画面、审片的单帧，三者逐像素一致。

## 结构

```
filmkit/     影片时钟 + 时间数学（Seg / EaseOut / Win / Pulse …）
film/        引擎：Film 配置 → 时钟循环、交叉溶解、音频载体、字幕、开始按钮
keyframe/    节点分镜：场景树 + 关键帧轨道 → core.Node，能播、能定格、能自检
shots/film.js  录制器：推时钟 → CDP 逐帧截图 → ffmpeg 编码（WORKERS 并行、断点续录）
tools/         serve.py（Range 静态服务）/ frames.py（抓帧）/ eq.py（LaTeX→路径）/ mux.py（混音合成）/ diff.py（确定性回归）
../examples/   例子工程：math-demo（数学动画）/ cylinder-demo（圆柱展开教学片）/ style-demo（风格画廊）/ take-demo（层与窗口）/ portrait-demo（竖屏）…
../films/      例子录好的成片 mp4
```

两种分镜共用同一个时钟（`filmkit.Now()`），所以录制器不必知道画面是怎么来的。这个时钟按权威顺序取三个来源：录制器的 `window.__filmT` > 影片自己发布的播放头（`film.Run` 挂载时发布给 `filmkit`）> 墙钟。第三条是页面的常态，第二条让挂在一支片子里的节点分镜跟着影片暂停、拖进度、变速、审片逐帧走——不需要第二套钩子。

| | 分镜是什么 | 适合 |
| --- | --- | --- |
| `film.Shot` | `func(*core.Drawing, *Frame)`——一笔一笔画 | 教学片、图解、公式、几何动画 |
| `keyframe.Take` | 一棵节点树 + 每通道一条轨道——动的是 `AnimatedStyle` | 文字/卡片/图表/图层/shader 场景，要逐帧确定、要能定格截图 |

## 一部片子的最小写法

```go
package main

import "nsgofilm/film"

func main() {
    film.Run(film.Film{
        W: 1920, H: 1080,
        Duration: 30, Fade: 0.5,
        Bg:    core.RGB(0x0b1220),
        Audio: "asset:film.mp4",          // resources/ 下只带声音的 mp4；无声片给 ""
        Shots: []film.Shot{
            {0, 6, shotTitle},
            {5.5, 16, shotFall},          // 相邻两幕重叠 0.5s → 交叉溶解
        },
        Subs: []film.Sub{
            {6.8, 10.2, "天黑了。", false},
            {14.8, 17.6, "是星星呀。", true},
        },
        SubStyle: film.DefaultSubStyle(cream, gold, ink),
    })
}
```

分镜函数只管画：

```go
func shotFall(d *core.Drawing, f *film.Frame) {
    d.Save()
    f.Cam(d, 30-60*f.Ease, 0, 1.06)   // 取景：zoom ≥ 1，推近靠放大角色
    film.DrawPic(d, bg, 0, 0, f.W, f.H, 1)
    film.PicAt(d, star, 1120, groundY, 132, film.A(1))
    film.Glow(d, 960, 520, 320, gold, 0.34)
    film.Vignette(d, f.W, f.H, 0.38)
    d.Restore()
}
```

## 两种跑法（共用同一份画面代码）

- **观看**（直接跑）：声音是主时钟，画面跟着声音位置走，音画结构性不漂。
- **录制**：`shots/film.js` 在页面脚本跑之前注入 `window.__filmT = 0`，
  引擎一看到这个全局就知道有人在外部推时间 → 切 filmkit 时钟、
  不播声音、不显示按钮。判据只有一条：`__filmT` 是不是数字。

## 审片模式（?audit=1）

URL 加 `?audit=1` 打开：画面底部出现控制条——播放/暂停、当前时间
（裸秒数，和分镜表/字幕表/frames.py 的时间戳同一口径）、可点可拖的
进度条。**跳转带着声音走**，用来核对台词、音效和画面的对位；拖动
即暂停，松手后按播放继续。无声片同样可用（进度条只推画面时钟）。

控制条照常用播放器的规矩：暂停时常驻；播放中静置 2.5 秒自动隐藏
（把字幕让出来），点画面、动鼠标、碰控制条任意一处唤出，且不打断播放。

键盘（审片模式）：空格 播放/暂停，← → 逐帧（±1/FPS，`Film.FPS` 可调），
Shift+← → ±1s，J 0.5× 慢速 / L 1× / Shift+L 2× / K 暂停——
无声片变速同样生效。

## 时间轴总览（?timeline=1）

URL 加 `?timeline=1`（native：`NSGO_FILM_TIMELINE=1`）：每幕一帧的
contact sheet，点卡片跳进该幕审片（时钟与媒体都 seek 到幕首）。
审片的「地图」——只读，不编辑。

```sh
python3 tools/serve.py "dist/<名字>-0.1.0-web" 8099   # 带 Range：审片跳转需要
# 浏览器打开 http://127.0.0.1:8099/?audit=1
```

## 命令

```sh
# 观看（web）
go tool nsgo build -target web .
python3 tools/serve.py "dist/<名字>-0.1.0-web" 8099   # 带 Range：审片跳转需要

# 观看（macOS 原生，AVPlayer 出声）
go tool nsgo build -target macos_arm64 .
open "dist/<名字>-0.1.0-macos_arm64/<名字>.app"

# 原生审片（控制条 + AVPlayer 精确 seek）
NSGO_FILM_AUDIT=1 "dist/<名字>-0.1.0-macos_arm64/<名字>.app/Contents/MacOS/<名字>"

# 录制（逐帧，无声；音频后期混）
python3 tools/serve.py "dist/<名字>-0.1.0-web" 8099   # 带 Range：审片跳转需要
CHROME=/path/to/chrome node shots/film.js http://127.0.0.1:8099/ out <时长秒> 30
# 并行：WORKERS=3 开 3 个 Chrome 各渲一段（默认 min(4, CPU 核数)）
# 续录：frames/ 里已有的帧自动跳过——中断后重跑同一条命令接着录

# 审片（按时间戳抓单帧）
CHROME_BIN=/path/to/chrome python3 tools/frames.py http://127.0.0.1:8099/ frames 3,8,15
```


## 分镜装配：Series / Seq / 转场（新写法，老写法兼容）

```go
film.Run(film.Film{
    // Duration 不填 = 自动取最晚的幕尾
    Fade: 0.5,
    Bg:   paper,
    Shots: film.Series(0.5,
        film.Shot{D: 6.5, Fn: sceneA},                          // 0–6.5
        film.Shot{D: 6, Fn: sceneB, In: film.FadeBlack},        // 6–12，经底色过渡
        film.Shot{D: 3.5, Fn: sceneC, In: film.Cut, Out: film.FadeBlack}, // 11.5–15，硬切进、淡出收
    ),
})
```

- **Series**：幕只声明时长 `D`，A/B 与转场重叠自动排布；时长改了全片重排。
- **Seq + 局部时间**：`Frame.LocalT/LocalD` 恒等于 `T-A`/`B-A`——场景函数
  只用 LocalT 写，就能整体搬家、复用（`film.Seq(start, dur, fn)` 是语法糖）。
- **转场**：`In`/`Out` ∈ `Cross`（默认，= 现状交叉溶解）/ `Cut`（硬切）/
  `FadeBlack`（经底色过渡，前半程出幕淡出+底色浮现，后半程入幕淡入）/
  `FadeWhite`（白色闪场）。最后一幕 `Out=FadeBlack` = 片尾淡出。

## 字幕：自动换行 + 逐条覆写

```go
Subs: []film.Sub{
    {A: 1.8, B: 5.4, Text: "长句会按 0.62×画幅宽自动换行，底板随行数加高。"},
    {A: 7.0, B: 8.4, Text: "顶部小字", Pos: film.SubTop, Size: 40, Color: cyan},
    {A: 12.4, B: 14.6, Text: "居中卡片", Pos: film.SubCenter, Size: 56, Color: gold},
}
```

零值兼容：`Size/Color` 为零用 `SubStyle` 的，`Pos` 为零在底部。

## 数学动画（mathkit）

manim 式原语见 `mathkit/`：坐标系、函数描画（Plot/Area/Tangent）、
公式（LaTeX 预生成 `tools/eq.py` → 逐笔写 EqWrite / 变形 EqMorph）、
3D 针孔投影与曲面线框（Axes3D）。演示片：`../examples/math-demo/`。


## 音效（SFX）与成片合成

```go
film.Run(film.Film{
    ...
    Sfx: []film.Sfx{
        {T: 3.2, Asset: "asset:ding.m4a", Gain: 0.8},
    },
})
```

影片时钟走到 T 的沿上触发（拖进度条跨过 T 也能预览，倒带自动重置）。
录制模式不出声——声音统一走后期：

```sh
# 音效表和 Film.Sfx 同源：film.SfxJSON 打印，存成 sfx.json
python3 tools/mux.py 录制.mp4 final_audio.m4a 成片.mp4 --sfx sfx.json
```

## 确定性回归（tools/diff.py）

```sh
CHROME=/path/to/chrome python3 tools/diff.py http://127.0.0.1:8099/ /tmp/diff 67 30
# OK: 2010 frames identical across two runs
```

同一部片录两遍、逐帧哈希比对——任何渲染层改动之后跑一遍，
退出码 0 = 「逐像素一致」的承诺还在。

## 运动质量（filmkit）

```go
filmkit.Spring{Stiffness: 130, Damping: 8}.At(t)  // 弹簧，欠阻尼会过冲
filmkit.Bezier(0.16, 1, 0.3, 1)(p)                // CSS 同款 cubic-bezier
filmkit.KF(t, []float64{2,5,8}, []float64{0,10,0}, filmkit.EaseInOut) // 多站关键帧
```


## PPT 式视频模版（tpl）

把 html-ppt-designer STYLES.md 的 17 种视觉风格映射成视频模版：
`tpl.Style` 是一组 token（调色板/字体/强调色），`tpl.Title / Section /
Bullets / KPI / Quote / ImageFull / Outro` 是预制分镜构造器——
一部「PPT 式视频」就是选风格 + 按内容排 Series：

```go
film.Run(film.Film{
    W: 1920, H: 1080, Fade: 0.5,
    Shots: film.Series(0.5,
        film.Shot{D: 5, Fn: tpl.Title(tpl.TED, "把想法拍成视频", "副标题")},
        film.Shot{D: 6.5, Fn: tpl.Bullets(tpl.Keynote, "要点", items), In: film.FadeBlack},
        film.Shot{D: 6, Fn: tpl.KPI(tpl.Consulting, "成果", []tpl.KPIItem{
            {Value: 3.2, Unit: "×", Label: "产出速度", Decimals: 1},
        }), In: film.FadeBlack},
        film.Shot{D: 5, Fn: tpl.Outro(tpl.NeoTokyo, "MAKE IDEAS MOVE", "…"), Out: film.FadeBlack},
    ),
})
```

预设：TED / Keynote / TypicalPPT / Gamma / Consulting / Editorial /
Swiss / Newspaper / Bauhaus / Kinfolk / Muji / Brutalist / NeoTokyo /
DarkMode / RedBlackTech / Cartoon / Education。风格画廊演示：`../examples/style-demo/`（成片 `../films/style-demo.mp4`）。

## 节点分镜（keyframe）

`film.Shot` 是「一笔一笔画」，`keyframe.Take` 是「一棵节点树 + 每通道一条轨道」。三种能力是它值得单独存在的原因：**能定格**（同一棵树播放和录像是同一个画面）、**能自检**（看不见屏幕也能审）、**能读**（时间表是文本）。

```go
take := keyframe.New("orbit", 8*time.Second)        // 一支 take = 一个时钟

take.Disc("moon").At(120, 60).Size(54, 54).Fill(core.RGB(0x2f5f86)).
	On("x", keyframe.Keys(0, 120, 8, 300)).          // 秒，不是 0..1
	On("scale", keyframe.Keys(0, 0.6, 3, 1.1, 8, 0.9))
take.Text("cap", "3 SHOTS · 1 CLOCK").Font(17).At(40, 84).
	On("opacity", keyframe.Win(2*time.Second, 7*time.Second, 400*time.Millisecond))

keyframe.Tie(cx)          // 跟时钟走：__filmT / 所在影片的播放头 / 墙钟
take.Frame(cx)            // 编译；播放 0 commit/帧，定格 1 commit/帧
```

节点类型与可控通道（写错不是静默失败，`Lint` 会指出来）：

| Kind | 画什么 | 通道 |
| --- | --- | --- |
| `Group` | 不画东西，分组＋统一淡出（满幅合成箱，`At` 无效，只能靠 `x`/`y` 轨道移动） | `x` `y` `opacity` |
| `Box` `Disc` `Text` | 矩形 / 圆盘 / 文字（自带底板、内距） | `x` `y` `scale` `scaleX` `scaleY` `rotate` `opacity` |
| `Camera` | 一个 `core.Shader` 场景，相机是 uniform | `eye.x/y/z` `look.x/y/z` `fov` `tint` |
| `Canvas` | 自己的 `Draw` 回调 | — |

```go
keyframe.Keys(0, 24, 3, 200)               // 关键帧：秒 → 值
keyframe.KeysAt(8300*ms, 24, 15100*ms, 60)  // 或者直接给 time.Duration
keyframe.Win(500*ms, 5*time.Second, 400*ms) // 窗口：0→1（带 ramp）→保持→0
keyframe.Keys(...).Ease(keyframe.EaseInOut) // 一条轨道一个缓动
keyframe.EaseOut                            // 进场：快出慢停（EaseIn 是它的镜像）
```

**没有屏幕也能审。** `Lint()` 报十二类真错：节点重名、`Track` 指向不存在的节点、通道写错（`eye.z` 写到 `Box` 上）、同一通道两条轨道、关键帧超出镜头、时间重复的帧被丢掉（硬切要两条通道）、`opacity` 越界、某条轨道 >¼ 时长完全不动（死区间），以及画不出来的节点——camera 没 `Program`、canvas 没 `Draw`、disc 没宽度（撑满成方块）、group 用 `At` 定位。`Timeline()` 把每条轨道按自己的缓动画成 48 列 ASCII 曲线，`Dump()` 再加节点清单；`Sample(t)` 是纯数值求值，用来写断言而不是比截图。

**录制走同一个录制器**：推 `__filmT` 就够，不需要第二套钩子。页面家具用 `HIDE` 藏、只录一块用 `CLIP`、先点进页签用 `PRESS`：

```sh
cd film-framework
READY='[data-testid="tab-Studio"]' PRESS=Studio HIDE='[data-testid^="chip-"]' \
CLIP='[data-testid="stage-S1"]' SIZE=1280x900 WORKERS=3 \
node shots/film.js http://127.0.0.1:8099/ ../lab/shots/rec/take-s1 18 10
```

录制器一看到 `[data-testid^="take-"]` 就多等一帧、再等 `SETTLE` ms（默认 60）：seek 先落进时钟，平台帧时钟才把动画值写进视图树。演示见 `../lab/page_studio.go`（三张卡：18s 相机长镜、12s 三幕分镜、时间表与 Lint 的自画像），成片 `../lab/docs/take-s1.mp4`、`../lab/docs/story-s2.mp4`。

## 把 take 放进影片：层与窗口

`film.Shot` 画在画布上，`keyframe.Take` 活在节点树里。两者要出现在同一帧，
就把 take 挂成一层：

```go
film.Run(film.Film{
    W: 1920, H: 1080, Duration: 21,
    Shots: film.Series(0.6, film.Shot{D: 21, Fn: backdrop}),   // 底图：一画到底
    Layers: []film.Layer{
        film.Over(title, 0.4, 7.0),    // 0.4s 挂上，7.0s 摘掉
        film.Over(cards, 7.0, 14.6),
        film.Over(outro, 14.6, 21),
    },
    Subs: []film.Sub{...},
})
```

- **窗口决定「在不在」，轨道决定「看不看得见」。** `Over` 的 `start`/`end` 是影片
  时间：窗口外这层根本不挂载（不画、不推进），窗口内它叠在画布上面；淡入淡出
  仍然写在各条 `opacity` 轨道里。两件事分开，暂停、拖进度、逐帧才对得上。
- **`film` 自己调 `keyframe.Tie(cx)`**，有层就有时钟跟随，不必在 `main` 里补一行。
- **`Take.Delay(d)`**（`Over` 内部用它）：把 take 自己的 0 对齐到影片第 d 秒。
  窗口整体搬家时关键帧一个都不用改——`SeekAll` 会把 d 减掉。单独播放的
  `Play()` 不受影响，它永远从自己的 0 开始。
- **层按画幅缩放并裁切**：1920×1080 的 take 在 1280×720 的窗口里等比缩到
  0.667 居中，越界的部分被裁掉（`overflow-hidden`），所以设计尺寸只管一处。
- **字节点量自己的宽度、不换行**：多行正文一行一个 `Text`，别指望自动折行。
- **别用 `core.Dynamic` 包窗口判断**：`Dynamic` 每帧重建，take 会每帧重新
  `bind` 一次（每帧造一个 Holder）。用 `core.Show`——它的 `When` 只在真假翻转
  时才换分支，播放期间零提交。`Over` 已经这么做了。
- 例：`../examples/take-demo/`（21 秒 16:9：一层 canvas 底图 + 三层字卡 + 三条字幕），
  成片 `../films/take-demo.mp4`，动图 `../lab/docs/take-demo.gif`；
  `../examples/portrait-demo/`（10 秒竖屏 1080×1920：一条 `Over(rail, 0, 0)` 铺满全片的层
  + 三张按窗口换的字卡，画布自己交叉溶解一次），成片 `../films/portrait-demo.mp4`。

## 拍一部新片

1. 复制一个例子工程（`../examples/cylinder-demo/` 最小，或任何故事工程）当模板，改 `nsgo.yaml` 的名字和 id
2. 生图/抠图 → `gen_assets.py` 生成 `assets/` + `assets_gen.go`
3. TTS + 配乐 → 混成 `final_audio.m4a` → 套 64×64 黑场壳成 `resources/film.mp4`
   （`ffmpeg -i final_audio.m4a -f lavfi -i color=black:s=64x64:d=<时长> -shortest ...`）
4. 写分镜：`main.go` 声明 Film，`shots.go` 写每一幕
5. 构建 → 看片 → 录制 → `story/assemble.py` 式的合成
