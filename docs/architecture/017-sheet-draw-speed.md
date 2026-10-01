# 17. Sheet draw speed trials

* **Status:** Accepted
* **Date:** 2026-10-01
* **Authors:** @oleg

---

## Context

1. `internal/render/sheet/paint` draws one cell at a time.
2. A full page is at most 10×7.
3. The back sits in the bottom-right cell.
4. These trials live under `internal/render/sheet/draw/`.
5. The app stays on `paint`.

## How they were timed

1. One page.
2. 69 faces plus one back.
3. Each source is a foursouls.com room card, 1312×962.
4. The sheet cell is 1000×733.
5. Lanczos resamples every face and the back.
6. Two warmup runs, then 10 measured runs.
7. Each round times every package before the next round.
8. `GOMAXPROCS` was 10.
9. Minimum is the fastest run.
10. Median is the mean of the 5th and 6th sorted runs.

## Considered options

1. **`seq`**
   1. Lanczos, then `images.Draw` one cell at a time.
   2. Left to right, top to bottom, then the back.
   3. JPEG quality 80.
   4. This is the byte baseline.
   5. It can win when JPEG encode is the whole cost.
   6. Min 1599.8 ms. Median 1609.7 ms.
2. **`row`**
   1. Same pixels as `seq`.
   2. One goroutine per row writes disjoint cells.
   3. Rows can fill the canvas together.
   4. JPEG still runs after the rows join.
   5. Min 1490.9 ms. Median 1499.9 ms.
   6. JPEG matches `seq`.
3. **`cell`**
   1. Same pixels as `seq`.
   2. One goroutine per cell, including the back.
   3. More writers than `row` on a 10×7 page.
   4. JPEG still runs after the cells join.
   5. Min 1488.8 ms. Median 1699.3 ms.
   6. JPEG matches `seq`.
4. **`resize`**
   1. Lanczos on every face and the back, in parallel.
   2. Then the `seq` draw.
   3. It helps when faces are not already the cell size.
   4. Every face is larger than the cell, so Lanczos runs.
   5. Min 1564.6 ms. Median 1582.0 ms.
   6. JPEG matches `seq`.
5. **`resize_row`**
   1. Parallel Lanczos, then one goroutine per row.
   2. It stacks the `resize` and `row` ideas.
   3. Parallel Lanczos resamples, then the rows draw.
   4. Min 1458.9 ms. Median 1474.4 ms.
   5. JPEG matches `seq`.
6. **`bilinear`**
   1. `ApproxBiLinear` instead of Lanczos.
   2. Then the `seq` draw.
   3. The filter is cheaper when a face is resampled.
   4. Pixels may differ after a resample.
   5. These faces are resampled, so this JPEG does not match `seq`.
   6. Min 1888.0 ms. Median 1932.8 ms.
7. **`pages`**
   1. Each page is `seq`, including its JPEG.
   2. Pages run at the same time.
   3. A multi-page deck can overlap encode.
   4. This trial was one page, so nothing overlapped.
   5. Min 1773.8 ms. Median 1808.0 ms.
   6. That one JPEG matches `seq`.

## Decision

1. Do not switch the app off `internal/render/sheet/paint`.
2. Fastest median is `resize_row` (1474.4 ms).
3. `row`, `cell`, `resize`, and `resize_row` match the `seq` JPEG bytes.

## Consequences

### Positive

1. The seven packages can be timed again on the same inputs.
2. `resize_row` is 135.3 ms under `seq` on this page.

### Negative and risks

1. `cell` and `row` write one `image.RGBA` from several goroutines.
2. The writes are disjoint cells.
3. JPEG encode is still one thread, and it is most of the time.
4. `bilinear` bytes change when a face is actually resampled.

### Neutral

1. `internal/render/page_drawer` is unchanged.
2. `internal/render/sheet/paint` is unchanged.

## Stage split

1. Same images, `GOMAXPROCS` 10, two warmups, then 10 runs.
2. Decode sits outside the timed call.
3. Four PNG decodes: min 145.1 ms, median 146.1 ms.
4. `png.Decode` returns `*image.NRGBA` or `*image.RGBA`.
5. Lanczos returns `*image.NRGBA`.
6. The sheet is `*image.RGBA`.
7. `draw.Src` from `*image.NRGBA` uses `drawNRGBASrc`.
8. That loop premultiplies each pixel.
9. It does not call `At` or `Set`.
10. `*image.RGBA` with `draw.Src` uses `copy` per row.
11. `internal/render/sheet/draw/rgba` converts each cell, then draws.
12. That JPEG matches `seq`.
13. `rgba` min 1752.2 ms. Median 1779.5 ms.
14. JPEG encode is the shared cost.

### `seq`

1. Lanczos min 445.0 ms. Median 460.2 ms.
2. Draw min 149.6 ms. Median 152.6 ms.
3. JPEG min 1134.4 ms. Median 1136.7 ms.

### `resize_row`

1. Lanczos min 406.0 ms. Median 419.0 ms.
2. Draw min 45.7 ms. Median 49.2 ms.
3. JPEG min 1146.1 ms. Median 1150.0 ms.

## JPEG encoders

1. Same 10000×5131 sheet. Quality 80.
2. `GOMAXPROCS` 10. Two warmups, then 10 runs.
3. Decode, Lanczos, and draw sit outside this timer.
4. `image/jpeg` min 1128.5 ms. Median 1136.8 ms.
5. That row is the byte baseline.
6. `libjpeg` uses libjpeg-turbo via `github.com/pixiv/go-libjpeg`.
7. It needs cgo and Homebrew `jpeg-turbo`.
8. `libjpeg` min 157.6 ms. Median 161.6 ms.
9. Those bytes do not match `image/jpeg`.
10. `jpegli` uses `github.com/gen2brain/jpegli`.
11. That build runs jpegli under wazero.
12. It needs no cgo and no system library.
13. `jpegli` min 882.3 ms. Median 908.9 ms.
14. Those bytes do not match `image/jpeg`.
15. libvips was left out. The first two answered the question.
16. Fastest encoder is `libjpeg` (161.6 ms).
17. The app stays on `image/jpeg`.
