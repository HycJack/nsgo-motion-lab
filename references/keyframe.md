# keyframe — 把一镜到底写成能编译的东西

`film-framework/film` 的分镜是**画布函数**：画什么由参数算出来。`keyframe` 是另一种分镜——**一棵场景树 + 每通道一条关键帧轨道**，编译成 `core.Node`：

> 一支 take = 一棵节点树 + 每通道一条轨道 + 一个时钟。

这就是 three.js 的形状（`Scene` + `KeyframeTrack` + `AnimationMixer`），区别在它编译成 NSGo 节点，而 `Seek` 之后同一个画面能**定格成帧**——所以它能被录制器逐帧拍下来。

什么时候用它而不是画布分镜：画面是**组件/图层**（文字、卡片、圆盘、图表、shader 场景），动的是位置/缩放/透明度/旋转；要逐帧确定、要能自检、要能定格截图。画布分镜更适合「一笔一笔画」的教学片。

```go
take := keyframe.New("shoot", 18*time.Second)          // 一支镜头＝一个时钟

cam := take.Camera("cam", rayProgram).Label("ray-marched take")
cam.On("eye.x", keyframe.Keys(0, -2.60, 5.4, -1.55, 18, 1.70))
cam.On("fov",   keyframe.Keys(0, 0.72, 10.8, 0.78, 18, 0.60))
cam.On("tint",  keyframe.Keys(0, 0.05, 12.6, 0.90, 18, 0.12).Ease(keyframe.EaseInOut))

take.Text("slate", "ONE TAKE").Font(20).Bold().Color(colText).
	Pad(7, 15).At(20, 18).Fill(core.RGBA(0x0b1119, 0.72)).
	On("opacity", keyframe.Win(500*ms, 5*time.Second, 400*ms))

keyframe.Tie(cx)                                      // 跟时钟走：__filmT / 所在影片的播放头 / 墙钟
take.Frame(cx)                                        // 编译成 core.Node
```

`Keys` 的单位是**秒**，不是 0..1 —— 脚本读起来就是时间表：

```
take shoot · 18.0s · 3 tracks
  cam.eye.x  #-----=======+++++++*********:::::::::        0.0 5.4 18.0
  cam.fov    #-=====+++++++********:::                      0.0 10.8 18.0
  slate.opacity  :-----====*****########                   0.5 5.0
```

## 为什么是编译器，不是运行时

`Take` 只在构建时运行一次，产出一棵普通的 NSGo 节点树；此后唯一会动的值就是主时钟 `p`。

| | 播放 | 定格 / 录制 |
| --- | --- | --- |
| 唯一活动值 | `core.Timing(p, 1, {Duration: 18s, Easing: Linear})` | `p.Set(t/18)` |
| 每帧成本 | **0 次 commit**（平台按帧推进） | 1 次 commit（重建这棵子树） |
| 谁在驱动 | 平台帧时钟 | `filmkit.Now()`：录制器的 `__filmT` > 所在影片发布的时钟 > 墙钟 |

关键性质：**同一棵树既播放又定格，所以两者不可能漂移**。`Seek(t)` 之后截到的帧，就是播到 `t` 时的那一帧；这条由 `TestTakeSeekMatchesSample` 守着（逐节点对比 `Animated.TranslateX/ScaleX/Opacity` 与 `Sample`，误差 0.01）。

## 与 three.js 的对应

| three.js | keyframe | 差别 |
| --- | --- | --- |
| `new Scene()` | `keyframe.New(name, dur)` | 场景自带时长，不需要第二个对象 |
| `Object3D` + `.add()` | `Node` + `.Add(kind, name)` / `take.Box(...)` | 类型在 `Kind` 里（`Box/Disc/Text/Canvas/Camera/Group`） |
| `new KeyframeTrack('.position[x]', times, values)` | `node.On("x", keyframe.Keys(0, 24, 3, 200))` | 时间用 `time.Duration` 或秒；通道名是短名（`x`/`eye.z`/`opacity`） |
| `new AnimationClip(name, dur, tracks)` | 隐式：`Take` 就是 clip，轨道挂在节点上 | — |
| `mixer.setTime(t)` | `take.Seek(t)` / 录制器推 `__filmT` | `Seek` 会**停住**播放 |
| `mixer.update(dt)` | `play()` → `core.Timing` | 平台推进，不 re-render |
| `renderer.setAnimationLoop` | 平台帧时钟 | Go 侧零 work |
| 放进一支影片的时间轴 | 挂进 `film.Film` 后一帧就是「seek 到影片时间」，见下一节 | — |
| — | `take.Lint()` | 看不见画面时的自检 |
| — | `take.Sample(t)` / `take.At(path, t)` | 纯数值，用来写断言 |
| — | `take.Timeline()` / `take.Dump()` | 把镜头读成文本 |

three.js 没有后三项，因为写 three.js 的人**看着屏幕调**。AI 写代码时看不到屏幕，所以自检和可读输出是必需品，不是装饰。

## 一个时钟，三种驱动

`filmkit.Now()` 是唯一读数入口，按权威顺序取三个来源：

| 顺序 | 来源 | 谁设置 | 用在哪 |
| --- | --- | --- | --- |
| 1 | 录制器的 `window.__filmT` | `shots/film.js` 在页面脚本之前注入 | 逐帧录制：帧率无关、可重跑复现 |
| 2 | 影片发布的播放头 | `film.Run` 挂载时 `filmkit.Drive(func() float64 { return filmT.Get() })`，卸载时 `cx.Cleanup` 停掉 | 片中的节点分镜跟着影片时间：暂停、拖进度、变速、`?audit=1` 逐帧全都对它生效 |
| 3 | 墙钟 | `filmkit.Wall()` | 页面里没人驱动时（lab 的页签就是这样） |

三条规矩：

- **`Take` 只读 `filmkit.Now()`，不读平台时钟。** 所以「第 7.98 秒那一帧」在录、在播、在审片三种情况下是同一个值。
- **推进时间的代码只读 `filmkit.Wall()`。** 影片自己的帧循环用它推进 `t`；在那里读 `Now()` 会读到上一帧发布出去的值，无声片就永远停在原地——这是 `film` 里唯一必须区分两者的地方。
- **`filmkit.Driven()` 判「有没有人在驱动」。** `Tie` 和 `start` 都看它：有驱动就 seek 过去（不自己播），没有才 `Play`。录制器和影片发布的时钟在这里是同一种东西，所以 `Tie` 不必知道谁在驱动。

发布是「最新的一次生效」：后挂载的影片接管时钟，被替换掉的那次 `stop` 是空操作（`atomic.Pointer` + `CompareAndSwap`），所以迟到的 cleanup 不会关掉新影片的时钟。

## 让 take 活在影片里：`Film.Layers`

take 是节点树，画布分镜是 canvas，两者要同框就挂成一层：

```go
film.Run(film.Film{
    Shots:  film.Series(0.6, film.Shot{D: 21, Fn: backdrop}),
    Layers: []film.Layer{
        film.Over(title, 0.4, 7.0),   // = film.Layer{Start: 0.4, End: 7.0, Node: title.Frame}
        film.Over(cards, 7.0, 14.6),
    },
})
```

| 概念 | 谁管 | 语义 |
| --- | --- | --- |
| `Layer.Start` / `End` | `film` | 窗口外**不挂载**（不画、不推进、不进树）；`End: 0` = 到片尾 |
| `Take.Delay(d)` | `keyframe` | take 自己的 0 对齐到影片第 d 秒；`SeekAll` 减掉它 |
| `opacity` 轨道 | 写分镜的人 | 进出场淡入淡出 |

三件事分开是有原因的：

- **`Delay` 而不是 `Start`。** 单独播放（`Play()`）时 take 从自己的 0 开始，这条路不能被窗口污染；`Delay` 只在被外部驱动时（`SeekAll`、`start` 的对表）把读数平移。所以「同一支 take 先给影片用、再拿去单播」不需要改。
- **窗口只是挂载判据，不是淡出。** 用 `Layer.End` 做淡出会让最后一帧被硬切掉；淡出走轨道，窗口在淡完之后才摘。
- **`Over` 顺手 `Delay`。** 层的窗口和 take 的对表是同一个数，两处都写必然写歪。

`film` 里有层时自己调 `keyframe.Tie(cx)`——跟随循环只需要一个，挂在影片 scope 上，窗口内的 take 各自 seek 到 `Now() - Delay`。

挂载用 `core.Show{When, Then}` 而不是 `core.Dynamic`：`Dynamic` 每帧重建整层（take 每帧重新 bind，每帧一个新 Holder），`Show` 的 `When` 只在真假翻转时才换分支，播放期间零提交。层本身是一个满幅 `Absolute` 的盒子，`OnLayout` 量出 `min(W/w, H/h)` 再用 `Transform.Scale` 缩一次，所以 take 只按设计尺寸写。

窗口、`Delay`、缩放都在 `film/layer.go`；测试 `TestLayerFollowsTheFilmClock` 用一次真实的 0.2 秒等待证明它跟的是播放头（把 `Tie` 注释掉，它会停在 x=20，而轨道要求 36）。

两个例子：`examples/take-demo/`（16:9，21 秒，三张字卡按窗口换）、`examples/portrait-demo/`（竖屏 1080×1920，10 秒，一条 `Over(rail, 0, 0)` 铺满全片的进度轨 + 三张字卡）。竖屏那个顺带证明设计画幅与窗口尺寸无关：同一份分镜先按 720×1280 的窗口录过一遍，等比缩放之后一个字都没挪。

## 通道

`On` 的通道名按节点类型限定；写错不是静默失败，`Lint` 会指出来。

| Kind | 通道 |
| --- | --- |
| `Group` | `x` `y` `opacity` |
| `Box` `Disc` `Text` | `x` `y` `scale` `scaleX` `scaleY` `rotate` `opacity` |
| `Camera`（shader） | `eye.x` `eye.y` `eye.z` `look.x` `look.y` `look.z` `fov` `tint` |
| `Canvas` | 无（自己 `Draw`，或用动画值重绘） |

`Group` 是**合成箱**：它不画东西，只用来把几段镜头分组，一组一个 `opacity` 窗口就是一次镜头交接——淡出一个组，等于淡出它所有孩子。

## 关键帧

```go
keyframe.Keys(0, 24, 3, 200)              // 秒；必须是递增的
keyframe.KeysAt(8300*ms, 24, 15100*ms, 60) // 或者用 time.Duration
keyframe.Win(500*ms, 5*time.Second, 400*ms) // 窗口：0→1（ramp）→保持→1→0
keyframe.FadeIn(from, ramp) / FadeOut(to, ramp)
keyframe.Keys(...).Ease(keyframe.EaseInOut)  // 整条轨道一个缓动
keyframe.EaseOut / keyframe.EaseIn            // 进场快出慢停 / 出场慢起快走
```

- **一条轨道只能有一个缓动**。平台是每条 `Interpolate` 一个 `Easing`，所以「先 ease-in 再 ease-out」不能靠一条轨道表达；要么拆成两条轨道（不同通道），要么接受一个整段缓动。`Lint` 会把「同一个通道写了两条轨道」报成 fatal（只有最后一条会渲染）。
- **`opacity`/`scale`/`rotate`/`fov`/`tint` 默认缓动**（`EaseInOut`），`x`/`y`/位置默认线性——位置类通道通常要的是匀速位移，而淡入淡出要的是柔和的边。
- **关键帧必须严格递增**，重复时间点会被 `dedupe` 丢掉（保留第一个），`Lint` 也会提示。
- 轨道端点之外是**夹住**的（clamp），不是循环：一个 8.3–14.8 s 的窗口在 15 s 时就是 0，不是回到开头。

## 自检：`Lint` 和 `Timeline`

```go
for _, issue := range take.Lint() {
	fmt.Println(issue)          // "error: city3.y: key 1 is at 13.00s, outside the take's 12.00s"
}
fmt.Println(take.Timeline())    // 每行一条轨道 + 48 列 ASCII 曲线（按该轨道自己的缓动画）
fmt.Println(take.Dump())        // 再加上节点清单
```

`Lint` 报的东西（十二类，都是真的踩过）：

| 检查 | 为什么 |
| --- | --- |
| 节点重名 | 选择器、`Lookup`、TestID 全靠名字，重名等于丢一个 |
| `Track` 指向不存在的节点 | 典型的手误，肉眼极难发现 |
| 通道不存在 | `"eye.z"` 写到 `Box` 上不会报错，只会什么都不动 |
| 同一通道两条轨道 | 只有最后一条生效 |
| 关键帧超出镜头时长 | 那一段永远播不到 |
| 时间重复的帧被丢掉 | `dedupe` 已经吃了它，Timeline 上看不出来；硬切要两条通道 |
| `opacity` 超出 0..1 | 有的平台直接画出透明度 2 |
| 死区间 | 一条轨道有 >¼ 时长完全不动，多半是写错了端点 |
| camera 没有 `Program` | 渲染成一个空节点，没有报错 |
| canvas 没有 `Draw` | 同上 |
| disc 没有宽度 | 撑满父容器，变成一个方块 |
| group 用 `At` 定位 | group 是满幅合成箱，只有 x/y 轨道能移动它 |

`Track` 的路径按「哪一个点能切出一个存在的节点」来切，所以 `"cam.eye.z"` 是节点 `cam` 的 `eye.z` 通道——通道名自己带点是允许的，camera 的八个通道全都这样。

`Timeline` 是给「看不见画面」的人看的：一行一条轨道，48 列 ASCII 曲线（`".:-=+*#@"`，常量轨道画一排 `-`），时间是秒且保留到毫秒（8 / 7.98 分得开）。曲线按**该轨道自己的缓动**画——`EaseInOut` 的行和匀速的行长得不一样，否则时间表画的是一条没人渲染的线。

## 导出：交给同一个录制器

节点分镜不需要自己的录制器。`film-framework/shots/film.js` 本来就是「推时钟 → 截图 → 编码」，而 take 只认 `filmkit.Now()`：

```sh
cd film-framework
CHROME=/path/to/chrome SIZE=1280x900 WORKERS=3 \
PRESS=Studio CLIP='[data-testid="stage-S1"]' HIDE='[data-testid^="chip-"]' \
node shots/film.js http://localhost:8099/ out/take-s1 18 10
# → out/take-s1/film.mp4：18 秒、180 帧、10 fps
```

- `PRESS=` 先点进页签（录制的是页面里的一块，不是整站首页）。
- `CLIP=` 只录那个元素的矩形：录一支 take，不是录一个网页（配合 `data-testid="stage-…"`）。
- `HIDE=` 先把页面家具（重播按钮、页签栏）藏掉。
- 录制器检测到 `[data-testid^="take-"]` 就多等一帧、再等 `SETTLE` ms（默认 60）。

三个坑，都不是「多试几次就好」的坑：

1. **页面外的写操作不会自己生效**。`p.Set(t)` 只改值节点：没有 flush、没有请求帧，树照旧画上一帧。所以 `Tie` 的循环跑在 `cx.RequestAnimationFrame` 的回调里，seek 就发生在一次真正的帧上（`Scope.Post` 在 web 上没有宿主，会 panic）。
2. **镜头必须是长寿对象**。`keyframe.New` 放在页面函数里意味着每次 commit 重建时钟，外部 seek 永远追不上；用 `sync.Once` 建一次，页面每次只 `take.Frame(cx)`。时钟本身挂在 `core.Holder` 上，锚在 `Tie(cx)` 收到的那个 app scope：**卸下页签再回来，画面必须还在原来的时刻**——不这么做的话平台的卸载会释放那个值节点，而释放后的值读出来是 0，镜头就静默地退回第 0 帧（`TestTakeNodeIsStable` 守着这条）。
3. **seek 需要一次 commit 才能落到视图上**。`Frame` 把树包在 `core.Dynamic` 里、读一个 version 信号，`Seek` 顺手把它 +1：一次 commit 只重建这支镜头，别的卡片不动。

还有一条更隐蔽的：**有人驱动时钟时不要自己播放**。`start` 在挂载时看到 `filmkit.Driven()`（录制器推了 `__filmT`，或者它所在的那支影片发布了自己的播放头）就只 seek 到那个时刻，不启动 `core.Timing`；否则本地时钟和外部时钟每帧都在写同一个值，画面会在两者之间抖。

## 测试

```go
// 图上写的数，和算子算出来的数，必须一致
take.Seek(9 * time.Second)
app.Settle()
if got := float64(app.Query("id:take-story-moon")[0].Animated.TranslateX); math.Abs(got-take.At("moon.x", 9*time.Second)) > 0.01 { … }
```

- `Sample(t)` 返回 `map[string]float64`（`"moon.x"`、`"wide.opacity"`），不碰任何平台，纯函数。
- TestID 约定 `take-<镜头名>-<节点路径>`，所以 `app.Query("id:take-shoot-cam")` 就是那个 shader。
- 镜头名全局唯一（`Takes()` 按注册顺序），`TakeNamed` 反查。

## 移植到别的项目

这个包依赖 `core` 和 `filmkit`（时钟那三个文件）、八个文件约 1600 行（`keyframe.go` 819、`lint.go` 324、`sample.go` 87、`clock.go` 90、`util.go` 18、两个时钟分支 23、测试 241）：

| 文件 | 内容 |
| --- | --- |
| `keyframe.go` | `Take`/`Node`/`Keyframes`/编译器（`build`/`place`/`look`/`box`/`uniforms`）＋注册表 |
| `clock.go` | `Tie(cx)`：一帧读一次 `filmkit.Now()`，动了就 `SeekAll`；也记住 app scope 供 `Holder` 锚定 |
| `clock_js.go` / `clock_other.go` | `Recording()`：转问 `filmkit`；真正的判据是 `filmkit.Driven()`（录制器或影片发布的时钟） |
| `sample.go` | `Sample`/`At`：纯数值求值，缓动直接复用 `Easing.At`，所以和渲染端逐位一致 |
| `lint.go` | `Lint`/`Timeline`/`Dump` |
| `keyframe_test.go` | 夹住、窗口形状、去重、六类 Lint、时间表 |

拷过去改两处即可：`rayProgram` 换成你的 shader，调色板换成你的颜色。它不 import 本仓库的 `theme.go`；包级状态只有镜头注册表和 `Tie` 记下的 app scope。
