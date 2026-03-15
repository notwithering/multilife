package gfx

import (
	"image/color"

	"github.com/notwithering/multilife/gfx/font"
)

func (b *Buffer) PixelOffset(x, y int) int {
	return (y*b.Width + x) * 3
}

func rgb8(c color.Color) (uint8, uint8, uint8, uint8) {
	switch v := c.(type) {
	case color.RGBA:
		return v.R, v.G, v.B, v.A
	case color.NRGBA:
		return v.R, v.G, v.B, v.A
	case color.Gray:
		return v.Y, v.Y, v.Y, 255
	case color.Alpha:
		return 0, 0, 0, v.A
	default:
		r, g, b, a := c.RGBA()
		return uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)
	}
}

func (b *Buffer) SetPixel(x, y int, col color.Color) {
	if x < 0 || x >= int(b.Width) || y < 0 || y >= int(b.Height) {
		return
	}

	newR, newG, newB, newA := rgb8(col)
	if newA == 0 {
		return
	}

	idx := b.PixelOffset(x, y)

	if newA == 255 {
		b.Data[idx+0] = newR
		b.Data[idx+1] = newG
		b.Data[idx+2] = newB
	} else {
		oldRf := float32(b.Data[idx+0]) / 255
		oldGf := float32(b.Data[idx+1]) / 255
		oldBf := float32(b.Data[idx+2]) / 255

		newRf := float32(newR) / 255
		newGf := float32(newG) / 255
		newBf := float32(newB) / 255
		newAf := float32(newA) / 255

		b.Data[idx+0] = uint8((newRf*newAf + oldRf*(1-newAf)) * 255)
		b.Data[idx+1] = uint8((newGf*newAf + oldGf*(1-newAf)) * 255)
		b.Data[idx+2] = uint8((newBf*newAf + oldBf*(1-newAf)) * 255)
	}
}

func (b *Buffer) DrawRect(x, y, w, h int, col color.Color) {
	for yy := y; yy < y+h; yy++ {
		for xx := x; xx < x+w; xx++ {
			b.SetPixel(xx, yy, col)
		}
	}
}

func (b *Buffer) DrawChar(x, y int, col color.Color, font *font.Font, character rune) {
	glyph := font.Get(character)

	for rowNum, row := range glyph {
		for colBit := range font.Width {
			if row&(1<<(font.Width-1-colBit)) != 0 {
				px := x + colBit
				py := y + rowNum + font.YOffset
				b.SetPixel(px, py, col)
			}
		}
	}
}

func (b *Buffer) DrawString(x, y int, col color.Color, font *font.Font, str string) {
	for i, ch := range str {
		b.DrawChar(x+i*(font.Width+font.HSpacing), y, col, font, ch)
	}
}
