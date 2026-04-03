package ffmpeg

import (
	"bytes"
	"testing"
)

func TestWriteProgress_ZeroTotal(t *testing.T) {
	var buf bytes.Buffer
	writeProgress(&buf, 1000, 0)
	if buf.Len() != 0 {
		t.Fatalf("expected no output when total is 0, got %q", buf.String())
	}
}

func TestMonitorProgress_NonProgressLines(t *testing.T) {
	r := bytes.NewBufferString("some random line\nanother line\n")
	var out bytes.Buffer
	MonitorProgress(r, 1000, &out)
	if out.Len() != 0 {
		t.Fatalf("expected no progress output for non-progress lines, got %q", out.String())
	}
}
