package stats

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (s *StatsPrinter) ShouldBasic() bool {
	return s.config.Basic.Enabled && s.currentFrame%s.config.Basic.Interval == 0
}

func (s *StatsPrinter) basicStatsText(sb *strings.Builder) {
	// Frame: 757/1800 (42.1%)
	sb.WriteString("Frame: ")
	if !s.config.Infinite {
		sb.WriteString(strconv.Itoa(s.currentFrame + 1))
		sb.WriteByte('/')
		sb.WriteString(strconv.Itoa(s.config.Basic.TotalFrames))
		sb.WriteString(" (")
		fmt.Fprintf(sb, "%.1f", float32(s.currentFrame+1)/float32(s.config.Basic.TotalFrames)*100)
		sb.WriteString("%)")
	} else {
		sb.WriteString(strconv.Itoa(s.currentFrame + 1))
	}
	sb.WriteByte('\n')

	// Render: 2.2398ms (500/s)
	sb.WriteString("Render: ")
	renderTime := s.renderEnd.Sub(s.renderStart)
	sb.WriteString(renderTime.String())
	sb.WriteString(" (")
	fmt.Fprintf(sb, "%.2f", 1.0/renderTime.Seconds())
	sb.WriteString("/s)")
	sb.WriteByte('\n')

	// UI: 50us (30000/s)
	sb.WriteString("UI: ")
	uiTime := s.uiEnd.Sub(s.uiStart)
	sb.WriteString(uiTime.String())
	sb.WriteString(" (")
	fmt.Fprintf(sb, "%.2f", 1.0/uiTime.Seconds())
	sb.WriteString("/s)")
	sb.WriteByte('\n')

	// Step: 30.2155342ms (120/s)
	sb.WriteString("Step: ")
	stepTime := s.stepEnd.Sub(s.stepStart)
	sb.WriteString(stepTime.String())
	sb.WriteString(" (")
	fmt.Fprintf(sb, "%.2f", 1.0/stepTime.Seconds())
	sb.WriteString("/s)")
	sb.WriteByte('\n')

	// Frame: 50.2155342ms (20/s)
	sb.WriteString("Frame: ")
	frameTime := s.frameEnd.Sub(s.frameStart)
	sb.WriteString(frameTime.String())
	sb.WriteString(" (")
	fmt.Fprintf(sb, "%.2f", 1.0/frameTime.Seconds())
	sb.WriteString("/s)")
	sb.WriteByte('\n')

	// Elapsed: 0m12s502ms (1m4s102ms)
	sb.WriteString("Elapsed: ")
	elapsedTime := time.Since(s.loopStart)
	sb.WriteString(durationToString(elapsedTime))
	sb.WriteString(" (")
	outputTime := time.Duration(float64(s.currentFrame) / float64(s.config.Basic.FPS) * float64(time.Second))
	sb.WriteString(durationToString(outputTime))
	sb.WriteByte(')')
	sb.WriteByte('\n')

	// Estimated: 0m50s324ms
	if !s.config.Infinite {
		sb.WriteString("Estimated: ")
		progress := float64(s.currentFrame) / float64(s.config.Basic.TotalFrames)
		if progress == 0 {
			sb.WriteString("N/A")
		} else {
			estimatedTimeLeft := float64(elapsedTime.Seconds())/progress - float64(elapsedTime.Seconds())
			sb.WriteString(durationToString(time.Duration(estimatedTimeLeft * float64(time.Second))))
		}
		sb.WriteByte('\n')
	}
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
