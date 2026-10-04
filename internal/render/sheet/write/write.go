// Draw one measured page and write a libjpeg-turbo quality-80 JPEG.
package write

import (
	"errors"
	"fmt"
	"image"
	"runtime"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/HardDie/DeckBuilder/internal/fs"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/back"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/draw/libjpeg"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/fit"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/paint"
)

// ErrUndecodable marks a face or back that does not decode; the decoder's text follows it.
var ErrUndecodable = errors.New("image does not decode")

// Timings is how long each part of one Draw took.
// Jobs is the wall time of the parallel decode and resize.
// Decode and Resize add up every job, so they can exceed Jobs.
type Timings struct {
	Jobs   time.Duration
	Decode time.Duration
	Resize time.Duration
	Paint  time.Duration
	Encode time.Duration
	Write  time.Duration
}

// Draw paints faces and the back into path and reports how long each part took.
// Cell width and height are already chosen. Faces and back are the original file bytes.
// Faces and the back are decoded and resized in parallel, at most GOMAXPROCS at a time.
func Draw(faces [][]byte, rawBack []byte, cellW, cellH int, shadow bool, path string) (Timings, error) {
	var t Timings
	var decodeSum, resizeSum atomic.Int64
	// decodeResize decodes raw, then resizes it with prepare, adding to both sums.
	decodeResize := func(raw []byte, prepare func(image.Image) image.Image) (image.Image, error) {
		start := time.Now()
		img, err := images.ImageFromBinary(raw)
		decoded := time.Now()
		decodeSum.Add(int64(decoded.Sub(start)))
		if err != nil {
			return nil, err
		}
		out := prepare(img)
		resizeSum.Add(int64(time.Since(decoded)))
		return out, nil
	}
	drawn := make([]image.Image, len(faces))
	var shaded image.Image
	jobs := make([]func() error, 0, len(faces)+1)
	for i, raw := range faces {
		jobs = append(jobs, func() error {
			img, err := decodeResize(raw, func(img image.Image) image.Image {
				return fit.Resize(img, cellW, cellH)
			})
			if err != nil {
				return fmt.Errorf("%w: face %d: %v", ErrUndecodable, i+1, err)
			}
			drawn[i] = img
			return nil
		})
	}
	jobs = append(jobs, func() error {
		img, err := decodeResize(rawBack, func(img image.Image) image.Image {
			return fit.Resize(back.Shade(img, shadow), cellW, cellH)
		})
		if err != nil {
			return fmt.Errorf("%w: back: %v", ErrUndecodable, err)
		}
		shaded = img
		return nil
	})
	start := time.Now()
	err := run(jobs)
	t.Jobs = time.Since(start)
	t.Decode = time.Duration(decodeSum.Load())
	t.Resize = time.Duration(resizeSum.Load())
	if err != nil {
		return t, err
	}
	start = time.Now()
	sheet, _, _ := paint.Canvas(cellW, cellH, drawn, shaded)
	t.Paint = time.Since(start)
	start = time.Now()
	body, err := libjpeg.Encode(sheet)
	t.Encode = time.Since(start)
	if err != nil {
		return t, err
	}
	start = time.Now()
	err = fs.CreateAndProcess(path, body, fs.BinToWriter)
	t.Write = time.Since(start)
	return t, err
}

// run calls every job on at most GOMAXPROCS goroutines and returns the first error in job order.
func run(jobs []func() error) error {
	errs := make([]error, len(jobs))
	next := make(chan int)
	var wg sync.WaitGroup
	workers := min(runtime.GOMAXPROCS(0), len(jobs))
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for i := range next {
				errs[i] = callJob(jobs[i])
			}
		}()
	}
	for i := range jobs {
		next <- i
	}
	close(next)
	wg.Wait()
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// callJob runs job and turns a panic into its error.
// A panic in a worker goroutine would otherwise end the whole app.
func callJob(job func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v\n%s", r, debug.Stack())
		}
	}()
	return job()
}
