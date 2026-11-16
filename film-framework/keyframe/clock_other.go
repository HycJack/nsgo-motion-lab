//go:build !(js && wasm)

package keyframe

// Recording is false off the web: there is no page to inject a clock into, and
// tests drive takes with Seek directly.
func Recording() bool { return false }
