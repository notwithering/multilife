package main

import "image/color"

const (
	legendEnabled bool = true
	legendX       int  = 1
	legendY       int  = 1
	legendPadding int  = 1
)

var (
	legendBackgroundColor = color.RGBA{0, 0, 0, 255 / 2}
	legendFontColor       = color.RGBA{170, 170, 170, 255}
)

type legend struct {
	species []*compiledSpecie
}

func newLegend(species []*compiledSpecie) *legend {
	return &legend{
		species: species,
	}
}

func (l *legend) draw(buf []byte) {
	numLines := len(species)

	maxLen := 0
	for _, specie := range l.species {
		if l := len(specie.name); l > maxLen {
			maxLen = l
		}
	}

	colorIndicatorXOffset := legendX + legendPadding
	colorIndicatorYOffset := legendY + legendPadding
	colorIndicatorSize := max(fontHeight+fontYOffset, fontWidth)

	textXOffset := colorIndicatorXOffset + colorIndicatorSize + 1
	textYOffset := colorIndicatorYOffset

	rectWidth := maxLen*(fontWidth+fontHSpacing) + legendPadding*2 + colorIndicatorSize + 1
	rectHeight := numLines*(fontHeight+fontVSpacing) - fontVSpacing + legendPadding*2

	drawRect(buf, legendX, legendY, rectWidth, rectHeight, legendBackgroundColor)

	for i, specie := range l.species {
		line := i * (fontHeight + fontVSpacing)
		drawRect(buf, colorIndicatorXOffset, colorIndicatorYOffset+line, colorIndicatorSize, colorIndicatorSize, specie.color)
		drawString(buf, textXOffset, textYOffset+line, legendFontColor, specie.name)
	}
}
