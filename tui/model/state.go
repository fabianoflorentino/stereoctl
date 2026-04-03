package model

// State represents the current TUI screen.
type State int

const (
	StateFileSelection State = iota
	StateDiagnosis
	StateProfileSelection
	StateConverting
	StateResult
)

func (s State) String() string {
	names := [...]string{
		"FileSelection",
		"Diagnosis",
		"ProfileSelection",
		"Converting",
		"Result",
	}
	if int(s) < len(names) {
		return names[s]
	}
	return "Unknown"
}
