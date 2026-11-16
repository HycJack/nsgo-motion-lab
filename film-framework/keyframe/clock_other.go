//go:build !(js && wasm)

package keyframe

import "nsgofilm/filmkit"

// Recording reports whether the recorder is driving: off the web there is no
// page to inject a clock into, and tests drive takes with Seek directly.
func Recording() bool { return filmkit.Recording() }
