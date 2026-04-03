package profiles

import (
	"fmt"

	"github.com/fabianoflorentino/stereoctl/internal/ffmpeg"
)

// Evaluation holds the check result for a profile.
type Evaluation struct {
	OK      bool
	Issues  []string
	Actions []string
}

// EvaluateResolveFree checks compatibility with DaVinci Resolve (free).
func EvaluateResolveFree(p *ffmpeg.ProbeOutput) Evaluation {
	ev := Evaluation{OK: true}

	if p == nil {
		ev.Issues = append(ev.Issues, "missing probe information")
		ev.OK = false
		return ev
	}

	// Video: prefer h264
	var hasVideo bool
	for _, s := range p.Streams {
		if s.CodecType == "video" {
			hasVideo = true
			if s.CodecName != "h264" && s.CodecName != "hevc" {
				ev.Issues = append(ev.Issues, fmt.Sprintf("video codec %s may be incompatible; prefer h264", s.CodecName))
				ev.Actions = append(ev.Actions, "re-encode video to H.264 (libx264)")
				ev.OK = false
			}
			break
		}
	}
	if !hasVideo {
		ev.Issues = append(ev.Issues, "no video stream found")
		ev.OK = false
	}

	// Audio: must be stereo AAC or PCM
	var audioFound bool
	for _, s := range p.Streams {
		if s.CodecType == "audio" {
			audioFound = true
			if s.CodecName == "aac" && s.Channels <= 2 {
				// OK
			} else if s.Channels > 2 {
				ev.Issues = append(ev.Issues, fmt.Sprintf("audio has %d channels; Resolve free needs stereo", s.Channels))
				ev.Actions = append(ev.Actions, "downmix to stereo (AAC or PCM)")
				ev.OK = false
			} else if s.CodecName != "aac" && s.CodecName != "pcm_s16le" {
				ev.Issues = append(ev.Issues, fmt.Sprintf("audio codec %s may be incompatible; prefer AAC stereo or PCM", s.CodecName))
				ev.Actions = append(ev.Actions, "transcode audio to AAC stereo")
				ev.OK = false
			}
			// check only first audio stream for recommendations
			break
		}
	}
	if !audioFound {
		ev.Issues = append(ev.Issues, "no audio stream found")
		ev.Actions = append(ev.Actions, "add an AAC stereo audio track or PCM")
		ev.OK = false
	}

	// Container: recommend MP4

	if ev.OK {
		ev.Actions = append(ev.Actions, "file looks compatible with Resolve free; consider remuxing to MP4 if not already")
	}

	return ev
}
