package main

import (
	"image/color"
	"unicode"
)

const (
	fontWidth    = 3
	fontHeight   = 4
	fontYOffset  = -1
	fontHSpacing = 1
	fontVSpacing = 1
)

// https://github.com/Michaelangel007/nanofont3x4
var font = map[rune][4]byte{
	'A':  {0b000, 0b010, 0b111, 0b101},
	'B':  {0b000, 0b110, 0b111, 0b111},
	'C':  {0b000, 0b011, 0b100, 0b011},
	'D':  {0b000, 0b110, 0b101, 0b110},
	'E':  {0b000, 0b111, 0b110, 0b111},
	'F':  {0b000, 0b111, 0b110, 0b100},
	'G':  {0b000, 0b110, 0b101, 0b111},
	'H':  {0b000, 0b101, 0b111, 0b101},
	'I':  {0b000, 0b010, 0b010, 0b010},
	'J':  {0b000, 0b001, 0b001, 0b110},
	'K':  {0b100, 0b101, 0b110, 0b101},
	'L':  {0b000, 0b100, 0b100, 0b111},
	'M':  {0b000, 0b111, 0b111, 0b101},
	'N':  {0b000, 0b111, 0b101, 0b101},
	'O':  {0b000, 0b111, 0b101, 0b111},
	'P':  {0b000, 0b111, 0b111, 0b100},
	'Q':  {0b000, 0b111, 0b101, 0b110},
	'R':  {0b000, 0b110, 0b111, 0b101},
	'S':  {0b000, 0b011, 0b010, 0b110},
	'T':  {0b000, 0b111, 0b010, 0b010},
	'U':  {0b000, 0b101, 0b101, 0b111},
	'V':  {0b000, 0b101, 0b101, 0b010},
	'W':  {0b000, 0b101, 0b111, 0b111},
	'X':  {0b000, 0b101, 0b010, 0b101},
	'Y':  {0b000, 0b101, 0b010, 0b010},
	'Z':  {0b000, 0b111, 0b010, 0b111},
	'\'': {0b000, 0b010, 0b010, 0b000},
	'/':  {0b000, 0b001, 0b010, 0b100},
	'0':  {0b000, 0b010, 0b101, 0b010},
	'1':  {0b000, 0b100, 0b010, 0b010},
	'2':  {0b000, 0b110, 0b010, 0b011},
	'3':  {0b000, 0b111, 0b011, 0b111},
	'4':  {0b000, 0b101, 0b111, 0b001},
	'5':  {0b000, 0b011, 0b010, 0b110},
	'6':  {0b000, 0b100, 0b111, 0b111},
	'7':  {0b000, 0b111, 0b001, 0b001},
	'8':  {0b000, 0b011, 0b111, 0b110},
	'9':  {0b000, 0b111, 0b111, 0b001},
}

func drawChar(buf []byte, ch rune, x, y int, col color.Color) {
	ch = unicode.ToUpper(ch)
	glyph, ok := font[ch]
	if !ok {
		return
	}
	r, g, b, _ := col.RGBA()
	for row := range len(glyph) {
		for colBit := range 3 {
			if glyph[row]&(1<<(2-colBit)) != 0 {
				px := x + colBit
				py := y + row + fontYOffset
				idx := (py*ecosystemSizeX + px) * 3
				buf[idx+0] = uint8(r >> 8)
				buf[idx+1] = uint8(g >> 8)
				buf[idx+2] = uint8(b >> 8)
			}
		}
	}
}

func drawString(buf []byte, startX, startY int, col color.Color, s string) {
	for i, ch := range s {
		drawChar(buf, ch, startX+i*4, startY, col) // 3px glyph + 1px spacing
	}
}

func drawRect(buf []byte, x0, y0, w, h int, col color.Color) {
	r1, g1, b1, a1 := col.RGBA()
	r1f := float32(r1>>8) / 255
	g1f := float32(g1>>8) / 255
	b1f := float32(b1>>8) / 255
	a1f := float32(a1>>8) / 255

	for y := y0; y < y0+h; y++ {
		for x := x0; x < x0+w; x++ {
			idx := (y*ecosystemSizeX + x) * 3

			r0f := float32(buf[idx+0]) / 255
			g0f := float32(buf[idx+1]) / 255
			b0f := float32(buf[idx+2]) / 255

			buf[idx+0] = uint8((r1f*a1f + r0f*(1-a1f)) * 255)
			buf[idx+1] = uint8((g1f*a1f + g0f*(1-a1f)) * 255)
			buf[idx+2] = uint8((b1f*a1f + b0f*(1-a1f)) * 255)
		}
	}
}
