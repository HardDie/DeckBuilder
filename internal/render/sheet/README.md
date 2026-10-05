# `internal/render/sheet`

Draws one TTS custom-deck page and writes it as a JPEG. Given the raw face images, the raw back image, a cell size, and a path, it decodes, resizes, paints, encodes, and writes.

## Where this fits

`render/generate` already decided what goes on the page and how big a cell is. `render/compose` calls `write.Draw` once per page that is not reused. This package never imports `render/generate`; see [compose](../compose/README.md) for the whole layer.

## The logic in short

```text
face bytes ─► decode ─► resize to cell ─┐
face bytes ─► decode ─► resize to cell ─┤   in parallel,
   …                                    ├─► paint ─► JPEG ─► file
back bytes ─► decode ─► darken ─► resize┘   up to one job per core
```

1. **Decode** every face and the back (PNG, JPEG, GIF, WebP, BMP, or TIFF).
2. **Resize** each to the cell with `stb_image_resize2` (Catmull-Rom).
3. **Darken** the back when the shadow setting is on.
4. **Paint** faces left to right, top to bottom; the back goes in the last cell.
5. **Encode** with libjpeg-turbo at quality 80 and write the file.

Steps 1–3 run in parallel. Painting and encoding run on one goroutine: in the trials, painting in parallel saved about 3%, and drawing several pages at once saved nothing and doubled memory ([ADR 017](../../../docs/architecture/017-sheet-draw-speed.md)).

## `write.Draw`

```go
func Draw(faces [][]byte, rawBack []byte, cellW, cellH int, shadow bool, path string) (Timings, error)
```

- `faces`, `rawBack`: the original file bytes, in page order.
- `cellW`, `cellH`: already chosen by `generate/layout`.
- One job per face plus one for the back runs on at most `GOMAXPROCS` goroutines. The first error in job order is returned.
- A panic in a job is recovered and returned as that job's error, so a bad image fails the page instead of the app.
- The canvas is `cellW × columns` by `cellH × rows`; the grid comes from `grid.Size(faces + 1)`.
- `compose` passes a temporary path and renames it afterwards; `Draw` itself writes the path it is given.
- `Timings` reports each part: the wall time of the parallel jobs, the summed decode and resize time of all jobs, then paint, encode, and write. `compose` logs it at Debug.

## Resize

`fit.Resize(img, w, h)` returns `img` unchanged when it already has the cell size. Otherwise it calls `stbresize.Resize`:

- **Filter:** Catmull-Rom. Visually identical to the Lanczos it replaced (PSNR 50 dB). The resize step is about 5× faster; a full page about 44% faster.
- **Any decoded layout** (NRGBA, RGBA, RGBA64, Gray, YCbCr, Paletted) is converted to 8-bit NRGBA first.
- **Opaque images** filter the four channels alike: the fast path.
- **Images with transparency** weight color by alpha, so the hidden color of transparent pixels does not bleed into edges as a dark fringe.
- **Static linking:** `stb_image_resize2.h` (public domain / MIT) is vendored and compiled into the binary by cgo. There is no library to link. SIMD comes from the target: NEON on arm64, SSE2 on amd64.

The previous pure-Go Lanczos path was removed once release builds passed on all four targets. Details: [ADR 022](../../../docs/architecture/022-stb-image-resize.md).

**When you change how a page looks** (filter, paint, shadow, encoder, quality), bump `sheetVersion` in `generate/layout`. Otherwise `compose` keeps reusing pages drawn by the old code.

## Packages

| Package | Job |
|---|---|
| `write` | `Draw`: the whole page, as above. The only entry point. |
| `fit` | `Resize` to the cell (stb). |
| `third_party/stbresize` | The cgo wrapper around `stb_image_resize2.h` (moved out of `sheet/`, ADR 028). |
| `back` | `Shade`: `imaging.AdjustBrightness` by −30% when shadow is on. With shadow off it still converts to NRGBA. |
| `grid` | `Size(n)`: smallest columns × rows from 2×2 to 10×7 that holds `n` cells. `Slot(i, cols)`: left to right, top to bottom. Same rule as `generate/layout`; keep them identical. |
| `paint` | `Canvas`: places the already-resized images on an RGBA page; the back in the bottom-right cell. |
| `draw/libjpeg` | `Encode`: libjpeg-turbo, quality 80, `DCTISlow`. Accepts `*image.RGBA`, `*image.YCbCr`, or `*image.Gray`. Linked as a static archive. |
| `bench` | Benchmarks of the live path (tests only). |

## Benchmarks

`bench` measures one full page: 69 faces and the back, 1312×962 PNGs from `bench/testdata/input`, resized to a 1000×733 cell.

| Benchmark | Measures |
|---|---|
| `AppWriteDraw`, `AppWriteDrawThreeSeq` | `write.Draw` end to end, for one page and three pages in a row |
| `Stage*` | One step at a time on one goroutine: decode, format check, resize, shade, paint, file write |
| `Par*` | Decode and resize spread over the cores, as `write.Draw` runs them |
| `Encode*` | `image/jpeg` against libjpeg-turbo on the same painted page |

`StageResize` runs one image at a time on one core; `ParResize` is the real cost in the app. They need the cgo and libjpeg environment from CLAUDE.md ("Commands"):

```bash
go test -tags=nomain -run '^$' -bench . -benchtime 5x ./internal/render/sheet/bench/
```
