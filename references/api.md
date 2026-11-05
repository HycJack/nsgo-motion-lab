# NSGo API reference

Verified against `github.com/oneo-me/nsgo` at commit `46c7d9c`. Confirm anything you are unsure of with `go doc`; this file is a map, not the source of truth.

## Imports an app may use

`github.com/oneo-me/nsgo` (`Run`, `NewView`, options) · `.../core` · `.../base` · `.../base/textview`, `base/textview/highlight` · `.../base/editor`, `base/diff` · `.../mobile` · `.../desktop` · `.../i18n` · `.../icons/lucide` · `.../tsx` (+ `tsx/compiler`, `tsx/bundler`) · `.../nsgotest`

Never import `.../internal/...`.

## Reactive graph

```go
func NewSignal[T comparable](cx *Scope, v T) *Signal[T]
func NewSignalFunc[T any](cx *Scope, v T, equal func(a, b T) bool) *Signal[T]
func NewMemo[T comparable](cx *Scope, fn func() T) *Memo[T]
func NewMemoFunc[T any](cx *Scope, fn func() T, equal func(a, b T) bool) *Memo[T]

func (s *Signal[T]) Get() T
func (s *Signal[T]) GetUntracked() T
func (s *Signal[T]) Set(v T)
func (s *Signal[T]) Update(fn func(T) T)
```

`T comparable` is a hard constraint — use the `Func` variants for slices, maps and structs holding them.

`Scope` is also the app's handle on everything else — the graph belongs to one goroutine, so `cx` is how code reaches it:

| Need | Calls |
|---|---|
| From another goroutine | `Post(fn)`, `Context() context.Context` |
| Ownership and lifetime | `Child()`, `Effect(fn)`, `Cleanup(fn)`, `Dispose()`, `Batch(fn)`, `Untrack(fn)` |
| Timers and frames | `AfterFunc(d, fn)`, `RequestAnimationFrame(fn)`, `Animate(d, curve, step)`, `Now()` |
| App state | `Storage()`, `DataDir()`, `CacheDir()`, `SetLocale(l)`, `Theme()`, `SetTheme(t)`, `Dark()`, `Palette()`, `Env()` |
| Windows and chrome | `SetMainView`, `OpenWindow`, `OpenPanel`, `NewTray`, `SetTitle`, `SetWindowMaterial`, `ToggleFullScreen`, `Quit`, `ExitApp`, `SetShutdownMode` |
| System services | `SetClipboard`, `ReadClipboard`, `PickFiles`, `SaveFile`, `OpenURL`, `Share`, `ShowAlert`, `PopupMenu`, `RegisterHotKey`, `RegisterFont` |

Asynchronous services take a `done` callback instead of returning a value. Full list: `go doc github.com/oneo-me/nsgo/core Scope`.

## Host elements (`core`)

`View`, `Text`, `Input`, `Scroll`, `List`, `Image`, `Canvas`, `Overlay`, `Modal`, `Screen`, `WebView`, `Video`, `Native`, `Shader`, `Icon`.

Apps normally take `Input`, `Screen` and `Modal` through `base.Input`, `base.NavPage` and `base.Dialog`, and `mobile.InputAccessoryView` on phones.

Common fields, all elements:

| Field | Meaning |
|---|---|
| `Class` / `ClassFunc` | Tailwind CSS 4 classes, under `Style` |
| `Style` / `StyleFunc` | typed style |
| `Children []Node` | child nodes |

`View` adds press and pointer handlers (`OnPress`, `OnLongPress`, `OnPressDown`, `OnPressUp func(PressEvent)`, `OnHoverIn`, `OnHoverOut func()`, `OnHoverMove func(HoverEvent)`, `OnMouseMove func(Point)`, `OnContextMenu func(ContextMenuEvent)`, `OnTouch func(TouchEvent)`), plus refs, tooltips and platform-specific sub-objects (`View.IOS.BackdropMaterial`, `View.Android.Ripple`, `View.Desktop.WindowDrag`).

`Text` adds `Content` / `ContentFunc`, `MaxLines`, `EllipsizeMode`, `NoFontScaling`, `MaxFontSizeMultiplier`.

**Every reactive field has a `…Func` twin.** Setting both panics at mount. Fields with no static form take only the function: `Disabled func() bool`, `Selected`, `Expanded`.

## Styles

```go
core.Style{
    FlexGrow: 1, AlignSelf: core.AlignStart,
    Padding: core.All(24),                 // Edges
    PaddingEnv: core.EnvAll(core.EnvSafeArea),
    Gap: 16,
    Background: core.RGB(0x0969da),        // Color
    CornerRadius: 8,
    FontSize: 28, FontWeight: core.FontWeightBold,
    Color: core.RGB(0xffffff), Opacity: core.Alpha(0.8),
    Active: &core.Style{...},              // state styles: Checked, Indeterminate,
                                           // Selected, Hover, Focus, Active, Disabled
}
```

Constructors seen in the wild: `All(p)`, `Symmetric(vertical, horizontal)`, `RGB(hex)`, `Alpha(f)`, `EnvAll(...)`, `EnvSafeArea`. Layout follows React Native's defaults as NSH does: column direction, no shrinking, no automatic minimum sizes, `flex-start` content alignment outside grids.

## `base` components (behavior only, no look)

`Button`, `Checkbox`, `Radio`, `RadioGroup`, `Switch`, `Tabs`/`Tab`, `Progress`, `Tooltip`, `Popover`, `Menu`/`MenuItem`/`Submenu`/`MenuSeparator`, `ContextMenu`, `Dialog`, `AlertDialog`, `Sheet`, `Toast`, `Collapsible`, `ResizablePanel`/`ResizablePanelGroup`/`ResizeHandle`, `DockArea`, `Table`, `DataTable`, `Plot`, plus `base.Input`, `base.NavPage`, `base.TextView` (`base/textview`), `base.Editor` (`base/editor`), `base.DiffView` (`base/diff`).

Not yet in `base` (see `docs/roadmap.md`): `Slider`, `NumberInput`, `OtpInput`, `Toggle`, `Accordion`, `Toolbar`, `Pagination`, `Tree`, `Carousel`, `Link`, `Avatar`, `Select`, `Combobox`, `DatePicker`, `ColorPicker`, `Command`, `HoverCard`. `examples/gallery-ui` has styled equivalents of some (Onevra's), which are not part of `base`.

```go
base.Button{
    Label:    "Press me",                     // accessible name, not the visible text
    OnPress:  func(core.PressEvent) { ... },
    Children: []core.Node{core.Text{Content: "Press me"}},  // what shows
    Disabled: func() bool { return n.Get() == 0 },
    Class:    "rounded-lg bg-blue-600 px-4 py-2 active:opacity-80",
}
```

## Components and regions

```go
type Component interface{ Render(cx *core.Scope) core.Node }

type Counter struct{ Start int }                       // fields are props
func (c Counter) Render(cx *core.Scope) core.Node { ... }
```

Regions rebuild only their own span: `core.Show`, `core.Each` (generic function), `core.Dynamic`, `core.Catch` (recovers panics from below it). `core.Func` adapts a plain function to a component. `core.NewHolder` + `Place` moves built content between places keeping state.

## Running

```go
nsgo.Run(Root)                              // func(cx *core.Scope) core.Node
nsgo.Run(App)                               // func(app *core.Scope): opens its own windows
nsgo.SetMainView(cx, Root)
nsgo.NewView(...)                           // embed in a host desktop window
```

Options: `WindowSize(w, h)`, `Appearance(core.AppearanceLight|Dark)`, `HideDock()`, `FPS(bool)`, `WindowMaterial(core.WindowMaterial)`, `Fonts(paths...)`, `Component(name, core.NativeFactory)`.

## Testing (`nsgotest`)

```go
func TestMain(m *testing.M) { nsgotest.Main(m) }
app := nsgotest.Start(t, Root, nsgotest.Size(390, 844), nsgotest.Environment(core.Environment{OS: core.OSIOS, Idiom: core.IdiomPhone, Locale: "en-US"}))
```

Selectors: `id:` · `label:` · `text:` (substring; `text:"Save"` or a bare string = whole text; `/regexp/flags`) · `role:<role> <name>` (substring unless quoted) · `view:<node>` · `>>` to descend. A bare string matches a test ID, else a label, else a whole text.

| Group | Methods |
|---|---|
| Find | `Find`, `FindMounted`, `Query`, `QueryMounted`, `Exists`, `Enabled`, `Texts`, `Dump`, `WaitFor`, `WaitForGone` |
| Act | `Press`, `PressDown`, `PressUp`, `LongPress`, `Key`, `KeyPress`, `Type`, `Focus`, `Hover`, `HoverAt`, `HoverOut`, `Scroll`, `ScrollTo`, `DragScroll`, `Drag`, `Touch`, `Back`, `HeaderBack` |
| Model | `Settle`, `Advance`, `FastForward`, `Resize`, `SetEnvironment`, `Storage`, `Post`, `Alerts`, `Answer`, `Screenshot`, `AnimatedValue` |

129 methods in all; `go doc github.com/oneo-me/nsgo/nsgotest` lists them.

## TSX

```go
//go:embed *.tsx *.ts
var src embed.FS

func main() { nsgo.Run(compiler.App(src, "App.tsx")) }
```

```tsx
import { mount, signal, Text, View } from "nsgo";
import { Button } from "nsgo/base";

function App() {
  const count = signal(0);
  return (
    <View class="flex-1 gap-4 p-6">
      <Text class="text-3xl font-bold">Pressed {count()} times</Text>
      <Button onPress={() => count.update((n) => n + 1)}>
        <Text>Press me</Text>
      </Button>
    </View>
  );
}

mount(App);
```

Modules: `nsgo`, `nsgo/base`, `nsgo/mobile`, `nsgo/desktop`, `nsgo/test`. Go components, functions and classes reach TSX through `tsx.Register("app/native", map[string]core.Node{"Badge": Badge{}})`, `tsx.RegisterFunctions`, `tsx.RegisterClass`. Tests are `*.test.tsx` and run as Go tests; type declarations come from `go tool nsgo check`.

## `nsgo.yaml`

```yaml
name: Hello
id: dev.oneo.nsgo.hello
version: 0.1.0
build_number: 1                    # CFBundleVersion / versionCode
appearance: dark                  # or light
theme: theme.json                 # NSGo's or a shadcn/ui theme
i18n: { dir: locales, source: en }
icon: { apple: assets/AppIcon.icon, png: assets/Icon.png }
dev:   { target: native, hot_reload: true, fps: false, icon: {...} }
build: { output: dist, open: false, fps: false, targets: [macos_arm64, windows_x64, linux_x64, ios, ios_simulator, android_arm64, android_x64, web] }
android: { permissions: [REQUEST_INSTALL_PACKAGES] }
sign: { profile: assets/AppStore.mobileprovision, keystore: assets/release.p12, key_alias: ... }
web: { host: ... }
```

Secrets come from `NSGO_KEYSTORE_PASSWORD` and `NSGO_KEY_PASSWORD`.
