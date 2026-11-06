//go:build !js

package filmkit

import "time"

// start is the process start; native films run on the wall clock.
var start = time.Now()

// Now is wall time since the program started, in seconds. Native runs have
// no recorder to step the film, so they play in real time — same shape as
// the web clock's fallback path.
func Now() float64 { return time.Since(start).Seconds() }
