# `internal/render/compose`

Runs one render of a game: catalog in, sheet images, back images, and a Tabletop Simulator (TTS) Saved Object JSON out. **Start here** to understand the render layer.

## The render layer at a glance

```text
window ──► bindings/generator.Game(gameID, sortOrder, scale)
                │
                ▼
     compose.GenerateGame        checks, game folder, then returns to the window
                │ goroutine
                ▼
     generate.Prepare            the plan: what goes on which sheet, file names, TTS JSON
                │
                ▼
     for each sheet:             reuse the file if its name exists,
       sheet/write.Draw          otherwise draw and encode it
                │
                ▼
     backs, JSON, cleanup,       then progress is "done"
     tts.SendToTTS
```

| Package | Question it answers | Pixels? | Writes files? |
|---|---|---|---|
| `render/generate` | What goes where, under which name? | No, only image headers | No |
| `render/sheet` | What does one page look like? | Yes | One JPEG per call |
| `render/compose` | In what order, and what survives a failure? | No | Backs, JSON, cleanup |
| `render/progress` | How far along is the render? | No | No |

`generate` and `sheet` never import each other. `compose` is the only caller of both, so planning is testable without pixels and drawing without a catalog.

## The logic in short

1. **One render at a time.** A second call while one runs gets `ErrRenderInProgress`.
2. **One folder per game.** Output goes to `result/<gameID>/`. Other games' folders are never touched.
3. **Unchanged pages are reused.** A sheet's file name holds a hash of everything drawn on it, so an existing name means identical content and the page is not drawn again.
4. **A failure keeps the last good render.** New files are written to a temporary name and renamed. Old files are removed only after the whole render succeeded.
5. **The window polls progress.** It steps once per sheet; reused sheets step instantly.

Decisions behind this: [ADR 018](../../../docs/architecture/018-generation-progress.md) (progress), [ADR 022](../../../docs/architecture/022-stb-image-resize.md) (resize), [ADR 023](../../../docs/architecture/023-per-game-results-and-page-reuse.md) (folders and reuse).

## Details

### Construction

`wire` builds the runner once and stores it as `compose.Generator`:

```go
compose.New(cfg, serviceGame, serviceCollection, serviceDeck, serviceCard, serviceSystem, serviceTTS)
```

| `New` argument | Used for |
|---|---|
| `cfg` | `Results()`; the game folder is `Results()/<gameID>`. |
| `serviceGame` | `Item` loads the game. |
| `serviceCollection`, `serviceDeck`, `serviceCard` | `List` while walking the catalog; deck and card `GetImage`, card `Item` inside `Prepare`. |
| `serviceSystem` | `GetSettings`: shadow on/off and the TTS card scale. |
| `serviceTTS` | `SendToTTS` receives the game bag at the end. |

The window calls `bindings/generator.Game(gameID, sortOrder, scale)`. It turns a scale below 1 into 1 and calls `GenerateGame` with a `GenerateGameRequest`:

| Field | Role |
|---|---|
| `SortOrder` | Passed to collection, deck, and card `List`. |
| `Scale` | Integer divisor of the cell size; 2 draws cells at half size. |

### `GenerateGame`: the synchronous part

These steps run before the call returns. Any error here goes straight back to the window. An error in steps 1–3 changes no file.

1. Take the "running" flag (`atomic.Bool`, compare-and-swap). If it is taken: `ErrRenderInProgress`.
2. Read settings. Load the game; an unknown id fails here.
3. `catalog.Collect` lists every collection, deck, and card.
4. Remove files lying directly in `result/` (left by versions before per-game folders).
5. Create `result/<gameID>/` if it is missing. Nothing inside it is deleted.
6. `progress.Begin`: status `in_progress`, 0%.
7. Start the goroutine. From here the goroutine owns the flag and releases it when it ends.

If any step before 7 fails, a `defer` releases the flag, so a failed start never blocks the next render.

### `run`: the goroutine

1. `generate.Prepare(dir, …)` returns the `Plan`. It reads image bytes and headers, nothing else.
2. `progress.Sheets(0, total)`.
3. **Backs.** Each `plan.Backs` file is written unless its name already exists. The bytes are the original back image.
4. **Sheets.** For each `plan.Sheets` entry: if the file exists, count it as reused; otherwise call `write.Draw`. Then `progress.Sheets(done, total)`.
5. **JSON.** `plan.Root` is written to `<gameID>.json` every time, since card text may have changed.
6. **Cleanup.** Every file in the game folder that steps 3–5 did not write or reuse is removed (old pages, old backs, leftover `.tmp` files). Failures here are logged, not fatal.
7. Log `render finished` with the game, sheet count, and reused count.
8. `SendToTTS(plan.Bag)`: best effort; TTS may not be running.

New files in steps 3–5 go through `fs.WriteAtomic`: written to `<name>.tmp`, then renamed. A crash leaves at most a `.tmp` file, which the next successful render removes.

### Errors, panics, and progress

| Outcome | Status | Window | Log | Files |
|---|---|---|---|---|
| Success | `done` | progress circle closes | `render finished` (game, sheets, reused) | Exactly this render's files |
| An image cannot be read | `error` | toast: "An image in deck "…" could not be read. The file may be damaged." | `render failed` plus the decoder's text | Previous files kept; new complete files may be added |
| Other error in `run` | `error` | toast with the error's message, or "Something went wrong…" | `render failed` (game, err) | Same |
| Panic in `run` or in a draw worker | `error` | toast "Something went wrong…" | `render failed` with `panic: … <stack>` | Same |

`safeRun` recovers a panic in the goroutine and turns it into an error. `sheet/write` does the same for its worker goroutines. A crash inside C code (libjpeg-turbo, stb) cannot be recovered.

The window reads status through `bindings/system.Status`. Reading `done` or `error` resets it to `empty`. An `error` status carries the readable message (`errfmt.Text` of the error kept by `progress.Fail`); the window shows it as a toast.

An image that does not decode is reported per deck: `layout` cannot read the first face's header (`*layout.ImageError`), or `write.Draw` cannot decode a face or the back (`write.ErrUndecodable`). Both become `generate.UnreadableImage(deckName)`.

### Reuse: when is a page drawn again?

A sheet name is `<deck>_<page>_<hash>.jpg`. `generate/layout` computes the hash from what is drawn: the face images in order, the back image, the cell size, the grid, the shadow flag, and `layout.sheetVersion`. So:

| Change | Sheet redrawn? | JSON changes? |
|---|---|---|
| Card name, description, variables, count | No | Yes |
| A face image, the card order | Only the pages whose faces changed | Yes |
| The deck's back image | Every page of that deck | Yes |
| Shadow setting, scale | Every page of the game | Yes |
| Drawing code (resize, paint, encoder) | Only if `sheetVersion` is bumped | — |

Bump `sheetVersion` in `generate/layout` whenever drawing changes pixels; otherwise old pages are reused.

## Tests

| File | Covers |
|---|---|
| `compose_test.go` | Unknown game, game folder and legacy cleanup, progress per sheet, overlap rejection, flag release after errors and panics |
| `reuse_test.go` | Reuse, text-only changes, one changed face, a failed render keeping the previous files |
| `golden_test.go` | Byte-exact output against `testdata/golden`. `WRITE_GOLDEN=1` rewrites the goldens; review that diff before committing. |

The tests use `generate/fake`, an in-memory catalog. The fixture faces are already cell-sized, so the goldens do not exercise resizing; `sheet/stbresize` tests that.
