package tools

import (
	edge "state-machine/edge_machine"
	"testing"
)

func TestFullMachine(t *testing.T) {
	alphabet := "ab"
	m := newMachine(
		[]uint{0},
		[]uint{1},
		edge.Edge{0, 1, "a"},
	)

	v := terminateVertex

	out := newMachine(
		[]uint{0},
		[]uint{1},
		edge.Edge{0, 1, "a"},
		edge.Edge{0, v, "b"},
		edge.Edge{1, v, "a"},
		edge.Edge{1, v, "b"},
		edge.Edge{v, v, "a"},
		edge.Edge{v, v, "b"},
	)

	ans := fullMachine(alphabet)(m)

	if !out.Equals(*ans) {
		t.Fail()
	}
}
