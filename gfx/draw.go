package gfx

import (
	"image/color"
	"main/gfx/font"
)

func (b *Buffer) pixelIndex(x, y int) int {
	return (y*b.Width + x) * 3
}

func (b *Buffer) DrawRect(x, y, w, h int, col color.Color) {
	r1, g1, b1, a1 := col.RGBA()
	r1f := float32(r1>>8) / 255
	g1f := float32(g1>>8) / 255
	b1f := float32(b1>>8) / 255
	a1f := float32(a1>>8) / 255

	for yy := y; yy < y+h; yy++ {
		for xx := x; xx < x+w; xx++ {
			if xx < 0 || xx >= int(b.Width) || yy < 0 || yy >= int(b.Height) {
				continue
			}
			idx := b.pixelIndex(xx, yy)
			r0f := float32(b.Data[idx+0]) / 255
			g0f := float32(b.Data[idx+1]) / 255
			b0f := float32(b.Data[idx+2]) / 255

			b.Data[idx+0] = uint8((r1f*a1f + r0f*(1-a1f)) * 255)
			b.Data[idx+1] = uint8((g1f*a1f + g0f*(1-a1f)) * 255)
			b.Data[idx+2] = uint8((b1f*a1f + b0f*(1-a1f)) * 255)
		}
	}

}

func (b *Buffer) DrawChar(x, y int, col color.Color, font *font.Font, character rune) {
	glyph := font.Get(character)

	r8, g8, b8, _ := col.RGBA()
	r8b := uint8(r8 >> 8)
	g8b := uint8(g8 >> 8)
	b8b := uint8(b8 >> 8)

	for row := range glyph {
		for colBit := range 3 {
			if glyph[row]&(1<<(2-colBit)) != 0 {
				px := x + colBit
				py := y + row + font.YOffset
				if px < 0 || px >= b.Width || py < 0 || py >= b.Height {
					continue
				}
				idx := b.pixelIndex(px, py)
				b.Data[idx+0] = r8b
				b.Data[idx+1] = g8b
				b.Data[idx+2] = b8b
			}
		}
	}
}

func (b *Buffer) DrawString(x, y int, col color.Color, font *font.Font, str string) {
	for i, ch := range str {
		b.DrawChar(x+i*(font.Width+1), y, col, font, ch)
	}
}
