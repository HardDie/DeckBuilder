# `internal/db`

Maps catalog aggregates to fsentry. Each package: `contract.go` (interface + request structs), `db.go`, sometimes `model.go` (JSON payload in `.info.json`).

## Two store types

| Package | Handle | Library |
|---|---|---|
| `core` | `fsentry.IFSEntry` | `internal/fsentry` (copied v0.0.11) |
| `settings`, `game`, `collection`, `deck`, `card` | `*fsentry.DB` | `github.com/HardDie/fsentry` v0.1.4 |

Callers: `application.Get` and service tests. Repositories depend only on the **interfaces** in `contract.go`.

## Field: `gamesPath`

On game/collection/deck/card/core:

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
| `Init` | legacy `Init` + `CreateFolder("games")` (ignore exist) | `application.Get`, tests |
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
| `UpdateInfo` | `UpdateFolderNameWithoutTimestamp` (renames folder, keeps timestamps) | repository Import |
| `ImageCreate` / `Get` / `Delete` | binary `"image"` | repository images |

`model`: `Description`, `Image` as `fsentry.QuotedString` (double-encoded JSON strings for display text).

`List` skips folders whose directory name ≠ `.info.json` id (logs “Corrupted game folder”).

## `collection`

Same as game, nested one level. Extra field `game dbGame.Game` so `Create`/`Get`/… first `game.Get` (missing game → game error, not a missing collection).

Path: `gamesPath, game.ID`. Entity `GameID` is filled from the request, not from disk.

Used by collection repository and as a parent lookup from `db/deck`.

## `deck`

| Field | Meaning |
|---|---|
| `db` | `*fsentry.DB` |
| `gamesPath` | `"games"` |
| `collection` | parent `db/collection` |

`Create` makes the deck folder then `CreateFolder[any]("cards", nil, …)` for the card list. Images sit on the deck folder (`"image"` binary = card back).

Used by deck repository and as parent from `db/card`.

## `card`

| Field | Meaning |
|---|---|
| `db` | `*fsentry.DB` |
| `gamesPath` | `"games"` |
| `deck` | parent deck (existence check) |

Cards are **not** folders. The `cards` folder payload is `map[id]*model`. Images are binaries named with the numeric id.

| Method | Disk | Used by |
|---|---|---|
| `Create`/`Update`/`Delete` | `UpdateFolder("cards", list, gamesPath, game, collection, deck)` | card repository |
| `List`/`Get` | `GetFolder[map[int64]*model]("cards", …)` | repository |
| `Image*` | `CreateBinary`/`GetBinary`/`RemoveBinary` under `…/cards` | repository |
