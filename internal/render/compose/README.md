# `internal/render/compose`

Runs one game render. It is the only caller of both `internal/render/generate` and `internal/render/sheet`.

`GenerateGame` lists the cards, clears `result/`, and returns. A goroutine then calls `generate.Prepare`, writes the raw back PNGs and the TTS JSON, and calls `internal/render/sheet/write.Draw` for each sheet.

Planning is described in `internal/render/generate`. Drawing and the libjpeg-turbo encode are described in `internal/render/sheet`.

`internal/render/deprecated/generator` and `internal/render/deprecated/page_drawer` are deprecated. compose does not call them.

## What the app passes in

`wire` builds the runner and stores it as `compose.Generator`:

```go
compose.New(cfg, serviceGame, serviceCollection, serviceDeck, serviceCard, serviceSystem, serviceTTS)
```

`main` passes that value to `bindings/generator.New`. The window calls `Game(gameID, sortOrder, scale)`.

`Game` turns a scale below 1 into 1, then calls:

```go
GenerateGame(gameID, compose.GenerateGameRequest{
    SortOrder: sortOrder,
    Scale:     scale,
})
```

| `New` argument | Role |
|---|---|
| `cfg` | `Results()` is the directory that is cleared and then filled. |
| `serviceGame` | `Item` loads the game. An unknown id returns before any folder change. |
| `serviceCollection` | `List` during `catalog.Collect`. |
| `serviceDeck` | `List` during the walk. `GetImage` later, inside `Prepare`. |
| `serviceCard` | `List` during the walk. `GetImage` and `Item` later, inside `Prepare`. |
| `serviceSystem` | `GetSettings`. Failure returns to the caller. The log line is `can't get config`. |
| `serviceTTS` | `SendToTTS` receives the inner bag after the JSON file is written. |

## `GenerateGameRequest`

| Field | Role |
|---|---|
| `SortOrder` | Passed to collection, deck, and card `List`. |
| `Scale` | Integer divisor of the cell, passed through to `generate.Prepare`. |

## `GenerateGame`

The call returns as soon as the goroutine is started. `internal/render/progress.Begin` runs before the return. Status becomes `in_progress` and the percent is 0.

Inside the goroutine, `run`:

1. Calls `generate.Prepare` with `cfg.Results()`, the game, the deck map, the deck order, the request scale, settings, the deck service, and the card service.
2. Calls `progress.Sheets(0, total)`. `total` is the sheet count.
3. Writes each `plan.Backs` entry with `fs.CreateAndProcess` and `fs.BinToWriter`. The bytes are the original back file.
4. For each `plan.Sheets` entry, calls `write.Draw` with the face bytes, back bytes, cell size, shadow flag, and JPEG path, then calls `progress.Sheets(done, total)`.
5. Writes `plan.Root` to `plan.JSONPath` with `fs.CreateAndProcess` and `fs.JsonToWriter`.
6. Calls `SendToTTS` with `plan.Bag`.

On success, status becomes `done`. On error, status becomes `error` and the log line starts with `Generator:`.

Before the goroutine, `fs.RemoveFolder` then `fs.CreateFolder` run on `cfg.Results()`. The previous render is gone.
