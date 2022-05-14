package image

import (
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
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

	var img image.Image
	img, _, err = image.Decode(imgFile)
	if err != nil {
		log.Fatal(err.Error())
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
func (img *Image) Draw(col, row int, imageSrc image.Image) {
	pos := image.Rect(
		col*img.Width,             // Start X
		row*img.Height,            // Start Y
		col*img.Width+img.Width,   // End X
		row*img.Height+img.Height, // End Y
	)
	draw.Draw(img.Img, pos, imageSrc, image.Point{0, 0}, draw.Src)
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
