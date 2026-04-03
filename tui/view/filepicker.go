package view

import (
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
)

// FileSelection renders the file picker screen.
func FileSelection(fp filepicker.Model) string {
	var sb strings.Builder
	sb.WriteString(titleStyle.Render("Select a file") + "\n\n")
	sb.WriteString(fp.View())
	sb.WriteString("\n" + dimStyle.Render("[q] quit"))
	return sb.String()
}
