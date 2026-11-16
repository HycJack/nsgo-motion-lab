package keyframe

import (
	"sync"
	"time"

	"github.com/oneo-me/nsgo/core"

	"nsgofilm/filmkit"
)

// Tie puts every live take under the film clock. One animation frame reads
// filmkit.Now() and seeks the takes when it moved.
//
// This is what lets the one recorder drive both halves of a film: a canvas shot
// reads the clock while it draws, and a node take is seeked by it — the recorder
// only has to move time (`window.__filmT` on the web target), and both land on
// the same frame. Nothing in a take reads the wall clock, so a frame is a
// function of the clock and the recording stays reproducible.
//
// Without a recorder there is no external clock. The takes then keep playing on
// their own (Play), which is what the live page wants; Tie only overrides that
// while something outside is driving time.
func Tie(cx *core.Scope) {
	tieMu.Lock()
	// cx is the app's own scope here (the one Root is called with), which is
	// also where a take anchors its clock so a page change cannot release it.
	anchor = cx
	if tied == cx {
		tieMu.Unlock()
		return
	}
	tied = cx
	tieGen++
	gen := tieGen
	tieMu.Unlock()

	last := -1.0
	var loop func(time.Time)
	loop = func(time.Time) {
		// A new scope means a new mount: this loop's scope is gone, so it
		// stops instead of running beside the one that replaced it (two
		// loops seeking the same takes would double the work per frame).
		tieMu.Lock()
		stale := gen != tieGen
		tieMu.Unlock()
		if stale {
			return
		}
		now := filmkit.Now()
		switch {
		case !Recording():
			// Nobody is driving: the take owns its clock.
			last = -1
		case now != last:
			last = now
			SeekAll(seconds(now))
		}
		cx.RequestAnimationFrame(loop)
	}
	cx.RequestAnimationFrame(loop)
}

var (
	tieMu  sync.Mutex
	tied   *core.Scope
	tieGen int
	// anchor is the app's own scope, the one Tie was called with: a take hangs
	// its clock there so unmounting a page does not release it.
	anchor *core.Scope
)

// anchorScope is the app's own scope, or nil before Tie has run.
func anchorScope() *core.Scope {
	tieMu.Lock()
	defer tieMu.Unlock()
	return anchor
}

// Progress reports the film clock each take was last seeked to, which is how a
// recorder proves the seek landed instead of trusting a picture.
func Progress() map[string]float64 {
	out := map[string]float64{}
	for _, t := range Takes() {
		out[t.Name()] = t.Progress() * t.Duration().Seconds()
	}
	return out
}
