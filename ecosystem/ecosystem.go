package ecosystem

import (
	"math"
	"runtime"
	"sync"

	"github.com/notwithering/multilife/gfx"
	"github.com/notwithering/multilife/rng"
	"github.com/notwithering/multilife/specie"
)

type Ecosystem struct {
	config          Config
	Species         []*specie.CompiledSpecie
	world           []specie.SpecieId
	nextWorld       []specie.SpecieId
	Stats           Stats
	neighborIndices [][8]int
}

type Stats struct {
	TotalPopulation    int
	PopulationBySpecie []int
}

func (e *Ecosystem) Render(buf *gfx.Buffer) {
	for y := range e.config.Height {
		idx := y * e.config.Width * 3
		for x := range e.config.Width {
			specieId := e.world[e.index(x, y)]

			col := e.config.Render.BackgroundColor
			if specieId != specie.NoSpecie {
				col = e.Species[specieId].Color
			}

			r, g, b, _ := col.RGBA()
			buf.Data[idx+0] = uint8(r >> 8)
			buf.Data[idx+1] = uint8(g >> 8)
			buf.Data[idx+2] = uint8(b >> 8)

			idx += 3
		}
	}
}

var neighborhood = [8][2]int{
	{-1, -1}, {0, -1}, {1, -1},
	{-1, 0} /*     */, {1, 0},
	{-1, 1}, {0, 1}, {1, 1},
}

func NewEcosystem(config Config, species []*specie.CompiledSpecie) Ecosystem {
	var eco Ecosystem

	if len(species) >= 255 {
		species = species[:255]
	}

	eco.config = config
	eco.Species = species
	eco.world = newEmptyWorld(config.Width, config.Height)
	eco.nextWorld = newEmptyWorld(config.Width, config.Height)

	columns, rows := computeGrid(len(species))

	var regionWidth int = config.Width / columns
	var regionHeight int = config.Height / rows

	eco.neighborIndices = make([][8]int, config.Width*config.Height)

	for x := range config.Width {
		for y := range config.Height {
			var indices [8]int

			for i, offset := range neighborhood {
				neighborX := (x + offset[0] + eco.config.Width) % eco.config.Width
				neighborY := (y + offset[1] + eco.config.Height) % eco.config.Height
				indices[i] = eco.index(neighborX, neighborY)
			}

			eco.neighborIndices[eco.index(x, y)] = indices
		}
	}

	for specieIndex, specie := range species {
		column := specieIndex % columns
		row := specieIndex / columns

		xStart := column*regionWidth + config.Region.Padding
		xEnd := (column+1)*regionWidth - config.Region.Padding
		yStart := row*regionHeight + config.Region.Padding
		yEnd := (row+1)*regionHeight - config.Region.Padding

		for x := xStart; x < xEnd; x++ {
			for y := yStart; y < yEnd; y++ {
				if rng.Rand.Intn(100) < config.Region.Density {
					eco.world[eco.index(x, y)] = specie.Id
				}
			}
		}
	}

	return eco
}

func newEmptyWorld(width, height int) []specie.SpecieId {
	var world = make([]specie.SpecieId, width*height)
	for i := range world {
		world[i] = specie.NoSpecie
	}
	return world
}

func computeGrid(n int) (rows, columns int) {
	if n <= 0 {
		return 0, 0
	}

	for columns = int(math.Ceil(math.Sqrt(float64(n)))); columns > 0; columns-- {
		if n%columns == 0 {
			rows = n / columns
			return
		}
	}

	return
}

func (e *Ecosystem) Step(collectStats bool) {
	numWorkers := max(runtime.NumCPU(), e.config.Height)
	rowsPerWorker := e.config.Height / numWorkers

	if collectStats {
		e.Stats.TotalPopulation = 0
		e.Stats.PopulationBySpecie = make([]int, len(e.Species))
	}

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for worker := range numWorkers {
		startY := worker * rowsPerWorker
		endY := startY + rowsPerWorker
		if worker == numWorkers-1 {
			endY = e.config.Height
		}

		go func() {
			defer wg.Done()
			workerPopulation := e.stepRange(startY, endY, collectStats)
			if collectStats {
				for cellId, count := range workerPopulation {
					e.Stats.TotalPopulation += count
					e.Stats.PopulationBySpecie[cellId] += count
				}
			}
		}()
	}

	wg.Wait()
	e.world, e.nextWorld = e.nextWorld, e.world
}

func (e *Ecosystem) stepRange(startY, endY int, collectStats bool) (workerPopulation []int) {
	width := e.config.Width
	speciesCount := len(e.Species)
	specieNeighbors := make([]int, speciesCount)
	neighborIndices := e.neighborIndices
	species := e.Species
	workerPopulation = make([]int, speciesCount)
	candidates := make([]specie.SpecieId, 0, speciesCount)

	for y := startY; y < endY; y++ {
		for x := range width {
			clear(specieNeighbors)

			cellIndex := e.index(x, y)

			var totalNeighbors int

			for _, index := range neighborIndices[cellIndex] {
				neighborSpecieId := e.world[index]
				if neighborSpecieId != specie.NoSpecie {
					specieNeighbors[neighborSpecieId]++
					totalNeighbors++
				}
			}

			cellId := e.world[cellIndex]

			var cell *specie.CompiledSpecie
			cellIsAlive := false
			cellWillLive := false

			if cellId != specie.NoSpecie {
				cell = species[cellId]
				cellIsAlive = true
				cellWillLive = cell.Rule.SurviveSet[totalNeighbors]

				if collectStats {
					workerPopulation[cellId]++
				}
			}

			candidates = candidates[:0]

			for specieIdInt, neighborsOfSpecie := range specieNeighbors {
				specieId := specie.SpecieId(specieIdInt)
				specie := species[specieId]
				shouldBirth := specie.Rule.BirthSet[neighborsOfSpecie]
				differentSpecie := specieId != cellId

				canCompete := shouldBirth && (!cellIsAlive || differentSpecie)

				if canCompete {
					candidates = append(candidates, specieId)
				}
			}

			winnerId := specie.NoSpecie
			maxWeight := -1

			for _, candidateId := range candidates {
				if candidateId == specie.NoSpecie {
					continue
				}

				neighbors := specieNeighbors[candidateId]
				if neighbors > maxWeight {
					maxWeight = neighbors
					winnerId = candidateId
				} else if neighbors == maxWeight {
					winnerId = specie.NoSpecie
				}
			}

			if winnerId == specie.NoSpecie && cellWillLive {
				winnerId = cellId
			}
			e.nextWorld[cellIndex] = winnerId
		}
	}

	return workerPopulation
}

func (e *Ecosystem) index(x, y int) int {
	return y*e.config.Width + x
}
