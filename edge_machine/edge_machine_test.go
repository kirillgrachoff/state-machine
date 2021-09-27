package edge_machine_test

import (
	e "state-machine/edge_machine"
	m "state-machine/machine"
	"testing"
)

func NewMachine(transfers... e.Edge) *e.Machine {
	ans := make([]e.Edge, 0, len(transfers))
	ans = append(ans, transfers...)
	return e.NewMachine(ans)
}

func Equals(a, b []m.State) bool {
	cnt := make(map[uint]int)
	for _, e := range a {
		cnt[e.Number()]++
	}
	for _, e := range b {
		cnt[e.Number()]--
	}
	for _, v := range cnt {
		if v != 0 {
			return false
		}
	}
	return true
}

func TestEmptyTransfers(t *testing.T) {
	m1 := NewMachine(
		e.Edge{From: e.State{Index: 0}, To: e.State{Index: 1}, With: ""},
		e.Edge{From: e.State{Index: 1}, To: e.State{Index: 2}, With: ""},
	)
	ans1 := []m.State{
		e.State{Index: 1},
	}
	if !Equals(m1.GoBy([]m.State{e.State{Index: 0}}, ""), ans1) {
		t.Fail()
	}
	ans2 := []m.State{
		e.State{Index: 2},
	}
	if !Equals(m1.GoBy([]m.State{e.State{Index: 1}}, ""), ans2) {
		t.Fail()
	}
}

func TestEquals(t *testing.T) {
	m1 := NewMachine(
		e.Edge{From: e.State{Index: 0}, To: e.State{Index: 1}, With: ""},
		e.Edge{From: e.State{Index: 1}, To: e.State{Index: 2}, With: ""},
	)
	m2 := NewMachine(
		e.Edge{From: e.State{Index: 1}, To: e.State{Index: 2}, With: ""},
		e.Edge{From: e.State{Index: 0}, To: e.State{Index: 1}, With: ""},
	)
	m3 := NewMachine(
		e.Edge{From: e.State{Index: 0}, To: e.State{Index: 1}, With: ""},
		e.Edge{From: e.State{Index: 1}, To: e.State{Index: 3}, With: ""},
	)
	m4 := NewMachine(
		e.Edge{From: e.State{Index: 0}, To: e.State{Index: 1}, With: ""},
		e.Edge{From: e.State{Index: 1}, To: e.State{Index: 2}, With: "a"},
	)
	if !(m1.Equals(*m2)) || m1.Equals(*m3) || m1.Equals(*m4) {
		t.Fail()
	}
}
