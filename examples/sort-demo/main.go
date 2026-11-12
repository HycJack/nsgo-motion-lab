// 《冒泡排序》—— 算法可视化教学片：离散事件驱动的纯函数动画。
//
// 与 cylinder-demo 的「连续运动学」不同，这一支的核心是**离散事件表**：
// init() 里把冒泡排序完整跑一遍，生成 28 个比较步骤（含是否交换、
// 每步之后的数组状态、完成的轮数）——画面是「时刻 t 之前发生过的
// 事件」的纯函数，帧率无关、逐帧确定。
//
// 四幕（Series，DarkMode 风格）：
//  1. 标题卡（tpl.Title）
//  2. 规则要点（tpl.Bullets）
//  3. 排序动画：8 根柱子逐轮冒泡，比较高亮、交换平移、已就位染绿
//  4. 代价：比较/交换次数 KPI + O(n²)，片尾淡出
package main

import (
	"fmt"
	"math"

	"nsgofilm/film"
	"nsgofilm/filmkit"
	"nsgofilm/tpl"

	"github.com/oneo-me/nsgo/core"
)

const (
	dw   float32 = 1920.0
	dh   float32 = 1080.0
	fade         = 0.5
	n            = 8 // 元素个数
)

var initArr = []int{5, 2, 8, 1, 9, 3, 7, 4}

// ---------------------------------------------------------------- 事件表

// step 是一次两两比较：j 与 j+1 比较，swap 说明是否交换，
// arr 是这一步之后的数组状态，pass 是这一步属于第几轮（0 起）。
type step struct {
	j    int
	swap bool
	arr  []int
	pass int
}

var (
	steps   []step
	swaps   int // 实际交换次数
)

// init 把冒泡排序完整跑一遍，生成事件表——确定性，无随机。
func init() {
	a := append([]int(nil), initArr...)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-1-i; j++ {
			sw := a[j] > a[j+1]
			if sw {
				a[j], a[j+1] = a[j+1], a[j]
				swaps++
			}
			steps = append(steps, step{j: j, swap: sw, arr: append([]int(nil), a...), pass: i})
		}
	}
}

// ---------------------------------------------------------------- 影片

func main() {
	film.Run(film.Film{
		W:    dw,
		H:    dh,
		Fade: fade,
		Bg:   tpl.DarkMode.Bg,
		Shots: film.Series(fade,
			film.Shot{D: 5, Fn: tpl.Title(tpl.DarkMode,
				"冒泡排序", "最直观的交换排序 · O(n²)")},
			film.Shot{D: 7, Fn: tpl.Bullets(tpl.DarkMode,
				"规则只有一句话",
				[]string{
					"每一轮从头到尾，两两比较相邻的两个元素",
					"大的往后换——一轮下来，最大的元素「冒」到末尾",
					"重复，直到某一轮没有任何交换",
				})},
			film.Shot{D: 15.6, Fn: sceneSort},
			film.Shot{D: 6.5, Fn: sceneCost},
			film.Shot{D: 4.5, Fn: tpl.Outro(tpl.DarkMode,
				"简单，但慢", "下一集：选择排序——少换，多看"),
				Out: film.FadeBlack,
			},
		),
		Subs: []film.Sub{
			{A: 1.0, B: 4.4, Text: "冒泡排序：最直观的交换排序。"},
			{A: 5.2, B: 10.6, Text: "每一轮从头两两比较，大的往后换——一轮结束，最大的元素就「冒」到了末尾。"},
			{A: 11.4, B: 13.6, Text: "重复几轮，直到没有任何交换。"},
			{A: 15.0, B: 17.4, Text: "比较 28 次，交换 13 次——这就是 O(n²) 的代价。"},
			{A: 17.8, B: 19.4, Text: "元素每多一个，代价平方级增长。", Dialogue: true},
		},
		SubStyle: film.DefaultSubStyle(tpl.DarkMode.Heading, tpl.DarkMode.Accent, core.RGB(0x000000)),
	})
}

// ---------------------------------------------------------------- 排序动画

// 排序动画的时间参数（局部时间）。
const (
	t0     = 0.9  // 第一步开始
	per    = 0.45 // 每步时长
	inAnim = 0.3  // 交换平移动画的时长（占每步前段）
)

var (
	cBar    = core.RGB(0x3B82F6) // 常规柱
	cCmp    = core.RGB(0xF59E0B) // 正在比较
	cSorted = core.RGB(0x10B981) // 已就位
)

// sceneSort：柱子入场 → 28 步冒泡（比较高亮/交换平移/已就位染绿）→ 收尾。
func sceneSort(d *core.Drawing, f *film.Frame) {
	t := f.LocalT

	// 当前步：t0 + k*per ≤ t。
	k := int(math.Floor((t - t0) / per))
	fin := k >= len(steps)
	if k < 0 {
		k = 0
	}
	if fin {
		k = len(steps)
	}

	// 数组状态：第 k 步之前的世界。
	arr := initArr
	if k > 0 {
		arr = steps[k-1].arr
	}

	// 已就位的尾巴：前 k 步里跑完的轮数。
	done := 0
	if k > 0 {
		done = steps[k-1].pass + 1
	}
	if fin {
		done = n
	}

	// 画布布局。
	base := dh * 0.78
	unit := dh * 0.55 / 9.0
	bw := float32(96)
	gap := (dw - float32(n)*bw) / float32(n+1)

	// 正在进行的那一步（若未结束）。
	curJ, swP := -1, 0.0
	if !fin && k < len(steps) {
		curJ = steps[k].j
		if steps[k].swap {
			swP = filmkit.Clamp01((t - (t0 + float64(k)*per)) / inAnim)
		}
	}

	for p := 0; p < n; p++ {
		v := arr[p]
		// 交换平移：比较对里的两根柱子在动画期间互换 x。
		xp := float32(p)
		if swP > 0 && curJ >= 0 {
			if p == curJ {
				xp = float32(curJ) + filmkit.F32(filmkit.EaseInOut(swP))
			} else if p == curJ+1 {
				xp = float32(curJ+1) - filmkit.F32(filmkit.EaseInOut(swP))
			}
		}
		x := gap + float32(xp)*(bw+gap)

		// 颜色：已就位尾巴绿色；比较对金色（优先于已就位色——
		// 末轮的比较对也要闪金）。
		c := cBar
		if p >= n-done {
			c = cSorted
		}
		if !fin && (p == curJ || p == curJ+1) {
			c = cCmp
		}
		h := float32(v) * unit
		d.SetFillColor(c)
		d.SetAlpha(film.A(1))
		d.BeginPath()
		d.Rect(x, base-h, bw, h)
		d.Fill()
		// 柱顶标数值。
		core2Text(d, fmt.Sprint(v), x+bw/2, base-h-14, 30, c, film.A(0.9))
	}

	// 状态行：轮次 / 比较 / 交换 计数（右下）。
	if !fin && k < len(steps) {
		st := steps[k]
		if swP > 0 && st.swap {
			core2Text(d, fmt.Sprintf("交换！ %d ↔ %d", initArrAt(arr, st.j, swP), initArrAt(arr, st.j+1, swP)), dw-320, 130, 36, cCmp, film.A(1), 1)
		} else if st.swap {
			core2Text(d, fmt.Sprintf("%d > %d ？是，交换", arr[st.j], arr[st.j+1]), dw-320, 130, 36, cCmp, film.A(1), 1)
		} else {
			core2Text(d, fmt.Sprintf("%d ≤ %d ？不换", arr[st.j], arr[st.j+1]), dw-320, 130, 36, tpl.DarkMode.Body, film.A(0.8), 1)
		}
		core2Text(d, fmt.Sprintf("第 %d 轮 · 比较 %d / %d · 交换 %d", st.pass+1, min(k+1, len(steps)), len(steps), swapsSoFar(k, swP > 0 && st.swap)), dw-320, 180, 30, tpl.DarkMode.Muted, film.A(0.8), 1)
	} else if fin {
		core2Text(d, fmt.Sprintf("完成 · 比较 %d 次 · 交换 %d 次", len(steps), swaps), dw-320, 130, 34, cSorted, film.A(1), 1)
	}

	// 收尾扫光：全部就位后，绿色从左到右扫一遍。
	if fin {
		sweep := filmkit.Seg(t, t0+float64(len(steps))*per+0.2, t0+float64(len(steps))*per+1.2)
		if sweep > 0 && sweep < 1 {
			x := gap + float32(sweep)*(dw-2*gap)
			d.SetFillColor(core.RGB(0x10B981))
			d.SetAlpha(film.A(0.18))
			d.BeginPath()
			d.Rect(x-30, base-9*unit-40, 60, 9*unit+40)
			d.Fill()
		}
	}
}

// initArrAt 是交换动画期间的语义辅助（演示用：直接读当前数组值）。
func initArrAt(arr []int, i int, _ float64) int { return arr[i] }

// swapsSoFar 统计前 k 步里的交换总数。
func swapsSoFar(k int, plusCur bool) int {
	c := 0
	for i := 0; i < k && i < len(steps); i++ {
		if steps[i].swap {
			c++
		}
	}
	if plusCur {
		c++
	}
	return c
}

// core2Text 是左对齐/右对齐文字（align 0 左 1 右）。
func core2Text(d *core.Drawing, s string, x, y, size float32, c core.Color, alpha float32, align ...float32) {
	if alpha <= 0.001 || s == "" {
		return
	}
	al := float32(0)
	if len(align) > 0 {
		al = align[0]
	}
	w := d.MeasureText(s, core.Font{Size: core.Points(size)})
	d.SetFillColor(c)
	d.SetAlpha(alpha)
	d.FillText(s, x-al*w, y, core.Font{Size: core.Points(size)})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ---------------------------------------------------------------- 代价卡

// sceneCost：比较/交换/元素 KPI + O(n²) 落版。
func sceneCost(d *core.Drawing, f *film.Frame) {
	tpl.KPI(tpl.DarkMode, "n = 8 的代价",
		[]tpl.KPIItem{
			{Value: float64(len(steps)), Label: fmt.Sprintf("两两比较 n(n−1)/2 = %d", len(steps))},
			{Value: float64(swaps), Label: "实际交换（这组数据）"},
			{Value: float64(n), Label: "元素个数"},
		})(d, f)

	if p := filmkit.Seg(f.LocalT, 2.6, 3.4); p > 0 {
		core2Text(d, "时间复杂度 O(n²) —— 元素每多一个，代价平方级增长",
			dw/2, dh*0.78, 52, tpl.DarkMode.Accent, film.A(float32(p)), 0.5)
	}
}
