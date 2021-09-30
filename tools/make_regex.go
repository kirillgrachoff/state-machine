package tools

import (
	"fmt"
	"log"
	"sort"
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

func filterVertexEdges(edges []m.Edge, v m.State) (left, loop, right, remain []m.Edge) {
	left = make([]m.Edge, 0)
	loop = make([]m.Edge, 0)
	right = make([]m.Edge, 0)
	remain = make([]m.Edge, 0)
	for _, e := range edges {
		if e.From.Number() == v.Number() && e.To.Number() == v.Number() {
			loop = append(loop, e)
		} else if e.From.Number() == v.Number() {
			right = append(right, e)
		} else if e.To.Number() == v.Number() {
			left = append(left, e)
		} else {
			remain = append(remain, e)
		}
	}
	return
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
	sort.Slice(states, func(i, j int) bool {
		return states[i].Number() < states[j].Number()
	})
	for _, s := range states {
		if s.Terminate() {
			continue
		}
		if s.Start() {
			continue
		}
		left, loop, right, edgesV := filterVertexEdges(edges, s)
		nEdges := deleteVertex(left, loop, right)
		edges = edgesV
		for _, v := range nEdges {
			edges = append(edges, v)
		}
	}
	ansCnt := make(map[string]struct{})
	for _, e := range edges {
		ansCnt[fmt.Sprintf("+(%v)", e.With)] = struct{}{}
	}
	ansArray := make([]string, 0)
	for s := range ansCnt {
		ansArray = append(ansArray, s)
	}
	sort.Slice(ansArray, func(i, j int) bool {
		return ansArray[i] < ansArray[j]
	})
	ans := ""
	for _, s := range ansArray {
		ans += s
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
	linkToTerminate := func(v m.State) {
		if used[v.Number()] || !v.Terminate() {
			return
		}
		used[v.Number()] = true
		ans = append(ans, m.Edge{
			innerState{v.Number(), false, v.Start()},
			terminate,
			"",
		})
	}
	for _, e := range edges {
		linkToTerminate(e.From)
		linkToTerminate(e.To)
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

func deleteVertex(ingoing, loops, outgoing []m.Edge) []m.Edge {
	middleRegexCnt := make(map[string]struct{})
	for _, e := range loops {
		if e.From.Number() != e.To.Number() {
			log.Fatalln("loop is not loop")
		}
		middleRegexCnt[e.With] = struct{}{}
	}
	var middleRegex string
	{
		all := make([]string, 0)
		for s := range middleRegexCnt {
			all = append(all, s)
		}
		sort.Slice(all, func(i, j int) bool {
			return all[i] < all[j]
		})
		for _, v := range all {
			middleRegex += fmt.Sprintf("+(%v)", v)
		}
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
	return ans
}
