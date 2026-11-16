package main

import (
	"strconv"
	"time"

	"github.com/oneo-me/nsgo/core"

	"nsgofilm/keyframe"
)

// This page is the direction layer in use: three takes written with
// keyframe.Take — the camera pass from the One take page, a three-shot storyboard
// cut on one clock, and the script check that reads a take back as text.
// Everything here animates with zero commits per frame.
func pageStudio(cx *core.Scope) core.Node {
	shoot, story, check := studioTakes()
	return grid(takeCard(cx, shoot, "S1 · Take", "eight camera tracks, four windows, one clock — the pass from T1, addressed by name", colAmber, 240),
		takeCard(cx, story, "S2 · Storyboard", "a cut is a 20 ms ramp, a dissolve is 400 ms — a handoff, not a mount", colViolet, 200),
		takeCard(cx, check, "S3 · Script check", "Take.Timeline and Take.Lint, drawn — a take as something to read, not only to watch", colCyan, 210))
}

// studioShoot is the one-take camera pass written as tracks: the same footage
// as T1, and the numbers now read as a clock instead of as progress.
func studioShoot() *keyframe.Take {
	take := keyframe.New("shoot", 18*time.Second)

	cam := take.Camera("cam", rayProgram).Label("ray-marched take")
	cam.On("eye.x", keyframe.Keys(0, -2.60, 2.5, -2.15, 5.4, -1.55, 8.1, -1.05, 10.8, -0.45, 13.5, 0.40, 15.8, 1.55, 18, 1.70))
	cam.On("eye.y", keyframe.Keys(0, 1.85, 2.5, 1.35, 5.4, 0.95, 8.1, 0.90, 10.8, 0.82, 13.5, 0.86, 15.8, 0.95, 18, 1.00))
	cam.On("eye.z", keyframe.Keys(0, 4.30, 2.5, 3.90, 5.4, 3.35, 8.1, 2.95, 10.8, 2.45, 13.5, 2.30, 15.8, 2.10, 18, 2.55))
	cam.On("look.x", keyframe.Keys(0, -0.05, 5.4, 0.02, 10.8, 0.12, 15.8, 0.15, 18, 0.05))
	cam.On("look.y", keyframe.Keys(0, 0.55, 8.1, 0.50, 13.5, 0.46, 18, 0.46))
	cam.On("look.z", keyframe.Keys(0, 0.10, 5.4, -0.02, 10.8, -0.05, 13.5, -0.10, 18, 0.18))
	cam.On("fov", keyframe.Keys(0, 0.72, 5.4, 0.66, 10.8, 0.78, 13.5, 0.70, 15.8, 0.66, 18, 0.60))
	cam.On("tint", keyframe.Keys(0, 0.05, 5.4, 0.30, 9.9, 0.70, 12.6, 0.90, 15.3, 0.35, 18, 0.12))

	take.Text("slate", "ONE TAKE").Font(20).Bold().Color(colText).
		Pad(7, 15).Radius(10).Fill(core.RGBA(0x0b1119, 0.72)).At(20, 18).
		On("opacity", keyframe.Win(500*ms, 5*time.Second, 400*ms))
	take.Text("line", "a whole shot on one clock").Font(12).Color(colMuted).
		At(24, 78).
		On("opacity", keyframe.Win(2900*ms, 8600*ms, 500*ms))
	// The lower third: a plate that names the take, then the line that replaces
	// it. Same slot, so the handover reads as one caption that changes.
	take.Text("panel", "TAKE 01 · 8 TRACKS").Font(12).Color(colAmber).
		Pad(8, 12).Radius(9).Fill(core.RGBA(0x0d131c, 0.88)).At(24, 172).
		On("opacity", keyframe.Win(8300*ms, 14800*ms, 450*ms)).
		On("x", keyframe.KeysAt(8300*ms, 24, 14800*ms, 46))
	take.Text("landing", "every beat is a track").Font(13).Bold().Color(colText).
		Pad(7, 12).Radius(9).Fill(core.RGBA(0x060a10, 0.74)).At(24, 172).
		On("opacity", keyframe.Win(15300*ms, 18*time.Second, 500*ms)).
		On("y", keyframe.KeysAt(15300*ms, 192, 18*time.Second, 176))

	return take
}

// studioStory is a three-shot storyboard: the shots are windows on one clock,
// so a cut is a 20 ms ramp and a dissolve is a 400 ms one. Nothing is mounted
// or removed; the shots hand the frame over.
func studioStory() *keyframe.Take {
	take := keyframe.New("story", 12*time.Second)

	wide := take.Group("wide")
	// Out over 400 ms, while the next shot comes up underneath: a dissolve is
	// an overlap, so nothing ever reaches an empty frame between two shots.
	wide.On("opacity", keyframe.Keys(0, 1, 3.8, 1, 4.2, 0).Ease(keyframe.EaseInOut))
	// Five blocks standing on the ground line at y=176: they slide up into
	// place, one after the other, and the shot ends on the finished row.
	for i := 0; i < 5; i++ {
		x := 40 + float64(i)*84
		h := 40 + float64(i%3)*18
		start := 0.3 + float64(i)*0.28
		wide.Add(keyframe.Box, "city"+strconv.Itoa(i+1)).At(x, 176).Size(30, h).
			Fill(core.RGB(0x24405c)).Radius(4).
			On("y", keyframe.Keys(0, 176, start, 176-h).Ease(keyframe.EaseInOut)).
			On("opacity", keyframe.Keys(0, 0, start, 1))
	}

	orbit := take.Group("orbit")
	// The second shot holds at full opacity until the frame the third arrives on,
	// then cuts out over 20 ms — the same 20 ms the third cut snaps in with. That
	// overlap is the point: two ramps that cross at one boundary sum to full
	// brightness, so there is no frame in between that shows neither shot. Give
	// the fade-out its own window ending before 8.0s and the stage goes empty for
	// a frame or two, which is what "dissolve on one clock" is supposed to prevent.
	orbit.On("opacity", keyframe.KeysAt(
		0*ms, 0, 3700*ms, 0, 4100*ms, 1, 7980*ms, 1, 8*time.Second, 0,
	).Ease(keyframe.Linear))
	orbit.Add(keyframe.Disc, "moon").At(120, 60).Size(54, 54).Fill(core.RGB(0x2f5f86)).
		On("x", keyframe.KeysAt(3700*ms, 120, 7*time.Second, 300)).
		On("scale", keyframe.KeysAt(3700*ms, 0.6, 5200*ms, 1.1, 7*time.Second, 0.9))
	orbit.Add(keyframe.Box, "rail").At(40, 176).Size(380, 2).Fill(core.RGBA(0x3a5f80, 0.7)).
		On("scaleX", keyframe.KeysAt(3700*ms, 0.2, 7*time.Second, 1))

	// The third shot snaps in on a 20 ms ramp instead of dissolving: the one
	// cut on this clock, and it happens while the second shot is still fading.
	land := take.Group("land")
	land.On("opacity", keyframe.KeysAt(7980*ms, 0, 8*time.Second, 1).Ease(keyframe.Linear))
	land.Add(keyframe.Text, "title").Text("3 SHOTS · 1 CLOCK").Font(17).Bold().Color(colText).
		At(40, 84).On("y", keyframe.KeysAt(8*time.Second, 96, 9*time.Second, 84, 12*time.Second, 80))

	return take
}

// studioCheck is a take read back as text: the schedule as a timeline and what
// Lint makes of it. This is the loop an author who cannot see the screen works
// in — write the take, print it, read the numbers.
func studioCheck() *keyframe.Take {
	take := keyframe.New("check", 4*time.Second)
	take.Box("plate").At(30, 30).Size(80, 40).Radius(8).Fill(core.RGBA(0x2a6b8f, 0.5)).
		On("opacity", keyframe.Win(0, 2*time.Second, 300*ms)).
		On("x", keyframe.Keys(0, 30, 4, 210))
	take.Text("cap", "checked").Font(13).Color(colText).At(30, 120).
		On("opacity", keyframe.FadeIn(2*time.Second, 300*ms))

	lines := scriptLines(take)
	take.Canvas("script").Label("timeline and lint").Draw(func(d *core.Drawing, size core.Size) {
		d.SetFillColor(core.RGB(0x0b1119))
		d.FillRect(0, 0, float32(size.Width), float32(size.Height))
		y := float32(12)
		for _, ln := range lines {
			col := colMuted
			if ln.hot {
				col = colCyan
			}
			d.SetFillColor(col)
			d.FillText(ln.text, 12, y, core.Font{Size: 8.2})
			y += 11.5
		}
	})

	return take
}

// scriptLine is one line of the check card.
type scriptLine struct {
	text string
	hot  bool
}

// scriptLines renders a take's schedule and its lint as lines.
func scriptLines(take *keyframe.Take) []scriptLine {
	out := []scriptLine{{"take " + take.Name() + " · " + take.Duration().String(), true}}
	for _, ln := range splitLines(take.Timeline()) {
		out = append(out, scriptLine{ln, false})
	}
	issues := take.Lint()
	if len(issues) == 0 {
		out = append(out, scriptLine{"lint: clean", true})
	}
	for _, is := range issues {
		out = append(out, scriptLine{"lint: " + is.String(), true})
	}
	return out
}

// splitLines breaks text into lines, dropping the empty tail a trailing newline
// leaves.
func splitLines(s string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			if i > start {
				out = append(out, s[start:i])
			}
			start = i + 1
		}
	}
	if start < len(s) {
		out = append(out, s[start:])
	}
	return out
}

// takeCard is the frame every studio take sits in: the stage, the name, what to
// watch, and a Replay chip bound to the take's own clock.
func takeCard(cx *core.Scope, take *keyframe.Take, title, note string, accent core.Color, height float32) core.Node {
	// The chip carries the card's number, the way the catalog tiles do, so
	// each chip on the page is addressable on its own.
	id := title
	for i := 0; i < len(title); i++ {
		if title[i] == ' ' {
			id = title[:i]
			break
		}
	}
	return card(title, note, accent,
		// The stage is named so an export can frame exactly the take and not
		// the card around it: [data-testid="stage-S1"] is a clip rectangle.
		core.View{TestID: "stage-" + id, Style: well(height), Children: []core.Node{
			take.Frame(cx),
			core.View{Style: core.Style{Position: core.Absolute, Right: core.Pt(10), Top: core.Pt(10)},
				Children: []core.Node{chip("Replay "+id, accent, take.Play)}},
		}},
	)
}
