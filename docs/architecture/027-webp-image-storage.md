# 27. Store every image as WebP

* **Status:** Accepted
* **Date:** 2026-10-05
* **Authors:** @oleg

---

## Context

1. Images are stored as uploaded: PNG, JPEG, or GIF.
2. Only those three formats can be uploaded.
3. Generation must decode every stored format.
   1. Go decodes JPEG and PNG slowly.
   2. Each decoder returns a different pixel layout; each needs its own conversion.
4. A full render of a 1025-card game took about 65 s.
   1. Decode and resize were about 80% of it.
5. Measured on 75 real images (single thread, libwebp):

| Source | WebP q95 size | WebP q95 decode | Original decode (Go) |
|---|---|---|---|
| Progressive JPEG 1987×2708 | 1.23× | 42 ms | 176 ms |
| Baseline JPEG 1987×2708 | 1.50× | 50 ms | 47 ms |
| PNG RGB 962×1312 | 0.17× | 10 ms | 27 ms |

6. WebP q95 loses 42–47 dB PSNR against the decoded original.
   1. Side by side at 3× zoom, text and edges look the same.
   2. Sheets are JPEG q80, which loses more.
7. WebP lossless was 3.5–4× larger for JPEG sources and not faster.
8. Lossy WebP always stores color at half resolution (4:2:0).
   1. Sharp colored detail drops to 17–21 dB: pixel art, thin colored lines.
   2. libwebp's sharp YUV option did not help, and doubled encode time.
   3. Lossless WebP of such images was smaller than lossy (166 B and 2.5 KB vs 22 KB and 19 KB).

## Considered options

1. **Keep every format; one fast decoder per format**
   1. No conversion, no migration.
   2. Lost. Generation keeps one path per format; each new format adds one.
2. **JPEG as is, everything else lossless**
   1. No size growth for JPEG sources.
   2. Lost. Two formats to decode, and a lossless format still needs a choice.
3. **WebP lossless for everything**
   1. Lost. 3.5–4× larger for JPEG sources, no faster.
4. **WebP lossy q95 for everything**
   1. One stored format, one decoder, alpha kept, shown by every webview.
   2. Won.

## Decision

1. Storage
   1. Every image in the catalog is WebP lossy, quality 95.
   2. Game, collection, deck, and card images alike.
   3. File names do not change (`image.bin`, `cards/<id>.bin`).
   4. Alpha is kept.
   5. Lossy unless lossy keeps less than 40 dB; then lossless.
      1. Real card art measured 45–48 dB, so it stays lossy.
      2. Lossless must decode exactly.
      3. libwebp decodes both with one call.
2. Input formats
   1. PNG, JPEG, GIF, WebP, BMP, TIFF.
   2. All six decode in pure Go (standard library and `golang.org/x/image`).
   3. GIF uses its first frame.
   4. HEIC, AVIF, and JPEG XL need C libraries; not now.
   5. ADR 020 limits apply before decoding.
   6. A side over 16383 pixels is rejected: WebP cannot hold it.
3. Conversion
   1. One function turns uploaded bytes into WebP.
   2. Upload, URL download, preview, and migration all call it.
   3. JPEG EXIF orientation is applied, so photos stand upright.
   4. A WebP upload is converted too, so every file has the same settings.
4. Encoder
   1. libwebp, linked as a static archive (`libwebp.a`, `libsharpyuv.a`).
   2. Pure Go has no WebP encoder.
   3. Release CI builds 1.6.0 from source on every target (`scripts/build-libwebp.sh`).
   4. CI checks it is not a shared dependency, as for libjpeg.
   5. The build uses libwebp's `makefile.unix`: only make and a C compiler.
   6. Threads are off, so no pthread (winpthread on Windows) is linked.
   7. macOS builds an arm64 and an x86_64 slice and joins them with `lipo`.
5. Display
   1. The image route already sends `Content-Type` from the file header.
   2. It sends `image/webp`; all three webviews show it.
   3. WKWebView shows WebP from macOS 11, so the minimum macOS becomes 12.0 (see 5.7).
   4. The upload preview no longer shows the local file.
   5. A binding converts the chosen file and returns the WebP to preview.
   6. A file that cannot be converted fails at preview, before save.
   7. 12.0, not 11.0: Go 1.25+ already needs macOS 12 (`LSMinimumSystemVersion`).
6. Game marker
   1. A converted game has `image_format: "webp"` in its `.info.json` data.
   2. A game created by this version gets the marker at once.
   3. Startup reads only the markers, never the images.
7. Migration of existing games
   1. At startup, games without the marker are converted.
   2. A modal names those games and their image counts first.
   3. The whole app is blocked until it ends.
   4. A progress bar shows the game and n of total images.
   5. Images are converted in parallel.
8. Safety of each image
   1. An image that is already WebP is skipped.
   2. Decode, apply orientation, encode WebP q95.
   3. Verify: decode the WebP; same width and height; lossy at least 40 dB, lossless exact.
   4. Copy the original to `DeckBuilderData/backup/webp/<game>/…`.
   5. Replace the file atomically: temporary file, fsync, rename.
   6. fsentry gains this atomic replace; `UpdateBinary` rewrote in place.
   7. A crash leaves each image either old or new, never broken.
9. Resume
   1. Closing the app mid-migration is safe.
   2. The next start finds the same games without the marker.
   3. Converted images are skipped, so it continues where it stopped.
10. Disk space
    1. Before a game, free space must hold its backup and its new files.
    2. If not, that game is not started and the modal says how much is needed.
11. Failures
    1. An image that does not convert keeps its original file.
    2. The game gets no marker.
    3. The modal lists the failed images.
    4. The next start tries that game again.
12. Backups
    1. They are kept.
    2. After the migration the modal says where they are and how large.
    3. The user deletes them by hand.
13. Generation
    1. A game without the marker cannot be generated.
    2. If every image is WebP already, the marker is written and the render starts.
       1. This covers a game fixed by uploading new images.
    3. Otherwise the render does not start.
       1. The message names up to 5 images that are not WebP, like missing images.
    4. Generation decodes only WebP, with libwebp.
       1. The other decoders are not on the render path.
    5. Pages are re-rendered once after migration: their bytes changed.
14. Archives
    1. Import writes the files as they are and removes the marker.
    2. The migration then converts the game, with its progress and safety.
    3. A WebP archive imports quickly: its images are skipped.
    4. An archive exported by this version does not import into an older one.

## Consequences

### Positive

1. Generation has one decoder and one pixel layout.
2. Six input formats instead of three.
3. Faster decode for most cards (42 ms instead of 176 ms for progressive JPEG).
4. Rotated phone photos stand upright.
5. Image writes become atomic for every caller.

### Negative and risks

1. A second lossy step for JPEG sources (42–47 dB).
2. JPEG-sourced games grow on disk by about 1.2–1.5×.
3. Backups double disk use until the user deletes them.
4. libwebp is a new static dependency on four targets.
5. macOS 10.13–11 are no longer listed as supported; Go 1.25+ never ran there.
6. A large collection blocks the app at its first start (about 1 minute per 1000 large cards on 10 cores).
7. Applying EXIF orientation can turn existing cards that were sideways on sheets.
8. An image that falls back to lossless is encoded twice.
   1. About 0.7 s per large card when lossy is kept, single thread.

### Neutral

1. File names, ids, and JSON references do not change.
2. The page hash covers image bytes, so reuse works as before after one re-render.

## Implementation steps

1. This ADR.
2. fsentry: atomic binary replace.
3. `internal/images`: decode the six formats, orientation, WebP encode, verify.
4. libwebp in the Makefile and release CI; minimum macOS 12.0.
5. Upload, URL download, and preview through the conversion.
6. Game marker; new games get it.
7. Migration service, backups, disk check, resume.
8. Migration modal and progress in the window.
9. Generation gate and WebP-only decoding.
10. Archive import removes the marker.
