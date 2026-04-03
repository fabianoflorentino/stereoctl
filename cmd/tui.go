package cmd

import (
	"fmt"
	"os"

	"github.com/fabianoflorentino/stereoctl/tui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch interactive terminal UI",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := tui.Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
