package tools

import (
	"math"
	edge "state-machine/edge_machine"
	m "state-machine/machine"
)

const terminateVertex = uint(math.MaxUint)

func FullMachine(alphabet string) func(machine m.FinalStateMachine) (m.FinalStateMachine, error) {
	return func(machine m.FinalStateMachine) (m.FinalStateMachine, error) {
		return fullMachine(alphabet)(machine), nil
	}
}

func fullMachine(alphabet string) func(machine m.FinalStateMachine) *edge.Machine {
	return func(machine m.FinalStateMachine) *edge.Machine {
		start := make([]uint, 0)
		terminate := make([]uint, 0)
		for _, state := range machine.States() {
			if state.Start() {
				start = append(start, state.Number())
			}
			if state.Terminate() {
				terminate = append(terminate, state.Number())
			}
		}
		edges := make([]edge.Edge, 0)
		used := make(map[uint]bool)
		for _, state := range machine.States() {
			used[state.Number()] = true
			hasEdges := make(map[string]bool)
			for _, e := range machine.OutgoingEdges([]m.State{state}) {
				hasEdges[e.With] = true
				edges = append(edges, *edge.NewCanonicalEdge(e))
			}
			for _, symbol := range alphabet {
				if hasEdges[string(symbol)] {
					continue
				}
				edges = append(edges, edge.Edge{
					From: state.Number(),
					To: terminateVertex,
					With: string(symbol),
				})
			}
		}
		for _, symbol := range alphabet {
			edges = append(edges, edge.Edge{
				From: terminateVertex,
				To: terminateVertex,
				With: string(symbol),
			})
		}
		return edge.BuildNewMachine(edges, start, terminate)
	}
}
