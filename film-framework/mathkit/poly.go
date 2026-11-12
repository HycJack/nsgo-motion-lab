// Package mathkit 给 film-framework 补上 manim 式数学动画的原语：
// 坐标系、函数图像、公式（LaTeX 预生成路径的逐笔写与变形）、3D 投影。
//
// 所有函数都是无状态的纯绘制：同一个进度画出来的东西永远一样，
// 进度一律由调用方从影片时钟算（filmkit.Seg / Ease*），透明度走
// film.A() 乘幕层透明度后传进来。
//
// 坐标约定：一律收屏幕坐标（y 向下），和 core.Drawing 一致。
package mathkit

import (
	"math"

	"github.com/oneo-me/nsgo/core"
)

// Poly 是一条折线：Pts 是扁平的 xy 对。Closed 的折线末段自动回起点，
// 用于公式轮廓（带孔的字形靠 even-odd 填充）；开放的用于函数曲线。
type Poly struct {
	Pts    []float32
	Closed bool
}

func (p Poly) X(i int) float32 { return p.Pts[2*i] }
func (p Poly) Y(i int) float32 { return p.Pts[2*i+1] }

// Length 是折线的弧长（Closed 把回起点的那段也算上）。
func (p Poly) Length() float64 {
	n := len(p.Pts) / 2
	if n < 2 {
		return 0
	}
	var l float64
	for i := 1; i < n; i++ {
		l += seg64(p, i-1, i)
	}
	if p.Closed {
		l += seg64(p, n-1, 0)
	}
	return l
}

// Stroke 按弧长画折线到 length 为止（length<0 画满）——
// 逐笔写（Write）动画就是把全长乘进度传进来。
// 自带 BeginPath/Stroke，一次调用一段完整路径。
func (p Poly) Stroke(d *core.Drawing, length float64) {
	n := len(p.Pts) / 2
	if n < 2 {
		return
	}
	d.BeginPath()
	d.MoveTo(p.X(0), p.Y(0))
	var acc float64
	for i := 1; i < n; i++ {
		seg := seg64(p, i-1, i)
		if length >= 0 && acc+seg > length && seg > 0 {
			t := (length - acc) / seg
			d.LineTo(mix(p.X(i-1), p.X(i), t), mix(p.Y(i-1), p.Y(i), t))
			d.Stroke()
			return
		}
		acc += seg
		d.LineTo(p.X(i), p.Y(i))
	}
	if p.Closed {
		d.ClosePath()
	}
	d.Stroke()
}

// Resample 把折线按弧长均分为 n 点（Closed 的第 n 点回到起点），
// Morph 的对位就建立在等长重采样上。
func (p Poly) Resample(n int) Poly {
	nPts := len(p.Pts) / 2
	if n < 2 || nPts < 2 {
		return Poly{Pts: append([]float32{}, p.Pts...), Closed: p.Closed}
	}
	segs := nPts - 1
	if p.Closed {
		segs = nPts
	}
	segLen := func(i int) float64 {
		j := i + 1
		if j == nPts {
			j = 0
		}
		return seg64(p, i, j)
	}
	total := 0.0
	for i := 0; i < segs; i++ {
		total += segLen(i)
	}
	if total <= 0 {
		return Poly{Pts: append([]float32{}, p.Pts...), Closed: p.Closed}
	}
	out := make([]float32, 0, (n+1)*2)
	out = append(out, p.X(0), p.Y(0))
	// 开放折线 n 点只有 n-1 个间隔；闭合的 n 点正好 n 段。
	divs := n
	if !p.Closed {
		divs = n - 1
	}
	step := total / float64(divs)
	last := nPts - 1
	seg, acc := 0, 0.0
	for k := 1; k <= divs; k++ {
		target := step * float64(k)
		for seg < segs && acc+segLen(seg) < target {
			acc += segLen(seg)
			seg++
		}
		if seg >= segs {
			out = append(out, p.X(last), p.Y(last))
			break
		}
		s := segLen(seg)
		t := 0.0
		if s > 0 {
			t = (target - acc) / s
		}
		j := seg + 1
		if j == nPts {
			j = 0
		}
		out = append(out, mix(p.X(seg), p.X(j), t), mix(p.Y(seg), p.Y(j), t))
	}
	return Poly{Pts: out, Closed: false}
}

// Centroid 是折线的平均点（Morph 按它从左到右配对轮廓）。
func (p Poly) Centroid() (float64, float64) {
	n := len(p.Pts) / 2
	if n == 0 {
		return 0, 0
	}
	var sx, sy float64
	for i := 0; i < n; i++ {
		sx += float64(p.X(i))
		sy += float64(p.Y(i))
	}
	return sx / float64(n), sy / float64(n)
}

func seg64(p Poly, i, j int) float64 {
	dx := float64(p.X(j) - p.X(i))
	dy := float64(p.Y(j) - p.Y(i))
	return math.Sqrt(dx*dx + dy*dy)
}

func mix(a, b float32, t float64) float32 {
	return a + (b-a)*float32(t)
}

// smooth 是 smoothstep，把 0..1 之外钳成 0..1。
func smooth(x float64) float64 {
	if x <= 0 {
		return 0
	}
	if x >= 1 {
		return 1
	}
	return x * x * (3 - 2*x)
}

// clamp01 把 x 钳到 0..1。
func clamp01(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}
