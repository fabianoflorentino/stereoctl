package ffmpeg

import (
	"bytes"
	"testing"
)

func TestMonitorProgress(t *testing.T) {
	// simulate ffmpeg progress lines (microseconds)
	data := "out_time_us=0\nout_time_us=1000000\nout_time_us=2000000\nout_time_us=3000000\nout_time_us=4000000\n"
	totalUs := int64(4000000)

	r := bytes.NewBufferString(data)
	var out bytes.Buffer

	MonitorProgress(r, totalUs, &out)

	got := out.String()
	if got == "" {
		t.Fatal("expected progress output, got empty string")
	}

	// The last progress line is out_time_us=4000000 which equals totalUs,
	// so the final output must contain 100.0%.
	if !bytes.Contains([]byte(got), []byte("100.0%")) {
		t.Fatalf("expected 100.0%% in final progress output, got: %q", got)
	}
	// Intermediate line at 2000000us (50%) must also appear.
	if !bytes.Contains([]byte(got), []byte("50.0%")) {
		t.Fatalf("expected 50.0%% in progress output, got: %q", got)
	}
}
