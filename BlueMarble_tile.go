package main

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/golang/snappy"
)

type Tile struct {
	Data  []byte
	Shape []int
}

func (t *Tile) Subset(x0, y0, x1, y1 int) *Tile {
	subset := []byte{}
	for i := x0 + y0*t.Shape[0]; i < x0+y1*(t.Shape[0]); i += t.Shape[0] {
		subset = append(subset, t.Data[i:i+(x1-x0)]...)
	}

	return &Tile{Data: subset, Shape: []int{x1 - x0, y1 - y0}}
}

type TiledImg map[int]Tile

func ReadRawFile(fName string) ([]byte, error) {
	start := time.Now()

	data, err := ioutil.ReadFile(fName)
	log.Printf("Reading Raw File from disk: %v\n", time.Since(start))

	return data, err
}

func WriteSnappyFile(fName string, data []byte) error {
	start := time.Now()

	err := ioutil.WriteFile(fName, snappy.Encode(nil, data), 0644)
	log.Printf("Writting Snappy File to disk: %v\n", time.Since(start))

	return err
}

func main() {
	data, err := ReadRawFile("../Downloads/dat1")
	if err != nil {
		panic(err)
	}

	tile := &Tile{Data: data, Shape: []int{21600, 10800}}

	log.Println(tile.Shape)

	subset := tile.Subset(15000, 5000, 20000, 10000)
	log.Println(len(subset.Data), subset.Shape)

	ch := image.Gray{Pix: subset.Data, Stride: subset.Shape[0], Rect: image.Rect(0, 0, subset.Shape[0], subset.Shape[1])}

	f, err := os.Create("../Downloads/tile.png")
	if err != nil {
		panic(err)
	}

	png.Encode(f, &ch)

	for j := 0; j < 10800; j += 400 {
		for i := 0; i < 21600; i += 400 {
			subset := tile.Subset(i, j, i+400, j+400)
			fName := fmt.Sprintf("../Downloads/snappy_tiles/tile_%02d_%02d.snp", i/400, j/400)
			WriteSnappyFile(fName, subset.Data)
		}
	}
}
