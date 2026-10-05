package images

import (
	"bytes"
	"encoding/binary"
	stderrors "errors"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math"
	"strings"
	"testing"

	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"

	"github.com/HardDie/DeckBuilder/internal/apperr"
	"github.com/HardDie/DeckBuilder/internal/images/webp"
)

// picture is detailed art; with alpha its left third is fully transparent.
func picture(alpha bool) *image.NRGBA {
	const w, h = 96, 64
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := color.NRGBA{uint8(x * 255 / w), uint8(y * 255 / h), uint8((x ^ y) * 4), 255}
			if alpha && x < w/3 {
				c.A = 0
			}
			img.SetNRGBA(x, y, c)
		}
	}
	return img
}

func encode(t *testing.T, f func(*bytes.Buffer) error) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := f(&buf); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// Every accepted format becomes a WebP of the same size that matches its decoded source.
func TestToWebPFormats(t *testing.T) {
	opaque, alpha := picture(false), picture(true)
	tests := []struct {
		name      string
		input     []byte
		wantAlpha bool
	}{
		{"png", encode(t, func(b *bytes.Buffer) error { return png.Encode(b, opaque) }), false},
		{"png with alpha", encode(t, func(b *bytes.Buffer) error { return png.Encode(b, alpha) }), true},
		{"jpeg", encode(t, func(b *bytes.Buffer) error { return jpeg.Encode(b, opaque, &jpeg.Options{Quality: 95}) }), false},
		{"gif", encode(t, func(b *bytes.Buffer) error {
			p := image.NewPaletted(opaque.Rect, palette.Plan9)
			draw.Draw(p, p.Rect, opaque, image.Point{}, draw.Src)
			return gif.Encode(b, p, nil)
		}), false},
		{"bmp", encode(t, func(b *bytes.Buffer) error { return bmp.Encode(b, opaque) }), false},
		{"tiff", encode(t, func(b *bytes.Buffer) error { return tiff.Encode(b, opaque, nil) }), false},
		{"webp", func() []byte {
			data, err := webp.Encode(opaque, 100)
			if err != nil {
				t.Fatal(err)
			}
			return data
		}(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := ToWebP(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if !webp.IsWebP(out) {
				t.Fatal("result is not WebP")
			}
			got, err := webp.Decode(out)
			if err != nil {
				t.Fatal(err)
			}
			src, err := decodeAny(tt.input, strings.TrimSpace(formatOf(t, tt.input)))
			if err != nil {
				t.Fatal(err)
			}
			if got.Rect != src.Rect {
				t.Fatalf("bounds %v, want %v", got.Rect, src.Rect)
			}
			if p := psnr(src, got); p < 40 {
				t.Fatalf("PSNR %.1f dB against the source, want >= 40", p)
			}
			if a := got.NRGBAAt(0, 0).A; tt.wantAlpha != (a == 0) {
				t.Fatalf("left pixel alpha %d, transparent wanted %v", a, tt.wantAlpha)
			}
		})
	}
}

func formatOf(t *testing.T, data []byte) string {
	t.Helper()
	_, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	return format
}

// A JPEG whose EXIF says "turn 90° clockwise" is stored upright: taller than wide,
// with the source's left (red) half on top.
func TestToWebPAppliesOrientation(t *testing.T) {
	src := image.NewNRGBA(image.Rect(0, 0, 80, 40))
	for y := 0; y < 40; y++ {
		for x := 0; x < 80; x++ {
			c := color.NRGBA{B: 220, A: 255}
			if x < 40 {
				c = color.NRGBA{R: 220, A: 255}
			}
			src.SetNRGBA(x, y, c)
		}
	}
	jpg := encode(t, func(b *bytes.Buffer) error { return jpeg.Encode(b, src, &jpeg.Options{Quality: 95}) })
	out, err := ToWebP(withAPP1(jpg, app1(binary.BigEndian, 6)))
	if err != nil {
		t.Fatal(err)
	}
	got, err := webp.Decode(out)
	if err != nil {
		t.Fatal(err)
	}
	if got.Rect.Size() != image.Pt(40, 80) {
		t.Fatalf("size %v, want 40x80", got.Rect.Size())
	}
	if top, bottom := got.NRGBAAt(20, 10), got.NRGBAAt(20, 70); top.R < 150 || bottom.B < 150 {
		t.Fatalf("top %v should be red, bottom %v blue", top, bottom)
	}
}

func TestToWebPRejects(t *testing.T) {
	truncated := encode(t, func(b *bytes.Buffer) error { return png.Encode(b, picture(false)) })
	truncated = truncated[:len(truncated)/2]
	tests := []struct {
		name    string
		input   []byte
		is      error
		message string
	}{
		{"not an image", []byte("hello"), apperr.ErrUnsupportedImage, ""},
		{"empty", nil, apperr.ErrUnsupportedImage, ""},
		{"damaged png", truncated, apperr.ErrUnsupportedImage, "damaged"},
		{"over the pixel cap", pngHeader(16000, 16000), apperr.ErrImageTooLarge, "megapixels"},
		{"side over the WebP limit", pngHeader(16384, 100), apperr.ErrImageTooLarge, "16383 pixels on a side"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ToWebP(tt.input)
			if !stderrors.Is(err, tt.is) {
				t.Fatalf("got %v, want %v", err, tt.is)
			}
			if !strings.Contains(apperr.Message(err), tt.message) {
				t.Fatalf("message %q, want it to mention %q", apperr.Message(err), tt.message)
			}
		})
	}
}

func TestPSNR(t *testing.T) {
	a := picture(true)
	if p := psnr(a, a); p < 1e9 {
		t.Fatalf("identical images: %v, want +Inf", p)
	}
	b := image.NewNRGBA(a.Rect)
	copy(b.Pix, a.Pix)
	b.Pix[0] = 0 // red of a fully transparent pixel: hidden, so it does not count
	if p := psnr(a, b); p < 1e9 {
		t.Fatalf("hidden color changed PSNR to %v", p)
	}
}

// smooth is photo-like art: gentle gradients, which lossy keeps well above 40 dB.
func smooth() *image.NRGBA {
	const w, h = 128, 96
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetNRGBA(x, y, color.NRGBA{uint8(x * 255 / w), uint8(y * 255 / h), 128, 255})
		}
	}
	return img
}

// thinLines is 1-px red strokes on blue: lossy halves color resolution and blurs them.
func thinLines() *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, 128, 96))
	for y := 0; y < 96; y++ {
		for x := 0; x < 128; x++ {
			c := color.NRGBA{20, 40, 200, 255}
			if x%7 == 0 || y%9 == 0 {
				c = color.NRGBA{230, 30, 30, 255}
			}
			img.SetNRGBA(x, y, c)
		}
	}
	return img
}

// pixelArt is 4×4 blocks of saturated colors.
func pixelArt() *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, 128, 96))
	pal := []color.NRGBA{{255, 0, 0, 255}, {0, 200, 0, 255}, {0, 0, 255, 255}, {255, 220, 0, 255}, {255, 255, 255, 255}, {0, 0, 0, 255}}
	for y := 0; y < 96; y++ {
		for x := 0; x < 128; x++ {
			img.SetNRGBA(x, y, pal[((x/4)*7+(y/4)*13)%len(pal)])
		}
	}
	return img
}

// Real-looking art stays lossy; sharp colored detail falls back to lossless and stays exact.
func TestToWebPChoosesLossyOrLossless(t *testing.T) {
	tests := []struct {
		name     string
		img      *image.NRGBA
		lossless bool
	}{
		{"smooth art", smooth(), false},
		{"thin colored lines", thinLines(), true},
		{"pixel art", pixelArt(), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := ToWebP(encode(t, func(b *bytes.Buffer) error { return png.Encode(b, tt.img) }))
			if err != nil {
				t.Fatal(err)
			}
			if got := bytes.Contains(out[:16], []byte("VP8L")); got != tt.lossless {
				t.Fatalf("lossless %v, want %v (header % x)", got, tt.lossless, out[:16])
			}
			got, err := webp.Decode(out)
			if err != nil {
				t.Fatal(err)
			}
			p := psnr(tt.img, got)
			if tt.lossless && !math.IsInf(p, 1) {
				t.Fatalf("lossless result is not exact: %.1f dB", p)
			}
			if !tt.lossless && p < lossyMinPSNR {
				t.Fatalf("lossy result %.1f dB, want >= %d", p, lossyMinPSNR)
			}
		})
	}
}
