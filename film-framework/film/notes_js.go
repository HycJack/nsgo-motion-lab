//go:build js && wasm

package film

import (
	"encoding/json"
	"syscall/js"
)

const notesKey = "nsgo-film-notes"

// loadMarks 从 localStorage 读标记（web 端持久化）。
func loadMarks() []float64 {
	v := js.Global().Get("localStorage").Call("getItem", notesKey)
	if v.Type() != js.TypeString {
		return nil
	}
	var ts []float64
	if json.Unmarshal([]byte(v.String()), &ts) != nil {
		return nil
	}
	return ts
}

// saveMarks 把标记写回 localStorage。
func saveMarks(ts []float64) {
	b, _ := json.Marshal(ts)
	js.Global().Get("localStorage").Call("setItem", notesKey, string(b))
}
