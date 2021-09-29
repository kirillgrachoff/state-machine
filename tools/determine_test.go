package tools

import (
	edge "state-machine/edge_machine"
	"testing"
)

func TestDetermine(t *testing.T) {
	m := newMachine(
		[]uint{0},
		[]uint{2},
		edge.Edge{0, 1, "a"},
		edge.Edge{1, 2, "a"},
		edge.Edge{0, 2, "a"},
	)

	out := newMachine(
		[]uint{0},
		[]uint{6, 4},
		edge.Edge{1, 6, "a"},
		edge.Edge{6, 4, "a"},
	)

	ans := determine(m)
	if !out.Equals(*ans) {
		t.Fail()
	}
}
