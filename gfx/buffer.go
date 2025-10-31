package gfx

type Buffer struct {
	Data   []byte
	Width  int
	Height int
}

func NewBuffer(width, height int) *Buffer {
	return &Buffer{
		Data:   make([]byte, width*height*3),
		Width:  width,
		Height: height,
	}
}

func (b *Buffer) Clear(color uint32) {
	for i := range b.Data {
		b.Data[i] = uint8(color >> (8 * (i % 3)))
	}
}
