package keyframe

import (
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
	got := dedupe([]Key{{0, 1}, {2, 1}, {2, 0}, {3, 0}})
	if len(got) != 3 || got[1].V != 1 {
		t.Fatalf("dedupe = %+v, want the 2s key kept at 1", got)
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

// TestSeekWithoutABuildIsSafe checks the exporter's failure mode: a take that
// was never mounted has no value to set, and seeking it must not panic.
func TestSeekWithoutABuildIsSafe(t *testing.T) {
	New("cold", time.Second).Seek(time.Second)
}
