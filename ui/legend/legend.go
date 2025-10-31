package legend

import (
	"main/gfx"
	"main/specie"
)

type Legend struct {
	config  Config
	species []*specie.CompiledSpecie
}

func NewLegend(config Config, species []*specie.CompiledSpecie) *Legend {
	return &Legend{
		config:  config,
		species: species,
	}
}

func (l *Legend) Draw(buf *gfx.Buffer) {
	maxLen := 0
	for _, specie := range l.species {
		if l := len(specie.Name); l > maxLen {
			maxLen = l
		}
	}

	fnt := l.config.Font
	colorIndicatorSize := max(fnt.Height+fnt.YOffset, fnt.Width)
	rectWidth := l.config.Padding +
		colorIndicatorSize +
		fnt.HorizontalSpacing +
		maxLen*(fnt.Width+fnt.HorizontalSpacing) - fnt.HorizontalSpacing +
		l.config.Padding

	numLines := len(l.species)
	rectHeight := l.config.Padding +
		numLines*(fnt.Height+fnt.YOffset+fnt.VerticalSpacing) - fnt.VerticalSpacing +
		l.config.Padding

	buf.DrawRect(l.config.X, l.config.Y,
		rectWidth, rectHeight,
		l.config.BackgroundColor,
	)

	contentX := l.config.X + l.config.Padding
	contentY := l.config.Y + l.config.Padding

	for i, specie := range l.species {
		colorIndicatorX := contentX
		colorIndicatorY := contentY +
			i*(colorIndicatorSize+fnt.VerticalSpacing)

		textX := colorIndicatorX + colorIndicatorSize + fnt.HorizontalSpacing
		textY := colorIndicatorY

		buf.DrawRect(colorIndicatorX, colorIndicatorY, colorIndicatorSize, colorIndicatorSize, specie.Color)
		buf.DrawString(textX, textY, l.config.FontColor, fnt, specie.Name)
	}
}
