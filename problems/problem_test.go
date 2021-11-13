package problems

import (
	"math"
	"state-machine/edge_machine"
	"testing"
)

func newMachine(start, terminate []uint, transfers ...edge_machine.Edge) *edge_machine.Machine {
	ans := make([]edge_machine.Edge, 0, len(transfers))
	ans = append(ans, transfers...)
	return edge_machine.BuildNewMachine(ans, start, terminate)
}

func assertEq(a, b int, t *testing.T) {
	if a != b {
		t.Fail()
	}
}

func TestProblemCaseAcceptable(t *testing.T) {
	m := newMachine(
		[]uint{0},
		[]uint{2},
		edge_machine.Edge{From: 0, To: 1, With: "a"},
		edge_machine.Edge{From: 1, To: 2, With: "b"},
	)
	ans, err := FindLongestSuffix(m, "b")
	if err != nil {
		t.Error(err)
	}
	assertEq(ans, 1, t)
}

func TestProblemCaseEdge(t *testing.T) {
	m := newMachine(
		[]uint{0},
		[]uint{1},
		edge_machine.Edge{From: 0, To: 1, With: "a"},
		edge_machine.Edge{From: 1, To: 1, With: "b"},
	)
	ans, err := FindLongestSuffix(m, "b")
	if err != nil {
		t.Error(err)
	}
	assertEq(ans, math.MaxInt, t)
}

func TestProblemCase1(t *testing.T) {
	m := newMachine(
		[]uint{0},
		[]uint{2},
		edge_machine.Edge{From: 0, To: 1, With: "a"},
		edge_machine.Edge{From: 1, To: 1, With: "b"},
		edge_machine.Edge{From: 1, To: 2, With: "a"},
	)
	ans, err := FindLongestSuffix(m, "a")
	if err != nil {
		t.Error(err)
	}
	assertEq(ans, 2, t)
}
