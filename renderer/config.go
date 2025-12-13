package renderer

type Config struct {
	Video struct {
		FPS int

		SourceWidth  int
		SourceHeight int

		OutputWidth  int
		OutputHeight int

		OutputFile string
	}
}
