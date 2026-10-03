# 22. Resize faces with stb_image_resize2

* **Status:** Accepted
* **Date:** 2026-10-03
* **Authors:** @oleg

---

## Context

1. Every face and the back are resized to the sheet cell.
2. The resize was `imaging.Resize` with Lanczos.
   1. Pure Go, float64, no SIMD.
   2. It already spreads one image over every core.
3. After [ADR 017](017-sheet-draw-speed.md), resize was the largest cost.
   1. Full 10×7 page, 69 faces of 1312×962 → 1000×733, Apple M4.
   2. Resize ~735 ms of a ~1535 ms page.
4. The app must stay statically linked (CLAUDE.md, Targets).
5. Inputs are PNG, JPEG, or GIF, decoded by Go.
   1. Decoded layouts: NRGBA, RGBA, RGBA64, Gray, YCbCr, Paletted.

## Considered options

Measured on the page above. Quality is PSNR against the Lanczos output.

| Option | Resize | Quality | Notes |
|---|---|---|---|
| `imaging` Lanczos | 757–840 ms | — | the previous path |
| `imaging` Catmull-Rom | ~500 ms | 50.2 dB | one-line change |
| `x/image/draw` Catmull-Rom | 616 ms | — | single-threaded per image |
| **stb Catmull-Rom, 4-channel** | **148 ms** | **50.2 dB** | opaque images |
| stb Catmull-Rom, alpha-aware | 360 ms | 50.2 dB | images with transparency |
| stb Mitchell, 4-channel | 162 ms | 43.3 dB | softer |
| libvips | — | — | rejected: large shared library |

1. Above about 40 dB the difference is not visible; 50 dB is rounding noise.
2. stb has no Lanczos; Catmull-Rom is the closest sharp filter.

## Decision

1. `fit.Resize` calls `sheet/stbresize.Resize`.
   1. Filter: Catmull-Rom.
   2. Any decoded layout is converted to 8-bit NRGBA first.
2. Alpha
   1. Opaque image: 4-channel path, about 2.4× faster.
   2. Image with transparency: alpha-aware path.
   3. Without it, hidden colors bleed into edges (31.6 dB, dark fringe).
3. Static linking
   1. `stb_image_resize2.h` v2.18 is vendored in `sheet/stbresize`.
   2. cgo compiles it into the binary; there is no library to link.
   3. `STB_IMAGE_RESIZE_STATIC` keeps its symbols internal.
   4. SIMD comes from the target: NEON on arm64, SSE2 on amd64.
   5. AVX2 is not enabled; it would need a flag older CPUs lack.
4. Fallback
   1. `fit.ResizeLanczos` keeps the previous path, marked deprecated.
   2. Remove it once release builds pass on all four targets.

## Consequences

### Positive

1. A full page dropped from ~1535 ms to ~854 ms (−44%) on an M4.
2. Output is visually identical (≥ 50 dB on every decoded layout).
3. No new runtime dependency; the binary links the same system libraries.

### Negative and risks

1. 10.7k lines of vendored C to keep.
   1. Upgrades mean replacing one header.
2. Windows (MinGW) and Linux (gcc) builds were not run locally.
   1. macOS arm64 was built; the amd64 SSE2 path was tested under Rosetta.
   2. CI tests run on ubuntu amd64; the release workflow builds all four.
3. Float rounding can differ between NEON and SSE2.
   1. Resize tests use PSNR, not exact bytes.
   2. The compose goldens use cell-sized faces, so they never resize.

### Neutral

1. `disintegration/imaging` stays for `back.Shade` and the deprecated fallback.
2. `sheet/bench` keeps `StageResizeLanczos` for comparison until removal.

## Later

1. 2026-10-03: release builds passed on all four targets.
2. `fit.ResizeLanczos` and its benchmark were removed (review item S15).
3. A test still compares stb with `imaging` Lanczos (≥ 40 dB).
