// Draw one measured page and write a libjpeg-turbo quality-80 JPEG.
package write

import (
	"errors"
	"fmt"
	"image"
	"runtime"
	"runtime/debug"
	"sync"

	"github.com/HardDie/DeckBuilder/internal/fs"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/back"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/draw/libjpeg"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/fit"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/paint"
)

// ErrUndecodable marks a face or back that does not decode; the decoder's text follows it.
var ErrUndecodable = errors.New("image does not decode")

// Draw paints faces and the back into path.
// Cell width and height are already chosen. Faces and back are the original file bytes.
// Faces and the back are decoded and resized in parallel, at most GOMAXPROCS at a time.
func Draw(faces [][]byte, rawBack []byte, cellW, cellH int, shadow bool, path string) error {
	drawn := make([]image.Image, len(faces))
	var shaded image.Image
	jobs := make([]func() error, 0, len(faces)+1)
	for i, raw := range faces {
		jobs = append(jobs, func() error {
			img, err := images.ImageFromBinary(raw)
			if err != nil {
				return fmt.Errorf("%w: face %d: %v", ErrUndecodable, i+1, err)
			}
			drawn[i] = fit.Resize(img, cellW, cellH)
			return nil
		})
	}
	jobs = append(jobs, func() error {
		decoded, err := images.ImageFromBinary(rawBack)
		if err != nil {
			return fmt.Errorf("%w: back: %v", ErrUndecodable, err)
		}
		shaded = fit.Resize(back.Shade(decoded, shadow), cellW, cellH)
		return nil
	})
	if err := run(jobs); err != nil {
		return err
	}
	sheet, _, _ := paint.Canvas(cellW, cellH, drawn, shaded)
	body, err := libjpeg.Encode(sheet)
	if err != nil {
		return err
	}
	return fs.CreateAndProcess(path, body, fs.BinToWriter)
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
