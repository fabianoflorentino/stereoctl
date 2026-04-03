package ffmpeg

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type ConvertOptions struct {
	Input      string
	Output     string
	CopyAudio  bool
	AudioCodec string
	Channels   int
	Bitrate    string
	// VideoCodec: if empty, copy video stream; otherwise re-encode using this codec (e.g. libx264)
	VideoCodec string
	// VideoExtraArgs: additional ffmpeg args for video encoding (e.g. -preset slow -crf 18)
	VideoExtraArgs []string
	// ExtraArgs: additional ffmpeg arguments to append before progress (e.g. -t 30)
	ExtraArgs []string
}

// Convert runs ffmpeg with the specified options and displays a progress bar based on the duration.
// It expects the duration in seconds as a string (from ffprobe) to calculate progress.
func Convert(opts ConvertOptions, durationStr string) error {
	totalSec, _ := strconv.ParseFloat(durationStr, 64)
	totalUs := int64(totalSec * 1_000_000)

	args := BuildFFmpegArgs(opts)

	stdout, waitFn, err := ffmpegCmdRunner(args)
	if err != nil {
		return fmt.Errorf("start ffmpeg: %w", err)
	}

	// use MonitorProgress to allow testing the progress parsing/formatting
	// write progress to stderr so it doesn't interfere with stdout capture
	MonitorProgress(stdout, totalUs, io.Discard)

	return waitFn()
}

// BuildFFmpegArgs builds the ffmpeg CLI arguments for the given ConvertOptions.
// This is separated out for easier unit testing.
func BuildFFmpegArgs(opts ConvertOptions) []string {
	args := []string{"-i", opts.Input}

	// video handling
	if opts.VideoCodec == "" {
		args = append(args, "-c:v", "copy")
	} else {
		args = append(args, "-c:v", opts.VideoCodec)
		if len(opts.VideoExtraArgs) > 0 {
			args = append(args, opts.VideoExtraArgs...)
		}
	}

	// audio handling
	if opts.CopyAudio {
		args = append(args, "-c:a", "copy")
	} else {
		args = append(args, "-c:a", opts.AudioCodec, "-ac", strconv.Itoa(opts.Channels), "-b:a", opts.Bitrate)
	}

	// append any extra args then progress and output
	if len(opts.ExtraArgs) > 0 {
		args = append(args, opts.ExtraArgs...)
	}

	args = append(args, "-progress", "pipe:1", "-nostats", "-loglevel", "error", "-y", opts.Output)
	return args
}

// ffmpegCmdRunner is a package-level variable pointing to the function that
// runs ffmpeg. It is set to a real runner by default but can be overridden
// in tests to avoid executing the external binary.
var ffmpegCmdRunner = func(args []string) (io.ReadCloser, func() error, error) {
	bin := os.Getenv("FFMPEG_BIN")
	if bin == "" {
		bin = "ffmpeg"
	}
	cmd := exec.Command(bin, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	cmd.Stderr = io.Discard

	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}

	return stdout, cmd.Wait, nil
}

// printProgress displays a progress bar in the terminal based on the current
// and total duration in microseconds.
// MonitorProgress reads ffmpeg -progress output from r and writes a simple
// progress bar to the provided writer. This function is testable by passing
// a custom reader and writer.
func MonitorProgress(r io.Reader, totalUs int64, out io.Writer) {
	scanner := bufio.NewScanner(r)
	var buf bytes.Buffer
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "out_time_us=") {
			val := strings.TrimPrefix(line, "out_time_us=")
			currentUs, err := strconv.ParseInt(val, 10, 64)
			if err == nil && currentUs >= 0 {
				// accumulate a single-line buffer and write to out
				buf.Reset()
				writeProgress(&buf, currentUs, totalUs)
				_, _ = out.Write(buf.Bytes())
			}
		}
	}
}

func writeProgress(out io.Writer, current, total int64) {
	if total <= 0 {
		return
	}

	pct := float64(current) / float64(total)
	if pct > 1 {
		pct = 1
	}

	const barWidth = 40

	filled := int(pct * barWidth)
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	elapsed := time.Duration(current) * time.Microsecond
	totalDur := time.Duration(total) * time.Microsecond

	_, _ = fmt.Fprintf(out, "\r[%s] %5.1f%%  %s / %s   ", bar, pct*100, formatDur(elapsed), formatDur(totalDur))
}

// ConvertWithProgress runs ffmpeg and calls onProgress with percentage values
// (0.0–1.0) as conversion proceeds. onProgress may be called from a goroutine
// started by the caller; it must be goroutine-safe.
func ConvertWithProgress(opts ConvertOptions, durationStr string, onProgress func(pct float64)) error {
	totalSec, _ := strconv.ParseFloat(durationStr, 64)
	totalUs := int64(totalSec * 1_000_000)

	args := BuildFFmpegArgs(opts)
	stdout, waitFn, err := ffmpegCmdRunner(args)
	if err != nil {
		return fmt.Errorf("start ffmpeg: %w", err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "out_time_us=") {
			val := strings.TrimPrefix(line, "out_time_us=")
			currentUs, parseErr := strconv.ParseInt(val, 10, 64)
			if parseErr == nil && currentUs >= 0 && totalUs > 0 {
				pct := float64(currentUs) / float64(totalUs)
				if pct > 1.0 {
					pct = 1.0
				}
				onProgress(pct)
			}
		}
	}

	return waitFn()
}

// formatDur converts a time.Duration to a string in HH:MM:SS format.
func formatDur(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60

	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
