# 20. Limits on incoming images

* **Status:** Accepted
* **Date:** 2026-10-03
* **Authors:** @oleg

---

## Context

1. Every entity can take an image from a URL or from uploaded bytes.
2. `network.DownloadBytes` fetches URL images.
   1. Game, collection, and deck reach it through `repositories.ImageBytes`.
   2. The card repository calls it directly.
3. The download runs inside the create / update binding.
   1. The dialog waits until it returns.
4. `http.Get` used the default client.
   1. The default client has no timeout.
   2. A stalled server kept the binding waiting forever.
5. `io.ReadAll` read the whole body.
   1. A wrong URL could load a huge file into memory.
6. `images.ValidateImage` fully decodes every new image.
   1. A tiny PNG can declare 50000×50000 pixels.
   2. Decoding it allocates about 10 GB.
   3. That ends the process.
   4. Uploaded bytes take the same path as downloads.

## Considered options

1. **No limits**
   1. The previous behavior.
   2. Lost. One bad URL or file can hang or kill the app.
2. **Timeout and size cap on the download only**
   1. Stops hangs and huge bodies.
   2. Lost on its own. A decompression bomb is a small file.
   3. Uploads skip the download path entirely.
3. **Download limits plus a pixel cap before decode**
   1. Download: timeout and byte cap.
   2. Validate: read the header with `image.DecodeConfig` first.
   3. Reject before any pixel memory is allocated.
   4. Covers URLs and uploads in one place.
   5. Won.

## Decision

Use option 3.

1. Download timeout
   1. 120 seconds.
   2. Set as `http.Client.Timeout`.
   3. It covers connect, headers, and body.
   4. Long enough for big images on slow links.
2. Download size cap
   1. 100 MB.
   2. A `Content-Length` above the cap fails before reading.
   3. Without a length, the body is read through `io.LimitReader(cap+1)`.
   4. One byte past the cap fails.
   5. Error: `NetworkBadResponse`, "image is larger than 100 MB".
3. Pixel cap
   1. 128 megapixels (`16384 × 8192`).
   2. The cap is on width × height, so any shape fits.
   3. `ValidateImage` checks the header before the full decode.
   4. Error: `ImageTooLarge`.
   5. The largest accepted image takes about 512 MB to decode.
4. Constants
   1. `downloadTimeout` and `maxImageBytes` in `internal/network/download.go`.
   2. `maxImagePixels` in `internal/images/utils.go`.
5. Values are about twice what real card art needs.
   1. Raise them only with a real image that fails.

## Consequences

### Positive

1. A stalled URL no longer hangs a dialog forever.
2. A wrong URL cannot fill memory with a huge body.
3. A decompression bomb is rejected from its header.
4. Uploads and downloads share the pixel check.

### Negative and risks

1. Art above 128 MP or 100 MB is refused.
2. A slow download can take up to 120 s.
   1. Progress is shown per [ADR 021](021-image-download-progress.md).
3. A failed image does not fail the save.
   1. The image part is not applied.
   2. The binding result carries `warning`, shown as a toast.

### Neutral

1. Images stored before these limits are not re-checked.
   1. They passed a full decode on upload.
2. A body read error is now `NetworkBadResponse`, not `InternalError`.
3. A timeout, before or during the body, is `NetworkTimeout`.
