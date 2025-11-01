package stats

import (
	"fmt"
	"strconv"
	"time"

	"github.com/notwithering/multilife/ecosystem"
	"github.com/notwithering/multilife/specie"
	"github.com/notwithering/sgr"
)

type StatsPrinter struct {
	config  Config
	species []*specie.CompiledSpecie

	loopStart time.Time
	loopEnd   time.Time

	frameStart time.Time
	frameEnd   time.Time

	currentFrame int

	renderStart time.Time
	renderEnd   time.Time

	uiStart time.Time
	uiEnd   time.Time

	stepStart time.Time
	stepEnd   time.Time
}

func NewStatsPrinter(config Config, species []*specie.CompiledSpecie) *StatsPrinter {
	return &StatsPrinter{
		config:  config,
		species: species,
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

func (s *StatsPrinter) Print(stats ecosystem.Stats) {
	var text string

	if s.config.Basic.Enabled {
		text += s.basicStats()
	}

	if s.config.Basic.Enabled && s.config.Ecosystem.Enabled {
		text += "----------\n"
	}

	if s.config.Ecosystem.Enabled {
		text += s.ecosystemStats(stats)
	}

	if s.config.Basic.Enabled || s.config.Ecosystem.Enabled {
		fmt.Print("\x1b[H\x1b[J" + text)
	}
}

func durationToString(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	milliseconds := int(d.Milliseconds()) % 1000

	return fmt.Sprintf("%dm%ds%03dms", minutes, seconds, milliseconds)
}

func (s *StatsPrinter) basicStats() string {
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

func (s *StatsPrinter) ecosystemStats(stats ecosystem.Stats) string {
	var text string

	// text += "warning: ecosystem stats can add up to 100us/frame (+1s/1,000,000 frames)"
	// text += "\n"

	var averageDensity float32
	for _, population := range stats.PopulationBySpecie {
		averageDensity += float32(population) / float32(stats.TotalPopulation)
	}
	averageDensity = averageDensity / float32(len(s.species)) * 100

	for specieId, population := range stats.PopulationBySpecie {
		specie := s.species[specieId]
		density := float32(population) / float32(stats.TotalPopulation) * 100

		if density == 0 {
			text += sgr.FgRed + sgr.Strike
		} else if density < averageDensity/2 {
			text += sgr.FgYellow
		} else if density < averageDensity/1.5 {
			text += sgr.FgHiYellow
		} else {
			text += sgr.FgGreen
		}

		// Conway's Life: 50423/99256 (50.8%)
		text += specie.Name + ": "
		text += strconv.Itoa(population) + "/" + strconv.Itoa(stats.TotalPopulation)
		text += " (" + fmt.Sprintf("%.1f", density) + "%)"
		text += sgr.Reset
		text += "\n"
	}

	// SPS: 245.66 (4.0707ms/step)
	// sps := 1.0 / stats.FrameTime.Seconds()
	// text += "SPS: "
	// text += fmt.Sprintf("%.2f", sps)
	// text += " (" + stats.FrameTime.String() + "/step)"
	// text += "\n"

	return text
}

func (s *StatsPrinter) PrintClosure() {

}

func (s *StatsPrinter) ShouldBasic() bool {
	return s.config.Basic.Enabled && s.currentFrame%s.config.Basic.Interval == 0
}

func (s *StatsPrinter) ShouldEcosystem() bool {
	return s.config.Ecosystem.Enabled && s.currentFrame%s.config.Ecosystem.Interval == 0
}
