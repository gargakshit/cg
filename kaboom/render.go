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
	sphereRadius = 1.5
)

func signedDistance(p r3.Vec) float64 {
	return r3.Norm(p) - sphereRadius
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
	for i := 0; i < 128; i++ {
		d := signedDistance(pos)
		if d < 0 {
			return true, pos
		}

		pos = r3.Add(pos, r3.Scale(math.Max(0.1*d, 0.01), direction))
	}

	return false, pos
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
		displacement := (math.Sin(16*hit.X)*math.Sin(16*hit.Y)*math.Sin(16*hit.Z) + 1) / 2

		img.Set(x, y, vec3ToColor(r3.Scale(
			intensity*displacement,
			r3.Vec{X: 1, Y: 1, Z: 1}),
		))
	} else {
		img.Set(x, y, vec3ToColor(r3.Vec{X: 0.2, Y: 0.7, Z: 0.8}))
	}
}
