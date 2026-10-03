# Helpers and other `internal/` packages

Developer reference. The map of the project is [Internals](Internals).

## `internal/apperr`

The errors the app can explain to the user. Each has a clear message and one kind: `ErrNotFound`, `ErrAlreadyExists`, `ErrInvalid`, or `ErrBusy`. Details go through `apperr.With` / `Withf`; `errors.Is` still matches the base error and its kind.

- Bound methods return them as-is. `bindings/errfmt.Format` (the Wails `ErrorFormatter`) shows `apperr.Message`; anything unexpected shows "Something went wrong" and is logged.
- The loopback HTTP server (`network.ResponseError`) maps the kind to a status: 404, 409, 400, else 500.
- Repositories map fsentry `ErrExist` / `ErrNotExist` / `ErrBadName` onto them in `repositories.MapFsentry`.
- Unexpected errors are plain errors wrapped with context (`fmt.Errorf("…: %w", err)`). See [ADR 025](../architecture/025-app-errors.md).

## `internal/network`

JSON envelope `{ data, meta, error }`: `ResponseError` for the loopback HTTP routes, `Meta` for binding list results. Image download: `DownloadBytes` (timeout, size cap) and `DownloadProgress`.

## `pkg/logger`

Sets up the standard `log/slog` logger: JSON (or text) lines on the console and in `DeckBuilderData/logs/deckbuilder.log`, rotated at startup when over 5 MB, 3 old files kept ([ADR 026](../architecture/026-log-file-and-rotation.md)). `wire` calls `logger.Init` first. Code logs with `slog` directly; Warn and Error lines carry their source file and line. Wails' own messages come through `wailsLogger` in the module root.

## `internal/fs`

OS files: write generate outputs, JSON helpers. Used by **render/compose**, image, and network. Game export/import uses fsentry `ExportFolder` / `ImportFolder`.

## `internal/images`

Validate upload/download bytes (png/jpeg/…). Repositories refuse invalid images (game still saved, warning log).

## `internal/render/progress`

Process-wide generate status (`empty`, `in_progress`, `done`, `error`) plus done/total sheets. Compose writes it; `bindings/system.Status` reads it. See [Generation](Generation).

## `internal/utils`

`NameToID`, sort helpers, `Allocate`, `NormalizeTimestamps`. Used across servers, services, repositories. `NormalizeTimestamps` fills preview null/empty `createdAt`/`updatedAt` when mapping to entities; it does not rewrite files ([ADR 009](../architecture/009-catalog-timestamps.md)).

## Logging and errors together

Log with a short lowercase message and fields: `slog.Warn("image download failed", "url", u, "err", err)`. User-facing text comes from `apperr` messages; technical detail goes to the log.
