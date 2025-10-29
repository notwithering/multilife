package main

import (
	"fmt"
	"io"
	"os/exec"
)

const (
	videoFPS    int = 120
	videoLength int = 20 //seconds
	videoFrames int = videoLength * videoFPS
	videoWidth  int = 1920
	videoHeight int = 1080
)

type renderer struct {
	cmd  *exec.Cmd
	pipe io.WriteCloser
}

func newRenderer() (*renderer, error) {
	var r = &renderer{}

	cmd := exec.Command("ffmpeg",
		"-y",
		"-f", "rawvideo",
		"-pixel_format", "rgb24",
		"-video_size", fmt.Sprintf("%dx%d", ecosystemSizeX, ecosystemSizeY),
		"-framerate", fmt.Sprintf("%d", videoFPS),
		"-i", "-",
		"-vf", fmt.Sprintf("scale=%d:%d:flags=neighbor", videoWidth, videoHeight),
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-preset", "slow",

		"output.mp4",
	)
	r.cmd = cmd

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	r.pipe = stdin

	return r, nil
}

func (r *renderer) start() error {
	return r.cmd.Start()
}
