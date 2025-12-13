package stats

import (
	"fmt"
	"time"

	"github.com/notwithering/multilife/ecosystem"
	"github.com/notwithering/multilife/specie"
)

type StatsPrinter struct {
	config         Config
	species        []*specie.CompiledSpecie
	ecosystemStats ecosystem.Stats

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

func (s *StatsPrinter) Print() {
	var text string

	if s.config.Basic.Enabled {
		text += s.basicStatsText()
	}

	if s.config.Basic.Enabled && s.config.Ecosystem.Enabled {
		text += "----------\n"
	}

	if s.config.Ecosystem.Enabled {
		text += s.ecosystemStatsText()
		text += "----------\n"
	}

	text += "Ctrl+C to finish.\n"

	if s.config.Basic.Enabled || s.config.Ecosystem.Enabled {
		fmt.Print("\x1b[H\x1b[J" + text)
	}
}

func (s *StatsPrinter) PrintClosure() {

}
