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
	var printFuncs []func(*strings.Builder)

	if s.config.Basic.Enabled {
		printFuncs = append(printFuncs, s.writeBasicStats)
	}
	if s.config.Timings.Enabled {
		printFuncs = append(printFuncs, s.writeTimingStats)
	}
	if s.config.Ecosystem.Enabled {
		printFuncs = append(printFuncs, s.writeEcosystemStats)
	}
	printFuncs = append(printFuncs, func(sb *strings.Builder) {
		sb.WriteString("Ctrl+C to finish.\n")
	})

	var sb strings.Builder
	sb.WriteString("\x1b[H\x1b[J")

	for i, fn := range printFuncs {
		if i != 0 {
			sb.WriteString("----------\n")
		}
		fn(&sb)
	}

	fmt.Fprint(os.Stderr, sb.String())
}

func (s *StatsPrinter) PrintClosure() {

}
