package filmkit

import "math"

// 运动质量三件套：弹簧（spring）、贝塞尔缓动（Bezier）、多站关键帧（KF）。
// 全部是纯函数——同一个 t 永远同一个值，录制和观看逐像素一致。

// Spring 是阻尼谐振子的参数。零值 = React Native 的默认手感
// （stiffness 170、damping 26、mass 1）。
type Spring struct {
	Stiffness float64 // 0 = 170
	Damping   float64 // 0 = 26
	Mass      float64 // 0 = 1
	Vel0      float64 // 初速度（进度/秒），一般留 0
}

// At 返回弹簧在 t 秒的进度：从 0 静止出发奔向 1。
// 欠阻尼会过冲超过 1 再回落——那是它好玩的全部理由；
// 想要"不超"就把 Damping 拉大（ζ ≥ 1 时单调）。
// 闭式解，确定性由纯算术保证。
func (s Spring) At(t float64) float64 {
	k := s.Stiffness
	if k <= 0 {
		k = 170
	}
	c := s.Damping
	if c <= 0 {
		c = 26
	}
	m := s.Mass
	if m <= 0 {
		m = 1
	}
	if t <= 0 {
		return 0
	}
	omega := math.Sqrt(k / m)              // 无阻尼角频率
	zeta := c / (2 * math.Sqrt(k*m))       // 阻尼比
	x0, v0 := -1.0, s.Vel0                 // 位移约定：目标 1，起点 0 → 偏移 -1

	if zeta < 1 { // 欠阻尼：衰减振荡
		wd := omega * math.Sqrt(1-zeta*zeta)
		e := math.Exp(-zeta * omega * t)
		return 1 + e*(x0*math.Cos(wd*t)+(v0+x0*zeta*omega)/wd*math.Sin(wd*t))
	}
	if zeta == 1 { // 临界阻尼：最快的不振荡
		e := math.Exp(-omega * t)
		return 1 + e*(x0+(v0+x0*omega)*t)
	}
	// 过阻尼：两个衰减指数。
	wd := omega * math.Sqrt(zeta*zeta-1)
	r1 := -zeta*omega + wd
	r2 := -zeta*omega - wd
	den := r1 - r2
	c1 := (v0 - r2*x0) / den
	c2 := (r1*x0 - v0) / den
	return 1 + c1*math.Exp(r1*t) + c2*math.Exp(r2*t)
}

// Bezier 返回 CSS 同款的 cubic-bezier(x1, y1, x2, y2) 缓动函数。
// x 分量单调（CSS 的约束），用二分反解参数再取 y——确定性二分，
// 迭代次数固定。
func Bezier(x1, y1, x2, y2 float64) func(float64) float64 {
	bx := func(t float64) float64 { return bezCoord(t, x1, x2) }
	by := func(t float64) float64 { return bezCoord(t, y1, y2) }
	return func(x float64) float64 {
		if x <= 0 {
			return 0
		}
		if x >= 1 {
			return 1
		}
		// 二分解 bx(t) = x（bx 单调）。
		lo, hi := 0.0, 1.0
		for i := 0; i < 30; i++ {
			mid := (lo + hi) / 2
			if bx(mid) < x {
				lo = mid
			} else {
				hi = mid
			}
		}
		return by((lo + hi) / 2)
	}
}

// bezCoord 是一维三次贝塞尔：P0=0, P1=u, P2=v, P3=1。
func bezCoord(t, u, v float64) float64 {
	o := 1 - t
	return 3*o*o*t*u + 3*o*t*t*v + t*t*t
}

// KF 是多站关键帧插值：times 严格递增、values 与之一一对应，
// 段内用 ease 缓动（nil = 线性），两端外钳。
// "起-升-停-落"这类节奏，一段 KF 顶过去三四段 Seg+Lerp。
func KF(t float64, times []float64, values []float64, ease func(float64) float64) float64 {
	n := len(times)
	if n == 0 || len(values) != n {
		return 0
	}
	if t <= times[0] {
		return values[0]
	}
	if t >= times[n-1] {
		return values[n-1]
	}
	i := 1
	for i < n-1 && t > times[i] {
		i++
	}
	p := (t - times[i-1]) / (times[i] - times[i-1])
	if ease != nil {
		p = ease(p)
	}
	return values[i-1] + (values[i]-values[i-1])*p
}
