package film

import (
	"strings"
	"unicode"

	"github.com/oneo-me/nsgo/core"

	"nsgofilm/filmkit"
)

// Sub 是一句字幕。零值兼容老工程：Size/Color 为零时用 SubStyle 的，
// Pos 为零时在底部。
type Sub struct {
	A, B     float64
	Text     string
	Dialogue bool
	Size     float32    // 0 = SubStyle.Size
	Color    core.Color // 0 = 旁白/对白色
	Pos      SubPos     // 0 = 底部
}

// SubPos 是字幕块的纵向位置。
type SubPos uint8

const (
	SubBottom SubPos = iota // 底部（默认）：末行基线在 SubStyle.Y
	SubTop                  // 顶部：首行基线在 topSubY
	SubCenter               // 居中：字块视觉中心在设计画幅正中
)

// SubStyle 是字幕的样式：配色由片子给，几何用默认值。
type SubStyle struct {
	Font       string     // 字体家族（容器里要有这个字体）
	Size       float32    // 字号
	Y          float32    // 底部字幕的文字基线（设计画幅内）
	Narration  core.Color // 旁白字幕的颜色
	Dialogue   core.Color // 对白字幕的颜色
	Plate      core.Color // 字幕底板的颜色（片子里有亮的东西，白字直接压上去读不出来）
	PlateAlpha float32    // 底板透明度
}

// DefaultSubStyle 生成一套标准字幕样式：46pt、y=952（按 1080p 设计）、
// 进出场各淡 0.35 秒。只把配色交给片子决定。
func DefaultSubStyle(narration, dialogue, plate core.Color) SubStyle {
	return SubStyle{
		Font:       "Noto Sans CJK SC",
		Size:       46,
		Y:          952,
		Narration:  narration,
		Dialogue:   dialogue,
		Plate:      plate,
		PlateAlpha: 0.52,
	}
}

// subFade 是每句字幕进出场淡入淡出的时长。
const subFade = 0.35

// topSubY 是顶部字幕的首行基线。
const topSubY = 170

// 字形度量：Noto Sans CJK SC 的视觉上/下沿（相对基线，em 单位）。
// 底板按字形的视觉中心对齐，不按基线铺——按基线铺字会浮在底板上沿。
const (
	subAscend  = 0.86
	subDescend = 0.14
)

// subMaxWFrac 是自动换行的行宽上限（相对设计画幅宽）。
const subMaxWFrac = 0.62

// subLineGap 是多行字幕的行距倍数。
const subLineGap = 1.3

// subTokens 把字幕拆成可断行的记号：CJK（含全角标点）逐字一断，
// 拉丁文按词一断（空格挂在词尾）。
func subTokens(s string) []string {
	var out []string
	var latin strings.Builder
	flush := func() {
		if latin.Len() > 0 {
			out = append(out, latin.String())
			latin.Reset()
		}
	}
	for _, r := range s {
		switch {
		case r >= 0x2E80: // CJK 与全角标点
			flush()
			out = append(out, string(r))
		case unicode.IsSpace(r):
			flush()
			// 行首不留空格：空格归到前一个词的词尾（见 wrapSub）。
			if len(out) > 0 {
				out[len(out)-1] += " "
			}
		default:
			latin.WriteRune(r)
		}
	}
	flush()
	return out
}

// wrapSub 贪心换行：measure 给出一段文字的宽度，超过 maxW 就断行。
// 纯函数，方便单测。
func wrapSub(tokens []string, measure func(string) float32, maxW float32) []string {
	var lines []string
	var cur strings.Builder
	for _, tok := range tokens {
		w := measure(cur.String() + tok)
		if cur.Len() > 0 && w > maxW {
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

// drawSubs 在屏幕空间画字幕（在镜头 restore 之后画，不受取景影响）。
// 长句按 0.62×画幅宽自动换行，底板随行数走。
func drawSubs(d *core.Drawing, t float64, f Film) {
	st := f.SubStyle
	for _, ln := range f.Subs {
		if t < ln.A || t > ln.B {
			continue
		}
		a := filmkit.Win(t, ln.A, ln.B, subFade)
		if a <= 0.001 {
			continue
		}
		size := st.Size
		if ln.Size > 0 {
			size = ln.Size
		}
		color := st.Narration
		if ln.Dialogue {
			color = st.Dialogue
		}
		if ln.Color != 0 {
			color = ln.Color
		}
		wt := core.FontWeightMedium
		font := core.Font{Family: st.Font, Size: core.Points(size), Weight: wt}
		measure := func(s string) float32 { return d.MeasureText(s, font) }

		// 排版：贪心断行，行距 1.3。
		tokens := subTokens(ln.Text)
		lines := wrapSub(tokens, measure, f.W*subMaxWFrac)
		lineH := size * subLineGap

		// 锚点：各位置的「末行基线」或「首行基线」。
		var firstY float32
		switch ln.Pos {
		case SubTop:
			firstY = topSubY
		case SubCenter:
			firstY = f.H/2 - float32(len(lines)-1)*lineH/2
		default:
			firstY = st.Y - float32(len(lines)-1)*lineH
		}

		// 底板：随行数走，按字形的视觉中心对齐。
		var maxW float32
		for _, l := range lines {
			if w := measure(l); w > maxW {
				maxW = w
			}
		}
		px := maxW/2 + size*0.91
		ph := float32(len(lines)-1)*lineH + size*1.62
		mid := firstY - (subAscend-subDescend)/2*size + float32(len(lines)-1)*lineH/2
		d.SetFillColor(st.Plate)
		d.SetAlpha(A(st.PlateAlpha * float32(a)))
		d.BeginPath()
		d.RoundRect(f.W/2-px, mid-ph/2, px*2, ph, 14)
		d.Fill()

		for i, l := range lines {
			Center(d, l, f.W/2, firstY+float32(i)*lineH, size, st.Font, wt, color, float32(a))
		}
	}
}
