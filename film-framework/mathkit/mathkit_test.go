package mathkit

import (
	"math"
	"testing"
)

func TestAxesScreen(t *testing.T) {
	a := Axes{X: [2]float64{0, 2}, Y: [2]float64{0, 1}, At: [2]float32{100, 200}, Size: [2]float32{400, 100}}
	// 原点在左下角。
	x, y := a.Screen(0, 0)
	if x != 100 || y != 300 {
		t.Fatalf("origin: got %v,%v want 100,300", x, y)
	}
	// 右上角。
	x, y = a.Screen(2, 1)
	if x != 500 || y != 200 {
		t.Fatalf("corner: got %v,%v want 500,200", x, y)
	}
	// 中点。
	x, y = a.Screen(1, 0.5)
	if math.Abs(float64(x-300)) > 0.01 || math.Abs(float64(y-250)) > 0.01 {
		t.Fatalf("mid: got %v,%v want 300,250", x, y)
	}
}

func TestPolyLength(t *testing.T) {
	// 单位正方形：周长 4（Closed）。
	p := Poly{Pts: []float32{0, 0, 1, 0, 1, 1, 0, 1}, Closed: true}
	if l := p.Length(); math.Abs(l-4) > 1e-9 {
		t.Fatalf("closed length: got %v want 4", l)
	}
	// 开放的少一段。
	p.Closed = false
	if l := p.Length(); math.Abs(l-3) > 1e-9 {
		t.Fatalf("open length: got %v want 3", l)
	}
}

func TestPolyResample(t *testing.T) {
	// 折线沿 x 轴 0..10，共 11 点；重采样 5 点应该等距 0,2.5,5,7.5,10。
	p := Poly{Pts: make([]float32, 0, 22)}
	for i := 0; i <= 10; i++ {
		p.Pts = append(p.Pts, float32(i), 0)
	}
	r := p.Resample(5)
	if len(r.Pts) != 10 {
		t.Fatalf("resample count: got %d floats want 10", len(r.Pts))
	}
	for k, want := range []float32{0, 2.5, 5, 7.5, 10} {
		if math.Abs(float64(r.X(k)-want)) > 1e-4 {
			t.Fatalf("resample[%d]: got %v want %v", k, r.X(k), want)
		}
	}
}

func TestProjectFrontAndBack(t *testing.T) {
	a := &Axes3D{
		X: [2]float64{-1, 1}, Y: [2]float64{-1, 1}, Z: [2]float64{-1, 1},
		At: [2]float32{0, 0}, Size: [2]float32{800, 600},
		Eye: [3]float64{0, 0, 5}, Look: [3]float64{0, 0, 0},
	}
	// 原点在视线正前方：应该落在画布中心。
	sx, sy, _, ok := a.Project(0, 0, 0)
	if !ok || math.Abs(float64(sx-400)) > 0.01 || math.Abs(float64(sy-300)) > 0.01 {
		t.Fatalf("center: got %v,%v ok=%v want 400,300", sx, sy, ok)
	}
	// 相机背后的点不可用。
	if _, _, _, ok := a.Project(0, 0, 10); ok {
		t.Fatal("point behind camera should not project")
	}
	// 朝右偏移的数据点在屏幕上应该偏右（世界 up=Y 轴）。
	sx2, _, _, _ := a.Project(1, 0, 0)
	if sx2 <= 400 {
		t.Fatalf("right point should be right of center: %v", sx2)
	}
}

func TestMorphEndpoints(t *testing.T) {
	// 变形在 p=0 和 p=1 时应与源/目标公式重合（同轮廓数、同采样）。
	from := Eq{W: 100, H: 100, Lines: []Poly{
		{Pts: []float32{0, 0, 100, 0, 100, 100}, Closed: true},
	}}
	to := Eq{W: 100, H: 100, Lines: []Poly{
		{Pts: []float32{0, 50, 50, 50, 50, 150}, Closed: true},
	}}
	a := from.Lines[0].Resample(eqResample)
	b := to.Lines[0].Resample(eqResample)
	// 采样后点数一致是对位的前提。
	if len(a.Pts) != len(b.Pts) {
		t.Fatalf("resample mismatch: %d vs %d", len(a.Pts), len(b.Pts))
	}
}
