package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fabianoflorentino/stereoctl/tui/model"
	"github.com/fabianoflorentino/stereoctl/tui/service"
	"github.com/fabianoflorentino/stereoctl/tui/view"
)

// Run starts the TUI application.
func Run() error {
	m := model.New(service.FFmpegProber{}, service.FFmpegConverter{}, view.DefaultRenderer{})
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("tui: %w", err)
	}
	return nil
}
