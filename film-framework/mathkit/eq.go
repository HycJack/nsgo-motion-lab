package mathkit

import (
	"math"
	"sort"

	"github.com/oneo-me/nsgo/core"
)

// Eq 是一条公式：轮廓折线的集合，坐标是屏幕朝向（y 向下），
// 原点在公式的左基线，字号 100 时的像素。由 tools/eq.py 从 LaTeX
// 预生成（matplotlib mathtext → 轮廓多边形），不要手写。
type Eq struct {
	W, H  float32
	Lines []Poly
}

const (
	// eqResample 是 Morph 时每条轮廓重采样的点数。
	eqResample = 48
	// eqScaleRef 是生成公式时用的字号，和 eq.py 对齐。
	eqScaleRef = 100
)

// EqDraw 以 (x,y) 为左基线画整条公式，scale 是相对字号 100 的倍数。
// 轮廓用 even-odd 填充——带孔的字形（e、o、∂）才能镂空。
func EqDraw(d *core.Drawing, e Eq, x, y, scale float32, c core.Color, alpha float32) {
	if alpha <= 0.001 {
		return
	}
	d.SetFillColor(c)
	d.SetAlpha(alpha)
	d.BeginPath()
	eqPath(d, e, x, y, scale)
	d.FillEvenOdd()
	d.SetAlpha(1)
}

// EqWrite 是逐笔写：按全部轮廓的总弧长依次描边，p 是进度 0..1。
// 预算跨轮廓累计——写完一条轮廓再写下一条，而不是所有笔画齐头并进。
func EqWrite(d *core.Drawing, e Eq, x, y, scale float32, p float64, c core.Color, alpha float32, width float32) {
	if alpha <= 0.001 || p <= 0 {
		return
	}
	budget := eqLength(e) * float64(scale) * p
	d.SetStrokeColor(c)
	d.SetAlpha(alpha)
	d.SetLineWidth(width)
	for _, ln := range e.Lines {
		if budget <= 0 {
			break
		}
		s := scalePoly(ln, x, y, scale)
		l := float64(scale) * ln.Length()
		if budget >= l {
			s.Stroke(d, -1)
			budget -= l
			continue
		}
		s.Stroke(d, budget)
		break
	}
	d.SetAlpha(1)
}

// EqMorph 是公式变形（manim 的 Transform）：轮廓做最优对应
// （贪心全局最小代价，见 pairContours），等弧长重采样后逐点插值。
// 多余轮廓复制对方——在 even-odd 填充下端点处自然消失/浮现。
// 细描边勾住轮廓，变形中段不会糊。
func EqMorph(d *core.Drawing, from, to Eq, x, y, scale float32, p float64, c core.Color, alpha float32, width float32) {
	if alpha <= 0.001 || p <= 0 {
		return
	}
	pairs := pairContours(sortedContours(from), sortedContours(to))
	if len(pairs) == 0 {
		return
	}
	d.SetFillColor(c)
	d.SetAlpha(alpha)
	d.BeginPath()
	for _, pr := range pairs {
		a := pr.From.Resample(eqResample)
		b := pr.To.Resample(eqResample)
		for k := 0; k < eqResample; k++ {
			sx := mix(lerpAt(a, k, x, scale), lerpAt(b, k, x, scale), p)
			sy := mix(lerpYAt(a, k, y, scale), lerpYAt(b, k, y, scale), p)
			if k == 0 {
				d.MoveTo(sx, sy)
			} else {
				d.LineTo(sx, sy)
			}
		}
		d.ClosePath()
	}
	d.FillEvenOdd()
	d.SetStrokeColor(c)
	d.SetLineWidth(width)
	d.BeginPath()
	for _, pr := range pairs {
		a := pr.From.Resample(eqResample)
		b := pr.To.Resample(eqResample)
		for k := 0; k < eqResample; k++ {
			sx := mix(lerpAt(a, k, x, scale), lerpAt(b, k, x, scale), p)
			sy := mix(lerpYAt(a, k, y, scale), lerpYAt(b, k, y, scale), p)
			if k == 0 {
				d.MoveTo(sx, sy)
			} else {
				d.LineTo(sx, sy)
			}
		}
		d.ClosePath()
	}
	d.Stroke()
	d.SetAlpha(1)
}

// EqWriteGradient 是带渐变的逐笔写：笔画颜色沿总弧长从 c1 渐到 c2。
// 已写部分的颜色不随进度变化（颜色"焊"在弧长位置上）。
func EqWriteGradient(d *core.Drawing, e Eq, x, y, scale float32, p float64, c1, c2 core.Color, alpha float32, width float32) {
	if alpha <= 0.001 || p <= 0 {
		return
	}
	total := eqLength(e) * float64(scale)
	budget := total * p
	done := 0.0
	d.SetAlpha(alpha)
	d.SetLineWidth(width)
	for _, ln := range e.Lines {
		if budget <= 0 {
			break
		}
		s := scalePoly(ln, x, y, scale)
		l := s.Length()
		draw := math.Min(l, budget)
		// 均匀重采样后按小段描，每段取其弧长位置对应的渐变色。
		const K = 64
		rs := s.Resample(K)
		m := len(rs.Pts) / 2
		step := l / float64(m-1)
		var prev [2]float32
		var acc float64
		for k := 0; k < m; k++ {
			cur := [2]float32{rs.X(k), rs.Y(k)}
			if k > 0 {
				end := math.Min(acc+step, draw)
				if end <= acc {
					break
				}
				mid := (acc + end) / 2
				t := clamp01((done + mid) / max(total, 0.001))
				d.SetStrokeColor(mixColor(c1, c2, t))
				d.BeginPath()
				d.MoveTo(prev[0], prev[1])
				d.LineTo(cur[0], cur[1])
				d.Stroke()
				if end < draw {
					acc += step
				} else {
					break
				}
			}
			prev = cur
			if k == 0 {
				acc = 0
			}
		}
		done += l
		budget -= l
		if budget <= 0 {
			break
		}
	}
	d.SetAlpha(1)
}

// mixColor 按 t 混合两个颜色（含 alpha）。
func mixColor(a, b core.Color, t float64) core.Color {
	ar, ag, ab, aa := a.Components()
	br, bg, bb, ba := b.Components()
	m := func(u, v float64) uint32 { return uint32((u + (v-u)*t)*255 + 0.5) }
	return core.Color(m(ar, br)<<24 | m(ag, bg)<<16 | m(ab, bb)<<8 | m(aa, ba))
}

// eqPath 把公式的全部轮廓铺进当前路径（缩放平移到 x,y 基线）。
func eqPath(d *core.Drawing, e Eq, x, y, scale float32) {
	for _, ln := range e.Lines {
		for k := 0; k < len(ln.Pts)/2; k++ {
			sx := x + ln.X(k)*scale
			sy := y + ln.Y(k)*scale // 生成侧已是 y 向下：负 y = 基线上方
			if k == 0 {
				d.MoveTo(sx, sy)
			} else {
				d.LineTo(sx, sy)
			}
		}
		if ln.Closed {
			d.ClosePath()
		}
	}
}

func scalePoly(ln Poly, x, y, scale float32) Poly {
	out := Poly{Pts: make([]float32, len(ln.Pts)), Closed: ln.Closed}
	for k := 0; k < len(ln.Pts)/2; k++ {
		out.Pts[2*k] = x + ln.X(k)*scale
		out.Pts[2*k+1] = y + ln.Y(k)*scale
	}
	return out
}

func eqLength(e Eq) float64 {
	var total float64
	for _, ln := range e.Lines {
		total += ln.Length()
	}
	return total
}

// sortedContours 把轮廓按质心 x 从左到右排——Morph 配对的依据。
func sortedContours(e Eq) []Poly {
	out := append([]Poly{}, e.Lines...)
	sort.Slice(out, func(i, j int) bool {
		ix, _ := out[i].Centroid()
		jx, _ := out[j].Centroid()
		return ix < jx
	})
	return out
}

// lerpAt/lerpYAt 取重采样后第 k 点换算到屏幕（等宽采样的行内展开）。
func lerpAt(p Poly, k int, x, scale float32) float32 { return x + p.X(k)*scale }
func lerpYAt(p Poly, k int, y, scale float32) float32 {
	return y + p.Y(k)*scale
}
