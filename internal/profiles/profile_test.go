package profiles

import (
	"strings"
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
	// both video and audio should be flagged
	if len(ev.Issues) < 2 {
		t.Fatalf("expected at least 2 issues (video + audio), got %d: %v", len(ev.Issues), ev.Issues)
	}
}

func TestEvaluateResolveFree_HEVC_OK(t *testing.T) {
	p := &ffmpeg.ProbeOutput{
		Streams: []struct {
			CodecType string "json:\"codec_type\""
			CodecName string "json:\"codec_name\""
			Channels  int    "json:\"channels\""
		}{
			{CodecType: "video", CodecName: "hevc"},
			{CodecType: "audio", CodecName: "aac", Channels: 2},
		},
		Format: struct {
			Duration string "json:\"duration\""
		}{Duration: "10"},
	}

	ev := EvaluateResolveFree(p)
	if !ev.OK {
		t.Fatalf("expected OK for HEVC+AAC stereo, got issues: %v", ev.Issues)
	}
}

func TestEvaluateResolveFree_NoVideo(t *testing.T) {
	p := &ffmpeg.ProbeOutput{
		Streams: []struct {
			CodecType string "json:\"codec_type\""
			CodecName string "json:\"codec_name\""
			Channels  int    "json:\"channels\""
		}{
			{CodecType: "audio", CodecName: "aac", Channels: 2},
		},
		Format: struct {
			Duration string "json:\"duration\""
		}{Duration: "10"},
	}

	ev := EvaluateResolveFree(p)
	if ev.OK {
		t.Fatalf("expected not OK when no video stream")
	}
	var found bool
	for _, issue := range ev.Issues {
		if strings.Contains(issue, "no video") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected 'no video stream' issue, got: %v", ev.Issues)
	}
}

func TestEvaluateResolveFree_NoAudio(t *testing.T) {
	p := &ffmpeg.ProbeOutput{
		Streams: []struct {
			CodecType string "json:\"codec_type\""
			CodecName string "json:\"codec_name\""
			Channels  int    "json:\"channels\""
		}{
			{CodecType: "video", CodecName: "h264"},
		},
		Format: struct {
			Duration string "json:\"duration\""
		}{Duration: "10"},
	}

	ev := EvaluateResolveFree(p)
	if ev.OK {
		t.Fatalf("expected not OK when no audio stream")
	}
	var found bool
	for _, issue := range ev.Issues {
		if strings.Contains(issue, "no audio") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected 'no audio stream' issue, got: %v", ev.Issues)
	}
}
