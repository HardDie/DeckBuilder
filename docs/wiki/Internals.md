# Internals

How the program is put together. The pages linked from each section are the field-level notes. Image generation has its own walkthrough: [Generation](Generation).

The window is Wails, at the repo root. The same process serves loopback HTTP for picture bytes and for the one-shot TTS payload. Catalog verbs are Wails bindings. They are not REST handlers.

```text
main.go
        │
        ├── wire.go           config, fsentry, services
        │
        └── internal/servers  image and TTS GET
                  │
                  ▼
        internal/services     rules
                  │
                  ▼
        internal/repositories images, zip, fsentry folders
```

Bindings call services directly. `internal/servers` is the HTTP edge that is still left: image GETs and `GET /api/tts/data`.

## `wire.go`

Composition root in package `main`. `wire` builds config, opens one `*fsentry.DB` on the data directory, creates the `games` folder, and constructs the services. `main.go` then builds `servers.New` and the Wails bindings.

Nothing in `wire` decides how a card is stored or how a sheet is drawn.

Details: [Application](Application).

## `servers`

`internal/servers` registers and serves image GETs and `GET /api/tts/data`. It calls services and writes image bytes or the JSON error envelope. It does not import fsentry.

An image handler streams the file a repository already resolved, with a content type. The TTS handler returns the last generated bag. The service then clears that buffer, so a second GET has nothing to serve.

Details: [Servers](Servers). [API](API) lists the routes.

## `services`

Rules live here. Catalog services are thin: list, filter, sort, and pass create/update/delete to the repository. Game adds duplicate, zip export, and zip import.

`generator` is the heavy one. It walks the tree, draws sheets, writes JSON, and asks `tts` to spawn. See [Generation](Generation).

`search` walks as many ids as the caller passed (all games, one game, or one collection). `replace` rewrites `FaceURL` and `BackURL` in a Saved Object file from a mapping. `system` reads settings. Saving settings today persists `lang` when it is `en` or `ru`. `tts` holds the one-shot bag and dials TTS.

Details: [Services](Services).

## `repositories`

Disk details that are not the folder layout itself. They store and load image bytes, check that an upload is a real PNG or JPEG, and they know the OS paths for zip export and import. A bad image is refused. The entity can still be saved. The failure is a warning in the log.

Details: [Repositories](Repositories).

## Folder layout

One `*fsentry.DB` handle, opened on the data root, shared by `internal/repositories`. Games, collections, and decks are folders. A card is a record inside the deck's `cards` folder, with an integer id. Images are sibling binary files. Generate output under `result/` is ordinary files, not fsentry documents.

Details: [DB](DB). Paths and the 10×7 limits: [Config](Config).

## Three shapes for one card

The same card is three different structs, and they are not interchangeable.

| Package | Who reads it | What a card is |
|---|---|---|
| `internal/entities` | Services and repositories | The in-memory catalog. Card id is `int64`. Timestamps are `time.Time`. |
| `internal/dto` | Wails and JSON to the window | The same fields for the GUI. Timestamps are RFC3339 strings, because Wails bindings cannot pass `time.Time`. `cachedImage` is `/api/.../image?<hash>`. |
| `internal/tts_entity` | The generator | TTS save JSON. A card's `Name` is `"Card"`. The printed name is `Nickname`. `CardID` is the sheet slot, not the catalog id. |

Details: [Entities and DTOs](Entities-and-DTOs).

## Bindings and the window

`main.go` starts Wails and the loopback server. `bindings/` is one package per area the window calls: game, collection, deck, card, system, search, generator, replace. A binding maps GUI fields onto a service call and returns the DTO envelope.

Details: [Cmd and tools](Cmd-and-tools).

## Helpers

`page_drawer` packs a sheet. `progress` is the process-wide render status. `images` decodes uploads and draws pixels. `fs` writes the result files. `errors` and `network` shape failures and the `{ data, meta, error }` envelope.

Details: [Helpers](Helpers).

## On disk

```text
<data>/                      DeckBuilderData, or ~/DeckBuilderData on macOS
  games/<gameId>/            name-derived folder id
    <collectionId>/
      <deckId>/
        image.bin            card back
        cards/               card records
  result/                    last render: sheets, backs, <gameId>.json
```

Game, collection, and deck ids come from the name (`NameToID`: lower case, spaces to `_`). Card ids are integers, one higher than the current max.
