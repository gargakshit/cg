package main

import (
	"gonum.org/v1/gonum/spatial/r3"
)

func clamp(min, max, f float64) float64 {
	if f < min {
		return min
	}

	if f > max {
		return max
	}

	return f
}

func lerp(f0, f1 float64, fac float64) float64 {
	return f0 + (f1-f0)*clamp(0, 1, fac)
}

func lerpVec(f0, f1 r3.Vec, fac float64) r3.Vec {
	return r3.Vec{
		X: lerp(f0.X, f1.X, fac),
		Y: lerp(f0.Y, f1.Y, fac),
		Z: lerp(f0.Z, f1.Z, fac),
	}
}
