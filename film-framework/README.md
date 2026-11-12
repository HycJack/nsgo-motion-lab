# nsgo-film · 用 NSGo 拍片子的通用框架

把一部影片拆成三层：**引擎**（这个包）、**素材**（贴图/音频）、**分镜**
（一部片子的全部内容声明）。分镜函数是影片时钟的纯函数——同一个时刻画
出来的东西永远一样，所以录制的片子、观看的画面、审片的单帧，三者逐像素一致。

## 结构

```
filmkit/     影片时钟 + 时间数学（Seg / EaseOut / Win / Pulse …）
film/        引擎：Film 配置 → 时钟循环、交叉溶解、音频载体、字幕、开始按钮
shots/film.js  录制器：推时钟 → CDP 逐帧截图 → ffmpeg 编码（WORKERS 并行、断点续录）
tools/         serve.py（Range 静态服务）/ frames.py（抓帧）/ eq.py（LaTeX→路径）/ mux.py（混音合成）/ diff.py（确定性回归）
../examples/   三个例子工程：math-demo（数学动画）/ cylinder-demo（圆柱展开教学片）/ style-demo（风格画廊）
../films/      三个例子录好的成片 mp4
```

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

## 拍一部新片

1. 复制一个例子工程（`../examples/cylinder-demo/` 最小，或任何故事工程）当模板，改 `nsgo.yaml` 的名字和 id
2. 生图/抠图 → `gen_assets.py` 生成 `assets/` + `assets_gen.go`
3. TTS + 配乐 → 混成 `final_audio.m4a` → 套 64×64 黑场壳成 `resources/film.mp4`
   （`ffmpeg -i final_audio.m4a -f lavfi -i color=black:s=64x64:d=<时长> -shortest ...`）
4. 写分镜：`main.go` 声明 Film，`shots.go` 写每一幕
5. 构建 → 看片 → 录制 → `story/assemble.py` 式的合成
