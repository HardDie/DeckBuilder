package page_drawer

import (
	"crypto/md5"
	"errors"
	"fmt"
	"image"
	"math"
	"path/filepath"

	"github.com/disintegration/imaging"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/fs"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type PageDrawer struct {
	images   []image.Image
	backside *image.NRGBA

	index       int
	commonIndex int
	title       string
	path        string

	scale      int
	innerScale float64
	width      int
	height     int

	config *entity.SettingInfo
}

func New(title, path string, scale, commonIndex int, config *entity.SettingInfo) *PageDrawer {
	return &PageDrawer{
		index:       1,
		commonIndex: commonIndex,
		title:       title,
		path:        path,
		scale:       scale,
		config:      config,
	}
}

func (d *PageDrawer) Inherit(d2 *PageDrawer) *PageDrawer {
	d.index = d2.index + 1
	d.commonIndex = d2.commonIndex + 1
	d.title, d.path = d2.title, d2.path
	d.backside = d2.backside
	d.width, d.height = d2.width, d2.height
	return d
}

func (d *PageDrawer) IsFull() bool {
	return len(d.images) >= config.MaxCount
}
func (d *PageDrawer) IsEmpty() bool {
	return len(d.images) == 0
}
func (d *PageDrawer) GetIndex() int {
	return d.index
}
func (d *PageDrawer) Size() int {
	return len(d.images)
}

func (d *PageDrawer) AddImage(img []byte) error {
	if d.IsFull() {
		return errors.New("page is full")
	}

	cardImg, err := images.ImageFromBinary(img)
	if err != nil {
		return err
	}

	cardWidth := cardImg.Bounds().Max.X
	cardHeight := cardImg.Bounds().Max.Y
	if (cardWidth * 10) > 10_000 {
		d.innerScale = 10_000 / 10 / float64(cardWidth)
		for {
			if int(math.Trunc(float64(cardWidth)*d.innerScale)) > 10_000 {
				logger.Debug.Println("Increase scale for width:", d.innerScale)
				d.innerScale += 0.01
			}
			break
		}
	}
	if (math.Trunc(float64(cardHeight)*d.innerScale) * 7) > 10_000 {
		d.innerScale = 10_000 / 7 / float64(cardHeight)
		for {
			if int(math.Trunc(float64(cardHeight)*d.innerScale)) > 10_000 {
				logger.Debug.Println("Increase scale for height:", d.innerScale)
				d.innerScale += 0.01
			}
			break
		}
	}

	if d.scale == 0 {
		d.scale = 1
	}
	if d.innerScale == 0 {
		d.innerScale = 1
	}

	if d.width == 0 && d.height == 0 {
		d.width = int(math.Trunc(float64(cardWidth)*d.innerScale)) / d.scale
		d.height = int(math.Trunc(float64(cardHeight)*d.innerScale)) / d.scale
	}

	if d.width != cardImg.Bounds().Max.X ||
		d.height != cardImg.Bounds().Max.Y {
		cardImg = imaging.Resize(cardImg, d.width, d.height, imaging.Lanczos)
	}
	d.images = append(d.images, cardImg)
	return nil
}
func (d *PageDrawer) SetBacksideImageAndSave(img []byte) (string, error) {
	// Save image on disk
	hash := md5.Sum(img)
	name := "backside_" + d.title + "_" + fmt.Sprintf("%x", hash[0:3]) + ".png"
	savePath := filepath.Join(d.path, name)
	err := fs.CreateAndProcess(savePath, img, fs.BinToWriter)
	if err != nil {
		return "", err
	}

	// Convert binary to image
	backsideImg, err := images.ImageFromBinary(img)
	if err != nil {
		return "", err
	}

	// Make image darker
	if d.config.EnableBackShadow {
		d.backside = imaging.AdjustBrightness(backsideImg, -30)
	} else {
		d.backside = imaging.AdjustBrightness(backsideImg, 0)
	}
	return fs.PathToAbsolutePath(savePath), nil
}
func (d *PageDrawer) Save() (string, int, int, error) {
	if d.IsEmpty() {
		return "", 0, 0, nil
	}

	if d.width != d.backside.Bounds().Max.X ||
		d.height != d.backside.Bounds().Max.Y {
		d.backside = imaging.Resize(d.backside, d.width, d.height, imaging.Lanczos)
	}

	// Calculate page size
	columns, rows := utils.CalculateGridSize(len(d.images) + 1)
	// Create image
	pageImage := images.CreateImage(d.width*columns, d.height*rows)
	// Draw cards
	for i, cardImg := range d.images {
		column, row := utils.CardIdToPageCoordinates(i, columns)
		images.Draw(pageImage, column, row, cardImg)
	}
	// Draw backside image
	images.Draw(pageImage, columns-1, rows-1, d.backside)

	// Filename
	pageName := fmt.Sprintf("%d_%s_%d_%d_%dx%d.jpg", d.commonIndex, d.title, d.index, len(d.images), columns, rows)
	savePath := filepath.Join(d.path, pageName)
	// Saving on disk
	err := fs.CreateAndProcess[image.Image](savePath, pageImage, images.JpegSaveToWriter)
	if err != nil {
		return "", 0, 0, err
	}
	return fs.PathToAbsolutePath(savePath), columns, rows, nil
}
