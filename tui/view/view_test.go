package view_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/fabianoflorentino/stereoctl/internal/profiles"
	"github.com/fabianoflorentino/stereoctl/tui/view"
)

// --- Diagnosis ---

func TestDiagnosisOKShowsCompatible(t *testing.T) {
	ev := profiles.Evaluation{OK: true, Actions: []string{"remux to MP4"}}
	out := view.Diagnosis(ev, "file.mp4")
	if !strings.Contains(out, "Compatible") {
		t.Error("expected compatible message")
	}
	if !strings.Contains(out, "remux to MP4") {
		t.Error("expected action in output")
	}
}

func TestDiagnosisWithIssuesShowsWarning(t *testing.T) {
	ev := profiles.Evaluation{
		OK:      false,
		Issues:  []string{"audio has 6 channels"},
		Actions: []string{"downmix to stereo"},
	}
	out := view.Diagnosis(ev, "movie.mkv")
	if !strings.Contains(out, "Issues found") {
		t.Error("expected issues section")
	}
	if !strings.Contains(out, "audio has 6 channels") {
		t.Error("expected issue detail")
	}
	if !strings.Contains(out, "downmix to stereo") {
		t.Error("expected action")
	}
}

func TestDiagnosisShowsFilename(t *testing.T) {
	ev := profiles.Evaluation{OK: true}
	out := view.Diagnosis(ev, "my_video.mp4")
	if !strings.Contains(out, "my_video.mp4") {
		t.Error("expected filename in output")
	}
}

func TestDiagnosisContainsKeyBindings(t *testing.T) {
	out := view.Diagnosis(profiles.Evaluation{OK: true}, "x.mp4")
	if !strings.Contains(out, "enter") || !strings.Contains(out, "esc") {
		t.Error("expected key binding hints")
	}
}

func TestDiagnosisNoActionsWhenEmpty(t *testing.T) {
	ev := profiles.Evaluation{OK: false, Issues: []string{"bad codec"}, Actions: nil}
	out := view.Diagnosis(ev, "f.mp4")
	if strings.Contains(out, "Actions:") {
		t.Error("should not show Actions section when empty")
	}
}

// --- Result ---

func TestResultSuccessShowsDone(t *testing.T) {
	out := view.Result("output.mp4", nil)
	if !strings.Contains(out, "Done") {
		t.Error("expected Done message")
	}
	if !strings.Contains(out, "output.mp4") {
		t.Error("expected output path")
	}
}

func TestResultErrorShowsMessage(t *testing.T) {
	out := view.Result("", errors.New("conversion failed"))
	if !strings.Contains(out, "Error") {
		t.Error("expected error message")
	}
	if !strings.Contains(out, "conversion failed") {
		t.Error("expected error detail")
	}
}

func TestResultContainsQuitHint(t *testing.T) {
	out := view.Result("out.mp4", nil)
	if !strings.Contains(out, "quit") {
		t.Error("expected quit hint")
	}
}

// --- Converting ---

func TestConvertingContainsTitle(t *testing.T) {
	bar := progress.New()
	out := view.Converting(bar)
	if !strings.Contains(out, "Converting") {
		t.Error("expected Converting title")
	}
}

func TestConvertingContainsQuitHint(t *testing.T) {
	out := view.Converting(progress.New())
	if !strings.Contains(out, "quit") {
		t.Error("expected quit hint")
	}
}

// --- Profile ---

func TestProfileContainsResolveName(t *testing.T) {
	out := view.Profile()
	if !strings.Contains(out, "Resolve") {
		t.Error("expected Resolve profile name")
	}
}

func TestProfileContainsKeyBindings(t *testing.T) {
	out := view.Profile()
	if !strings.Contains(out, "enter") || !strings.Contains(out, "esc") {
		t.Error("expected key binding hints")
	}
}

// --- FileSelection ---

func TestFileSelectionContainsTitle(t *testing.T) {
	fp := filepicker.New()
	out := view.FileSelection(fp)
	if !strings.Contains(out, "Select") {
		t.Error("expected Select title")
	}
}

func TestFileSelectionContainsQuitHint(t *testing.T) {
	out := view.FileSelection(filepicker.New())
	if !strings.Contains(out, "quit") {
		t.Error("expected quit hint")
	}
}

// --- DefaultRenderer ---

func TestDefaultRendererFileSelection(t *testing.T) {
	r := view.DefaultRenderer{}
	out := r.FileSelection(filepicker.New())
	if !strings.Contains(out, "Select") {
		t.Error("expected Select in FileSelection render")
	}
}

func TestDefaultRendererDiagnosis(t *testing.T) {
	r := view.DefaultRenderer{}
	out := r.Diagnosis(profiles.Evaluation{OK: true}, "f.mp4")
	if !strings.Contains(out, "Compatible") {
		t.Error("expected Compatible in Diagnosis render")
	}
}

func TestDefaultRendererProfile(t *testing.T) {
	r := view.DefaultRenderer{}
	out := r.Profile()
	if !strings.Contains(out, "Resolve") {
		t.Error("expected Resolve in Profile render")
	}
}

func TestDefaultRendererConverting(t *testing.T) {
	r := view.DefaultRenderer{}
	out := r.Converting(progress.New())
	if !strings.Contains(out, "Converting") {
		t.Error("expected Converting in Converting render")
	}
}

func TestDefaultRendererResultSuccess(t *testing.T) {
	r := view.DefaultRenderer{}
	out := r.Result("out.mp4", nil)
	if !strings.Contains(out, "Done") {
		t.Error("expected Done in Result render")
	}
}

func TestDefaultRendererResultError(t *testing.T) {
	r := view.DefaultRenderer{}
	out := r.Result("", errors.New("boom"))
	if !strings.Contains(out, "Error") {
		t.Error("expected Error in Result render")
	}
}
