package problems

import (
	"errors"
	"math"
	"state-machine/machine"
	"state-machine/tools"
)

func FindLongestSuffix(stateMachine machine.FiniteStateMachine, character string) (int, error) {
	if len(character) != 1 {
		return 0, errors.New("character size not equals to 1")
	}

	var err error

	stateMachine, err = tools.RemoveEmptySymbols(stateMachine)
	if err != nil {
		return 0, err
	}

	stateMachine, err = tools.Determine(stateMachine)
	if err != nil {
		return 0, err
	}

	states := stateMachine.States()
	current := make([]machine.State, 0)
	for _, state := range states {
		if !state.Terminate() {
			continue
		}
		current = append(current, state)
	}

	used := make(map[uint]bool)
	for _, state := range current {
		used[state.Number()] = true
	}

	answer := 0

	for {
		current = stateMachine.GoBackwardBy(current, character)
		for _, state := range current {
			if used[state.Number()] {
				return math.MaxInt, nil
			}
		}
		for _, state := range current {
			used[state.Number()] = true
		}
		if len(current) > 0 {
			answer++
		} else {
			break
		}
	}

	return answer, nil
}
