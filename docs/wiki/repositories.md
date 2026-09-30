# `internal/repositories`

Developer reference. The map of the project is [Internals](Internals).

Glue between services and db: image files, zip, HTTP download. Each catalog repo holds `cfg` and the matching `db/*` interface.

## Fields (catalog)

| Field | Meaning |
|---|---|
| `cfg` | `Games()`, image validation helpers |
| `game` / `collection` / `deck` / `card` | `internal/db/*` interfaces |

## Game repository methods

| Method | Disk / db | Used by |
|---|---|---|
| `Create` | `db.Create` then `ImageCreate` if URL or bytes | game service |
| `GetByID` / `GetAll` | `db.Get` / `List` | service Item/List |
| `Update` | `Move` if name changed, then `Update` payload, maybe replace image | service |
| `Delete` | `db.Delete` | service |
| `Duplicate` | `db.Duplicate` | service |
| `Export` | `db.ExportFolder` of that game under `games/` | service, Wails `game.Export` |
| `Import` | `db.ImportFolder` into `games/`, then `GetByID` | service, Wails `game.Import` |
| `GetImage` | `db.ImageGet` | service GetImage |

`Import` with a name asks fsentry to rewrite the folder id and `.info.json` name. Timestamps stay. An existing game with that destination id returns `GameExist` and is left unchanged. Importing under a different id leaves the archive's original game in place.

## Collection / deck / card

Same pattern: CRUD on db, `createImage` / `createImageFromByte` using `network.DownloadBytes` + `images.ValidateImage`, `GetImage` for the image server.

Card repository additionally maps `variables` and `count`.

## Settings repository

| Method | db | Used by |
|---|---|---|
| `Get` | `db/settings.Get`; missing file → defaults | system service |
| `Save` | `Set` | system `UpdateSettings` |

## What repositories are not

They do not parse HTTP. They do not draw sprite sheets (generator + `page_drawer`). They should not open a second fsentry root; they use the db handle and `cfg` OS paths for zip.
