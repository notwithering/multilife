package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/notwithering/multilife/ecosystem"
	"github.com/notwithering/multilife/gfx"
	"github.com/notwithering/multilife/renderer"
	"github.com/notwithering/multilife/rng"
	"github.com/notwithering/multilife/specie"
	"github.com/notwithering/multilife/ui"

	"github.com/notwithering/sgr"
)

func main() {
	config := newConfig()
	rng.InitRNG(config.RNG)

	compiledSpecies := specie.CompileSpecies(config.Ecosystem.Species)
	buf := gfx.NewBuffer(config.Renderer.Video.SourceWidth, config.Renderer.Video.SourceHeight)
	eco := ecosystem.NewEcosystem(config.Ecosystem, compiledSpecies)
	ui := ui.NewUI(config.UI, compiledSpecies)
	ren := renderer.NewRenderer(config.Renderer)
	if err := ren.Start(); err != nil {
		fmt.Println(err)
		return
	}

	const resetCursor = "\x1b[H"
	const clearBelow = "\x1b[J"

	startTime := time.Now()

	for i := range config.Renderer.Video.Length {
		if config.Stats.Enabled {
			fmt.Print(resetCursor + clearBelow)

			for specieId, population := range eco.Stats.PopulationBySpecie {
				specie := eco.Species[specieId]
				density := float32(population) / float32(eco.Stats.TotalPopulation) * 100
				normalDensity := 100.0 / float32(len(eco.Species))

				text := ""

				if density == 0 {
					text += sgr.FgRed + sgr.Strike
				} else if density < normalDensity/2 {
					text += sgr.FgYellow
				} else if density < normalDensity/1.5 {
					text += sgr.FgHiYellow
				} else {
					text += sgr.FgGreen
				}

				text += specie.Name + ": "                                                       // Conway's Life:
				text += strconv.Itoa(population) + "/" + strconv.Itoa(eco.Stats.TotalPopulation) // 50423/99256
				text += " (" + fmt.Sprintf("%.1f", density) + "%)"                               // (50.8%)
				text += sgr.Reset
				fmt.Println(text)
			}

			progress := "Frame: "                                                                                 // Frame:
			progress += strconv.Itoa(i+1) + "/" + strconv.Itoa(config.Renderer.Video.Length)                      // 757/1800
			progress += " (" + fmt.Sprintf("%.1f", float32(i+1)/float32(config.Renderer.Video.Length)*100) + "%)" // (42.1%)
			fmt.Println(progress)

			fps := 1.0 / eco.Stats.FrameTime.Seconds()
			fpsText := "FPS: "
			fpsText += fmt.Sprintf("%.2f", fps)
			fpsText += " (" + fmt.Sprintf("%v", eco.Stats.FrameTime) + "/frame)" // (3.0002ms)
			fmt.Println(fpsText)

			elapsed := time.Since(startTime)
			minutes := int(elapsed.Minutes())
			seconds := int(elapsed.Seconds()) % 60
			milliseconds := int(elapsed.Milliseconds()) % 1000

			elapsedText := "Elapsed: "
			if minutes > 0 {
				elapsedText += strconv.Itoa(minutes) + "m"
			}
			elapsedText += strconv.Itoa(seconds) + "s"
			elapsedText += strconv.Itoa(milliseconds) + "ms"

			fmt.Println(elapsedText)
		}

		eco.Render(buf)
		ui.Draw(buf)

		ren.Write(buf)

		// start := time.Now()
		// fmt.Printf("%v\n", time.Since(start))

		collectStats := config.Stats.Enabled && i%config.Stats.Interval == 0
		eco.Step(collectStats)
	}
}
