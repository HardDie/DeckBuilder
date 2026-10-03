# Internal Go review

Review of `internal/` from 2026-10-02.
Items are discussed one at a time.
Each item goes through: plan → approve → change → review → done.

Status values: `todo`, `planned`, `in review`, `done`, `skipped`.

## Stability

S1. Rename plus image change fails halfway (deck, collection)
 1. Status: `done`.
 2. `Update` calls `move`, so the folder id changes.
 3. `imageDelete` then uses the old `deckID` / `collectionID`.
 4. It returns `DeckNotExists`, but the rename is already on disk.
 5. Game is correct: it uses `newGame.ID`.
 6. Files: `internal/repositories/deck/repository.go`, `internal/repositories/collection/repository.go`.

S2. Card list read-modify-write race
 1. Status: `done`.
 2. Create, update, and delete read the `cards` folder, edit it, and write it back.
 3. Wails runs bindings on separate goroutines.
 4. Two parallel creates can get the same `maxID`, and one write is lost.
 5. File: `internal/repositories/card/repository.go`.

S3. Data race on `tts.dataForTTS`
 1. Status: `done`.
 2. The generate goroutine writes it.
 3. The HTTP handler reads and clears it.
 4. There is no lock.
 5. File: `internal/services/tts/service.go`.

S4. Overlapping generates are not blocked
 1. Status: `done`.
 2. CLAUDE.md says overlap needs a product decision.
 3. The code does not prevent it.
 4. A second call removes `result/` while the first goroutine writes there.
 5. File: `internal/render/compose/compose.go`.

S5. Panic in the generate goroutine kills the app
 1. Status: `done`.
 2. There is no `recover`.
 3. A bad image that panics in decode or resize crashes the app.
 4. Progress never reaches `error`.
 5. File: `internal/render/compose/compose.go`.

S6. Image download has no timeout and no size limit
 1. Status: `done`.
 2. `http.Get` uses the default client, so it has no timeout.
 3. A stalled URL hangs the binding forever.
 4. `io.ReadAll` has no cap, so a huge body can run out of memory.
 5. File: `internal/network/download.go`.
 6. Decision: [ADR 020](../architecture/020-image-input-limits.md).

S7. Failed image replace loses the old image
 1. Status: `done`.
 2. `Update` deletes the old image before downloading the new one.
 3. If the download fails, the entity keeps the new URL but has no image file.
 4. The user gets only a log warning.
 5. Files: game, collection, deck, and card repositories.
 6. Fix: `repositories.ResolveImage` downloads and validates before any write.
 7. The save succeeds; the binding `Result` carries `warning` for a toast.
 8. `make generate` did not pass the libjpeg flags; fixed in the Makefile.

S8. Small robustness fixes
 1. Status: `done`.
 2. `fs.CreateAndProcess` drops the `Close` error.
    1. A full disk can leave a truncated JPEG or JSON silently.
 3. `network.ResponseError` uses a type assertion, not `errors.As`.
    1. Wrapped `*Err` values become HTTP 500.
 4. Card `Update` does not clamp `Count < 1` like `Create` does.
 5. `tts.DataForTTS` returns a plain `errors.New`.
    1. `/api/tts/data` then logs "unhandled error" and answers 500.
 6. Count is now at least 1 everywhere: service writes, repository reads, GUI form.

S9. Only one app instance may run
 1. Status: `done`.
 2. Today a second launch works and binds `:5001`.
 3. Both instances then share the same `DeckBuilderData`.
 4. fsentry's lock file is taken per call, not for the process lifetime.
 5. So S2's read-modify-write can still race across processes.
 6. Wails v2.16 has `options.SingleInstanceLock`.
 7. File: `main.go`.

S10. Per-game result folder, build then swap
 1. Status: `done`.
 2. Today every generate removes all of `result/`.
 3. Generating game B breaks game A's saved object in TTS.
    1. The JSON holds absolute `file:///` paths into `result/`.
 4. A failed run also deletes the last good result.
 5. Idea: `result/<gameID>/`, one copy per game.
    1. Build in `result/.tmp-<gameID>-<random>/`.
    2. On success, swap the temp folder into place.
    3. On failure, drop the temp folder; the old result stays.
 6. JSON paths must point at the final folder, not the temp one.
 7. Cleanup: game delete, game rename, and leftover temp folders at startup.
 8. Check by hand whether TTS caches `file:///` images by path.
 9. Touches compose, `generate/layout`, CLAUDE.md, README, UC-08.
 10. Merged with content-hashed page names and reuse; no temporary folder.
 11. Decision: [ADR 023](../architecture/023-per-game-results-and-page-reuse.md).

S11. Progress for image downloads
 1. Status: `done`.
 2. A URL image downloads inside the create / update binding.
 3. Only a spinner shows until the call returns, up to 120 s (S6).
 4. Idea: report bytes read and total so the GUI can draw a bar.
 5. Needs a design: transport, which dialog, and cancel.
 6. Decision: [ADR 021](../architecture/021-image-download-progress.md).

S12. A file-uploaded image cannot be removed
 1. Status: `skipped`.
 2. A file upload stores the URL as `""`.
 3. "No URL, no file" then looks unchanged, so nothing is cleared.
 4. Same as before S7.
 5. Needs an explicit "remove image" flag from the GUI.
 6. Skipped by design: decks and cards must have an image ([ADR 024](../architecture/024-required-deck-and-card-images.md)).

S13. Large local uploads are slow over Wails IPC
 1. Status: `todo`.
 2. `imageFileBytes` sends the file as a JSON array of numbers.
 3. That is about 4× the file size, parsed on both sides.
 4. Options: base64 string, or a loopback upload route.

S14. Faster face resize
 1. Status: `done`.
 2. Resize was ~735 ms of a ~1535 ms page.
 3. `fit.Resize` now uses `stb_image_resize2` (Catmull-Rom, SIMD).
 4. Opaque images use the 4-channel path; transparent ones are alpha-aware.
 5. Page now ~854 ms on an M4.
 6. Decision: [ADR 022](../architecture/022-stb-image-resize.md).

S15. Remove the deprecated Lanczos resize
 1. Status: `done`.
 2. Wait for release builds on all four targets with stb.
 3. Check a render on Windows, Linux amd64, Linux arm64, and macOS.
 4. Then delete `fit.ResizeLanczos` and `BenchmarkStageResizeLanczos`.
 5. Then `imaging` stays only for `back.Shade`.
 6. 2026-10-03: release builds pass on all four targets with stb.
 7. Still open: a render checked on each platform.
 8. Removed after the four release builds; a test still compares with `imaging` Lanczos.

S16. Stable JSON order
 1. Status: `done`.
 2. Collection bags and a card's Lua lines followed Go map order.
 3. Two renders of an unchanged game could give different JSON.
 4. Now bags are sorted by collection id, Lua lines by variable name.
 5. Files: `generate/script`, `tts_entity/card.go`.

S17. Card scale and back shadow cannot be changed in the app
 1. Status: `todo`.
 2. `UpdateSettingsRequest` has only `Lang`; the GUI sends only `lang`.
 3. CLAUDE.md lists both as persisted settings.
 4. They stay at the defaults unless `settings.json` is edited by hand.
 5. Needs binding fields, validation, and a settings form.

S18. Missing settings fields read as zero
 1. Status: `done`.
 2. A `settings.json` without `CardSize` gives scale 0.
 3. TTS would spawn cards at zero size.
 4. Fix: fill missing fields from `entitiesSettings.Default()`.
 5. Today only a hand-edited or very old file triggers it.
 6. Done: `system.GetSettings` returns `Settings.Normalize()`; the file is not rewritten on read.

## Readability

R1. Merge the game / collection / deck repositories
 1. Status: `done`.
 2. The three files are about 90% the same, roughly 1,000 lines.
 3. Only the path depth and the error values differ.
 4. Idea: one shared folder helper keyed by path and error values.
 5. The card repo should reuse `MapFsentry` and `ImageBytes`.
 6. With one shared helper, S1 cannot happen again.

R2. Fold `er.MissingAncestor` into `MapFsentry`
 1. Status: `done`.
 2. Today the same check comes before almost every `MapFsentry` call.

R3. Remove one-line wrappers in repositories
 1. Status: `done`.
 2. `GetByID → get`, `GetAll → list`, `DeleteByID → delete`, `Duplicate → duplicate`.

R4. Delete dead code
 1. Status: `done`.
 2. These `sheet/draw/*` variants are used only by the bench and their own tests.
    1. `seq`, `bilinear`, `rgba`, `cell`, `row`, `resize`, `resize_row`, `pages`.
    2. `jpegli` is used only by the bench.
    3. `libjpeg` is live and stays.
 3. `render/deprecated/*` is used only as an oracle in `compose/golden_test.go`.
 4. `images.ImageSize`, `fs.IsFolderExist`, and `fs.CreateFolderIfNotExist` are unused.
 5. These error values are unused: `GameInfoNotExists`, `CollectionInfoNotExists`, `CardExists`.
 6. Also removed: `sheet/page`, `fit.Apply`, `back.Write`, `paint.WriteJPEG`, `utils/generator.go`, REST leftovers in `network` and `utils`.
 7. `go.mod` drops `gen2brain/jpegli` (and its `wazero` runtime).
 8. The compose goldens are now the byte reference on their own.

R5. `GetName()` silently lowercases
 1. Status: `done`.
 2. Search only works because of this hidden behavior.
 3. Make it explicit inside `FilterByName` and `Sort`.
 4. Done with R6.2: `Sort(items, field)` uses `slices.SortStableFunc`.

R6. Small cleanups
 1. Status: `done`.
 2. `utils.Sort(&items, …)` could become `slices.SortStableFunc` on a plain slice. Done in R5.
 3. `write.run` could become `errgroup` with `SetLimit`.
 4. `system.GetSettings` copies fields that the repository already filled with defaults.
 5. `system.UpdateSettings` uses `log.Println` instead of `logger`.
 6. Item 3 skipped: `errgroup` needs `golang.org/x/sync`.
    1. It returns the first error in time, not in job order.
    2. It does not recover panics, so `callJob` stays anyway.
 7. Items 4 and 5 done; `system` service got its first test.

## Order

1. S1–S5.
2. S6–S8.
3. R1–R3.
4. R4–R6.
5. S10, after S4.
6. S11, after S7.
7. S14, then S10 with page hashes, then S15 after a release.

## Added later

S19. Remind about missing deck and card images
 1. Status: `done`.
 2. List: `hasImage` in deck and card DTOs; "No image" placeholder and badge.
 3. Save: warning toast when a deck or card has no image.
 4. Render: `GenerateMissingImages` before start; nothing written.
 5. Fixed: a deck suggestion without an image sent the 404 body as a file.
 6. Decision: [ADR 024](../architecture/024-required-deck-and-card-images.md).
 7. A deck is also badged when any of its cards has no image (`cardsMissingImage`).

S20. Show render failures in the window
 1. Status: `todo`.
 2. A render that ends with `error` only hides the progress circle.
 3. The reason is only in the log.
 4. Idea: carry the reason in the polled status and show it as a toast.

S21. Error toasts start with "HTTP[400]"
 1. Status: `done`.
 2. `*Err.Error()` prints the HTTP code first; bindings send that text.
 3. Every error toast shows it, e.g. "HTTP[400] game exist".
 4. Warnings already strip it (`catalog.ImageWarning`).
 5. Scope grown: rework errors fully into `internal/apperr`.
    1. No HTTP codes in errors; the HTTP server maps error kinds to status codes.
    2. `Err…` names, readable messages, one `With` helper for details.
    3. Unexpected errors show "Something went wrong"; details go to the log.
    4. Done in steps, reviewed one by one.
 5. Service image tests no longer download from github.com; a local server serves the images.
 6. `internal/errors` deleted. Decision: [ADR 025](../architecture/025-app-errors.md).

S22. Rework logging
 1. Status: `in review` (pkg/logger: step 3 of 3).
 2. S21 sends unexpected error details only to the log.
 3. The log must then be easy to find and read.
 4. Today: `internal/logger` writes to stdout (Info, Warn) and stderr (Error) only.
 5. A packaged app has no visible console, so those lines are lost.
 6. Ideas: a log file in the data folder, rotation, a "copy log" action in the window.
 7. Done: `DeckBuilderData/logs/deckbuilder.log`, rotated at startup over 5 MB, 3 kept.
 8. "App started" line per launch; Wails messages included; no window link (for developers).
 9. Decision: [ADR 026](../architecture/026-log-file-and-rotation.md).
 10. Generalized into `pkg/logger` on `log/slog`, JSON by default, text optional.
    1. Step 1: the package and its tests.
    2. Step 2: switch the app and every call site to `slog`.
    3. Step 3: delete `internal/logger`; update ADR 026, CLAUDE.md, the wiki.
 11. The Wails adapter moved to the module root (`wailslog.go`); `pkg/logger` has no framework code.
