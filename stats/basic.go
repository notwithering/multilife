package stats

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (s *StatsPrinter) writeBasicStats(sb *strings.Builder) {
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
