package specie

import (
	"image/color"
	"main/rng"
	"main/rule"
)

type SpecieConfig struct {
	Name  string
	Rule  rule.RuleConfig
	Color color.Color
}

type SpecieId uint8

type CompiledSpecie struct {
	Id    SpecieId
	Name  string
	Rule  rule.CompiledRule
	Color color.Color
}

var currentId SpecieId

func (c SpecieConfig) Compile() *CompiledSpecie {
	specie := &CompiledSpecie{
		Id:    currentId,
		Name:  c.Name,
		Rule:  c.Rule.Compile(),
		Color: c.Color,
	}
	currentId++

	if randomColors {
		specie.Color = randomColor()
	}

	return specie
}

func CompileSpecies(specieConfigs []SpecieConfig) []*CompiledSpecie {
	var compiledSpecies []*CompiledSpecie

	for _, specie := range specieConfigs {
		compiledSpecies = append(compiledSpecies, specie.Compile())
	}

	return compiledSpecies
}

const (
	randomColors bool = true
)

func randomColor() color.RGBA {
	h := rng.Rand.Float64()
	s := 0.5 + 0.5*rng.Rand.Float64()
	v := 0.5 + 0.5*rng.Rand.Float64()

	r, g, b := hsvToRgb(h, s, v)
	return color.RGBA{r, g, b, 255}
}

func hsvToRgb(h, s, v float64) (uint8, uint8, uint8) {
	i := int(h * 6)
	f := h*6 - float64(i)
	p := v * (1 - s)
	q := v * (1 - f*s)
	t := v * (1 - (1-f)*s)
	var r, g, b float64
	switch i % 6 {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	case 5:
		r, g, b = v, p, q
	}
	return uint8(r * 255), uint8(g * 255), uint8(b * 255)
}
