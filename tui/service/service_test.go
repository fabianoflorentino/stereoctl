package service_test

import (
	"errors"
	"testing"

	"github.com/fabianoflorentino/stereoctl/internal/ffmpeg"
	"github.com/fabianoflorentino/stereoctl/tui/service"
)

// --- Prober tests ---

func TestFFmpegProberImplementsProberInterface(t *testing.T) {
	var _ service.Prober = service.FFmpegProber{}
}

type stubProber struct {
	out *ffmpeg.ProbeOutput
	err error
}

func (s stubProber) Probe(_ string) (*ffmpeg.ProbeOutput, error) { return s.out, s.err }

func TestStubProberReturnsConfiguredOutput(t *testing.T) {
	want := &ffmpeg.ProbeOutput{}
	p := stubProber{out: want}
	got, err := p.Probe("any.mp4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestStubProberReturnsConfiguredError(t *testing.T) {
	want := errors.New("probe failed")
	p := stubProber{err: want}
	_, err := p.Probe("any.mp4")
	if err != want {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

// --- Converter tests ---

func TestFFmpegConverterImplementsConverterInterface(t *testing.T) {
	var _ service.Converter = service.FFmpegConverter{}
}

type stubConverter struct {
	err        error
	progValues []float64
}

func (s stubConverter) Convert(_ ffmpeg.ConvertOptions, _ string) error { return s.err }
func (s stubConverter) ConvertWithProgress(_ ffmpeg.ConvertOptions, _ string, onProgress func(float64)) error {
	for _, p := range s.progValues {
		onProgress(p)
	}
	return s.err
}

func TestStubConverterReturnsNilOnSuccess(t *testing.T) {
	c := stubConverter{err: nil}
	if err := c.Convert(ffmpeg.ConvertOptions{}, "10.0"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStubConverterReturnsConfiguredError(t *testing.T) {
	want := errors.New("convert failed")
	c := stubConverter{err: want}
	err := c.Convert(ffmpeg.ConvertOptions{}, "10.0")
	if err != want {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

func TestFFmpegConverterImplementsConvertWithProgress(t *testing.T) {
	var _ service.Converter = service.FFmpegConverter{}
}

func TestStubConverterConvertWithProgressCallsCallback(t *testing.T) {
	c := stubConverter{progValues: []float64{0.25, 0.75, 1.0}}
	var got []float64
	err := c.ConvertWithProgress(ffmpeg.ConvertOptions{}, "10.0", func(p float64) {
		got = append(got, p)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 || got[0] != 0.25 || got[1] != 0.75 || got[2] != 1.0 {
		t.Fatalf("unexpected progress values: %v", got)
	}
}

func TestStubConverterConvertWithProgressReturnsError(t *testing.T) {
	want := errors.New("ffmpeg boom")
	c := stubConverter{err: want}
	err := c.ConvertWithProgress(ffmpeg.ConvertOptions{}, "10.0", func(_ float64) {})
	if err != want {
		t.Fatalf("expected %v, got %v", want, err)
	}
}
