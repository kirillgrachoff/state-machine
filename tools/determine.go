package tools

import (
	edge "state-machine/edge_machine"
	m "state-machine/machine"
	"state-machine/mask_acceptor"
)

// Determine is a function which removes all ambiguous transfers
func Determine(machine m.FiniteStateMachine) (m.FiniteStateMachine, error) {
	return determine(machine), nil
}

func determine(machine m.FiniteStateMachine) *edge.Machine {
	states := machine.States()
	start, terminate := m.SeparateStates(states)

	newEdges := dfsWithMasks(machine, start)

	startSet := mask_acceptor.New(start)
	terminateSet := mask_acceptor.New(terminate)

	newStart := make([]uint, 0)
	newTerminate := make([]uint, 0)

	for _, edge := range newEdges {
		if startSet.In(edge.From) {
			newStart = append(newStart, edge.From)
		}
		if startSet.In(edge.To) {
			newStart = append(newStart, edge.To)
		}
		if terminateSet.In(edge.From) {
			newTerminate = append(newTerminate, edge.From)
		}
		if terminateSet.In(edge.To) {
			newTerminate = append(newTerminate, edge.To)
		}
	}
	return edge.BuildNewMachine(newEdges, newStart, newTerminate)
}

func dfsWithMasks(machine m.FiniteStateMachine, start []uint) []edge.Edge {
	newEdges := make([]edge.Edge, 0)
	queue := make([][]uint, 0)
	queue = append(queue, start)
	used := make(map[uint]bool)
	used[mask_acceptor.ToMask(start)] = true
	for len(queue) > 0 {
		vertex := queue[0]
		queue = queue[1:]
		alphabet := make(map[string]struct{})
		vStates := toStates(vertex)
		vertexMask := mask_acceptor.ToMask(vertex)
		for _, e := range machine.OutgoingEdges(vStates) {
			alphabet[e.With] = struct{}{}
		}
		for symbol := range alphabet {
			to := machine.GoBy(vStates, symbol)
			if len(to) == 0 {
				continue
			}
			toUint := toUints(to)
			toMask := mask_acceptor.ToMask(toUint)
			newEdges = append(newEdges, edge.Edge{
				From: vertexMask,
				To:   toMask,
				With: symbol,
			})
			if used[toMask] {
				continue
			}
			used[toMask] = true
			queue = append(queue, toUint)
		}
	}
	return newEdges
}

func toStates(args []uint) []m.State {
	ans := make([]m.State, 0)
	for _, vertex := range args {
		ans = append(ans, innerState{Index: vertex})
	}
	return ans
}

func toUints(args []m.State) []uint {
	ans := make([]uint, 0)
	for _, state := range args {
		ans = append(ans, state.Number())
	}
	return ans
}
