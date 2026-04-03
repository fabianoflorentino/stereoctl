package model

import (
	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fabianoflorentino/stereoctl/internal/ffmpeg"
	"github.com/fabianoflorentino/stereoctl/internal/profiles"
	"github.com/fabianoflorentino/stereoctl/tui/service"
)

// --- messages ---

// ProbeCompletedMsg is sent when ffprobe finishes successfully.
type ProbeCompletedMsg struct{ Result *ffmpeg.ProbeOutput }

// ProbeErrorMsg is sent when ffprobe returns an error.
type ProbeErrorMsg struct{ Err error }

// ConvertProgressMsg is sent periodically during conversion with the current
// percentage (0.0–1.0).
type ConvertProgressMsg struct{ Percent float64 }

// ConvertDoneMsg is sent when ffmpeg conversion finishes successfully.
type ConvertDoneMsg struct{ Output string }

// ConvertErrorMsg is sent when ffmpeg conversion returns an error.
type ConvertErrorMsg struct{ Err error }

// --- Renderer ---

// Renderer builds the string shown for each state. Injected so that the model
// remains testable without depending on lipgloss rendering.
type Renderer interface {
	FileSelection(fp filepicker.Model) string
	Diagnosis(ev profiles.Evaluation, file string) string
	Profile() string
	Converting(bar progress.Model) string
	Result(output string, err error) string
}

// --- AppModel ---

// AppModel is the root Bubble Tea model for the TUI.
type AppModel struct {
	prober    service.Prober
	converter service.Converter
	renderer  Renderer
	state     State

	fp          filepicker.Model
	progressBar progress.Model

	// progressCh receives ConvertProgressMsg / ConvertDoneMsg / ConvertErrorMsg
	// from the conversion goroutine. Channels are reference types so copies of
	// AppModel all share the same underlying channel.
	progressCh <-chan tea.Msg

	selectedFile string
	probe        *ffmpeg.ProbeOutput
	evaluation   profiles.Evaluation
	output       string
	err          error
}

// New creates a new AppModel.
func New(prober service.Prober, converter service.Converter, renderer Renderer) AppModel {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".mp4", ".mkv", ".mov", ".avi", ".mxf"}

	return AppModel{
		prober:      prober,
		converter:   converter,
		renderer:    renderer,
		state:       StateFileSelection,
		fp:          fp,
		progressBar: progress.New(progress.WithDefaultGradient()),
	}
}

// Init implements tea.Model.
func (m AppModel) Init() tea.Cmd {
	return m.fp.Init()
}

// Update implements tea.Model.
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}

	case ProbeCompletedMsg:
		m.probe = msg.Result
		m.evaluation = profiles.EvaluateResolveFree(m.probe)
		m.state = StateDiagnosis
		return m, nil

	case ProbeErrorMsg:
		m.err = msg.Err
		m.state = StateResult
		return m, nil

	case ConvertProgressMsg:
		pbCmd := m.progressBar.SetPercent(msg.Percent)
		return m, tea.Batch(pbCmd, readProgress(m.progressCh))

	case ConvertDoneMsg:
		m.output = msg.Output
		m.state = StateResult
		return m, nil

	case ConvertErrorMsg:
		m.err = msg.Err
		m.state = StateResult
		return m, nil
	}

	switch m.state {
	case StateFileSelection:
		return m.updateFileSelection(msg)
	case StateDiagnosis:
		return m.updateDiagnosis(msg)
	case StateProfileSelection:
		return m.updateProfileSelection(msg)
	case StateConverting:
		newPB, cmd := m.progressBar.Update(msg)
		if pb, ok := newPB.(progress.Model); ok {
			m.progressBar = pb
		}
		return m, cmd
	}

	return m, nil
}

// View implements tea.Model.
func (m AppModel) View() string {
	if m.renderer == nil {
		return ""
	}
	switch m.state {
	case StateFileSelection:
		return m.renderer.FileSelection(m.fp)
	case StateDiagnosis:
		return m.renderer.Diagnosis(m.evaluation, m.selectedFile)
	case StateProfileSelection:
		return m.renderer.Profile()
	case StateConverting:
		return m.renderer.Converting(m.progressBar)
	case StateResult:
		return m.renderer.Result(m.output, m.err)
	}
	return ""
}

// --- accessors for tests ---

// CurrentState returns the current State (exported for tests).
func (m AppModel) CurrentState() State { return m.state }

// SelectedFile returns the path of the currently selected file.
func (m AppModel) SelectedFile() string { return m.selectedFile }

// ProbeCmd returns the probe command for path. Exported so tests can execute it directly.
func (m AppModel) ProbeCmd(path string) tea.Cmd { return m.probeCmd(path) }

// StartConvert initialises the progress channel, starts the conversion goroutine
// and returns the updated model together with the first read command.
// Exported so tests can exercise the full streaming flow.
func (m AppModel) StartConvert() (AppModel, tea.Cmd) { return m.startConvert() }

// ProbeResult returns the last ProbeOutput received.
func (m AppModel) ProbeResult() *ffmpeg.ProbeOutput { return m.probe }

// Err returns the last error recorded.
func (m AppModel) Err() error { return m.err }

// Output returns the output path after a successful conversion.
func (m AppModel) Output() string { return m.output }

// Evaluation returns the last profile evaluation result.
func (m AppModel) Evaluation() profiles.Evaluation { return m.evaluation }

// --- private helpers ---

func (m AppModel) updateFileSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.fp, cmd = m.fp.Update(msg)
	if didSelect, path := m.fp.DidSelectFile(msg); didSelect {
		m.selectedFile = path
		return m, tea.Batch(cmd, m.probeCmd(path))
	}
	return m, cmd
}

func (m AppModel) probeCmd(path string) tea.Cmd {
	return func() tea.Msg {
		result, err := m.prober.Probe(path)
		if err != nil {
			return ProbeErrorMsg{Err: err}
		}
		return ProbeCompletedMsg{Result: result}
	}
}

func (m AppModel) updateDiagnosis(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch key.String() {
	case "enter", "c":
		m.state = StateConverting
		return m.startConvert()
	case "p":
		m.state = StateProfileSelection
	case "esc":
		m.state = StateFileSelection
	}
	return m, nil
}

func (m AppModel) updateProfileSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	if key.String() == "enter" {
		m.state = StateConverting
		return m.startConvert()
	}
	return m, nil
}

// startConvert creates a buffered channel, starts the conversion goroutine and
// returns an updated model (with progressCh set) plus the initial read command.
func (m AppModel) startConvert() (AppModel, tea.Cmd) {
	opts := profiles.BuildConvertOptionsForResolveFree(m.probe, m.selectedFile)
	dur := m.probe.Format.Duration
	out := opts.Output

	ch := make(chan tea.Msg, 64)
	m.progressCh = ch

	conv := m.converter
	go func() {
		err := conv.ConvertWithProgress(opts, dur, func(pct float64) {
			ch <- ConvertProgressMsg{Percent: pct}
		})
		if err != nil {
			ch <- ConvertErrorMsg{Err: err}
		} else {
			ch <- ConvertDoneMsg{Output: out}
		}
	}()

	return m, readProgress(ch)
}

// readProgress returns a Cmd that blocks until the next message arrives on ch.
func readProgress(ch <-chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return <-ch
	}
}
