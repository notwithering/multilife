package legend

import (
	"image/color"
	"main/gfx/font"
)

type Config struct {
	Enabled         bool
	X               int
	Y               int
	Padding         int
	Font            *font.Font
	BackgroundColor color.Color
	FontColor       color.Color
}
