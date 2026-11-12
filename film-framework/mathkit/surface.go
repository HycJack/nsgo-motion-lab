package mathkit

import (
	"sort"

	"github.com/oneo-me/nsgo/core"
)

// SurfaceShaded 画参数曲面的实心面片：每个四边形按法线朝向选色打光，
// 按深度排序（painter's sort）。
//
// base 是正面（法线朝向光源那侧）的颜色，back 是反面——
// 卷曲的面（展开中的柱面）正反两色，内外立现。
// light 是光照方向（世界坐标，内部归一化，指向光源）；
// 亮度 = 0.35 环境光 + 0.65 的漫反射。
// p 是铺开进度，面片按网格行序展开，铺好的部分做深度排序。
// 想要「面 + 线」就再叠一层低透明度的 SurfaceWireframe。
func (a *Axes3D) SurfaceShaded(d *core.Drawing, fn func(u, v float64) [3]float64,
	u0, u1, v0, v1 float64, nu, nv int, p float64, base, back core.Color, light [3]float64, alpha float32) {
	if alpha <= 0.001 || p <= 0 || nu < 1 || nv < 1 {
		return
	}
	L := norm3(light)
	if len3(L) == 0 {
		L = [3]float64{0.4, 0.8, 0.45}
	}
	pt := func(i, j int) [3]float64 {
		return fn(u0+(u1-u0)*float64(i)/float64(nu), v0+(v1-v0)*float64(j)/float64(nv))
	}

	type quad struct {
		pts    [4][2]float32
		depth  float64
		bright float64
		back   bool
	}
	// 铺开进度按行走：先铺满的行先出现。
	rowF := p * float64(nu)
	rows := int(rowF)
	var quads []quad
	for i := 0; i < nu; i++ {
		if i > rows {
			break
		}
		// 本行内部的列进度（最后一行铺一半）。
		colF := 1.0
		if i == rows {
			colF = rowF - float64(i)
			if colF <= 0 {
				break
			}
		}
		cols := nv
		if colF < 1 {
			cols = int(float64(nv) * colF)
		}
		for j := 0; j < cols; j++ {
			p00, p10, p11, p01 := pt(i, j), pt(i+1, j), pt(i+1, j+1), pt(i, j+1)
			var q quad
			s00, d00, ok0 := a.project(p00)
			s10, d10, ok1 := a.project(p10)
			s11, d11, ok2 := a.project(p11)
			s01, d01, ok3 := a.project(p01)
			if !ok0 || !ok1 || !ok2 || !ok3 {
				continue
			}
			q.pts = [4][2]float32{s00, s10, s11, s01}
			q.depth = (d00 + d10 + d11 + d01) / 4
			// 法线：两条邻边叉积（世界坐标）。
			e1 := sub3(p10, p00)
			e2 := sub3(p01, p00)
			n := norm3(cross3(e1, e2))
			// 有向光：正对光源用 base，背对用 back，各侧独立打光。
			if s := dot(n, L); s >= 0 {
				q.bright = 0.35 + 0.65*s
			} else {
				q.bright = 0.35 - 0.65*s // 反面：亮度取 |dot|
				q.back = true
			}
			quads = append(quads, q)
		}
	}
	sort.SliceStable(quads, func(i, j int) bool { return quads[i].depth > quads[j].depth })

	for _, q := range quads {
		col := base
		if q.back {
			col = back
		}
		d.SetFillColor(shade(col, q.bright))
		d.SetAlpha(alpha)
		d.BeginPath()
		d.MoveTo(q.pts[0][0], q.pts[0][1])
		for k := 1; k < 4; k++ {
			d.LineTo(q.pts[k][0], q.pts[k][1])
		}
		d.ClosePath()
		d.Fill()
	}
	d.SetAlpha(1)
}

// shade 把颜色按亮度系数缩放 RGB（alpha 不动）。
func shade(c core.Color, f float64) core.Color {
	r, g, b, a := c.Components()
	cl := func(v float64) uint32 { return uint32(min(max(v, 0), 1)*255 + 0.5) }
	return core.Color(cl(r*f)<<24 | cl(g*f)<<16 | cl(b*f)<<8 | cl(a))
}

// Label3D 在世界点 pos 写一行字（billboard：永远是屏幕朝向）。
// align 0 左 / 0.5 中 / 1 右 相对锚点；点在相机后不画。
// 相机转，标签跟着走——给轴、曲线端点做标注用。
func (a *Axes3D) Label3D(d *core.Drawing, s string, pos [3]float64, size float32, c core.Color, alpha float32, align float32) {
	if alpha <= 0.001 || s == "" {
		return
	}
	sp, _, ok := a.project(pos)
	if !ok {
		return
	}
	Label(d, s, sp[0], sp[1], size, c, alpha, align)
}
