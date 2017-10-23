package main

import (
	"github.com/golang/snappy"
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

	cmpr := snappy.Encode(nil, img.Pix)

	err = ioutil.WriteFile("../Downloads/dat.snp", cmpr, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
