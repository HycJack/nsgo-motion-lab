package filmkit

import (
	"sync"
	"testing"
)

// TestDriveBecomesTheClock checks the middle rung of Now's precedence: a film
// that publishes its playhead makes Now answer from it, so everything drawn in
// that film reads one clock.
func TestDriveBecomesTheClock(t *testing.T) {
	stop := Drive(func() float64 { return 1.5 })
	defer stop()

	if got := Now(); got != 1.5 {
		t.Errorf("Now() = %v, want the published 1.5", got)
	}
	if !Driven() {
		t.Error("a published clock does not count as driven; a take would play on its own instead of following")
	}
}

// TestStopRestoresTheWallClock checks that a film which goes away takes its
// clock with it: nothing left to follow means the wall clock again.
func TestStopRestoresTheWallClock(t *testing.T) {
	stop := Drive(func() float64 { return 42 })
	stop()

	if Driven() {
		t.Error("a stopped clock still counts as driving")
	}
	if got := Now(); got == 42 {
		t.Errorf("Now() still answers from the stopped clock")
	}
	stop() // stopping twice must not panic or unpublish somebody else's clock
}

// TestTheNewestClockWins checks the replacement rule, which is what makes it
// safe for a stale cleanup to fire late: stopping an already-replaced clock does
// nothing, so the film that mounted on top keeps its time.
func TestTheNewestClockWins(t *testing.T) {
	old := Drive(func() float64 { return 1 })
	newer := Drive(func() float64 { return 2 })
	defer newer()

	if got := Now(); got != 2 {
		t.Fatalf("Now() = %v, want the newest clock's 2", got)
	}
	old() // the film that mounted first goes away afterwards
	if got := Now(); got != 2 {
		t.Errorf("a stale stop silenced the newer clock: Now() = %v", got)
	}
}

// TestWallIgnoresTheClock is the one that keeps a film moving. The frame loop
// advances time with Wall: reading Now there would return the loop's own output
// from the frame before, and a silent film would stand still forever.
func TestWallIgnoresTheClock(t *testing.T) {
	stop := Drive(func() float64 { return 7 })
	defer stop()

	if got := Wall(); got == 7 {
		t.Error("Wall() returned the published clock; the film's frame loop would freeze")
	}
}

// TestConcurrentReadsAndPublishes exercises the ordering the frame loop relies
// on: the loop reads every frame while a control (a scrub, a pause) writes the
// playhead, and a film may publish or stop from another goroutine entirely.
func TestConcurrentReadsAndPublishes(t *testing.T) {
	var mu sync.Mutex
	playhead := 0.0
	stop := Drive(func() float64 {
		mu.Lock()
		defer mu.Unlock()
		return playhead
	})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 200; i++ {
			mu.Lock()
			playhead = float64(i)
			mu.Unlock()
			Now()
			Driven()
		}
	}()
	for i := 0; i < 200; i++ {
		Now()
		Wall()
	}
	stop()
	wg.Wait()
}
