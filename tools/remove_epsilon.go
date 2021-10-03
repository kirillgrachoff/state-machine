package tools

import (
	edge "state-machine/edge_machine"
	m "state-machine/machine"
)

func RemoveEpsilon(machine m.FinalStateMachine) (m.FinalStateMachine, error) {
	ans, err := removeEpsilon(machine)
	if err != nil {
		return nil, err
	}
	return ans, nil
}

func removeEpsilon(machine m.FinalStateMachine) (*edge.Machine, error) {
	ans := make([]edge.Edge, 0)
	start := make([]uint, 0, 1)
	terminate := make([]uint, 0, 0)
	for _, state := range machine.States() {
		if state.Terminate() {
			terminate = append(terminate, state.Number())
		}
		if state.Start() {
			start = append(start, state.Number())
		}
		to := goByEmptyTransfers(machine, state)
		ans = append(ans, to...)
	}
	terminate = findTerminate(machine, terminate)
	return edge.BuildNewMachine(ans, start, terminate), nil
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
				finalEdges = append(finalEdges, edge.Edge{From: state.Number(), To: e.To.Number(), With: e.With})
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

// findTerminate finds all terminate vertices
func findTerminate(machine m.FinalStateMachine, terminate []uint) []uint {
	queue := make([]uint, 0, len(terminate))
	for _, vertex := range terminate {
		queue = append(queue, vertex)
	}
	used := make(map[uint]bool)
	for len(queue) > 0 {
		vertex := queue[0]
		queue = queue[1:]
		for _, edge := range machine.IngoingEdges([]m.State{&edge.State{Index: vertex}}) {
			if used[edge.From.Number()] {
				continue
			}
			used[edge.From.Number()] = true
			if edge.With != "" {
				continue
			}
			terminate = append(terminate, edge.From.Number())
		}
	}
	return terminate
}
