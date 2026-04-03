package view

import (
	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/fabianoflorentino/stereoctl/internal/profiles"
)

// DefaultRenderer implements model.Renderer using the view package functions.
type DefaultRenderer struct{}

func (DefaultRenderer) FileSelection(fp filepicker.Model) string {
	return FileSelection(fp)
}

func (DefaultRenderer) Diagnosis(ev profiles.Evaluation, file string) string {
	return Diagnosis(ev, file)
}

func (DefaultRenderer) Profile() string {
	return Profile()
}

func (DefaultRenderer) Converting(bar progress.Model) string {
	return Converting(bar)
}

func (DefaultRenderer) Result(output string, err error) string {
	return Result(output, err)
}
