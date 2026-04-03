package profiles

import (
	"testing"

	"github.com/fabianoflorentino/stereoctl/internal/ffmpeg"
)

func TestEvaluateResolveFree_OK(t *testing.T) {
	p := &ffmpeg.ProbeOutput{
		Streams: []struct {
			CodecType string "json:\"codec_type\""
			CodecName string "json:\"codec_name\""
			Channels  int    "json:\"channels\""
		}{
			{CodecType: "video", CodecName: "h264"},
			{CodecType: "audio", CodecName: "aac", Channels: 2},
		},
		Format: struct {
			Duration string "json:\"duration\""
		}{Duration: "10"},
	}

	ev := EvaluateResolveFree(p)
	if !ev.OK {
		t.Fatalf("expected OK, got issues: %v", ev.Issues)
	}
}

func TestEvaluateResolveFree_Issues(t *testing.T) {
	p := &ffmpeg.ProbeOutput{
		Streams: []struct {
			CodecType string "json:\"codec_type\""
			CodecName string "json:\"codec_name\""
			Channels  int    "json:\"channels\""
		}{
			{CodecType: "video", CodecName: "vp9"},
			{CodecType: "audio", CodecName: "eac3", Channels: 6},
		},
		Format: struct {
			Duration string "json:\"duration\""
		}{Duration: "10"},
	}

	ev := EvaluateResolveFree(p)
	if ev.OK {
		t.Fatalf("expected issues, got OK")
	}
	if len(ev.Issues) == 0 {
		t.Fatalf("expected at least one issue")
	}
}
