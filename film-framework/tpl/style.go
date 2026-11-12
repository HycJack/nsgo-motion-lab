// Package tpl 把 html-ppt-designer 的 17 种视觉风格映射成视频模版：
// 每种风格是一组 token（调色板、字体、动效性格），配上预制的分镜
// 构造器（Title / Bullets / KPI / Quote / …），一部「PPT 式视频」
// 就是选一个 Style + 按内容排 Series。
//
// 设计原则沿用 STYLES.md：配色必须有出处、留白是有意图的静默、
// 动效服务叙事节奏、禁止无根据的 glow 和渐变。
package tpl

import (
	"github.com/oneo-me/nsgo/core"
)

// Style 是一种视觉风格的全部 token（对应 STYLES.md 的 CSS 变量）。
type Style struct {
	Name string // 人名（"TED 演讲"）

	Bg      core.Color // 画面底色
	Card    core.Color // 卡片底
	Heading core.Color // 标题字
	Body    core.Color // 正文字
	Muted   core.Color // 弱化字
	Accent  core.Color // 强调色（TED 红 / Apple 蓝 / 瑞士红 …）

	TitleFont string // 标题字体栈（缺字自动回退）
	BodyFont  string // 正文字体栈
}

func rgb(hex uint32) core.Color { return core.RGB(hex) }

// 十七种风格，色值逐条对照 STYLES.md v5.2。

var TED = Style{
	Name: "TED 演讲", Bg: rgb(0x1A1A1A), Card: rgb(0x2C2C2C),
	Heading: rgb(0xFFFFFF), Body: rgb(0xE0E0E0), Muted: rgb(0x999999),
	Accent: rgb(0xE62B1E), TitleFont: "Noto Sans SC", BodyFont: "Noto Sans SC",
}

var Keynote = Style{
	Name: "Apple Keynote", Bg: rgb(0xFFFFFF), Card: rgb(0xF5F5F7),
	Heading: rgb(0x1D1D1F), Body: rgb(0x424245), Muted: rgb(0x86868B),
	Accent: rgb(0x0071E3), TitleFont: "Noto Sans SC", BodyFont: "Noto Sans SC",
}

var TypicalPPT = Style{
	Name: "Typical PPT", Bg: rgb(0xF0F4FA), Card: rgb(0xFFFFFF),
	Heading: rgb(0x1F3864), Body: rgb(0x404040), Muted: rgb(0x808080),
	Accent: rgb(0xFFC000), TitleFont: "Noto Sans SC", BodyFont: "Noto Sans SC",
}

var Gamma = Style{
	Name: "Gamma", Bg: rgb(0xFAFEFF), Card: rgb(0xFFFFFF),
	Heading: rgb(0x2D3436), Body: rgb(0x545454), Muted: rgb(0xA0A0A0),
	Accent: rgb(0x0EA5E9), TitleFont: "Noto Sans SC", BodyFont: "Noto Sans SC",
}

var Consulting = Style{
	Name: "Consulting", Bg: rgb(0xFAFBFC), Card: rgb(0xFFFFFF),
	Heading: rgb(0x0F2137), Body: rgb(0x3A4F66), Muted: rgb(0x7A8FA3),
	Accent: rgb(0xD4A843), TitleFont: "Noto Sans SC", BodyFont: "Noto Sans SC",
}

var Editorial = Style{
	Name: "Editorial", Bg: rgb(0xFEFEFE), Card: rgb(0xFFFFFF),
	Heading: rgb(0x1A1A1A), Body: rgb(0x333333), Muted: rgb(0x888888),
	Accent: rgb(0xC0392B), TitleFont: "Noto Serif SC", BodyFont: "Noto Sans SC",
}

var Swiss = Style{
	Name: "Swiss Style", Bg: rgb(0xFFFFFF), Card: rgb(0xF2F2F2),
	Heading: rgb(0x2C2C2C), Body: rgb(0x404040), Muted: rgb(0x808080),
	Accent: rgb(0xD0021B), TitleFont: "Noto Sans SC", BodyFont: "Noto Sans SC",
}

var Newspaper = Style{
	Name: "Newspaper", Bg: rgb(0xF5F1EB), Card: rgb(0xFFFDF7),
	Heading: rgb(0x1A1A1A), Body: rgb(0x333333), Muted: rgb(0x777777),
	Accent: rgb(0xB22222), TitleFont: "Noto Serif SC", BodyFont: "Noto Sans SC",
}

var Bauhaus = Style{
	Name: "Bauhaus", Bg: rgb(0xFFFFFF), Card: rgb(0xF7F7F7),
	Heading: rgb(0x1A1A1A), Body: rgb(0x333333), Muted: rgb(0x888888),
	Accent: rgb(0xE63B2E), TitleFont: "Noto Sans SC", BodyFont: "Noto Sans SC",
}

// BauhausBlue/BauhausYellow 是包豪斯三原色里的另外两支，拼色块用。
var BauhausBlue = rgb(0x2B4FA2)
var BauhausYellow = rgb(0xF5C300)

var Kinfolk = Style{
	Name: "Kinfolk", Bg: rgb(0xFAF7F2), Card: rgb(0xFFFCF7),
	Heading: rgb(0x5A4A3A), Body: rgb(0x6B5E52), Muted: rgb(0xA09585),
	Accent: rgb(0x9B7B5E), TitleFont: "Noto Serif SC", BodyFont: "Noto Sans SC",
}

var Muji = Style{
	Name: "Muji", Bg: rgb(0xF7F5F0), Card: rgb(0xFFFFFF),
	Heading: rgb(0x3A3A3A), Body: rgb(0x5C5C5C), Muted: rgb(0x999999),
	Accent: rgb(0x8B7355), TitleFont: "Noto Sans SC", BodyFont: "Noto Sans SC",
}

var Brutalist = Style{
	Name: "Brutalist", Bg: rgb(0xFFFFFF), Card: rgb(0xF0F0F0),
	Heading: rgb(0x000000), Body: rgb(0x333333), Muted: rgb(0x777777),
	Accent: rgb(0xFF3D00), TitleFont: "Noto Sans SC", BodyFont: "Noto Sans SC",
}

var NeoTokyo = Style{
	Name: "Neo-Tokyo", Bg: rgb(0x0A0A0F), Card: rgb(0x141420),
	Heading: rgb(0xFFFFFF), Body: rgb(0xC0C0D0), Muted: rgb(0x6A6A80),
	Accent: rgb(0xFF2D6B), TitleFont: "Noto Sans SC", BodyFont: "Noto Sans SC",
}

// NeonCyan/NeonGreen 是 Neo-Tokyo 的另外两支霓虹色。
var NeonCyan = rgb(0x00B4D8)
var NeonGreen = rgb(0x39FF14)

var DarkMode = Style{
	Name: "Dark Mode", Bg: rgb(0x0F172A), Card: rgb(0x1E293B),
	Heading: rgb(0xF1F5F9), Body: rgb(0x94A3B8), Muted: rgb(0x475569),
	Accent: rgb(0x3B82F6), TitleFont: "Noto Sans SC", BodyFont: "Noto Sans SC",
}

var RedBlackTech = Style{
	Name: "红黑白科技", Bg: rgb(0xFFFFFF), Card: rgb(0xF5F5F5),
	Heading: rgb(0x000000), Body: rgb(0x1A1A1A), Muted: rgb(0x808080),
	Accent: rgb(0xE63946), TitleFont: "Noto Sans SC", BodyFont: "Noto Sans SC",
}

var Cartoon = Style{
	Name: "卡通 2.5D", Bg: rgb(0xF0F7FF), Card: rgb(0xFFFFFF),
	Heading: rgb(0x2C3E50), Body: rgb(0x4A5568), Muted: rgb(0xA0AEC0),
	Accent: rgb(0xF39C12), TitleFont: "Noto Sans SC", BodyFont: "Noto Sans SC",
}

var Education = Style{
	Name: "Education", Bg: rgb(0xFFF7ED), Card: rgb(0xFFFFFF),
	Heading: rgb(0x1E293B), Body: rgb(0x4B5563), Muted: rgb(0x9CA3AF),
	Accent: rgb(0xF59E0B), TitleFont: "Noto Sans SC", BodyFont: "Noto Sans SC",
}

// Education 的语义色：蓝=概念 绿=案例 橙=重点 红=注意。
var (
	EduConcept = rgb(0x2563EB)
	EduCase    = rgb(0x10B981)
	EduFocus   = rgb(0xF59E0B)
	EduWarn    = rgb(0xEF4444)
)
