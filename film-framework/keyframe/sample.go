package keyframe

import (
	"time"
)

// Sample evaluates every track at d, in the take's own units, without asking
// the platform. It is the pure half of the take: the platform solves the same
// curve — core.Easing.At — so these are the numbers the render tree shows,
// which is what makes an export checkable and a timeline testable.
//
// Values are the track's own: x and y in points, opacity a fraction, scale a
// factor, and the camera channels whatever the shader expects.
func (t *Take) Sample(d time.Duration) map[string]float64 {
	f := d.Seconds() / t.dur.Seconds()
	out := map[string]float64{}
	t.root.walk(func(n *Node) {
		for _, tr := range n.tracks {
			out[tr.path] = t.eval(tr, f)
		}
	})
	return out
}

// At is Sample for one track, or 0 when the path has no track.
func (t *Take) At(path string, d time.Duration) float64 {
	f := d.Seconds() / t.dur.Seconds()
	var out float64
	t.root.walk(func(n *Node) {
		if n.tracks == nil {
			return
		}
		for _, tr := range n.tracks {
			if tr.path == path {
				out = t.eval(tr, f)
			}
		}
	})
	return out
}

// eval is the piecewise-linear evaluation an Interpolation does: clamp outside
// the keys, ease inside each segment.
func (t *Take) eval(tr *track, f float64) float64 {
	keys := tr.keys
	if len(keys) == 0 {
		return 0
	}
	if len(keys) == 1 {
		return keys[0].V
	}
	at := func(k Key) float64 { return k.T / t.dur.Seconds() }
	if f <= at(keys[0]) {
		return keys[0].V
	}
	for i := 1; i < len(keys); i++ {
		a, b := keys[i-1], keys[i]
		x0, x1 := at(a), at(b)
		if f > x1 {
			continue
		}
		if x1 == x0 {
			return b.V
		}
		u := tr.ease.At((f - x0) / (x1 - x0))
		return a.V + (b.V-a.V)*u
	}
	return keys[len(keys)-1].V
}

func (n *Node) walk(fn func(*Node)) {
	fn(n)
	for _, k := range n.kids {
		k.walk(fn)
	}
}

// Nodes lists every named node in the take, for a dump or a test.
func (t *Take) Nodes() []*Node {
	var out []*Node
	t.root.walk(func(n *Node) {
		if n.path() != "" {
			out = append(out, n)
		}
	})
	return out
}
