package film

import (
	"testing"
	"time"

	"github.com/oneo-me/nsgo/core"
	"github.com/oneo-me/nsgo/nsgotest"

	"nsgofilm/filmkit"
)

// TestFilmPublishesItsPlayhead checks the live half of "one clock": a mounted
// film publishes the time it draws with, so a take inside it can follow the
// film's playhead instead of running its own animation.
//
// The second half is the trap. The frame loop advances time with filmkit.Wall,
// not Now: Now answers from the clock the film published, so a loop that read
// Now would be reading its own output one frame late and a silent film would
// stand still. That is exactly what this asserts — the playhead advances on
// real time, unmounting releases the clock, and the next film can take over.
func TestFilmPublishesItsPlayhead(t *testing.T) {
	if filmkit.Driven() {
		t.Fatal("a clock was already published before any film mounted")
	}

	f := Film{
		W: 320, H: 180,
		Duration: 4,
		Shots:    []Shot{{D: 4, Fn: func(*core.Drawing, *Frame) {}}},
	}

	t.Run("playing", func(t *testing.T) {
		app := nsgotest.Start(t, New(f), nsgotest.Size(360, 200))
		app.Settle()

		if !filmkit.Driven() {
			t.Fatal("a mounted film published no clock; takes inside it would play on their own")
		}
		if now := filmkit.Now(); now > 0.05 {
			t.Errorf("a fresh film is at %.3fs, want the head of the film", now)
		}

		// Real time on purpose: a silent film's playhead is the wall clock, and a
		// virtual advance would not prove the loop is reading it.
		time.Sleep(150 * time.Millisecond)
		app.Advance(20 * time.Millisecond)
		app.Settle()

		if now := filmkit.Now(); now < 0.1 {
			t.Errorf("the playhead is at %.3fs after 150ms of playing — the frame loop is reading the clock it publishes", now)
		}
	})

	// The subtest's cleanup has disposed the film by now.
	if filmkit.Driven() {
		t.Error("a disposed film left its clock published; the next film could not take over")
	}
}
