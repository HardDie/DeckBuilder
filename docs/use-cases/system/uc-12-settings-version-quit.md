# UC-12: Settings, version, and quit

**Module:** `internal/services/system`  
**Status:** Implemented  
**Actors:** GUI chrome (language, about, window close)  
**Goal:** Persist settings, show build version, stop the process from the UI  
**Preconditions:** Server started with a version string and debug flag

## Main scenario (happy path)

1. `GET /api/system/settings` returns defaults (`en`, back shadow off, scale 1) merged with fsentry settings.
2. `PATCH /api/system/settings` with `{ "lang": "en"|"ru" }` saves language when it changed.
3. `GET /api/system/version` returns the ldflag / `go install` version string.
4. `DELETE /api/system/quit` (non-debug) exits immediately (`os.Exit(0)`).

## Alternative scenarios and errors

* **2a. Unknown lang:** ignored; settings unchanged.
* **4a. `-debug`:** quit handler returns without exiting.
* **Settings fields `enable_back_shadow` and `card_size`:** stored and used by generate when present on disk; **PATCH currently only applies `lang`.**

## Postconditions

* Debug sessions stay up for Swagger/GUI work. Production GUI can exit the backend when the page unloads.
