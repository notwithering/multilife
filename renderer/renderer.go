package renderer

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/notwithering/multilife/gfx"
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
		"-c:v", "ffv1",
		"-crf", "0",
		"-preset", "slow",
		"-pix_fmt", "yuv444p",
		config.Video.OutputFile,
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
