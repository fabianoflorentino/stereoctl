package ffmpeg

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// Integration test that uses a real sample video in testdata/sample.mp4.
// The test will be skipped if `ffmpeg` is not available on PATH.
func TestConvertIntegration(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not found, skipping integration test")
	}

	in := filepath.Join("..", "..", "testdata", "sample.mp4")
	if _, err := os.Stat(in); err != nil {
		t.Fatalf("sample file not found: %v", err)
	}

	// Probe
	p, err := Probe(in)
	if err != nil {
		t.Fatalf("Probe failed: %v", err)
	}

	// choose audio handling based on probe result
	var copyAudio bool
	var codec string
	var channels int
	for _, s := range p.Streams {
		if s.CodecType == "audio" {
			codec = s.CodecName
			channels = s.Channels
			break
		}
	}

	if codec == "aac" && channels <= 2 {
		copyAudio = true
	}

	out := filepath.Join(os.TempDir(), "stereoctl_test_output.mp4")
	defer os.Remove(out)

	opts := ConvertOptions{
		Input:      in,
		Output:     out,
		CopyAudio:  copyAudio,
		AudioCodec: "aac",
		Channels:   2,
		Bitrate:    "192k",
	}

	if err := Convert(opts, p.Format.Duration); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	if _, err := os.Stat(out); err != nil {
		t.Fatalf("expected output file created, but stat failed: %v", err)
	}
}
