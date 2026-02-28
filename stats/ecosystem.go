package stats

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"github.com/notwithering/multilife/ecosystem"
	"github.com/notwithering/sgr"
)

func (s *StatsPrinter) ShouldEcosystem() bool {
	return s.config.Ecosystem.Enabled && s.currentFrame%s.config.Ecosystem.Interval == 0
}

func (s *StatsPrinter) UpdateEcosystemStats(stats ecosystem.Stats) {
	s.ecosystemStats = stats
}

func (s *StatsPrinter) ecosystemStatsText(sb *strings.Builder) {
	// sb.WriteString("warning: ecosystem stats can add up to 100us/frame (+1s/1,000,000 frames)")
	// sb.WriteByte('\n')

	var averageDensity float64
	for _, population := range s.ecosystemStats.PopulationBySpecie {
		averageDensity += float64(population) / float64(s.ecosystemStats.TotalPopulation)
	}
	averageDensity = averageDensity / float64(len(s.species))

	for specieId, population := range s.ecosystemStats.PopulationBySpecie {
		specie := s.species[specieId]
		density := float64(population) / float64(s.ecosystemStats.TotalPopulation)

		c := gradient(
			density/averageDensity,
			color.RGBA{255, 0, 0, 255},
			color.RGBA{255, 255, 0, 255},
			color.RGBA{0, 255, 0, 255},
		).(color.RGBA)

		if density == 0 {
			sb.WriteString(sgr.Strike)
		}

		fmt.Fprintf(sb, "%s%d;%d;%dm", sgr.FgColorRGB, c.R, c.G, c.B)

		// Conway's Life: 50423/99256 (50.8%)
		sb.WriteString(specie.Name + ": ")
		sb.WriteString(strconv.Itoa(population) + "/" + strconv.Itoa(s.ecosystemStats.TotalPopulation))
		sb.WriteString(" (")
		fmt.Fprintf(sb, "%.1f", density*100)
		sb.WriteString("%)")
		sb.WriteString(sgr.Reset)
		sb.WriteByte('\n')
	}
}
