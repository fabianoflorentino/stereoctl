package ffmpeg

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestWriteProgress_FullAndPartial(t *testing.T) {
	var buf bytes.Buffer

	total := int64(4_000_000)

	// test partial progress (50%)
	writeProgress(&buf, 2_000_000, total)
	out := buf.String()
	if !strings.Contains(out, "50.0%") {
		t.Fatalf("expected 50.0%% in output, got %q", out)
	}

	// test full progress (100%)
	buf.Reset()
	writeProgress(&buf, 4_000_000, total)
	out = buf.String()
	if !strings.Contains(out, "100.0%") {
		t.Fatalf("expected 100.0%% in output, got %q", out)
	}
}

func TestFormatDur(t *testing.T) {
	d := time.Hour + 2*time.Minute + 3*time.Second
	s := formatDur(d)
	if s != "01:02:03" {
		t.Fatalf("expected 01:02:03, got %s", s)
	}
}
