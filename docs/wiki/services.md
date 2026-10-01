# `internal/services`

Developer reference. The map of the project is [Internals](Internals).

Rules that are not HTTP and not disk. Each package: `contract.go` + `service.go` (+ `*_test.go` / fuzz).

## Catalog (`game`, `collection`, `deck`, `card`)

Typical fields: `cfg`, `repository*`.

| Method | Behavior | Called from |
|---|---|---|
| `Create` / `Update` / `Delete` / `Item` | Thin pass-through to repository | servers |
| `List(sortField, search)` | `GetAll` then in-memory filter + `utils.Sort` | servers |
| `GetImage` | repository binary + content type | `internal/servers` |
| `Duplicate` (game) | repository duplicate | game server |
| `Export` (game) | zip bytes of the game folder | `game.Export` binding |
| `Import` (game) | create a game from zip bytes; optional name | `game.Import` binding |

Collection/deck/card methods take parent ids (`gameID`, …) from the URL.

## `system`

| Method | Behavior | Called from |
|---|---|---|
| `GetSettings` | defaults merged with stored settings | `GetSettings` binding; generator (scale, shadow) |
| `UpdateSettings` | today only persists `lang` if `en`/`ru` | `UpdateSettings` binding |

Holds `repositorySettings` from `internal/repositories/settings`.

## `generator`

| Method | Behavior |
|---|---|
| `GenerateGame` | Starts a goroutine; the binding returns immediately. Walks game → collections → decks → cards, draws sheets via `page_drawer`, writes PNGs + TTS JSON under `cfg.Results()`, updates `progress`, optionally `services/tts.SendToTTS`. |

Depends on game/collection/deck/card **services** plus system (settings) and TTS. Do not start a second generate; progress is a process singleton.

## `search`

`RecursiveSearch(sort, search, gameID, collectionID)` walks services according to how many ids are set (all games vs one game vs one collection). Used by the search binding.

## `replace`

| Method | Behavior |
|---|---|
| `Prepare` | Unique FaceURL/BackURL keys from a Saved Object JSON |
| `Replace` | Apply mapping file; used so hosted URLs replace `file://` |

Used by `bindings/replace`. Operates on `tts_entity` JSON, not the catalog.

## `tts`

In-memory last generated **bag** JSON for `GET /api/tts/data` (cleared after read). `SendToTTS` is TCP Execute Lua on `127.0.0.1:39999`. Called from generator after a successful render.
