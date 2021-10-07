package tools

import (
	edge "state-machine/edge_machine"
	"testing"
)

func TestRemoveUnused(t *testing.T) {
	m := newMachine(
		[]uint{0},
		[]uint{3},
		edge.Edge{From: 0, To: 3, With: "a"},
		edge.Edge{From: 0, To: 3, With: "b"},
		edge.Edge{From: 1, To: 3, With: "hello"},
	)

	out := newMachine(
		[]uint{0},
		[]uint{3},
		edge.Edge{From: 0, To: 3, With: "a"},
		edge.Edge{From: 0, To: 3, With: "b"},
	)

	ans := removeUnused(m)

	if !out.Equals(*ans) {
		t.Fail()
	}
}

func TestRepaint(t *testing.T) {
	m := newMachine(
		[]uint{0},
		[]uint{3},
		edge.Edge{From: 0, To: 3, With: "a"},
		edge.Edge{From: 0, To: 3, With: "b"},
	)

	out := newMachine(
		[]uint{0},
		[]uint{1},
		edge.Edge{From: 0, To: 1, With: "a"},
		edge.Edge{From: 0, To: 1, With: "b"},
	)

	ans := repaintVertices(m)

	if !out.Equals(*ans) {
		t.Fail()
	}
}
