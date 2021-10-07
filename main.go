package main

import (
	"fmt"
	"log"
	"state-machine/edge_machine"
	m "state-machine/machine"
	"state-machine/tools"
)

type Job func(m.FiniteStateMachine) (m.FiniteStateMachine, error)

func NewMachine(start, terminate []uint, transfers ...edge_machine.Edge) m.FiniteStateMachine {
	ans := make([]edge_machine.Edge, 0, len(transfers))
	ans = append(ans, transfers...)
	return edge_machine.BuildNewMachine(ans, start, terminate)
}

func main() {
	machine := NewMachine(
		[]uint{0},
		[]uint{1},
		edge_machine.Edge{0, 0, "a"},
		edge_machine.Edge{0, 1, "b"},
	)

	var err error
	machine, err = PipelineSeq(
		machine,
		tools.RemoveEmptySymbols,
		tools.RemoveUnused,
		tools.Determine,
		tools.FullMachine("ab"),
		tools.RemoveUnused,
		//tools.Invert,
	)

	if err != nil {
		log.Fatalln(err)
	}

	ans, err := tools.MakeRegex(machine)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(ans)
}

func PipelineSeq(machine m.FiniteStateMachine, seq ...Job) (m.FiniteStateMachine, error) {
	args := make([]Job, 0)
	args = append(args, seq...)
	return Pipeline(machine, args)
}

func Pipeline(machine m.FiniteStateMachine, seq []Job) (m.FiniteStateMachine, error) {
	for _, command := range seq {
		var err error
		machine, err = command(machine)
		if err != nil {
			return nil, err
		}
	}
	return machine, nil
}
