package images

import (
	"bytes"
	stderrors "errors"
	"fmt"
	"image"
	"image/draw"
	"math"

	_ "golang.org/x/image/bmp"  // ToWebP accepts BMP
	_ "golang.org/x/image/tiff" // and TIFF
	_ "golang.org/x/image/webp" // DecodeConfig reads WebP headers

	"github.com/HardDie/DeckBuilder/internal/apperr"
	"github.com/HardDie/DeckBuilder/internal/images/webp"
)

const (
	// webpQuality is the stored lossy quality (ADR 027): 45–48 dB on real card art.
	webpQuality = 95
	// lossyMinPSNR is the least a lossy result may keep. Lossy WebP stores color at
	// half resolution, so sharp colored detail (pixel art, thin colored lines)
	// drops to about 20 dB; such an image is stored lossless instead.
	lossyMinPSNR = 40
)

// ToWebP converts an image in any accepted format (PNG, JPEG, GIF, WebP, BMP, TIFF)
// to the stored format, WebP (ADR 027): lossy at quality 95, or lossless when
// lossy would keep less than 40 dB.
// A JPEG is turned upright by its EXIF orientation. GIF uses its first frame.
// The result is decoded again and compared with the source before it is returned.
func ToWebP(input []byte) ([]byte, error) {
	cfg, format, err := image.DecodeConfig(bytes.NewReader(input))
	if stderrors.Is(err, image.ErrFormat) {
		return nil, apperr.ErrUnsupportedImage
	}
	if err != nil {
		return nil, damaged(err)
	}
	if int64(cfg.Width)*int64(cfg.Height) > maxImagePixels {
		return nil, apperr.Withf(apperr.ErrImageTooLarge, "image is too large: %dx%d, max %d megapixels", cfg.Width, cfg.Height, maxImagePixels>>20)
	}
	if cfg.Width > webp.MaxSide || cfg.Height > webp.MaxSide {
		return nil, apperr.Withf(apperr.ErrImageTooLarge, "image is too large: %dx%d, max %d pixels on a side", cfg.Width, cfg.Height, webp.MaxSide)
	}

	src, err := decodeAny(input, format)
	if err != nil {
		return nil, damaged(err)
	}
	if format == "jpeg" {
		src = orient(src, jpegOrientation(input))
	}

	out, err := encodeWebP(src)
	if err != nil {
		return nil, fmt.Errorf("convert %s %dx%d to webp: %w", format, cfg.Width, cfg.Height, err)
	}
	return out, nil
}

// encodeWebP tries lossy q95 and keeps it if it holds at least lossyMinPSNR.
// Otherwise it encodes lossless, which must decode exactly.
func encodeWebP(src *image.NRGBA) ([]byte, error) {
	lossy, err := webp.Encode(src, webpQuality)
	if err != nil {
		return nil, err
	}
	p, err := compare(src, lossy)
	if err != nil {
		return nil, fmt.Errorf("verify lossy: %w", err)
	}
	if p >= lossyMinPSNR {
		return lossy, nil
	}
	lossless, err := webp.EncodeLossless(src)
	if err != nil {
		return nil, err
	}
	p, err = compare(src, lossless)
	if err != nil {
		return nil, fmt.Errorf("verify lossless: %w", err)
	}
	if !math.IsInf(p, 1) {
		return nil, fmt.Errorf("verify lossless: PSNR %.1f dB, want an exact match", p)
	}
	return lossless, nil
}

// decodeAny decodes input to NRGBA. WebP goes through libwebp, the rest through Go.
func decodeAny(input []byte, format string) (*image.NRGBA, error) {
	if format == "webp" {
		return webp.Decode(input)
	}
	img, _, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, err
	}
	return toNRGBA(img), nil
}

// toNRGBA returns img as NRGBA whose Pix starts at its top-left pixel.
func toNRGBA(img image.Image) *image.NRGBA {
	b := img.Bounds()
	if n, ok := img.(*image.NRGBA); ok && b.Min == (image.Point{}) {
		return n
	}
	dst := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Rect, img, b.Min, draw.Src)
	return dst
}

// compare decodes out and returns its PSNR against src. A different size is an error.
func compare(src *image.NRGBA, out []byte) (float64, error) {
	got, err := webp.Decode(out)
	if err != nil {
		return 0, err
	}
	if got.Rect.Size() != src.Rect.Size() {
		return 0, fmt.Errorf("size %v, want %v", got.Rect.Size(), src.Rect.Size())
	}
	return psnr(src, got), nil
}

// psnr compares two same-size NRGBA images in premultiplied color plus alpha,
// so the hidden color of a fully transparent pixel does not count.
// Identical images give +Inf.
func psnr(a, b *image.NRGBA) float64 {
	var sum float64
	for i := 0; i+3 < len(a.Pix); i += 4 {
		aa, ba := float64(a.Pix[i+3]), float64(b.Pix[i+3])
		for c := 0; c < 3; c++ {
			d := float64(a.Pix[i+c])*aa/255 - float64(b.Pix[i+c])*ba/255
			sum += d * d
		}
		sum += (aa - ba) * (aa - ba)
	}
	if sum == 0 {
		return math.Inf(1)
	}
	return 10 * math.Log10(255*255/(sum/float64(len(a.Pix))))
}
