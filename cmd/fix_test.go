package cmd

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Test preview mode prints ffmpeg command without executing ffmpeg.
func TestFixPreview(t *testing.T) {
	// create a fake ffprobe script that prints minimal JSON
	tmpdir := t.TempDir()
	script := filepath.Join(tmpdir, "fake-ffprobe.sh")
	content := "#!/bin/sh\ncat <<'JSON'\n{\"streams\":[{\"codec_type\":\"audio\",\"codec_name\":\"eac3\",\"channels\":6}],\"format\":{\"duration\":\"10\"}}\nJSON\n"
	if err := os.WriteFile(script, []byte(content), 0700); err != nil {
		t.Fatalf("write script: %v", err)
	}

	// set FFPROBE_BIN to our script
	orig := os.Getenv("FFPROBE_BIN")
	_ = os.Setenv("FFPROBE_BIN", script)
	defer func() { _ = os.Setenv("FFPROBE_BIN", orig) }()

	// set preview flag with defer to ensure cleanup even on early failure
	flagPreview = true
	defer func() { flagPreview = false }()

	// capture stderr with defer to ensure restoration even on early failure
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	defer func() {
		_ = w.Close()
		os.Stderr = old
	}()

	// run fixCmd with a dummy filename
	err := fixCmd.RunE(nil, []string{"input.mkv"})
	_ = w.Close()
	os.Stderr = old

	out, _ := io.ReadAll(r)
	s := string(out)

	if err != nil {
		t.Fatalf("fixCmd.RunE returned error: %v; output: %s", err, s)
	}

	if !strings.Contains(s, "Preview ffmpeg command:") {
		t.Fatalf("expected preview output, got: %s", s)
	}

	// sanity: command should include ffmpeg and -c:a aac
	if !strings.Contains(s, "-c:a aac") {
		t.Fatalf("expected ffmpeg args to include -c:a aac, got: %s", s)
	}
}
