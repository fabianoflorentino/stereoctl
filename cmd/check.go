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
			_, _ = fmt.Fprintln(os.Stdout, "OK: file seems compatible with Resolve free")
		} else {
			_, _ = fmt.Fprintln(os.Stdout, "Issues found:")
			for _, it := range ev.Issues {
				_, _ = fmt.Fprintf(os.Stdout, " - %s\n", it)
			}
		}

		if len(ev.Actions) > 0 {
			_, _ = fmt.Fprintln(os.Stdout, "Recommended actions:")
			for _, a := range ev.Actions {
				_, _ = fmt.Fprintf(os.Stdout, " - %s\n", a)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
