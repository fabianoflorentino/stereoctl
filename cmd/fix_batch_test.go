package cmd

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFixBatch(t *testing.T) {
	tmp := t.TempDir()
	// create sample files
	files := []string{"a.mkv", "b.mp4", "c.txt"}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(tmp, f), []byte("dummy"), 0644); err != nil {
			t.Fatalf("write sample: %v", err)
		}
	}

	// create fake ffprobe and ffmpeg scripts and set env vars
	ffprobeScript := filepath.Join(tmp, "fake-ffprobe.sh")
	ffprobeContent := "#!/bin/sh\ncat <<'JSON'\n{" + "\"streams\": [{\"codec_type\":\"audio\",\"codec_name\":\"eac3\",\"channels\":6}],\"format\":{\"duration\":\"1\"}}\nJSON\n"
	if err := os.WriteFile(ffprobeScript, []byte(ffprobeContent), 0700); err != nil {
		t.Fatalf("write ffprobe script: %v", err)
	}
	_ = os.Setenv("FFPROBE_BIN", ffprobeScript)
	defer func() { _ = os.Unsetenv("FFPROBE_BIN") }()

	ffmpegScript := filepath.Join(tmp, "fake-ffmpeg.sh")
	ffmpegContent := "#!/bin/sh\n# emit a single progress line and exit\necho out_time_us=0\nexit 0\n"
	if err := os.WriteFile(ffmpegScript, []byte(ffmpegContent), 0700); err != nil {
		t.Fatalf("write ffmpeg script: %v", err)
	}
	_ = os.Setenv("FFMPEG_BIN", ffmpegScript)
	defer func() { _ = os.Unsetenv("FFMPEG_BIN") }()

	// capture stderr
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// run batch
	flagBatch = true
	err := fixCmd.RunE(nil, []string{tmp})

	// restore
	_ = w.Close()
	os.Stderr = old
	out, _ := io.ReadAll(r)
	s := string(out)

	if err != nil {
		t.Fatalf("fix batch returned error: %v; output: %s", err, s)
	}

	if !strings.Contains(s, "Processing") || !strings.Contains(s, "Done:") {
		t.Fatalf("unexpected batch output: %s", s)
	}

	flagBatch = false
}
