package tools

import (
	"sort"
	m "state-machine/machine"
)
import edge "state-machine/edge_machine"

func RemoveUnused(machine m.FinalStateMachine) m.FinalStateMachine {
	return repaintVertices(removeUnused(machine))
}

func repaintVertices(machine m.FinalStateMachine) *edge.Machine {
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
	edges := make([]edge.Edge, 0)
	for _, s := range states {
		if s.Start() {
			start = append(start, index[s.Number()])
		}
		if s.Terminate() {
			terminate = append(terminate, index[s.Number()])
		}
		for _, e := range machine.OutgoingEdges([]m.State{s}) {
			edges = append(edges, edge.Edge{From: index[e.From.Number()], To: index[e.To.Number()], With: e.With})
		}
	}
	return edge.NewMachine(edges, start, terminate)
}

func removeUnused(machine m.FinalStateMachine) *edge.Machine {
	start := make([]uint, 0, 1)
	terminate := make([]uint, 0, 1)
	states := machine.States()
	for _, v := range states {
		if v.Start() {
			start = append(start, v.Number())
		}
		if v.Terminate() {
			terminate = append(terminate, v.Number())
		}
	}
	queue := make([]uint, 0, 1)
	used := make(map[uint]bool)
	for _, s := range start {
		queue = append(queue, s)
		used[s] = true
	}
	newEdges := make([]edge.Edge, 0)
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		used[v] = true
		for _, e := range machine.OutgoingEdges([]m.State{&edge.State{Index: v}}) {
			newEdges = append(newEdges, *edge.NewEdgeMachine(e))
			if !used[e.To.Number()] {
				queue = append(queue, e.To.Number())
			}
		}
	}
	return edge.NewMachine(newEdges, start, terminate)
}
