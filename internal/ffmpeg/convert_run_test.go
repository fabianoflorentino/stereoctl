package ffmpeg

import (
	"bytes"
	"io"
	"testing"
)

func TestConvert_WithFakeRunner(t *testing.T) {
	// prepare fake progress output (microseconds)
	data := "out_time_us=0\nout_time_us=1000000\nout_time_us=2000000\nout_time_us=3000000\nout_time_us=4000000\n"
	r := io.NopCloser(bytes.NewBufferString(data))

	// fake runner returns the reader and a wait function
	orig := ffmpegCmdRunner
	ffmpegCmdRunner = func(args []string) (io.ReadCloser, func() error, error) {
		return r, func() error { return nil }, nil
	}
	defer func() { ffmpegCmdRunner = orig }()

	opts := ConvertOptions{
		Input:     "in.mp4",
		Output:    "out.mp4",
		CopyAudio: true,
	}

	if err := Convert(opts, "4"); err != nil {
		t.Fatalf("Convert returned error with fake runner: %v", err)
	}
}
