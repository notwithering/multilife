package legend

import "main/gfx/font"

type Config struct {
	Enabled bool
	X       int
	Y       int
	Padding int
	Font    *font.Font
}
