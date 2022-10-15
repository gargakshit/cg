package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"time"
)

const (
	width  = 640
	height = 480
	fov    = math.Pi / 3
)

func main() {
	fmt.Println("Kaboom!")

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	timeStart := time.Now()
	render(img)
	fmt.Println("Rendering took", time.Since(timeStart))

	f, err := os.Create("./out/out.png")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
