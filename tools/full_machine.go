package tools

import (
	"math"
	edge "state-machine/edge_machine"
	m "state-machine/machine"
)

const terminateVertex = uint(math.MaxUint)

func FullMachine(alphabet string) func(machine m.FinalStateMachine) m.FinalStateMachine {
	return func(machine m.FinalStateMachine) m.FinalStateMachine {
		return fullMachine(alphabet)(machine)
	}
}

func fullMachine(alphabet string) func(machine m.FinalStateMachine) *edge.Machine {
	return func(machine m.FinalStateMachine) *edge.Machine {
		start := make([]uint, 0)
		terminate := make([]uint, 0)
		for _, s := range machine.States() {
			if s.Start() {
				start = append(start, s.Number())
			}
			if s.Terminate() {
				terminate = append(terminate, s.Number())
			}
		}
		edges := make([]edge.Edge, 0)
		used := make(map[uint]bool)
		for _, s := range machine.States() {
			used[s.Number()] = true
			hasEdges := make(map[string]bool)
			for _, e := range machine.OutgoingEdges([]m.State{s}) {
				hasEdges[e.With] = true
				edges = append(edges, *edge.NewEdgeMachine(e))
			}
			for _, c := range alphabet {
				if hasEdges[string(c)] {
					continue
				}
				edges = append(edges, edge.Edge{From: s.Number(), To: terminateVertex, With: string(c)})
			}
		}
		for _, c := range alphabet {
			edges = append(edges, edge.Edge{From: terminateVertex, To: terminateVertex, With: string(c)})
		}
		return edge.NewMachine(edges, start, terminate)
	}
}
