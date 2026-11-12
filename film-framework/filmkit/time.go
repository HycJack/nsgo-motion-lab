// Package filmkit turns an NSGo scene into a film: one injectable clock and
// the time-window math every frame-driven scene needs.
//
// A film scene is a pure function of filmkit.Now(). The recorder (shots/film.js)
// steps a JS global (__filmT) one frame at a time and screenshots, so the
// rendered film is deterministic — same input, same pixels, every run. Native
// builds have no recorder and fall back to wall time.
package filmkit

import "math"

// Clamp01 pins x to [0, 1].
func Clamp01(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}

// Seg is the progress of t through [a, b], clamped to [0, 1]. It is the
// one-lane primitive every timed beat decomposes into: start, ramp, hold,
// exit — each its own Seg.
func Seg(t, a, b float64) float64 { return Clamp01((t - a) / (b - a)) }

// Lerp blends a→b by p (unclamped).
func Lerp(a, b, p float64) float64 { return a + (b-a)*p }

// Lerp32 is Lerp in float32, for drawing coordinates.
func Lerp32(a, b, p float32) float32 { return a + (b-a)*p }

// F64 truncates to float32 for the drawing API.
func F32(x float64) float32 { return float32(x) }

// Smooth is smoothstep: an S-curve that eases both ends, for fades and
// camera moves that should read as mechanical, not springy.
func Smooth(x float64) float64 {
	x = Clamp01(x)
	return x * x * (3 - 2*x)
}

// EaseOut decelerates into the target — the default entrance curve.
func EaseOut(x float64) float64 {
	x = Clamp01(x)
	return 1 - math.Pow(1-x, 3)
}

// EaseInOut is the symmetric S-curve, for rotations and traversals.
func EaseInOut(x float64) float64 {
	x = Clamp01(x)
	if x < 0.5 {
		return 4 * x * x * x
	}
	return 1 - math.Pow(-2*x+2, 3)/2
}

// EaseOutBack overshoots by about ten percent before settling — for pops
// that should feel a little too eager. The endpoints short-circuit so the
// curve lands exactly on 0 and 1 (the polynomial leaves float dust there).
func EaseOutBack(x float64) float64 {
	const c1, c3 = 1.70158, 2.70158
	if x <= 0 {
		return 0
	}
	if x >= 1 {
		return 1
	}
	m := x - 1
	return 1 + c3*m*m*m + c1*m*m
}

// Win is a hold window with soft edges: 0 before a, ramping to 1 over ramp,
// holding through the middle, ramping back to 0 over the last ramp before b.
// Fades the whole scene through one call.
func Win(t, a, b, ramp float64) float64 {
	return Seg(t, a, a+ramp) * (1 - Seg(t, b-ramp, b))
}

// Pulse rises and falls inside [a, b]: 0 at the ends, peaking at 1 in the
// middle. For one-shot emphasis that must return to rest.
func Pulse(t, a, b float64) float64 {
	p := Seg(t, a, b)
	return Smooth(1 - math.Abs(2*p-1))
}
