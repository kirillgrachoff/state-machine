package tools

import (
	edge "state-machine/edge_machine"
	"testing"
)

func TestInvert(t *testing.T) {
	m := newMachine(
		[]uint{0, 1, 5},
		[]uint{10, 20},
		edge.Edge{4, 5, "empty"},
	)

	out := newMachine(
		[]uint{0, 1, 5},
		[]uint{0, 1, 4, 5},
		edge.Edge{4, 5, "empty"},
	)

	ans := invert(m)

	if !out.Equals(*ans) {
		t.Fail()
	}
}
