---
name: nsgo
description: Use when writing, running, testing, packaging or debugging an NSGo app. NSGo (github.com/oneo-me/nsgo) is a cross-platform native UI engine with a Go backend — the Go port of Native Shell (NSH) — targeting macOS, iOS, Windows, Linux, Android and the web. Triggers include creating a new NSGo project, `go tool nsgo` commands, nsgo.yaml, the core/base API (signals, View, Text, Button, Style, Tailwind classes), animation and motion (core.Motion, presets, AnimatedStyle, keyframes, easing, spring, shader-driven motion, 3D scenes), TSX apps on goja, nsgotest selectors, embedding or packaging, or NSGo code that will not compile.
---

# NSGo

A Go UI engine: app logic, the reactive graph, styles, layout and text run in Go; each platform draws with its own native views. Apps are written in Go (or TSX compiled to goja). Design, names and behavior are ported from [NSH](https://github.com/longbridge/nsh) (Rust).

## Iron rule: read the API, never guess it

NSGo is pre-release, has no stable API and almost no public documentation on the web, so recalled identifiers are almost always wrong. Go fails to compile on one wrong field name. Before writing an app, get the truth from the module itself:

```sh
go doc github.com/oneo-me/nsgo/core      # or: go doc -all ./core
go doc github.com/oneo-me/nsgo/base
```

In a checkout of the repo, `ARCHITECTURE.md`, `docs/architecture/*.md` (25 subject docs) and `examples/` are authoritative; `examples/gallery` exercises every element, component and service. If you have neither the module nor the repo, say so — do not substitute another framework (Fyne, Gio, Wails) for the user's actual request.

## Setup

Requires **Go 1.26.0+** (`go.mod` says `go 1.26.0`).

| OS | Native target | Notes |
|---|---|---|
| macOS | ✅ | Full toolchain: builds/signs every target with Go + Xcode |
| Windows | ✅ | Native + web build and run |
| Linux | ✅ | Needs GTK 4.14+, Pango, Cairo; WebKitGTK 6.0 and GStreamer for web views/video |

Add the module and its tool to the app once — `go tool nsgo` then runs it without installing:

```sh
go get -tool github.com/oneo-me/nsgo/cmd/nsgo
```

NSGo publishes **no tagged versions**, so the module resolves to a pseudo-version (`v0.0.0-2026…-46c7d9c`). Some regional mirrors 404 on an untagged module; if `go get` cannot find it, pin the proxy:

```sh
GOPROXY=https://proxy.golang.org,direct GOSUMDB=off go get -tool github.com/oneo-me/nsgo/cmd/nsgo
```

`GOSUMDB=off` is needed only where `sum.golang.org` is unreachable. Working from a local checkout instead, point at it and tidy:

```sh
go mod edit -replace github.com/oneo-me/nsgo=/path/to/nsgo && go mod tidy
```

## Create a project

```sh
mkdir counter && cd counter
go mod init counter
go get -tool github.com/oneo-me/nsgo/cmd/nsgo   # the library and the tool
go tool nsgo init -name counter .              # writes nsgo.yaml, app.go, app_test.go
go mod tidy
go tool nsgo dev .                             # runs it; rebuilds on every save
go tool nsgo test .                            # headless: 2 passed
```

**`init` needs a name.** Bare `go tool nsgo init` in an existing directory fails with `"." is not an app name`, because the default name comes from the directory — and in `.` that is `.`. Pass `-name`, or pass a directory to create:

| Command | Effect |
|---|---|
| `go tool nsgo init -name Counter .` | writes into the current directory, app named `Counter` |
| `go tool nsgo init Counter` | creates `./Counter/` holding the app |
| `go tool nsgo init -template tsx Counter` | creates `./Counter/` as a TSX app |

`init` prints its own next step (`next: go tool nsgo dev Counter`) — trust it over memory.

`-template` is one of `go` (default), `tsx`, `mixed`. The app name takes letters, digits, spaces, `-` and `_`; the bundle id becomes `dev.nsgo.<name>`. Pass `-module` to have `go.mod` written too.

| Template | Files |
|---|---|
| `go` | `nsgo.yaml`, `app.go`, `app_test.go` |
| `tsx` | `nsgo.yaml`, `main.go`, `App.tsx`, `app.test.tsx`, `tsconfig.json`, `.gitignore` |
| `mixed` | the same, with `App.tsx` importing Go of its own and a `native.go` |

`nsgo.yaml` names the app and its targets:

```yaml
name: Hello
id: dev.oneo.nsgo.hello
version: 0.1.0
dev:
  target: native        # native | web | ios | ios_simulator | android
  hot_reload: true
  fps: false
build:
  output: dist
  targets: [macos_arm64, windows_x64, ios_simulator, android_arm64, web]
```

Also valid: `build_number`, `appearance`, `theme`, `i18n`, `icon.{apple,png}`, `dev.icon`, `android.permissions`, `sign.{profile,keystore}`, `web.host`. Beside it, `resources/` holds shipped files (loaded as `asset:`), `assets/` holds project-only files (icons), `locales/` holds `en.yaml`, `de.yaml`… in ICU MessageFormat.

## The app you get (`-template go`)

```go
import (
	"fmt"

	"github.com/oneo-me/nsgo"
	"github.com/oneo-me/nsgo/base"
	"github.com/oneo-me/nsgo/core"
)

func Root(cx *core.Scope) core.Node {
	count := core.NewSignal(cx, 0)
	return core.View{Style: core.Style{FlexGrow: 1, Padding: core.All(24), Gap: 16}, Children: []core.Node{
		core.Text{ContentFunc: func() string { return fmt.Sprintf("Pressed %d times", count.Get()) }},
		base.Button{
			Label:   "Press me",
			OnPress: func(core.PressEvent) { count.Update(func(n int) int { return n + 1 }) },
			Style:   core.Style{Background: core.RGB(0x0969da), CornerRadius: 8, Padding: core.Symmetric(10, 16)},
			Children: []core.Node{core.Text{Content: "Press me", Style: core.Style{Color: core.RGB(0xffffff)}}},
		},
	}}
}

func main() { nsgo.Run(Root) }
```

`nsgo.Run` takes either a root view (`func(cx *core.Scope) core.Node`) or an app function (`func(app *core.Scope)`) that opens its own windows via `desktop`. Options: `WindowSize(w, h)`, `Appearance(...)`, `HideDock()`, `FPS(bool)`, `WindowMaterial(...)`, `Fonts(paths...)`, `Component(name, factory)`.

## Mental model

- **Signals.** `core.NewSignal[T comparable](cx, v)` → `.Get()`, `.Set(v)`, `.Update(fn)`. `NewMemo`, `NewEffect`. The graph belongs to the engine goroutine and never locks; other goroutines enter it with `Scope.Post`.
- **Elements are structs; the zero value is the default** (like `http.Server`). You set only what you want to differ.
- **Every reactive field has a twin**: `Content` + `ContentFunc func() string`, `Style` + `StyleFunc`, `Class` + `ClassFunc`. `XFunc` wins over `X`; setting **both panics at mount**. Fields with no static form take only the function (`Disabled func() bool`, `Selected`, `Expanded`).
- **Components** are any type with `Render(cx *core.Scope) core.Node`; their fields are props. `core.Func` adapts a plain function.
- **Regions** — `Show`, `Each`, `Dynamic`, `Catch` — rebuild only their own span. `Each` is a generic function, not a method.
- **Styling**: typed `core.Style` values, or Tailwind classes in `Class` (`"flex-row gap-4 p-4 bg-white dark:bg-zinc-900"`). `base`'s components are deliberately unstyled — behavior, state and accessibility only.

## Motion

Two clocks. **Commit** = Go re-renders a span (signals, `Show`/`Each`/`Dynamic`). **Frame** = the platform's animator steps `AnimatedStyle` fields and shader uniforms, with no commit. All motion belongs on the frame clock; a value read during render re-renders the subtree and restarts every animation in it.

| Need | Use |
|---|---|
| A state change (hover, press, selection, resize) | `Style.Transition` — a duration, `core.EaseOut`, and a `Properties` bitmask; no graph at all, and skipped under reduce-motion |
| Standard entrance/exit/attention/loop | `core.Motion` + a `core.MotionPreset` (38 of them, all honour reduce-motion) |
| Custom choreography, overlap, holds, swaps | `AnimatedStyle` + `AnimatedNode` graph (`Timing`, `Spring`, `Interpolate`, `Sequence`, `Stagger`, `Loop`) |
| One long take: many beats on one clock, no cuts | one **linear** `0→1` lane + a window (`Interpolate`) per channel — recipe 7 in `references/animation.md` |
| Per-pixel / infinite field / 3D | `core.Shader` uniforms as nodes — zero commits |
| Procedural drawing | `core.Canvas` `Draw` (re-records per change) or `d.Bind` (the platform replays it, no re-record) |

```go
spec := core.FadeInUp.Spec()             // or your own core.MotionSpec
m := core.NewMotion(cx, &spec)           // nil = build at rest, Play later
return core.View{Style: core.Style{Animated: m.Style()}, Children: []core.Node{card()}}
```

Five rules worth more than the API: **(1)** never read a signal or animated value in render — closures or `AnimatedStyle` only; **(2)** restyle instead of rebuild; **(3)** a finite animation *holds* its last value, so the frame at the end must be the readable one (a loop that returns to its start looks blank once it settles); **(4)** an animated element needs a measured parent — an all-absolute `View` measures 0×0 and clips everything; **(5)** `AnimatedStyle` fields are `AnimatedNode`s, not numbers (`core.Constant(1)`, not `1`), and it has no `Width`/`Height`/`Perspective`.

For anything longer than one beat, the whole take is a single **linear** `0→1` value and every channel is a window on it. Hide by `Opacity 0`, never by unmounting: a re-mounted node restarts from zero, and that discontinuity is exactly a cut. Ease inside each window, not on the master lane — an eased master drifts every window boundary.

Read `references/animation.md` before writing any animation: preset table, easing/duration numbers that worked, choreography recipes, shader-and-canvas motion, 3D (ray marching), testing, and a symptom→cause table for motion bugs.

## Dev loop, build, test

| Task | Command |
|---|---|
| Run, rebuild on save | `go tool nsgo dev <dir>` |
| Run in a browser | `go tool nsgo dev -target web <dir>` |
| Run on simulator/device | `go tool nsgo dev -target ios_simulator\|android\|ios <dir>` |
| Trace every commit | `... -trace` |
| Release packages | `go tool nsgo build <dir>` (`-target all`, `-output DIR`, `-run`, `-device D`) |
| Embed in a host app | `... build -embed` → XCFramework (iOS) / AAR (Android) / `nsgo.js` (web); desktops use `nsgo.NewView` |
| Tests | `go tool nsgo test <dir>` (`-run REGEXP`, `-junit FILE`, `-target macos\|windows\|linux\|web\|ios_simulator\|android`, `-fps`) |
| Type-check TSX | `go tool nsgo check <dir>` |
| Regenerate message catalogs | `go tool nsgo i18n <dir>` |

Flags override the `nsgo.yaml` values they name. `build.targets` and `-target` accept `all`, or names like `windows_x64`, `macos_arm64`, `linux_x64`, `ios`, `web`. `build` writes one directory per target under `-output` (default `dist`), named `<name>-<version>-<target>/`. A desktop test target runs only on that OS (`-target windows` on Windows, `-target macos` on a Mac); `web`, `ios_simulator` and `android` run from a Mac.

Tests act like a user. `nsgotest` runs them headless with no window; `-target` runs the same tests in a real window with native input and screenshots (`app.Screenshot("name")`).

```go
func TestMain(m *testing.M) { nsgotest.Main(m) }   // gives the platform the main thread

func TestPress(t *testing.T) {
	app := nsgotest.Start(t, Root)
	app.Press("text:Press me")
	app.Find("text:Pressed 1 times")
}
```

Selectors are Playwright's: `id:`, `label:`, `text:` (a substring; `text:"Save"` or a bare string = whole text), `role:<role> <name>`, `view:<n>`, `/regexp/flags`, and `>>` to descend. Steps: `Press`, `PressDown`/`PressUp`, `LongPress`, `Key`, `KeyPress`, `Type`, `Focus`, `Hover`, `Scroll`, `ScrollTo`, `Drag`, `Touch`, `Back`, `Resize`, `Settle`, `Advance`, `WaitFor`, `WaitForGone`, `Exists`, `Query`, `Texts`, `Dump`. Full list: `go doc github.com/oneo-me/nsgo/nsgotest`.

TSX apps are the same shape: `App.tsx` imports `mount, signal, Text, View` from `"nsgo"` (components from `"nsgo/base"`) and ends with `mount(App)`; `main.go` embeds the sources and calls `nsgo.Run(compiler.App(src, "App.tsx"))`. `*.test.tsx` files run as Go tests. See `references/api.md`.

## Common mistakes

| Mistake | Reality |
|---|---|
| Guessing field names | NSGo is pre-release with almost no public documentation, so recalled identifiers are wrong. `go doc` first — one wrong name is a compile error |
| Setting both `Content` and `ContentFunc` | Panics at mount. Pick one |
| `NewSignal(cx, []string{...})` | `T` must be `comparable`. Use `NewSignalFunc` with your own `equal` |
| `NewSignal` outside a scope | Signals need a `cx`; the graph is goroutine-local, so `cx` is the only handle. From another goroutine use `Scope.Post` |
| Expecting `Button.Label` to render | `Label` is the accessible name (what `role:`/`label:` selectors match). Visible text comes from `Children` |
| Importing `internal/...` | Apps may import only `nsgo`, `core`, `base`, `mobile`, `desktop`, `i18n`, `icons`, `tsx`, `nsgotest` |
| `go build ./...` on Windows or Linux | Fails: `internal/platform/apple/edge.go` is missing its `//go:build darwin` tag and imports darwin-only `purego/objc`. Build your own package (`go build .`), or add the tag in a checkout |
| `JustifyContent: core.AlignCenter` | Two enums: `AlignItems` takes `core.AlignStart/Center/End/Stretch/Baseline`, `JustifyContent` takes `core.JustifyStart/Center/End/SpaceBetween/...`. Mixing them is a compile error |
| `core.Shader{OnMouseMove: ...}` | A shader has `OnPress`, `OnLayout`, `Style` — no pointer-move handler. Put `OnMouseMove func(Point)` / `OnHoverOut func()` on a `core.View` around it |
| Trusting `MustCompileShader` alone | Kage compiles at run time, so a bad shader panics on start-up, not at build. Test each source with `core.CompileShader` and report its error |
| A `Scroll` that never scrolls and shows no bar | `Scroll` defaults to `FlexShrink: 1`, but a `View` parent does not: a column with only `FlexGrow: 1` grows to its content, the scroll view inherits that height and nothing overflows. Give the column `FlexGrow: 1, FlexShrink: 1` (React Native's `flex: 1`). Assert it: the scroll's frame must be shorter than the window |
| Looking for a permanent scrollbar | The indicator is an overlay: native paints a rounded thumb that fades in while scrolling, holds 700ms, then fades in 300ms; `HideIndicator` removes it. On web it is the browser's own bar (`overflow-y: auto`, `scrollbar-width: auto`) |
| An overlay pinned with `Left` only | An absolutely positioned view with one horizontal edge takes the *remaining* width of its containing block, not its content width: a tooltip at `Left: 34` inside a 26pt round anchor is left with nothing and wraps to one character per line. Give it an explicit `Width` |
| Alpha inside an interpolated colour | `AnimatedStyle.Background` fed by `core.Interpolate` with `Colors: {RGBA(0xffffff, 0.12), ...}` paints as near-solid white, not 12%. Use opaque colours for interpolated fills and keep alpha for `Style.Background` |
| `MotionPreset.Name()` as a label | Returns camelCase (`fadeIn`, `zoomIn`, `breathe`), not the Go constant's name |
| An animation "jumps back to the start" | A signal or animated value read during render, or `NewAnimatedValue` built inside the render of a `Dynamic`: the span rebuilds and the values restart. Read inside closures only |
| An animated card shows nothing, no error | Its enclosing `View` has only absolute children and no size, so it measures 0×0 and clips them. Give the stage `Position: Absolute` plus all four insets, or a real height |
| A card is blank after its animation settles | A finite animation holds its final value, and the last frame equals the first. End the track on the readable frame, or loop for real (`core.Loop(a, 0)` is forever) |
| `Scale: 1` / `Opacity: 0.8` will not compile | `AnimatedStyle` fields are `AnimatedNode`: `core.Constant(1)`. Colours need a node too, and interpolated colours must be opaque |
| Expecting multi-window on Windows/Linux | Only macOS: `desktop.OpenWindow`/`OpenPanel`/`NewTray` return nil there |
| Depending on the API being stable | Pre-release: no tags, no releases, no stable API. Pin a commit. macOS also uses private Core Animation APIs for scroll-edge blur, which Mac App Store review may reject |

## Where the truth is

| Question | Read |
|---|---|
| Elements, regions, commit, mount protocol | `docs/architecture/core.md` |
| Styles, state styles, Tailwind | `docs/architecture/styles.md` |
| Text, input, selection | `docs/architecture/text.md` |
| Scroll views, `List`, virtualization | `docs/architecture/scrolling.md` |
| Accessibility, pointers, focus, gestures | `docs/architecture/interaction.md` |
| Services, menus, windows, lifecycle | `docs/architecture/services.md` |
| Components in `base` | `docs/architecture/components.md` |
| TSX, bundling, hot reload, sandbox | `docs/architecture/tsx.md` |
| The `nsgo` tool, packaging, signing | `docs/architecture/tool.md` |
| Testing and selectors | `docs/architecture/testing.md` |
| Animation, motion presets, keyframes, shader and canvas motion, 3D | `references/animation.md` (this skill); `core/motion.go`, `core/animated.go`, `core/canvas.go` |
| What is missing | `docs/roadmap.md` |
