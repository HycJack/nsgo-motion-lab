//go:build js && wasm

package keyframe

import "nsgofilm/filmkit"

// Recording reports whether the recorder is driving: the recorder sets the JS
// global window.__filmT to a number of seconds before the app's own scripts run,
// and filmkit.Now() then returns it instead of anything else.
//
// A take does not need this to follow a film, only to know which of the two
// drivers is in charge — see Driven on the clock side. Both answers are the same
// question asked of filmkit, and this is the name the take code reads better as.
func Recording() bool { return filmkit.Recording() }
