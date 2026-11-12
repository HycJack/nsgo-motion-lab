package flowkit

import (
	"math"

	"github.com/oneo-me/nsgo/core"
)

// clamp 来自引擎的 camera.go（布局/路由共用）。
func clamp(v, lo, hi float64) float64 { return math.Max(lo, math.Min(v, hi)) }

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

// filmkitA 是 float64 alpha 到 float32 的便捷转换。
func filmkitA(a float64) float32 { return float32(a) }

// dimFill 把填充色调暗（状态色底）。
func dimFill(c core.Color) core.Color {
	r, g, b, a := c.Components()
	ch := func(v float64) uint32 { return uint32(v*0.35*255 + 0.5) }
	return core.Color(ch(r)<<24 | ch(g)<<16 | ch(b)<<8 | uint32(a*255+0.5))
}
