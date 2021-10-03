package edge_machine

import (
	m "state-machine/machine"
)

// Machine is represented With its list of edges
type Machine struct {
	states  map[uint]State
	mapping []Edge
}

func NewCanonicalMachine(transfers []m.Edge) *Machine {
	edges := make([]Edge, 0)
	terminate := make([]uint, 0)
	start := make([]uint, 0)
	insert := func(state m.State) {
		if state.Terminate() {
			terminate = append(terminate, state.Number())
		}
		if state.Start() {
			start = append(start, state.Number())
		}
	}
	for _, edge := range transfers {
		insert(edge.To)
		insert(edge.From)
		edges = append(edges, Edge{edge.From.Number(), edge.To.Number(), edge.With})
	}
	return BuildNewMachine(edges, start, terminate)
}

func BuildNewMachine(transfers []Edge, startVertices, terminateVertices []uint) *Machine {
	machine := &Machine{
		states:  make(map[uint]State),
		mapping: make([]Edge, len(transfers), len(transfers)),
	}
	copy(machine.mapping, transfers)
	const (
		isTerminate = 1 << iota
		isStart
		isUsed
	)
	status := make(map[uint]uint)
	for _, edge := range transfers {
		status[edge.From] |= isUsed
		status[edge.To] |= isUsed
	}
	for _, startVertex := range startVertices {
		status[startVertex] |= isStart
	}
	for _, terminateVertex := range terminateVertices {
		status[terminateVertex] |= isTerminate
	}
	for number, summary := range status {
		machine.states[number] = State{
			Index: number,
			End: (summary & isTerminate) != 0,
			Begin: (summary & isStart) != 0,
		}
	}
	return machine
}

func (machine Machine) Equals(other Machine) bool {
	const usedTwice = 2
	edgeCnt := make(map[Edge]uint)
	for _, edge := range machine.mapping {
		edgeCnt[edge]++
	}
	for _, edge := range other.mapping {
		edgeCnt[edge]++
	}
	for _, cnt := range edgeCnt {
		if cnt != usedTwice {
			return false
		}
	}
	return true
}

// goByRune is a procedure which receives a key To make transfer With
// it is guaranteed that its size not more than 1
func (machine Machine) goByRune(from []m.State, with string) []m.State {
	if len(with) > 1 {
		panic("transfer size is > 1")
	}
	fromCnt := make(map[uint]bool)
	for _, f := range from {
		fromCnt[f.Number()] = true
	}

	ans := make([]m.State, 0, 0)
	for _, edge := range machine.mapping {
		if fromCnt[edge.From] && edge.With == with {
			ans = append(ans, machine.states[edge.To])
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
	for _, symbol := range with {
		ans = machine.goByRune(ans, string(symbol))
	}
	return ans
}

func (machine Machine) States() []m.State {
	ans := make([]m.State, 0, len(machine.states))
	for _, state := range machine.states {
		ans = append(ans, state)
	}
	return ans
}

func (machine Machine) OutgoingEdges(from []m.State) []m.Edge {
	ans := make([]m.Edge, 0)
	in := make(map[uint]bool)
	for _, from := range from {
		in[from.Number()] = true
	}
	for _, edge := range machine.mapping {
		if in[edge.From] {
			ans = append(ans, m.Edge{
				From: machine.states[edge.From],
				To: machine.states[edge.To],
				With: edge.With,
			})
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
	for _, edge := range machine.mapping {
		if in[edge.To] {
			ans = append(ans, m.Edge{
				From: machine.states[edge.From],
				To: machine.states[edge.To],
				With: edge.With,
			})
		}
	}
	return ans
}

type Edge struct {
	From uint
	To   uint
	With string
}

func NewCanonicalEdge(e m.Edge) *Edge {
	return &Edge{From: e.From.Number(), To: e.To.Number(), With: e.With}
}

func NewEdge(from State, to State, with string) *Edge {
	if len(with) > 1 {
		panic("too big jump")
	}
	return &Edge{from.Number(), to.Number(), with}
}

func (e Edge) Equals(another Edge) bool {
	return e.From == another.From && e.To == another.To && e.With == another.With
}

type State struct {
	Index uint
	End   bool
	Begin bool
}

func (s State) Equals(another State) bool {
	return s.Index == another.Index
}

func (s State) Number() uint {
	return s.Index
}

func (s State) Terminate() bool {
	return s.End
}

func (s State) Start() bool {
	return s.Begin
}
