package ffmpeg

import (
	"bufio"
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
}

// Convert runs ffmpeg with the specified options and displays a progress bar based on the duration.
// It expects the duration in seconds as a string (from ffprobe) to calculate progress.
func Convert(opts ConvertOptions, durationStr string) error {
	totalSec, _ := strconv.ParseFloat(durationStr, 64)
	totalUs := int64(totalSec * 1_000_000)

	args := []string{"-i", opts.Input, "-c:v", "copy"}

	if opts.CopyAudio {
		args = append(args, "-c:a", "copy")
	} else {
		args = append(args, "-c:a", opts.AudioCodec, "-ac", strconv.Itoa(opts.Channels), "-b:a", opts.Bitrate)
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

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		line := scanner.Text()

		if after, ok := strings.CutPrefix(line, "out_time_us="); ok {
			val := after
			currentUs, err := strconv.ParseInt(val, 10, 64)
			if err == nil && currentUs >= 0 {
				printProgress(currentUs, totalUs)
			}
		}
	}

	return cmd.Wait()
}

// printProgress displays a progress bar in the terminal based on the current
// and total duration in microseconds.
func printProgress(current, total int64) {
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

	fmt.Printf("\r[%s] %5.1f%%  %s / %s   ", bar, pct*100, formatDur(elapsed), formatDur(totalDur))
}

// formatDur converts a time.Duration to a string in HH:MM:SS format.
func formatDur(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60

	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
