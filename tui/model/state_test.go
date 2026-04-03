package model_test

import (
	"testing"

	"github.com/fabianoflorentino/stereoctl/tui/model"
)

func TestStateString(t *testing.T) {
	cases := []struct {
		state model.State
		want  string
	}{
		{model.StateFileSelection, "FileSelection"},
		{model.StateDiagnosis, "Diagnosis"},
		{model.StateProfileSelection, "ProfileSelection"},
		{model.StateConverting, "Converting"},
		{model.StateResult, "Result"},
	}
	for _, c := range cases {
		if got := c.state.String(); got != c.want {
			t.Errorf("State(%d).String() = %q, want %q", c.state, got, c.want)
		}
	}
}

func TestStateStringUnknown(t *testing.T) {
	s := model.State(99)
	if got := s.String(); got != "Unknown" {
		t.Errorf("expected Unknown, got %q", got)
	}
}

func TestStateIota(t *testing.T) {
	if model.StateFileSelection != 0 {
		t.Error("StateFileSelection must be 0")
	}
	if model.StateResult != 4 {
		t.Error("StateResult must be 4")
	}
}
