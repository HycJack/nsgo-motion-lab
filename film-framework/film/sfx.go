package film

import (
	"encoding/json"
	"fmt"
)

// Sfx 是一枚音效：影片时钟走到 T 的那一刻响起。
// 观看和审片模式能听到（录制模式不出声，声音走 tools/mux.py 后期合成）。
type Sfx struct {
	T     float64 // 触发时刻（秒，影片时钟）
	Asset string  // asset: 路径（resources/ 下的短音频）
	Gain  float32 // 音量，0 = 1
}

// SfxJSON 把音效表排成 JSON——交给 tools/mux.py 做成片合成，
// 和 Film.Sfx 保持同一份事实源。
func SfxJSON(s []Sfx) string {
	type row struct {
		T    float64 `json:"t"`
		File string  `json:"file"`
		Gain float32 `json:"gain"`
	}
	out := make([]row, 0, len(s))
	for _, x := range s {
		g := x.Gain
		if g <= 0 {
			g = 1
		}
		out = append(out, row{T: x.T, File: x.Asset, Gain: g})
	}
	b, _ := json.MarshalIndent(out, "", "  ")
	return fmt.Sprintf("%s", b)
}
