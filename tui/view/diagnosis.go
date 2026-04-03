package view

import (
	"fmt"
	"strings"

	"github.com/fabianoflorentino/stereoctl/internal/profiles"
)

// Diagnosis renders the diagnosis screen.
func Diagnosis(ev profiles.Evaluation, file string) string {
	var sb strings.Builder
	sb.WriteString(titleStyle.Render("Diagnosis") + "\n")
	sb.WriteString(dimStyle.Render(file) + "\n\n")

	if ev.OK {
		sb.WriteString(okStyle.Render("✔ Compatible with DaVinci Resolve (free)") + "\n")
	} else {
		sb.WriteString(warnStyle.Render("Issues found:") + "\n")
		for _, issue := range ev.Issues {
			sb.WriteString(errorStyle.Render("  • "+issue) + "\n")
		}
	}

	if len(ev.Actions) > 0 {
		sb.WriteString("\n" + titleStyle.Render("Actions:") + "\n")
		for _, a := range ev.Actions {
			fmt.Fprintf(&sb, "  → %s\n", a)
		}
	}

	sb.WriteString("\n" + dimStyle.Render("[enter/c] convert   [esc] back   [q] quit"))
	return sb.String()
}
