package tools

import (
	"state-machine/edge_machine"
	m "state-machine/machine"
	"state-machine/predefines"
)

func Minimize(machine m.FiniteStateMachine) (m.FiniteStateMachine, error) {
		return minimize(machine), nil
}

func minimize(machine m.FiniteStateMachine) *edge_machine.Machine {
	states := machine.States()
	equalityClass := make(map[uint]uint)
	for _, state := range states {
		if state.Terminate() {
			equalityClass[state.Number()] = 0
		} else {
			equalityClass[state.Number()] = 1
		}
	}
	oldEqualityClassesCount := 2
	newEdges := make([]transfersHolder, 0)
	for {
		transfers := make(map[transfersHolder][]uint)
		for _, state := range states {
			this := transfersHolder{
				VertexClass: equalityClass[state.Number()],
			}
			for i, symbol := range predefines.Alphabet {
				toNumber := machine.GoBy([]m.State{state}, string(symbol))[0].Number()
				this.Transfers[i] = equalityClass[toNumber]
			}
			transfers[this] = append(transfers[this], state.Number())
		}

		if oldEqualityClassesCount == len(transfers) {
			equalityClass = buildNewEqualityClasses(transfers)
			for top := range transfers {
				newEdges = append(newEdges, top)
			}
			break
		}

		oldEqualityClassesCount = len(transfers)

		equalityClass = buildNewEqualityClasses(transfers)
	}
	newStart := make([]uint, 0)
	newTerminate := make([]uint, 0)
	for _, state := range states {
		if state.Start() {
			newStart = append(newStart, equalityClass[state.Number()])
		}
		if state.Terminate() {
			newTerminate = append(newTerminate, equalityClass[state.Number()])
		}
	}
	newEdgesInMachine := make([]edge_machine.Edge, 0)
	for _, transfer := range newEdges {
		for symbolNumber, to := range transfer.Transfers {
			newEdgesInMachine = append(newEdgesInMachine, edge_machine.Edge{
				From: transfer.VertexClass,
				To: to,
				With: string(predefines.Alphabet[symbolNumber]),
			})
		}
	}
	return edge_machine.BuildNewMachine(newEdgesInMachine, newStart, newTerminate)
}

func buildNewEqualityClasses(transfers map[transfersHolder][]uint) map[uint]uint {
	newEqualityClass := make(map[uint]uint)
	classCounter := uint(0)
	for _, class := range transfers {
		for _, vertex := range class {
			newEqualityClass[vertex] = classCounter
		}
		classCounter++
	}
	return newEqualityClass
}

type transfersHolder struct {
	VertexClass uint
	Transfers [predefines.AlphabetLength]uint
}
