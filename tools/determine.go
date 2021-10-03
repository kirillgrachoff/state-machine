package tools

import (
	edge "state-machine/edge_machine"
	m "state-machine/machine"
	"unsafe"
)

// Determine is a function which removes all ambiguous transfers
func Determine(machine m.FinalStateMachine) (m.FinalStateMachine, error) {
	return determine(machine), nil
}

func determine(machine m.FinalStateMachine) *edge.Machine {
	states := machine.States()
	start := make([]uint, 0)
	terminate := make([]uint, 0)
	for _, state := range states {
		if state.Terminate() {
			terminate = append(terminate, state.Number())
		}
		if state.Start() {
			start = append(start, state.Number())
		}
	}
	newEdges := make([]edge.Edge, 0)
	queue := make([][]uint, 0)
	queue = append(queue, start)
	used := make(map[uint]bool)
	used[mask(start)] = true
	for len(queue) > 0 {
		vertex := queue[0]
		queue = queue[1:]
		alphabet := make(map[string]struct{})
		vStates := toStates(vertex)
		vertexMask := mask(vertex)
		for _, e := range machine.OutgoingEdges(vStates) {
			alphabet[e.With] = struct{}{}
		}
		for symbol := range alphabet {
			to := machine.GoBy(vStates, symbol)
			if len(to) == 0 {
				continue
			}
			toUint := toUints(to)
			toMask := mask(toUint)
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
	startMap := make(map[uint]bool)
	for _, state := range start {
		startMap[state] = true
	}
	terminateMap := make(map[uint]bool)
	for _, state := range terminate {
		terminateMap[state] = true
	}
	newStart := make([]uint, 0)
	newTerminate := make([]uint, 0)

	isStartMask := func(m uint) bool {
		for _, i := range fromMask(m) {
			if startMap[i] {
				return true
			}
		}
		return false
	}

	isTerminateMask := func(m uint) bool {
		for _, i := range fromMask(m) {
			if terminateMap[i] {
				return true
			}
		}
		return false
	}

	for _, edge := range newEdges {
		if isStartMask(edge.From) {
			newStart = append(newStart, edge.From)
		}
		if isStartMask(edge.To) {
			newStart = append(newStart, edge.To)
		}
		if isTerminateMask(edge.From) {
			newTerminate = append(newTerminate, edge.From)
		}
		if isTerminateMask(edge.To) {
			newTerminate = append(newTerminate, edge.To)
		}
	}
	return edge.BuildNewMachine(newEdges, newStart, newTerminate)
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

func mask(args []uint) uint {
	ans := uint(0)
	for _, v := range args {
		ans |= 1 << v
	}
	return ans
}

func fromMask(mask uint) []uint {
	ans := make([]uint, 0)
	for i := uint(0); i < 8*uint(unsafe.Sizeof(mask)); i++ {
		if (mask & (1 << i)) != 0 {
			ans = append(ans, i)
		}
	}
	return ans
}
