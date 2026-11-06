//go:build js && wasm

package filmkit

import "syscall/js"

// start is performance.now() (in seconds) at first use without an override.
var start = -1.0

// Now is the film time in seconds.
//
// When the recorder sets the JS global __filmT (a number of seconds), that
// value wins, so frames can be stepped deterministically — the recorder
// drives the film, the wall clock only fills in when nobody is recording.
// Every animation a film scene draws must read Now (directly or through a
// signal fed from it) and nothing else, or the recording will not be
// reproducible.
func Now() float64 {
	g := js.Global()
	if v := g.Get("__filmT"); v.Type() == js.TypeNumber {
		return v.Float()
	}
	perf := g.Get("performance")
	if perf.Type() == js.TypeUndefined {
		return 0
	}
	now := perf.Call("now").Float() / 1000
	if start < 0 {
		start = now
	}
	return now - start
}
