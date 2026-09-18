# `internal/repositories`

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
| `Export` | zip `cfg.Games()/id` via `internal/fs` | service/server |
| `Import` | unzip into games dir, `GetByID`, `UpdateInfo` if renamed | service/server |
| `GetImage` | `db.ImageGet` | service GetImage |

`UpdateInfo` is the import path: rewrite folder display name without bumping timestamps (`db/game.UpdateInfo` → fsentry `UpdateFolderNameWithoutTimestamp`).

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
