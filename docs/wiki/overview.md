# Overview

DeckBuilder’s Go process is a loopback HTTP server (`127.0.0.1:5000`). The Vue GUI is a separate repo; production assets are embedded from `web/dist`.

## Request flow (catalog example)

Creating a game:

1. Mux hits a handler registered in `internal/api` (`POST /api/games` → `servers/game.CreateHandler`).
2. The **server** parses multipart (`name`, `description`, `image`, `imageFile`) and calls the **service**.
3. The **service** forwards to the **repository** (no extra rules on create).
4. The **repository** calls **db** `Create`, then optionally downloads or stores a binary image.
5. **db** talks to fsentry: `CreateFolder` under `games/`.
6. The server maps the **entity** to a **dto** (`CachedImage` URL) and `network.Response` wraps `{ "data": … }`.

List/search add filtering and sorting in the **service**, not in db.

## Two fsentry handles

`application.Get` builds both:

| Variable | Type | Module | Used by |
|---|---|---|---|
| `fs` | legacy `internal/fsentry` `IFSEntry` (`NewFSEntry`) | vendored v0.0.11 | `db/core`, `db/card` |
| `db` | `github.com/HardDie/fsentry` `*DB` (`New` + `Init`) | v0.1.x | `db/settings`, `db/game`, `db/collection`, `db/deck` |

Both point at `cfg.Data`. Production `db.Init()` takes the lock file `.fsentry.lock`. Tests that mix old `core` with new aggregates pass `WithNoLockFile()`.

## On-disk layout (production)

DB root = `cfg.Data` (`DeckBuilderData`, or `~/DeckBuilderData` on macOS).

```text
<data>/
  settings.json          # db/settings entry (v0.1 GetEntry/CreateEntry)
  .fsentry.lock          # v0.1 inter-process lock
  games/                 # created by db/core.Init (legacy CreateFolder)
    <gameId>/            # folder + .info.json (description, image URL)
      image.bin
      <collectionId>/
        image.bin
        <deckId>/
          image.bin      # card back
          cards/         # folder whose payload is the card map
            <cardId>.bin
  result/                # generate output (ordinary files, not fsentry)
  cache/                 # reserved in config; not the catalog
```

`cfg.Games()` is `filepath.Join(Data, "games")`. That is a **filesystem path** for zip export/import in repositories, not an fsentry path segment.

Service tests often open fsentry at `cfg.Games()` and still pass the `"games"` segment, which nests `data/games/games`. Production opens at `cfg.Data` and **must** pass `"games"`.

## IDs

- Game / collection / deck folder ids: `NameToID(name)` (lowercase, spaces → `_`, strip punctuation). Display `Name` is stored separately (QuotedString in `.info.json`).
- Card ids: `int64`, allocated as max existing id + 1, stored in the `cards` folder payload, not as one folder per card.
