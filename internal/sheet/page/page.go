// One TTS custom-deck page: faces, a back file, and a JPEG sheet.
package page

import (
	"errors"
	"fmt"
	"image"
	"path/filepath"

	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesSettings "github.com/HardDie/DeckBuilder/internal/entities/settings"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/sheet/back"
	"github.com/HardDie/DeckBuilder/internal/sheet/fit"
	"github.com/HardDie/DeckBuilder/internal/sheet/paint"
)

// Page is one sheet. index starts at 1.
type Page struct {
	faces []image.Image
	back  image.Image

	index       int
	commonIndex int
	title       string
	path        string

	scale int
	size  fit.Size

	settings *entitiesSettings.Settings
}

func New(title, path string, scale, commonIndex int, settings *entitiesSettings.Settings) *Page {
	return &Page{
		index:       1,
		commonIndex: commonIndex,
		title:       title,
		path:        path,
		scale:       scale,
		settings:    settings,
	}
}

// Inherit turns a blank page into the next page of prev.
// Faces, scale, and settings stay on prev.
func (p *Page) Inherit(prev *Page) *Page {
	p.index = prev.index + 1
	p.commonIndex = prev.commonIndex + 1
	p.title, p.path = prev.title, prev.path
	p.back = prev.back
	p.size.Width = prev.size.Width
	p.size.Height = prev.size.Height
	p.size.Set = prev.size.Set
	return p
}

func (p *Page) IsFull() bool {
	return len(p.faces) >= config.MaxCount
}

func (p *Page) IsEmpty() bool {
	return len(p.faces) == 0
}

func (p *Page) GetIndex() int {
	return p.index
}

func (p *Page) Size() int {
	return len(p.faces)
}

func (p *Page) AddImage(img []byte) error {
	if p.IsFull() {
		return errors.New("page is full")
	}
	card, err := images.ImageFromBinary(img)
	if err != nil {
		return err
	}
	bounds := card.Bounds().Max
	p.scale, p.size = fit.Apply(bounds.X, bounds.Y, p.scale, p.size)
	card = fit.Resize(card, p.size.Width, p.size.Height)
	p.faces = append(p.faces, card)
	return nil
}

// SetBacksideImageAndSave writes the original bytes, then keeps a shaded copy.
func (p *Page) SetBacksideImageAndSave(img []byte) (string, error) {
	abs, err := back.Write(p.path, p.title, img)
	if err != nil {
		return "", err
	}
	decoded, err := images.ImageFromBinary(img)
	if err != nil {
		return "", err
	}
	p.back = back.Shade(decoded, p.settings.EnableBackShadow)
	return abs, nil
}

// Save writes the JPEG. An empty page writes nothing.
func (p *Page) Save() (string, int, int, error) {
	if p.IsEmpty() {
		return "", 0, 0, nil
	}
	p.back = fit.Resize(p.back, p.size.Width, p.size.Height)
	sheet, cols, rows := paint.Canvas(p.size.Width, p.size.Height, p.faces, p.back)
	name := fmt.Sprintf("%d_%s_%d_%d_%dx%d.jpg", p.commonIndex, p.title, p.index, len(p.faces), cols, rows)
	abs, err := paint.WriteJPEG(filepath.Join(p.path, name), sheet)
	if err != nil {
		return "", 0, 0, err
	}
	return abs, cols, rows, nil
}
