package main

import (
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"
)

type (
	Image struct {
		Img    *image.RGBA
		Width  int
		Height int
	}
)

func OpenImage(path string) image.Image {
	imgFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err.Error())
	}

	tokens := strings.Split(path, ".")
	ext := strings.ToLower(tokens[len(tokens)-1])

	var img image.Image
	switch ext {
	case "png":
		img, err = png.Decode(imgFile)
	case "jpg", "jpeg":
		img, err = jpeg.Decode(imgFile)
	default:
		log.Fatal("Unknown file format:", ext)
	}
	if err != nil {
		log.Fatal(err.Error(), path)
	}
	return img
}
func CreateImage(width, height, cols, rows int) *Image {
	return &Image{
		Img:    image.NewRGBA(image.Rect(0, 0, width*cols, height*rows)),
		Width:  width,
		Height: height,
	}
}
func (img *Image) Draw(col int, row int, imageSrc image.Image) {
	pos := image.Rect(
		col*img.Width,             // Start X
		row*img.Height,            // Start Y
		col*img.Width+img.Width,   // End X
		row*img.Height+img.Height, // End Y
	)
	draw.Draw(img.Img, pos, imageSrc, image.Point{0, 0}, draw.Src)
	return
}
func (img *Image) SaveImage(title string) {
	out, err := os.Create(title)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = png.Encode(out, img.Img)
	if err != nil {
		log.Fatal(err.Error())
	}
}
