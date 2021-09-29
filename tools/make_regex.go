package tools

import (
	"fmt"
	"log"
	m "state-machine/machine"
)

//func MakeRegex(machine m.FinalStateMachine) string {
//
//}

type innerState struct {
	number uint
	terminate bool
	start bool
}

func (i innerState) Number() uint {
	return i.number
}

func (i innerState) Start() bool {
	return i.start
}

func (i innerState) Terminate() bool {
	return i.terminate
}

func edgeToTerminate(edges []m.Edge) (ans []m.Edge, newTerminate []uint) {
	terminate := innerState{terminateVertex, true, false}
	ans = make([]m.Edge, 0)
	newTerminate = []uint{terminate.number}
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

func deleteVertex(ingoing, loops, outgoing []m.Edge) []m.Edge {
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
	ans := make([]m.Edge, 0)
	for _, in := range ingoing {
		for _, out := range outgoing {
			if in.To.Number() != out.From.Number() {
				log.Fatalln("input error")
			}
			ans = append(ans, m.Edge{in.From, out.To, fmt.Sprintf("%v(%v)*%v", in.With, middleRegex, out.With)})
		}
	}
	return ans
}
