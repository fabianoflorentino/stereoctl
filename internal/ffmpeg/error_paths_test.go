package ffmpeg

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"
)

func TestFFmpegCmdRunner_BinaryMissing(t *testing.T) {
	orig := os.Getenv("FFMPEG_BIN")
	_ = os.Setenv("FFMPEG_BIN", "nonexistent-ffmpeg-binary-xyz")
	defer func() { _ = os.Setenv("FFMPEG_BIN", orig) }()

	_, _, err := ffmpegCmdRunner([]string{"-version"})
	if err == nil {
		t.Fatalf("expected error when ffmpeg binary missing")
	}
}

func TestProbeCmdRunner_BinaryMissing(t *testing.T) {
	orig := os.Getenv("FFPROBE_BIN")
	_ = os.Setenv("FFPROBE_BIN", "nonexistent-ffprobe-binary-xyz")
	defer func() { _ = os.Setenv("FFPROBE_BIN", orig) }()

	_, err := probeCmdRunner([]string{"-version"})
	if err == nil {
		t.Fatalf("expected error when ffprobe binary missing")
	}
}

func TestConvert_ReturnsErrorWhenRunnerFails(t *testing.T) {
	orig := ffmpegCmdRunner
	ffmpegCmdRunner = func(args []string) (io.ReadCloser, func() error, error) {
		return nil, nil, errors.New("runner failed")
	}
	defer func() { ffmpegCmdRunner = orig }()

	opts := ConvertOptions{Input: "in", Output: "out", CopyAudio: true}
	if err := Convert(opts, "1"); err == nil {
		t.Fatalf("expected Convert to return error when runner fails")
	}
}

func TestMonitorProgress_InvalidNumber(t *testing.T) {
	var out bytes.Buffer
	data := "out_time_us=notanumber\n"
	MonitorProgress(bytes.NewBufferString(data), 1000000, &out)
	if out.Len() != 0 {
		t.Fatalf("expected no output for invalid out_time_us, got %q", out.String())
	}
}
