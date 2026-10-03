# 26. Log file and rotation

* **Status:** Accepted
* **Date:** 2026-10-04
* **Authors:** @oleg

---

## Context

1. `internal/logger` wrote only to the console.
2. The packaged app has no visible console.
   1. Every log line was lost outside `make dev`.
3. Since [ADR 025](025-app-errors.md), unexpected errors say "Details are in the log".
4. The log is for developers, not shown in the window.
5. A file that only grows would fill the disk over months of use.

## Considered options

1. **Console only**
   1. The previous behavior.
   2. Lost. Nothing survives a packaged run.
2. **A rotation library (e.g. lumberjack)**
   1. Rotates by size while running.
   2. Lost. A new dependency for a problem a few lines solve.
3. **One file, rotated at startup**
   1. Check the size once per launch; rename when too large.
   2. Won. One session never writes megabytes.

## Decision

1. Location
   1. `DeckBuilderData/logs/deckbuilder.log` (`config.Logs()`).
   2. Next to `games/` and `result/`; on macOS under `~/DeckBuilderData`.
2. Writing
   1. Every logger writes to the console and the file.
   2. The file is opened for appending; each line is written at once.
   3. Wails' own messages go through `logger.Wails`, so they land there too.
3. Each launch
   1. `wire` calls `logger.Init` before anything else.
   2. Then one line: `===== App started: version …, os/arch, data folder … =====`.
   3. That line separates runs when reading the file.
4. Rotation at startup
   1. If the log is over 5 MB: `.2` → `.3`, `.1` → `.2`, the log → `.1`.
   2. The oldest (`.3`) is deleted, so at most about 20 MB stays on disk.
   3. Names keep the extension: `deckbuilder.1.log`.
5. Failure
   1. If the folder or file cannot be opened, logging stays on the console.
   2. The app starts anyway; a warning says the log file is off.
6. No link in the window; developers open the folder themselves.

## Consequences

### Positive

1. Unexpected errors can be investigated after the fact.
2. Each run is easy to find by its "App started" line.
3. Disk use is bounded.

### Negative and risks

1. A single very long session can grow the file past 5 MB until the next start.
2. Users must be told where the file is when a developer asks for it.

### Neutral

1. Tests do not call `logger.Init`; they log to the console as before.

## Later

1. 2026-10-04: logging moved to the standard `log/slog`.
   1. `pkg/logger` sets it up: console plus file, rotation as above.
   2. Lines are JSON by default; `FormatText` is the option.
   3. Warn and Error lines carry their source file and line.
2. The launch line is now `"msg":"app started"` with fields, without `=====`.
3. The Wails adapter lives in the module root (`wailslog.go`), not in `pkg/`.
   1. `pkg/logger` stays free of framework code.
4. The console gets every level on stderr.
5. `internal/logger` was removed.
