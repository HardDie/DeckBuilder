// Lanczos, then *image.RGBA cells so draw.Src copies rows.
package rgba

import (
	"bytes"
	"image"
	"image/draw"

	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/sheet/fit"
	"github.com/HardDie/DeckBuilder/internal/sheet/grid"
)

// Image draws one page. It does not encode.
// Each cell is *image.RGBA before it is drawn, so draw.Src uses a row copy.
func Image(faces []image.Image, back image.Image) *image.RGBA {
	faces, back = resizeAll(faces, back)
	cols, rows := grid.Size(len(faces) + 1)
	cellW := faces[0].Bounds().Dx()
	cellH := faces[0].Bounds().Dy()
	page := images.CreateImage(cellW*cols, cellH*rows)
	for i, face := range faces {
		col, row := grid.Slot(i, cols)
		images.Draw(page, col, row, face)
	}
	images.Draw(page, cols-1, rows-1, back)
	return page
}

// JPEG matches seq pixels.
func JPEG(faces []image.Image, back image.Image) ([]byte, error) {
	return encode(Image(faces, back))
}

func resizeAll(faces []image.Image, back image.Image) ([]image.Image, image.Image) {
	b := faces[0].Bounds().Max
	_, size := fit.Apply(b.X, b.Y, 1, fit.Size{})
	out := make([]image.Image, len(faces))
	for i, face := range faces {
		out[i] = asRGBA(fit.Resize(face, size.Width, size.Height))
	}
	return out, asRGBA(fit.Resize(back, size.Width, size.Height))
}

// asRGBA makes the cell the type draw.Src copies with copy.
// imaging.Resize returns *image.NRGBA, which takes the premultiply loop instead.
func asRGBA(src image.Image) *image.RGBA {
	if dst, ok := src.(*image.RGBA); ok && dst.Bounds().Min == (image.Point{}) {
		return dst
	}
	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)
	return dst
}

func encode(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if err := images.JpegSaveToWriter(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
