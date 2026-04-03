package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

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
