# Helpers and other `internal/` packages

## `internal/errors`

Sentinel errors (`GameExist`, `BadName`, …) with optional HTTP status. `.AddMessage` / `.HTTP` for wrapping. `IfErrorLog` logs and continues.

Servers call `network.ResponseError` with these. db maps fsentry `ErrExist` / `ErrNotExist` / `ErrBadName` onto them.

## `internal/network`

JSON envelope `{ data, meta, error }`. `Response`, `ResponseError`, `RequestToObject`. Also opens the browser on non-debug start (`cmd` via application/network).

## `internal/logger`

`Info` / `Warn` / `Error` std loggers. db `List` logs corrupt folders here.

## `internal/fs`

OS files: zip archive/unarchive game folders, write generate outputs, JSON helpers. Used by **game repository** export/import and **generator**, not by `internal/db`.

## `internal/images`

Validate upload/download bytes (png/jpeg/…). Repositories refuse invalid images (game still saved, warning log).

## `internal/page_drawer`

Lays out TTS 10×7 sheets, last cell hidden/back. Used only by **generator**.

## `internal/progress`

Process-wide generate status (`empty`, `in_progress`, `done`, `error`). `GET /api/system/status` **clears** `done`/`error` after read. Used by generator (write) and system server (read).

## `internal/utils`

`NameToID`, multipart file extract, sort helpers, `Allocate`. Used across servers, services, db (legacy card timestamps).

## `internal/fsentry`

Vendored v0.0.11 (`NewFSEntry`, `IFSEntry`, `pkg/fsentry_error`, `pkg/fsentry_types`). Only `db/core` and `db/card` (and their tests) should keep importing this. New aggregates use `github.com/HardDie/fsentry`.

## `internal/logger`, `internal/network` callers

Almost every server. Keep messages user-safe; put fsentry detail in `AddMessage` for logs/HTTP body as existing code does.
