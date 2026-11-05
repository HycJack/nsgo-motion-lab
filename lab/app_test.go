package main

import (
	"testing"
	"time"

	"github.com/oneo-me/nsgo/core"
	"github.com/oneo-me/nsgo/nsgotest"
)

// TestMain gives the platform the main thread under nsgo test -target.
func TestMain(m *testing.M) { nsgotest.Main(m) }

// TestShadersCompile checks every program with the compiler's own errors,
// which MustCompileShader would raise as a panic at start-up.
func TestShadersCompile(t *testing.T) {
	for _, tc := range []struct {
		name string
		src  string
	}{
		{"plasma", plasmaSource},
		{"pointer", pointerSource},
		{"sdf", sdfSource},
		{"ripple", rippleSource},
		{"ray", raySource},
	} {
		if _, err := core.CompileShader(tc.src); err != nil {
			t.Errorf("%s: %v", tc.name, err)
		}
	}
}

// TestTabsSwitchPages walks the three tabs and checks the page each one
// builds.
func TestTabsSwitchPages(t *testing.T) {
	app := nsgotest.Start(t, Root, nsgotest.Size(1180, 860))

	app.Find("id:canvas-wave")
	app.Find("id:canvas-gauge")

	app.Press("id:tab-Animation")
	app.Find("id:preset-fadeIn")
	app.Find("id:preset-breathe")
	app.Press("id:preset-zoomIn")
	app.Press("id:chip-Replay curves")

	app.Press("id:tab-UI Kit")
	app.Find("id:toggle")
	app.Find("id:hamburger")

	app.Press("id:tab-Shader")
	app.Find("id:shader-plasma")
	app.Find("id:shader-pointer")
	app.Find("id:shader-sdf")
	app.Find("id:shader-ripple")

	app.Press("id:tab-3D")
	app.Find(`text:"3D 01 · Ray-marched scene"`)
	app.Find(`text:"3D 09 · CSS 3D"`)
	app.Press("id:chip-Replay 3D 08")

	app.Press("id:tab-Canvas")
	app.Find("id:canvas-bars")
}

// TestUIKitControls drives the ui kit page: every control it offers has
// to answer, and the state it changes has to show.
func TestUIKitControls(t *testing.T) {
	app := nsgotest.Start(t, Root, nsgotest.Size(1180, 860))
	app.Press("id:tab-UI Kit")

	// The toast is mounted only while it shows, so it appears on press
	// and leaves on its own two seconds later.
	app.Press("id:chip-Show a toast")
	app.Find(`text:"Saved"`)
	app.Advance(3 * time.Second)
	app.WaitForGone(`text:"Saved"`, time.Second)

	text := func(selector string) string {
		t.Helper()
		vs := app.Query(selector)
		if len(vs) != 1 {
			t.Fatalf("%s: %d views match", selector, len(vs))
		}
		return vs[0].Text
	}

	app.Press("id:toggle")
	app.Find(`text:"notifications on"`)

	app.Press("id:seg-Week")
	app.Press("id:step-plus")
	if got := text("id:step-value"); got != "2" {
		t.Errorf("stepper shows %q after one press, want 2", got)
	}
	app.Press("id:step-minus")
	if got := text("id:step-value"); got != "1" {
		t.Errorf("stepper shows %q after undoing, want 1", got)
	}

	app.Press("id:chip-Add one")
	if got := text("id:badge-count"); got != "1" {
		t.Errorf("badge shows %q after one add, want 1", got)
	}
	app.Press("id:chip-Add one")
	if got := text("id:badge-count"); got != "2" {
		t.Errorf("badge shows %q after two adds, want 2", got)
	}
	app.Press("id:chip-Mark read")
	app.WaitForGone("id:badge-count", time.Second)

	app.Press("id:accordion")
	app.Press("id:hamburger")
	app.Find(`text:"Projects"`)

	app.Press("id:utab-Billing")
	app.Find(`text:"One invoice, paid."`)

	app.Press("id:chip-Next")
	app.Press("id:chip-Prev")

	app.Press("id:chip-Finish")
	app.Find(`text:"sent"`)
}

// TestControlsAreLive presses the chips that change a demo's input.
func TestControlsAreLive(t *testing.T) {
	app := nsgotest.Start(t, Root, nsgotest.Size(1180, 860))

	app.Press("id:chip-Shuffle")
	app.Find("id:canvas-bars")

	app.Press("id:tab-Shader")
	app.Press("id:chip-Spacing")
	app.Find("id:shader-pointer")
}

// TestPackPages walks the eight motion-prompt packs: every tab has to
// build, and the first and last tile of its range has to be on it, so a
// pack that silently dropped its tail fails here.
func TestPackPages(t *testing.T) {
	app := nsgotest.Start(t, Root, nsgotest.Size(1180, 860))

	for _, tc := range []struct {
		tab   string
		id    string
		first string
		last  string
	}{
		{"Motion", "001", `text:"001 · Pop-in"`, `text:"024 · Impact shake"`},
		{"Camera", "025", `text:"025 · Pan"`, `text:"048 · L-cut"`},
		{"Landing", "073", `text:"073 · Hero entrance"`, `text:"096 · Loading skeleton"`},
		{"Photo", "097", `text:"097 · Backlighting"`, `text:"120 · Double exposure"`},
		{"Motion2", "121", `text:"121 · Anticipation"`, `text:"144 · Split screen"`},
		{"Landing2", "145", `text:"145 · Timeline on scroll"`, `text:"168 · Background video"`},
		{"Camera2", "169", `text:"169 · Pedestal"`, `text:"192 · Montage"`},
	} {
		app.Press("id:tab-" + tc.tab)
		app.Find(tc.first)
		app.Find(tc.last)
		// The tile's own Replay chip carries the catalog number, which is
		// a second, independent check that the id matches the title.
		app.Press("id:chip-Replay " + tc.id)
	}

	// The ui pack (049–072) sits at the foot of the UI Kit page.
	app.Press("id:tab-UI Kit")
	app.Find(`text:"049 · Modal dialog"`)
	app.Find(`text:"072 · Empty state"`)
}
