package stbresize

import (
	"bytes"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math"
	"testing"

	"github.com/disintegration/imaging"
)

const (
	srcW, srcH = 262, 192
	dstW, dstH = 200, 147
)

// pattern is detailed art, so resampling differences show.
// With transparentLeft, the left third is fully transparent and holds black RGB,
// and the next sixth is half transparent.
func pattern(transparentLeft bool) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, srcW, srcH))
	for y := 0; y < srcH; y++ {
		for x := 0; x < srcW; x++ {
			c := color.NRGBA{uint8(x * 255 / srcW), uint8(y * 255 / srcH), uint8((x ^ y) & 255), 255}
			if transparentLeft && x < srcW/3 {
				c = color.NRGBA{}
			} else if transparentLeft && x < srcW/2 {
				c.A = 128
			}
			img.SetNRGBA(x, y, c)
		}
	}
	return img
}

func decode(t *testing.T, enc func(*bytes.Buffer) error) image.Image {
	t.Helper()
	var buf bytes.Buffer
	if err := enc(&buf); err != nil {
		t.Fatal(err)
	}
	img, _, err := image.Decode(&buf)
	if err != nil {
		t.Fatal(err)
	}
	return img
}

// psnr compares premultiplied colors, so hidden colors of transparent pixels do not count.
func psnr(a, b *image.NRGBA) float64 {
	var sum float64
	for i := 0; i < len(a.Pix); i += 4 {
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

// Every layout the PNG, JPEG, and GIF decoders produce resizes like imaging's Catmull-Rom.
func TestResizeMatchesImaging(t *testing.T) {
	opaque, alpha := pattern(false), pattern(true)
	premul := image.NewRGBA(alpha.Bounds())
	draw.Draw(premul, premul.Bounds(), alpha, image.Point{}, draw.Src)
	gray := image.NewGray(opaque.Bounds())
	draw.Draw(gray, gray.Bounds(), opaque, image.Point{}, draw.Src)
	deep := image.NewNRGBA64(opaque.Bounds())
	draw.Draw(deep, deep.Bounds(), opaque, image.Point{}, draw.Src)

	tests := []struct {
		name string
		img  image.Image
	}{
		{"png_opaque", decode(t, func(b *bytes.Buffer) error { return png.Encode(b, opaque) })},
		{"png_transparent", decode(t, func(b *bytes.Buffer) error { return png.Encode(b, alpha) })},
		{"png_16bit", decode(t, func(b *bytes.Buffer) error { return png.Encode(b, deep) })},
		{"png_gray", decode(t, func(b *bytes.Buffer) error { return png.Encode(b, gray) })},
		{"premultiplied_rgba", premul},
		{"jpeg", decode(t, func(b *bytes.Buffer) error { return jpeg.Encode(b, opaque, &jpeg.Options{Quality: 90}) })},
		{"gif", decode(t, func(b *bytes.Buffer) error {
			p := image.NewPaletted(opaque.Bounds(), palette.Plan9)
			draw.Draw(p, p.Bounds(), opaque, image.Point{}, draw.Src)
			return gif.Encode(b, p, nil)
		})},
		{"sub_image", opaque.SubImage(image.Rect(10, 10, srcW-10, srcH-10))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Resize(tt.img, dstW, dstH)
			if got.Bounds() != image.Rect(0, 0, dstW, dstH) {
				t.Fatalf("bounds %v", got.Bounds())
			}
			want := imaging.Resize(tt.img, dstW, dstH, imaging.CatmullRom)
			if p := psnr(want, got); p < 50 {
				t.Fatalf("PSNR %.1f dB vs imaging, want >= 50", p)
			}
		})
	}
}

// White art on a transparent background whose hidden pixels are black
// must not get a dark fringe.
func TestResizeTransparentEdgesStayClean(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, srcW, srcH))
	for y := 0; y < srcH; y++ {
		for x := srcW / 2; x < srcW; x++ {
			img.SetNRGBA(x, y, color.NRGBA{255, 255, 255, 255})
		}
	}
	got := Resize(img, dstW, dstH)
	for i := 0; i < len(got.Pix); i += 4 {
		if got.Pix[i+3] > 0 && got.Pix[i] < 250 {
			t.Fatalf("dark fringe: RGB %d at alpha %d", got.Pix[i], got.Pix[i+3])
		}
	}
}

func TestResizeEmpty(t *testing.T) {
	if got := Resize(image.NewNRGBA(image.Rect(0, 0, 0, 0)), dstW, dstH); got.Bounds().Dx() != dstW {
		t.Fatalf("empty source: %v", got.Bounds())
	}
	if got := Resize(pattern(false), 0, 0); !got.Bounds().Empty() {
		t.Fatalf("zero target: %v", got.Bounds())
	}
}
