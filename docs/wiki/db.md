# `internal/repositories`

The code lives in `internal/repositories`.

Developer reference. The map of the project is [Internals](Internals).

Maps catalog aggregates to fsentry. Each package: `contract.go` (interface + request structs) and `repository.go`. Game, collection, and deck store `repositories.FolderModel` through `repositories.Folder` (see [Repositories](Repositories)); card and settings keep their own `model.go`.

## Handle

All packages take `*fsentry.DB` from `github.com/HardDie/fsentry` v0.1.7 (`New` + `Init`). Callers: `wire` and service tests. Repositories depend only on the **interfaces** in `contract.go`.

## Field: `gamesPath`

On deck/game/card/core as a field, and as a constant inside `repositories.Folder`:

```go
gamesPath string // always "games"
```

This is an **fsentry path segment** (first argument after the object name), not `cfg.Data` and not `cfg.Games()`.

With DB rooted at `cfg.Data`:

| Call | Path args | Directory |
|---|---|---|
| settings `GetEntry("settings")` | *(none)* | `<data>/settings.json` |
| game `CreateFolder(name, …, gamesPath)` | `"games"` | `<data>/games/<id>` |
| collection | `"games", game.ID` | `<data>/games/<game>/<collection>` |
| deck | `"games", gameID, collection.ID` | `…/<deck>` |
| cards folder | `"games", gameID, collectionID, deckID` | `…/cards` |

Omitting `gamesPath` would put games next to `settings.json`. Do not confuse it with `cfg.Games()`, which is the **OS path** `filepath.Join(Data, "games")` for zip.

`core.Init` creates the `games` folder so later `CreateFolder` parent dirs exist.

## `core`

| Method | Meaning | Used by |
|---|---|---|
| `Init` | `Init` + `CreateFolder[any]("games", nil)` (ignore exist) | `wire`, tests |
| `Drop` | delete whole store root | tests cleanup |

## `settings`

No `gamesPath`. One JSON **entry** named `settings`.

| Field | Meaning |
|---|---|
| `db` | `*fsentry.DB` |

| Method | fsentry | Used by |
|---|---|---|
| `Get` | `GetEntry[SettingInfo]("settings")` — payload is `info.Data`, not a field on `SettingInfo` | settings repository |
| `Set` | `CreateEntry` or `UpdateEntry` if exists | settings repository |

`SettingInfo`: `Lang`, `EnableBackShadow`, `CardSize.{ScaleX,Y,Z}`.

## `game`

| Field | Meaning |
|---|---|
| `db` | `*fsentry.DB` |
| `gamesPath` | `"games"` |

| Method | fsentry | Used by |
|---|---|---|
| `Create` / `Get` / `List` / `Move` / `Update` / `Delete` | folder CRUD under `games/` | game repository |
| `Duplicate` | `DuplicateFolder` | repository Duplicate |
| `Export` / `Import` | `ExportFolder` / `ImportFolder` under `games/` | repository Export/Import |
| `ImageCreate` / `Get` / `Delete` | binary `"image"` | repository images |

`model`: `Description`, `Image` as `fsentry.QuotedString` (double-encoded JSON strings for display text).

`List` skips folders whose directory name ≠ `.info.json` id (logs “Corrupted game folder”).

## `collection`

Same as game, nested one level. Fields are `cfg`, `*fsentry.DB`, and `gamesPath`.

Path: `gamesPath, gameID`. Entity `GameID` is the caller's id.

Package: `internal/repositories/collection`. The deck repository passes the collection id as a path segment.

## `deck`

| Field | Meaning |
|---|---|
| `db` | `*fsentry.DB` |
| `gamesPath` | `"games"` |

`Create` makes the deck folder then `CreateFolder[any]("cards", nil, …)` for the card list. Images sit on the deck folder (`"image"` binary = card back).

Package: `internal/repositories/deck`. The card repository passes the deck id as a path segment.

## `card`

| Field | Meaning |
|---|---|
| `db` | `*fsentry.DB` |
| `gamesPath` | `"games"` |

Cards are **not** folders. The `cards` folder payload is `map[id]*model`. Images are binaries named with the numeric id.

| Method | Disk | Used by |
|---|---|---|
| `Create`/`Update`/`Delete` | `UpdateFolder("cards", list, gamesPath, game, collection, deck)` | card repository |
| `List`/`Get` | `GetFolder[map[int64]*model]("cards", …)` | repository |
| `Image*` | `CreateBinary`/`GetBinary`/`RemoveBinary` under `…/cards` | repository |
