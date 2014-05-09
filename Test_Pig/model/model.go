package pig

import (
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"strconv"
)

func Print(url string, n int) {
	sf, err := os.Open(url)
	if err != nil {
		log.Fatal(err)
	}
	defer sf.Close()
	src, err := png.Decode(sf)
	if err != nil {
		log.Fatal(err)
	}
	df, err := os.Create(strconv.Itoa(n) + url)
	if err != nil {
		log.Fatal(err)
	}
	defer df.Close()

	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)
	err = png.Encode(df, dst)
	if err != nil {
		log.Fatal(err)
	}
}
