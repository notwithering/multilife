package main

import (
	"image/color"

	"github.com/notwithering/multilife/ecosystem"
	"github.com/notwithering/multilife/gfx/font"
	"github.com/notwithering/multilife/renderer"
	"github.com/notwithering/multilife/rng"
	"github.com/notwithering/multilife/specie"
	"github.com/notwithering/multilife/stats"
	"github.com/notwithering/multilife/ui"
)

func newConfig() config {
	config := config{}

	// ecosystem

	config.Ecosystem.Species = []specie.SpecieConfig{
		specie.SpecieSlowBlob,
		specie.SpecieWalledCities,
		specie.SpecieVote4_5,
		specie.SpecieDiamoeba,
		specie.SpecieVote,
		specie.SpecieAmoeba,
		specie.SpecieConwaysLife,
		specie.SpecieBacteria,
	}

	sizeDivisor := 4
	config.Ecosystem.Width = 1920 / sizeDivisor
	config.Ecosystem.Height = 1080 / sizeDivisor

	config.Ecosystem.Render.BackgroundColor = color.Black

	config.Ecosystem.Region.Density = 50 //%
	config.Ecosystem.Region.Padding = 10 //px

	// renderer

	config.Renderer.Video.FPS = 60
	config.Renderer.Video.SourceWidth = config.Ecosystem.Width
	config.Renderer.Video.SourceHeight = config.Ecosystem.Height
	config.Renderer.Video.OutputWidth = 1920
	config.Renderer.Video.OutputHeight = 1080
	config.Renderer.Video.OutputFile = "output.mp4"

	// main

	config.Main.Infinite = true
	videoLengthInSeconds := 30 //seconds
	config.Main.VideoLength = config.Renderer.Video.FPS * videoLengthInSeconds

	// legend

	config.UI.Legend.Enabled = true
	config.UI.Legend.X = 1
	config.UI.Legend.Y = 1
	config.UI.Legend.Padding = 1
	config.UI.Legend.Font = font.Nanofont3x4
	config.UI.Legend.BackgroundColor = color.RGBA{0, 0, 0, 255 / 2}
	config.UI.Legend.FontColor = color.RGBA{170, 170, 170, 255}

	// rng

	config.RNG.Seed = 0

	// stats

	config.Stats.Infinite = config.Main.Infinite

	config.Stats.Basic.Enabled = true
	config.Stats.Basic.Interval = 30
	config.Stats.Basic.FPS = config.Renderer.Video.FPS

	config.Stats.Ecosystem.Enabled = true
	config.Stats.Ecosystem.Interval = 50

	config.Stats.Basic.TotalFrames = config.Main.VideoLength

	return config
}

type config struct {
	Ecosystem ecosystem.Config
	Renderer  renderer.Config
	Main      MainConfig
	UI        ui.Config
	RNG       rng.Config
	Stats     stats.Config
}

type MainConfig struct {
	VideoLength int
	Infinite    bool
}
