package cmd

import (
	"fmt"
	"os"

	"github.com/fabianoflorentino/stereoctl/internal/ffmpeg"
	"github.com/fabianoflorentino/stereoctl/internal/profiles"
	"github.com/spf13/cobra"
)

var (
	flagProfile string
)

var fixCmd = &cobra.Command{
	Use:   "fix <file>",
	Short: "Apply profile-based fixes to make a file compatible (e.g. resolve-free)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]

		probe, err := ffmpeg.Probe(input)
		if err != nil {
			return fmt.Errorf("probe failed: %w", err)
		}

		switch flagProfile {
		case "resolve-free":
			opts := profiles.BuildConvertOptionsForResolveFree(probe, input)
			// override output if provided
			if flagOutput != "" {
				opts.Output = flagOutput
			}

			fmt.Fprintf(os.Stderr, "Running ffmpeg with options: %+v\n", opts)

			if err := ffmpeg.Convert(opts, probe.Format.Duration); err != nil {
				return fmt.Errorf("conversion failed: %w", err)
			}

			fmt.Fprintf(os.Stderr, "Done: %s\n", opts.Output)
			return nil
		default:
			return fmt.Errorf("unknown profile: %s", flagProfile)
		}
	},
}

func init() {
	fixCmd.Flags().StringVarP(&flagProfile, "profile", "p", "resolve-free", "profile to apply (resolve-free)")
	// `flagOutput` is defined in convert.go and reused here via the global variable
	rootCmd.AddCommand(fixCmd)
}
