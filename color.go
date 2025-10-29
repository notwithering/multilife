package main

import (
	"image/color"
)

const (
	randomColors bool = true
)

func randomColor() color.RGBA {
	h := rng.Float64()
	s := 0.5 + 0.5*rng.Float64()
	v := 0.5 + 0.5*rng.Float64()

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
