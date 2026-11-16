package filmkit

import "sync/atomic"

// Clock is a clock a film publishes for everything drawn inside it.
//
// Without one, a film and the shots in it read the wall clock, and a take inside
// a film runs its own animation — so pausing, scrubbing or stepping the film
// would leave the take playing. With one, `Now` answers from the film's own
// playhead and a take seeks to it instead of playing, which is what makes "one
// clock" true live as well as while recording.
type Clock struct {
	read func() float64
}

// published is the live clock, if a film has published one. A pointer, so
// reading it needs no lock: a film publishes once at mount and stops once at
// unmount, and every frame in between only reads.
var published atomic.Pointer[Clock]

// Drive publishes read as the film clock and returns the func that stops it.
//
// Inside a film, call it once, with a read of the film's own playhead:
//
//	stop := filmkit.Drive(func() float64 { return playhead.Get() })
//	cx.Cleanup(stop)
//
// The newest publication wins, so a film that mounts on top of another takes
// the clock over; stopping one that has already been replaced does nothing,
// which is what keeps a stale cleanup from silencing the film that replaced it.
//
// The recorder's window.__filmT still outranks this: a recording is stepped by
// something outside and must not be second-guessed by the film playing under
// it.
func Drive(read func() float64) func() {
	c := &Clock{read: read}
	published.Store(c)
	return func() { published.CompareAndSwap(c, nil) }
}

// Driven reports whether something outside is driving film time: a recorder's
// window.__filmT on the web target, or a film that has published a clock.
// Shots use it to tell "play yourself" from "follow along".
func Driven() bool { return Recording() || hasClock() }

// hasClock reports whether a film published a clock, without counting a
// recorder: the film's own frame loop asks this when deciding how to advance.
func hasClock() bool { return published.Load() != nil }

// clockNow is the published clock's reading, or false when there is none. It is
// the middle of Now's precedence order, between the recorder and the wall.
func clockNow() (float64, bool) {
	if c := published.Load(); c != nil {
		return c.read(), true
	}
	return 0, false
}
