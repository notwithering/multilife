package gfx

import (
	"image/color"

	"github.com/notwithering/multilife/gfx/font"
)

func (b *Buffer) PixelOffset(x, y int) int {
	return (y*b.Width + x) * 3
}

func (b *Buffer) SetPixel(x, y int, col color.Color) {
	if x < 0 || x >= int(b.Width) || y < 0 || y >= int(b.Height) {
		return
	}

	newR, newG, newB, newA := col.RGBA()
	if newA == 0 {
		return
	}

	idx := b.PixelOffset(x, y)

	if newA >= 255 {
		b.Data[idx+0] = uint8(newR)
		b.Data[idx+1] = uint8(newG)
		b.Data[idx+2] = uint8(newB)
	} else {
		oldRf := float32(b.Data[idx+0]) / 255
		oldGf := float32(b.Data[idx+1]) / 255
		oldBf := float32(b.Data[idx+2]) / 255

		newRf := float32(newR>>8) / 255
		newGf := float32(newG>>8) / 255
		newBf := float32(newB>>8) / 255
		newAf := float32(newA>>8) / 255

		b.Data[idx+0] = uint8((newRf*newAf + oldRf*(1-newAf)) * 255)
		b.Data[idx+1] = uint8((newGf*newAf + oldGf*(1-newAf)) * 255)
		b.Data[idx+2] = uint8((newBf*newAf + oldBf*(1-newAf)) * 255)
	}
}

func (b *Buffer) DrawRect(x, y, w, h int, col color.Color) {
	for yy := y; yy < y+h; yy++ {
		for xx := x; xx < x+w; xx++ {
			b.SetPixel(x, y, col)
		}
	}
}

func (b *Buffer) DrawChar(x, y int, col color.Color, font *font.Font, character rune) {
	glyph := font.Get(character)

	r8, g8, b8, _ := col.RGBA()
	r8b := uint8(r8 >> 8)
	g8b := uint8(g8 >> 8)
	b8b := uint8(b8 >> 8)

	for rowNum, row := range glyph {
		for colBit := range 3 {
			if row&(1<<(2-colBit)) != 0 {
				px := x + colBit
				py := y + rowNum + font.YOffset
				if px < 0 || px >= b.Width || py < 0 || py >= b.Height {
					continue
				}
				idx := b.PixelIndex(px, py)
				b.Data[idx+0] = r8b
				b.Data[idx+1] = g8b
				b.Data[idx+2] = b8b
			}
		}
	}
}

func (b *Buffer) DrawString(x, y int, col color.Color, font *font.Font, str string) {
	for i, ch := range str {
		b.DrawChar(x+i*(font.Width+font.HSpacing), y, col, font, ch)
	}
}
