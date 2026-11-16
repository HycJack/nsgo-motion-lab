//go:build !js

package filmkit

import "time"

// start is the process start; native films run on the wall clock.
var start = time.Now()

// Recording reports whether a recorder is driving: native runs have no browser
// to set __filmT, so the recorder never drives them.
func Recording() bool { return false }

// Now is the film time in seconds: the clock a film published for the shots
// inside it, or the wall clock when nothing is driving.
func Now() float64 {
	if t, ok := clockNow(); ok {
		return t
	}
	return Wall()
}

// Wall is the platform's own time since the program started, ignoring any
// published clock. A film's frame loop advances time with this: reading Now
// there would return the loop's own output from the frame before, and time
// would stand still.
func Wall() float64 { return time.Since(start).Seconds() }
