package stats

import (
	"fmt"
	"os"
	"strings"
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
	if !s.config.Basic.Enabled && !s.config.Ecosystem.Enabled {
		return
	}

	var sb strings.Builder
	sb.WriteString("\x1b[H\x1b[J")

	if s.config.Basic.Enabled {
		s.basicStatsText(&sb)
	}

	if s.config.Basic.Enabled && s.config.Ecosystem.Enabled {
		sb.WriteString("----------\n")
	}

	if s.config.Ecosystem.Enabled {
		s.ecosystemStatsText(&sb)
		sb.WriteString("----------\n")
	}

	sb.WriteString("Ctrl+C to finish.\n")
	fmt.Fprint(os.Stderr, sb.String())
}

func (s *StatsPrinter) PrintClosure() {

}
