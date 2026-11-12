package mathkit

import "testing"

func contourAt(x, y float64) Poly {
	return Poly{Pts: []float32{float32(x), float32(y), float32(x + 10), float32(y), float32(x + 10), float32(y + 10)}, Closed: true}
}

func TestPairContoursSquare(t *testing.T) {
	// 等量轮廓：应按质心就近一一对应。
	from := []Poly{contourAt(0, 0), contourAt(100, 0)}
	to := []Poly{contourAt(100, 50), contourAt(0, 50)}
	pairs := pairContours(from, to)
	if len(pairs) != 2 {
		t.Fatalf("pairs: %d", len(pairs))
	}
	// from[0]（质心 5,5）应配 to 里质心 (5,55) 那条 = to[1]。
	fx, fy := from[0].Centroid()
	tx, ty := to[1].Centroid()
	ok := false
	for _, pr := range pairs {
		gx, gy := pr.From.Centroid()
		hx, hy := pr.To.Centroid()
		if gx == fx && gy == fy && hx == tx && hy == ty {
			ok = true
		}
	}
	if !ok {
		t.Fatal("nearest pairing violated")
	}
}

func TestPairContoursExtraFrom(t *testing.T) {
	// from 多一条：三个源对两个目标——目标不重复出现（even-odd 抵消靠
	// 重复源实现，重复目标不允许）。
	from := []Poly{contourAt(0, 0), contourAt(50, 0), contourAt(100, 0)}
	to := []Poly{contourAt(25, 40), contourAt(75, 40)}
	pairs := pairContours(from, to)
	seenTo := map[[2]float64]int{}
	for _, pr := range pairs {
		if pr.Fade {
			continue
		}
		cx, cy := pr.To.Centroid()
		seenTo[[2]float64{cx, cy}]++
	}
	for k, n := range seenTo {
		if n > 1 {
			t.Fatalf("duplicate target contour %v drawn %d times", k, n)
		}
	}
	// from 的每一条都有归宿。
	if len(pairs) != 3 {
		t.Fatalf("pairs: %d want 3", len(pairs))
	}
}

func TestPairContoursExtraTo(t *testing.T) {
	// to 多一条：两个源对三个目标——多余的 to 复制源（重复源在
	// p=0 抵消 = 淡入），但每个目标只画一次。
	from := []Poly{contourAt(25, 0), contourAt(75, 0)}
	to := []Poly{contourAt(0, 40), contourAt(50, 40), contourAt(100, 40)}
	pairs := pairContours(from, to)
	seenTo := map[[2]float64]int{}
	for _, pr := range pairs {
		if pr.Fade {
			continue
		}
		cx, cy := pr.To.Centroid()
		seenTo[[2]float64{cx, cy}]++
	}
	for k, n := range seenTo {
		if n > 1 {
			t.Fatalf("duplicate target contour %v drawn %d times", k, n)
		}
	}
	if len(seenTo) != 3 {
		t.Fatalf("targets covered: %d want 3", len(seenTo))
	}
	if len(pairs) != 3 {
		t.Fatalf("pairs: %d want 3", len(pairs))
	}
}
