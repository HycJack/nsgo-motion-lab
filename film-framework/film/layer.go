package film

import (
	"time"

	"github.com/nsgo-dev/nsgo/core"

	"nsgofilm/keyframe"
)

// Layer 是叠在画布上的一层**节点树**，跟影片同一个时钟。
//
// Shot 往画布里画，Layer 往树里放 —— 一个 keyframe.Take 只能以节点树存在，
// 所以它必须是 Layer。Start/End 是影片秒数（End 为 0 表示到片尾），
// 层在自己的窗口之外不挂载：窗口之外不画，也不占帧。
//
// 坐标：图层和画布用同一个设计画幅（Film.W × Film.H）——画布把设计框等比
// 缩放并居中，图层照做（见 layerBox），所以窗口大小变了两者也不脱节。
type Layer struct {
	Start, End float64
	Node       func(*core.Scope) core.Node
}

// Over 把一支 take 挂到画布上：从 start 秒显示到 end 秒（end 为 0 表示到
// 片尾），并且让这支 take 自己的 0 秒落在影片的 start 秒上 —— 也就是把它当
// 一幕用，而不是当作从片头就在跑的底图。
//
// 两个数一起设是有意的：Delay 定时间（seek 到哪里），Start/End 定存亡
// （这一层在不在）。所以拖动进度条、暂停、逐帧，画面都能对上。
func Over(t *keyframe.Take, start, end float64) Layer {
	t.Delay(time.Duration(start * float64(time.Second)))
	return Layer{Start: start, End: end, Node: t.Frame}
}

// shows 报告影片时间 t 是否落在这一层的窗口里；两个边界都为 0 表示铺满全片。
func (l Layer) shows(t float64) bool {
	if l.Start == 0 && l.End == 0 {
		return true
	}
	return t >= l.Start && (l.End == 0 || t <= l.End)
}

// layerBox 把一层放进画布的坐标系：设计框 W×H 在可用区域里等比缩放、居中，
// 和 canvas 用的 min(W/…, H/…) 是同一个算法，所以节点和画布上画的东西对得上。
// inner 由调用方建好（通常是一个 core.Show），缩放值只在这里建一次。
func layerBox(cx *core.Scope, f Film, inner core.Node) core.Node {
	scale := core.NewSignal(cx, 1.0)
	// 量到窗口就立刻定缩放；只在新值不同时写，免得布局回调把自己再叫一遍。
	fit := func(b core.Bounds) {
		s := min(float64(b.Size.Width)/float64(f.W), float64(b.Size.Height)/float64(f.H))
		if s > 0 && s != scale.Get() {
			scale.Set(s)
		}
	}
	return core.View{
		Style: core.Style{
			Position:       core.Absolute,
			Left:           core.Pt(0),
			Top:            core.Pt(0),
			Right:          core.Pt(0),
			Bottom:         core.Pt(0),
			JustifyContent: core.JustifyCenter,
			AlignItems:     core.AlignCenter,
		},
		OnLayout: fit,
		Children: []core.Node{core.View{
			Class: "overflow-hidden",
			StyleFunc: func() core.Style {
				return core.Style{
					Width:     core.Pt(core.Points(f.W)),
					Height:    core.Pt(core.Points(f.H)),
					Transform: core.Transform{Scale: float32(scale.Get())},
				}
			},
			Children: []core.Node{inner},
		}},
	}
}
