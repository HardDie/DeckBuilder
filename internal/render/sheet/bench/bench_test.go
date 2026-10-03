// Benchmarks for one full sheet: 69 faces plus the back on a 10×7 page.
// Inputs are the 1312×962 PNGs in testdata/input. The cell is 1000×733.
// Only the code the app runs is measured; ADR 017 records the earlier draw trials.
//
//	go test -tags=nomain -run '^$' -bench . -benchtime 5x ./internal/render/sheet/bench/
package bench

import (
	"bytes"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/HardDie/DeckBuilder/internal/fs"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/back"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/draw/libjpeg"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/fit"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/paint"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/write"
)

const (
	faceCount = 69
	// The cell layout picks for 1312×962 faces on a 10-wide page.
	cellW = 1000
	cellH = 733
)

type input struct {
	faces    []image.Image
	back     image.Image
	rawFaces [][]byte
	rawBack  []byte
}

var (
	loadOnce sync.Once
	loaded   input
)

func load(b *testing.B) input {
	b.Helper()
	loadOnce.Do(func() {
		dir := filepath.Join("testdata", "input")
		read := func(name string) []byte {
			body, err := os.ReadFile(filepath.Join(dir, name))
			if err != nil {
				b.Fatal(err)
			}
			return body
		}
		names := []string{"face_a.png", "face_b.png", "face_wide.png"}
		raw := make([][]byte, len(names))
		dec := make([]image.Image, len(names))
		for i, name := range names {
			raw[i] = read(name)
			img, err := images.ImageFromBinary(raw[i])
			if err != nil {
				b.Fatal(err)
			}
			dec[i] = img
		}
		loaded.rawBack = read("back.png")
		back, err := images.ImageFromBinary(loaded.rawBack)
		if err != nil {
			b.Fatal(err)
		}
		loaded.back = back
		for i := 0; i < faceCount; i++ {
			loaded.faces = append(loaded.faces, dec[i%len(dec)])
			loaded.rawFaces = append(loaded.rawFaces, raw[i%len(raw)])
		}
	})
	return loaded
}

// parallel runs fn(i) for i in [0, n) on GOMAXPROCS workers, like write.run.
func parallel(n int, fn func(i int)) {
	next := make(chan int)
	var wg sync.WaitGroup
	workers := min(runtime.GOMAXPROCS(0), n)
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for i := range next {
				fn(i)
			}
		}()
	}
	for i := 0; i < n; i++ {
		next <- i
	}
	close(next)
	wg.Wait()
}

// resized returns the faces and the shaded back at the cell size.
func resized(b *testing.B) ([]image.Image, image.Image) {
	in := load(b)
	faces := make([]image.Image, len(in.faces))
	for i, face := range in.faces {
		faces[i] = fit.Resize(face, cellW, cellH)
	}
	return faces, fit.Resize(back.Shade(in.back, true), cellW, cellH)
}

// sheet returns the painted 10000×5131 page, ready to encode.
func sheet(b *testing.B) *image.RGBA {
	faces, shaded := resized(b)
	page, _, _ := paint.Canvas(cellW, cellH, faces, shaded)
	return page
}

// The app path: decode, resize, paint, libjpeg-turbo, and the file write.

func BenchmarkAppWriteDraw(b *testing.B) {
	in := load(b)
	path := filepath.Join(b.TempDir(), "sheet.jpg")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := write.Draw(in.rawFaces, in.rawBack, cellW, cellH, true, path); err != nil {
			b.Fatal(err)
		}
	}
}

// Three pages, one after another, as compose runs them.
func BenchmarkAppWriteDrawThreeSeq(b *testing.B) {
	in := load(b)
	path := filepath.Join(b.TempDir(), "sheet.jpg")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for p := 0; p < 3; p++ {
			if err := write.Draw(in.rawFaces, in.rawBack, cellW, cellH, true, path); err != nil {
				b.Fatal(err)
			}
		}
	}
}

// Stages, one goroutine each.

func BenchmarkStageDecode(b *testing.B) {
	in := load(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, raw := range in.rawFaces {
			if _, err := images.ImageFromBinary(raw); err != nil {
				b.Fatal(err)
			}
		}
		if _, err := images.ImageFromBinary(in.rawBack); err != nil {
			b.Fatal(err)
		}
	}
}

// Full decode, as on upload, against the header read in repository GetImage.
func BenchmarkStageFormatValidate(b *testing.B) {
	in := load(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, raw := range in.rawFaces {
			if _, err := images.ValidateImage(raw); err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkStageFormatHeader(b *testing.B) {
	in := load(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, raw := range in.rawFaces {
			if _, err := images.ImageType(raw); err != nil {
				b.Fatal(err)
			}
		}
	}
}

// One goroutine. stb resizes one image per core, while imaging's Lanczos spreads
// one image over every core, so compare BenchmarkParResize for the app's real cost.
func BenchmarkStageResize(b *testing.B) {
	in := load(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, face := range in.faces {
			fit.Resize(face, cellW, cellH)
		}
		fit.Resize(in.back, cellW, cellH)
	}
}

// The previous Lanczos resize, kept until it is removed (ADR 022).
func BenchmarkStageResizeLanczos(b *testing.B) {
	in := load(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, face := range in.faces {
			fit.ResizeLanczos(face, cellW, cellH) //nolint:staticcheck // comparison with the deprecated path
		}
		fit.ResizeLanczos(in.back, cellW, cellH) //nolint:staticcheck // comparison with the deprecated path
	}
}

func BenchmarkStageShadeBack(b *testing.B) {
	in := load(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fit.Resize(back.Shade(in.back, true), cellW, cellH)
	}
}

func BenchmarkStagePaint(b *testing.B) {
	faces, shaded := resized(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		paint.Canvas(cellW, cellH, faces, shaded)
	}
}

func BenchmarkStageFileWrite(b *testing.B) {
	body, err := libjpeg.Encode(sheet(b))
	if err != nil {
		b.Fatal(err)
	}
	path := filepath.Join(b.TempDir(), "sheet.jpg")
	b.SetBytes(int64(len(body)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := fs.CreateAndProcess(path, body, fs.BinToWriter); err != nil {
			b.Fatal(err)
		}
	}
}

// Encoders on the same finished page.

func BenchmarkEncodeImageJPEG(b *testing.B) {
	img := sheet(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80}); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeLibjpeg(b *testing.B) {
	img := sheet(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := libjpeg.Encode(img); err != nil {
			b.Fatal(err)
		}
	}
}

// Stages of the parallel write.Draw, timed one at a time.

func BenchmarkParDecode(b *testing.B) {
	in := load(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parallel(len(in.rawFaces), func(j int) {
			if _, err := images.ImageFromBinary(in.rawFaces[j]); err != nil {
				b.Error(err)
			}
		})
	}
}

func BenchmarkParResize(b *testing.B) {
	in := load(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parallel(len(in.faces), func(j int) {
			fit.Resize(in.faces[j], cellW, cellH)
		})
	}
}

func BenchmarkParDecodeResize(b *testing.B) {
	in := load(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parallel(len(in.rawFaces), func(j int) {
			img, err := images.ImageFromBinary(in.rawFaces[j])
			if err != nil {
				b.Error(err)
				return
			}
			fit.Resize(img, cellW, cellH)
		})
	}
}
