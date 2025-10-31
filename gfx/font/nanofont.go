package font

import "unicode"

var Nanofont3x4 = &Font{
	Width:             3,
	Height:            4,
	YOffset:           -1,
	HorizontalSpacing: 1,
	VerticalSpacing:   1,
	Bitmap:            nanofont3x4Bitmap,
	Get:               nanofont3x4Get,
}

// https://github.com/Michaelangel007/nanofont3x4
var nanofont3x4Bitmap = map[rune][]byte{
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
	' ':  {0b000, 0b000, 0b000, 0b000},
}

func nanofont3x4Get(r rune) []byte {
	r = unicode.ToUpper(r)
	glyph, ok := nanofont3x4Bitmap[r]
	if !ok {
		return nanofont3x4Bitmap[' ']
	}
	return glyph
}
