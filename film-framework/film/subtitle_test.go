package film

import (
	"strings"
	"testing"
)

// fakeMeasure：每字符宽 10，空格宽 5。
func fakeMeasure(s string) float32 {
	w := float32(0)
	for _, r := range s {
		if r == ' ' {
			w += 5
		} else {
			w += 10
		}
	}
	return w
}

func TestWrapSubCJK(t *testing.T) {
	// 六个汉字，每行最多 3 个（30 宽）。
	tokens := subTokens("一二三四五六")
	lines := wrapSub(tokens, fakeMeasure, 30)
	if len(lines) != 2 || lines[0] != "一二三" || lines[1] != "四五六" {
		t.Fatalf("got %v", lines)
	}
}

func TestWrapSubLatin(t *testing.T) {
	// 拉丁词不拆断：hello world(105) 一行，foo 一行。
	tokens := subTokens("hello world foo")
	lines := wrapSub(tokens, fakeMeasure, 110)
	if len(lines) != 2 || lines[0] != "hello world" || lines[1] != "foo" {
		t.Fatalf("got %v", lines)
	}
}

func TestWrapSubMixed(t *testing.T) {
	// 中英混排：CJK 逐字断，拉丁整词断。
	tokens := subTokens("星空star字段")
	lines := wrapSub(tokens, fakeMeasure, 50)
	joined := strings.Join(lines, "|")
	if !strings.Contains(joined, "star") {
		t.Fatalf("latin word should stay whole: %v", lines)
	}
	for _, l := range lines {
		if w := fakeMeasure(l); w > 50 {
			t.Fatalf("line over maxW: %q = %v", l, w)
		}
	}
}

func TestSubTokensSpaces(t *testing.T) {
	// 空格挂词尾；行首不带空格（wrapSub 负责 TrimLeft）。
	tokens := subTokens("a b")
	if len(tokens) != 2 || tokens[0] != "a " || tokens[1] != "b" {
		t.Fatalf("got %q", tokens)
	}
}
