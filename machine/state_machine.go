package machine

// State is an inner representation of machine's state
type State interface {
	Number() uint
	Start() bool
	Terminate() bool
}

// Edge is for transfers representation
type Edge struct {
	From State
	To State
	With string
}

// FiniteStateMachine is a finite automaton
type FiniteStateMachine interface {
	GoBackwardBy(to []State, with string) []State
    GoForwardBy(from []State, with string) []State
	OutgoingEdges(from []State) []Edge
	IngoingEdges(to []State) []Edge
	States() []State
}

// DeterminedStateMachine is a FiniteStateMachine which could be used to match strings
type DeterminedStateMachine interface {
	FiniteStateMachine
	Match(str string) bool
}

func SeparateStates(states []State) (start []uint, terminate []uint) {
	start = make([]uint, 0)
	terminate = make([]uint, 0)
	for _, state := range states {
		if state.Terminate() {
			terminate = append(terminate, state.Number())
		}
		if state.Start() {
			start = append(start, state.Number())
		}
	}
	return
}
