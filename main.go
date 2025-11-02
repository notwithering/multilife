package main

import (
	"fmt"

	"github.com/notwithering/multilife/ecosystem"
	"github.com/notwithering/multilife/gfx"
	"github.com/notwithering/multilife/renderer"
	"github.com/notwithering/multilife/rng"
	"github.com/notwithering/multilife/specie"
	"github.com/notwithering/multilife/stats"
	"github.com/notwithering/multilife/ui"
)

func main() {
	// cpuFile, _ := os.Create("cpu.prof")
	// pprof.StartCPUProfile(cpuFile)
	// defer func() {
	// 	pprof.StopCPUProfile()
	// 	cpuFile.Close()
	// }()

	config := newConfig()
	rng.InitRNG(config.RNG)

	compiledSpecies := specie.CompileSpecies(config.Ecosystem.Species)
	statsPrinter := stats.NewStatsPrinter(config.Stats, compiledSpecies)
	buffer := gfx.NewBuffer(config.Renderer.Video.SourceWidth, config.Renderer.Video.SourceHeight)
	eco := ecosystem.NewEcosystem(config.Ecosystem, compiledSpecies)
	ui := ui.NewUI(config.UI, compiledSpecies)
	ren := renderer.NewRenderer(config.Renderer)
	if err := ren.Start(); err != nil {
		fmt.Println(err)
		return
	}

	statsPrinter.StartedLoop()
	for range config.Renderer.Video.Length {
		statsPrinter.StartedFrame()

		statsPrinter.StartedRender()
		eco.Render(buffer)
		statsPrinter.EndedRender()

		statsPrinter.StartedUI()
		ui.Draw(buffer)
		statsPrinter.EndedUI()

		statsPrinter.StartedStep()
		ren.Write(buffer)
		statsPrinter.EndedStep()

		statsPrinter.StartedStep()
		eco.Step(statsPrinter.ShouldEcosystem())
		statsPrinter.EndedStep()

		statsPrinter.EndedFrame()

		statsPrinter.Print(eco.Stats)
	}
	statsPrinter.EndedLoop()

	statsPrinter.PrintClosure()
}
