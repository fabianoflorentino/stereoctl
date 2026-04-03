package ffmpeg

import (
	"testing"
)

func TestParseProbeOutput(t *testing.T) {
	sample := []byte(`{"streams":[{"codec_type":"video","codec_name":"h264"},{"codec_type":"audio","codec_name":"eac3","channels":6}],"format":{"duration":"1285.184"}}`)

	p, err := ParseProbeOutput(sample)
	if err != nil {
		t.Fatalf("ParseProbeOutput error: %v", err)
	}

	if len(p.Streams) != 2 {
		t.Fatalf("expected 2 streams, got %d", len(p.Streams))
	}

	// check audio stream
	var found bool
	for _, s := range p.Streams {
		if s.CodecType == "audio" {
			found = true
			if s.CodecName != "eac3" {
				t.Fatalf("expected audio codec eac3, got %s", s.CodecName)
			}
			if s.Channels != 6 {
				t.Fatalf("expected 6 channels, got %d", s.Channels)
			}
		}
	}
	if !found {
		t.Fatalf("audio stream not found")
	}

	if p.Format.Duration != "1285.184" {
		t.Fatalf("unexpected duration: %s", p.Format.Duration)
	}
}
