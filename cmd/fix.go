package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fabianoflorentino/stereoctl/internal/ffmpeg"
	"github.com/fabianoflorentino/stereoctl/internal/profiles"
	"github.com/spf13/cobra"
)

var (
	flagProfile string
	flagPreview bool
	flagBatch   bool
)

var fixCmd = &cobra.Command{
	Use:  "fix <file>",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]

		probe, err := ffmpeg.Probe(input)
		if err != nil {
			return fmt.Errorf("probe failed: %w", err)
		}

		switch flagProfile {
		case "resolve-free":
			// batch mode: expand pattern or directory into file list
			var files []string
			if flagBatch {
				// if input contains glob chars, use Glob
				if strings.ContainsAny(input, "*?[") {
					matches, _ := filepath.Glob(input)
					files = append(files, matches...)
				} else {
					// if it's a directory, walk and collect known video extensions
					fi, err := os.Stat(input)
					if err == nil && fi.IsDir() {
						extset := map[string]bool{".mp4": true, ".mkv": true, ".mov": true, ".webm": true, ".mxf": true}
						_ = filepath.Walk(input, func(p string, info os.FileInfo, err error) error {
							if err != nil || info.IsDir() {
								return nil
							}
							if extset[strings.ToLower(filepath.Ext(p))] {
								files = append(files, p)
							}
							return nil
						})
					} else {
						// treat input as single file pattern
						files = append(files, input)
					}
				}

				if len(files) == 0 {
					fmt.Fprintf(os.Stderr, "no files matched for batch input: %s\n", input)
					return nil
				}

				for _, f := range files {
					probeF, err := ffmpeg.Probe(f)
					if err != nil {
						fmt.Fprintf(os.Stderr, "probe failed for %s: %v\n", f, err)
						continue
					}
					opts := profiles.BuildConvertOptionsForResolveFree(probeF, f)
					if flagOutput != "" {
						opts.Output = flagOutput
					}
					if flagPreview {
						args := ffmpeg.BuildFFmpegArgs(opts)
						fmt.Fprintf(os.Stderr, "Preview for %s:\n%s %s\n", f, "ffmpeg", strings.Join(args, " "))
						continue
					}
					fmt.Fprintf(os.Stderr, "Processing %s...\n", f)
					if err := ffmpeg.Convert(opts, probeF.Format.Duration); err != nil {
						fmt.Fprintf(os.Stderr, "conversion failed for %s: %v\n", f, err)
						continue
					}
					fmt.Fprintf(os.Stderr, "Done: %s\n", opts.Output)
				}
				return nil
			}

			// single-file path
			opts := profiles.BuildConvertOptionsForResolveFree(probe, input)
			// override output if provided
			if flagOutput != "" {
				opts.Output = flagOutput
			}

			if flagPreview {
				args := ffmpeg.BuildFFmpegArgs(opts)
				fmt.Fprintf(os.Stderr, "Preview ffmpeg command:\n%s %s\n", "ffmpeg", strings.Join(args, " "))
				return nil
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
	fixCmd.Flags().BoolVarP(&flagPreview, "preview", "n", false, "preview ffmpeg command without running it")
	fixCmd.Flags().BoolVarP(&flagBatch, "batch", "B", false, "treat argument as directory or glob and process multiple files")
	// `flagOutput` is defined in convert.go and reused here via the global variable
	rootCmd.AddCommand(fixCmd)
}
