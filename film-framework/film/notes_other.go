//go:build !js || !wasm

package film

// loadMarks native 端不持久化（仅本次会话内存）。
func loadMarks() []float64 { return nil }

// saveMarks native 端不落盘。
func saveMarks(ts []float64) {}
