package profiles

import (
	"testing"

	"github.com/fabianoflorentino/stereoctl/internal/ffmpeg"
)

func TestBuildConvertOptionsForResolveFree_H264Aac(t *testing.T) {
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

	opts := BuildConvertOptionsForResolveFree(p, "sample.mkv")

	if opts.Output != "sample.mp4" {
		t.Fatalf("expected output sample.mp4, got %s", opts.Output)
	}
	if opts.VideoCodec != "" {
		t.Fatalf("expected video copy (empty VideoCodec), got %s", opts.VideoCodec)
	}
	if !opts.CopyAudio {
		t.Fatalf("expected CopyAudio true for AAC stereo")
	}
}

func TestBuildConvertOptionsForResolveFree_ReencodeAndDownmix(t *testing.T) {
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

	opts := BuildConvertOptionsForResolveFree(p, "clip.webm")

	if opts.VideoCodec != "libx264" {
		t.Fatalf("expected video re-encode to libx264, got %s", opts.VideoCodec)
	}
	if len(opts.VideoExtraArgs) == 0 {
		t.Fatalf("expected VideoExtraArgs for encoding settings")
	}
	if opts.CopyAudio {
		t.Fatalf("expected CopyAudio false for multichannel audio")
	}
	if opts.AudioCodec != "aac" {
		t.Fatalf("expected AudioCodec aac, got %s", opts.AudioCodec)
	}
	if opts.Channels != 2 {
		t.Fatalf("expected Channels 2, got %d", opts.Channels)
	}
}

func TestBuildConvertOptionsForResolveFree_HEVCAac(t *testing.T) {
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

	opts := BuildConvertOptionsForResolveFree(p, "sample.mkv")

	// HEVC is Resolve-compatible; must copy video, not re-encode
	if opts.VideoCodec != "" {
		t.Fatalf("expected video copy for HEVC input (empty VideoCodec), got %s", opts.VideoCodec)
	}
	if len(opts.VideoExtraArgs) != 0 {
		t.Fatalf("expected no VideoExtraArgs for HEVC copy, got %v", opts.VideoExtraArgs)
	}
	if !opts.CopyAudio {
		t.Fatalf("expected CopyAudio true for AAC stereo")
	}
}
