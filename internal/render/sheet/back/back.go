// Deck-back file and the copy drawn on the sheet.
package back

import (
	"crypto/md5"
	"fmt"
	"image"
	"path/filepath"

	"github.com/disintegration/imaging"

	"github.com/HardDie/DeckBuilder/internal/fs"
)

// Write stores the original bytes as backside_<title>_<hash>.png.
func Write(dir, title string, raw []byte) (string, error) {
	sum := md5.Sum(raw)
	name := "backside_" + title + "_" + fmt.Sprintf("%x", sum[0:3]) + ".png"
	path := filepath.Join(dir, name)
	if err := fs.CreateAndProcess(path, raw, fs.BinToWriter); err != nil {
		return "", err
	}
	return fs.PathToAbsolutePath(path), nil
}

// Shade darkens by 30 when shadow is on.
// Brightness 0 still converts the image.
func Shade(img image.Image, shadow bool) *image.NRGBA {
	if shadow {
		return imaging.AdjustBrightness(img, -30)
	}
	return imaging.AdjustBrightness(img, 0)
}
