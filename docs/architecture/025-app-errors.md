# 25. App errors in `internal/apperr`

* **Status:** Accepted
* **Date:** 2026-10-04
* **Authors:** @oleg

---

## Context

1. Errors lived in `internal/errors`, built for the old REST API.
   1. Each error stored an HTTP status code.
   2. `Error()` printed it: every toast started with "HTTP[400]".
2. The app moved to Wails bindings.
   1. Only a small loopback HTTP server remains: images and TTS data.
3. `AddMessage` often replaced the app's text with fsentry's.
   1. "game not found" became "does not exist".
4. The package hid Go's `errors`, so code used `er` and `stderrors` aliases.
5. Unexpected failures showed raw library text, e.g. "dial tcp: …".

## Considered options

1. **Keep `internal/errors`, strip the prefix**
   1. Lost. Codes would stay in every error for one small server.
2. **Specific errors only, no kinds**
   1. The HTTP server would list every "not found" error by hand.
   2. Lost. A forgotten one would answer 500.
3. **Kinds plus specific errors in `internal/apperr`**
   1. Won.

## Decision

1. Package `internal/apperr`.
   1. The name does not hide Go's `errors`.
   2. One shared package: the same error crosses repositories, services, and bindings.
2. Kinds
   1. `ErrNotFound`, `ErrAlreadyExists`, `ErrInvalid`, `ErrBusy`.
   2. Every specific error belongs to exactly one kind (tested).
3. Specific errors
   1. Named `Err…`, the Go convention (staticcheck ST1012).
   2. Each has a clear lowercase message, e.g. "a deck with this name already exists".
   3. `apperr.With` / `Withf` add details; `errors.Is` still matches.
4. Unexpected errors
   1. Plain errors wrapped with context: `fmt.Errorf("read settings: %w", err)`.
   2. They are logged and shown as "Something went wrong. Details are in the log."
5. Edges
   1. Wails: `bindings/errfmt.Format` is the `ErrorFormatter`; it shows `apperr.Message`, first letter capitalized.
   2. HTTP: `network.ResponseError` maps kinds to 404, 409, 400, else 500.
   3. Errors themselves know neither Wails nor HTTP.
6. `logger.IfError` replaces `errors.IfErrorLog`.

## Consequences

### Positive

1. Toasts read as sentences: "Game not found".
2. The useful part of a message is no longer replaced by library text.
3. Technical detail goes to the log, not to the user.
4. A new error needs one line and gets the right HTTP status automatically.

### Negative and risks

1. The log is only on the console today; review item S22 covers that.
2. Messages are English only, like the rest of the backend.

### Neutral

1. A missing image on `/api/…/image` answers 404 instead of 400.

## Later

1. 2026-10-04: `logger.IfError` was removed with `internal/logger`.
   1. Code logs errors with `slog`: `slog.Warn("…", "err", err)`.
   2. See [ADR 026](026-log-file-and-rotation.md), section "Later".
