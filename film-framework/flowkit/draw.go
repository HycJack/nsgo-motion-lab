package flowkit

// 流程图绘制：形状按 Kind（start/end 药丸、decision 菱形、data 平行
// 四边形、task 圆角矩形），颜色按 Status（done 绿 / doing 金 / blocked
// 红 / todo 灰）。reveal 动画：节点浮现 + 边沿折线生长 + 箭头。

import (
	"math"

	"github.com/oneo-me/nsgo/core"
)

// Palette 是流程图的配色。
type Palette struct {
	NodeFill   core.Color
	NodeStroke core.Color
	Text       core.Color
	Muted      core.Color
	Accent     core.Color // 边与箭头
	Done       core.Color // status=done
	Doing      core.Color // status=doing
	Blocked    core.Color // status=blocked
	Font       string
}

// DefaultPalette 返回 DarkMode 风格的配色（和 tpl.DarkMode 同族）。
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
		Font:       "Noto Sans SC",
	}
}

// WithPalette 换一套配色（链式便捷）。
func (c Chart) WithPalette(p Palette) Chart { c.pal = p; return c }

// Duration 返回整张图动画播完的时刻。
func (c Chart) Duration() float64 {
	d := 0.0
	for i := range c.Nodes {
		if t := c.Nodes[i].RevealAt + c.nodeDur; t > d {
			d = t
		}
	}
	for i := range c.Edges {
		if t := c.Edges[i].RevealAt + c.edgeDur; t > d {
			d = t
		}
	}
	return d + 0.3
}

// Draw 画出时刻 t 的整张图。
func (c Chart) Draw(d *core.Drawing, t float64) {
	pal := c.pal
	for i := range c.Nodes {
		nd := &c.Nodes[i]
		a := filmkitF(clamp01((t - nd.RevealAt) / c.nodeDur))
		if a <= 0.001 {
			continue
		}
		drawNode(d, nd, pal, float64(a))
		if a >= 1 && nd.Note != "" {
			noteAlpha := clamp01((t - nd.RevealAt - c.nodeDur) / 0.4)
			if noteAlpha > 0 {
				coreText(d, nd.Note, float32(nd.X), float32(nd.Y+nd.H+30), 26, pal.Muted, float32(noteAlpha), 0.5)
			}
		}
	}
	for i := range c.Edges {
		e := &c.Edges[i]
		p := clamp01((t - e.RevealAt) / c.edgeDur)
		if p <= 0.001 {
			continue
		}
		drawEdge(d, e, pal, p)
	}
}

// drawNode 按形状与状态画一个节点。
func drawNode(d *core.Drawing, nd *Node, pal Palette, alpha float64) {
	a := filmkitA(alpha)
	x, y := float32(nd.X), float32(nd.Y)
	w, h := float32(nd.W), float32(nd.H)
	stroke, fill := pal.NodeStroke, pal.NodeFill
	switch nd.Status {
	case "done":
		stroke, fill = pal.Done, dimFill(pal.Done)
	case "doing":
		stroke, fill = pal.Doing, dimFill(pal.Doing)
	case "blocked":
		stroke, fill = pal.Blocked, dimFill(pal.Blocked)
	case "todo":
		stroke, fill = pal.Muted, pal.NodeFill
	}

	d.SetFillColor(fill)
	d.SetAlpha(a)
	d.SetStrokeColor(stroke)
	d.SetLineWidth(3)
	d.BeginPath()
	switch nd.Kind {
	case "start", "end":
		// 药丸：全圆角。
		d.RoundRect(x, y, w, h, h/2)
	case "decision":
		// 菱形。
		d.MoveTo(x, y+h/2)
		d.LineTo(x+w/2, y)
		d.LineTo(x+w, y+h/2)
		d.LineTo(x+w/2, y+h)
		d.ClosePath()
	case "data":
		// 平行四边形。
		sk := float32(18)
		d.MoveTo(x+sk, y)
		d.LineTo(x+w, y)
		d.LineTo(x+w-sk, y+h)
		d.LineTo(x, y+h)
		d.ClosePath()
	default:
		d.RoundRect(x, y, w, h, 12)
	}
	d.Fill()
	d.Stroke()

	// 文字：居中（过长缩小）。
	size := float32(30)
	for d.MeasureText(nd.Label, core.Font{Family: pal.Font, Size: core.Points(size)}) > w-24 && size > 20 {
		size -= 2
	}
	lw := d.MeasureText(nd.Label, core.Font{Family: pal.Font, Size: core.Points(size)})
	d.SetFillColor(pal.Text)
	d.SetAlpha(a)
	d.FillText(nd.Label, x+w/2-lw/2, y+h/2+float32(size)*0.34,
		core.Font{Family: pal.Font, Size: core.Points(size)})
	d.SetAlpha(1)

	// 状态角标：done ✓ / doing ● / blocked ✕。
	mark := ""
	mc := pal.Muted
	switch nd.Status {
	case "done":
		mark, mc = "✓", pal.Done
	case "doing":
		mark, mc = "●", pal.Doing
	case "blocked":
		mark, mc = "✕", pal.Blocked
	}
	if mark != "" {
		d.SetFillColor(mc)
		d.SetAlpha(a)
		d.FillText(mark, x+w-30, y+26, core.Font{Family: pal.Font, Size: core.Points(24)})
		d.SetAlpha(1)
	}
}

// drawEdge 画正交折线 + 箭头 + 标签，reveal 沿折线生长。
func drawEdge(d *core.Drawing, e *Edge, pal Palette, p float64) {
	if len(e.Pts) < 2 {
		return
	}
	pts := make([][2]float32, len(e.Pts))
	for i, w := range e.Pts {
		pts[i] = [2]float32{float32(w[0]), float32(w[1])}
	}
	// 总长与描画终点。
	var total float64
	for i := 1; i < len(pts); i++ {
		total += dist(pts[i-1], pts[i])
	}
	end := total * p

	// 箭头大小：终点方向由最后一段给出（即使未画到头）。
	lastSeg := pts[len(pts)-2]
	lastEnd := pts[len(pts)-1]

	d.SetStrokeColor(pal.Accent)
	d.SetAlpha(filmkitA(0.9))
	d.SetLineWidth(3)
	if e.Dashed {
		d.SetLineDash(12, 9)
		defer d.SetLineDash()
	}
	d.BeginPath()
	d.MoveTo(pts[0][0], pts[0][1])
	var acc float64
	var tip [2]float32
	tipFound := false
	for i := 1; i < len(pts); i++ {
		seg := dist(pts[i-1], pts[i])
		if acc+seg <= end {
			d.LineTo(pts[i][0], pts[i][1])
			acc += seg
			continue
		}
		t := (end - acc) / seg
		tip = [2]float32{
			pts[i-1][0] + (pts[i][0]-pts[i-1][0])*float32(t),
			pts[i-1][1] + (pts[i][1]-pts[i-1][1])*float32(t),
		}
		d.LineTo(tip[0], tip[1])
		lastSeg, lastEnd = pts[i-1], pts[i]
		// 箭头方向用到达点之前的方向。
		lastSeg, lastEnd = segAt(pts, i, t)
		tipFound = true
		break
	}
	d.Stroke()
	d.SetLineDash()

	// 箭头（描完才出现）。
	if p >= 1 || tipFound {
		dirx := lastEnd[0] - lastSeg[0]
		diry := lastEnd[1] - lastSeg[1]
		l := float32(math.Hypot(float64(dirx), float64(diry)))
		if l > 0.001 {
			ux, uy := dirx/l, diry/l
			hx, hy := lastEnd[0], lastEnd[1]
			d.SetFillColor(pal.Accent)
			d.SetAlpha(filmkitA(0.9))
			d.BeginPath()
			d.MoveTo(hx+ux*14, hy+uy*14)
			d.LineTo(hx-ux*4+(-uy)*7, hy-uy*4+(ux)*7)
			d.LineTo(hx-ux*4-(-uy)*7, hy-uy*4-(ux)*7)
			d.ClosePath()
			d.Fill()
		}
	}
	d.SetAlpha(1)

	// 标签：折线中点上方。
	if e.Label != "" && p >= 1 {
		mid, midD := midOf(pts)
		_ = midD
		d.SetFillColor(pal.Text)
		d.SetAlpha(filmkitA(0.85))
		d.FillText(e.Label, mid[0], mid[1]-12,
			core.Font{Family: pal.Font, Size: core.Points(26)})
		d.SetAlpha(1)
	}
	_ = lastSeg
	_ = lastEnd
}

func segAt(pts [][2]float32, i int, t float64) ([2]float32, [2]float32) {
	if i >= 2 {
		return pts[i-2], pts[i-1]
	}
	return pts[0], pts[1]
}

func midOf(pts [][2]float32) ([2]float32, float64) {
	var total float64
	for i := 1; i < len(pts); i++ {
		total += dist(pts[i-1], pts[i])
	}
	half := total / 2
	var acc float64
	for i := 1; i < len(pts); i++ {
		seg := dist(pts[i-1], pts[i])
		if acc+seg >= half {
			t := (half - acc) / seg
			return [2]float32{
				pts[i-1][0] + (pts[i][0]-pts[i-1][0])*float32(t),
				pts[i-1][1] + (pts[i][1]-pts[i-1][1])*float32(t),
			}, half
		}
		acc += seg
	}
	return pts[len(pts)-1], total
}

func dist(a, b [2]float32) float64 {
	dx := float64(b[0] - a[0])
	dy := float64(b[1] - a[1])
	return math.Sqrt(dx*dx + dy*dy)
}

func dimFill(c core.Color) core.Color {
	r, g, b, a := c.Components()
	ch := func(v float64) uint32 { return uint32(v*0.35*255 + 0.5) }
	return core.Color(ch(r)<<24 | ch(g)<<16 | ch(b)<<8 | uint32(a*255+0.5))
}

func clamp01(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}

func filmkitF(x float64) float32 { return float32(x) }
func filmkitA(a float64) float32 { return float32(a) }

func coreText(d *core.Drawing, s string, x, y, size float32, c core.Color, alpha float32, align float32) {
	if alpha <= 0.001 || s == "" {
		return
	}
	w := d.MeasureText(s, core.Font{Size: core.Points(size)})
	d.SetFillColor(c)
	d.SetAlpha(alpha)
	d.FillText(s, x-align*w, y, core.Font{Size: core.Points(size)})
}
