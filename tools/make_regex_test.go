package tools

import (
	edge "state-machine/edge_machine"
	m "state-machine/machine"
	"testing"
)

func newEdge(from, to uint, with string) m.Edge {
	return m.Edge{From: innerState{number: from}, To: innerState{number: to}, With: with}
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

func TestEdgeToTerminate(t *testing.T) {
	terminate := innerState{terminateVertex, true, false}
	input := []m.Edge{
		m.Edge{innerState{number: 0, terminate: true, start: true}, innerState{number: 1}, "0-1"},
		newEdge(1, 2, "1-2"),
		m.Edge{innerState{number: 2}, innerState{number: 3, terminate: true}, "2-3"},
	}

	out := []m.Edge{
		m.Edge{innerState{number: 0, terminate: false, start: true}, innerState{number: 1}, "0-1"},
		newEdge(1, 2, "1-2"),
		m.Edge{innerState{number: 2}, innerState{number: 3, terminate: false}, "2-3"},
		m.Edge{innerState{number: 0, terminate: false, start: true}, terminate, ""},
		m.Edge{innerState{number: 3, terminate: false}, terminate, ""},
	}

	fakeOut := []m.Edge{
		m.Edge{innerState{number: 0, terminate: true, start: true}, innerState{number: 1}, "0-1"},
		newEdge(1, 2, "1-2"),
		m.Edge{innerState{number: 2}, innerState{number: 3, terminate: true}, "2-3"},
		m.Edge{innerState{number: 0, terminate: false, start: true}, terminate, ""},
		m.Edge{innerState{number: 3, terminate: true}, terminate, ""},
	}

	ans := edgeToTerminate(input)

	if !Equals(out, ans) || Equals(fakeOut, ans) {
		t.Fail()
	}
}

func TestMakeRegex(t *testing.T) {
	m := newMachine(
		[]uint{0},
		[]uint{1},
		edge.Edge{0, 1, "a"},
		edge.Edge{1, 1, "b"},
	)

	out := "a((b))*"

	ans := MakeRegex(m)

	if ans != out {
		t.Fail()
	}
}
