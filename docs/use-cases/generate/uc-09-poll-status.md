# UC-09: Poll generate progress

**Module:** `internal/render/progress`  
**Status:** Implemented  
**Actors:** GUI polling system `Status`  
**Goal:** Show type, message, 0–1 progress, and status of the current job  
**Preconditions:** None (singleton always exists)

## Main scenario (happy path)

1. Generate (UC-08) calls `Begin`, then `Sheets(done, total)` per page, then `Finish` or `Fail`.
2. Client `System.Status` receives `dto.Status`.
3. While `in_progress`, further polls return updated progress.
4. When `done` or `error`, **this call flushes** the singleton back to `empty` (`Status`).

## Alternative scenarios and errors

* **1a. Overlapping generate:** rejected with `ErrRenderInProgress` (UC-08 1a); the running job keeps the singleton.
* **4a. Client misses the terminal poll:** next poll may already be `empty`.

## Postconditions

* UI must treat a single `done`/`error` payload as the completion signal and not expect it on the next request.
