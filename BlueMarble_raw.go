package main

import (
	"image"
	"io/ioutil"
	"log"
)

func main() {
	content, err := ioutil.ReadFile("../Downloads/dat")
	if err != nil {
		log.Fatal(err)
	}
	img := image.RGBA{Pix: content, Stride: 86400, Rect: image.Rect(0, 0, 21600, 10800)}

	log.Println(img.At(10, 10))

	bounds := img.Bounds()

	log.Println(bounds)

	canvas := image.NewAlpha(bounds)

	// is this image opaque
	op := canvas.Opaque()

	log.Println(op)
}
