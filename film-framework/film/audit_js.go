//go:build js && wasm

package film

import (
	"strings"
	"syscall/js"
)

// auditMode 报告是否以审片模式打开：URL 带 ?audit=1。
//
// 审片模式在画面下方挂一条控制条：播放/暂停、当前时间、可点可拖的
// 进度条——跳转带着声音走，用来核对台词、音效和画面的对位。
// 录制器不受影响：它开的是干净页面，不带这个参数。
func auditMode() bool {
	return strings.Contains(js.Global().Get("location").Get("search").String(), "audit=1")
}

// timelineMode 报告是否以时间轴总览打开：URL 带 ?timeline=1。
// 每幕一帧的 contact sheet，点哪张卡就跳进审片模式的那一幕。
func timelineMode() bool {
	return strings.Contains(js.Global().Get("location").Get("search").String(), "timeline=1")
}
