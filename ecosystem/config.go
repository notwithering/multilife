package ecosystem

import (
	"image/color"
	"main/specie"
)

type Config struct {
	Species []specie.SpecieConfig
	Width   int
	Height  int

	Region struct {
		Density int
		Padding int
	}

	Render struct {
		BackgroundColor color.Color
	}
}
