package stats

import (
	"fmt"
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

	var averageDensity float32
	for _, population := range s.ecosystemStats.PopulationBySpecie {
		averageDensity += float32(population) / float32(s.ecosystemStats.TotalPopulation)
	}
	averageDensity = averageDensity / float32(len(s.species)) * 100

	for specieId, population := range s.ecosystemStats.PopulationBySpecie {
		specie := s.species[specieId]
		density := float32(population) / float32(s.ecosystemStats.TotalPopulation) * 100

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
		text += strconv.Itoa(population) + "/" + strconv.Itoa(s.ecosystemStats.TotalPopulation)
		text += " (" + fmt.Sprintf("%.1f", density) + "%)"
		text += sgr.Reset
		text += "\n"
	}

	return text
}
