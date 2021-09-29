package tools

import (
	edge "state-machine/edge_machine"
	m "state-machine/machine"
)

func Invert(machine m.FinalStateMachine) m.FinalStateMachine {
	return invert(machine)
}

func invert(machine m.FinalStateMachine) *edge.Machine {
	edges := make([]edge.Edge, 0)
	start := make([]uint, 0)
	terminate := make([]uint, 0)
	for _, s := range machine.States() {
		if s.Start() {
			start = append(start, s.Number())
		}
		if !s.Terminate() {
			terminate = append(terminate, s.Number())
		}
		for _, e := range machine.OutgoingEdges([]m.State{s}) {
			edges = append(edges, *edge.NewEdgeMachine(e))
		}
	}
	return edge.NewMachine(edges, start, terminate)
}
