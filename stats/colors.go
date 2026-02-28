package stats

import (
	"image/color"
)

func lerp(a, b uint8, t float64) uint8 {
	return uint8(float64(a) + t*(float64(b)-float64(a)))
}

func lerpColor(c0, c1 color.Color, t float64) color.Color {
	a := c0.(color.RGBA)
	b := c1.(color.RGBA)

	return color.RGBA{
		R: lerp(a.R, b.R, t),
		G: lerp(a.G, b.G, t),
		B: lerp(a.B, b.B, t),
		A: lerp(a.A, b.A, t),
	}
}

func gradient(x float64, colors ...color.Color) color.Color {
	n := len(colors)
	if n == 0 {
		return color.RGBA{}
	}
	if n == 1 {
		return colors[0]
	}

	if x <= 0 {
		return colors[0]
	}
	if x >= 1 {
		return colors[n-1]
	}

	segmentSize := 1.0 / float64(n-1)

	index := int(x / segmentSize)
	if index >= n-1 {
		index = n - 2
	}

	localX := (x - float64(index)*segmentSize) / segmentSize

	return lerpColor(colors[index], colors[index+1], localX)
}
