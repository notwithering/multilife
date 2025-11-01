package specie

import (
	"image/color"

	"github.com/notwithering/multilife/rule"
)

type SpecieConfig struct {
	Name string
	Rule rule.RuleConfig
}

type SpecieId uint8

type CompiledSpecie struct {
	Id    SpecieId
	Name  string
	Rule  rule.CompiledRule
	Color color.Color
}

var currentId SpecieId

func (c SpecieConfig) Compile(numberOfSpecies int) *CompiledSpecie {
	specie := &CompiledSpecie{
		Id:    currentId,
		Name:  c.Name,
		Rule:  c.Rule.Compile(),
		Color: makeColor(int(currentId), numberOfSpecies),
	}
	currentId++

	return specie
}

func CompileSpecies(specieConfigs []SpecieConfig) []*CompiledSpecie {
	var compiledSpecies []*CompiledSpecie

	for _, specie := range specieConfigs {
		compiledSpecies = append(compiledSpecies, specie.Compile(len(specieConfigs)))
	}

	return compiledSpecies
}
