package main

import (
	"bytes"
	"compress/gzip"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/golang/snappy"
	"github.com/pierrec/lz4"
)

func WriteRawFile(fName string, data []byte) error {
	start := time.Now()

	err := ioutil.WriteFile(fName, data, 0644)
	log.Printf("Writting Raw File to disk: %v\n", time.Since(start))

	return err
}

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

func ReadSnappyFile(fName string) ([]byte, error) {
	start := time.Now()

	data, err := ioutil.ReadFile(fName)
	cdata, err := snappy.Decode(nil, data)
	log.Printf("Reading Snappy File from disk: %v\n", time.Since(start))

	return cdata, err
}

func WriteLZ4File(fName string, data []byte) error {
	start := time.Now()

	comp := make([]byte, len(data))

	//compress
	l, err := lz4.CompressBlock(data, comp, 0)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fName, comp[:l], 0644)
	log.Printf("Writting LZ4 File to disk: %v\n", time.Since(start))

	return err
}

func ReadLZ4File(fName string) ([]byte, error) {
	start := time.Now()

	data, err := ioutil.ReadFile(fName)

	//decompress
	decomp := make([]byte, len(data)*3)
	l, err := lz4.UncompressBlock(data, decomp, 0)
	if err != nil {
		return []byte{}, err
	}
	log.Printf("Reading LZ4 File from disk: %v\n", time.Since(start))

	return decomp[:l], err
}

func WriteZipFile(fName string, data []byte) error {
	start := time.Now()

	outFile, err := os.Create(fName)
	if err != nil {
		return err
	}

	gzipWriter := gzip.NewWriter(outFile)

	if _, err = gzipWriter.Write(data); err != nil {
		return err
	}

	if err = gzipWriter.Close(); err != nil {
		return err
	}

	log.Printf("Writting Zip File to disk: %v\n", time.Since(start))

	return nil
}

func ReadZipFile(fName string) ([]byte, error) {
	start := time.Now()

	data, err := ioutil.ReadFile(fName)
	if err != nil {
		return []byte{}, err
	}

	buf := bytes.NewBuffer(data)

	gzipReader, err := gzip.NewReader(buf)
	if err != nil {
		return []byte{}, err
	}

	var resB bytes.Buffer
	_, err = resB.ReadFrom(gzipReader)
	if err != nil {
		return []byte{}, err
	}

	log.Printf("Reading Zip File from disk: %v\n", time.Since(start))

	return resB.Bytes(), nil
}

func main() {
	src, err := os.Open("../Downloads/world.topo.bathy.200401.3x21600x10800.png")
	if err != nil {
		panic(err)
	}
	defer src.Close()

	img, err := png.Decode(src)
	if err != nil {
		panic(err)
	}

	rgba := img.(*image.RGBA)

	err = WriteRawFile("../Downloads/dat", rgba.Pix)
	if err != nil {
		panic(err)
	}

	err = WriteSnappyFile("../Downloads/dat.snp", rgba.Pix)
	if err != nil {
		panic(err)
	}

	err = WriteLZ4File("../Downloads/dat.lz4", rgba.Pix)
	if err != nil {
		panic(err)
	}

	err = WriteZipFile("../Downloads/dat.gz", rgba.Pix)
	if err != nil {
		panic(err)
	}

	_, err = ReadRawFile("../Downloads/dat")
	if err != nil {
		panic(err)
	}

	_, err = ReadSnappyFile("../Downloads/dat.snp")
	if err != nil {
		panic(err)
	}

	_, err = ReadLZ4File("../Downloads/dat.lz4")
	if err != nil {
		panic(err)
	}

	_, err = ReadZipFile("../Downloads/dat.gz")
	if err != nil {
		panic(err)
	}

	red := make([]byte, len(rgba.Pix)/4)
	green := make([]byte, len(rgba.Pix)/4)
	blue := make([]byte, len(rgba.Pix)/4)

	for i, _ := range red {
		red[i] = rgba.Pix[i*4]
		green[i] = rgba.Pix[1+i*4]
		blue[i] = rgba.Pix[2+i*4]
	}

	ch1 := image.Gray{Pix: red, Stride: 21600, Rect: image.Rect(0, 0, 21600, 10800)}
	ch2 := image.Gray{Pix: green, Stride: 21600, Rect: image.Rect(0, 0, 21600, 10800)}
	ch3 := image.Gray{Pix: blue, Stride: 21600, Rect: image.Rect(0, 0, 21600, 10800)}

	f1, err := os.Create("../Downloads/red.png")
	if err != nil {
		panic(err)
	}
	f2, err := os.Create("../Downloads/green.png")
	if err != nil {
		panic(err)
	}
	f3, err := os.Create("../Downloads/blue.png")
	if err != nil {
		panic(err)
	}

	png.Encode(f1, &ch1)
	png.Encode(f2, &ch2)
	png.Encode(f3, &ch3)

	err = WriteRawFile("../Downloads/dat1", red)
	if err != nil {
		panic(err)
	}
	err = WriteRawFile("../Downloads/dat2", green)
	if err != nil {
		panic(err)
	}
	err = WriteRawFile("../Downloads/dat3", blue)
	if err != nil {
		panic(err)
	}

	_, err = ReadRawFile("../Downloads/dat1")
	if err != nil {
		panic(err)
	}
	_, err = ReadRawFile("../Downloads/dat2")
	if err != nil {
		panic(err)
	}
	_, err = ReadRawFile("../Downloads/dat3")
	if err != nil {
		panic(err)
	}

	err = WriteSnappyFile("../Downloads/dat1.snp", red)
	if err != nil {
		panic(err)
	}
	err = WriteSnappyFile("../Downloads/dat2.snp", green)
	if err != nil {
		panic(err)
	}
	err = WriteSnappyFile("../Downloads/dat3.snp", blue)
	if err != nil {
		panic(err)
	}

	_, err = ReadSnappyFile("../Downloads/dat1.snp")
	if err != nil {
		panic(err)
	}
	_, err = ReadSnappyFile("../Downloads/dat2.snp")
	if err != nil {
		panic(err)
	}
	_, err = ReadSnappyFile("../Downloads/dat3.snp")
	if err != nil {
		panic(err)
	}

	err = WriteZipFile("../Downloads/dat1.gz", red)
	if err != nil {
		panic(err)
	}
	err = WriteZipFile("../Downloads/dat2.gz", green)
	if err != nil {
		panic(err)
	}
	err = WriteZipFile("../Downloads/dat3.gz", blue)
	if err != nil {
		panic(err)
	}

	_, err = ReadZipFile("../Downloads/dat1.gz")
	if err != nil {
		panic(err)
	}
	_, err = ReadZipFile("../Downloads/dat2.gz")
	if err != nil {
		panic(err)
	}
	_, err = ReadZipFile("../Downloads/dat3.gz")
	if err != nil {
		panic(err)
	}
}
