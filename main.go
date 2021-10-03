package main

import (
	"fmt"
	"log"
	"state-machine/edge_machine"
	m "state-machine/machine"
	"state-machine/tools"
)

type Job func(m.FinalStateMachine) (m.FinalStateMachine, error)

func NewMachine(start, terminate []uint, transfers ...edge_machine.Edge) m.FinalStateMachine {
	ans := make([]edge_machine.Edge, 0, len(transfers))
	ans = append(ans, transfers...)
	return edge_machine.BuildNewMachine(ans, start, terminate)
}

func main() {
	machine := NewMachine(
		[]uint{0, 1},
		[]uint{2, 3},
		edge_machine.Edge{0, 1, ""},
		edge_machine.Edge{1, 2, "a"},
		edge_machine.Edge{1, 2, "b"},
		edge_machine.Edge{1, 3, "b"},
	)

	var err error
	machine, err = PipelineSeq(
		machine,
		tools.RemoveEpsilon,
		tools.RemoveUnused,
		tools.Determine,
		tools.FullMachine("ab"),
		tools.RemoveUnused,
		tools.Invert,
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

func PipelineSeq(machine m.FinalStateMachine, seq ...Job) (m.FinalStateMachine, error) {
	args := make([]Job, 0)
	args = append(args, seq...)
	return Pipeline(machine, args)
}

func Pipeline(machine m.FinalStateMachine, seq []Job) (m.FinalStateMachine, error) {
	for _, command := range seq {
		var err error
		machine, err = command(machine)
		if err != nil {
			return nil, err
		}
	}
	return machine, nil
}
