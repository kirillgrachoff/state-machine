package tools

import (
	"fmt"
	"sort"
	"state-machine/edge_machine"
	m "state-machine/machine"
)

func machineFromEdges(edges []m.Edge) m.FinalStateMachine {
	return edge_machine.NewCanonicalMachine(edges)
}

func getEdges(machine m.FinalStateMachine) []m.Edge {
	ans := make([]m.Edge, 0)
	for _, state := range machine.States() {
		for _, edge := range machine.OutgoingEdges([]m.State{state}) {
			ans = append(ans, edge)
		}
	}
	return ans
}

func filterVertexEdges(edges []m.Edge, vertex m.State) (ingoing, loop, outgoing, remain []m.Edge) {
	ingoing = make([]m.Edge, 0)
	loop = make([]m.Edge, 0)
	outgoing = make([]m.Edge, 0)
	// remain is for that edges which are not connected with vertex
	remain = make([]m.Edge, 0)
	for _, edge := range edges {
		if edge.From.Number() == vertex.Number() && edge.To.Number() == vertex.Number() {
			loop = append(loop, edge)
		} else if edge.From.Number() == vertex.Number() {
			outgoing = append(outgoing, edge)
		} else if edge.To.Number() == vertex.Number() {
			ingoing = append(ingoing, edge)
		} else {
			remain = append(remain, edge)
		}
	}
	return
}

func MakeRegex(machine m.FinalStateMachine) (ansRegex string, err error) {
	defer func() {
		if r := recover(); r != nil {
			ansRegex = ""
			err = fmt.Errorf("%v", r)
		}
	}()
	edges := getEdges(machine)
	edges = edgeToTerminate(edges)
	machine = machineFromEdges(edges)
	states := machine.States()
	sort.Slice(states, func(i, j int) bool {
		return states[i].Number() < states[j].Number()
	})
	for _, state := range states {
		if state.Terminate() {
			continue
		}
		if state.Start() {
			continue
		}
		left, loop, right, edgesV := filterVertexEdges(edges, state)
		newEdges := deleteVertex(left, loop, right)
		edges = edgesV
		for _, vertex := range newEdges {
			edges = append(edges, vertex)
		}
	}
	ansRegexMap := make(map[string]struct{})
	for _, edge := range edges {
		ansRegexMap[fmt.Sprintf("+(%v)", edge.With)] = struct{}{}
	}
	ansRegexArray := make([]string, 0)
	for state := range ansRegexMap {
		ansRegexArray = append(ansRegexArray, state)
	}
	sort.Slice(ansRegexArray, func(i, j int) bool {
		return ansRegexArray[i] < ansRegexArray[j]
	})
	ansRegex = ""
	for _, s := range ansRegexArray {
		ansRegex += s
	}
	if len(ansRegex) > 0 {
		ansRegex = ansRegex[1:]
	}
	return ansRegex, nil
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
	for _, edge := range edges {
		linkToTerminate(edge.From)
		linkToTerminate(edge.To)
		ans = append(ans, m.Edge{
			From: innerState{
				edge.From.Number(),
				false,
				edge.From.Start(),
			},
			To: innerState{
				edge.To.Number(),
				false,
				edge.To.Start(),
			},
			With: edge.With,
		})
	}
	return
}

func deleteVertex(ingoing, loops, outgoing []m.Edge) []m.Edge {
	middleRegexCnt := make(map[string]struct{})
	for _, edge := range loops {
		if edge.From.Number() != edge.To.Number() {
			panic("loop is not a loop")
		}
		middleRegexCnt[edge.With] = struct{}{}
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
				panic("input error")
			}
			ans = append(ans, m.Edge{
				in.From,
				out.To,
				fmt.Sprintf("%v%v%v", in.With, middleRegex, out.With),
			})
		}
	}
	return ans
}
