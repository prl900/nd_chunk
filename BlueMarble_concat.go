package main

import (
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

func Concat(tiles ...*Tile) *Tile {
	width := 400 * len(tiles)
	height := 400
	subset := []byte{}

	for i := 0; i < 400; i++ {
		for _, t := range tiles {
			subset = append(subset, t.Data[i*400:(i+1)*400]...)
		}
	}

	return &Tile{Data: subset, Shape: []int{width, height}}
}

func ReadSnappyFile(fName string) ([]byte, error) {
	start := time.Now()

	data, err := ioutil.ReadFile(fName)
	cdata, err := snappy.Decode(nil, data)
	log.Printf("Reading Snappy File from disk: %v\n", time.Since(start))

	return cdata, err
}

func WriteSnappyFile(fName string, data []byte) error {
	start := time.Now()

	err := ioutil.WriteFile(fName, snappy.Encode(nil, data), 0644)
	log.Printf("Writting Snappy File to disk: %v\n", time.Since(start))

	return err
}

func main() {
	t1, err := ReadSnappyFile("../Downloads/snappy_tiles/tile_26_07.snp")
	if err != nil {
		panic(err)
	}
	t2, err := ReadSnappyFile("../Downloads/snappy_tiles/tile_27_07.snp")
	if err != nil {
		panic(err)
	}
	t3, err := ReadSnappyFile("../Downloads/snappy_tiles/tile_28_07.snp")
	if err != nil {
		panic(err)
	}
	t4, err := ReadSnappyFile("../Downloads/snappy_tiles/tile_29_07.snp")
	if err != nil {
		panic(err)
	}

	agg := Concat(&Tile{Data: t1, Shape: []int{400, 400}},
		&Tile{Data: t2, Shape: []int{400, 400}},
		&Tile{Data: t3, Shape: []int{400, 400}},
		&Tile{Data: t4, Shape: []int{400, 400}})

	ch := image.Gray{Pix: agg.Data, Stride: agg.Shape[0], Rect: image.Rect(0, 0, agg.Shape[0], agg.Shape[1])}

	f, err := os.Create("../Downloads/agg_tile.png")
	if err != nil {
		panic(err)
	}

	png.Encode(f, &ch)

}
