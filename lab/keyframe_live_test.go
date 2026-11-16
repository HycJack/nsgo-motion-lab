package main

import (
	"math"
	"testing"
	"time"

	"github.com/oneo-me/nsgo/nsgotest"

	"nsgofilm/keyframe"
)

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
// page commits several times over a session and a take that replayed on every
// commit would make a seek impossible.
func TestTakeNodeIsStable(t *testing.T) {
	app := nsgotest.Start(t, Root, nsgotest.Size(1180, 860))
	app.Press("id:tab-Studio")

	shoot := keyframe.TakeNamed("shoot")
	if shoot == nil {
		t.Fatal("the camera take did not register")
	}
	shoot.Seek(9 * time.Second)
	app.Settle()

	// Switching tabs away and back rebuilds the page, and the take with it.
	app.Press("id:tab-Canvas")
	app.Press("id:tab-Studio")
	app.Settle()

	app.Query("id:take-shoot-cam")
	if got := float64(app.Query("id:take-shoot-cam")[0].Animated.TranslateX); math.Abs(got) > 0.01 {
		t.Errorf("camera moved to %.2f after a rebuild, want its base position", got)
	}
}
