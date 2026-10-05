package webp

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"testing"
)

// gradient is an opaque test picture, or one whose left half is transparent.
func gradient(w, h int, alpha bool) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := color.NRGBA{uint8(x * 255 / w), uint8(y * 255 / h), 128, 255}
			if alpha && x < w/2 {
				c.A = 0
			}
			img.SetNRGBA(x, y, c)
		}
	}
	return img
}

func TestEncodeDecode(t *testing.T) {
	tests := []struct {
		name      string
		img       *image.NRGBA
		wantAlpha bool
	}{
		{"opaque", gradient(64, 48, false), false},
		{"transparent half", gradient(64, 48, true), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Encode(tt.img, 95)
			if err != nil {
				t.Fatal(err)
			}
			if !IsWebP(data) {
				t.Fatalf("not a WebP header: % x", data[:12])
			}
			if got := bytes.Contains(data, []byte("ALPH")); got != tt.wantAlpha {
				t.Fatalf("alpha chunk %v, want %v", got, tt.wantAlpha)
			}
			w, h, err := Size(data)
			if err != nil || w != 64 || h != 48 {
				t.Fatalf("Size = %d, %d, %v", w, h, err)
			}
			got, err := Decode(data)
			if err != nil {
				t.Fatal(err)
			}
			if got.Rect != tt.img.Rect {
				t.Fatalf("bounds %v", got.Rect)
			}
			if a := got.NRGBAAt(0, 0).A; tt.wantAlpha && a != 0 {
				t.Fatalf("transparent pixel came back with alpha %d", a)
			}
		})
	}
}

// A sub-image starts at its own top-left pixel, not at the parent's.
func TestEncodeSubImage(t *testing.T) {
	parent := gradient(64, 48, false)
	sub := parent.SubImage(image.Rect(32, 0, 64, 48)).(*image.NRGBA)
	data, err := Encode(sub, 100)
	if err != nil {
		t.Fatal(err)
	}
	got, err := Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	want := sub.NRGBAAt(32, 0).R
	if r := got.NRGBAAt(0, 0).R; r < want-8 || r > want+8 {
		t.Fatalf("first pixel red %d, want about %d (the sub-image's, not the parent's)", r, want)
	}
}

func TestEncodeBadSize(t *testing.T) {
	for _, r := range []image.Rectangle{image.Rect(0, 0, 0, 10), image.Rect(0, 0, MaxSide+1, 1)} {
		if _, err := Encode(image.NewNRGBA(r), 95); !errors.Is(err, ErrEncode) {
			t.Errorf("size %v: got %v, want ErrEncode", r.Size(), err)
		}
	}
}

func TestDecodeBadData(t *testing.T) {
	for _, data := range [][]byte{nil, []byte("RIFF\x00\x00\x00\x00WEBPjunk"), []byte("not an image")} {
		if _, err := Decode(data); !errors.Is(err, ErrDecode) {
			t.Errorf("%q: got %v, want ErrDecode", data, err)
		}
	}
}

func TestIsWebP(t *testing.T) {
	for data, want := range map[string]bool{
		"RIFF\x10\x00\x00\x00WEBPVP8 ": true,
		"RIFF\x10\x00\x00\x00WAVEfmt ": false,
		"\x89PNG\r\n\x1a\n":            false,
		"RIFF":                         false,
	} {
		if got := IsWebP([]byte(data)); got != want {
			t.Errorf("IsWebP(%q) = %v, want %v", data, got, want)
		}
	}
}

// Lossless keeps every visible pixel, including partly transparent ones.
func TestEncodeLossless(t *testing.T) {
	img := gradient(64, 48, true)
	img.SetNRGBA(40, 10, color.NRGBA{200, 100, 50, 128})
	data, err := EncodeLossless(img)
	if err != nil {
		t.Fatal(err)
	}
	got, err := Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	for y := 0; y < 48; y++ {
		for x := 0; x < 64; x++ {
			want := img.NRGBAAt(x, y)
			if want.A == 0 {
				if a := got.NRGBAAt(x, y).A; a != 0 {
					t.Fatalf("(%d,%d) alpha %d, want 0", x, y, a)
				}
				continue
			}
			if g := got.NRGBAAt(x, y); g != want {
				t.Fatalf("(%d,%d) = %v, want %v", x, y, g, want)
			}
		}
	}
}
