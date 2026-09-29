# UC-12: Settings and version

**Module:** `internal/services/system`  
**Status:** Implemented  
**Actors:** GUI chrome (language, about)  
**Goal:** Persist settings and show the build version  
**Preconditions:** Server started with a version string

## Main scenario (happy path)

1. `System.GetSettings` returns defaults (`en`, back shadow off, scale 1) merged with fsentry settings.
2. `System.UpdateSettings` with `{ "lang": "en"|"ru" }` saves language when it changed.
3. `System.GetVersion` returns this repository's version.

## Alternative scenarios and errors

* **2a. Unknown lang:** ignored; settings unchanged.
* **Settings fields `enable_back_shadow` and `card_size`:** stored and used by generate when present on disk; **`UpdateSettings` currently only applies `lang`.**

## Postconditions

* Language is stored in the fsentry settings file when it changed.
* The about chrome shows the process version string.
