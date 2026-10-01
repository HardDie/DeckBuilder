package jpegli

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"testing"
)

func TestEncode(t *testing.T) {
	img := pattern()
	got, err := Encode(img)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) < 2 || got[0] != 0xff || got[1] != 0xd8 {
		t.Fatalf("not a jpeg: %d bytes", len(got))
	}
	if _, err := jpeg.Decode(bytes.NewReader(got)); err != nil {
		t.Fatal(err)
	}
	var std bytes.Buffer
	if err := jpeg.Encode(&std, img, &jpeg.Options{Quality: 80}); err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(got, std.Bytes()) {
		t.Fatal("bytes match image/jpeg")
	}
}

func pattern() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, 16, 16))
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			img.SetRGBA(x, y, color.RGBA{R: uint8(x * 15), G: uint8(y * 12), B: 40, A: 255})
		}
	}
	return img
}
