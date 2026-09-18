# UC-09: Poll generate progress

**Module:** `internal/progress`  
**Status:** Implemented  
**Actors:** GUI polling `/api/system/status`  
**Goal:** Show type, message, 0–1 progress, and status of the current job  
**Preconditions:** None (singleton always exists)

## Main scenario (happy path)

1. Generate (UC-08) calls `SetType` / `SetMessage` / `SetProgress` / `SetStatus`.
2. Client `GET /api/system/status` receives `dto.Status`.
3. While `in_progress`, further polls return updated progress.
4. When `done` or `error`, **this GET flushes** the singleton back to `empty` (`StatusHandler`).

## Alternative scenarios and errors

* **1a. Overlapping generate:** one singleton; later job overwrites fields (unsupported).
* **4a. Client misses the terminal poll:** next poll may already be `empty`.

## Postconditions

* UI must treat a single `done`/`error` payload as the completion signal and not expect it on the next request.
