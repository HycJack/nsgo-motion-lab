package filmkit

import (
	"math"
	"testing"
)

func TestSegClampsAndEndpoints(t *testing.T) {
	cases := []struct {
		t, a, b, want float64
	}{
		{t: -1, a: 0, b: 1, want: 0},
		{t: 0, a: 0, b: 1, want: 0},
		{t: 0.5, a: 0, b: 1, want: 0.5},
		{t: 1, a: 0, b: 1, want: 1},
		{t: 2, a: 0, b: 1, want: 1},
		// Reversed window (a > b) climbs from 0 to 1 as t decreases.
		{t: 1, a: 1, b: 0, want: 0},
		{t: 0, a: 1, b: 0, want: 1},
	}
	for _, c := range cases {
		if got := Seg(c.t, c.a, c.b); math.Abs(got-c.want) > 1e-9 {
			t.Errorf("Seg(%v,%v,%v) = %v, want %v", c.t, c.a, c.b, got, c.want)
		}
	}
}

func TestEaseCurvesEndpointsAndBounds(t *testing.T) {
	fns := map[string]func(float64) float64{
		"smooth": Smooth, "easeOut": EaseOut, "easeInOut": EaseInOut, "easeOutBack": EaseOutBack,
	}
	for name, fn := range fns {
		if fn(0) != 0 || fn(1) != 1 {
			t.Errorf("%s endpoints: f(0)=%v f(1)=%v, want 0 and 1", name, fn(0), fn(1))
		}
		for i := 0; i <= 20; i++ {
			x := float64(i) / 20
			v := fn(x)
			if v < -0.001 || v > 1.105 { // easeOutBack may overshoot ~10%
				t.Errorf("%s(%v) = %v out of range", name, x, v)
			}
		}
	}
	// Everything but the back-ease is monotonic and never exceeds 1.
	for name, fn := range fns {
		if name == "easeOutBack" {
			continue
		}
		prev := fn(0)
		for i := 1; i <= 20; i++ {
			v := fn(float64(i) / 20)
			if v < prev-1e-9 || v > 1+1e-9 {
				t.Errorf("%s not monotonic within [0,1]: f(%v)=%v after %v", name, i, v, prev)
			}
			prev = v
		}
	}
}

func TestWinAndPulse(t *testing.T) {
	if Win(0, 1, 3, 0.5) != 0 {
		t.Error("Win before the window should be 0")
	}
	if Win(2, 1, 3, 0.5) != 1 {
		t.Error("Win mid-window should be 1")
	}
	if Win(4, 1, 3, 0.5) != 0 {
		t.Error("Win after the window should be 0")
	}
	if Pulse(0, 0, 2) > 0.001 || Pulse(2, 0, 2) > 0.001 {
		t.Errorf("Pulse at the edges = %v/%v, want ~0", Pulse(0, 0, 2), Pulse(2, 0, 2))
	}
	if math.Abs(Pulse(1, 0, 2)-1) > 1e-9 {
		t.Errorf("Pulse mid-window = %v, want 1", Pulse(1, 0, 2))
	}
}
