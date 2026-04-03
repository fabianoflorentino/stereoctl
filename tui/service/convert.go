package service

import "github.com/fabianoflorentino/stereoctl/internal/ffmpeg"

// Converter defines the contract for converting a media file.
type Converter interface {
	Convert(opts ffmpeg.ConvertOptions, duration string) error
	// ConvertWithProgress runs conversion and calls onProgress (0.0–1.0) periodically.
	ConvertWithProgress(opts ffmpeg.ConvertOptions, duration string, onProgress func(float64)) error
}

// FFmpegConverter is the real implementation that delegates to ffmpeg.
type FFmpegConverter struct{}

func (FFmpegConverter) Convert(opts ffmpeg.ConvertOptions, duration string) error {
	return ffmpeg.Convert(opts, duration)
}

func (FFmpegConverter) ConvertWithProgress(opts ffmpeg.ConvertOptions, duration string, onProgress func(float64)) error {
	return ffmpeg.ConvertWithProgress(opts, duration, onProgress)
}
