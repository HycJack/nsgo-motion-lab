package flowkit

import (
	"testing"

	"github.com/oneo-me/nsgo-flow/flowfile"
	"github.com/oneo-me/nsgo/core"
)

const sample = `title "登录流程"
direction down
node start "打开登录页" kind=start
node input "输入账号密码"
node check "校验密码" kind=decision
node lock "锁定账号" kind=data
node home "进入首页" kind=end status=done

start -> input
input -> check
check -> home "成功"
check ..> lock "失败 5 次"
`

// fakeMeasure 是测试用的等宽测量（CJK 每字 30、ASCII 每字 15）。
type fakeMeasure struct{}

func (fakeMeasure) MeasureText(s string, font core.Font) float32 {
	w := float32(0)
	for _, r := range s {
		if r > 0x2E80 {
			w += 30
		} else {
			w += 15
		}
	}
	return w
}

func TestLayoutRanks(t *testing.T) {
	doc, err := flowfile.Parse(sample)
	if err != nil {
		t.Fatal(err)
	}
	c := Layout(doc, TopToBottom, 0, 0, 1200, 900, fakeMeasure{})
	if len(c.Nodes) != 5 || len(c.Edges) != 4 {
		t.Fatalf("nodes/edges: %d/%d", len(c.Nodes), len(c.Edges))
	}
	byID := map[string]Node{}
	for _, nd := range c.Nodes {
		byID[nd.ID] = nd
	}
	if !(byID["start"].Y < byID["input"].Y && byID["input"].Y < byID["check"].Y) {
		t.Fatalf("rows out of order: %+v %+v %+v", byID["start"], byID["input"], byID["check"])
	}
	if byID["lock"].Y != byID["home"].Y {
		t.Fatalf("lock/home should share a row: %+v %+v", byID["lock"], byID["home"])
	}
	if byID["lock"].X >= byID["home"].X {
		t.Fatalf("lock should precede home: %+v %+v", byID["lock"], byID["home"])
	}
}

func TestEdgesHaveRoutes(t *testing.T) {
	doc, _ := flowfile.Parse(sample)
	c := Layout(doc, TopToBottom, 0, 0, 1200, 900, fakeMeasure{})
	for _, e := range c.Edges {
		if e.route == nil || len(e.route.points) < 2 {
			t.Fatalf("edge %s->%s has no route", e.FromID, e.ToID)
		}
		if e.route.points[0] == e.route.points[len(e.route.points)-1] {
			t.Fatalf("edge %s->%s is a point", e.FromID, e.ToID)
		}
	}
}
