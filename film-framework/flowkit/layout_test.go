package flowkit

import (
	"testing"

	"github.com/oneo-me/nsgo-flow/flowfile"
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

func TestLayoutRanks(t *testing.T) {
	doc, err := flowfile.Parse(sample)
	if err != nil {
		t.Fatal(err)
	}
	c := Layout(doc, TopToBottom, 0, 0, 1200, 800)
	if len(c.Nodes) != 5 || len(c.Edges) != 4 {
		t.Fatalf("nodes/edges: %d/%d", len(c.Nodes), len(c.Edges))
	}
	// rank：start=0，input=1，check=2，lock 与 home=3。
	byID := map[string]Node{}
	for _, nd := range c.Nodes {
		byID[nd.ID] = nd
	}
	if byID["start"].Rank != 0 || byID["input"].Rank != 1 || byID["check"].Rank != 2 {
		t.Fatalf("ranks: %+v %+v %+v", byID["start"], byID["input"], byID["check"])
	}
	if byID["lock"].Rank != 3 || byID["home"].Rank != 3 {
		t.Fatalf("tail ranks: %+v %+v", byID["lock"], byID["home"])
	}
}

func TestLayoutNoOverlap(t *testing.T) {
	doc, _ := flowfile.Parse(sample)
	c := Layout(doc, TopToBottom, 0, 0, 1200, 800)
	// 同 rank 的节点中心距不小于较小宽度的 1.2 倍（不重叠）。
	for r := 0; r <= 3; r++ {
		var inRank []Node
		for _, nd := range c.Nodes {
			if nd.Rank == float64(r) {
				inRank = append(inRank, nd)
			}
		}
		for i := 1; i < len(inRank); i++ {
			dx := inRank[i].X - inRank[i-1].X
			w := (inRank[i].W + inRank[i-1].W) / 2
			if dx < w {
				t.Fatalf("rank %d: nodes overlap (dx %v < half-widths %v)", r, dx, w)
			}
		}
	}
}

func TestEdgesHaveRoutes(t *testing.T) {
	doc, _ := flowfile.Parse(sample)
	c := Layout(doc, TopToBottom, 0, 0, 1200, 800)
	for _, e := range c.Edges {
		if len(e.Pts) < 2 {
			t.Fatalf("edge %s->%s has no route", e.FromID, e.ToID)
		}
		if e.Pts[0] == e.Pts[len(e.Pts)-1] {
			t.Fatalf("edge %s->%s is a point", e.FromID, e.ToID)
		}
	}
}
