// Package flowkit 把 nsgo-flow 的 .flow 流程图渲染成视频：
// 解析（flowfile，上游包）→ 分层布局 → 逐节点/逐边的确定性动画。
//
// 时序模型：节点按（rank, 出现序）依次浮现（stagger）；一条边在它的
// 起点节点出现后开始描画（沿线生长），终点节点随后浮现。
// 全部是影片时钟的纯函数。
package flowkit

import (
	"github.com/oneo-me/nsgo-flow/flowfile"
)

// Chart 是布局完成的流程图。
type Chart struct {
	Nodes []Node
	Edges []Edge
	Title string
	pal   Palette

	W, H    float64 // 内容包围盒（在调用方给的区域内）
	nReveal []float64
	// 时序：节点 i 在 i*nodeStagger 出现，历时 nodeDur；
	// 边 e 在其 from 节点出现完 + edgeLag 后描画，历时 edgeDur。
	nodeStagger, nodeDur, edgeLag, edgeDur float64
}

// Node 是布局后的一个节点。
type Node struct {
	ID, Label, Note      string
	Kind, Status         string
	X, Y, W, H           float64 // 中心 X、顶 Y、宽、高
	Rank, Order, RevealAt float64
}

// Edge 是布局后的一条连线：Pts 是正交折线（世界坐标）。
type Edge struct {
	FromID, ToID, Label string
	Dashed              bool
	Pts                 [][2]float64
	RevealAt            float64
}

// Direction 是布局方向。
type Direction string

const (
	TopToBottom Direction = "down"
	LeftToRight Direction = "right"
)

// region 是流程图可用的区域。
type region struct {
	X, Y, W, H float64
}

// nodeSize 按标签估节点尺寸。
func nodeSize(label string) (w, h float64) {
	runs := len([]rune(label))
	w = 150
	if cw := float64(runs)*17 + 44; cw > w {
		w = cw
	}
	if w > 260 {
		w = 260
	}
	if float64(runs) > 10 {
		h = 88
	} else {
		h = 64
	}
	return w, h
}

// Layout 把文档排成分层布局（Sugiyama 简化版）：
// rank = 起点到该点的最长路（入度 0 为起点；环回退到出现序），
// 同 rank 内按出现序排开，坐标在区域内均分。
func Layout(doc flowfile.Doc, dir Direction, x, y, w, h float64) Chart {
	c := Chart{Title: doc.Title, pal: DefaultPalette()}
	n := len(doc.Nodes)
	if n == 0 {
		return c
	}

	// 邻接表与前驱表。
	next := make(map[string][]string, n)
	prev := make(map[string][]string, n)
	for _, e := range doc.Edges {
		next[e.From] = append(next[e.From], e.To)
		prev[e.To] = append(prev[e.To], e.From)
	}
	byID := make(map[string]flowfile.Node, n)
	for _, nd := range doc.Nodes {
		byID[nd.ID] = nd
	}

	// rank：最长路（拓扑序；环用入度保护——图有限，简单兜底）。
	rank := make(map[string]int, n)
	for _, nd := range doc.Nodes {
		rank[nd.ID] = 0
	}
	for round := 0; round < n; round++ {
		changed := false
		for _, e := range doc.Edges {
			if rank[e.To] < rank[e.From]+1 {
				rank[e.To] = rank[e.From] + 1
				changed = true
			}
		}
		if !changed {
			break
		}
	}

	// maxRank → 行/列数。
	maxRank := 0
	for _, r := range rank {
		if r > maxRank {
			maxRank = r
		}
	}
	rows := maxRank + 1

	// 同 rank 的出现序。
	inRank := make([][]string, rows)
	for _, nd := range doc.Nodes {
		r := rank[nd.ID]
		inRank[r] = append(inRank[r], nd.ID)
	}

	// 坐标：内容驱动的紧凑布局——行距/列距按节点实际尺寸 + 间距，
	// 整体在区域内居中（视频里不出现大片空白）。
	const gap = 76.0
	c.Nodes = make([]Node, 0, n)
	pos := make(map[string][2]float64, n)
	size := make(map[string][2]float64, n)
	vertical := dir != LeftToRight

	for _, nd := range doc.Nodes {
		w2, h2 := nodeSize(nd.Label)
		size[nd.ID] = [2]float64{w2, h2}
	}

	// 每 rank 的跨度：down 方向是行高（该行最高节点+gap），right 方向是列宽。
	span := make([]float64, rows)
	for r := 0; r < rows; r++ {
		var maxS float64
		for _, id := range inRank[r] {
			sz := size[id]
			if vertical && sz[1] > maxS {
				maxS = sz[1]
			}
			if !vertical && sz[0] > maxS {
				maxS = sz[0]
			}
		}
		span[r] = maxS + gap
	}
	total := 0.0
	for r := 0; r < rows; r++ {
		total += span[r]
	}
	startOff := 0.0
	if vertical {
		startOff = y + (h-total)/2
	} else {
		startOff = x + (w-total)/2
	}

	for r := 0; r < rows; r++ {
		cnt := len(inRank[r])
		if cnt == 0 {
			continue
		}
		for k, id := range inRank[r] {
			nd := byID[id]
			w2, h2 := size[id][0], size[id][1]
			var cx, cy float64
			if vertical {
				cx = x + w*(float64(k)+0.5)/float64(cnt)
				cy = startOff
				for q := 0; q < r; q++ {
					cy += span[q]
				}
				cy += span[r] / 2
			} else {
				cy = y + h*(float64(k)+0.5)/float64(cnt)
				cx = startOff
				for q := 0; q < r; q++ {
					cx += span[q]
				}
				cx += span[r] / 2
			}
			pos[id] = [2]float64{cx, cy}
			c.Nodes = append(c.Nodes, Node{
				ID: nd.ID, Label: nd.Label, Note: nd.Note,
				Kind: nd.Kind, Status: nd.Status,
				X: cx, Y: cy - h2/2, W: w2, H: h2,
				Rank: float64(r), Order: float64(k),
			})
		}
	}

	// 时序：节点按（rank, order）排序后的出现序 stagger。
	c.nodeStagger, c.nodeDur, c.edgeLag, c.edgeDur = 0.4, 0.45, 0.12, 0.5
	idx := make(map[string]int, n)
	for i := range c.Nodes {
		idx[c.Nodes[i].ID] = i
		c.Nodes[i].RevealAt = float64(i) * c.nodeStagger
	}

	// 边：正交折线 + 时序（from 节点出现完 + lag）。
	for _, e := range doc.Edges {
		pf, okf := pos[e.From]
		pt, okt := pos[e.To]
		if !okf || !okt {
			continue
		}
		ew, eh := size[e.From][0], size[e.From][1]
		tw, th := size[e.To][0], size[e.To][1]
		var pts [][2]float64
		if vertical {
			sx, sy := pf[0], pf[1]+eh/2
			tx, ty := pt[0], pt[1]-th/2
			mid := (sy + ty) / 2
			if sx == tx {
				pts = [][2]float64{{sx, sy}, {tx, ty}}
			} else {
				pts = [][2]float64{{sx, sy}, {sx, mid}, {tx, mid}, {tx, ty}}
			}
		} else {
			sx, sy := pf[0]+ew/2, pf[1]
			tx, ty := pt[0]-tw/2, pt[1]
			mid := (sx + tx) / 2
			if sy == ty {
				pts = [][2]float64{{sx, sy}, {tx, ty}}
			} else {
				pts = [][2]float64{{sx, sy}, {mid, sy}, {mid, ty}, {tx, ty}}
			}
		}
		fromAt := 0.0
		if i, ok := idx[e.From]; ok {
			fromAt = c.Nodes[i].RevealAt + c.nodeDur
		}
		toAt := 0.0
		if i, ok := idx[e.To]; ok {
			toAt = c.Nodes[i].RevealAt
		}
		at := mathMax(fromAt, toAt-c.nodeStagger/2) + c.edgeLag
		c.Edges = append(c.Edges, Edge{
			FromID: e.From, ToID: e.To, Label: e.Label, Dashed: e.Dashed,
			Pts: pts, RevealAt: at,
		})
	}

	// 包围盒。
	c.W, c.H = w, h
	return c
}

func mathMax(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
