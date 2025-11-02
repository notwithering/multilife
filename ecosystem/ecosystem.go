package ecosystem

import (
	"math"
	"runtime"
	"sync"
	"time"

	"github.com/notwithering/multilife/gfx"
	"github.com/notwithering/multilife/rng"
	"github.com/notwithering/multilife/specie"
)

const (
	EmptyCell specie.SpecieId = 255
)

type Ecosystem struct {
	config    Config
	Species   []*specie.CompiledSpecie
	world     []specie.SpecieId
	nextWorld []specie.SpecieId
	Stats     Stats
}

type Stats struct {
	TotalPopulation    int
	PopulationBySpecie []int
	FrameTime          time.Duration
}

func (e *Ecosystem) Render(buf *gfx.Buffer) {
	for y := range e.config.Height {
		idx := y * e.config.Width * 3
		for x := range e.config.Width {
			specieId := e.world[e.index(x, y)]

			col := e.config.Render.BackgroundColor
			if specieId != EmptyCell {
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

	for x := range config.Width {
		for y := range config.Height {
			eco.world[eco.index(x, y)] = EmptyCell
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
		world[i] = EmptyCell
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

var neighborhood = [8][2]int{
	{-1, -1}, {0, -1}, {1, -1},
	{-1, 0} /*     */, {1, 0},
	{-1, 1}, {0, 1}, {1, 1},
}

func (e *Ecosystem) Step(collectStats bool) {
	numWorkers := max(runtime.NumCPU(), e.config.Height)
	rowsPerWorker := e.config.Height / numWorkers

	startTime := time.Now()
	if collectStats {
		e.Stats.TotalPopulation = 0
		e.Stats.PopulationBySpecie = make([]int, len(e.Species))
	}

	workerPopulations := make([][]int, numWorkers)
	for worker := range numWorkers {
		workerPopulations[worker] = make([]int, len(e.Species))
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
			specieNeighbors := make([]int, len(e.Species))
			candidates := make([]specie.SpecieId, len(e.Species))

			for y := startY; y < endY; y++ {
				for x := range e.config.Width {
					clear(specieNeighbors)

					var totalNeighbors int
					for _, offset := range neighborhood {
						neighborX := (x + offset[0] + e.config.Width) % e.config.Width
						neighborY := (y + offset[1] + e.config.Height) % e.config.Height
						neighborSpecieId := e.world[e.index(neighborX, neighborY)]
						if neighborSpecieId != EmptyCell {
							specieNeighbors[neighborSpecieId]++
							totalNeighbors++
						}
					}

					cellId := e.world[e.index(x, y)]

					var cell *specie.CompiledSpecie
					cellIsAlive := false
					cellWillLive := false

					if cellId != EmptyCell {
						cell = e.Species[cellId]
						cellIsAlive = true
						cellWillLive = cell.Rule.SurviveSet[totalNeighbors]

						if collectStats {
							workerPopulations[worker][cellId]++
						}
					}

					clear(candidates)

					for specieIdInt, neighborsOfSpecie := range specieNeighbors {
						specieId := specie.SpecieId(specieIdInt)
						specie := e.Species[specieId]
						shouldBirth := specie.Rule.BirthSet[neighborsOfSpecie]
						differentSpecie := specieId != cellId

						canCompete := shouldBirth &&
							((cellIsAlive && differentSpecie) ||
								(!cellIsAlive))

						if canCompete {
							candidates = append(candidates, specieId)
						}
					}

					winnerId := EmptyCell
					maxWeight := -1

					for _, candidateId := range candidates {
						neighbors := specieNeighbors[candidateId]
						if neighbors > maxWeight {
							maxWeight = neighbors
							winnerId = candidateId
						} else if neighbors == maxWeight {
							winnerId = EmptyCell
						}
					}

					if winnerId == EmptyCell && cellWillLive {
						winnerId = cellId
					}
					e.nextWorld[e.index(x, y)] = winnerId
				}
			}
		}()
	}

	wg.Wait()
	e.world, e.nextWorld = e.nextWorld, e.world

	if collectStats {
		e.Stats.FrameTime = time.Since(startTime)

		for _, workerPopulation := range workerPopulations {
			for cellId, count := range workerPopulation {
				e.Stats.TotalPopulation += count
				e.Stats.PopulationBySpecie[cellId] += count
			}
		}
	}
}

func (e *Ecosystem) index(x, y int) int {
	return y*e.config.Width + x
}
