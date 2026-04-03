package view

import (
	"fmt"
	"strings"
)

// Result renders the final result screen.
func Result(output string, err error) string {
	var sb strings.Builder
	if err != nil {
		sb.WriteString(errorStyle.Render("✖ Error: "+err.Error()) + "\n")
	} else {
		sb.WriteString(okStyle.Render("✔ Done!") + "\n")
		fmt.Fprintf(&sb, "Output: %s\n", output)
	}
	sb.WriteString("\n" + dimStyle.Render("[q] quit"))
	return sb.String()
}
