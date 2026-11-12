//go:build js && wasm

package film

import "syscall/js"

// recordingMode 报告这部片子是不是正被录制器驱动。
//
// 录制器（shots/film.js）会在页面脚本运行前把 window.__filmT 设成 0，
// 然后逐帧改它。所以「__filmT 是个数字」= 有人在外部推时间。
//
// 是录制模式时：按 filmkit 的时钟走，声音不播，「开始」按钮也不出现
// —— 录出来的画面必须和纯画面版一模一样。
func recordingMode() bool {
	return js.Global().Get("__filmT").Type() == js.TypeNumber
}
