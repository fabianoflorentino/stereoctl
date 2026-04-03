package cmd

import (
	"fmt"
	"os"

	"github.com/fabianoflorentino/stereoctl/internal/ffmpeg"
	"github.com/fabianoflorentino/stereoctl/internal/profiles"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check <file>",
	Short: "Check a file for Resolve compatibility",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]

		p, err := ffmpeg.Probe(input)
		if err != nil {
			return err
		}

		ev := profiles.EvaluateResolveFree(p)

		if ev.OK {
			fmt.Fprintf(os.Stdout, "OK: file seems compatible with Resolve free\n")
		} else {
			fmt.Fprintf(os.Stdout, "Issues found:\n")
			for _, it := range ev.Issues {
				fmt.Fprintf(os.Stdout, " - %s\n", it)
			}
		}

		if len(ev.Actions) > 0 {
			fmt.Fprintf(os.Stdout, "Recommended actions:\n")
			for _, a := range ev.Actions {
				fmt.Fprintf(os.Stdout, " - %s\n", a)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
