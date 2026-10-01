# `internal/render/sheet`

Draws one TTS custom-deck page and encodes the sheet JPEG.

`internal/render/compose` is the caller that sits outside both this package and `internal/render/generate`. It already has cell size, face bytes, back bytes, and the destination path. It calls `write.Draw` for each page.

## `write.Draw`

```go
func Draw(faces [][]byte, rawBack []byte, cellW, cellH int, shadow bool, path string) error
```

`faces` and `rawBack` are the original file bytes. `cellW` and `cellH` are already chosen.

1. Decode each face.
2. Lanczos-resize it when `Bounds().Max` differs from the cell (`fit.Resize`). A face that is already the cell size stays as it is.
3. Decode the back. `back.Shade` darkens by 30 brightness when `shadow` is set. Brightness 0 still converts the image. Then resize that copy to the cell.
4. `paint.Canvas` places faces left to right, top to bottom, and the back in the bottom-right cell. The canvas is `cellW * columns` by `cellH * rows`.
5. `libjpeg.Encode` writes a quality-80 JPEG.
6. Those bytes are stored at `path`.

`Draw` calls `paint.Canvas` directly. The drawer switch below is in `page.Sheet`.

## `page.Sheet`

`page` keeps one page: faces, a shaded back, a cell, and a file name.

| Call | What it does |
|---|---|
| `New` | Page index starts at 1. `commonIndex`, title, directory, scale, and settings come from the caller. |
| `AddImage` | Decode one face, run `fit.Apply` for the cell, resize to that cell, append. A full page returns `"page is full"`. |
| `SetBacksideImageAndSave` | Write the original bytes with `back.Write`, then keep a shaded copy. |
| `Inherit` | Next page: index and `commonIndex` each increase by 1. The back and the locked cell are copied. Faces are not. |
| `Sheet` | Resize the back to the cell, choose the grid, run the live drawer, and return the RGBA image plus columns, rows, and the JPEG path. An empty page returns a nil image. |
| `Save` | `Sheet`, then `paint.WriteJPEG`. That encoder is `image/jpeg` at quality 80. |

`Sheet` does not encode. The file name it returns:

```text
<commonIndex>_<title>_<page>_<face count>_<columns>x<rows>.jpg
```

A page holds at most `config.MaxCount` faces (69). The back uses the extra cell.

### Switching the drawer

The live call in `Sheet` is `paint.Canvas`. The other drawers are comments, and their packages are not imported, so they are not compiled.

To switch, comment the live call and uncomment one other. Add that package's import. The import path is the comment on the line.

```go
// To switch, comment the live call and uncomment one other.
sheet, _, _ = paint.Canvas(p.size.Width, p.size.Height, p.faces, p.back)
// sheet = seq.Image(p.faces, p.back) // internal/render/sheet/draw/seq
// sheet = row.Image(p.faces, p.back) // internal/render/sheet/draw/row
// sheet = cell.Image(p.faces, p.back) // internal/render/sheet/draw/cell
// sheet = resize.Image(p.faces, p.back) // internal/render/sheet/draw/resize
// sheet = resize_row.Image(p.faces, p.back) // internal/render/sheet/draw/resize_row
// sheet = bilinear.Image(p.faces, p.back) // internal/render/sheet/draw/bilinear
// sheet = pages.Image(p.faces, p.back) // internal/render/sheet/draw/pages
// sheet = rgba.Image(p.faces, p.back) // internal/render/sheet/draw/rgba
```

Each `Image` function returns `*image.RGBA` and does not encode. Those drawers also resize. `paint.Canvas` places the images it is given.

## Subpackages

| Package | Job |
|---|---|
| `fit` | Cell rule (`Apply`) and Lanczos resize (`Resize`) when `Bounds().Max` differs from the cell. |
| `grid` | Smallest columns×rows that can hold *n* cells, from 2×2 through 10×7. `Slot` is left to right, top to bottom. |
| `back` | `Write` stores the original bytes as `backside_<title>_<hash>.png`. `Shade` is the copy drawn on the sheet. |
| `paint` | `Canvas` is the live drawer. `WriteJPEG` encodes with `image/jpeg` at quality 80. |
| `page` | One page, and the drawer switch in `Sheet`. |
| `write` | `Draw`: paint one measured page and encode it with libjpeg-turbo. |
| `draw/libjpeg` | Quality-80 JPEG through libjpeg-turbo (`github.com/pixiv/go-libjpeg`), `DCTISlow`. `Encode` accepts `*image.RGBA`, `*image.YCbCr`, or `*image.Gray`. |
| `draw/jpegli` | Quality-80 comparison encoder. `write.Draw` does not call it. |
| `draw/seq` | Lanczos, then one cell at a time. |
| `draw/row` | One goroutine per row. |
| `draw/cell` | One goroutine per cell. |
| `draw/resize` | Lanczos-resize every face and the back in parallel, then draw in order. |
| `draw/resize_row` | Parallel Lanczos, then one goroutine per row. |
| `draw/bilinear` | `ApproxBiLinear` resize, then the same sequential cell draw as `seq`. |
| `draw/pages` | One page uses `seq`. `JPEGs` encodes several pages at the same time. |
| `draw/rgba` | Lanczos, then `*image.RGBA` cells so `draw.Src` copies rows. |

The `draw/*` packages other than `libjpeg` and `jpegli` each export `Image` (the picture) and `JPEG` (`image/jpeg` at quality 80). The switch calls `Image`.
