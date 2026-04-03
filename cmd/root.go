package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev" // overridden at build time via -ldflags "-X ...cmd.version=X.Y.Z"

var rootCmd = &cobra.Command{
	Use:     "stereoctl",
	Short:   "Convert video audio to AAC stereo",
	Long:    "stereoctl inspects video files and converts incompatible audio tracks to AAC stereo, keeping the video stream untouched.",
	Version: version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(convertCmd)
}
