// Benchmarks for one full sheet: 69 faces plus the back on a 10×7 page.
// Inputs are the 1312×962 PNGs in ../page/testdata/input. The cell is 1000×733.
//
//	go test -tags=nomain -run '^$' -bench . -benchtime 5x ./internal/render/sheet/bench/
package bench

import (
	"bytes"
	"image"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/HardDie/DeckBuilder/internal/fs"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/back"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/draw/bilinear"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/draw/cell"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/draw/jpegli"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/draw/libjpeg"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/draw/pages"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/draw/resize"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/draw/resize_row"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/draw/rgba"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/draw/row"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/draw/seq"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/fit"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/grid"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/paint"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/write"
)

const faceCount = 69

type input struct {
	faces    []image.Image
	back     image.Image
	rawFaces [][]byte
	rawBack  []byte
	cellW    int
	cellH    int
}

var (
	loadOnce sync.Once
	loaded   input
)

func load(b *testing.B) input {
	b.Helper()
	loadOnce.Do(func() {
		dir := filepath.Join("..", "page", "testdata", "input")
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
		max := dec[0].Bounds().Max
		_, size := fit.Apply(max.X, max.Y, 1, fit.Size{})
		loaded.cellW, loaded.cellH = size.Width, size.Height
	})
	return loaded
}

// Draw trials: Lanczos, draw, then image/jpeg. Decode is outside the timer.

func benchJPEG(b *testing.B, fn func([]image.Image, image.Image) ([]byte, error)) {
	in := load(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := fn(in.faces, in.back); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDrawSeq(b *testing.B)       { benchJPEG(b, seq.JPEG) }
func BenchmarkDrawRow(b *testing.B)       { benchJPEG(b, row.JPEG) }
func BenchmarkDrawCell(b *testing.B)      { benchJPEG(b, cell.JPEG) }
func BenchmarkDrawResize(b *testing.B)    { benchJPEG(b, resize.JPEG) }
func BenchmarkDrawResizeRow(b *testing.B) { benchJPEG(b, resize_row.JPEG) }
func BenchmarkDrawBilinear(b *testing.B)  { benchJPEG(b, bilinear.JPEG) }
func BenchmarkDrawRGBA(b *testing.B)      { benchJPEG(b, rgba.JPEG) }

// Pages: several full pages encoded at the same time.
func BenchmarkDrawPagesThree(b *testing.B) {
	in := load(b)
	sheets := []pages.Sheet{{Faces: in.faces, Back: in.back}, {Faces: in.faces, Back: in.back}, {Faces: in.faces, Back: in.back}}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := pages.JPEGs(sheets); err != nil {
			b.Fatal(err)
		}
	}
}

// Same three pages, one after another, for comparison with BenchmarkDrawPagesThree.
func BenchmarkDrawSeqThree(b *testing.B) {
	in := load(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for p := 0; p < 3; p++ {
			if _, err := seq.JPEG(in.faces, in.back); err != nil {
				b.Fatal(err)
			}
		}
	}
}

// Stages of the seq path, timed one at a time.

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

// Format check in repository GetImage, once per face during generate.Prepare.
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

func BenchmarkStageLanczos(b *testing.B) {
	in := load(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, face := range in.faces {
			fit.Resize(face, in.cellW, in.cellH)
		}
		fit.Resize(in.back, in.cellW, in.cellH)
	}
}

func BenchmarkStageDraw(b *testing.B) {
	in := load(b)
	resized := make([]image.Image, len(in.faces))
	for i, face := range in.faces {
		resized[i] = fit.Resize(face, in.cellW, in.cellH)
	}
	back := fit.Resize(in.back, in.cellW, in.cellH)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		seq.Image(resized, back)
	}
}

// Encoders on the same finished 10000×5131 sheet.

func sheet(b *testing.B) *image.RGBA {
	in := load(b)
	return seq.Image(in.faces, in.back)
}

func BenchmarkEncodeImageJPEG(b *testing.B) {
	img := sheet(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		if err := images.JpegSaveToWriter(&buf, img); err != nil {
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

func BenchmarkEncodeJpegli(b *testing.B) {
	img := sheet(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := jpegli.Encode(img); err != nil {
			b.Fatal(err)
		}
	}
}

// The app path: decode, Lanczos, paint, libjpeg-turbo, and the file write.
func BenchmarkAppWriteDraw(b *testing.B) {
	in := load(b)
	path := filepath.Join(b.TempDir(), "sheet.jpg")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := write.Draw(in.rawFaces, in.rawBack, in.cellW, in.cellH, true, path); err != nil {
			b.Fatal(err)
		}
	}
}

// Three app pages, one after another, as compose runs them today.
func BenchmarkAppWriteDrawThreeSeq(b *testing.B) {
	in := load(b)
	dir := b.TempDir()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for p := 0; p < 3; p++ {
			path := filepath.Join(dir, "sheet.jpg")
			if err := write.Draw(in.rawFaces, in.rawBack, in.cellW, in.cellH, true, path); err != nil {
				b.Fatal(err)
			}
		}
	}
}

// Stages of the parallel write.Draw, timed one at a time.

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

func BenchmarkParLanczos(b *testing.B) {
	in := load(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parallel(len(in.faces), func(j int) {
			fit.Resize(in.faces[j], in.cellW, in.cellH)
		})
	}
}

func BenchmarkParDecodeLanczos(b *testing.B) {
	in := load(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parallel(len(in.rawFaces), func(j int) {
			img, err := images.ImageFromBinary(in.rawFaces[j])
			if err != nil {
				b.Error(err)
				return
			}
			fit.Resize(img, in.cellW, in.cellH)
		})
	}
}

func BenchmarkStageShadeBack(b *testing.B) {
	in := load(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fit.Resize(back.Shade(in.back, true), in.cellW, in.cellH)
	}
}

func BenchmarkStagePaint(b *testing.B) {
	in := load(b)
	resized := make([]image.Image, len(in.faces))
	for i, face := range in.faces {
		resized[i] = fit.Resize(face, in.cellW, in.cellH)
	}
	shaded := fit.Resize(back.Shade(in.back, true), in.cellW, in.cellH)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		paint.Canvas(in.cellW, in.cellH, resized, shaded)
	}
}

func BenchmarkStageFileWrite(b *testing.B) {
	in := load(b)
	body, err := libjpeg.Encode(seq.Image(in.faces, in.back))
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

// The write.Draw before cells were painted inside the workers:
// parallel decode and Lanczos, then paint.Canvas on one goroutine, then libjpeg-turbo.
func previousWriteDraw(faces [][]byte, rawBack []byte, cellW, cellH int, shadow bool, path string) error {
	drawn := make([]image.Image, len(faces))
	var shaded image.Image
	errs := make([]error, len(faces)+1)
	parallel(len(faces)+1, func(i int) {
		if i == len(faces) {
			decoded, err := images.ImageFromBinary(rawBack)
			if err != nil {
				errs[i] = err
				return
			}
			shaded = fit.Resize(back.Shade(decoded, shadow), cellW, cellH)
			return
		}
		img, err := images.ImageFromBinary(faces[i])
		if err != nil {
			errs[i] = err
			return
		}
		drawn[i] = fit.Resize(img, cellW, cellH)
	})
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	sheet, _, _ := paint.Canvas(cellW, cellH, drawn, shaded)
	body, err := libjpeg.Encode(sheet)
	if err != nil {
		return err
	}
	return fs.CreateAndProcess(path, body, fs.BinToWriter)
}

func BenchmarkAppWriteDrawPrevious(b *testing.B) {
	in := load(b)
	path := filepath.Join(b.TempDir(), "sheet.jpg")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := previousWriteDraw(in.rawFaces, in.rawBack, in.cellW, in.cellH, true, path); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAppWriteDrawThreeSeqPrevious(b *testing.B) {
	in := load(b)
	dir := b.TempDir()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for p := 0; p < 3; p++ {
			if err := previousWriteDraw(in.rawFaces, in.rawBack, in.cellW, in.cellH, true, filepath.Join(dir, "sheet.jpg")); err != nil {
				b.Fatal(err)
			}
		}
	}
}

// Every draw method on equal terms: decoded inputs, then the method's Image, then libjpeg-turbo.
// Decode is outside the timer. The app encoder replaces image/jpeg, which hid the draw cost.

// appImage is option A from decoded images: parallel Lanczos, each worker draws its cell.
// write.Draw no longer does this. It stays here as one more method to compare.
func appImage(faces []image.Image, backImg image.Image) *image.RGBA {
	b := faces[0].Bounds().Max
	_, size := fit.Apply(b.X, b.Y, 1, fit.Size{})
	cols, rows := grid.Size(len(faces) + 1)
	page := images.CreateImage(size.Width*cols, size.Height*rows)
	parallel(len(faces)+1, func(i int) {
		if i == len(faces) {
			images.Draw(page, cols-1, rows-1, fit.Resize(backImg, size.Width, size.Height))
			return
		}
		col, row := grid.Slot(i, cols)
		images.Draw(page, col, row, fit.Resize(faces[i], size.Width, size.Height))
	})
	return page
}

func benchMethod(b *testing.B, draw func([]image.Image, image.Image) *image.RGBA) {
	in := load(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := libjpeg.Encode(draw(in.faces, in.back)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMethodSeq(b *testing.B)       { benchMethod(b, seq.Image) }
func BenchmarkMethodRow(b *testing.B)       { benchMethod(b, row.Image) }
func BenchmarkMethodCell(b *testing.B)      { benchMethod(b, cell.Image) }
func BenchmarkMethodResize(b *testing.B)    { benchMethod(b, resize.Image) }
func BenchmarkMethodResizeRow(b *testing.B) { benchMethod(b, resize_row.Image) }
func BenchmarkMethodBilinear(b *testing.B)  { benchMethod(b, bilinear.Image) }
func BenchmarkMethodRGBA(b *testing.B)      { benchMethod(b, rgba.Image) }
func BenchmarkMethodApp(b *testing.B)       { benchMethod(b, appImage) }

// Three pages: one after another, or all at once like the pages trial.
func benchPages(b *testing.B, draw func([]image.Image, image.Image) *image.RGBA, together bool) {
	in := load(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		one := func() {
			if _, err := libjpeg.Encode(draw(in.faces, in.back)); err != nil {
				b.Error(err)
			}
		}
		if !together {
			for p := 0; p < 3; p++ {
				one()
			}
			continue
		}
		var wg sync.WaitGroup
		wg.Add(3)
		for p := 0; p < 3; p++ {
			go func() { defer wg.Done(); one() }()
		}
		wg.Wait()
	}
}

func BenchmarkPagesSeqOneByOne(b *testing.B) { benchPages(b, seq.Image, false) }
func BenchmarkPagesSeqTogether(b *testing.B) { benchPages(b, seq.Image, true) }
func BenchmarkPagesAppOneByOne(b *testing.B) { benchPages(b, appImage, false) }
func BenchmarkPagesAppTogether(b *testing.B) { benchPages(b, appImage, true) }
