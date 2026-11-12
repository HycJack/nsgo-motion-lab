package flowkit

// Chart：.flow 文档 → 引擎布局（base/flow 的 Sugiyama）→ 逐节点/逐边
// 的确定性动画。尺寸用 Drawing 的真实文本测量（每帧重排，小图无感）。

import (
	"math"

	"github.com/oneo-me/nsgo-flow/flowfile"
	"github.com/oneo-me/nsgo/core"
)

// Chart 是布局完成的流程图。
type Chart struct {
	Title string
	Nodes []Node
	Edges []EdgeRoute

	Total    float64 // 整图动画播完的时刻
	pal      Palette
}

// Node 是布局后的节点卡片（X/Y 为左上角）。
type Node struct {
	ID, Label, Note string
	Kind, Status    string
	X, Y, W, H      float64
	RevealAt        float64
}

// EdgeRoute 是一条布好线的边。
type EdgeRoute struct {
	FromID, ToID, Label string
	Dashed              bool
	route               *route
	RevealAt            float64
}

// 引擎默认间距。
const (
	nodeGap = 40.0
	rankGap = 72.0
	edgeGap = 24.0
)

// Palette 是流程图配色。
type Palette struct {
	NodeFill   core.Color
	NodeStroke core.Color
	Text       core.Color
	Muted      core.Color
	Accent     core.Color
	Done       core.Color
	Doing      core.Color
	Blocked    core.Color
	LabelBg    core.Color
	LabelLine  core.Color
	Font       string
}

// DefaultPalette 返回 DarkMode 风格配色（与 tpl.DarkMode 同族）。
func DefaultPalette() Palette {
	return Palette{
		NodeFill:   core.RGB(0x1E293B),
		NodeStroke: core.RGB(0x3B82F6),
		Text:       core.RGB(0xF1F5F9),
		Muted:      core.RGB(0x64748B),
		Accent:     core.RGB(0x3B82F6),
		Done:       core.RGB(0x34D399),
		Doing:      core.RGB(0xF59E0B),
		Blocked:    core.RGB(0xF87171),
		LabelBg:    core.RGB(0x0F172A),
		LabelLine:  core.RGB(0x334155),
		Font:       "Noto Sans SC",
	}
}

// WithPalette 换配色。
func (c Chart) WithPalette(p Palette) Chart { c.pal = p; return c }

// Duration 返回整图动画播完的时刻。
func (c Chart) Duration() float64 { return c.Total }

// Layout 把文档排进区域。dir 选方向；d 用来实测文字（卡片尺寸）。
func Layout(doc flowfile.Doc, dir Direction, x, y, w, h float64, d textMeasurer) Chart {
	c := Chart{Title: doc.Title, pal: DefaultPalette()}

	nodes := make([]layoutNode, len(doc.Nodes))
	byID := make(map[string]flowfile.Node, len(doc.Nodes))
	ids := make([]string, len(doc.Nodes))
	index := make(map[string]int, len(doc.Nodes))
	for i, nd := range doc.Nodes {
		byID[nd.ID] = nd
		ids[i] = nd.ID
		index[nd.ID] = i
		lw := d.MeasureText(nd.Label, font(30, 500))
		nw := 0.0
		if nd.Note != "" {
			nw = float64(d.MeasureText(nd.Note, font(24, 400)))
		}
		w := maxF(210, math.Max(float64(lw)+112, nw+112))
		h := 62.0
		if nd.Note != "" {
			h = 96
		}
		nodes[i] = layoutNode{width: w, height: h}
	}
	edges := make([]layoutEdge, 0, len(doc.Edges))
	for _, e := range doc.Edges {
		i, okI := index[e.From]
		j, okJ := index[e.To]
		if !okI || !okJ {
			continue
		}
		edges = append(edges, layoutEdge{source: i, target: j, dashed: e.Dashed})
	}

	dir2 := TopToBottom
	if dir == LeftToRight {
		dir2 = LeftToRight
	}
	res := layout(nodes, edges, dir2, spacing{node: nodeGap, rank: rankGap, edge: edgeGap})
	// 引擎从 0,0 开始排：整体平移到区域内居中（节点与路由一起）。
	offX := x + (w-res.width)/2
	offY := y + (h-res.height)/2
	for i := range res.x {
		res.x[i] += offX
		res.y[i] += offY
	}
	for k := range res.routes {
		for j := range res.routes[k] {
			res.routes[k][j].x += offX
			res.routes[k][j].y += offY
		}
	}

	// 时序：节点按文档顺序 stagger；边在两端就位后描画。
	c.Nodes = make([]Node, len(doc.Nodes))
	reveal := make([]float64, len(doc.Nodes))
	for i, nd := range doc.Nodes {
		reveal[i] = float64(i) * 0.38
		if reveal[i]+0.45 > c.Total {
			c.Total = reveal[i] + 0.45
		}
		c.Nodes[i] = Node{
			ID: nd.ID, Label: nd.Label, Note: nd.Note,
			Kind: nd.Kind, Status: nd.Status,
			X: res.x[i], Y: res.y[i], W: nodes[i].width, H: nodes[i].height,
			RevealAt: reveal[i],
		}
	}
	for k, e := range doc.Edges {
		i, okI := index[e.From]
		j, okJ := index[e.To]
		if !okI || !okJ || len(res.routes[k]) < 2 {
			continue
		}
		at := maxF(reveal[i]+0.45, reveal[j]) + 0.1
		if at+0.6 > c.Total {
			c.Total = at + 0.6
		}
		c.Edges = append(c.Edges, EdgeRoute{
			FromID: e.From, ToID: e.To, Label: e.Label, Dashed: e.Dashed,
			route:    newRoute(res.routes[k]),
			RevealAt: at,
		})
	}
	return c
}

func font(size float32, weight core.FontWeight) core.Font {
	return core.Font{Family: "Noto Sans SC", Size: core.Points(size), Weight: weight}
}

func maxF(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
