# 3. Persist the catalog with fsentry under DeckBuilderData

* **Status:** Accepted
* **Date:** 2026-09-18
* **Authors:** @oleg

---

## Context

The catalog includes nested folders, JSON metadata, and image binaries. A typical desktop user expects a visible data folder they can back up, not a hidden SQLite file they cannot inspect.

## Considered options

1. **SQLite / other SQL** in a user data dir.
2. **Ad-hoc JSON files** written by each service.
3. **fsentry** (`github.com/HardDie/fsentry`) as a folder-tree database with pretty JSON.

## Decision

Use option 3.

`application.Get` constructs one `fsentry.New(cfg.Data, fsentry.WithPretty())` and `Init()`. `internal/db/*` maps aggregates to folders. Repositories attach image files using `cfg` paths.

Data root is `DeckBuilderData` in the working directory, except **macOS uses `~/DeckBuilderData`** because the app bundle is not writable. Generated output is `cfg.Results()` (`…/result`). Tests may call `Config.SetDataPath`.

Settings live in the same store (`internal/db/settings`), not in `os.UserConfigDir`.

## Consequences

### Positive

* Users can copy or zip `DeckBuilderData` as a backup.
* Export/import of a game can zip a folder subtree.
* Pretty JSON is git- and diff-friendly for developers.

### Negative and risks

* Concurrent writers are not a SQL transaction model; this is a single-user app.
* Folder names constrain ids (`BadName` / `GameExist` from fsentry errors).

### Neutral

* Do not introduce a second persistence library for the catalog. Generate output files in `result/` are ordinary PNGs/JSON via `internal/fs`, not fsentry documents.
