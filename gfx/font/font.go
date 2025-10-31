package font

type Font struct {
	Width             int
	Height            int
	YOffset           int
	HorizontalSpacing int
	VerticalSpacing   int
	Bitmap            map[rune][]byte
	Get               func(rune) []byte
}
