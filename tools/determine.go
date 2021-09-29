package tools

import (
	"sort"
	edge "state-machine/edge_machine"
	m "state-machine/machine"
	"unsafe"
)

// Determine is a function which removes all ambiguous transfers
func Determine(machine m.FinalStateMachine) m.FinalStateMachine {
	return determine(machine)
}

func determine(machine m.FinalStateMachine) *edge.Machine {
	states := machine.States()
	start := make([]uint, 0)
	terminate := make([]uint, 0)
	for _, s := range states {
		if s.Terminate() {
			terminate = append(terminate, s.Number())
		}
		if s.Start() {
			start = append(start, s.Number())
		}
	}
	newEdges := make([]edge.Edge, 0)
	queue := make([][]uint, 0)
	queue = append(queue, start)
	used := make(map[uint]bool)
	used[mask(start)] = true
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		alphabet := make(map[string]struct{})
		vStates := toStates(v)
		vMask := mask(v)
		for _, e := range machine.OutgoingEdges(vStates) {
			alphabet[e.With] = struct{}{}
		}
		for w := range alphabet {
			to := machine.GoBy(vStates, w)
			if len(to) == 0 {
				continue
			}
			toUint := toUints(to)
			toMask := mask(toUint)
			newEdges = append(newEdges, edge.Edge{From: vMask, To: toMask, With: w})
			if used[toMask] {
				continue
			}
			used[toMask] = true
			queue = append(queue, toUint)
		}
	}
	startMap := make(map[uint]bool)
	for _, v := range start {
		startMap[v] = true
	}
	terminateMap := make(map[uint]bool)
	for _, v := range start {
		terminateMap[v] = true
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

	for _, e := range newEdges {
		if isStartMask(e.From) {
			newStart = append(newStart, e.From)
		}
		if isStartMask(e.To) {
			newStart = append(newStart, e.To)
		}
		if isTerminateMask(e.From) {
			newTerminate = append(newTerminate, e.From)
		}
		if isTerminateMask(e.To) {
			newTerminate = append(newTerminate, e.To)
		}
	}
	return edge.NewMachine(newEdges, newStart, newTerminate)
}

func toStates(args []uint) []m.State {
	ans := make([]m.State, 0)
	for _, v := range args {
		ans = append(ans, innerState{Index: v})
	}
	return ans
}

func toUints(args []m.State) []uint {
	ans := make([]uint, 0)
	for _, s := range args {
		ans = append(ans, s.Number())
	}
	return ans
}

func sorted(args []uint) []uint {
	ans := make([]uint, len(args))
	copy(ans, args)
	sort.Slice(ans, func(i, j int) bool {
		return ans[i] < ans[j]
	})
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
