# Helpers and other `internal/` packages

Developer reference. The map of the project is [Internals](Internals).

## `internal/errors`

Sentinel errors (`GameExist`, `BadName`, …) with optional HTTP status. `.AddMessage` / `.HTTP` for wrapping. `IfErrorLog` logs and continues.

Servers call `network.ResponseError` with these. Repositories map fsentry `ErrExist` / `ErrNotExist` / `ErrBadName` onto them.

## `internal/network`

JSON envelope `{ data, meta, error }`. `Response`, `ResponseError`, `RequestToObject`.

## `internal/logger`

`Info` / `Warn` / `Error` std loggers. Repository `List` logs corrupt folders here.

## `internal/fs`

OS files: write generate outputs, JSON helpers. Used by **generator**, image, and network. Game export/import uses fsentry `ExportFolder` / `ImportFolder`.

## `internal/images`

Validate upload/download bytes (png/jpeg/…). Repositories refuse invalid images (game still saved, warning log).

## `internal/render/page_drawer`

The generator's JSON walk uses it so page index and slot match the sheets. Methods, cell size, file names: [Page drawer](Page-drawer). Face sheets are written by `internal/render/sheet/page`.

## `internal/progress`

Process-wide generate status (`empty`, `in_progress`, `done`, `error`). System `Status` **clears** `done`/`error` after read. Used by generator (write) and the system binding (read).

## `internal/utils`

`NameToID`, multipart file extract, sort helpers, `Allocate`, `NormalizeTimestamps`. Used across servers, services, repositories. `NormalizeTimestamps` fills preview null/empty `createdAt`/`updatedAt` when mapping to entities; it does not rewrite files ([ADR 009](../architecture/009-catalog-timestamps.md)).

## `internal/logger`, `internal/network` callers

Almost every server. Keep messages user-safe; put fsentry detail in `AddMessage` for logs/HTTP body as existing code does.
