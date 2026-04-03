package ffmpeg

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

// ProbeOutput represents the relevant parts of ffprobe's JSON output.
type ProbeOutput struct {
	Streams []struct {
		CodecType string `json:"codec_type"`
		CodecName string `json:"codec_name"`
		Channels  int    `json:"channels"`
	} `json:"streams"`
	Format struct {
		Duration string `json:"duration"`
	} `json:"format"`
}

// Probe runs ffprobe on the given input file and returns the parsed output.
func Probe(input string) (*ProbeOutput, error) {
	args := []string{"-v", "quiet", "-print_format", "json", "-show_streams", "-show_format", input}

	out, err := probeCmdRunner(args)
	if err != nil {
		return nil, fmt.Errorf("ffprobe: %w", err)
	}

	return ParseProbeOutput(out)
}

// probeCmdRunner allows tests to inject fake ffprobe output. By default it
// executes the `ffprobe` binary (or the path set in FFPROBE_BIN).
var probeCmdRunner = func(args []string) ([]byte, error) {
	bin := os.Getenv("FFPROBE_BIN")
	if bin == "" {
		bin = "ffprobe"
	}
	cmd := exec.Command(bin, args...)
	return cmd.Output()
}

// ParseProbeOutput parses ffprobe JSON output into a ProbeOutput struct.
// This is exposed to allow unit testing without executing ffprobe.
func ParseProbeOutput(data []byte) (*ProbeOutput, error) {
	var probe ProbeOutput
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, fmt.Errorf("parse ffprobe output: %w", err)
	}
	return &probe, nil
}
