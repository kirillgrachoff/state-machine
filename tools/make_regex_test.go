package tools

import (
	edge "state-machine/edge_machine"
	m "state-machine/machine"
	"testing"
)

func newEdge(from, to uint, with string) m.Edge {
	return m.Edge{From: innerState{Index: from}, To: innerState{Index: to}, With: with}
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
		m.Edge{innerState{Index: 0, End: true, Begin: true}, innerState{Index: 1}, "0-1"},
		newEdge(1, 2, "1-2"),
		m.Edge{innerState{Index: 2}, innerState{Index: 3, End: true}, "2-3"},
	}

	out := []m.Edge{
		m.Edge{innerState{Index: 0, End: false, Begin: true}, innerState{Index: 1}, "0-1"},
		newEdge(1, 2, "1-2"),
		m.Edge{innerState{Index: 2}, innerState{Index: 3, End: false}, "2-3"},
		m.Edge{innerState{Index: 0, End: false, Begin: true}, terminate, ""},
		m.Edge{innerState{Index: 3, End: false}, terminate, ""},
	}

	fakeOut := []m.Edge{
		m.Edge{innerState{Index: 0, End: true, Begin: true}, innerState{Index: 1}, "0-1"},
		newEdge(1, 2, "1-2"),
		m.Edge{innerState{Index: 2}, innerState{Index: 3, End: true}, "2-3"},
		m.Edge{innerState{Index: 0, End: false, Begin: true}, terminate, ""},
		m.Edge{innerState{Index: 3, End: true}, terminate, ""},
	}

	ans := edgeToTerminate(input)

	if !Equals(out, ans) || Equals(fakeOut, ans) {
		t.Fail()
	}
}

func TestMakeRegex(t *testing.T) {
	{
		m := newMachine(
			[]uint{0},
			[]uint{1},
			edge.Edge{0, 1, "a"},
			edge.Edge{1, 1, "b"},
		)

		out := "(a((b))*)"

		ans, err := MakeRegex(m)
		if err != nil {
			t.Error(err)
		}

		if ans != out {
			t.Fail()
		}
	}
	{
		m := newMachine(
			[]uint{0, 1},
			[]uint{1, 2},
			edge.Edge{0, 1, "a"},
			edge.Edge{1, 2, "b"},
		)

		out := "()+(a)+(ab)+(b)"

		ans, err := MakeRegex(m)
		if err != nil {
			t.Error(err)
		}

		if ans != out {
			t.Fail()
		}
	}
}

func BenchmarkMakeRegex(b *testing.B) {
	m := newMachine(
		[]uint{0, 1},
		[]uint{1, 2},
		edge.Edge{0, 1, "a"},
		edge.Edge{1, 2, "b"},
	)

	out := "()+(a)+(ab)+(b)"
	for testNo := 0; testNo < b.N; testNo++ {
		ans, err := MakeRegex(m)
		if err != nil {
			b.Error(err)
		}
		if ans != out {
			b.Fail()
		}
	}
}
