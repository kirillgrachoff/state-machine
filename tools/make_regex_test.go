package tools

import (
	m "state-machine/machine"
	"testing"
)

type innerState struct {
	number uint
}

func (i innerState) Number() uint {
	return i.number
}

func (i innerState) Start() bool {
	return false
}

func (i innerState) Terminate() bool {
	return false
}

func newEdge(from, to uint, with string) m.Edge {
	return m.Edge{From: innerState{from}, To: innerState{to}, With: with}
}

func Equals(a, b []m.Edge) bool {
	cnt := make(map[m.Edge]int)
	for _, e := range a {
		cnt[e]++
	}
	for _, e := range b {
		cnt[e]--
	}
	for _, v := range cnt {
		if v != 0 {
			return false
		}
	}
	return true
}

func TestDeleteVertex(t *testing.T) {
	left := []m.Edge{
		newEdge(0, 1, "0"),
		newEdge(0, 1, "1"),
	}
	loops := []m.Edge{
		newEdge(1, 1, "a"),
		newEdge(1, 1, "b"),
	}
	right := []m.Edge{
		newEdge(1, 2, "2"),
		newEdge(1, 3, "3"),
	}

	out := []m.Edge{
		newEdge(0, 2, "1((a)+(b))*2"),
		newEdge(0, 2, "0((a)+(b))*2"),
		newEdge(0, 3, "1((a)+(b))*3"),
		newEdge(0, 3, "0((a)+(b))*3"),
	}

	ans := deleteVertex(left, loops, right)
	if !Equals(out, ans) {
		t.Fail()
	}
}
