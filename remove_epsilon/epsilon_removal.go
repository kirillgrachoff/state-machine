package remove_epsilon

import (
	edge "state-machine/edge_machine"
	m "state-machine/machine"
)

func RemoveEpsilon(machine m.FinalStateMachine) (m.FinalStateMachine, error) {
	return removeEpsilon(machine)
}

func removeEpsilon(machine m.FinalStateMachine) (*edge.Machine, error) {
	ans := make([]edge.Edge, 0)
	for _, state := range machine.States() {
		to := goByEmptyTransfers(machine, state)
		ans = append(ans, to...)
	}
	return edge.NewMachine(ans), nil
}

// goByEmptyTransfers transforms all empty-string transfers from the state
func goByEmptyTransfers(machine m.FinalStateMachine, state m.State) []edge.Edge {

	finalEdges := make([]edge.Edge, 0)
	used := make(map[uint]bool)

	queue := make([]m.State, 0, 1)
	queue = append(queue, state)
	used[state.Number()] = true

	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for _, e := range machine.OutgoingEdges([]m.State{v}) {
			if e.With != "" {
				finalEdges = append(finalEdges, edge.Edge{edge.State{state.Number(), state.Terminate()}, edge.State{e.To.Number(), e.To.Terminate()}, e.With})
				continue
			}
			if used[e.To.Number()] {
				continue
			}
			queue = append(queue, e.To)
			used[e.To.Number()] = true
		}
	}
	return finalEdges
}
