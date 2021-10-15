package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	Souls     []string
	Character []string
	Loot      []string
	Monster   []string
	Room      []string
	Eternal   []string
	Treasure  []string
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
	img, err := png.Decode(imgFile)
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
	fmt.Printf("Building %s image...\n", title)
	out, err := os.Create(title)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = png.Encode(out, img.Img)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Image %s ready!\n", title)
}

func BestSize(count int) (images, cols, rows int) {
	cols = 10
	rows = 7
	images = cols * rows
	for r := 2; r <= 7; r++ {
		for c := 2; c <= 10; c++ {
			posible := c * r
			if posible < images && posible >= count {
				images = posible
				cols = c
				rows = r
			}
		}
	}
	images = count
	return
}
func BuildDeck(slice *[]string, path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err.Error())
	}

	if len(*slice) > 0 {
		*slice = append(*slice, "EMPTY")
	}

	for _, img := range files {
		switch img.Name() {
		case "r-the_enigma_back.png": // Back enigma image
			continue
		case "r-the_revenant.png": // Back anima sola
			continue
		}
		*slice = append(*slice, path+"/"+img.Name())
	}
}
func Crawl(path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, file := range files {
		switch file.Name() {
		case "Bonus Soul":
			BuildDeck(&Souls, path+"/"+file.Name())
		case "Character":
			BuildDeck(&Character, path+"/"+file.Name())
		case "Loot":
			BuildDeck(&Loot, path+"/"+file.Name())
		case "Monster":
			BuildDeck(&Monster, path+"/"+file.Name())
		case "Room":
			BuildDeck(&Room, path+"/"+file.Name())
		case "Starting Item":
			BuildDeck(&Eternal, path+"/"+file.Name())
		case "Treasure":
			BuildDeck(&Treasure, path+"/"+file.Name())
		default:
			Crawl(path + "/" + file.Name())
		}
	}
}
func BuildDeckImage(prefix string, slice []string, postfix string) {
	if len(slice) > 70 {

	}

	images, cols, rows := BestSize(len(slice))

	var imgSlice []image.Image
	for _, path := range slice {
		if path == "EMPTY" {
			imgSlice = append(imgSlice, nil)
			continue
		}
		imgSlice = append(imgSlice, OpenImage(path))
	}

	bound := imgSlice[0].Bounds().Max

	rgba := CreateImage(bound.X, bound.Y, cols, rows)
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			if len(imgSlice) <= (row*cols + col) {
				continue
			}
			img := imgSlice[row*cols+col]
			if img == nil {
				continue
			}
			rgba.Draw(col, row, img)
		}
	}
	rgba.SaveImage(fmt.Sprintf("result_png/%s%d_%dx%d%s.png", prefix, images, cols, rows, postfix))
}

func main() {
	root := "/path/to/googledrive/Four Souls"
	Crawl(root + "/Second Edition")
	Crawl(root + "/Promo")

	slices := []struct {
		Title string
		Slice []string
	}{
		{"bonus_soul", Souls},
		{"character", Character},
		{"loot", Loot},
		{"monster", Monster},
		{"room", Room},
		{"eternal", Eternal},
		{"treasure", Treasure},
	}

	for _, slice := range slices {
		fmt.Println(slice.Title, len(slice.Slice))
	}

	fmt.Println()
	for _, slice := range slices {
		wg := sync.WaitGroup{}

		if len(slice.Slice) < 70 {
			wg.Add(1)
			go func(title string, slice []string) {
				BuildDeckImage(title, slice, "_v2")
				wg.Done()
			}(slice.Title+"_", slice.Slice)
		} else {
			left := 0
			right := 70
			num := 1
			for true {
				if right > len(slice.Slice) {
					right = len(slice.Slice)
				}
				wg.Add(1)
				go func(title string, slice []string) {
					BuildDeckImage(title, slice, "_v2")
					wg.Done()
				}(fmt.Sprintf("%s_%d_", slice.Title, num), slice.Slice[left:right])

				if right == len(slice.Slice) {
					break
				}
				left += 70
				right += 70
				num++
			}
		}

		fmt.Println("Start wait...")
		wg.Wait()
		fmt.Println("Done!")
	}
}
