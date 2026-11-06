// Command showreel is Motion Lab's film: sixteen seconds of the lab's own
// moves, drawn as one pure function of film time.
//
// Every animation reads filmkit.Now(); the recorder (shots/film.js) steps
// the clock one frame at a time, so the rendered film is deterministic —
// run the capture twice and you get the same pixels. Chapter beats follow
// recipe 7 from the One take page: many effects, one clock.
package main

import (
	"fmt"
	"math"
	"time"

	"github.com/oneo-me/nsgo"
	"github.com/oneo-me/nsgo/core"
	"motionlab/filmkit"
)

// ---------------------------------------------------------------- palette

var (
	bg     = core.RGB(0x0a0e14)
	text   = core.RGB(0xeaf1fa)
	muted  = core.RGBA(0xa2b3c9, 0.78)
	faint  = core.RGBA(0xa2b3c9, 0.45)
	line   = core.RGBA(0x93a7c4, 0.35)
	cyan   = core.RGB(0x35d6de)
	violet = core.RGB(0xa78bfa)
	green  = core.RGB(0x39d98a)
	amber  = core.RGB(0xffb454)
	pink   = core.RGB(0xff6b9d)
	blue   = core.RGB(0x4d8dff)
)

const (
	display = "Helvetica Neue"
	mono    = "Menlo"
	dw, dh  = 1920.0, 1080.0
	total   = 16.0
)

// curAlpha multiplies every drawn alpha, so a scene fades as a whole even
// though its parts set their own opacities.
var curAlpha float32 = 1

// A scales a for the current scene fade.
func A(a float32) float32 { return curAlpha * a }

func main() {
	nsgo.Run(Root, nsgo.WindowSize(1920, 1080))
}

// Root redraws a full-window canvas every frame from the film clock —
// wall time when played, recorder-stepped time when filmed.
func Root(cx *core.Scope) core.Node {
	clock := core.NewSignal(cx, 0.0)
	var loop func(time.Time)
	loop = func(time.Time) {
		clock.Set(filmkit.Now())
		cx.RequestAnimationFrame(loop)
	}
	cx.RequestAnimationFrame(loop)
	return core.Canvas{
		Style: core.Style{Width: core.Pct(100), Height: core.Pct(100), Background: bg},
		Draw: func(d *core.Drawing, size core.Size) {
			render(d, float32(size.Width), float32(size.Height), clock.Get())
		},
	}
}

// ---------------------------------------------------------------- helpers

func txt(d *core.Drawing, s string, x, y, size float32, fam string, wt core.FontWeight, c core.Color, a float32) {
	d.SetFillColor(c)
	d.SetAlpha(A(a))
	d.FillText(s, x, y, core.Font{Family: fam, Size: core.Points(size), Weight: wt})
}

func center(d *core.Drawing, s string, cx, y, size float32, fam string, wt core.FontWeight, c core.Color, a float32) {
	w := d.MeasureText(s, core.Font{Family: fam, Size: core.Points(size), Weight: wt})
	txt(d, s, cx-w/2, y, size, fam, wt, c, a)
}

// strokeLine draws one segment.
func strokeLine(d *core.Drawing, x0, y0, x1, y1, width float32, c core.Color, a float32) {
	d.SetStrokeColor(c)
	d.SetAlpha(A(a))
	d.SetLineWidth(width)
	d.BeginPath()
	d.MoveTo(x0, y0)
	d.LineTo(x1, y1)
	d.Stroke()
}

// dot fills a circle.
func dot(d *core.Drawing, x, y, r float32, c core.Color, a float32) {
	d.SetFillColor(c)
	d.SetAlpha(A(a))
	d.BeginPath()
	d.Arc(x, y, r, 0, 2*math.Pi, false)
	d.Fill()
}

// glow turns a soft shadow on for the next fill, and off again after.
func glow(d *core.Drawing, c core.Color, blur float32, fn func()) {
	d.SetShadow(c, 0, 0, blur)
	fn()
	d.SetShadow(core.RGB(0x000000), 0, 0, 0)
}

// frameProgress strokes a rectangular frame walking its perimeter from the
// top-left corner, drawing fraction p of the whole way round.
func frameProgress(d *core.Drawing, x0, y0, x1, y1, p, width float32, c core.Color, a float32) {
	w, h := x1-x0, y1-y0
	per := 2 * (w + h)
	rem := p * per
	d.SetStrokeColor(c)
	d.SetAlpha(A(a))
	d.SetLineWidth(width)
	d.BeginPath()
	cx, cy := x0, y0
	d.MoveTo(cx, cy)
	for _, e := range [4][2]float32{{w, 0}, {0, h}, {-w, 0}, {0, -h}} {
		seg := e[0] + e[1]
		if rem >= seg { // full edge
			cx, cy = cx+e[0], cy+e[1]
			d.LineTo(cx, cy)
			rem -= seg
			continue
		}
		var f float32
		if seg > 0 {
			f = rem / seg
		}
		d.LineTo(cx+e[0]*f, cy+e[1]*f)
		break
	}
	d.Stroke()
}

// ---------------------------------------------------------------- render

// scene is one beat of the film: draw is called for every frame inside
// [a, b]; adjacent scenes overlap by fade so cuts cross-dissolve.
type scene struct {
	a, b float64
	fn   func(*core.Drawing, float64)
}

var scenes = []scene{
	{0.0, 3.5, sceneTitle},
	{3.2, 7.4, sceneCanvas},
	{7.1, 11.2, sceneMotion},
	{10.9, 14.2, scene3D},
	{13.9, total, sceneEnd}, // the end card holds; no fade-out
}

func render(d *core.Drawing, w, h float32, t float64) {
	s := float32(math.Min(float64(w)/dw, float64(h)/dh))
	ox := (w - float32(dw)*s) / 2
	oy := (h - float32(dh)*s) / 2
	d.Save()
	d.Translate(ox, oy)
	d.Scale(s, s)

	d.SetAlpha(1)
	d.SetFillColor(bg)
	d.BeginPath()
	d.Rect(0, 0, float32(dw), float32(dh))
	d.Fill()

	for _, sc := range scenes {
		if t < sc.a || t > sc.b {
			continue
		}
		// Scenes cross-dissolve over fade seconds at each edge; the last
		// scene skips its fade-out so the final frame is the end card.
		fade := 0.3
		curAlpha = float32(filmkit.Win(t, sc.a, sc.b, fade))
		if sc.b >= total {
			curAlpha = float32(filmkit.Seg(t, sc.a, sc.a+fade))
		}
		if curAlpha <= 0 {
			continue
		}
		sc.fn(d, t)
	}
	curAlpha = 1
	d.Restore()
}

// chapterTag is the label every chapter opens with: a dash, the number and
// the name, sliding in from the left.
func chapterTag(d *core.Drawing, num, name string, c core.Color, t, a float64) {
	p := float32(filmkit.EaseOut(filmkit.Seg(t, a, a+0.45)))
	x := float32(180) + (1-p)*(-40)
	txt(d, "—", x, 260, 30, mono, core.FontWeightNormal, c, 0.9)
	txt(d, num+" · "+name, x+46, 260, 30, mono, core.FontWeightNormal, muted, 0.9)
	ghost := num
	d.SetAlpha(A(0.05))
	d.SetFillColor(c)
	d.FillText(ghost, float32(dw)-420, float32(dh)-120, core.Font{Family: mono, Size: core.Points(430), Weight: core.FontWeightBold})
}

// sceneTitle: the frame writes itself, the name types on, the census counts.
func sceneTitle(d *core.Drawing, t float64) {
	frameProgress(d, 120, 100, 1800, 980, float32(filmkit.EaseOut(filmkit.Seg(t, 0.25, 1.2))), 1.5, line, 1)

	name := "NSGo Motion Lab"
	n := int(filmkit.Seg(t, 0.7, 2.0) * float64(len(name)))
	if n > 0 {
		txt(d, name[:n], 180, 470, 128, display, core.FontWeightBold, text, 1)
	}
	// Caret blinks at 1.2 Hz while typing, then fades.
	if t < 2.5 {
		blink := 0.5 + 0.5*math.Sin(t*2.4*math.Pi)
		ca := filmkit.Seg(t, 0.7, 0.8) * (1 - filmkit.Seg(t, 2.3, 2.5)) * blink
		w := d.MeasureText(name[:max(n, 1)], core.Font{Family: display, Size: core.Points(128), Weight: core.FontWeightBold})
		strokeLine(d, 180+w+14, 380, 180+w+14, 486, 5, cyan, float32(ca))
	}

	if p := filmkit.Seg(t, 2.0, 2.5); p > 0 {
		center(d, "canvas · motion · shader · 3d", 960, 580, 34, mono, core.FontWeightNormal, muted, float32(p))
	}

	// The census counts up under the name.
	census := int(234 * filmkit.EaseOut(filmkit.Seg(t, 2.2, 3.1)))
	line := fmt.Sprintf("13 tabs · %d cards · one clock", census)
	txt(d, line, 180, 880, 30, mono, core.FontWeightNormal, faint, float32(filmkit.Seg(t, 2.2, 2.5)))
}

// sceneCanvas: two waves draw themselves across a ticked baseline.
func sceneCanvas(d *core.Drawing, t float64) {
	chapterTag(d, "01", "CANVAS", cyan, t, 3.5)

	x0, x1, base := float32(340), float32(1580), float32(620)
	// Ticks fade in first.
	tickA := float32(filmkit.Seg(t, 3.6, 3.9)) * 0.5
	for i := 0; i <= 24; i++ {
		x := x0 + (x1-x0)*float32(i)/24
		dot(d, x, base, 2.5, line, tickA)
	}

	// wave draws k of N sample points along a sine.
	wave := func(x0, x1, base, amp, cycles float32, ph float64, prog, a float32, c core.Color, headGlow bool) {
		n := 110
		k := int(prog * float32(n))
		if k < 1 {
			return
		}
		d.SetStrokeColor(c)
		d.SetAlpha(A(a))
		d.SetLineWidth(3)
		d.BeginPath()
		var hx, hy float32
		for i := 0; i <= k; i++ {
			f := float32(i) / float32(n)
			x := x0 + (x1-x0)*f
			y := base - amp*float32(math.Sin(float64(f*cycles)*2*math.Pi+ph))
			if i == 0 {
				d.MoveTo(x, y)
			} else {
				d.LineTo(x, y)
			}
			hx, hy = x, y
		}
		d.Stroke()
		if headGlow && prog < 1 {
			glow(d, c, 18, func() { dot(d, hx, hy, 7, c, a) })
		}
	}

	p1 := float32(filmkit.EaseInOut(filmkit.Seg(t, 3.7, 5.9)))
	wave(x0, x1, base, 110, 2.2, 0, p1, 1, cyan, true)
	p2 := float32(filmkit.EaseInOut(filmkit.Seg(t, 4.6, 6.6)))
	wave(x0, x1, base+40, 70, 3.1, math.Pi/3, p2, 0.8, violet, false)

	// Crest markers pop once the first wave is complete.
	if filmkit.Seg(t, 6.1, 6.3) > 0 {
		pop := float32(filmkit.EaseOutBack(filmkit.Seg(t, 6.1, 6.5)))
		for i := 1; i <= 2; i++ {
			f := (float64(i) - 0.25) / 2.2 / 2 // crests of the 2.2-cycle wave
			x := x0 + (x1-x0)*float32(f)
			y := base - 110
			dot(d, x, y, 6*pop, cyan, 0.9)
		}
	}
}

// sceneMotion: pills stagger in, clear out, and a card flips.
func sceneMotion(d *core.Drawing, t float64) {
	chapterTag(d, "02", "MOTION", violet, t, 7.3)

	cols := []core.Color{blue, cyan, green, amber, pink, violet, blue}
	for i, c := range cols {
		in := filmkit.EaseOutBack(filmkit.Seg(t, 7.5+float64(i)*0.09, 8.05+float64(i)*0.09))
		out := filmkit.EaseOut(filmkit.Seg(t, 9.7+float64(i)*0.07, 10.2+float64(i)*0.07))
		if in <= 0 || out >= 1 {
			continue
		}
		x := float32(300+190*i) + float32(-140*(1-in)) + float32(360*out)
		y := float32(480) + float32(60*(1-in)) - float32(120*out)
		sc := float32(0.6 + 0.4*in)
		a := float32(filmkit.Clamp01(1.4*in) * (1 - out))
		drawPill(d, x, y, 150, 54, sc, c, a)
	}
	if p := filmkit.Seg(t, 8.3, 8.6); p > 0 && t < 9.7 {
		center(d, "stagger · 90 ms", 960, 585, 26, mono, core.FontWeightNormal, faint, float32(p))
	}

	// The flip card: absent until 9.75, then a quick entrance fade, then
	// scaleX is cos(θ) — front shows while cos > 0, back after.
	if t < 9.75 {
		return
	}
	θ := math.Pi * filmkit.EaseInOut(filmkit.Seg(t, 9.9, 10.9))
	cos := float32(math.Cos(θ))
	ac := float32(math.Abs(float64(cos)))
	a := ac * float32(filmkit.Seg(t, 9.75, 9.95))
	if a < 0.001 {
		return
	}
	cx, cy, cw, ch := float32(960), float32(740), float32(360), float32(180)
	d.Save()
	d.Translate(cx, cy)
	d.Scale(ac, 1)
	d.Translate(-cx, -cy)
	rrectPath(d, cx-cw/2, cy-ch/2, cw, ch, 20)
	if cos > 0 {
		d.SetFillColor(core.RGBA(0xffb454, 0.14))
		d.SetAlpha(A(a))
		d.Fill()
		d.SetStrokeColor(amber)
		d.Stroke()
		center(d, "192", cx, cy-6, 76, display, core.FontWeightBold, amber, a)
		center(d, "effects", cx, cy+52, 26, mono, core.FontWeightNormal, faint, a)
	} else {
		d.SetFillColor(core.RGBA(0xa78bfa, 0.16))
		d.SetAlpha(A(a))
		d.Fill()
		d.SetStrokeColor(violet)
		d.Stroke()
		center(d, "1", cx, cy-6, 76, display, core.FontWeightBold, violet, a)
		center(d, "clock", cx, cy+52, 26, mono, core.FontWeightNormal, faint, a)
	}
	d.Restore()
	if p := filmkit.Seg(t, 10.3, 10.6); p > 0 {
		center(d, "flip · cos", 960, 880, 26, mono, core.FontWeightNormal, faint, float32(p))
	}
}

func drawPill(d *core.Drawing, x, y, w, h, sc float32, c core.Color, a float32) {
	d.Save()
	d.Translate(x+w/2, y+h/2)
	d.Scale(sc, sc)
	d.Translate(-(x + w/2), -(y + h/2))
	rrectPath(d, x, y, w, h, h/2)
	d.SetFillColor(core.RGBA(uint32(rgbOf(c)), 0.14))
	d.SetAlpha(A(a))
	d.Fill()
	d.SetStrokeColor(c)
	d.SetAlpha(A(a))
	d.SetLineWidth(1.5)
	d.Stroke()
	d.Restore()
}

func rrectPath(d *core.Drawing, x, y, w, h, r float32) {
	d.BeginPath()
	d.RoundRect(x, y, w, h, r)
}

// rgbOf re-derives the 24-bit rgb of a lab accent for translucent fills.
func rgbOf(c core.Color) uint32 { return uint32(c >> 8 & 0xffffff) }

// scene3D: the wireframe cube turns, projected on the CPU — the same moves
// the 3D page marches on the GPU.
func scene3D(d *core.Drawing, t float64) {
	chapterTag(d, "03", "3D", green, t, 11.1)

	yaw := 25 + 300*filmkit.EaseInOut(filmkit.Seg(t, 11.3, 13.6))
	pitch := -20 + 18*math.Sin(t*1.5)
	dolly := filmkit.Lerp(1.14, 1.0, filmkit.EaseOut(filmkit.Seg(t, 11.0, 12.0)))

	cx, cy, s := float32(960), float32(600), float32(210*dolly)
	cosY, sinY := float32(math.Cos(yaw*math.Pi/180)), float32(math.Sin(yaw*math.Pi/180))
	cosX, sinX := float32(math.Cos(pitch*math.Pi/180)), float32(math.Sin(pitch*math.Pi/180))

	project := func(ix, iy, iz float32) (float32, float32) {
		// yaw about Y, then pitch about X, then a cheap perspective.
		x := ix*cosY + iz*sinY
		z := -ix*sinY + iz*cosY
		y := iy*cosX - z*sinX
		z = iy*sinX + z*cosX
		f := float32(2.2)
		zz := z + 3.2
		return cx + x*f/zz*s, cy + y*f/zz*s
	}

	corners := [8][3]float32{{-1, -1, -1}, {1, -1, -1}, {1, 1, -1}, {-1, 1, -1}, {-1, -1, 1}, {1, -1, 1}, {1, 1, 1}, {-1, 1, 1}}
	edges := [12][2]int{{0, 1}, {1, 2}, {2, 3}, {3, 0}, {4, 5}, {5, 6}, {6, 7}, {7, 4}, {0, 4}, {1, 5}, {2, 6}, {3, 7}}

	var pts [8][2]float32
	for i, c := range corners {
		pts[i][0], pts[i][1] = project(c[0], c[1], c[2])
	}
	for _, e := range edges {
		strokeLine(d, pts[e[0]][0], pts[e[0]][1], pts[e[1]][0], pts[e[1]][1], 2, green, 0.85)
	}
	for _, p := range pts {
		dot(d, p[0], p[1], 5, green, 0.9)
	}
	// Ground shadow under the cube, breathing with the dolly.
	d.SetFillColor(core.RGBA(0x39d98a, 0.08))
	d.SetAlpha(A(1))
	d.BeginPath()
	d.EllipseArc(cx, cy+340, float32(240*dolly), float32(26*dolly), 0, 0, 2*math.Pi, false)
	d.Fill()

	if p := filmkit.Seg(t, 12.2, 12.6); p > 0 {
		center(d, "same moves · projected, not marched", 960, 960, 26, mono, core.FontWeightNormal, faint, float32(p))
	}
}

// sceneEnd: the mark, the repo, and a clean last frame.
func sceneEnd(d *core.Drawing, t float64) {
	pulse := 1 + 0.12*math.Sin(t*2*math.Pi*0.8)
	glow(d, cyan, 24, func() {
		dot(d, 960, 460, float32(10*pulse), cyan, float32(filmkit.Seg(t, 14.0, 14.4)))
	})
	center(d, "Motion Lab", 960, 610, 96, display, core.FontWeightBold, text, float32(filmkit.Seg(t, 14.1, 14.6)))
	center(d, "github.com/HycJack/nsgo-motion-lab", 960, 690, 30, mono, core.FontWeightNormal, faint, float32(filmkit.Seg(t, 14.5, 15.0)))

	// One hairline sweeps out under the repo line.
	p := float32(filmkit.EaseOut(filmkit.Seg(t, 14.6, 15.2)))
	half := 160 * p
	strokeLine(d, 960-half, 740, 960+half, 740, 1, line, float32(filmkit.Seg(t, 14.6, 14.8))*0.6)
}
