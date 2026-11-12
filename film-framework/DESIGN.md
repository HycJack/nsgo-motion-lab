# film-framework 表达能力提升 · 设计

原则不变：**分镜是影片时钟的纯函数**。所有新原语都是声明式的、无状态的，
同一个时刻画出来的东西永远一样；不做可视化编辑器——编辑发生在 Go 代码里，
框架只提供更锋利的词。

---

## 0. 已完成（本次基线）

- **审片模式** `?audit=1`：播放/暂停、裸秒时间、带声音跳转的进度条；
  播放中静置 2.5s 自动隐藏控制条（把字幕让出来），点画面/动鼠标唤出。
- **字幕底板**按字形视觉中心对齐（基线 ± em 度量），不再浮在底板上沿。
- **tools/serve.py**：带 Range 的静态服务——`python3 -m http.server` 会让
  浏览器判定媒体不可 seek，跳转静默失效。

## 1. 字幕富化（优先级最高）

字幕是每部片子的硬需求，先补齐。

- **自动换行**：`Sub` 长句按 `MaxW`（默认 0.62×画幅宽）贪心断行，
  逐字测宽（CJK 逐字、Latin 按词），底板高度随行数走，行距 1.25×。
  `drawSubs` 现在是单行假设，改为先排版后绘制。
- **逐条覆写**：`Sub` 增加可选项，不用为一句台词动全局样式：
  ```go
  type Sub struct {
      A, B     float64
      Text     string
      Dialogue bool
      Size     float32  // 0 = 用 SubStyle.Size
      Color    core.Color // 0 = 用旁白/对白色
      Pos      SubPos   // Bottom(默认) / Top / Center
  }
  ```
- **位置档**解决旁白和对白同屏打架；`Center` 档顺带覆盖「字幕卡」
  （片头标题这种居中大字）。

## 2. 转场原语

现状只有交叉溶解（`Win` 控制幕层透明度）。设计成幕的属性：

```go
type Transition uint8
const (
    Cross Transition = iota // 现状：重叠区双方 alpha 反相
    Cut                     // 显式硬切（重叠 0 秒）
    FadeBlack               // 出幕淡出 → Bg 停留 → 入幕淡入
    FadeWhite               // 同上，白色闪场
)

type Shot struct {
    A, B  float64
    Fn    func(d *core.Drawing, f *Frame)
    In, Out Transition // 默认 Cross，零值兼容老工程
}
```

实现要点：FadeBlack/FadeWhite 需要「间隙」语义——重叠区的前半段出幕
`1→0`、中点两幕全透明露出底色、后半段入幕 `0→1`。引擎按
`In/Out` 分别算两侧的包络，替代现在单一的 `Win`。

## 3. SFX 与成片合成 ✅（已完成：Film.Sfx 沿触发 + 倒带重置 + tools/mux.py）

现状音效必须在外部混音时对时间戳，改分镜 = 重混。设计：

```go
type Sfx struct {
    T     float64  // 触发时刻（影片时钟）
    Asset string   // resources/ 下的短音频
    Gain  float32
}
Film{ ..., Sfx: []film.Sfx{...} }
```

- **观看模式**：每个 Sfx 一个隐藏 `<audio>`，在帧循环里检测
  filmT 跨越 T 的沿（上一帧 < T ≤ 这一帧）触发 `play()`——
  挂影片时钟而不是墙钟定时器，跳转/暂停语义才正确。
- **录制模式**：不发声（和主音轨一样后期处理）。
- **tools/mux.py**：录制的无声 mp4 + 主音轨 + Sfx 表 → ffmpeg
  `adelay`/`amix` 一条命令出成片。流水线最后一环收进框架，
  各故事工程不再自备 assemble 脚本。

## 4. 确定性回归 ✅（已完成：tools/diff.py，math-demo 实跑 20 帧逐字节一致）

「逐像素一致」目前没有验证手段。设计 `tools/diff.py`：

```
同一 URL 录两遍 → 逐帧哈希 → 逐帧比对 → 全等才退出 0
```

配合 `shots/film.js` 的续录能力（frames/ 里已有的帧跳过），
CI 里一条命令守住确定性。没有它，任何渲染层改动都可能
悄悄打破「拍两遍不一样」而不自知。

## 5. 镜头组合原语

`Frame.Cam` 之上加常用组合，全部纯函数：

- `Frame.Push(d, zoomFrom, zoomTo)`——整幕推近（内部仍受 zoom≥1 约束）
- `Frame.Shake(d, amp)`——确定性抖动，位移由 `filmkit.Hash` 驱动
- `Frame.Parallax(d, layers []Layer)`——多层视差，层速度系数声明式给出

## 6. A() 纪律的机械检查

溶解依赖每个绘制走 `A()`，漏一处就会在溶解瞬间"跳"出来。
`render()` 进出时断言 `curAlpha == 1`（非 1 说明某幕漏 Restore 或
手动 SetAlpha 没走 A()），panic 带上幕的序号——审查期就炸，
不要等到成片里看见。

---

## 优先级与理由

| 序 | 项 | 理由 |
|---|---|---|
| 1 | 字幕富化 | 每部片子的硬需求，长台词现在会溢出画面 |
| 2 | 转场原语 | 表达力提升最大的一块，实现小 |
| 3 | SFX + mux | 改变迭代方式：改分镜不再重混音 |
| 4 | diff.py | 兑现确定性承诺，改动渲染层后的安全网 |
| 5 | 镜头原语 | 锦上添花，现有 Cam 已能覆盖大部分 |
| 6 | A() 断言 | 半小时的活，收益是省掉一类难查的闪帧 |

---

## 7. 参考 Remotion（2026-11 调研）

> 更新：spring/Bezier/KF 三件套已完成（filmkit/motion.go，闭式解+固定迭代二分，含单测）。

Remotion 与本框架同构（帧 = 时间的纯函数），以下按价值排序：

**平移过来的**

1. **Series 自动排布**：分镜只声明时长 `Shot{D, Fn}`，`film.Series(fade, ...)`
   计算 A/B 与转场重叠——消掉手写时间戳与 Fade 手动同步这类簿记 bug。
2. **嵌套 Sequence**：场景即子时间轴。`film.Seq(start, sceneFn)`，Frame 增加
   `localT`（场景内从 0 计时）。场景可复用、可整体移动，60s 片子不再把所有
   坐标挂在全片绝对时间上。实现极小，收益最大。
3. **spring() + 贝塞尔缓动**：`filmkit.Spring(t, Spring{Stiffness, Damping})`
   （确定性解析解/模拟，不碰墙钟）、`filmkit.Bezier(t, x1, y1, x2, y2)`。
4. **多站插值**：`filmkit.KF(t, []float64{t0,t1,t2}, []float32{v0,v1,v2}, easing)`
   ——一段代码吃掉"起-升-停-落"的多段 Seg+Lerp 组合。
5. ✅ **录制并行化**：WORKERS 个 Chrome 进程各渲一段连续帧区间，按全局
   帧号落盘、ffmpeg 顺序装配；断点续录（已有帧跳过）+ 缺帧检查。
   实测 3 workers 2.7× 提速，跨 worker 逐帧比对 0 差异。
6. **音量包络**：BGM 增益按 `Seg` 随对白区间自动 ducking，替代外部混音拉线。
7. **审片键盘 + 分镜刻度**：空格/←→逐帧/J K L；进度条上画 Shots/Subs 边界刻度。
8. **只读时间轴总览 `?timeline=1`**：每幕一帧的 contact sheet，点哪跳哪
   （Studio 的"读"部分，不是编辑器）。

**明确不抄**：React 组件树（NSGo View 已声明式）、整数帧时间（秒制与
分镜表/TTS/frames.py 同口径）、Studio 式可视化编辑（编辑留在代码里）。

并入总排期：§7.2（嵌套 Sequence）与 §7.1（Series）插到字幕富化之后做，
§7.5（并行录制）独立可先行——它是纯工具层改动，不碰引擎。

## 8. mathkit：manim 式数学动画（3D 与 LaTeX 的 NSGo 现状）

查证结论（nsgo-motion-lab/lab/page_3d.go、nsgo base/textview）：

- **3D 有三条官方路径**：`core.Shader` 光线步进（场景漫游，uniform 动画零
  提交）、CPU 投影 + Canvas（Go 算顶点 + painter 排序，可命中测试）、
  CSS 3D（仅卡片）。manim 式数学 3D 走 **CPU 投影**——manim 自己也是
  CPU 软件渲染 + painter sort，思路同构。`mathkit.Axes3D` = 点投影 +
  空间折线/曲线 + 深度排序，Motion Lab 的旋转立方体已给出骨架。
- **LaTeX**：`base/textview.Math` 解析 Markdown 的 $...$/$$...$$，
  `RenderMath` 钩子由应用供渲染（UI 排版路线，web 端可接 KaTeX）。
  拍片要走路径动画（Write/Transform），仍用预生成管线：
  `tools/eq.py` LaTeX → SVG path d → `core.ParsePath`（NSGo 原生）→
  `mathkit.Write`（按路径长度逐段描）/ `mathkit.Morph`（重采样配对插值）。

mathkit 清单：Axes / Axes3D / Plot / Tangent / Write / Morph / eq.py。

## 9. mathkit 优化清单（2026-11-12 复查后）

已做：2D/3D 坐标轴穿原点 + 箭头 + 轴上刻度（`Nums` 可选数字）、
niceStep 步长；演示片加相机环绕 + 螺旋线。

- ✅ **3D 面片着色**：SurfaceShaded（四边形填充 + 法线打光 + 深度排序，
  与 SurfaceWireframe 可叠加）。
- ✅ **3D 拾取标签**：Label3D（世界点锚定 billboard，相机转标签跟着走）。
- ✅ **分段深度排序**：Poly3 段按深度从远到近，螺旋/环绕曲线遮挡正确。
- ✅ **Morph 配对升级**：贪心全局最小代价（质心+周长+面积），多余轮廓
  走 even-odd 语义（多余 to 复制源浮现 / 多余 from 收缩到质心消失）。
- **spring/Bezier/KF**（filmkit）：运动质量三件套，纯函数、确定性。
- ✅ **EqWrite 着色**：EqWriteGradient，颜色沿总弧长渐变且焊在位置上。
- **性能**：数学场景每帧全量重算目前没问题（几百个投影），若上到
  万级线段再谈缓存。

## 10. tpl：风格模版系统（2026-11-12）

按 html-ppt-designer STYLES.md v5.2 的 17 种视觉风格，落到
`tpl/` 包：Style token（调色板/字体/强调色，逐条对照 CSS 变量）+
预制分镜构造器（Title/Section/Bullets/KPI/Quote/ImageFull/Outro，
全部 LocalT 纯函数、入场统一 rise、每幕自铺风格底色、alpha 走 film.A）。

一部「PPT 式视频」= 选 Style + 按内容排 Series。演示片：style-demo/
（六幕六风格画廊）。README 有完整写法。
