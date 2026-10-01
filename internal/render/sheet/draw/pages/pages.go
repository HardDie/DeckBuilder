// Draw and JPEG-encode each page at the same time. One page uses seq.
package pages

import (
	"image"
	"sync"

	"github.com/HardDie/DeckBuilder/internal/render/sheet/draw/seq"
)

// Sheet is one page of faces plus its back.
type Sheet struct {
	Faces []image.Image
	Back  image.Image
}

// Image draws one page with seq. It does not encode.
// The app saves one page at a time, so this is that one page.
func Image(faces []image.Image, back image.Image) *image.RGBA {
	return seq.Image(faces, back)
}

// JPEGs returns one quality-80 JPEG per sheet, in order.
func JPEGs(sheets []Sheet) ([][]byte, error) {
	out := make([][]byte, len(sheets))
	errs := make([]error, len(sheets))
	var wg sync.WaitGroup
	wg.Add(len(sheets))
	for i, sheet := range sheets {
		i, sheet := i, sheet
		go func() {
			defer wg.Done()
			out[i], errs[i] = seq.JPEG(sheet.Faces, sheet.Back)
		}()
	}
	wg.Wait()
	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}
