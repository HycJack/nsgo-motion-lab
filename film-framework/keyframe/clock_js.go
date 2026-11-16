//go:build js && wasm

package keyframe

import "syscall/js"

// Recording reports whether something outside is driving the film clock: the
// recorder sets window.__filmT to a number of seconds before the app's own
// scripts run, and filmkit.Now() then returns it instead of the wall clock.
//
// The distinction matters at mount time. A take that mounts while a recorder is
// driving must not start its own playback clock, or the two would write the same
// value every frame and the recording would flicker between them.
func Recording() bool {
	return js.Global().Get("__filmT").Type() == js.TypeNumber
}
