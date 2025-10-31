package main

import (
	"fmt"
	"main/ecosystem"
	"main/gfx"
	"main/renderer"
	"main/rng"
	"main/specie"
	"main/ui/legend"
)

func main() {
	config := newConfig()
	rng.InitRNG(config.RNG)

	compiledSpecies := specie.CompileSpecies(config.Ecosystem.Species)
	// rng.Rand.Shuffle(len(compiledSpecies), func(i, j int) {
	// 	compiledSpecies[i], compiledSpecies[j] = compiledSpecies[j], compiledSpecies[i]
	// })

	ren := renderer.NewRenderer(config.Renderer)
	if err := ren.Start(); err != nil {
		fmt.Println(err)
		return
	}

	eco := ecosystem.NewEcosystem(config.Ecosystem, compiledSpecies)
	leg := legend.NewLegend(config.UI.Legend, compiledSpecies)
	buf := gfx.NewBuffer(config.Renderer.Video.SourceWidth, config.Renderer.Video.SourceHeight)

	for i := range config.Renderer.Video.Length {
		fmt.Printf("\r%d/%d (%.1f%%)", i+1, config.Renderer.Video.Length, float32(i+1)/float32(config.Renderer.Video.Length)*100)

		eco.Render(buf)
		if config.UI.Legend.Enabled {
			leg.Draw(buf)
		}

		ren.Write(buf)

		// start := time.Now()
		eco.Step()
		// fmt.Printf("%v\n", time.Since(start))
	}
}
