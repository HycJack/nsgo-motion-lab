package film

// 审片批注：C 键在当前时刻打一个标记，X 删最近的，O 导出 JSON
// （时间戳清单——交给后续的修改者：人或 AI）。
// web 端存 localStorage（刷新不丢）；native 端仅内存。

import "fmt"

// notesJSON 把标记清单排成导出 JSON。
func notesJSON(ts []float64, dur float64) string {
	out := fmt.Sprintf(`{"duration": %.2f, "markers": [`, dur)
	for i, t := range ts {
		if i > 0 {
			out += ", "
		}
		out += fmt.Sprintf(`{"t": %.2f, "note": ""}`, t)
	}
	return out + "]}"
}

// addMark 加标记：0.2s 内去重，按时间排序。
func addMark(ts []float64, t float64) []float64 {
	for _, x := range ts {
		if absF(x-t) < 0.2 {
			return ts
		}
	}
	t += 0 // 保持 float64
	out := make([]float64, 0, len(ts)+1)
	inserted := false
	for _, x := range ts {
		if !inserted && x > t {
			out = append(out, t)
			inserted = true
		}
		out = append(out, x)
	}
	if !inserted {
		out = append(out, t)
	}
	return out
}

// dropNearest 删掉距 t 最近的标记（0.3s 内才算近）。
func dropNearest(ts []float64, t float64) ([]float64, bool) {
	best, bd := -1, 0.31
	for i, x := range ts {
		if d := absF(x - t); d < bd {
			bd, best = d, i
		}
	}
	if best < 0 {
		return ts, false
	}
	out := append(ts[:best:best], ts[best+1:]...)
	return out, true
}

func absF(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
