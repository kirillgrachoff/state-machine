package edge_machine

import (
	m "state-machine/machine"
	"testing"
)

func newMachine(start, terminate []uint, transfers ...Edge) *Machine {
	ans := make([]Edge, 0, len(transfers))
	ans = append(ans, transfers...)
	return NewMachine(ans, start, terminate)
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
	m1 := newMachine(
		[]uint{0},
		[]uint{},
		Edge{From: 0, To: 1, With: ""},
		Edge{From: 1, To: 2, With: ""},
	)
	ans1 := []m.State{
		State{Index: 1},
	}
	if !Equals(m1.GoBy([]m.State{State{Index: 0}}, ""), ans1) {
		t.Fail()
	}
	ans2 := []m.State{
		State{Index: 2},
	}
	if !Equals(m1.GoBy([]m.State{State{Index: 1}}, ""), ans2) {
		t.Fail()
	}
}

func TestEquals(t *testing.T) {
	m1 := newMachine(
		[]uint{0},
		[]uint{},
		Edge{From: 0, To: 1, With: ""},
		Edge{From: 1, To: 2, With: ""},
	)
	m2 := newMachine(
		[]uint{0},
		[]uint{},
		Edge{From: 1, To: 2, With: ""},
		Edge{From: 0, To: 1, With: ""},
	)
	m3 := newMachine(
		[]uint{0},
		[]uint{},
		Edge{From: 0, To: 1, With: ""},
		Edge{From: 1, To: 3, With: ""},
	)
	m4 := newMachine(
		[]uint{0},
		[]uint{},
		Edge{From: 0, To: 1, With: ""},
		Edge{From: 1, To: 2, With: "a"},
	)
	if !(m1.Equals(*m2)) || m1.Equals(*m3) || m1.Equals(*m4) {
		t.Fail()
	}
}
