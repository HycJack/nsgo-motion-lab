package flowkit

// 路由绘制：移植自引擎 base/flow/edges.go，去掉交互态的参数绑定
// （AnimatedNode/shifted），视频里节点不动，路由画一次即可。
// 保留：圆角拐弯（corners + quadTo）、虚线（按轴锚定相位）、
// 长段中点的药丸标签、箭头。

import (
	"math"

	"github.com/oneo-me/nsgo/core"
)

const (
	defaultEdgeWidth = 2.0
	cornerRadius     = 12
	cornerSteps      = 6
	arrowLength      = 12
	arrowHalfWidth   = 7
	dashLength       = 6
	dashGap          = 8
	labelSize        = 26
	labelPadX        = 12
	labelPadY        = 4
)

// route 是一条已布局的边：折线点 + 各点沿线权重（0..1）。
type route struct {
	points  []point
	weights []float64
}

// newRoute 由折线构造路由并算出各点权重。
func newRoute(points []point) *route {
	r := &route{points: points, weights: make([]float64, len(points))}
	total := 0.0
	for k := 1; k < len(points); k++ {
		total += math.Hypot(points[k].x-points[k-1].x, points[k].y-points[k-1].y)
	}
	along := 0.0
	for k := range points {
		if k > 0 {
			along += math.Hypot(points[k].x-points[k-1].x, points[k].y-points[k-1].y)
		}
		if total > 0 {
			r.weights[k] = along / total
		}
	}
	return r
}

// cut 把路由按权重 w 截断（reveal 动画用）：返回截断后的路由。
func (r *route) cut(w float64) *route {
	if w >= 1 {
		return r
	}
	pts := []point{r.points[0]}
	ws := []float64{0}
	for k := 1; k < len(r.points); k++ {
		if r.weights[k] <= w {
			pts = append(pts, r.points[k])
			ws = append(ws, r.weights[k])
			continue
		}
		seg := r.weights[k] - r.weights[k-1]
		t := 0.0
		if seg > 0 {
			t = (w - r.weights[k-1]) / seg
		}
		a, b := r.points[k-1], r.points[k]
		pts = append(pts, point{a.x + (b.x-a.x)*t, a.y + (b.y-a.y)*t})
		ws = append(ws, w)
		break
	}
	return &route{points: pts, weights: ws}
}

// draw 描出路由：圆角拐弯 + 箭头（reveal 时箭头跟着笔尖走）。
func (r *route) draw(d *core.Drawing, width float32, c core.Color, dashed bool, reveal float64) {
	rr := r.cut(reveal)
	points := rr.points
	if len(points) < 2 {
		return
	}
	n := len(points)
	end, before := points[n-1], points[n-2]
	dx, dy := end.x-before.x, end.y-before.y
	length := math.Hypot(dx, dy)
	if length == 0 {
		return
	}
	ux, uy := dx/length, dy/length
	head := math.Min(arrowLength, length)
	line := append([]point(nil), points...)
	line[n-1] = point{end.x - ux*head*.8, end.y - uy*head*.8}

	d.SetLineCap(core.CapRound).SetLineJoin(core.JoinRound)
	d.SetStrokeColor(c).SetLineWidth(width)
	if dashed {
		r.dashes(d, line, rr.weights)
	} else {
		d.BeginPath()
		d.MoveTo(float32(line[0].x), float32(line[0].y))
		corners(line, rr.weights, func(a, control, b point, wa, wc, wb float64) {
			d.LineTo(float32(a.x), float32(a.y))
			d.QuadTo(float32(control.x), float32(control.y), float32(b.x), float32(b.y))
		})
		d.LineTo(float32(line[n-1].x), float32(line[n-1].y))
		d.Stroke()
	}
	d.SetLineDash()

	// 箭头（reveal 不足一圈时画在笔尖，随生长前进）。
	d.SetFillColor(c).SetStrokeColor(c).SetLineWidth(1)
	d.BeginPath().MoveTo(float32(end.x), float32(end.y)).
		LineTo(float32(end.x-ux*head-uy*arrowHalfWidth), float32(end.y-uy*head+ux*arrowHalfWidth)).
		LineTo(float32(end.x-ux*head+uy*arrowHalfWidth), float32(end.y-uy*head-ux*arrowHalfWidth)).
		ClosePath().Fill().Stroke()
}

// dashes 沿轴锚定相位的虚线（一段一段描，拐弯延续花型）。
func (r *route) dashes(d *core.Drawing, line []point, weights []float64) {
	const period = dashLength + dashGap
	d.SetLineDash(dashLength, dashGap)
	piece := func(a point, b point) {
		length := math.Hypot(b.x-a.x, b.y-a.y)
		if length == 0 {
			return
		}
		ux, uy := (b.x-a.x)/length, (b.y-a.y)/length
		if math.Abs(ux) > .999 || math.Abs(uy) > .999 {
			anchor := ux*a.x + uy*a.y
			if ux < 0 || uy < 0 {
				anchor += dashLength
			}
			d.SetLineDashOffset(float32(math.Mod(math.Mod(anchor, period)+period, period)))
		}
		d.BeginPath()
		d.MoveTo(float32(a.x), float32(a.y))
		d.LineTo(float32(b.x), float32(b.y))
		d.Stroke()
	}
	from := line[0]
	corners(line, weights, func(a, control, b point, wa, wc, wb float64) {
		piece(from, a)
		length, prev := 0.0, a
		for step := 1; step <= cornerSteps; step++ {
			t := float64(step) / cornerSteps
			u := 1 - t
			q := point{u*u*a.x + 2*u*t*control.x + t*t*b.x, u*u*a.y + 2*u*t*control.y + t*t*b.y}
			length += math.Hypot(q.x-prev.x, q.y-prev.y)
			prev = q
		}
		d.SetLineDashOffset(float32(math.Mod(length, period)))
		d.BeginPath()
		d.MoveTo(float32(a.x), float32(a.y))
		d.QuadTo(float32(control.x), float32(control.y), float32(b.x), float32(b.y))
		d.Stroke()
		from = b
	})
	piece(from, line[len(line)-1])
}

// label 在路由最长段的中点画药丸标签。
func (r *route) label(d *core.Drawing, text string, label, labelBg, labelBorder core.Color, alpha float32) {
	points := r.points
	best, length := 1, -1.0
	for k := 1; k < len(points); k++ {
		if l := math.Hypot(points[k].x-points[k-1].x, points[k].y-points[k-1].y); l > length {
			best, length = k, l
		}
	}
	a, b := points[best-1], points[best]
	x, y := (a.x+b.x)/2, (a.y+b.y)/2
	font := core.Font{Size: labelSize, Weight: 500}
	m := d.TextMetrics(text, font)
	tw, th := float64(m.Width), float64(m.LineHeight)
	box := struct{ x, y, w, h float32 }{float32(x - tw/2 - labelPadX), float32(y - th/2 - labelPadY), float32(tw + 2*labelPadX), float32(th + 2*labelPadY)}
	d.BeginPath().RoundRect(box.x, box.y, box.w, box.h, box.h/2)
	d.SetFillColor(labelBg.Faded(alpha)).Fill()
	d.SetStrokeColor(labelBorder.Faded(alpha)).SetLineWidth(1).Stroke()
	d.SetFillColor(label.Faded(alpha)).SetTextBaseline(core.BaselineTop).
		FillText(text, box.x+labelPadX, box.y+labelPadY, font)
}

// corners 对每个内部点调用 bend：曲线离开入段的位置、该点作控制点、
// 加入出段的位置。拐弯最多占相邻段的一半，短促的折跳变成平滑台阶。
func corners(points []point, weights []float64, bend func(a, control, b point, wa, wc, wb float64)) {
	for k := 1; k+1 < len(points); k++ {
		prev, at, next := points[k-1], points[k], points[k+1]
		in, out := math.Hypot(at.x-prev.x, at.y-prev.y), math.Hypot(next.x-at.x, next.y-at.y)
		if in == 0 || out == 0 {
			continue
		}
		r := math.Min(cornerRadius, math.Min(in, out)/2)
		wa := weights[k] - (weights[k]-weights[k-1])*r/in
		wb := weights[k] + (weights[k+1]-weights[k])*r/out
		bend(point{at.x - (at.x-prev.x)/in*r, at.y - (at.y-prev.y)/in*r}, at,
			point{at.x + (next.x-at.x)/out*r, at.y + (next.y-at.y)/out*r}, wa, weights[k], wb)
	}
}
