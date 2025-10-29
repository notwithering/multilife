package main

import "math/rand"

const (
	seed int64 = 0
)

var rng *rand.Rand = rand.New(rand.NewSource(seed))
