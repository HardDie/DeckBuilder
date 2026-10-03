# `internal/repositories`

Developer reference. The map of the project is [Internals](Internals).

Glue between services and fsentry: image files, zip, HTTP download.

## Shared folder store (`folder.go`)

Game, collection, and deck are the same kind of data: a folder with a description, an image URL, and an optional `image` binary. `repositories.Folder` stores one such level; each repository wraps it.

| Item | Meaning |
|---|---|
| `FolderModel` | `.info.json` payload: `description`, `image` |
| `FolderErrors` | that level's `Exist`, `NotExist`, `ImageExist`, `ImageNotExist` |
| `parent` | path below `games/`: `nil`, `{gameID}`, or `{gameID, collectionID}` |
| `Create` / `Update` | `ResolveImage` first (download + validate), then write; a bad image is returned as `Saved.ImageError` and not applied |
| `Get` / `List` / `Delete` / `Image` | `List` skips unreadable or mismatched folders |

Errors go through `MapFsentry`, which maps a missing parent (`missingParent`) before the level's own errors. All of them are `apperr` errors.

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

`Import` with a name asks fsentry to rewrite the folder id and `.info.json` name. Timestamps stay. An existing game with that destination id returns `ErrGameExists` and is left unchanged. Importing under a different id leaves the archive's original game in place.

## Collection / deck

Thin wrappers over `Folder` that add the parent ids (`GameID`, `CollectionID`) to the entity.

Deck also creates the `cards` folder on create and has `GetAllDecksInGame` (unique by name and image URL).

## Card

All cards of a deck live in one `cards/.info.json` map, so the card repository does not use `Folder`. It uses `ResolveImage` and `MapFsentry` the same way. A mutex serializes the read-modify-write of that map. It also maps `variables` and `count` (at least 1).

## Settings repository

| Method | Call | Used by |
|---|---|---|
| `Get` | `internal/repositories/settings` `GetEntry("settings")`; missing file → defaults | system service |
| `Save` | `Set` | system `UpdateSettings` |

## What repositories are not

They do not parse HTTP. They do not draw sprite sheets (`internal/render`). They should not open a second fsentry root; they use the `*fsentry.DB` handle and `cfg` OS paths for zip.
