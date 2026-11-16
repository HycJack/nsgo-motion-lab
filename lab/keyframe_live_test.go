package main

import (
	"math"
	"strings"
	"testing"
	"time"

	"github.com/oneo-me/nsgo/nsgotest"

	"nsgofilm/keyframe"
)

// TestStudioTakesAreClean checks the three demo takes against their own lint.
// The S3 card prints "lint: clean", so the page is lying if it does not hold —
// and a demo shot is the one place where a hidden mistake teaches the wrong
// thing.
func TestStudioTakesAreClean(t *testing.T) {
	shoot, story, check := studioTakes()
	for _, take := range []*keyframe.Take{shoot, story, check} {
		issues := take.Lint()
		if len(issues) == 0 {
			continue
		}
		var b strings.Builder
		for _, is := range issues {
			b.WriteString("\n  " + is.String())
		}
		t.Errorf("%s lints dirty:%s", take.Name(), b.String())
	}
}

// TestTakePlaysWhenNoOneDrives checks the other half of the clock: with no
// recorder around the take owns its time, so it plays itself instead of sitting
// on frame zero. It starts that play from inside the mount, which is where the
// engine's frame loop can be re-entered, so this pins that the shot moves.
//
// nsgotest settles an animation in one go, so what is checked is the end state:
// a take that never started still shows frame zero (y 176), and one that played
// shows its last key (y 136).
func TestTakePlaysWhenNoOneDrives(t *testing.T) {
	app := nsgotest.Start(t, Root, nsgotest.Size(1180, 860))
	app.Press("id:tab-Studio")

	story := keyframe.TakeNamed("story")
	if story == nil {
		t.Fatal("the storyboard take did not register")
	}
	app.Settle()

	rendered := float64(app.Query("id:take-story-city1")[0].Animated.TranslateY)
	if want := story.At("city1.y", 2*time.Second); math.Abs(rendered-want) > 0.01 {
		t.Errorf("city1 renders %.2f, want its last key %.2f — the take never played (frame zero is %.2f)",
			rendered, want, story.At("city1.y", 0))
	}
	if story.Progress() != 0 {
		t.Errorf("nothing seeked this take, yet Progress says %.2f", story.Progress())
	}
}

// TestTakeSeekMatchesSample drives a mounted take through the platform and
// checks the rendered values against Take.Sample. The whole export path rests
// on these two agreeing: a frozen frame is only a frame of the take if the
// clock and the pure evaluation are the same function of the same value.
func TestTakeSeekMatchesSample(t *testing.T) {
	app := nsgotest.Start(t, Root, nsgotest.Size(1180, 860))
	app.Press("id:tab-Studio")
	app.Find(`text:"S2 · Storyboard"`)

	story := keyframe.TakeNamed("story")
	if story == nil {
		t.Fatal("the storyboard take did not register")
	}

	for _, at := range []time.Duration{0, 6 * time.Second, 9 * time.Second, 12 * time.Second} {
		story.Seek(at)
		app.Settle()

		sample := story.Sample(at)
		for _, tc := range []struct {
			selector string
			path     string
			got      float64
		}{
			{"id:take-story-moon", "moon.x", float64(app.Query("id:take-story-moon")[0].Animated.TranslateX)},
			{"id:take-story-rail", "rail.scaleX", float64(app.Query("id:take-story-rail")[0].Animated.ScaleX)},
			{"id:take-story-wide", "wide.opacity", float64(app.Query("id:take-story-wide")[0].Animated.Opacity)},
		} {
			if math.Abs(tc.got-sample[tc.path]) > 0.01 {
				t.Errorf("at %v, %s renders %.3f, Sample says %.3f", at, tc.path, tc.got, sample[tc.path])
			}
		}
	}
}

// TestTakeNodeIsStable checks that rebuilding a take does not restart it: the
// page commits several times over a session, and a take that replayed on every
// commit would make a seek impossible. The proof is the rendered value — after
// a tab round-trip the moon still sits where the seek put it, not where frame
// zero would.
func TestTakeNodeIsStable(t *testing.T) {
	app := nsgotest.Start(t, Root, nsgotest.Size(1180, 860))
	app.Press("id:tab-Studio")

	story := keyframe.TakeNamed("story")
	if story == nil {
		t.Fatal("the storyboard take did not register")
	}
	at := 9 * time.Second
	story.Seek(at)
	app.Settle()

	want := story.Sample(at)["moon.x"]
	before := float64(app.Query("id:take-story-moon")[0].Animated.TranslateX)
	if math.Abs(before-want) > 0.01 {
		t.Fatalf("moon renders %.2f after the seek, Sample says %.2f", before, want)
	}

	// Switching tabs away and back unmounts the page and asks the take for its
	// tree again. A take that restarted there would show frame zero's x.
	app.Press("id:tab-Canvas")
	app.Press("id:tab-Studio")
	app.Settle()

	after := float64(app.Query("id:take-story-moon")[0].Animated.TranslateX)
	if math.Abs(after-want) > 0.01 {
		t.Errorf("after remounting, moon.x renders %.2f, want %.2f (the take restarted)", after, want)
	}
	if got := story.Progress() * story.Duration().Seconds(); math.Abs(got-at.Seconds()) > 0.05 {
		t.Errorf("clock moved to %.2fs after remounting, want %.2fs", got, at.Seconds())
	}
}
