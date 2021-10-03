package tools

import (
	edge "state-machine/edge_machine"
	"testing"
)

func newMachine(start, terminate []uint, transfers ...edge.Edge) *edge.Machine {
	ans := make([]edge.Edge, 0, len(transfers))
	ans = append(ans, transfers...)
	return edge.BuildNewMachine(ans, start, terminate)
}

func TestEpsilonRemoval(t *testing.T) {
	m := newMachine(
		[]uint{0},
		[]uint{},
		edge.Edge{0, 1, ""},
		edge.Edge{0, 2, ""},
		edge.Edge{1, 3, "a"},
		edge.Edge{2, 3, "b"},
	)

	out := newMachine(
		[]uint{0},
		[]uint{},
		edge.Edge{0, 3, "a"},
		edge.Edge{0, 3, "b"},
		edge.Edge{1, 3, "a"},
		edge.Edge{2, 3, "b"},
	)
	ans, err := removeEpsilon(m)
	if err != nil {
		t.Error(err)
	}
	if !out.Equals(*ans) {
		t.Fail()
	}
}

func TestEmptyGraph(t *testing.T) {
	m := newMachine([]uint{0}, []uint{})
	ans, err := removeEpsilon(m)
	if err != nil {
		t.Error(err)
	}
	if !m.Equals(*ans) {
		t.Fail()
	}
}

func TestRemoveTerminate(t *testing.T) {
	m := newMachine(
		[]uint{0},
		[]uint{3, 4},
		edge.Edge{0, 1, ""},
		edge.Edge{1, 2, ""},
		edge.Edge{2, 3, ""},
		edge.Edge{3, 4, "a"},
	)

	ans, err := removeEpsilon(m)

	if err != nil {
		t.Error(err)
	}

	out := newMachine(
		[]uint{0},
		[]uint{0, 1, 2, 3, 4},
		edge.Edge{0, 4, "a"},
		edge.Edge{1, 4, "a"},
		edge.Edge{2, 4, "a"},
		edge.Edge{3, 4, "a"},
	)

	if !out.Equals(*ans) {
		t.Fail()
	}
}
