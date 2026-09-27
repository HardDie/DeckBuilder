# UC-12: Settings, version, and quit

**Module:** `internal/services/system`  
**Status:** Implemented  
**Actors:** GUI chrome (language, about, window close)  
**Goal:** Persist settings, show build version, stop the process from the UI  
**Preconditions:** Server started with a version string and debug flag

## Main scenario (happy path)

1. `System.GetSettings` returns defaults (`en`, back shadow off, scale 1) merged with fsentry settings.
2. `System.UpdateSettings` with `{ "lang": "en"|"ru" }` saves language when it changed.
3. `System.GetVersion` returns the ldflag / `go install` version string.
4. `System.Quit` (non-debug) exits immediately (`os.Exit(0)`).

## Alternative scenarios and errors

* **2a. Unknown lang:** ignored; settings unchanged.
* **4a. `-debug`:** `Quit` returns without exiting.
* **Settings fields `enable_back_shadow` and `card_size`:** stored and used by generate when present on disk; **`UpdateSettings` currently only applies `lang`.**

## Postconditions

* Debug sessions stay up for Swagger/GUI work. Production GUI can exit the backend when the page unloads.
