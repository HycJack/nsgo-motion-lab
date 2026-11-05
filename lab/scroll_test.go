package main

import (
	"testing"

	"github.com/oneo-me/nsgo/nsgotest"
)

// TestScrollFitsWindow holds the scroll view inside the window. When the
// column around it forgot FlexShrink, it grew to its content (894pt in an
// 860pt window), the scroll view inherited that height, and nothing ever
// overflowed — so no bar and no scrolling, on every target.
func TestScrollFitsWindow(t *testing.T) {
	const width, height = 1180, 860
	app := nsgotest.Start(t, Root, nsgotest.Size(width, height))

	vs := app.Query("id:page-scroll")
	if len(vs) != 1 {
		t.Fatalf("found %d scroll views, want 1", len(vs))
	}
	v := vs[0]
	if h := v.Frame.Size.Height; h >= height {
		t.Errorf("scroll view is %vpt tall in a %vpt window: it grew past the window, so it can never scroll", h, height)
	}
	// The title bar keeps its room above it.
	if y := v.Frame.Origin.Y; y <= 0 {
		t.Errorf("scroll view starts at %v, on top of the title bar", y)
	}
	t.Logf("scroll view: %v at y=%v in a %vpt window", v.Frame.Size.Height, v.Frame.Origin.Y, height)
}
