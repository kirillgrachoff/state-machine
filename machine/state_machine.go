package machine

// State is an inner representation of machine's state
type State interface {
	Number() uint
	Terminate() bool
}

// Edge is for transfers representation
type Edge struct {
	From State
	To State
	With string
}

// FinalStateMachine is a finite automaton
type FinalStateMachine interface {
    GoBy(from []State, with string) []State
	OutgoingEdges(from []State) []Edge
	States() []State
}

// DeterminedStateMachine is a FinalStateMachine which could be used to match strings
type DeterminedStateMachine interface {
	FinalStateMachine
	Match(str string) bool
}

