# 23. Per-game result folders and page reuse

* **Status:** Accepted
* **Date:** 2026-10-03
* **Authors:** @oleg

---

## Context

1. Every render deleted all of `result/`, then drew every page again.
2. The JSON holds absolute `file:///` paths into `result/`.
   1. Users copy it into TTS Saved Objects.
   2. Rendering game B deleted game A's images, so A broke in TTS.
3. A failed render had already deleted the last good one.
4. Users render again and again while editing a few cards.
   1. Every page was redrawn, about 0.85 s each after [ADR 022](022-stb-image-resize.md).
5. TTS caches images; a changed page under an old name could show stale art.

## Considered options

1. **`result/<game>_<timestamp>/` per render**
   1. Nothing breaks, but disk grows with every click.
   2. Any cleanup brings the breakage back.
   3. Lost.
2. **`result/<gameID>/`, replaced through a temporary folder**
   1. One copy per game; a failed run keeps the old one.
   2. Still redraws every page.
   3. Lost to option 3, which keeps the same guarantees.
3. **`result/<gameID>/` with content-hashed names**
   1. A page name carries a hash of what is drawn on it.
   2. An existing name means identical content, so it is reused.
   3. Won.

## Decision

1. Folder
   1. Each game renders into `result/<gameID>/`.
   2. Other games' folders are never touched.
   3. Files left directly in `result/` by older versions are removed.
2. Names
   1. Sheet: `<deck>_<page>_<hash>.jpg`.
   2. Back: `backside_<deck>_<hash>.png`, a hash of the back bytes.
   3. JSON: `<gameID>.json`, rewritten on every render.
3. Page hash
   1. xxHash64 (`cespare/xxhash/v2`); no security need, so the fastest wins.
   2. Inputs in order: `sheetVersion`, cell width and height, columns, rows, shadow.
   3. Then the back bytes, then each face's bytes in page order.
   4. Each image is prefixed with its length.
   5. Card name, description, variables, and count are not drawn, so not hashed.
   6. A swapped card order changes the face order, so the hash.
4. `layout.sheetVersion`
   1. Bump it when drawing changes pixels: resize, paint, shadow, encoder, quality.
   2. Then no page from an older build is reused.
5. Writes
   1. New files go to `<name>.tmp`, then are renamed (`fs.WriteAtomic`).
   2. A crash leaves only a temporary file, never a broken page.
   3. Files the plan did not use are removed only after a successful run.
6. Game lifecycle
   1. Renaming a game moves its folder, so its pages stay reusable.
   2. Deleting a game removes its folder.

## Consequences

### Positive

1. Rendering one game no longer breaks another game's saved object.
2. A failed render keeps the last good one.
3. Unchanged pages cost a hash, not a draw.
4. A changed page gets a new URL, so TTS cannot show a cached old image.

### Negative and risks

1. Forgetting to bump `sheetVersion` reuses pages drawn by the old code.
2. One copy per game stays on disk until the game is deleted.
3. Saved objects copied before a rename still point at the old folder.
4. Renaming or deleting a game while it renders races with the render.
   1. Only one render runs at a time; the window does both rarely.

### Neutral

1. Two decks with the same id and identical pages share one file.
2. The compose goldens hold the new names; their image bytes did not change.
