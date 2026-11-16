package main

import (
	"sync"

	"nsgofilm/keyframe"
)

// The Studio page's takes are long-lived objects, not values a page builds: a
// take owns a clock, and a clock has to survive a commit for an export to be
// able to set it. This is the one place the layer differs from the catalog
// tiles, whose effects are rebuilt on every commit and simply replay.
var (
	takesOnce sync.Once
	shootTake *keyframe.Take
	storyTake *keyframe.Take
	checkTake *keyframe.Take
)

// studioTakes builds the page's three takes once and hands them out.
func studioTakes() (*keyframe.Take, *keyframe.Take, *keyframe.Take) {
	takesOnce.Do(func() {
		shootTake = studioShoot()
		storyTake = studioStory()
		checkTake = studioCheck()
	})
	return shootTake, storyTake, checkTake
}
