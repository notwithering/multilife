package main

import (
	"fmt"
)

func main() {
	var compiledSpecies []*compiledSpecie
	for i, specie := range species {
		compiled := specie.compile()
		compiled.id = i
		if randomColors {
			compiled.color = randomColor()
		}
		compiledSpecies = append(compiledSpecies, compiled)
	}

	rng.Shuffle(len(compiledSpecies), func(i, j int) {
		compiledSpecies[i], compiledSpecies[j] = compiledSpecies[j], compiledSpecies[i]
	})

	renderer, err := newRenderer()
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := renderer.start(); err != nil {
		fmt.Println(err)
		return
	}

	ecosystem := newEcosystem(compiledSpecies)
	legend := newLegend(compiledSpecies)

	buf := make([]byte, ecosystemSizeX*ecosystemSizeY*3)

	for i := range videoFrames {
		fmt.Printf("\r%d/%d (%.1f%%)", i+1, videoFrames, float32(i+1)/float32(videoFrames)*100)

		ecosystem.render(buf)
		if legendEnabled {
			legend.draw(buf)
		}

		renderer.pipe.Write(buf)

		// start := time.Now()
		ecosystem.step()
		// fmt.Printf("%v\n", time.Since(start))
	}
}
