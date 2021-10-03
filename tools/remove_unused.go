package tools

import (
	"sort"
	"state-machine/edge_machine"
	m "state-machine/machine"
)

func RemoveUnused(machine m.FinalStateMachine) (m.FinalStateMachine, error) {
	return repaintVertices(removeUnused(machine)), nil
}

func repaintVertices(machine m.FinalStateMachine) *edge_machine.Machine {
	states := machine.States()
	sort.Slice(states, func(i, j int) bool {
		return states[i].Number() < states[j].Number()
	})
	index := make(map[uint]uint)
	for i, s := range states {
		index[s.Number()] = uint(i)
	}
	terminate := make([]uint, 0)
	start := make([]uint, 0)
	edges := make([]edge_machine.Edge, 0)
	for _, state := range states {
		if state.Start() {
			start = append(start, index[state.Number()])
		}
		if state.Terminate() {
			terminate = append(terminate, index[state.Number()])
		}
		for _, edge := range machine.OutgoingEdges([]m.State{state}) {
			edges = append(edges, edge_machine.Edge{
				From: index[edge.From.Number()],
				To: index[edge.To.Number()],
				With: edge.With,
			})
		}
	}
	return edge_machine.BuildNewMachine(edges, start, terminate)
}

func removeUnused(machine m.FinalStateMachine) *edge_machine.Machine {
	start := make([]uint, 0, 1)
	terminate := make([]uint, 0, 1)
	states := machine.States()
	for _, vertex := range states {
		if vertex.Start() {
			start = append(start, vertex.Number())
		}
		if vertex.Terminate() {
			terminate = append(terminate, vertex.Number())
		}
	}
	queue := make([]uint, 0, 1)
	used := make(map[uint]bool)
	for _, state := range start {
		queue = append(queue, state)
		used[state] = true
	}
	newEdges := make([]edge_machine.Edge, 0)
	for len(queue) > 0 {
		vertex := queue[0]
		queue = queue[1:]
		used[vertex] = true
		for _, edge := range machine.OutgoingEdges([]m.State{&edge_machine.State{Index: vertex}}) {
			newEdges = append(newEdges, *edge_machine.NewCanonicalEdge(edge))
			if !used[edge.To.Number()] {
				queue = append(queue, edge.To.Number())
			}
		}
	}
	return edge_machine.BuildNewMachine(newEdges, start, terminate)
}
