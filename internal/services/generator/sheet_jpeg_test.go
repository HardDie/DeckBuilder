package generator

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"testing"

	entitiesSettings "github.com/HardDie/DeckBuilder/internal/entities/settings"
	"github.com/HardDie/DeckBuilder/internal/sheet/draw/libjpeg"
	sheetpage "github.com/HardDie/DeckBuilder/internal/sheet/page"
)

func TestSaveSheetUsesLibjpeg(t *testing.T) {
	dir := t.TempDir()
	cfg := entitiesSettings.Default()
	appPage := sheetpage.New("deck", dir, 1, 1, &cfg)
	same := sheetpage.New("deck", dir, 1, 1, &cfg)
	face := pngBytes(t, 8, 12, color.RGBA{R: 220, G: 20, B: 40, A: 255})
	back := pngBytes(t, 8, 12, color.RGBA{G: 180, B: 40, A: 255})
	for _, page := range []*sheetpage.Page{appPage, same} {
		if _, err := page.SetBacksideImageAndSave(back); err != nil {
			t.Fatal(err)
		}
		if err := page.AddImage(face); err != nil {
			t.Fatal(err)
		}
	}

	abs, cols, rows, err := saveSheet(appPage)
	if err != nil {
		t.Fatal(err)
	}
	if cols != 2 || rows != 2 {
		t.Fatalf("grid %dx%d", cols, rows)
	}
	got, err := os.ReadFile(abs)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) < 2 || got[0] != 0xff || got[1] != 0xd8 {
		t.Fatalf("not a jpeg: %d bytes", len(got))
	}

	sheet, _, _, _, err := same.Sheet()
	if err != nil {
		t.Fatal(err)
	}
	want, err := libjpeg.Encode(sheet)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("file %d bytes, libjpeg %d bytes", len(got), len(want))
	}
	var std bytes.Buffer
	if err := jpeg.Encode(&std, sheet, &jpeg.Options{Quality: 80}); err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(got, std.Bytes()) {
		t.Fatal("sheet jpeg matches image/jpeg")
	}
}

func pngBytes(t *testing.T, w, h int, c color.RGBA) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, c)
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}
