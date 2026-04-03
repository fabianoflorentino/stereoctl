package profiles

import (
	"path/filepath"

	"github.com/fabianoflorentino/stereoctl/internal/ffmpeg"
)

// BuildConvertOptionsForResolveFree returns ffmpeg.ConvertOptions suitable for
// making the input compatible with DaVinci Resolve (free).
func BuildConvertOptionsForResolveFree(p *ffmpeg.ProbeOutput, input string) ffmpeg.ConvertOptions {
	// default output: same name with .mp4
	ext := filepath.Ext(input)
	output := input[:len(input)-len(ext)] + ".mp4"

	opts := ffmpeg.ConvertOptions{
		Input:   input,
		Output:  output,
		Bitrate: "192k",
	}

	// video: re-encode if not h264
	for _, s := range p.Streams {
		if s.CodecType == "video" {
			if s.CodecName != "h264" {
				opts.VideoCodec = "libx264"
				opts.VideoExtraArgs = []string{"-preset", "slow", "-crf", "20"}
			}
			break
		}
	}

	// audio: check first audio stream
	for _, s := range p.Streams {
		if s.CodecType == "audio" {
			if s.CodecName == "aac" && s.Channels <= 2 {
				opts.CopyAudio = true
			} else {
				opts.CopyAudio = false
				opts.AudioCodec = "aac"
				opts.Channels = 2
				opts.Bitrate = "192k"
			}
			break
		}
	}

	return opts
}
