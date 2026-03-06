package stats

import (
	"fmt"
	"strings"
	"time"
)

func (s *StatsPrinter) ShouldTimings() bool {
	return s.config.Timings.Enabled && s.currentFrame%s.config.Timings.Interval == 0
}

// Render: 2.2398ms (500/s)
func writeTiming(sb *strings.Builder, name string, start, end time.Time) {
	timing := end.Sub(start)

	sb.WriteString(name)
	sb.WriteString(": ")
	sb.WriteString(timing.String())
	sb.WriteString(" (")
	fmt.Fprintf(sb, "%.2f", 1.0/timing.Seconds())
	sb.WriteString("/s)")
	sb.WriteByte('\n')
}

func (s *StatsPrinter) writeTimingStats(sb *strings.Builder) {
	writeTiming(sb, "Render", s.renderStart, s.uiEnd)
	writeTiming(sb, "UI", s.uiStart, s.uiEnd)
	writeTiming(sb, "Step", s.stepStart, s.stepEnd)
	writeTiming(sb, "Frame", s.frameStart, s.frameEnd)
}

func (s *StatsPrinter) StartedFrame() {
	s.currentFrame++
	if s.ShouldTimings() {
		s.frameStart = time.Now()
	}
}
func (s *StatsPrinter) EndedFrame() {
	if s.ShouldTimings() {
		s.frameEnd = time.Now()
	}
}

func (s *StatsPrinter) StartedRender() {
	if s.ShouldTimings() {
		s.renderStart = time.Now()
	}
}
func (s *StatsPrinter) EndedRender() {
	if s.ShouldTimings() {
		s.renderEnd = time.Now()
	}
}

func (s *StatsPrinter) StartedUI() {
	if s.ShouldTimings() {
		s.uiStart = time.Now()
	}
}
func (s *StatsPrinter) EndedUI() {
	if s.ShouldTimings() {
		s.uiEnd = time.Now()
	}
}

func (s *StatsPrinter) StartedStep() {
	if s.ShouldTimings() {
		s.stepStart = time.Now()
	}
}
func (s *StatsPrinter) EndedStep() {
	if s.ShouldTimings() {
		s.stepEnd = time.Now()
	}
}
