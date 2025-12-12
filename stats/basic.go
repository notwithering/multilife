package stats

import (
	"fmt"
	"strconv"
	"time"
)

func (s *StatsPrinter) ShouldBasic() bool {
	return s.config.Basic.Enabled && s.currentFrame%s.config.Basic.Interval == 0
}

func (s *StatsPrinter) basicStatsText() string {
	var text string

	// Frame: 757/1800 (42.1%)
	text += "Frame: "
	text += strconv.Itoa(s.currentFrame+1) + "/" + strconv.Itoa(s.config.Basic.TotalFrames)
	text += " (" + fmt.Sprintf("%.1f", float32(s.currentFrame+1)/float32(s.config.Basic.TotalFrames)*100) + "%)"
	text += "\n"

	// Render: 2.2398ms (500/s)
	text += "Render: "
	renderTime := s.renderEnd.Sub(s.renderStart)
	text += renderTime.String()
	text += " (" + fmt.Sprintf("%.2f", 1.0/renderTime.Seconds()) + "/s)"
	text += "\n"

	// UI: 50us (30000/s)
	text += "UI: "
	uiTime := s.uiEnd.Sub(s.uiStart)
	text += uiTime.String()
	text += " (" + fmt.Sprintf("%.2f", 1.0/uiTime.Seconds()) + "/s)"
	text += "\n"

	// Step: 30.2155342ms (120/s)
	text += "Step: "
	stepTime := s.stepEnd.Sub(s.stepStart)
	text += stepTime.String()
	text += " (" + fmt.Sprintf("%.2f", 1.0/stepTime.Seconds()) + "/s)"
	text += "\n"

	// Frame: 50.2155342ms (20/s)
	text += "Frame: "
	frameTime := s.frameEnd.Sub(s.frameStart)
	text += frameTime.String()
	text += " (" + fmt.Sprintf("%.2f", 1.0/frameTime.Seconds()) + "/s)"
	text += "\n"

	// Elapsed: 0m12s502ms
	text += "Elapsed: "
	elapsedTime := time.Since(s.loopStart)
	text += durationToString(elapsedTime)
	text += "\n"

	// Estimated: 0m50s324ms
	text += "Estimated: "
	progress := float64(s.currentFrame) / float64(s.config.Basic.TotalFrames)
	if progress == 0 {
		text += "N/A"
	} else {
		estimatedTimeLeft := float64(elapsedTime.Seconds())/progress - float64(elapsedTime.Seconds())
		text += durationToString(time.Duration(estimatedTimeLeft * float64(time.Second)))
	}
	text += "\n"

	return text
}

func (s *StatsPrinter) StartedLoop() {
	s.loopStart = time.Now()
}
func (s *StatsPrinter) EndedLoop() {
	s.loopEnd = time.Now()
}

func (s *StatsPrinter) StartedFrame() {
	s.currentFrame++
	if s.ShouldBasic() {
		s.frameStart = time.Now()
	}
}
func (s *StatsPrinter) EndedFrame() {
	if s.ShouldBasic() {
		s.frameEnd = time.Now()
	}
}

func (s *StatsPrinter) StartedRender() {
	if s.ShouldBasic() {
		s.renderStart = time.Now()
	}
}
func (s *StatsPrinter) EndedRender() {
	if s.ShouldBasic() {
		s.renderEnd = time.Now()
	}
}

func (s *StatsPrinter) StartedUI() {
	if s.ShouldBasic() {
		s.uiStart = time.Now()
	}
}
func (s *StatsPrinter) EndedUI() {
	if s.ShouldBasic() {
		s.uiEnd = time.Now()
	}
}

func (s *StatsPrinter) StartedStep() {
	if s.ShouldBasic() {
		s.stepStart = time.Now()
	}
}
func (s *StatsPrinter) EndedStep() {
	if s.ShouldBasic() {
		s.stepEnd = time.Now()
	}
}
