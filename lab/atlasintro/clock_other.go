//go:build !js

package main

import "time"

var atlasStart = time.Now()

// atlasTime is wall time since the program started, for native runs.
func atlasTime() float64 {
	return time.Since(atlasStart).Seconds()
}
