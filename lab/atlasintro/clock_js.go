//go:build js && wasm

package main

import (
	"syscall/js"
)

// atlasStart is the performance.now() at first use, in seconds.
var atlasStart = -1.0

// atlasTime is the film time. When the capture recorder sets the JS
// global __atlasT (a number of seconds), that wins, so frames can be
// stepped deterministically; otherwise it is wall time on the page clock.
func atlasTime() float64 {
	g := js.Global()
	if v := g.Get("__atlasT"); v.Type() == js.TypeNumber {
		return v.Float()
	}
	perf := g.Get("performance")
	if perf.Type() == js.TypeUndefined {
		return 0
	}
	now := perf.Call("now").Float() / 1000
	if atlasStart < 0 {
		atlasStart = now
	}
	return now - atlasStart
}
