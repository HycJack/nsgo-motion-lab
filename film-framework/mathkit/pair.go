package mathkit

import (
	"math"
	"sort"
)

// polyArea 是折线的鞋带公式绝对面积（Closed 按闭合算）。
func polyArea(p Poly) float64 {
	n := len(p.Pts) / 2
	if n < 3 {
		return 0
	}
	var a float64
	m := n
	if !p.Closed {
		m = n - 1
	}
	for i := 0; i < m; i++ {
		j := (i + 1) % n
		a += float64(p.X(i))*float64(p.Y(j)) - float64(p.X(j))*float64(p.Y(i))
	}
	return math.Abs(a / 2)
}

// contourCost 是两条轮廓的对应代价：质心距离为主，周长差与面积差为辅。
// 同一坐标空间里，形状相近的轮廓代价自然小。
func contourCost(a, b Poly) float64 {
	ax, ay := a.Centroid()
	bx, by := b.Centroid()
	d := math.Hypot(ax-bx, ay-by)
	return d + 0.3*math.Abs(a.Length()-b.Length()) + 0.3*math.Abs(polyArea(a)-polyArea(b))
}

// Pair 是一对轮廓的 morph 任务。
type Pair struct {
	From, To Poly
	// Fade 标记 To 是退化点：多余的 from 收缩消失用。
	// even-odd 填充下零面积轮廓不影响结果。
	Fade bool
}

// pairContours 把 from/to 的轮廓做最优对应（贪心全局最小代价：
// 质心距离 + 周长差 + 面积差），返回 morph 绘制对。
//
// 数量不等时的语义（manim Transform 的分裂/合并）：
//   - 多余的 to：复制代价最近的 from 作源——p=0 时重复源在 even-odd
//     填充下相互抵消（不可见），morph 中「从无到有浮现」；
//   - 多余的 from：收缩到最近 to 的质心（退化点，Fade=true）——
//     「收缩消失」，且不占用 to（目标各画一次，结果不缺笔画）。
func pairContours(from, to []Poly) []Pair {
	nf, nt := len(from), len(to)
	type cand struct{ i, j int }
	cands := make([]cand, 0, nf*nt)
	cost := func(i, j int) float64 { return contourCost(from[i], to[j]) }
	for i := 0; i < nf; i++ {
		for j := 0; j < nt; j++ {
			cands = append(cands, cand{i, j})
		}
	}
	sort.Slice(cands, func(a, b int) bool {
		return cost(cands[a].i, cands[a].j) < cost(cands[b].i, cands[b].j)
	})

	assignI := make([]int, nf) // from[i] 配到的 to 下标
	assignJ := make([]int, nt) // to[j] 配到的 from 下标
	for i := range assignI {
		assignI[i] = -1
	}
	for j := range assignJ {
		assignJ[j] = -1
	}
	usedI := make([]bool, nf)
	usedJ := make([]bool, nt)
	for _, c := range cands {
		if usedI[c.i] || usedJ[c.j] {
			continue
		}
		usedI[c.i], usedJ[c.j] = true, true
		assignI[c.i], assignJ[c.j] = c.j, c.i
	}

	// 剩余的 to：复制代价最近的 from 作源（重复源 p=0 抵消）。
	for j := 0; j < nt; j++ {
		if assignJ[j] >= 0 {
			continue
		}
		best, bd := -1, math.Inf(1)
		for i := 0; i < nf; i++ {
			if c := cost(i, j); c < bd {
				bd, best = c, i
			}
		}
		assignJ[j] = best
	}

	var out []Pair
	// 一一对应的部分。
	for i := 0; i < nf; i++ {
		if j := assignI[i]; j >= 0 {
			out = append(out, Pair{From: from[i], To: to[j]})
		}
	}
	// 多余的 to：复制代价最近的 from 作源（重复源在 p=0 抵消浮现）；
	// 若那个 from 本身也是多余的，直接消耗它，成为一对真实 morph。
	for j := 0; j < nt; j++ {
		if usedJ[j] {
			continue
		}
		best, bd := -1, math.Inf(1)
		for i := 0; i < nf; i++ {
			if c := cost(i, j); c < bd {
				bd, best = c, i
			}
		}
		if best < 0 {
			continue
		}
		out = append(out, Pair{From: from[best], To: to[j]})
		assignI[best] = j
	}
	// 多余的 from：收缩到最近 to 的质心（退化点）。
	for i := 0; i < nf; i++ {
		if assignI[i] >= 0 {
			continue
		}
		best, bd := -1, math.Inf(1)
		for j := 0; j < nt; j++ {
			if c := cost(i, j); c < bd {
				bd, best = c, j
			}
		}
		if best < 0 {
			continue
		}
		bx, by := to[best].Centroid()
		out = append(out, Pair{
			From: from[i],
			To:   Poly{Pts: []float32{float32(bx), float32(by), float32(bx), float32(by), float32(bx), float32(by)}, Closed: true},
			Fade: true,
		})
	}
	return out
}
