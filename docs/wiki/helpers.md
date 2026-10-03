# Helpers and other `internal/` packages

Developer reference. The map of the project is [Internals](Internals).

## `internal/errors`

Sentinel errors (`GameExist`, `BadName`, …) with optional HTTP status. `.AddMessage` / `.HTTP` for wrapping. `IfErrorLog` logs and continues.

Servers call `network.ResponseError` with these. Repositories map fsentry `ErrExist` / `ErrNotExist` / `ErrBadName` onto them.

## `internal/network`

JSON envelope `{ data, meta, error }`: `ResponseError` for the loopback HTTP routes, `Meta` for binding list results. Image download: `DownloadBytes` (timeout, size cap) and `DownloadProgress`.

## `internal/logger`

`Info` / `Warn` / `Error` std loggers. Repository `List` logs corrupt folders here.

## `internal/fs`

OS files: write generate outputs, JSON helpers. Used by **render/compose**, image, and network. Game export/import uses fsentry `ExportFolder` / `ImportFolder`.

## `internal/images`

Validate upload/download bytes (png/jpeg/…). Repositories refuse invalid images (game still saved, warning log).

## `internal/render/progress`

Process-wide generate status (`empty`, `in_progress`, `done`, `error`) plus done/total sheets. Compose writes it; `bindings/system.Status` reads it. See [Generation](Generation).

## `internal/utils`

`NameToID`, sort helpers, `Allocate`, `NormalizeTimestamps`. Used across servers, services, repositories. `NormalizeTimestamps` fills preview null/empty `createdAt`/`updatedAt` when mapping to entities; it does not rewrite files ([ADR 009](../architecture/009-catalog-timestamps.md)).

## `internal/logger`, `internal/network` callers

Almost every server. Keep messages user-safe; put fsentry detail in `AddMessage` for logs/HTTP body as existing code does.
