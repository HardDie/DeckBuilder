# 18. Generation progress in `internal/render/progress`

* **Status:** Accepted
* **Date:** 2026-10-01
* **Authors:** @oleg

---

## Context

1. One goroutine draws every sheet.
2. `internal/render/compose` owns that loop.
3. `internal/render/generate` only plans.
4. `internal/render/sheet` only draws one page.
5. The heavy work is the sheet loop.
6. The window should show that a run is in progress.
7. `MainView.vue` already polls `bindings/system.Status`.
8. That poll feeds the `n-progress` circle.
9. One sheet stays at 0% until that sheet finishes.
10. The circle must stay up for that stretch.
11. Many sheets should step the same bar.
12. `Prepare` and `write.Draw` must stay unaware of progress.

## Considered options

1. **Callback field on `GenerateGameRequest`**
   1. This was the previous implementation.
   2. The request carried `Report func(done, total int)`.
   3. Compose called it from the sheet loop.
   4. A nil `Report` skipped the update.
   5. `bindings/generator` passed `reportSheets`.
   6. That function did `done/total*100`.
   7. It wrote the percent into `internal/render/deprecated/progress`.
   8. What it takes here: a field on every `GenerateGame` call.
   9. The generator binding must remember to pass the function.
   10. Percent math lives in the binding, not next to the loop.
   11. Lost. Two owners of one number.
   12. The request mixes a UI hook into the render call.
   13. A forgotten `Report` leaves the bar stuck.
2. **`internal/render/progress` as the only writer**
   1. This is the current implementation.
   2. The package stores done, total, percent, and status.
   3. Compose calls `Begin` before `GenerateGame` returns.
   4. `Begin` sets `in_progress` and zeroes the bar.
   5. After the plan, compose calls `Sheets(0, total)`.
   6. `total` is `len(plan.Sheets)`.
   7. After each sheet, compose calls `Sheets(done, total)`.
   8. Percent is `done/total*100`.
   9. A total of 0 keeps the percent at 0.
   10. Success calls `Finish`. Failure calls `Fail`.
   11. `bindings/system.Status` reads `Get`.
   12. A `done` or `error` read then calls `Reset`.
   13. `MainView.vue` keeps the 500 ms poll.
   14. The circle shows while status is `in_progress`.
   15. That includes a 0% value.
   16. One sheet moves from 0 to 100.
   17. Many sheets step once per sheet.
   18. What it takes here: one package and calls in the existing loop.
   19. No new binding method. No new route. No Wails import in compose.
   20. Won. One owner of the numbers.
   21. The poll that already drives the circle stays the only transport.
3. **Wails events from that same state**
   1. Keep `internal/render/progress` as the numbers.
   2. After each `Sheets` call, emit a window event.
   3. `runtime.EventsEmit` needs a Wails context.
   4. Compose must not import Wails.
   5. The emit would sit in a binding or in `main`.
   6. Vue would `EventsOn` and drop the status interval.
   7. What it takes here: a context in the goroutine, an event name, and a listener in `MainView.vue`.
   8. Lost. The circle already polls `Status`.
   9. Events would be a second channel for the same state.
   10. Push is the usual choice when nothing is already polling.
4. **A new HTTP poll just for progress**
   1. Add a loopback route that returns done, total, and percent.
   2. Vue would `fetch` it on an interval.
   3. Or it would replace the Wails `Status` call.
   4. What it takes here: a handler, a route, and a second client in the window.
   5. Images and TTS stay on HTTP.
   6. Desktop status already left that server.
   7. Lost. A new route duplicates `bindings/system.Status`.
5. **A wrapper outside compose**
   1. A decorator around `compose.Generator`.
   2. It sees `GenerateGame` start and return.
   3. It does not enter the sheet loop.
   4. What it takes here: another type in front of `New`.
   5. Per-sheet progress needs a second catalog walk.
   6. Or the plan has to leak out of compose.
   7. Lost. Sheet boundaries are visible only inside the loop.
6. **Reuse `internal/render/deprecated/progress`**
   1. The old singleton stores type, message, percent, and status.
   2. The deprecated generator still writes it.
   3. That writer uses a card count, not the sheet count.
   4. What it takes here: compose calls `SetProgress` on that singleton.
   5. `Status` would keep reading the old package.
   6. Lost. The new path would share a mailbox with the old generator.
   7. `internal/render/deprecated/progress` is not part of this path.

## Decision

Use option 2.

1. The bar is determinate.
2. The unit is the sheet count.
3. `internal/render/progress` is the only writer of done, total, and percent.
4. Compose calls `Begin`, `Sheets`, `Finish`, and `Fail`.
5. `bindings/system.Status` is the only reader for the window.
6. `MainView.vue` keeps polling that status.
7. A second transport is extra.
8. Wails events are the usual push when no poll exists yet.
9. This app already polls.
10. An outside wrapper cannot see sheet boundaries.
11. `Prepare` and `write.Draw` stay unaware.
12. `internal/render/deprecated/progress` is not imported on this path.

## Why the package beats the request callback

1. Option 1 put the formula in `bindings/generator`.
2. Option 2 puts done, total, and percent in one package.
3. The generator binding no longer stores a percent.
4. `GenerateGameRequest` is only sort order and scale.
5. Callers cannot forget the hook.
6. `Sheets` is the same `done, total` call.
7. It is no longer a field the caller must pass.

## Best practice for this app

1. One state owner.
2. One transport the UI already has.
3. Push events when the UI has no poll.
4. Do not add HTTP for a value `Status` already returns.
5. Do not count sheets from outside the loop that owns them.
6. Keep planning and drawing free of progress calls.

## Consequences

### Positive

1. Planning and drawing stay free of progress calls.
2. The window keeps one poll.
3. The percent has one formula.

### Negative and risks

1. One process-wide value. A second render overwrites it.
2. The bar moves by sheets, not by pixels inside one sheet.

### Neutral

1. `internal/render/deprecated/progress` stays for the deprecated generator.
2. `Prepare` and `write.Draw` are unchanged.
