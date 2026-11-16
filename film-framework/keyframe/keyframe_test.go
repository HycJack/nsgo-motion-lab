package keyframe

import (
	"math"
	"strings"
	"testing"
	"time"
)

// TestSampleClampsAndEases checks the pure half of a take: the same clamp
// outside the keys and the same per-segment easing the platform solves.
func TestSampleClampsAndEases(t *testing.T) {
	take := New("sample", 4*time.Second)
	take.Box("a").On("x", Keys(0, 0, 2, 100, 4, 0).Ease(Linear))

	for _, tc := range []struct {
		at   time.Duration
		want float64
	}{
		{-time.Second, 0},
		{0, 0},
		{time.Second, 50},
		{2 * time.Second, 100},
		{3 * time.Second, 50},
		{9 * time.Second, 0},
	} {
		if got := take.At("a.x", tc.at); got != tc.want {
			t.Errorf("x at %v = %.1f, want %.1f", tc.at, got, tc.want)
		}
	}
}

// TestWinShape checks the window helper: hold at 1 between the ramps, 0 outside.
func TestWinShape(t *testing.T) {
	k := Win(time.Second, 3*time.Second, 500*time.Millisecond)
	keys := k.Keys()
	if len(keys) != 4 {
		t.Fatalf("window has %d keys, want 4", len(keys))
	}
	want := []Key{{1, 0}, {1.5, 1}, {2.5, 1}, {3, 0}}
	for i, w := range want {
		if keys[i] != w {
			t.Errorf("key %d = %+v, want %+v", i, keys[i], w)
		}
	}
}

// TestDedupeKeepsFirst checks that a repeated time keeps the first value: an
// interpolation needs strictly increasing input, and a hard cut written as two
// keys at once has to survive as one.
func TestDedupeKeepsFirst(t *testing.T) {
	got, dropped := dedupe([]Key{{0, 1}, {2, 1}, {2, 0}, {3, 0}})
	if len(got) != 3 || got[1].V != 1 {
		t.Fatalf("dedupe = %+v, want the 2s key kept at 1", got)
	}
	if dropped != 1 {
		t.Errorf("dedupe dropped %d keys, want 1", dropped)
	}
	// And Lint has to say it out loud: a key that vanished is invisible in the
	// timeline, so it would otherwise be found only in the render.
	take := New("cut", 4*time.Second)
	take.Box("a").On("x", Keys(0, 0, 2, 100, 2, 0, 4, 100))
	if said := issues(take); !strings.Contains(said, "repeated a time") {
		t.Errorf("Lint stayed silent about a dropped key; it said:\n%s", said)
	}
}

// TestLintCatchesMistakes walks the mistakes this layer exists to catch: a
// channel the kind does not have, a path that names nothing, keys outside the
// take, and two tracks fighting over one channel.
func TestLintCatchesMistakes(t *testing.T) {
	take := New("lint", 3*time.Second)
	take.Box("a").
		On("w", Keys(0, 10, 2, 20)).     // Box has no width channel
		On("opacity", Keys(0, 0, 9, 1)). // 9s is past the take
		On("opacity", Keys(0, 0, 2, 1))  // and a second track on it
	take.Track("nothing.opacity", Keys(0, 0)) // no node named "nothing"
	take.Box("a").On("x", Keys(0, 0, 3, 10))  // two nodes share a name

	issues := take.Lint()
	if len(issues) == 0 {
		t.Fatal("Lint found nothing in a take with four mistakes")
	}
	want := []string{
		"has no channel",
		"outside the take",
		"only the last one renders",
		"no node has this name",
		"share this name",
	}
	all := make([]string, len(issues))
	for i, is := range issues {
		all[i] = is.String()
	}
	joined := strings.Join(all, "\n")
	for _, w := range want {
		if !strings.Contains(joined, w) {
			t.Errorf("Lint missed %q; it said:\n%s", w, joined)
		}
	}
}

// TestLintCleanTake checks the other direction: a well-formed take reports
// nothing but what it should — here, a dead tail, because nothing moves in the
// second half.
func TestLintCleanTake(t *testing.T) {
	take := New("clean", 4*time.Second)
	take.Box("a").On("opacity", Win(0, 4*time.Second, 400*time.Millisecond))
	take.Box("b").At(10, 0).Size(20, 20).
		On("x", Keys(1, 10, 3.5, 120).Ease(EaseInOut))
	for _, is := range take.Lint() {
		if is.Fatal {
			t.Errorf("clean take reported an error: %s", is)
		}
	}
}

// TestDeadSpan pins the gap detector: a take whose only track stops moving a
// quarter of the way in has to report the rest as dead.
func TestDeadSpan(t *testing.T) {
	take := New("dead", 8*time.Second)
	take.Box("a").On("x", Keys(0, 0, 1, 100))
	spans := take.deadSpans()
	if len(spans) == 0 {
		t.Fatal("no dead span found after a track ends at 1s of 8s")
	}
	if spans[0][0] < 0.9 || spans[0][1] < 7 {
		t.Errorf("dead span = %v, want roughly 1s..8s", spans[0])
	}
}

// TestTimelineHasARowPerTrack checks the text view of a take, which is what an
// author reads instead of the screen.
func TestTimelineHasARowPerTrack(t *testing.T) {
	take := New("rows", 2*time.Second)
	take.Box("a").On("opacity", Win(0, time.Second, 200*time.Millisecond))
	take.Text("t", "hi").On("y", Keys(0, 0, 2, 40))
	out := take.Timeline()
	for _, want := range []string{"a.opacity", "t.y", "2.0s", "2 tracks"} {
		if !strings.Contains(out, want) {
			t.Errorf("timeline is missing %q:\n%s", want, out)
		}
	}
	t.Logf("\n%s", take.Dump())
}

// TestLintCatchesUndrawable walks the nodes that are in the tree but can never
// draw: a camera with no program, a canvas with no Draw, a disc that is a
// rectangle because it has no width, and a group moved by At (which build
// pins at 0,0 — a group moves with a track).
func TestLintCatchesUndrawable(t *testing.T) {
	take := New("undrawable", 2*time.Second)
	take.Add(Camera, "cam")
	take.Add(Canvas, "sheet")
	take.Add(Disc, "moon").At(10, 10) // no Size: fills the stage as a slab
	take.Group("plate").At(20, 20)    // base position is ignored on a group

	joined := issues(take)
	for _, want := range []string{
		"camera has no Program",
		"canvas has no Draw",
		"disc has no width",
		"group ignores At",
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("Lint missed %q; it said:\n%s", want, joined)
		}
	}
}

// TestTimelineDrawsTheTrackEasing checks that the sparkline is the curve the
// renderer will play, not time itself: the same keys eased in and out have to
// look different from a straight ramp, or the timeline lies about the take.
func TestTimelineDrawsTheTrackEasing(t *testing.T) {
	straight := New("straight", 4*time.Second)
	straight.Box("a").On("x", Keys(0, 0, 4, 100).Ease(Linear))
	curved := New("curved", 4*time.Second)
	curved.Box("a").On("x", Keys(0, 0, 4, 100).Ease(EaseInOut))

	if straight.Timeline() == curved.Timeline() {
		t.Errorf("both easings draw one row:\n%s", straight.Timeline())
	}
}

// TestCameraChannelsWithDots checks the one place where a channel name contains
// the separator: "cam.eye.x" is node "cam", channel "eye.x". Splitting on the
// last dot made every camera track a fatal "no channel x" — on the kind that
// needs channels with dots the most.
func TestCameraChannelsWithDots(t *testing.T) {
	take := New("cam", 2*time.Second)
	cam := take.Add(Camera, "cam")
	cam.On("eye.x", Keys(0, -2, 2, 2))
	cam.On("look.z", Keys(0, 0, 2, 1))

	if said := issues(take); strings.Contains(said, "no channel") {
		t.Fatalf("a camera's own channels were rejected:\n%s", said)
	}

	// Track takes a path, and a path with a dotted channel must find its node.
	byPath := New("path", 2*time.Second)
	byPath.Add(Camera, "cam2")
	byPath.Track("cam2.eye.z", Keys(0, 1, 2, 3))
	if got := byPath.At("cam2.eye.z", 2*time.Second); math.Abs(got-3) > 1e-9 {
		t.Errorf("Track(cam2.eye.z) evaluated to %v, want 3 — it did not bind to the camera", got)
	}
	if said := issues(byPath); strings.Contains(said, "no channel") || strings.Contains(said, "no node") {
		t.Errorf("Track left something to complain about:\n%s", said)
	}
}

// TestTimelineDrawsAFlatTrack checks that a track that never changes still
// shows up as a row: a row of spaces reads as "no track here", which is the
// opposite of a channel that holds one value for the whole take.
func TestTimelineDrawsAFlatTrack(t *testing.T) {
	take := New("flat", 2*time.Second)
	take.Box("a").On("opacity", Keys(0, 0.5, 2, 0.5))
	row := ""
	for _, line := range strings.Split(take.Timeline(), "\n") {
		if strings.Contains(line, "a.opacity") {
			row = line
		}
	}
	if !strings.Contains(row, "----") {
		t.Errorf("a flat track draws no curve: %q", row)
	}
}

// issues is Lint as one string, for the "does it say this" tests.
func issues(take *Take) string {
	lines := make([]string, 0, 8)
	for _, is := range take.Lint() {
		lines = append(lines, is.String())
	}
	return strings.Join(lines, "\n")
}

// TestSeekWithoutABuildIsSafe checks the exporter's failure mode: a take that
// was never mounted has no value to set, and seeking it must not panic.
func TestSeekWithoutABuildIsSafe(t *testing.T) {
	New("cold", time.Second).Seek(time.Second)
}
