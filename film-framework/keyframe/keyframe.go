// Package keyframe is a direction layer over NSGo's animation graph: a scene
// graph, one clock per take, and tracks addressed by path.
//
// It borrows three ideas from three.js, because they are what makes a
// renderer reasonable to write code against:
//
//   - a scene graph (Take.Add / Node.Add) whose nodes carry a name, so anything
//     can be addressed rather than captured in a closure;
//   - keyframe tracks (Node.On) whose times are seconds into the take, like
//     three.js's KeyframeTrack, so timing reads as a clock instead of as
//     progress arithmetic;
//   - a clock you can set (Take.Seek), like AnimationMixer.setTime, because a
//     video is something you render frame by frame, not something you wait for.
//
// Everything in a take is a function of one value between 0 and 1, so the same
// take plays live (zero commits per frame) and freezes correctly for export.
//
// The one rule it inherits from NSGo: the master clock is LINEAR in time.
// Easing belongs inside a track, never on the clock, or every track boundary
// drifts against the wall clock.
package keyframe

import (
	"fmt"
	"sync"
	"time"

	"github.com/oneo-me/nsgo/core"
)

// Kind is what a node draws.
type Kind int

const (
	// Group is a container: it has no size of its own beyond the stage, and
	// its children position themselves from the stage's top-left corner.
	Group Kind = iota
	// Box is a filled rectangle with an optional corner radius.
	Box
	// Disc is a circle, its diameter the node's width.
	Disc
	// Text is a label; it measures itself unless the node is padded, in
	// which case the padding becomes a plate behind the text.
	Text
	// Camera is a shader node whose uniforms a take animates: the ray-marched
	// scene's eye, look-at, field of view and tint.
	Camera
	// Canvas is a 2D drawing node. Its Draw runs when a signal it read
	// changed, so a take leaves it still and animates the nodes around it.
	Canvas
)

// channels lists what each kind can be told to animate, in the order a
// timeline reads best. Anything else is a lint error, reported with this
// list, because a track on a channel that does not exist is silently nothing.
var channels = map[Kind][]string{
	Group:  {"x", "y", "opacity"},
	Box:    {"x", "y", "scale", "scaleX", "scaleY", "rotate", "opacity"},
	Disc:   {"x", "y", "scale", "scaleX", "scaleY", "rotate", "opacity"},
	Text:   {"x", "y", "scale", "scaleX", "scaleY", "rotate", "opacity"},
	Camera: {"eye.x", "eye.y", "eye.z", "look.x", "look.y", "look.z", "fov", "tint"},
	Canvas: nil,
}

// easedByDefault are the channels whose default easing is ease-in-out: a
// caption or a fade that starts and stops abruptly reads as a bug. Position
// channels default to linear, because a curved camera path is a polyline of
// straight segments — easing each one makes the move stutter at every key.
var easedByDefault = map[string]bool{
	"opacity": true, "scale": true, "scaleX": true, "scaleY": true, "rotate": true,
	"fov": true, "tint": true,
}

// Key is one keyframe: seconds into the take, and the value there.
type Key struct {
	T float64
	V float64
}

// Keyframes is a channel's keys plus the easing that curves each segment — the
// layer's version of three.js's KeyframeTrack, minus the path, which is the
// argument that takes it.
//
// The zero Keyframes is empty. Build one with Keys, KeysAt, Win, FadeIn or
// FadeOut, and curve it with Ease.
type Keyframes struct {
	keys []Key
	ease core.Easing
	set  bool
}

// Keys builds keyframes from (seconds, value) pairs.
func Keys(pairs ...float64) Keyframes {
	if len(pairs)%2 != 0 {
		panic("keyframe: Keys needs (seconds, value) pairs")
	}
	out := Keyframes{keys: make([]Key, 0, len(pairs)/2)}
	for i := 0; i < len(pairs); i += 2 {
		out.keys = append(out.keys, Key{T: pairs[i], V: pairs[i+1]})
	}
	return out
}

// Ease curves every segment with c, instead of the channel's default (linear
// for position and camera channels, ease-in-out for fades and scales).
func (k Keyframes) Ease(c core.Easing) Keyframes {
	k.ease, k.set = c, true
	return k
}

// Linear is the identity easing: this is time itself, no curve.
var Linear = core.Easing{Bezier: core.Linear}

// EaseInOut is the usual choice for a window: slow in, slow out.
var EaseInOut = core.Easing{Bezier: core.EaseInOut}

// Len is how many keyframes the channel has.
func (k Keyframes) Len() int { return len(k.keys) }

// Keys returns the keyframes themselves, for a test or a dump.
func (k Keyframes) Keys() []Key { return k.keys }

// KeysAt is Keys for callers who write durations rather than seconds.
func KeysAt(pairs ...any) Keyframes {
	if len(pairs)%2 != 0 {
		panic("keyframe: KeysAt needs (duration, value) pairs")
	}
	out := Keyframes{keys: make([]Key, 0, len(pairs)/2)}
	for i := 0; i < len(pairs); i += 2 {
		switch v := pairs[i].(type) {
		case time.Duration:
			out.keys = append(out.keys, Key{T: v.Seconds(), V: toFloat(pairs[i+1])})
		case float64:
			out.keys = append(out.keys, Key{T: v, V: toFloat(pairs[i+1])})
		default:
			panic(fmt.Sprintf("keyframe: KeysAt time %d is %T, not time.Duration", i/2, pairs[i]))
		}
	}
	return out
}

func toFloat(v any) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case int:
		return float64(n)
	case time.Duration:
		return n.Seconds()
	default:
		panic(fmt.Sprintf("keyframe: value is %T, not a number", v))
	}
}

// Win is a hold window: 0 until from, 1 across the hold, 0 after to, with a
// ramp on each side. The workhorse of a take — a caption, a panel, a landing.
func Win(from, to, ramp time.Duration) Keyframes {
	f, t, r := from.Seconds(), to.Seconds(), ramp.Seconds()
	if r <= 0 {
		r = 0.2
	}
	if f+r > t-r {
		r = (t - f) / 2
	}
	return Keyframes{keys: []Key{{f, 0}, {f + r, 1}, {t - r, 1}, {t, 0}}}
}

// FadeIn is 0 to 1 over ramp, starting at from.
func FadeIn(from, ramp time.Duration) Keyframes {
	f, r := from.Seconds(), ramp.Seconds()
	return Keyframes{keys: []Key{{f, 0}, {f + r, 1}}}
}

// FadeOut is 1 to 0 over ramp, ending at to.
func FadeOut(to, ramp time.Duration) Keyframes {
	t, r := to.Seconds(), ramp.Seconds()
	return Keyframes{keys: []Key{{t - r, 1}, {t, 0}}}
}

// track is one channel's keyframes, before they are compiled to a node.
type track struct {
	path string
	keys []Key
	ease core.Easing
}

// Ease replaces the track's easing.
func (t *track) Ease(c core.Easing) *track {
	t.ease = c
	return t
}

// Node is one object in a take: what it draws, where it starts, what moves.
type Node struct {
	take *Take
	kind Kind
	name string

	x, y float64
	w, h float64
	r    float64

	fill  core.Color
	color core.Color
	text  string
	label string
	font  float64
	bold  bool
	padH  float64
	padV  float64

	prog *core.ShaderProgram
	draw func(d *core.Drawing, size core.Size)

	// missing is set when a path-addressed track named a node that does not
	// exist: the node still holds the track so Lint can report it, but it
	// draws nothing.
	missing string

	kids   []*Node
	tracks []*track
}

// Take is one continuous piece of footage.
type Take struct {
	name string
	dur  time.Duration
	root *Node

	p       *core.AnimatedValue
	cx      *core.Scope
	started bool
	seekTo  float64
	// version is bumped by a seek so the Frame dynamic rebuilds this take
	// alone: one commit per frozen frame, none while it plays.
	version *core.Signal[int]
}

// New starts a take of length d. A take with no duration is one second long,
// which is almost certainly not what the caller meant.
func New(name string, d time.Duration) *Take {
	if d <= 0 {
		d = time.Second
	}
	t := &Take{name: name, dur: d}
	t.root = &Node{take: t, kind: Group, name: ""}
	return t
}

// Name is the take's identity, used by the seek hook.
func (t *Take) Name() string { return t.name }

// Duration is how long the take runs.
func (t *Take) Duration() time.Duration { return t.dur }

// Add puts a node at the top of the scene graph.
func (t *Take) Add(kind Kind, name string) *Node {
	n := &Node{take: t, kind: kind, name: name, w: 0, h: 0}
	t.root.kids = append(t.root.kids, n)
	return n
}

// Box adds a filled rectangle.
func (t *Take) Box(name string) *Node { return t.Add(Box, name) }

// Disc adds a circle; its diameter is the node's width.
func (t *Take) Disc(name string) *Node { return t.Add(Disc, name) }

// Text adds a label.
func (t *Take) Text(name, text string) *Node { return t.Add(Text, name).Text(text) }

// Canvas adds a 2D drawing node.
func (t *Take) Canvas(name string) *Node { return t.Add(Canvas, name) }

// Camera adds the ray-marched scene's camera.
func (t *Take) Camera(name string, prog *core.ShaderProgram) *Node {
	return t.Add(Camera, name).Program(prog)
}

// Group adds a container.
func (t *Take) Group(name string) *Node { return t.Add(Group, name) }

// Lookup finds a node by name.
func (t *Take) Lookup(name string) *Node { return t.root.find(name) }

// Track records a channel by path, e.g. "camera.fov" or "title.opacity". The
// node methods are shorter, but a path is what data-driven takes have. A path
// naming a node that does not exist keeps the track and leaves the reporting
// to Lint, so a typo shows up as a message rather than a panic in a renderer.
func (t *Take) Track(path string, keys Keyframes) *Node {
	name, ch := splitPath(path)
	if n := t.root.find(name); n != nil {
		return n.On(ch, keys)
	}
	n := &Node{take: t, kind: Text, name: name, missing: path}
	t.root.kids = append(t.root.kids, n)
	return n.On(ch, keys)
}

// On records a channel's keyframes on this node. It returns the node, so a
// take reads as a chain of one-line statements.
func (n *Node) On(channel string, keys Keyframes) *Node {
	tr := &track{path: n.path() + "." + channel, keys: dedupe(keys.keys)}
	tr.ease = core.Easing{Bezier: core.Linear}
	if easedByDefault[channel] {
		tr.ease = core.Easing{Bezier: core.EaseInOut}
	}
	if keys.set {
		tr.ease = keys.ease
	}
	n.tracks = append(n.tracks, tr)
	return n
}

// At moves the node's base position.
func (n *Node) At(x, y float64) *Node { n.x, n.y = x, y; return n }

// Size sets width and height.
func (n *Node) Size(w, h float64) *Node { n.w, n.h = w, h; return n }

// Radius rounds a box's corners, or thins a canvas ring.
func (n *Node) Radius(r float64) *Node { n.r = r; return n }

// Fill sets the background.
func (n *Node) Fill(c core.Color) *Node { n.fill = c; return n }

// Color sets a text's colour.
func (n *Node) Color(c core.Color) *Node { n.color = c; return n }

// Text sets a label's content.
func (n *Node) Text(s string) *Node { n.text = s; return n }

// Label names the node for assistive technology.
func (n *Node) Label(s string) *Node { n.label = s; return n }

// Font sets a label's size in points.
func (n *Node) Font(size float64) *Node { n.font = size; return n }

// Bold makes a label bold.
func (n *Node) Bold() *Node { n.bold = true; return n }

// Pad gives a label a plate: horizontal and vertical padding.
func (n *Node) Pad(h, v float64) *Node { n.padH, n.padV = h, v; return n }

// Program sets a camera's shader.
func (n *Node) Program(p *core.ShaderProgram) *Node { n.prog = p; return n }

// Draw sets a canvas's drawing function.
func (n *Node) Draw(fn func(d *core.Drawing, size core.Size)) *Node { n.draw = fn; return n }

// Add puts a child inside this node.
func (n *Node) Add(kind Kind, name string) *Node {
	kid := &Node{take: n.take, kind: kind, name: name}
	n.kids = append(n.kids, kid)
	return kid
}

// Path is the node's address in the take, used by tracks and by the seek API.
func (n *Node) Path() string { return n.path() }

func (n *Node) path() string {
	if n.name == "" {
		return ""
	}
	return n.name
}

func (n *Node) find(name string) *Node {
	if n.name == name {
		return n
	}
	for _, k := range n.kids {
		if got := k.find(name); got != nil {
			return got
		}
	}
	return nil
}

func (n *Node) track(channel string) *track {
	for _, tr := range n.tracks {
		if tr.path == n.path()+"."+channel {
			return tr
		}
	}
	return nil
}

func splitPath(p string) (string, string) {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '.' {
			return p[:i], p[i+1:]
		}
	}
	return p, ""
}

// dedupe drops keyframes that repeat a time, keeping the first: an
// interpolation needs strictly increasing input.
func dedupe(keys []Key) []Key {
	out := make([]Key, 0, len(keys))
	for _, k := range keys {
		if len(out) > 0 && k.T <= out[len(out)-1].T {
			continue
		}
		out = append(out, k)
	}
	return out
}

// Node compiles the take into an NSGo node and starts the clock on first
// build. Everything under it is a function of the take's one value, so the
// whole take costs no commits per frame.
func (t *Take) Node(cx *core.Scope) core.Node {
	t.attach(cx)
	return t.root.build(cx, nil)
}

// Frame is Node for a take that outlives a page: the graph sits in a dynamic
// over the take's own version, so a Seek commits this subtree and nothing else,
// while playback stays on the frame clock with no commits at all.
//
// This is the piece that makes an export possible. The platform writes animated
// values into the view tree when it builds them; a clock set from outside has to
// ask for a build to show the change, and one commit per rendered frame is
// exactly what an export should cost.
func (t *Take) Frame(cx *core.Scope) core.Node {
	t.attach(cx)
	return core.Dynamic(func(cx *core.Scope) core.Node {
		if t.version != nil {
			t.version.Get() // a seek bumps this, and only this subtree rebuilds
		}
		return t.root.build(cx, nil)
	})
}

// attach binds the take to a scope and starts it the first time it is mounted.
// A different scope means the page was rebuilt under it, so the value node has
// to be made again: the take object survives a rebuild, its binding does not.
func (t *Take) attach(cx *core.Scope) {
	if t.p == nil || t.cx != cx {
		t.cx = cx
		t.version = core.NewSignal(cx, 0)
		t.p = core.NewAnimatedValue(cx, t.seekTo)
		t.started = false
	}
	if !t.started {
		t.started = true
		switch {
		case Recording():
			// A recorder owns the clock: starting playback here would write the
			// same value from two places, every frame. Land on frame zero
			// instead — but not from inside this build. A seek starts an
			// animation, and starting one while the tree is being built
			// re-enters the engine's frame loop (wasm dies of OOM there).
			cx.RequestAnimationFrame(func(time.Time) { t.Seek(0) })
		default:
			t.Play()
		}
	}
	register(t)
}

// Play restarts the take from zero on the frame clock.
func (t *Take) Play() {
	if t.p == nil || t.cx == nil {
		return
	}
	t.p.Set(0)
	cfg := core.TimingConfig{Duration: t.dur, Easing: core.Easing{Bezier: core.Linear}}
	t.cx.StartAnimation(core.Timing(t.p, 1, cfg), nil)
}

// Seek freezes the take at d. This is what an exporter calls: set the clock,
// render one frame, repeat. It stops playback, so a seeked take stays put.
//
// The clock is moved by a one-nanosecond timing rather than by a bare Set. A
// Set changes the value node but not the view tree: the platform writes
// animated properties as an animation steps. A nanosecond timing lands on the
// value the same frame and writes it, which is what makes a frozen frame a
// frame of the take instead of whatever the tree last committed.
func (t *Take) Seek(d time.Duration) {
	if d < 0 {
		d = 0
	}
	if d > t.dur {
		d = t.dur
	}
	t.seekTo = d.Seconds() / t.dur.Seconds()
	if t.p == nil {
		return
	}
	t.p.Set(t.seekTo)
	if t.cx != nil {
		cfg := core.TimingConfig{Duration: time.Nanosecond, Easing: core.Easing{Bezier: core.Linear}}
		t.cx.StartAnimation(core.Timing(t.p, t.seekTo, cfg), nil)
	}
	// And a commit: it costs one build, and it keeps a seek correct on a
	// backend that writes the animated properties at build time instead.
	if t.version != nil {
		t.version.Set(t.version.Get() + 1)
	}
}

// Progress is the fraction the clock was last told to take, so an exporter can
// prove its seek landed instead of trusting a picture.
func (t *Take) Progress() float64 { return t.seekTo }

// Bound reports whether the take has a live clock on a live scope, which is
// what an exporter needs to know before it trusts a frame.
func (t *Take) Bound() bool { return t.p != nil && t.cx != nil }

// Node builds one object of the graph. inherit is the opacity a group has
// already applied to everything inside it, or nil at the top.
func (n *Node) build(cx *core.Scope, inherit core.AnimatedNode) core.Node {
	if n.missing != "" {
		// A track addressed a node that does not exist: draw nothing, and
		// let Lint say so.
		return core.View{Style: core.Style{Width: core.Pt(0), Height: core.Pt(0)}}
	}

	st := core.Style{Position: core.Absolute, Left: core.Pt(0), Top: core.Pt(0)}
	anim := &core.AnimatedStyle{}
	n.place(&st, anim)
	n.look(&st, anim)
	n.box(&st)

	// A group's fade multiplies into its children. The web target does not
	// composite a container's opacity, so leaning on it would make a take
	// target-dependent; core.Mul keeps the fade on the frame clock instead.
	inherited := inherit
	if tr := n.track("opacity"); tr != nil {
		anim.Opacity = n.take.node(tr)
		if inherit != nil {
			anim.Opacity = core.Mul(anim.Opacity, inherit)
		}
		inherited = anim.Opacity
	} else if inherit != nil {
		anim.Opacity = inherit
	}

	var kids []core.Node
	for _, k := range n.kids {
		kids = append(kids, k.build(cx, inherited))
	}

	switch n.kind {
	case Group:
		// A group is a compositing box: fade it and everything inside fades,
		// which is how a shot change is written on one clock.
		st.Position = core.Absolute
		st.Left, st.Top = core.Pt(0), core.Pt(0)
		st.Right, st.Bottom = core.Pt(0), core.Pt(0)
		st.Animated = anim
		return core.View{Style: st, TestID: testID(n), Children: kids}
	case Text:
		// A label keeps the position place() worked out: it is a node like any
		// other, and starting from a fresh style would drop its base.
		st.FontSize = core.Points(n.font)
		st.Color = n.color
		st.AlignItems = core.AlignCenter
		st.JustifyContent = core.JustifyCenter
		if n.bold {
			st.FontWeight = core.FontWeightBold
		}
		if n.padH > 0 || n.padV > 0 {
			st.Background = n.fill
			st.CornerRadius = core.Points(n.r)
			st.Padding = core.Symmetric(core.Points(n.padV), core.Points(n.padH))
		}
		st.Animated = anim
		return core.Text{Content: n.text, Style: st, TestID: testID(n)}
	case Camera:
		st.Animated = anim
		return core.Shader{
			Program:  n.prog,
			Label:    n.label,
			TestID:   testID(n),
			Style:    st,
			Uniforms: n.uniforms(),
		}
	case Canvas:
		st.Animated = anim
		return core.Canvas{Label: n.label, TestID: testID(n), Style: st, Draw: n.draw}
	default:
		st.Background = n.fill
		st.CornerRadius = core.Points(n.r)
		if n.kind == Disc {
			st.CornerRadius = core.Points(n.w / 2)
		}
		st.Animated = anim
		return core.View{Style: st, TestID: testID(n), Children: kids}
	}
}

// box sizes a node: a width or height of zero fills the parent instead, because
// "as big as the stage" is what a full-bleed element wants and what a caller
// expects from leaving it out. A label measures itself, so it is never given a
// right or bottom edge.
func (n *Node) box(st *core.Style) {
	if n.kind == Text {
		return
	}
	if n.w > 0 {
		st.Width = core.Pt(core.Points(n.w))
	} else {
		st.Right = core.Pt(0)
	}
	if n.h > 0 {
		st.Height = core.Pt(core.Points(n.h))
	} else {
		st.Bottom = core.Pt(0)
	}
}

// place puts the node where it starts and animates it from there. A track on
// x or y is an absolute position, so it replaces the base rather than adding
// to it — that is what someone writing keys expects.
func (n *Node) place(st *core.Style, anim *core.AnimatedStyle) {
	if tr := n.track("x"); tr != nil {
		anim.TranslateX = n.take.node(tr)
	} else {
		st.Left = core.Pt(core.Points(n.x))
	}
	if tr := n.track("y"); tr != nil {
		anim.TranslateY = n.take.node(tr)
	} else {
		st.Top = core.Pt(core.Points(n.y))
	}
}

// look adds the channels every drawable node shares. Opacity is handled by
// build, because a group's fade has to be multiplied into its children.
func (n *Node) look(_ *core.Style, anim *core.AnimatedStyle) {
	if tr := n.track("scale"); tr != nil {
		anim.Scale = n.take.node(tr)
	}
	if tr := n.track("scaleX"); tr != nil {
		anim.ScaleX = n.take.node(tr)
	}
	if tr := n.track("scaleY"); tr != nil {
		anim.ScaleY = n.take.node(tr)
	}
	if tr := n.track("rotate"); tr != nil {
		anim.Rotate = n.take.node(tr)
	}
}

// uniforms turns the camera channels into shader uniforms. A channel with no
// track is a plain number, which the platform accepts as a constant lane.
func (n *Node) uniforms() core.Uniforms {
	u := core.Uniforms{}
	vec := func(prefix string, base [3]float64) []any {
		out := make([]any, 3)
		for i, axis := range [...]string{"x", "y", "z"} {
			if tr := n.track(prefix + "." + axis); tr != nil {
				out[i] = n.take.node(tr)
			} else {
				out[i] = base[i]
			}
		}
		return out
	}
	u["Eye"] = vec("eye", [3]float64{0, 0.8, 2.6})
	u["Look"] = vec("look", [3]float64{0, 0.5, 0})
	if tr := n.track("fov"); tr != nil {
		u["Fov"] = n.take.node(tr)
	} else {
		u["Fov"] = 0.8
	}
	if tr := n.track("tint"); tr != nil {
		u["Tint"] = n.take.node(tr)
	} else {
		u["Tint"] = 0.2
	}
	return u
}

// node compiles one track into an interpolation on the take's clock. The
// clock is linear in time, so a key at 4.5s sits at 4.5s.
func (t *Take) node(tr *track) core.AnimatedNode {
	keys := tr.keys
	if len(keys) == 0 {
		keys = []Key{{0, 0}}
	}
	if len(keys) == 1 {
		// A single key is a constant: duplicate it so the input increases.
		keys = []Key{{0, keys[0].V}, {t.dur.Seconds(), keys[0].V}}
	}
	in := make([]float64, len(keys))
	out := make([]float64, len(keys))
	for i, k := range keys {
		in[i] = k.T / t.dur.Seconds()
		out[i] = k.V
	}
	return core.Interpolate(t.p, core.Interpolation{Input: in, Output: out, Easing: tr.ease})
}

// testID names a node in the render tree so tests and the seek hook can find
// it: "take-<name>-<path>".
func testID(n *Node) string {
	if n.path() == "" {
		return ""
	}
	return "take-" + n.take.name + "-" + n.path()
}

// registry lets the platform-side seek hook reach every live take without the
// caller threading a list through the app.
var (
	mu       sync.Mutex
	registry = map[string]*Take{}
	order    []string
)

func register(t *Take) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := registry[t.name]; !ok {
		order = append(order, t.name)
	}
	registry[t.name] = t
}

// Takes lists the live takes, in the order they were first built.
func Takes() []*Take {
	mu.Lock()
	defer mu.Unlock()
	out := make([]*Take, 0, len(order))
	for _, name := range order {
		if t := registry[name]; t != nil {
			out = append(out, t)
		}
	}
	return out
}

// Take finds a live take by name.
func TakeNamed(name string) *Take {
	mu.Lock()
	defer mu.Unlock()
	return registry[name]
}

// SeekAll freezes every live take at d. A take shorter than d clamps, which is
// what an exporter wants at the tail of a sequence.
func SeekAll(d time.Duration) {
	for _, t := range Takes() {
		t.Seek(d)
	}
}

// SeekFraction freezes every live take at f of its own duration, for the case
// where takes of different lengths should be compared at the same beat.
func SeekFraction(f float64) {
	for _, t := range Takes() {
		t.Seek(time.Duration(f * float64(t.dur)))
	}
}
