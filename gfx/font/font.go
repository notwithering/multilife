package font

type Font struct {
	Width   int
	Height  int // height from baseline to standard letters like `A` or `l`
	YOffset int

	Ascent  int // height of area above cap line for letters like `backtick` and accents
	Descent int // height of area below baseline for lowest letters like `y` or `q`

	HSpacing int
	VSpacing int
	Bitmap   map[rune][]byte
	Get      func(rune) []byte
}
