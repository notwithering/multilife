package ui

import (
	"github.com/notwithering/multilife/gfx"
	"github.com/notwithering/multilife/specie"
	"github.com/notwithering/multilife/ui/legend"
)

type UI struct {
	legend *legend.Legend
}

func NewUI(config Config, species []*specie.CompiledSpecie) UI {
	return UI{
		legend: legend.NewLegend(config.Legend, species),
	}
}

func (ui *UI) Draw(buf *gfx.Buffer) {
	ui.legend.Draw(buf)
}
