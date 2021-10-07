package tools

import (
	"state-machine/edge_machine"
	m "state-machine/machine"
)

func Invert(machine m.FiniteStateMachine) (m.FiniteStateMachine, error) {
	return invert(machine), nil
}

func invert(machine m.FiniteStateMachine) *edge_machine.Machine {
	edges := make([]edge_machine.Edge, 0)
	start := make([]uint, 0)
	terminate := make([]uint, 0)
	for _, state := range machine.States() {
		if state.Start() {
			start = append(start, state.Number())
		}
		if !state.Terminate() {
			terminate = append(terminate, state.Number())
		}
		for _, edge := range machine.OutgoingEdges([]m.State{state}) {
			edges = append(edges, *edge_machine.NewCanonicalEdge(edge))
		}
	}
	return edge_machine.BuildNewMachine(edges, start, terminate)
}
