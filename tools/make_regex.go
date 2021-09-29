package tools

import (
	"fmt"
	"log"
	edge "state-machine/edge_machine"
	m "state-machine/machine"
)

func machineFromEdges(edges []m.Edge) m.FinalStateMachine {
	return edge.NewCanonicalMachine(edges)
}

func getEdges(machine m.FinalStateMachine) []m.Edge {
	ans := make([]m.Edge, 0)
	for _, s := range machine.States() {
		for _, e := range machine.OutgoingEdges([]m.State{s}) {
			ans = append(ans, e)
		}
	}
	return ans
}

func getIngoingLoopsOutgoing(machine m.FinalStateMachine, v m.State) (left, loop, right []m.Edge) {
	left = make([]m.Edge, 0)
	loop = make([]m.Edge, 0)
	right = make([]m.Edge, 0)
	for _, e := range machine.IngoingEdges([]m.State{v}) {
		if e.From.Number() != e.To.Number() {
			left = append(left, e)
		} else {
			loop = append(loop, e)
		}
	}
	for _, e := range machine.OutgoingEdges([]m.State{v}) {
		if e.From.Number() != e.To.Number() {
			right = append(right, e)
		}
	}
	return
}

func MakeRegex(machine m.FinalStateMachine) string {
	edges := getEdges(machine)
	edges = edgeToTerminate(edges)
	machine = machineFromEdges(edges)
	states := machine.States()
	for _, s := range states {
		if s.Terminate() {
			continue
		}
		if s.Start() {
			continue
		}
		left, loop, right := getIngoingLoopsOutgoing(machineFromEdges(edges), s)
		edges = deleteVertex(edges, s, left, loop, right)
	}
	ans := ""
	for _, e := range edges {
		ans += fmt.Sprintf("+(%v)", e.With)
	}
	if len(ans) > 0 {
		ans = ans[1:]
	}
	return ans
}

type innerState struct {
	Index uint
	End   bool
	Begin bool
}

func (i innerState) Number() uint {
	return i.Index
}

func (i innerState) Start() bool {
	return i.Begin
}

func (i innerState) Terminate() bool {
	return i.End
}

func edgeToTerminate(edges []m.Edge) (ans []m.Edge) {
	terminate := innerState{terminateVertex, true, false}
	ans = make([]m.Edge, 0)
	used := make(map[uint]bool)
	for _, e := range edges {
		if !used[e.To.Number()] && e.To.Terminate() {
			used[e.To.Number()] = true
			ans = append(ans, m.Edge{
				innerState{e.To.Number(), false, e.To.Start()},
				terminate,
				"",
			})
		}
		if !used[e.From.Number()] && e.From.Terminate() {
			used[e.To.Number()] = true
			ans = append(ans, m.Edge{
				innerState{e.From.Number(), false, e.From.Start()},
				terminate,
				"",
			})
		}
		ans = append(ans, m.Edge{
			From: innerState{
				e.From.Number(),
				false,
				e.From.Start(),
			},
			To: innerState{
				e.To.Number(),
				false,
				e.To.Start(),
			},
			With: e.With,
		})
	}
	return
}

func deleteVertex(edges []m.Edge, vertex m.State, ingoing, loops, outgoing []m.Edge) []m.Edge {
	var middleRegex string
	for _, e := range loops {
		if e.From.Number() != e.To.Number() {
			log.Fatalln("loop is not loop")
		}
		middleRegex = middleRegex + fmt.Sprintf("+(%v)", e.With)
	}
	if len(middleRegex) > 0 {
		middleRegex = middleRegex[1:]
	}
	if len(middleRegex) > 0 {
		middleRegex = fmt.Sprintf("(%v)*", middleRegex)
	}
	ans := make([]m.Edge, 0)
	for _, in := range ingoing {
		for _, out := range outgoing {
			if in.To.Number() != out.From.Number() {
				log.Fatalln("input error")
			}
			ans = append(ans, m.Edge{in.From, out.To, fmt.Sprintf("%v%v%v", in.With, middleRegex, out.With)})
		}
	}
	for _, e := range edges {
		if e.From.Number() == vertex.Number() || e.To.Number() == vertex.Number() {
			continue
		}
		ans = append(ans, e)
	}
	return ans
}
