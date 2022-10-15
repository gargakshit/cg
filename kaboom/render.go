package main

import (
	"fmt"
	"image"
	"math"
	"runtime"
	"sync"

	"gonum.org/v1/gonum/spatial/r3"
)

const (
	sphereRadius   = 1.5
	noiseAmplitude = 1.0
)

func hash(n float64) float64 {
	x := math.Sin(n) * 43758.5453
	return x - math.Floor(x)
}

func noise(x r3.Vec) float64 {
	p := r3.Vec{X: math.Floor(x.X), Y: math.Floor(x.Y), Z: math.Floor(x.Z)}
	f := r3.Sub(x, p)
	f = vec3Mul(f, vec3Mul(f, r3.Sub(r3.Vec{X: 3, Y: 3, Z: 3}, r3.Scale(2, f))))
	n := vec3XYZ(vec3Mul(p, r3.Vec{X: 1, Y: 57, Z: 113}))

	return lerp(
		lerp(
			lerp(hash(n), hash(n+1), f.X),
			lerp(hash(n+57), hash(n+58), f.X),
			f.Y,
		),
		lerp(
			lerp(hash(n+113), hash(n+114), f.X),
			lerp(hash(n+170), hash(n+171), f.X),
			f.Y,
		),
		f.Z,
	)
}

func rotate(v r3.Vec) r3.Vec {
	return r3.Vec{
		X: vec3XYZ(vec3Mul(r3.Vec{Y: 0.8, Z: 0.6}, v)),
		Y: vec3XYZ(vec3Mul(r3.Vec{X: -0.80, Y: 0.36, Z: -0.48}, v)),
		Z: vec3XYZ(vec3Mul(r3.Vec{X: -0.6, Y: -0.48, Z: 0.64}, v)),
	}
}

func fractalBrownianMotion(v r3.Vec) float64 {
	p := rotate(v)

	f := 0.0
	f += 0.5 * noise(p)
	p = r3.Scale(2.32, p)
	f += 0.25 * noise(p)
	p = r3.Scale(3.03, p)
	f += 0.125 * noise(p)
	p = r3.Scale(2.61, p)
	f += 0.0625 * noise(p)

	return f / 0.9375
}

func signedDistance(p r3.Vec) float64 {
	displacement := -fractalBrownianMotion(r3.Scale(3.4, p)) * noiseAmplitude
	return r3.Norm(p) - (sphereRadius + displacement)
}

const eps = 0.1

func distanceFieldNormal(pos r3.Vec) r3.Vec {
	d := signedDistance(pos)

	nx := signedDistance(r3.Add(pos, r3.Vec{X: eps})) - d
	ny := signedDistance(r3.Add(pos, r3.Vec{Y: eps})) - d
	nz := signedDistance(r3.Add(pos, r3.Vec{Z: eps})) - d

	return r3.Unit(r3.Vec{X: nx, Y: ny, Z: nz})
}

func sphereTrace(origin, direction r3.Vec) (bool, r3.Vec) {
	pos := origin

	if vec3XYZ(vec3Mul(origin, origin))-math.Pow(vec3XYZ(vec3Mul(origin, direction)), 2) >
		sphereRadius*sphereRadius {
		return false, pos
	}

	for i := 0; i < 128; i++ {
		d := signedDistance(pos)
		if d < 0 {
			return true, pos
		}

		pos = r3.Add(pos, r3.Scale(math.Max(0.1*d, 0.01), direction))
	}

	return false, pos
}

var (
	yellow   = r3.Vec{X: 1.7, Y: 1.3, Z: 1}
	orange   = r3.Vec{X: 1, Y: 0.6}
	red      = r3.Vec{X: 1}
	darkGray = r3.Vec{X: 0.2, Y: 0.2, Z: 0.2}
	gray     = r3.Vec{X: 0.4, Y: 0.4, Z: 0.4}
)

func paletteFire(d float64) r3.Vec {
	x := clamp(0, 1, d)

	if x < 0.25 {
		return lerpVec(gray, darkGray, 4*x)
	}

	if x < 0.5 {
		return lerpVec(darkGray, red, x*4-1)
	}

	if x < 0.75 {
		return lerpVec(red, orange, x*4-2)
	}

	return lerpVec(orange, yellow, x*4-3)
}

func render(img *image.RGBA) {
	numPartitions := runtime.GOMAXPROCS(0)
	partitionSize := height / numPartitions

	fmt.Println("Running in parallel with", numPartitions, "partitions")

	var wg sync.WaitGroup
	wg.Add(numPartitions)

	for i := 0; i < numPartitions; i++ {
		go renderPartition(&wg, i, partitionSize, img)
	}

	wg.Wait()
}

func renderPartition(wg *sync.WaitGroup, i, size int, img *image.RGBA) {
	base := i * size
	for y := base; y < base+size; y++ {
		for x := 0; x < width; x++ {
			renderPixel(x, y, img)
		}
	}

	wg.Done()
}

func renderPixel(x, y int, img *image.RGBA) {
	dir := r3.Vec{
		X: (float64(x) + 0.5) - (float64(width) / 2),
		Y: -(float64(y) + 0.5) + (float64(height) / 2),
		Z: -float64(height) / (2 * math.Tan(fov/2)),
	}

	didHit, hit := sphereTrace(r3.Vec{Z: 3}, r3.Unit(dir))
	if didHit {
		lightDir := r3.Unit(r3.Sub(r3.Vec{X: 10, Y: 10, Z: 10}, hit))
		intensity := math.Max(
			0.4,
			vec3XYZ(vec3Mul(lightDir, distanceFieldNormal(hit))),
		)

		noiseLevel := (sphereRadius - r3.Norm(hit)) / noiseAmplitude
		img.Set(x, y, vec3ToColor(r3.Scale(intensity, paletteFire((-0.2+noiseLevel)*2))))
	} else {
		img.Set(x, y, vec3ToColor(r3.Vec{X: 0.2, Y: 0.7, Z: 0.8}))
	}
}
