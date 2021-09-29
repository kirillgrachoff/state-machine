package tools

import (
	"fmt"
	"log"
	m "state-machine/machine"
)

//func MakeRegex(machine m.FinalStateMachine) string {
//
//}

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
