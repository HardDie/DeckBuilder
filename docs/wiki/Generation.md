# Generation (`internal/render`)

Developer reference. The map of the project is [Internals](Internals). The product steps are in the [user guide](Guide).

Render turns one game into files under `cfg.Results()/<gameID>` (`<data>/result/<gameID>`): a PNG copy of each deck back, a JPEG face sheet for each page, and one TTS Saved Object JSON. It then asks `tts` to spawn that object.

The URLs written today are local paths: [ADR 005](https://github.com/HardDie/DeckBuilder/blob/master/docs/architecture/005-tts-generate-local-paths.md). Hosting them later is [ADR 016](https://github.com/HardDie/DeckBuilder/blob/master/docs/architecture/016-auto-upload-generated-images.md).

## Packages

| Package | Role |
|---|---|
| `render/compose` | Runs one render; the only caller of `generate` and `sheet`. Field-level notes: `internal/render/compose/README.md`. |
| `render/generate` | Plans: walks the catalog, lays out pages, builds the TTS JSON. Draws nothing. |
| `render/sheet` | Draws one page: decode, resize to the cell (`sheet/stbresize`), paint, libjpeg-turbo quality 80. |
| `render/progress` | Process-wide status the window polls: `empty`, `in_progress`, `done`, `error`. |

## Who starts it

The window calls `bindings/generator.Game`. Scale below 1 becomes 1, then `compose.GenerateGame` runs.

`GenerateGame` returns as soon as the game exists, the cards are listed, and `result/<gameID>/` exists. Drawing continues in a goroutine. The window polls system `Status` about twice a second.

Only one render runs at a time; a second call gets `GenerateInProgress`. Reading `done` or `error` resets the status to `empty`.

## Steps

1. `GetSettings`, then the game. An unknown game returns before any folder change.
2. `catalog.Collect` walks collections, decks, and cards in the requested sort order.
3. `result/<gameID>/` is created if missing. Nothing is deleted yet.
4. In the goroutine, `generate.Prepare` reads every back and face, picks the cell size from the first face of each deck, splits decks into pages of at most 69 faces, and builds the JSON.
5. Back files are written as-is. Each page goes through `sheet/write.Draw`: faces are decoded and resized in parallel with `stb_image_resize2` (Catmull-Rom), then painted and encoded. Progress steps once per page.
6. The JSON is written, files the plan no longer uses are removed, then `SendToTTS` receives the inner bag.

## Reuse

A sheet is named `<deck>_<page>_<hash>.jpg`. The hash (xxHash64) covers everything drawn: the faces in order, the back, the cell, the grid, the shadow flag, and `layout.sheetVersion`. Card text, variables, and count are not drawn, so editing them does not redraw anything. A page whose name already exists is reused. New files are written to `<name>.tmp` and renamed, so a failed run leaves the previous render intact. See [ADR 023](https://github.com/HardDie/DeckBuilder/blob/master/docs/architecture/023-per-game-results-and-page-reuse.md).

A panic while planning or drawing is recovered and reported as `error`.

## Page layout

TTS Custom Deck defaults: at most 10 wide × 7 high. The last cell holds the back (`BackIsHidden` is false), so one page carries 69 faces. A deck with one card is emitted as a Card object, not a Deck.

## Speed

Measured trials and their numbers: [ADR 017](https://github.com/HardDie/DeckBuilder/blob/master/docs/architecture/017-sheet-draw-speed.md). Benchmarks for the live path live in `internal/render/sheet/bench`.
