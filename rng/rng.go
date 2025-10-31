package rng

import (
	"math/rand"
)

var (
	Source rand.Source
	Rand   *rand.Rand
)

func InitRNG(config Config) {
	Source = rand.NewSource(config.Seed)
	Rand = rand.New(Source)
}
