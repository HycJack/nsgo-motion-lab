//go:build !js || !wasm

package film

// recordingMode：非 web 平台没有 __filmT，永远是「观看模式」。
func recordingMode() bool { return false }
