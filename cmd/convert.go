package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fabianoflorentino/stereoctl/internal/ffmpeg"
	"github.com/spf13/cobra"
)

var (
	flagOutput  string
	flagBitrate string
)

var convertCmd = &cobra.Command{
	Use:   "convert <file>",
	Short: "Convert video audio to AAC stereo",
	Long: `Probe the given video file and convert its audio track to AAC stereo.
The video stream is copied without re-encoding.

If the audio is already AAC stereo the file is simply remuxed.`,
	Args:    cobra.ExactArgs(1),
	Example: "  stereoctl convert movie.mkv\n  stereoctl convert movie.mkv --output fixed.mp4 --bitrate 256k",
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]

		output := flagOutput
		if output == "" {
			ext := filepath.Ext(input)
			output = strings.TrimSuffix(input, ext) + ".mp4"
		}

		probe, err := ffmpeg.Probe(input)
		if err != nil {
			return fmt.Errorf("probe failed: %w", err)
		}

		var codec string
		var channels int
		for _, s := range probe.Streams {
			if s.CodecType == "audio" {
				codec = s.CodecName
				channels = s.Channels
				break
			}
		}

		fmt.Fprintf(os.Stderr, "Detected audio: %s - %d channels\n", codec, channels)

		alreadyOk := codec == "aac" && channels <= 2
		if alreadyOk {
			fmt.Fprintln(os.Stderr, "✔ Already compatible. Remuxing...")
		} else {
			fmt.Fprintln(os.Stderr, "⚠ Converting audio to AAC stereo...")
		}

		opts := ffmpeg.ConvertOptions{
			Input:      input,
			Output:     output,
			CopyAudio:  alreadyOk,
			AudioCodec: "aac",
			Channels:   2,
			Bitrate:    flagBitrate,
		}

		if err := ffmpeg.Convert(opts, probe.Format.Duration); err != nil {
			return fmt.Errorf("conversion failed: %w", err)
		}

		fmt.Fprintf(os.Stderr, "\n✅ Done: %s\n", output)
		return nil
	},
}

func init() {
	convertCmd.Flags().StringVarP(&flagOutput, "output", "o", "", "output file path (default: same name as input with .mp4 extension)")
	convertCmd.Flags().StringVarP(&flagBitrate, "bitrate", "b", "192k", "audio bitrate for conversion")
}
