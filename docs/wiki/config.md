# `internal/config`

Developer reference. The map of the project is [Internals](Internals).

Process-wide paths and generate limits. Constructed once via `config.Get(version)`.

## Fields

| Field | Meaning |
|---|---|
| `Version` | Exact git tag, 12-character commit hash, or `dev` |
| `Data` | fsentry root (`DeckBuilderData` or `~/DeckBuilderData` on darwin) |
| `Game` | Relative catalog folder name, always `"games"` |
| `Cache` | `"cache"` (config only today) |
| `Result` | `"result"` — generate output under Data |
| `CardImagePath` etc. | `printf` templates for `CachedImage` URLs on DTOs |

## Methods

| Method | Meaning | Used by |
|---|---|---|
| `Get` | Defaults + macOS home-dir Data | `application.Get`, tests |
| `Games()` | `Data/games` **OS path** | repositories (zip export/import), some tests as fsentry root |
| `Results()` | `Data/result` | generator, replace, TTS buffer |
| `SetDataPath` | Tests only; retarget Data | `*_test.go` temp dirs |

## Constants

`MinWidth`/`MaxWidth` (2–10), `MinHeight`/`MaxHeight` (2–7), `MaxCount = 10*7 - 1` (69). TTS uses the last face-sheet cell as the hidden image when `BackIsHidden` is false. `MaxFilenameLength` (200) matches fsentry id limits.

`Game` (`"games"`) is the same string stored as `gamesPath` on `internal/repositories` structs. Changing one without the other splits the tree.
