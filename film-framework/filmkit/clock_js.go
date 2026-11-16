//go:build js && wasm

package filmkit

import "syscall/js"

// start is performance.now() (in seconds) at first use without an override.
var start = -1.0

// Recording reports whether the recorder is driving: it sets the JS global
// __filmT to a number of seconds before the page's scripts run, and that value
// then outranks every other source of film time.
func Recording() bool {
	return js.Global().Get("__filmT").Type() == js.TypeNumber
}

// Now is the film time in seconds, in order of authority: the recorder's
// __filmT, then the clock the film published for the shots inside it, then the
// wall clock. Every animation a film draws must read Now (directly or through a
// signal fed from it) and nothing else, or the recording will not be
// reproducible.
func Now() float64 {
	if v := js.Global().Get("__filmT"); v.Type() == js.TypeNumber {
		return v.Float()
	}
	if t, ok := clockNow(); ok {
		return t
	}
	return Wall()
}

// Wall is the platform's own time, ignoring both the recorder and any published
// clock. A film's frame loop advances time with this: reading Now there would
// return the loop's own output from the frame before, and time would stand
// still.
func Wall() float64 {
	perf := js.Global().Get("performance")
	if perf.Type() == js.TypeUndefined {
		return 0
	}
	now := perf.Call("now").Float() / 1000
	if start < 0 {
		start = now
	}
	return now - start
}
