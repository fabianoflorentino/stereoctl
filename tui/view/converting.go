package view

import "github.com/charmbracelet/bubbles/progress"

// Converting renders the conversion progress screen.
func Converting(bar progress.Model) string {
	return titleStyle.Render("Converting...") + "\n\n" +
		bar.View() + "\n\n" +
		dimStyle.Render("[q] quit")
}
