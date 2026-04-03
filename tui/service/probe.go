package service

import "github.com/fabianoflorentino/stereoctl/internal/ffmpeg"

// Prober defines the contract for probing a media file.
type Prober interface {
	Probe(path string) (*ffmpeg.ProbeOutput, error)
}

// FFmpegProber is the real implementation that delegates to ffprobe.
type FFmpegProber struct{}

func (FFmpegProber) Probe(path string) (*ffmpeg.ProbeOutput, error) {
	return ffmpeg.Probe(path)
}
