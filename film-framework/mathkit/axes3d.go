package mathkit

import (
	"math"
	"sort"

	"github.com/oneo-me/nsgo/core"
)

// Axes3D 是 3D 坐标系：针孔投影在 CPU 上做（和 manim 同一个思路），
// 画到 core.Canvas 上。没有深度缓冲——线和面按深度排序（painter's
// sort），对线框足够。
type Axes3D struct {
	X, Y, Z   [2]float64 // 数据范围
	At        [2]float32 // 画布左上角
	Size      [2]float32 // 画布宽高
	Eye, Look [3]float64 // 相机位置与注视点
	Lens      float64    // 焦距系数，0 = 1
}

// Project 把世界点投到屏幕。cz 是相机坐标的深度（用于排序）；
// 点在相机背后或太近时 ok=false。
func (a *Axes3D) Project(x, y, z float64) (sx, sy float32, cz float64, ok bool) {
	const near = 0.05
	fwd, right, up := a.basis()
	rel := [3]float64{x - a.Eye[0], y - a.Eye[1], z - a.Eye[2]}
	cx := dot(rel, right)
	cy := dot(rel, up)
	cz = dot(rel, fwd)
	if cz <= near {
		return 0, 0, cz, false
	}
	focal := a.Lens * float64(a.Size[0]+a.Size[1]) / 2
	if a.Lens <= 0 {
		focal = float64(a.Size[0]+a.Size[1]) / 2
	}
	hx := float64(a.Size[0]) / 2
	hy := float64(a.Size[1]) / 2
	sx = a.At[0] + float32(hx+cx/cz*focal)
	sy = a.At[1] + float32(hy-cy/cz*focal)
	return sx, sy, cz, true
}

// basis 是相机坐标系：前、右、上三个单位向量。
func (a *Axes3D) basis() (fwd, right, up [3]float64) {
	fwd = norm3(sub3(a.Look, a.Eye))
	if len3(fwd) == 0 {
		fwd = [3]float64{0, 0, 1}
	}
	// 世界 up 取 (0,1,0)；视线贴着它时退化成 (0,0,1)。
	worldUp := [3]float64{0, 1, 0}
	if math.Abs(dot(fwd, worldUp)) > 0.999 {
		worldUp = [3]float64{0, 0, 1}
	}
	right = norm3(cross3(fwd, worldUp))
	up = cross3(right, fwd)
	return fwd, right, up
}

// Draw 画三根轴：穿过数据原点（不在范围内就钳到最近的边），
// 从原点向两端长出来，正方向带箭头，p 是建轴进度。
func (a *Axes3D) Draw(d *core.Drawing, p float64, c core.Color, alpha float32, width float32) {
	if alpha <= 0.001 || p <= 0 {
		return
	}
	origin := [3]float64{clampd(0, a.X[0], a.X[1]), clampd(0, a.Y[0], a.Y[1]), clampd(0, a.Z[0], a.Z[1])}
	ends := [][3]float64{
		{a.X[1], origin[1], origin[2]},
		{origin[0], a.Y[1], origin[2]},
		{origin[0], origin[1], a.Z[1]},
	}
	starts := [][3]float64{
		{a.X[0], origin[1], origin[2]},
		{origin[0], a.Y[0], origin[2]},
		{origin[0], origin[1], a.Z[0]},
	}
	d.SetStrokeColor(c)
	d.SetAlpha(alpha)
	d.SetLineWidth(width)
	for k := range ends {
		// 从原点向两端各长 p（半程）。
		half := [][3]float64{origin, interp3pt(origin, ends[k], p)}
		a.strokePts(d, half, 0)
		neg := [][3]float64{origin, interp3pt(origin, starts[k], p)}
		a.strokePts(d, neg, 0)
	}
	// 箭头：正方向两端（轴长出来之后浮现）。
	arrA := alpha * float32(smooth((p-0.7)/0.3))
	if arrA > 0.001 {
		for _, e := range ends {
			axisEnd := interp3pt(origin, e, 1)
			tip, _, ok := a.project(axisEnd)
			if !ok {
				continue
			}
			inner, _, ok2 := a.project(interp3pt(origin, e, 0.94))
			if !ok2 {
				continue
			}
			a.arrow3(d, inner, tip, c, arrA)
		}
	}

	// 刻度：三根轴上按 niceStep 取点，画屏幕垂直的小段
	// （刻度的屏幕方向取轴的屏幕垂线，相机转了也成立）。
	tickA := alpha * float32(smooth((p-0.6)/0.4))
	if tickA > 0.001 {
		d.SetAlpha(tickA)
		const tick = 6.0
		axes := []struct {
			rng  [2]float64
			fixe [2]float64 // 另两维固定在原点
			axis int
		}{
			{a.X, [2]float64{origin[1], origin[2]}, 0},
			{a.Y, [2]float64{origin[0], origin[2]}, 1},
			{a.Z, [2]float64{origin[0], origin[1]}, 2},
		}
		for _, ax := range axes {
			step := niceStep(ax.rng[1]-ax.rng[0], 6)
			// 轴的屏幕方向：原点 → 正端。
			p0, _, ok0 := a.project(origin)
			p1, _, ok1 := a.project(func() [3]float64 {
				q := origin
				q[ax.axis] = ax.rng[1]
				return q
			}())
			if !ok0 || !ok1 {
				continue
			}
			dx, dy := p1[0]-p0[0], p1[1]-p0[1]
			l := float32(math.Hypot(float64(dx), float64(dy)))
			if l < 1 {
				continue
			}
			// 屏幕垂线方向。
			px, py := -dy/l, dx/l
			lo := math.Ceil(ax.rng[0]/step) * step
			for v := lo; v <= ax.rng[1]+1e-9; v += step {
				if math.Abs(v) < step/2 {
					continue // 原点不给自己画刻度
				}
				q := origin
				q[ax.axis] = v
				sp, _, ok := a.project(q)
				if !ok {
					continue
				}
				d.BeginPath()
				d.MoveTo(sp[0]-float32(px)*tick, sp[1]-float32(py)*tick)
				d.LineTo(sp[0]+float32(px)*tick, sp[1]+float32(py)*tick)
				d.Stroke()
			}
		}
	}
	d.SetAlpha(1)
}

// arrow3 在屏幕平面上画一个小实心三角：从 inner 指向 tip。
func (a *Axes3D) arrow3(d *core.Drawing, inner, tip [2]float32, c core.Color, alpha float32) {
	vx, vy := tip[0]-inner[0], tip[1]-inner[1]
	l := float32(math.Hypot(float64(vx), float64(vy)))
	if l < 1 {
		return
	}
	vx, vy = vx/l, vy/l
	px, py := -vy, vx
	const s = 11.0
	d.SetFillColor(c)
	d.SetAlpha(alpha)
	d.BeginPath()
	d.MoveTo(tip[0], tip[1])
	d.LineTo(inner[0]-vx*s+px*s*0.6, inner[1]-vy*s+py*s*0.6)
	d.LineTo(inner[0]-vx*s-px*s*0.6, inner[1]-vy*s-py*s*0.6)
	d.ClosePath()
	d.Fill()
}

func interp3pt(a, b [3]float64, t float64) [3]float64 {
	return [3]float64{a[0] + (b[0]-a[0])*t, a[1] + (b[1]-a[1])*t, a[2] + (b[2]-a[2])*t}
}

// Poly3 画空间折线（参数曲线就密集采样成折线）：p 是描画进度，
// 相机背后的一段自动跳过。段按深度从远到近排序——
// 相交/盘绕的曲线（螺旋、环绕轨道）遮挡关系才正确。
func (a *Axes3D) Poly3(d *core.Drawing, pts [][3]float64, p float64, c core.Color, alpha float32, width float32) {
	if alpha <= 0.001 || p <= 0 || len(pts) < 2 {
		return
	}
	nSeg := len(pts) - 1
	nDraw := int(p * float64(nSeg))
	type seg struct {
		from, to [2]float32
		depth    float64
	}
	segs := make([]seg, 0, nDraw)
	for i := 0; i < nDraw && i < nSeg; i++ {
		from, d0, ok0 := a.project(pts[i])
		to, d1, ok1 := a.project(pts[i+1])
		if !ok0 || !ok1 {
			continue
		}
		segs = append(segs, seg{from, to, (d0 + d1) / 2})
	}
	sort.Slice(segs, func(i, j int) bool { return segs[i].depth > segs[j].depth })
	d.SetStrokeColor(c)
	d.SetAlpha(alpha)
	d.SetLineWidth(width)
	for _, s := range segs {
		d.BeginPath()
		d.MoveTo(s.from[0], s.from[1])
		d.LineTo(s.to[0], s.to[1])
		d.Stroke()
	}
	d.SetAlpha(1)
}

type surfaceLine struct {
	pts   [][2]float32
	depth float64
}

// SurfaceWireframe 画参数曲面的线框：u、v 两个方向的网格线，
// 按平均深度排序（远的后画会被近的压住……painter 是先画远的），
// p 是铺开进度——行先、列后。
func (a *Axes3D) SurfaceWireframe(d *core.Drawing, fn func(u, v float64) [3]float64,
	u0, u1, v0, v1 float64, nu, nv int, p float64, c core.Color, alpha float32, width float32) {
	if alpha <= 0.001 || p <= 0 || nu < 2 || nv < 2 {
		return
	}
	var lines []surfaceLine
	// 固定 u、沿 v 扫的线。
	for i := 0; i <= nu; i++ {
		u := u0 + (u1-u0)*float64(i)/float64(nu)
		pts := make([][3]float64, 0, nv+1)
		for j := 0; j <= nv; j++ {
			pts = append(pts, fn(u, v0+(v1-v0)*float64(j)/float64(nv)))
		}
		lines = append(lines, a.line3(pts))
	}
	// 固定 v、沿 u 扫的线。
	for j := 0; j <= nv; j++ {
		v := v0 + (v1-v0)*float64(j)/float64(nv)
		pts := make([][3]float64, 0, nu+1)
		for i := 0; i <= nu; i++ {
			pts = append(pts, fn(u0+(u1-u0)*float64(i)/float64(nu), v))
		}
		lines = append(lines, a.line3(pts))
	}
	sort.SliceStable(lines, func(i, j int) bool { return lines[i].depth > lines[j].depth })
	nDraw := p * float64(len(lines))
	d.SetStrokeColor(c)
	d.SetAlpha(alpha)
	d.SetLineWidth(width)
	for k, ln := range lines {
		frac := min(max(nDraw-float64(k), 0), 1)
		if frac <= 0 {
			break
		}
		a.stroke2(d, ln, frac)
	}
	d.SetAlpha(1)
}

// line3 把一串世界点投好、算好平均深度。
func (a *Axes3D) line3(pts [][3]float64) surfaceLine {
	ln := surfaceLine{pts: make([][2]float32, 0, len(pts))}
	var sum float64
	for _, p := range pts {
		s, depth, ok := a.project(p)
		if !ok {
			continue
		}
		ln.pts = append(ln.pts, [2]float32{s[0], s[1]})
		sum += depth
	}
	if len(ln.pts) > 0 {
		ln.depth = sum / float64(len(ln.pts))
	}
	return ln
}

func (a *Axes3D) stroke2(d *core.Drawing, ln surfaceLine, frac float64) {
	if len(ln.pts) < 2 {
		return
	}
	pl := Poly{Pts: flatten2(ln.pts)}
	pl.Stroke(d, pl.Length()*frac)
}

func (a *Axes3D) project(p [3]float64) ([2]float32, float64, bool) {
	sx, sy, cz, ok := a.Project(p[0], p[1], p[2])
	return [2]float32{sx, sy}, cz, ok
}

func (a *Axes3D) strokePts(d *core.Drawing, pts [][3]float64, _ float64) {
	for i := 1; i < len(pts); i++ {
		from, _, fromOK := a.project(pts[i-1])
		to, _, toOK := a.project(pts[i])
		if !fromOK || !toOK {
			continue
		}
		d.BeginPath()
		d.MoveTo(from[0], from[1])
		d.LineTo(to[0], to[1])
		d.Stroke()
	}
}

func flatten2(pts [][2]float32) []float32 {
	out := make([]float32, 0, len(pts)*2)
	for _, p := range pts {
		out = append(out, p[0], p[1])
	}
	return out
}

func sub3(a, b [3]float64) [3]float64 { return [3]float64{a[0] - b[0], a[1] - b[1], a[2] - b[2]} }
func dot3(a, b [3]float64) float64    { return a[0]*b[0] + a[1]*b[1] + a[2]*b[2] }
func cross3(a, b [3]float64) [3]float64 {
	return [3]float64{a[1]*b[2] - a[2]*b[1], a[2]*b[0] - a[0]*b[2], a[0]*b[1] - a[1]*b[0]}
}
func len3(a [3]float64) float64 { return math.Sqrt(dot3(a, a)) }
func norm3(a [3]float64) [3]float64 {
	l := len3(a)
	if l == 0 {
		return a
	}
	return [3]float64{a[0] / l, a[1] / l, a[2] / l}
}
func dot(a, b [3]float64) float64 { return dot3(a, b) }

// Poly3Dashed 画空间虚线：dash/gap 按屏幕像素弧长交替，虚线段按深度
// 从远到近绘制（遮挡正确），p 是描画进度（沿总弧长）。
// 3D 里「看不见的后半圈」「辅助投影线」的标准画法。
func (a *Axes3D) Poly3Dashed(d *core.Drawing, pts [][3]float64, p float64, dash, gap float32, c core.Color, alpha float32, width float32) {
	if alpha <= 0.001 || p <= 0 || len(pts) < 2 || dash <= 0 {
		return
	}
	cycle := float64(dash + gap)

	// 投影：保留 ok 标记，相机后的点断开（不跨界连画）。
	type spt struct {
		p     [2]float32
		depth float64
		ok    bool
	}
	sp := make([]spt, len(pts))
	for i, w := range pts {
		s, depth, ok := a.project(w)
		sp[i] = spt{s, depth, ok}
	}

	// 总屏幕弧长（只统计两端都 ok 的相邻对）。
	total := 0.0
	for i := 1; i < len(sp); i++ {
		if sp[i].ok && sp[i-1].ok {
			total += dist2(sp[i-1].p, sp[i].p)
		}
	}
	limit := p * total

	// 沿折线走，按 cycle 切 on/off，收集虚线段。
	type dseg struct {
		a, b  [2]float32
		depth float64
	}
	var segs []dseg
	acc := 0.0
	for i := 1; i < len(sp) && acc < limit; i++ {
		if !sp[i].ok || !sp[i-1].ok {
			continue
		}
		a2, b2 := sp[i-1].p, sp[i].p
		d2 := (sp[i-1].depth + sp[i].depth) / 2
		segLen := dist2(a2, b2)
		if segLen <= 0 {
			continue
		}
		sf := 0.0
		for sf < segLen && acc < limit {
			phase := math.Mod(acc, cycle)
			if phase >= float64(dash) { // 在 gap 里：跳到下个周期
				jump := math.Min(float64(segLen)-sf, cycle-phase)
				acc += jump
				sf += jump
				continue
			}
			run := math.Min(float64(dash)-phase, float64(segLen)-sf)
			run = math.Min(run, limit-acc)
			if run <= 0 {
				sf = segLen
				break
			}
			e := sf + run
			segs = append(segs, dseg{
				lerp2(a2, b2, sf/segLen),
				lerp2(a2, b2, e/segLen),
				d2,
			})
			acc += run
			sf = e
		}
	}
	sort.Slice(segs, func(i, j int) bool { return segs[i].depth > segs[j].depth })

	d.SetStrokeColor(c)
	d.SetAlpha(alpha)
	d.SetLineWidth(width)
	for _, g := range segs {
		d.BeginPath()
		d.MoveTo(g.a[0], g.a[1])
		d.LineTo(g.b[0], g.b[1])
		d.Stroke()
	}
	d.SetAlpha(1)
}

func dist2(a, b [2]float32) float64 {
	dx := float64(b[0] - a[0])
	dy := float64(b[1] - a[1])
	return math.Sqrt(dx*dx + dy*dy)
}

func lerp2(a, b [2]float32, t float64) [2]float32 {
	return [2]float32{a[0] + (b[0]-a[0])*float32(t), a[1] + (b[1]-a[1])*float32(t)}
}
