package main

import (
	"gonum.org/v1/gonum/spatial/r3"
)

func vec3Mul(v1, v2 r3.Vec) r3.Vec {
	return r3.Vec{
		X: v1.X * v2.X,
		Y: v1.Y * v2.Y,
		Z: v1.Z * v2.Z,
	}
}

func vec3XYZ(vec r3.Vec) float64 {
	return vec.X + vec.Y + vec.Z
}
