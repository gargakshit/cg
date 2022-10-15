package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"runtime"
	"time"
)

const (
	width  = 1500
	height = 1500
	fov    = math.Pi / 3
)

func main() {
	if os.Args[len(os.Args)-1] == "render-animation" {
		err := animationMain()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	} else {
		err := imageMain()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	}
}

var frameTime = 4.0

func imageMain() error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	numPartitions := runtime.GOMAXPROCS(0)
	fmt.Println("Partitions:", numPartitions)
	fmt.Println("Parallel: true")
	fmt.Println("Width:", width)
	fmt.Println("Height:", height)

	timeStart := time.Now()
	render(numPartitions, img)
	fmt.Println("Rendering took", time.Since(timeStart))

	f, err := os.Create("./out/out.png")
	if err != nil {
		return err
	}

	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		return err
	}

	return nil
}

func animationMain() error {
	frameTime = 0.0
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	animationTime := 9.2
	fps := 60.0
	timeIncrement := 1 / fps
	numFrames := int(animationTime * fps)

	numPartitions := runtime.GOMAXPROCS(0)
	fmt.Println("Partitions:", numPartitions)
	fmt.Println("Parallel: true")
	fmt.Println("Width:", width)
	fmt.Println("Height:", height)
	fmt.Println("FPS:", fps)
	fmt.Println("Animation time:", animationTime)
	fmt.Println("Frames:", numFrames)

	lastFrameTime := time.Now()
	for i := 0; i < numFrames; i++ {
		lastFrameTime = time.Now()
		fmt.Printf("\rFrame: %d | Last frame time: %v", i+1, time.Since(lastFrameTime))

		render(numPartitions, img)

		f, err := os.Create(fmt.Sprintf("./out/%d.png", i+1))
		if err != nil {
			return err
		}

		err = png.Encode(f, img)
		if err != nil {
			return err
		}

		err = f.Close()
		if err != nil {
			return err
		}

		frameTime += timeIncrement
	}

	fmt.Println("Rendered all frames")

	return nil
}
