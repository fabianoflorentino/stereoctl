package ffmpeg

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
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

	cmd := exec.Command("ffmpeg", args...)
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return fmt.Errorf("stdout pipe: %w", err)
	}

	cmd.Stderr = io.Discard

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start ffmpeg: %w", err)
	}

	// use MonitorProgress to allow testing the progress parsing/formatting
	// write progress to stderr so it doesn't interfere with stdout capture
	MonitorProgress(stdout, totalUs, io.Discard)

	return cmd.Wait()
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
				out.Write(buf.Bytes())
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

	fmt.Fprintf(out, "\r[%s] %5.1f%%  %s / %s   ", bar, pct*100, formatDur(elapsed), formatDur(totalDur))
}

// formatDur converts a time.Duration to a string in HH:MM:SS format.
func formatDur(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60

	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
