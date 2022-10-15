package main

import (
	"image/color"

	"gonum.org/v1/gonum/spatial/r3"
)

func vec3ToColor(vec r3.Vec) color.RGBA {
	v := r3.Scale(255, vec)
	return color.RGBA{
		R: uint8(clamp(0, 255, v.X)),
		G: uint8(clamp(0, 255, v.Y)),
		B: uint8(clamp(0, 255, v.Z)),
		A: 255,
	}
}
