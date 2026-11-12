package filmkit

import (
	"math"
	"testing"
)

func TestSpringEndpoints(t *testing.T) {
	s := Spring{}
	if s.At(0) != 0 {
		t.Fatalf("At(0) = %v want 0", s.At(0))
	}
	if math.Abs(s.At(30)-1) > 1e-6 {
		t.Fatalf("At(30) = %v want ~1", s.At(30))
	}
}

func TestSpringUnderdampedOvershoots(t *testing.T) {
	s := Spring{Stiffness: 120, Damping: 8}
	peak := 0.0
	for t := 0.0; t < 3; t += 0.01 {
		if v := s.At(t); v > peak {
			peak = v
		}
	}
	if peak <= 1.0 {
		t.Fatalf("underdamped should overshoot, peak = %v", peak)
	}
}

func TestSpringOverdampedMonotone(t *testing.T) {
	s := Spring{Stiffness: 100, Damping: 60} // ζ > 1
	prev := 0.0
	for ts := 0.0; ts < 5; ts += 0.01 {
		v := s.At(ts)
		if v < prev-1e-9 || v > 1+1e-9 {
			t.Fatalf("overdamped must be monotone ≤ 1: t=%v v=%v prev=%v", ts, v, prev)
		}
		prev = v
	}
}

func TestBezierEndpointsAndEasing(t *testing.T) {
	ease := Bezier(0.16, 1, 0.3, 1) // easeOutExpo 风格
	if ease(0) != 0 || ease(1) != 1 {
		t.Fatalf("endpoints: %v %v", ease(0), ease(1))
	}
	// 前段陡：x=0.1 时应已过半。
	if v := ease(0.1); v < 0.45 {
		t.Fatalf("should be steep early: ease(0.1) = %v", v)
	}
	linear := Bezier(0, 0, 1, 1)
	if math.Abs(linear(0.25)-0.25) > 1e-6 {
		t.Fatalf("linear bezier: %v", linear(0.25))
	}
}

func TestKF(t *testing.T) {
	times := []float64{2, 5, 8}
	vals := []float64{0, 10, 0}
	// 端点外钳。
	if KF(0, times, vals, nil) != 0 || KF(99, times, vals, nil) != 0 {
		t.Fatal("clamp outside range")
	}
	// 中点线性。
	if v := KF(3.5, times, vals, nil); math.Abs(v-5) > 1e-9 {
		t.Fatalf("mid segment: %v", v)
	}
	// 缓动作用于段内。
	if v := KF(3.5, times, vals, Smooth); math.Abs(v-5) > 1e-9 {
		t.Fatalf("mid with symmetric ease: %v", v)
	}
	if v := KF(3.0, times, vals, EaseOut); v <= 5 {
		t.Fatalf("easeOut should be past linear: %v", v)
	}
	// 第二段回落。
	if v := KF(6.5, times, vals, nil); math.Abs(v-5) > 1e-9 {
		t.Fatalf("second segment: %v", v)
	}
}
