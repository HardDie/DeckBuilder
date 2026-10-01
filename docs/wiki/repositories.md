# `internal/repositories`

Developer reference. The map of the project is [Internals](Internals).

Glue between services and fsentry: image files, zip, HTTP download. Each catalog repository holds `cfg`, `*fsentry.DB`, and `gamesPath`.

## Fields (catalog)

| Field | Meaning |
|---|---|
| `cfg` | `Games()`, image validation helpers |
| `db` | `*fsentry.DB` |
| `gamesPath` | `"games"` |

## Game repository methods

| Method | Disk | Used by |
|---|---|---|
| `Create` | `CreateFolder`, then image if URL or bytes | game service |
| `GetByID` / `GetAll` | `GetFolder` / `List` | service Item/List |
| `Update` | `MoveFolder` if name changed, then payload, maybe replace image | service |
| `Delete` | `RemoveFolder` | service |
| `Duplicate` | `DuplicateFolder` | service |
| `Export` | `ExportFolder` of that game under `games/` | service, Wails `game.Export` |
| `Import` | `ImportFolder` into `games/`, then `GetByID` | service, Wails `game.Import` |
| `GetImage` | `GetBinary` `"image"` | service GetImage |

`Import` with a name asks fsentry to rewrite the folder id and `.info.json` name. Timestamps stay. An existing game with that destination id returns `GameExist` and is left unchanged. Importing under a different id leaves the archive's original game in place.

## Collection / deck / card

Same pattern: folder CRUD, `createImage` / `createImageFromByte` using `network.DownloadBytes` + `images.ValidateImage`, `GetImage` for the image server.

Card repository additionally maps `variables` and `count`.

## Settings repository

| Method | Call | Used by |
|---|---|---|
| `Get` | `internal/repositories/settings` `GetEntry("settings")`; missing file → defaults | system service |
| `Save` | `Set` | system `UpdateSettings` |

## What repositories are not

They do not parse HTTP. They do not draw sprite sheets (generator + `internal/render/page_drawer`). They should not open a second fsentry root; they use the `*fsentry.DB` handle and `cfg` OS paths for zip.
