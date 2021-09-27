package remove_epsilon

import (
	edge "state-machine/edge_machine"
	"testing"
)

func NewMachine(transfers... edge.Edge) *edge.Machine {
	ans := make([]edge.Edge, 0, len(transfers))
	ans = append(ans, transfers...)
	return edge.NewMachine(ans)
}

func TestEpsilonRemoval(t *testing.T) {
	m := NewMachine(
		edge.Edge{edge.State{0, false}, edge.State{1, false}, ""},
		edge.Edge{edge.State{0, false}, edge.State{2, false}, ""},
		edge.Edge{edge.State{1, false}, edge.State{3, false}, "a"},
		edge.Edge{edge.State{2, false}, edge.State{3, false}, "b"},
	)

	out := NewMachine(
		edge.Edge{edge.State{0, false}, edge.State{3, false}, "a"},
		edge.Edge{edge.State{0, false}, edge.State{3, false}, "b"},
		edge.Edge{edge.State{1, false}, edge.State{3, false}, "a"},
		edge.Edge{edge.State{2, false}, edge.State{3, false}, "b"},
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
	m := edge.NewMachine(make([]edge.Edge, 0))
	ans, err := removeEpsilon(m)
	if err != nil {
		t.Error(err)
	}
	if !m.Equals(*ans) {
		t.Fail()
	}
}

func TestRemoveTerminate(t *testing.T) {
	m := NewMachine(
		edge.Edge{edge.State{0, false}, edge.State{1, false}, ""},
		edge.Edge{edge.State{1, false}, edge.State{2, false}, ""},
		edge.Edge{edge.State{2, false}, edge.State{3, true}, ""},
		edge.Edge{edge.State{3, true}, edge.State{4, true}, "a"},
	)

	ans, err := removeEpsilon(m)

	if err != nil {
		t.Error(err)
	}

	out := NewMachine(
		edge.Edge{edge.State{0, true}, edge.State{3, true}, ""},
		edge.Edge{edge.State{3, true}, edge.State{4, true}, "a"},
	)

	if !out.Equals(*ans) {
		t.Fail()
	}
}
