package film

import (
	"math"
	"sync"
	"testing"
	"time"

	"github.com/nsgo-dev/nsgo/core"
	"github.com/nsgo-dev/nsgo/nsgotest"

	"nsgofilm/filmkit"
	"nsgofilm/keyframe"
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

// TestLayerFollowsTheFilmClock mounts a take the way a film does and reads the
// rendered node: the layer is in the tree, and it sits where the take's own
// sample says at the film's playhead — the whole point of Over.
func TestLayerFollowsTheFilmClock(t *testing.T) {
	beat := keyframe.New("beat", 4*time.Second)
	beat.Box("moon").Size(20, 20).At(0, 90).Fill(core.RGB(0xffffff)).
		On("x", keyframe.Keys(0, 20, 4, 300))

	f := Film{
		W: 320, H: 180,
		Duration: 4,
		Shots:    []Shot{{D: 4, Fn: func(*core.Drawing, *Frame) {}}},
		Layers:   []Layer{Over(beat, 0, 4)},
	}
	app := nsgotest.Start(t, New(f), nsgotest.Size(320, 180))
	// A take schedules its one mount-time seek for the next frame. Let it land
	// now: if it lands after the sleep below it would look like following.
	app.Advance(20 * time.Millisecond)
	app.Settle()

	// A layer is in the tree only while the film is inside its window, and the
	// whole film is this one's window.
	if !app.Exists("id:take-beat-moon") {
		t.Fatal("the layer's take is not in the tree; Over mounted nothing")
	}

	// Let the film run on real time, then read the picture again. A take seeks
	// once when it mounts; only the follow loop moves it afterwards, so this is
	// what separates a layer that follows the playhead from one frozen on it.
	time.Sleep(220 * time.Millisecond)
	app.Advance(20 * time.Millisecond)
	app.Settle()

	// Same clock, same frame: what the take samples at the film's playhead is
	// what it renders.
	now := filmkit.Now()
	want := beat.Sample(time.Duration(now * float64(time.Second)))["moon.x"]
	if now < 0.15 {
		t.Fatalf("the playhead is at %.3fs; the test proved nothing", now)
	}
	got := float64(app.Query("id:take-beat-moon")[0].Animated.TranslateX)
	if math.Abs(got-want) > 0.5 {
		t.Errorf("the moon renders at x=%.1f at %.3fs, want %.2f from the track", got, now, want)
	}
}

// TestDelayedTakeShowsItsOwnTime pins Take.Delay on both paths that seek a
// take: the one-shot seek at mount, and the follow loop that moves it when the
// clock advances. A take delayed by 2s shows its own 1s when the clock reads
// 3s and its own 2s at 4s; without the shift the follow loop would clamp it to
// the last keyframe and the picture would jump 200 points.
//
// The clock is published by hand instead of mounting a film, so the test does
// not have to wait four seconds for the picture to mean anything.
func TestDelayedTakeShowsItsOwnTime(t *testing.T) {
	beat := keyframe.New("late", 4*time.Second)
	beat.Box("dot").Size(20, 20).At(0, 90).Fill(core.RGB(0xffffff)).
		On("x", keyframe.Keys(0, 0, 4, 400))
	beat.Delay(2 * time.Second)

	var mu sync.Mutex
	playhead := 3.0
	stop := filmkit.Drive(func() float64 {
		mu.Lock()
		defer mu.Unlock()
		return playhead
	})
	defer stop()

	app := nsgotest.Start(t, func(cx *core.Scope) core.Node {
		keyframe.Tie(cx)
		return beat.Frame(cx)
	}, nsgotest.Size(320, 180))
	frame := func() {
		app.Advance(20 * time.Millisecond)
		app.Settle()
	}
	frame()

	nodes := app.Query("id:take-late-dot")
	if len(nodes) == 0 {
		t.Fatal("the take is not in the tree")
	}
	dot := func() float64 { return float64(app.Query("id:take-late-dot")[0].Animated.TranslateX) }
	check := func(at time.Duration, what string) {
		t.Helper()
		want := beat.Sample(at)["dot.x"]
		if want <= 0 || want >= 400 {
			t.Fatalf("the fixture is wrong: dot.x at %v is %v, not mid-ramp", at, want)
		}
		if got := dot(); math.Abs(got-want) > 0.5 {
			t.Errorf("%s: the dot renders at x=%.1f, want %.2f (the take's own %v)", what, got, want, at)
		}
	}

	// Mount seek: clock 3s minus the 2s delay is the take's 1s.
	check(1*time.Second, "at mount")

	// Follow loop: the clock moves one second, so the take must move one second.
	mu.Lock()
	playhead = 4
	mu.Unlock()
	frame()
	check(2*time.Second, "after the clock advanced")
}

// TestLayerWindowHidesIt checks the predicate that keeps a layer unmounted
// outside its window — the reason a take authored for one beat does not sit
// frozen through the rest of the film.
func TestLayerWindowHidesIt(t *testing.T) {
	whole := Layer{}
	beat := Layer{Start: 2, End: 5}
	tail := Layer{Start: 2}

	for _, c := range []struct {
		l    Layer
		t    float64
		want bool
	}{
		{whole, 0, true}, {whole, 9, true},
		{beat, 1.99, false}, {beat, 2, true}, {beat, 5, true}, {beat, 5.01, false},
		{tail, 1, false}, {tail, 6, true},
	} {
		if got := c.l.shows(c.t); got != c.want {
			t.Errorf("layer[%v,%v] at %.2fs shows %v, want %v", c.l.Start, c.l.End, c.t, got, c.want)
		}
	}
}
