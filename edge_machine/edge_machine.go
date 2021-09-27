package edge_machine

import (
	"log"
	m "state-machine/machine"
)

// Machine is represented With its list of edges
type Machine struct {
	states []State
	mapping []Edge
}

func NewMachine(transfers []Edge, terminateVertices []uint) *Machine {
	machine := &Machine{
		states: make([]State, 0, 0),
		mapping: make([]Edge, len(transfers), len(transfers)),
	}
	copy(machine.mapping, transfers)
	cnt := make(map[uint]bool)
	for _, e := range transfers {
		cnt[e.From] = true
		cnt[e.To] = true
	}
	for i := uint(0); i < uint(len(cnt)); i++ {
		machine.states = append(machine.states, State{Index: i})
	}
	for _, t := range terminateVertices {
		machine.states[t].Finish = true
	}
	return machine
}

func (machine Machine) Equals(other Machine) bool {
	edgeCnt := make(map[Edge]uint)
	for _, e := range machine.mapping {
		edgeCnt[e]++
	}
	for _, e := range other.mapping {
		edgeCnt[e]++
	}
	for _, cnt := range edgeCnt {
		if cnt != 2 {
			return false
		}
	}
	return true
}

// goByRune is a procedure which receives a key To make transfer With
// it is guaranteed that its size not more than 1
func (machine Machine) goByRune(from []m.State, with string) []m.State {
	if len(with) > 1 {
		log.Fatalln("transfer size is > 1")
	}
	fromCnt := make(map[uint]bool)
	for _, f := range from {
		fromCnt[f.Number()] = true
	}

	ans := make([]m.State, 0, 0)
	for _, e := range machine.mapping {
		if fromCnt[e.From] && e.With == with {
			ans = append(ans, machine.states[e.To])
		}
	}
	return ans
}

// GoBy works incorrect if `with` == ""
// it does only one step for every rune
func (machine Machine) GoBy(from []m.State, with string) []m.State {
	if with == "" {
		return machine.goByRune(from, "")
	}
	ans := from
	for _, c := range with {
		ans = machine.goByRune(ans, string(c))
	}
	return ans
}

func (machine Machine) States() []m.State {
	ans := make([]m.State, 0, len(machine.states))
	for _, v := range machine.states {
		ans = append(ans, v)
	}
	return ans
}

func (machine Machine) OutgoingEdges(from []m.State) []m.Edge {
	ans := make([]m.Edge, 0)
	in := make(map[uint]bool)
	for _, from := range from {
		in[from.Number()] = true
	}
	for _, e := range machine.mapping {
		if in[e.From] {
			ans = append(ans, m.Edge{From: machine.states[e.From], To: machine.states[e.To], With: e.With})
		}
	}
	return ans
}

func (machine Machine) IngoingEdges(to []m.State) []m.Edge {
	ans := make([]m.Edge, 0)
	in := make(map[uint]bool)
	for _, to := range to {
		in[to.Number()] = true
	}
	for _, e := range machine.mapping {
		if in[e.To] {
			ans = append(ans, m.Edge{From: machine.states[e.From], To: machine.states[e.To], With: e.With})
		}
	}
	return ans
}

type Edge struct {
	From uint
	To   uint
	With string
}

func NewEdge(from State, to State, with string) *Edge {
	if len(with) > 1 {
		log.Fatalln("too big jump")
	}
	return &Edge{from.Number(), to.Number(), with}
}

func (e Edge) Equals(another Edge) bool {
	return e.From == another.From && e.To == another.To && e.With == another.With
}

type State struct {
	Index  uint
	Finish bool
}

func (s State) Equals(another State) bool {
	return s.Index == another.Index
}

func (s State) Number() uint {
	return s.Index
}

func (s State) Terminate() bool {
	return s.Finish
}
