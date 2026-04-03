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

	// expect at least a 100% or 0/100% markers somewhere
	if !bytes.Contains([]byte(got), []byte("100.0%")) && !bytes.Contains([]byte(got), []byte("50.0%")) {
		t.Fatalf("unexpected progress output: %q", got)
	}
}
