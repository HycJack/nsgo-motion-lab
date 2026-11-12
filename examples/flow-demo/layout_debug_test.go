package main

import (
	"testing"

	"github.com/oneo-me/nsgo/core"

	"nsgofilm/flowkit"

	"github.com/oneo-me/nsgo-flow/flowfile"
)

type fakeM struct{}

func (fakeM) MeasureText(s string, f core.Font) float32 { return float32(len([]rune(s)) * 30) }

func TestDebugLayout(t *testing.T) {
	doc, _ := flowfile.Parse(loginFlow)
	c := flowkit.Layout(doc, flowkit.TopToBottom, 307, 151, 1229, 777, fakeM{})
	for _, nd := range c.Nodes {
		t.Logf("%s: x=%v y=%v w=%v h=%v", nd.ID, nd.X, nd.Y, nd.W, nd.H)
	}
	for _, e := range c.Edges {
		_ = e
	}
}
