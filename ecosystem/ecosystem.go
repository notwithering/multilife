package ecosystem

import (
	"image/color"
	"main/gfx"
	"main/rng"
	"main/specie"
	"math"
)

var (
	ecosystemBackgroundColor color.Color = color.Black
)

const (
	EmptyCell int = -1
)

type Ecosystem struct {
	config  Config
	World   [][]int
	Species []*specie.CompiledSpecie
}

func (e Ecosystem) Render(buf *gfx.Buffer) {
	i := 0
	for y := range e.config.Height {
		for x := range e.config.Width {
			col := ecosystemBackgroundColor
			specieId := e.World[x][y]
			if specieId != EmptyCell {
				specie := e.Species[specieId]
				col = specie.Color
			}
			r, g, b_, _ := col.RGBA()
			buf.Data[i+0] = uint8(r >> 8)
			buf.Data[i+1] = uint8(g >> 8)
			buf.Data[i+2] = uint8(b_ >> 8)
			i += 3
		}
	}
}

func NewEcosystem(config Config, species []*specie.CompiledSpecie) Ecosystem {
	var eco Ecosystem
	eco.config = config
	eco.World = newEmptyWorld(config.Width, config.Height)
	eco.Species = species
	columns, rows := computeGrid(len(species))

	var regionWidth int = config.Width / columns
	var regionHeight int = config.Height / rows

	for x := range config.Width {
		for y := range config.Height {
			eco.World[x][y] = EmptyCell
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
					eco.World[x][y] = specie.Id
				}
			}
		}
	}

	return eco
}

func newEmptyWorld(width, height int) [][]int {
	var world = make([][]int, width)
	for x := range world {
		world[x] = make([]int, height)
		for y := range world[x] {
			world[x][y] = EmptyCell
		}
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

func (e *Ecosystem) Step() {
	var next Ecosystem
	next.World = newEmptyWorld(e.config.Width, e.config.Height)
	specieNeighbors := make([]int, len(e.Species))

	for x := range e.config.Width {
		for y := range e.config.Height {
			for id := range specieNeighbors {
				specieNeighbors[id] = 0
			}

			var totalNeighbors int
			for _, offset := range neighborhood {
				neighborX := (x + offset[0] + e.config.Width) % e.config.Width
				neighborY := (y + offset[1] + e.config.Height) % e.config.Height
				neighborSpecieId := e.World[neighborX][neighborY]
				if neighborSpecieId != EmptyCell {
					specieNeighbors[neighborSpecieId]++
					totalNeighbors++
				}
			}

			type candidate struct {
				specieId int
				weight   int
			}
			var candidates []candidate

			cellId := e.World[x][y]

			var cell *specie.CompiledSpecie
			cellIsAlive := false
			cellWillLive := false

			if cellId != EmptyCell {
				cell = e.Species[cellId]
				cellIsAlive = true
				cellWillLive = cell.Rule.SurviveSet[totalNeighbors]
			}

			// cellWillDie := cellIsAlive && !cellShouldSurvive
			// cellWillLive := cellIsAlive && cellShouldSurvive

			for _, specie := range e.Species {
				differentSpecie := specie.Id != cellId
				neighborsOfSpecie := specieNeighbors[specie.Id]
				shouldBirth := specie.Rule.BirthSet[neighborsOfSpecie]

				canCompete := shouldBirth &&
					((cellIsAlive && differentSpecie && shouldBirth) ||
						(!cellIsAlive))

				if canCompete {
					weight := neighborsOfSpecie * 100 / (totalNeighbors + 1)
					candidates = append(candidates, candidate{specie.Id, weight})
				}
			}

			winnerId := EmptyCell
			maxWeight := 0

			for _, candidate := range candidates {
				if candidate.weight > maxWeight {
					maxWeight = candidate.weight
					winnerId = candidate.specieId
				} else if candidate.weight == maxWeight {
					winnerId = EmptyCell
				}
			}

			if winnerId == EmptyCell && cellWillLive {
				winnerId = cellId
			}
			next.World[x][y] = winnerId
		}
	}

	e.World = next.World
}
