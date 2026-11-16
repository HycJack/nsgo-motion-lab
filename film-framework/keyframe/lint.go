package keyframe

import (
	"fmt"
	"sort"
	"strings"

	"github.com/oneo-me/nsgo/core"
)

// Issue is one thing wrong (or suspicious) about a take. Lint is what makes
// this layer usable by something that cannot see the screen: a track on a
// channel that does not exist, a key past the end, two tracks fighting over
// one channel, or a stretch of seconds where nothing at all moves.
type Issue struct {
	// Path is the node or channel the issue is about; empty for the take.
	Path string
	// Msg says what is wrong, in the imperative.
	Msg string
	// Fatal marks issues that make the take render wrongly rather than
	// uninterestingly.
	Fatal bool
}

func (i Issue) String() string {
	kind := "warn"
	if i.Fatal {
		kind = "error"
	}
	if i.Path == "" {
		return fmt.Sprintf("%s: %s", kind, i.Msg)
	}
	return fmt.Sprintf("%s: %s: %s", kind, i.Path, i.Msg)
}

// Lint checks the take and returns what it found, worst first. An empty result
// means every track names a real channel, every key is inside the take, and
// something moves for most of the running time.
func (t *Take) Lint() []Issue {
	var out []Issue
	seen := map[string]int{} // node name -> how many nodes share it
	byPath := map[string][]*track{}

	t.root.walk(func(n *Node) {
		if n.path() == "" {
			return
		}
		seen[n.path()]++
		for _, tr := range n.tracks {
			byPath[tr.path] = append(byPath[tr.path], tr)
		}
	})

	// Names address nodes, so two nodes with one name make every track on
	// it ambiguous.
	for name, count := range seen {
		if count > 1 {
			out = append(out, Issue{Path: name, Fatal: true,
				Msg: fmt.Sprintf("%d nodes share this name; tracks cannot tell them apart", count)})
		}
	}

	// A channel that does not exist on the kind is silently nothing, which is
	// the worst kind of bug in a take: it renders, it just never moves.
	t.root.walk(func(n *Node) {
		if n.missing != "" {
			out = append(out, Issue{Path: n.missing, Fatal: true,
				Msg: "no node has this name"})
		}
		allowed := channels[n.kind]
		for _, tr := range n.tracks {
			_, ch := splitPath(tr.path)
			if !hasChannel(allowed, ch) {
				out = append(out, Issue{Path: tr.path, Fatal: true,
					Msg: fmt.Sprintf("%s has no channel %q; it has %s", kindName(n.kind), ch, list(allowed))})
			}
		}
	})

	for path, list := range byPath {
		if len(list) > 1 {
			out = append(out, Issue{Path: path, Fatal: true,
				Msg: fmt.Sprintf("%d tracks target one channel; only the last one renders", len(list))})
		}
		tr := list[0]
		for i, k := range tr.keys {
			if k.T < 0 || k.T > t.dur.Seconds() {
				out = append(out, Issue{Path: path, Fatal: true,
					Msg: fmt.Sprintf("key %d is at %.2fs, outside the take's %.2fs", i, k.T, t.dur.Seconds())})
			}
			if i > 0 && k.T <= tr.keys[i-1].T {
				out = append(out, Issue{Path: path, Fatal: true,
					Msg: fmt.Sprintf("key %d at %.2fs does not come after key %d at %.2fs", i, k.T, i-1, tr.keys[i-1].T)})
			}
		}
		if _, ch := splitPath(path); ch == "opacity" {
			for _, k := range tr.keys {
				if k.V < 0 || k.V > 1 {
					out = append(out, Issue{Path: path,
						Msg: fmt.Sprintf("opacity %.2f is outside 0..1", k.V)})
					break
				}
			}
		}
	}

	// A dead stretch is not an error, but it is the thing a take is usually
	// missing: two seconds where nothing moves is a still frame in a film.
	for _, gap := range t.deadSpans() {
		out = append(out, Issue{
			Msg: fmt.Sprintf("%.1fs–%.1fs nothing moves (%.0f%% of the take)",
				gap[0], gap[1], 100*(gap[1]-gap[0])/t.dur.Seconds())})
	}

	sort.SliceStable(out, func(i, j int) bool { return out[i].Fatal && !out[j].Fatal })
	return out
}

// deadSpans finds the stretches longer than a quarter of the take where no
// track changes by more than a rounding error.
func (t *Take) deadSpans() [][2]float64 {
	const step = 0.1
	dur := t.dur.Seconds()
	n := int(dur/step) + 1
	prev := t.Sample(0)
	var spans [][2]float64
	start := -1.0
	for i := 1; i < n; i++ {
		at := float64(i) * step
		cur := t.Sample(seconds(at))
		moved := false
		for path, v := range cur {
			if old, ok := prev[path]; ok && absf(v-old) > 1e-3 {
				moved = true
				break
			}
		}
		if moved {
			if start >= 0 && at-step-start > dur/4 {
				spans = append(spans, [2]float64{start, at - step})
			}
			start = -1
		} else if start < 0 {
			start = at - step
		}
		prev = cur
	}
	if start >= 0 && dur-start > dur/4 {
		spans = append(spans, [2]float64{start, dur})
	}
	return spans
}

// Timeline renders the take's schedule as text: one row per track, with the
// seconds its keys fall on and a sparkline of the curve between them. This is
// the view a take wants when it is wrong — it shows at a glance which rows
// overlap (a handoff) and which rows have a hole (a gap you would have to cut
// to hide).
func (t *Take) Timeline() string {
	type row struct {
		path string
		keys []Key
	}
	var rows []row
	t.root.walk(func(n *Node) {
		for _, tr := range n.tracks {
			rows = append(rows, row{tr.path, tr.keys})
		}
	})
	sort.SliceStable(rows, func(i, j int) bool { return rows[i].path < rows[j].path })

	var b strings.Builder
	fmt.Fprintf(&b, "%s · %.1fs · %d tracks\n", t.name, t.dur.Seconds(), len(rows))
	width := 0
	for _, r := range rows {
		if len(r.path) > width {
			width = len(r.path)
		}
	}
	const cols = 48
	for _, r := range rows {
		var spark []byte
		lo, hi := 0.0, 0.0
		for _, k := range r.keys {
			lo = minf(lo, k.V)
			hi = maxf(hi, k.V)
		}
		if hi-lo < 1e-9 {
			hi = lo + 1
		}
		for c := 0; c < cols; c++ {
			v := t.eval(&track{keys: r.keys, ease: core.Easing{Bezier: core.Linear}},
				float64(c)/float64(cols-1))
			// An ASCII ramp, not block characters: a timeline gets read in a
			// console as often as on a canvas, and block glyphs mangle there.
			level := int(7 * (v - lo) / (hi - lo))
			spark = append(spark, " .:-=+*#"[level])
		}
		times := make([]string, len(r.keys))
		for i, k := range r.keys {
			times[i] = fmt.Sprintf("%.1f", k.T)
		}
		fmt.Fprintf(&b, "  %-*s %s  %s\n", width, r.path, spark, strings.Join(times, " "))
	}
	return b.String()
}

// Dump is the timeline plus a node list: a take as text, so a change can be
// reviewed as a diff instead of as a screenshot.
func (t *Take) Dump() string {
	var b strings.Builder
	b.WriteString(t.Timeline())
	b.WriteString("  nodes\n")
	for _, n := range t.Nodes() {
		anim := len(n.tracks)
		fmt.Fprintf(&b, "    %-16s %-7s %5.1f×%-5.1f %s (%d tracks)\n",
			n.path(), kindName(n.kind), n.w, n.h, colour(n.fill), anim)
	}
	return b.String()
}

func kindName(k Kind) string {
	switch k {
	case Group:
		return "group"
	case Box:
		return "box"
	case Disc:
		return "disc"
	case Text:
		return "text"
	case Camera:
		return "camera"
	case Canvas:
		return "canvas"
	}
	return "node"
}

func colour(c core.Color) string {
	if c == 0 {
		return "-"
	}
	return fmt.Sprintf("#%06x", uint32(c)&0xffffff)
}

func hasChannel(list []string, ch string) bool {
	for _, c := range list {
		if c == ch {
			return true
		}
	}
	return false
}

func list(items []string) string {
	if len(items) == 0 {
		return "no channels"
	}
	return strings.Join(items, ", ")
}

func absf(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

func minf(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func maxf(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
