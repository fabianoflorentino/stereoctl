package ffmpeg

import (
	"strings"
	"testing"
)

func TestBuildFFmpegArgs_VideoCopy_AudioCopy(t *testing.T) {
	opts := ConvertOptions{
		Input:     "in.mkv",
		Output:    "out.mp4",
		CopyAudio: true,
	}

	args := BuildFFmpegArgs(opts)
	s := strings.Join(args, " ")

	if !strings.Contains(s, "-c:v copy") {
		t.Fatalf("expected video copy in args, got %s", s)
	}
	if !strings.Contains(s, "-c:a copy") {
		t.Fatalf("expected audio copy in args, got %s", s)
	}
}

func TestBuildFFmpegArgs_ReencodeAudioAndVideo(t *testing.T) {
	opts := ConvertOptions{
		Input:          "in.webm",
		Output:         "out.mp4",
		VideoCodec:     "libx264",
		VideoExtraArgs: []string{"-preset", "slow", "-crf", "20"},
		CopyAudio:      false,
		AudioCodec:     "aac",
		Channels:       2,
		Bitrate:        "192k",
		ExtraArgs:      []string{"-t", "30"},
	}

	args := BuildFFmpegArgs(opts)
	s := strings.Join(args, " ")

	if !strings.Contains(s, "-c:v libx264") {
		t.Fatalf("expected video codec libx264 in args, got %s", s)
	}
	if !strings.Contains(s, "-preset slow") || !strings.Contains(s, "-crf 20") {
		t.Fatalf("expected video extra args in args, got %s", s)
	}
	if !strings.Contains(s, "-c:a aac") || !strings.Contains(s, "-ac 2") || !strings.Contains(s, "-b:a 192k") {
		t.Fatalf("expected audio transcode args, got %s", s)
	}
	if !strings.Contains(s, "-t 30") {
		t.Fatalf("expected extra args present, got %s", s)
	}
}
