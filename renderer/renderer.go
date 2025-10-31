package renderer

import (
	"fmt"
	"io"
	"main/gfx"
	"os/exec"
)

type renderer struct {
	cmd  *exec.Cmd
	pipe io.WriteCloser
}

func NewRenderer(config Config) *renderer {
	r := &renderer{}

	cmd := exec.Command("ffmpeg",
		"-y",
		"-f", "rawvideo",
		"-pixel_format", "rgb24",
		"-video_size", fmt.Sprintf("%dx%d", config.Video.SourceWidth, config.Video.SourceHeight),
		"-framerate", fmt.Sprintf("%d", config.Video.FPS),
		"-i", "-",
		"-vf", fmt.Sprintf("scale=%d:%d:flags=neighbor", config.Video.OutputWidth, config.Video.OutputHeight),
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-preset", "slow",

		"output.mp4",
	)
	r.cmd = cmd

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil
	}
	r.pipe = stdin

	return r
}

func (r *renderer) Start() error {
	return r.cmd.Start()
}

func (r *renderer) Write(buf *gfx.Buffer) {
	r.pipe.Write(buf.Data)
}
