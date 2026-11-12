package tpl

// 预制分镜构造器：每个返回标准分镜函数 func(*core.Drawing, *film.Frame)，
// 内部只读 Frame.LocalT（场景局部时间），可直接放进 film.Series / film.Seq。
//
// 所有场景约定 1920×1080 设计画幅（Frame.W/H 自适应），边距 8%。

import (
	"fmt"
	"math"
	"strings"

	"nsgofilm/film"
	"nsgofilm/filmkit"

	"github.com/oneo-me/nsgo/core"
)

const (
	margin    = 0.08 // 页边距（相对画幅宽）
	titleSize = 96
	bodySize  = 42
)

// put 是 tpl 自己的文字绘制（对齐方式：0 左 / 0.5 中 / 1 右）。
func put(d *core.Drawing, s string, x, y, size float32, font string, wt core.FontWeight, c core.Color, alpha float32, align float32) {
	if alpha <= 0.001 || s == "" {
		return
	}
	w := d.MeasureText(s, core.Font{Family: font, Size: core.Points(size), Weight: wt})
	d.SetFillColor(c)
	d.SetAlpha(film.A(alpha))
	d.FillText(s, x-float32(align)*w, y, core.Font{Family: font, Size: core.Points(size), Weight: wt})
	d.SetAlpha(film.A(1))
}

// wrapText 贪心换行：CJK 逐字断、拉丁按词断（和 film 字幕同一策略）。
func wrapText(s string, measure func(string) float32, maxW float32) []string {
	var toks []string
	var latin strings.Builder
	flush := func() {
		if latin.Len() > 0 {
			toks = append(toks, latin.String())
			latin.Reset()
		}
	}
	for _, r := range s {
		switch {
		case r >= 0x2E80:
			flush()
			toks = append(toks, string(r))
		case r == ' ':
			flush()
			if len(toks) > 0 {
				toks[len(toks)-1] += " "
			}
		default:
			latin.WriteRune(r)
		}
	}
	flush()

	var lines []string
	var cur strings.Builder
	for _, tok := range toks {
		if cur.Len() > 0 && measure(cur.String()+tok) > maxW {
			lines = append(lines, strings.TrimRight(cur.String(), " "))
			cur.Reset()
			tok = strings.TrimLeft(tok, " ")
			if tok == "" {
				continue
			}
		}
		cur.WriteString(tok)
	}
	if cur.Len() > 0 {
		lines = append(lines, strings.TrimRight(cur.String(), " "))
	}
	if len(lines) == 0 {
		lines = []string{""}
	}
	return lines
}

// bg 铺风格底色——每个模版场景自己铺，同一部片才能每幕换风格。
func bg(s Style, f *film.Frame, d *core.Drawing) {
	d.SetFillColor(s.Bg)
	d.SetAlpha(film.A(1))
	d.BeginPath()
	d.Rect(0, 0, f.W, f.H)
	d.Fill()
}

// rise 是模板统一的入场：淡入 + 上浮 30px。at 是开始时刻，dur 时长。
func rise(t, at, dur float64) (float64, float32) {
	p := filmkit.Clamp01((t - at) / dur)
	e := filmkit.EaseOut(p)
	return e, float32((1 - e) * 30)
}

// Title 是开场标题卡：巨号标题 + 副标题 + 强调色下划线。
// 适合片头（配 FadeBlack 进）。
func Title(s Style, title, sub string) func(*core.Drawing, *film.Frame) {
	return func(d *core.Drawing, f *film.Frame) {
		bg(s, f, d)
		t := f.LocalT
		cx, cy := f.W / 2, f.H/2-60

		e, dy := rise(t, 0.2, 0.9)
		put(d, title, cx, cy+dy, titleSize, s.TitleFont, core.FontWeightMedium, s.Heading, float32(e), 0.5)

		// 强调色下划线：从中心向两边生长。
		barP := filmkit.Seg(t, 0.9, 1.6)
		if barP > 0 {
			w := 180.0 * float32(barP)
			d.SetFillColor(s.Accent)
			d.SetAlpha(film.A(1))
			d.BeginPath()
			d.Rect(cx-w/2, cy+50, w, 8)
			d.Fill()
		}

		e2, dy2 := rise(t, 1.2, 0.9)
		put(d, sub, cx, cy+150+dy2, 46, s.BodyFont, core.FontWeightNormal, s.Muted, float32(e2), 0.5)
	}
}

// Section 是章节过渡卡：弱化的巨大编号 + 章节标题 + 强调色短线。
// 适合当 FadeBlack 之间的「呼吸」幕。
func Section(s Style, num int, title string) func(*core.Drawing, *film.Frame) {
	return func(d *core.Drawing, f *film.Frame) {
		bg(s, f, d)
		t := f.LocalT
		cx, cy := f.W/2, f.H/2

		// 巨大弱化编号（背景装饰，右移错位——不对称的平衡）。
		e, _ := rise(t, 0.1, 0.8)
		put(d, fmt.Sprintf("%02d", num), cx+f.W*0.14, cy+120, 340, s.TitleFont, core.FontWeightMedium, s.Card, float32(e), 0.5)

		e2, dy2 := rise(t, 0.5, 0.9)
		put(d, title, cx-f.W*0.06, cy+dy2, 88, s.TitleFont, core.FontWeightMedium, s.Heading, float32(e2), 0.5)

		barP := filmkit.Seg(t, 1.0, 1.6)
		if barP > 0 {
			w := 120.0 * float32(barP)
			d.SetFillColor(s.Accent)
			d.SetAlpha(film.A(1))
			d.BeginPath()
			d.Rect(cx-f.W*0.06, cy+50, w, 8)
			d.Fill()
		}
	}
}

// Bullet 是一条要点：强调色标记 + 文字（自动换行）。
func Bullet(s Style, d *core.Drawing, x, y, w float32, alpha float32, s2 string) {
	// 强调色小方块标记。
	d.SetFillColor(s.Accent)
	d.SetAlpha(film.A(alpha))
	d.BeginPath()
	d.Rect(x, y-30, 14, 14)
	d.Fill()
	wrapTextDraw(s2, s, x+36, y, w-36, bodySize, s.Body, alpha, d)
}

func wrapTextDraw(s2 string, st Style, x, y, w float32, size float32, c core.Color, alpha float32, d *core.Drawing) {
	font := core.Font{Family: st.BodyFont, Size: core.Points(size), Weight: core.FontWeightNormal}
	measure := func(s string) float32 { return d.MeasureText(s, font) }
	lines := wrapText(s2, measure, w)
	lh := size * 1.45
	for i, l := range lines {
		d.SetFillColor(c)
		d.SetAlpha(film.A(alpha))
		d.FillText(l, x, y+float32(i)*lh, font)
	}
}

// Bullets 是要点列表卡：标题 + 逐条浮现的要点（stagger 0.55s）。
func Bullets(s Style, title string, items []string) func(*core.Drawing, *film.Frame) {
	return func(d *core.Drawing, f *film.Frame) {
		bg(s, f, d)
		t := f.LocalT
		mx := f.W * margin
		top := f.H * 0.2

		e, dy := rise(t, 0.2, 0.8)
		put(d, title, mx, top+dy, 76, s.TitleFont, core.FontWeightMedium, s.Heading, float32(e), 0)
		// 标题下的强调短线。
		barP := filmkit.Seg(t, 0.7, 1.2)
		if barP > 0 {
			d.SetFillColor(s.Accent)
			d.SetAlpha(film.A(1))
			d.BeginPath()
			d.Rect(mx, top+34, 90*float32(barP), 7)
			d.Fill()
		}

		start := 1.3
		for i, it := range items {
			e, _ := rise(t, start+float64(i)*0.55, 0.7)
			Bullet(s, d, mx+8, top+180+float32(i)*120, f.W*(1-2*margin), float32(e), it)
		}
	}
}

// KPI 是数据卡：n 个大数字（弹簧计数上冲）+ 标签，横排。
// 数字从 0 数到 Value（EaseOut），单位后缀可选。
type KPIItem struct {
	Value  float64
	Unit   string // "%", "万", "×" …
	Label  string
	Decimals int
}

func KPI(s Style, title string, items []KPIItem) func(*core.Drawing, *film.Frame) {
	return func(d *core.Drawing, f *film.Frame) {
		bg(s, f, d)
		t := f.LocalT
		mx := f.W * margin
		top := f.H * 0.2

		e, dy := rise(t, 0.2, 0.8)
		put(d, title, mx, top+dy, 76, s.TitleFont, core.FontWeightMedium, s.Heading, float32(e), 0)

		n := float64(len(items))
		gap := (f.W - 2*mx) / float32(max(n, 1))
		for i, it := range items {
			at := 1.1 + float64(i)*0.4
			p := filmkit.EaseOut(filmkit.Seg(t, at, at+1.4))
			if p <= 0 {
				continue
			}
			v := it.Value * p
			num := fmt.Sprintf("%s%s", trimNum(v, it.Decimals), it.Unit)
			x := mx + gap*(float32(i)+0.5)
			y := f.H * 0.52
			// 数字（大）+ 标签（小）居中于各自格子。
			put(d, num, x, y, 150, s.TitleFont, core.FontWeightMedium, s.Accent, 1, 0.5)
			put(d, it.Label, x, y+80, 40, s.BodyFont, core.FontWeightNormal, s.Body, 1, 0.5)
			// 格子分隔线（弱）。
			if i > 0 {
				d.SetFillColor(s.Card)
				d.SetAlpha(film.A(1))
				d.BeginPath()
				d.Rect(mx+gap*float32(i)-1, y-140, 2, 240)
				d.Fill()
			}
		}
	}
}

func trimNum(v float64, decimals int) string {
	if decimals <= 0 {
		return fmt.Sprintf("%d", int64(math.Round(v)))
	}
	return fmt.Sprintf("%.*f", decimals, v)
}

// Quote 是大引用卡：衬线大字 + 出处，配一个巨大的强调色引号。
func Quote(s Style, text, source string) func(*core.Drawing, *film.Frame) {
	return func(d *core.Drawing, f *film.Frame) {
		bg(s, f, d)
		t := f.LocalT
		cx, cy := f.W/2, f.H/2-40

		// 巨大引号（装饰，弱化）。
		e0, _ := rise(t, 0.1, 0.8)
		put(d, "“", cx-f.W*0.30, cy-120, 260, s.TitleFont, core.FontWeightMedium, s.Accent, 0.5*float32(e0), 0.5)

		font := core.Font{Family: s.TitleFont, Size: core.Points(84), Weight: core.FontWeightNormal}
		measure := func(s string) float32 { return d.MeasureText(s, font) }
		lines := wrapText(text, measure, f.W*(1-2*margin))
		lh := float32(84 * 1.5)
		y0 := cy - float32(len(lines)-1)*lh/2
		for i, l := range lines {
			ei, dyi := rise(t, 0.3+float64(i)*0.25, 0.8)
			put(d, l, cx, y0+float32(i)*lh+dyi, 84, s.TitleFont, core.FontWeightNormal, s.Heading, float32(ei), 0.5)
		}

		e2, _ := rise(t, 0.9+float64(len(lines))*0.2, 0.8)
		put(d, "—— "+source, cx+f.W*0.22, cy+float32(len(lines))*lh+60, 40, s.BodyFont, core.FontWeightNormal, s.Muted, float32(e2), 0)
	}
}

// ImageFull 是全幅图片卡：图片铺满 + 缓慢推近（cinematic 1→1.08）+
// 底部压暗 + 白字说明。img 用 gen_assets 生成的 Image 都行。
func ImageFull(s Style, img film.Image, caption string) func(*core.Drawing, *film.Frame) {
	return func(d *core.Drawing, f *film.Frame) {
		t := f.LocalT
		// 缓慢推近：zoom 1 → 1.08。
		zoom := float32(1 + 0.08*filmkit.Seg(t, 0, 6))
		d.Save()
		f.Cam(d, 0, 0, zoom)
		film.DrawPic(d, img, 0, 0, f.W, f.H, 1)
		d.Restore()
		// 底部压暗（自下而上的黑色渐变），说明字浮在上面。
		g := core.CanvasLinear(0, f.H, 0, f.H*0.45,
			core.GradientStop{Offset: 0, Color: core.RGB(0x000000).Faded(0.55)},
			core.GradientStop{Offset: 1, Color: core.RGB(0x000000).Faded(0)},
		)
		d.SetFillGradient(g)
		d.BeginPath()
		d.Rect(0, f.H*0.4, f.W, f.H*0.6)
		d.Fill()
		e, dy := rise(t, 0.5, 0.9)
		put(d, caption, f.W/2, f.H*0.9+dy, 56, s.BodyFont, core.FontWeightMedium, core.RGB(0xFFFFFF), float32(e), 0.5)
	}
}

// Outro 是片尾卡：居中一句话 + 弱化副行，最后定格。
func Outro(s Style, text, sub string) func(*core.Drawing, *film.Frame) {
	return func(d *core.Drawing, f *film.Frame) {
		bg(s, f, d)
		t := f.LocalT
		cx, cy := f.W/2, f.H/2-40
		e, dy := rise(t, 0.3, 1.0)
		put(d, text, cx, cy+dy, 110, s.TitleFont, core.FontWeightMedium, s.Heading, float32(e), 0.5)
		e2, dy2 := rise(t, 1.0, 1.0)
		put(d, sub, cx, cy+160+dy2, 44, s.BodyFont, core.FontWeightNormal, s.Muted, float32(e2), 0.5)
	}
}
