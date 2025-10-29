package main

import (
	"image/color"
	"math"
)

const (
	ecosystemSizeDiv int = 4
	ecosystemSizeX   int = 1920 / ecosystemSizeDiv
	ecosystemSizeY   int = 1080 / ecosystemSizeDiv

	ecosystemRegionDensity int = 50 //%
	ecosystemRegionPadding int = 20 //px
)

var (
	ecosystemBackgroundColor color.Color = color.Black
)

type ecosystem struct {
	world   [ecosystemSizeX][ecosystemSizeY]*compiledSpecie
	species []*compiledSpecie
}

func (e ecosystem) render(buf []byte) {
	i := 0
	for y := range ecosystemSizeY {
		for x := range ecosystemSizeX {
			c := ecosystemBackgroundColor
			if e.world[x][y] != nil {
				c = e.world[x][y].color
			}
			r, g, b_, _ := c.RGBA()
			buf[i+0] = uint8(r >> 8)
			buf[i+1] = uint8(g >> 8)
			buf[i+2] = uint8(b_ >> 8)
			i += 3
		}
	}
}

func newEcosystem(species []*compiledSpecie) ecosystem {
	var eco ecosystem
	eco.species = species
	columns, rows := computeGrid(len(species))

	var regionWidth int = ecosystemSizeX / columns
	var regionHeight int = ecosystemSizeY / rows

	for specieIndex, specie := range species {
		column := specieIndex % columns
		row := specieIndex / columns

		xStart := column*regionWidth + ecosystemRegionPadding
		xEnd := (column+1)*regionWidth - ecosystemRegionPadding
		yStart := row*regionHeight + ecosystemRegionPadding
		yEnd := (row+1)*regionHeight - ecosystemRegionPadding

		for x := xStart; x < xEnd; x++ {
			for y := yStart; y < yEnd; y++ {
				if rng.Int()%100 < ecosystemRegionDensity {
					eco.world[x][y] = specie
				}
			}
		}
	}

	return eco
}

func computeGrid(n int) (rows, columns int) {
	if n <= 0 {
		return 0, 0
	}

	for columns = int(math.Round(math.Sqrt(float64(n)))); columns > 0; columns-- {
		if n%columns == 0 {
			rows = n / columns
			return
		}
	}

	return
}

var neighborOffsets = [8][2]int{
	{-1, -1}, {0, -1}, {1, -1},
	{-1, 0}, {1, 0},
	{-1, 1}, {0, 1}, {1, 1},
}

func (e *ecosystem) step() {
	var next ecosystem
	specieNeighbors := make([]int, len(species))

	for x := range ecosystemSizeX {
		for y := range ecosystemSizeY {
			for id := range specieNeighbors {
				specieNeighbors[id] = 0
			}

			var totalNeighbors int
			for _, offset := range neighborOffsets {
				neighborX := (x + offset[0] + ecosystemSizeX) % ecosystemSizeX
				neighborY := (y + offset[1] + ecosystemSizeY) % ecosystemSizeY
				if specie := e.world[neighborX][neighborY]; specie != nil {
					specieNeighbors[specie.id]++
					totalNeighbors++
				}
			}

			if totalNeighbors == 0 {
				continue
			}

			type candidate struct {
				specie *compiledSpecie
				weight int
			}
			var candidates []candidate
			totalWeight := 0

			addCandidate := func(specie *compiledSpecie) {
				weight := specieNeighbors[specie.id] + 1*100/totalNeighbors + 1
				totalWeight += weight
				candidates = append(candidates, candidate{specie, weight})
			}

			me := e.world[x][y]
			isAlive := me != nil
			shouldSurvive := isAlive && me.rule.surviveSet[totalNeighbors]
			if shouldSurvive {
				addCandidate(me)
			}

			for _, specie := range e.species {
				neighborsOfSpecie := specieNeighbors[specie.id]
				if neighborsOfSpecie == 0 {
					continue
				}

				shouldBirth := specie.rule.birthSet[neighborsOfSpecie]
				differentSpecie := specie != me
				if shouldBirth && differentSpecie {
					addCandidate(specie)
				}
			}

			var winnerSpecie *compiledSpecie

			if totalWeight > 0 {
				sum := 0
				for _, candidate := range candidates {
					sum += candidate.weight
					if rng.Intn(totalWeight) < sum {
						winnerSpecie = candidate.specie
						break
					}
				}
			}

			next.world[x][y] = winnerSpecie
		}
	}

	e.world = next.world
}
