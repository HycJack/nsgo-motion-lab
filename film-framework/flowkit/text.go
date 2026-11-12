package flowkit

import "github.com/oneo-me/nsgo/core"

// textMeasurer 是布局需要的文字测量（由 core.Drawing 满足）。
type textMeasurer interface {
	MeasureText(s string, font core.Font) float32
}
