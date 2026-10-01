# `internal/page_drawer`

Developer reference. The map of the project is [Internals](Internals). How Render walks a game is [Generation](Generation).

One type, `PageDrawer`, fills one TTS custom-deck page and writes two files: a JPEG face sheet and a PNG copy of the deck back. The generator builds the TTS JSON. This package does not.

## Who calls it

`internal/services/generator` only.

`generateImages` draws real faces and writes the files. `generateJson` walks the same cards with a 10×10 dummy JPEG so page index and slot stay aligned with those files. The dummy walk uses `New` with an empty directory and scale `1`. It never calls `Save`.

## `New`

```go
func New(title, path string, scale, commonIndex int, settings *entitiesSettings.Settings) *PageDrawer
```

| Field | Starting value |
|---|---|
| `index` | `1` (page inside this deck) |
| `commonIndex` | the caller's value (page count across the game, used in the file name) |
| `title` | deck id, used in both file names |
| `path` | directory for the JPEG and the PNG (`cfg.Results()` while drawing; `""` on the JSON walk) |
| `scale` | settings card scale, integer divisor of the cell. `0` becomes `1` on the first `AddImage` |
| `settings` | `EnableBackShadow` only |

`index` is not zero. The first sheet is page 1.

## Methods

| Method | What it does |
|---|---|
| `IsEmpty()` | No faces yet. The back image does not count. |
| `IsFull()` | `len(faces) >= config.MaxCount` (69). |
| `Size()` | Face count. The back is not included. |
| `GetIndex()` | This deck's page number, starting at 1. |
| `AddImage(img []byte) error` | Decode one face and append it. Error `"page is full"` when already full. Decode errors are returned as-is. |
| `SetBacksideImageAndSave(img []byte) (string, error)` | Write the original bytes, then keep a brightness-adjusted copy for the sheet. |
| `Save() (string, int, int, error)` | Write the JPEG. Returns the absolute path, columns, and rows. |
| `Inherit(prev *PageDrawer) *PageDrawer` | Turn a blank drawer into the next page of `prev`. |

`Save` on an empty page returns `"", 0, 0, nil` and writes nothing. `Save` reads the back image. The generator always calls `SetBacksideImageAndSave` before the first face.

## Cell size (`AddImage`)

The decoded image's `Bounds().Max` is the card width and height.

1. If `width * 10 > 10000`, `innerScale` becomes `10000 / 10 / width`.
2. If `trunc(height * innerScale) * 7 > 10000`, `innerScale` becomes `10000 / 7 / height`.
3. A following loop adds `0.01` when the truncated scaled side is still above 10,000, then stops. The two formulas above land near 1,000 and near `10000/7`, so that add does not run.
4. Scale `0` becomes `1`. `innerScale` `0` becomes `1`.
5. When the cell is still 0×0, `width = trunc(cardWidth * innerScale) / scale` and the same for height. Integer division.
6. Later faces do not change the cell. They are resized to it with Lanczos when `Bounds().Max` differs.

Step 2 reads the scale step 1 just wrote. When step 1 does not run, `innerScale` is still 0, the height product is 0, and a tall face is not shrunk. The first face that sets the cell locks it. A later face can still overwrite `innerScale`, and the locked width and height stay.

The resize check compares `Bounds().Max`, not `Dx`/`Dy`.

## Back file

`SetBacksideImageAndSave` writes first, then decodes.

File name:

```text
backside_<title>_<6 hex chars>.png
```

The hex is `md5(raw)[:3]` printed with `%x`. The file bytes are the argument, unchanged. The return value is `fs.PathToAbsolutePath` of that path.

The copy drawn on the sheet is `imaging.AdjustBrightness`. Shadow on uses `-30`. Shadow off uses `0`, which still converts the image. A decode error after the write returns that error and an empty path. The PNG is already on disk.

The sheet copy is resized to the cell inside `Save`, not here.

## `Save`

Grid size is `utils.CalculateGridSize(faceCount + 1)`: the smallest `columns * rows` that can hold the faces plus the back, from 2×2 through 10×7. A short page stays a tight grid. 69 faces plus the back is 10×7.

`images.CreateImage` builds an RGBA canvas of `cellWidth * columns` by `cellHeight * rows`. Faces are drawn left to right, top to bottom (`utils.CardIdToPageCoordinates`). The back is then drawn at `(columns-1, rows-1)`, the bottom-right cell. That cell is not the next free slot when the rectangle has gaps. A 2×2 page with one face leaves the other two cells empty (zero RGBA). JPEG stores those as black.

File name:

```text
<commonIndex>_<title>_<index>_<faceCount>_<columns>x<rows>.jpg
```

Encoded with `images.JpegSaveToWriter` (JPEG quality 80). The return path is absolute. Columns and rows are what the generator stores as `NumWidth` and `NumHeight`.

## `Inherit`

Called as `(&PageDrawer{}).Inherit(prev)`.

| Copied from `prev` | Not copied |
|---|---|
| `index + 1`, `commonIndex + 1` | faces |
| `title`, `path` | `scale`, `innerScale` |
| back image | `settings` |
| cell width and height | |

The next page keeps the cell and the darkened back. `AddImage` on that page still runs the scale math, then resizes into the copied cell because the cell is already set. The generator does not call `SetBacksideImageAndSave` again.

## What this package does not write

No TTS JSON. The generator turns `Save`'s path, columns, and rows into `FaceURL`, `BackURL`, `NumWidth`, and `NumHeight`, with `file:///` in front of the absolute paths. Card `CardID` is `pageIndex * 100 + slot`, and the slot is `Size() - 1` after that face was added.
