package main

import (
	"image"
	"os"

	"image/color"
	"image/png"
)

func load(fname string) ([][]Color, int, int) {
	file, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, err2 := png.Decode(file)

	if err2 != nil {
		panic(err2)
	}

	var new_img [][]Color

	for y := 0; y < img.Bounds().Max.Y; y++ {
		var row []Color
		for x := 0; x < img.Bounds().Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			row = append(row, Color{
				int(r),
				int(g),
				int(b),
			})
		}

		new_img = append(new_img, row)
	}

	return new_img, img.Bounds().Max.X, img.Bounds().Max.Y
}

func save(fname string, img [][]Color) {
	new_img := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{len(img[0]), len(img)},
	})

	for y, row := range img {
		for x, pix := range row {
			new_img.Set(x, y, color.RGBA{uint8(pix.r), uint8(pix.g), uint8(pix.b), 255})
		}
	}

	file, err := os.Create(fname)
	if err != nil {
		panic(err)
	}

	if err := png.Encode(file, new_img); err != nil {
		file.Close()
		panic(err)
	}

	if err := file.Close(); err != nil {
		panic(err)
	}

	defer file.Close()
}
