package hls

import (
	"context"
	"os"
	"os/exec"
	"path"
)

func Transcode(_ context.Context, input, outputDir string) error {
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			return err
		}
	}

	args := []string{
		"-i", input,
		"-c:a", "aac",
		"-ar", "48000",
		"-b:a", "192k",
		"-c:v", "h264", "-profile:v", "main",
		"-b:v", "5000k", "-maxrate", "5350k", "-bufsize", "7500k",
		"-crf", "20",
		"-g", "48", "-keyint_min", "48", "-sc_threshold", "0",
		"-hls_time", "4", "-hls_playlist_type", "vod",
		"-hls_segment_filename", path.Join(outputDir, "1080p_%03d.ts"),
		path.Join(outputDir, "1080p.m3u8"),
	}

	cmd := exec.Command("ffmpeg", args...)
	return cmd.Run()
}
