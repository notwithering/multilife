package stats

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/notwithering/multilife/ecosystem"
	"github.com/notwithering/sgr"
)

func (s *StatsPrinter) ShouldEcosystem() bool {
	return s.config.Ecosystem.Enabled && s.currentFrame%s.config.Ecosystem.Interval == 0
}

func (s *StatsPrinter) UpdateEcosystemStats(stats ecosystem.Stats) {
	s.ecosystemStats = stats
}

func (s *StatsPrinter) ecosystemStatsText() string {
	var text string

	// text += "warning: ecosystem stats can add up to 100us/frame (+1s/1,000,000 frames)"
	// text += "\n"

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
			text += sgr.Strike
		}

		text += fmt.Sprintf("%s%d;%d;%dm", sgr.FgColorRGB, c.R, c.G, c.B)

		// Conway's Life: 50423/99256 (50.8%)
		text += specie.Name + ": "
		text += strconv.Itoa(population) + "/" + strconv.Itoa(s.ecosystemStats.TotalPopulation)
		text += " (" + fmt.Sprintf("%.1f", density*100) + "%)"
		text += sgr.Reset
		text += "\n"
	}

	return text
}
