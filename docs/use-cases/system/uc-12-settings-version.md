# UC-12: Settings and version

**Module:** `internal/services/system`  
**Status:** Implemented  
**Actors:** GUI settings dialog (gear in the top bar), about chrome  
**Goal:** Persist settings and show the build version  
**Preconditions:** Process started. Version is the stamped git tag, the commit hash, or `dev`

## Main scenario (happy path)

1. `System.GetSettings` returns defaults (`en`, back shadow off, scale 1) merged with fsentry settings.
2. The user clicks the gear in the top bar on any screen (game, collection, deck, card).
3. The settings dialog loads `System.GetSettings` and shows card scale and back shadow.
4. The user saves; `System.UpdateSettings` with `{ "card_scale", "enable_back_shadow", "lang"? }` stores the values that changed.
5. The next render uses them: the card scale goes to `scaleX` and `scaleZ` of every deck and card, `scaleY` stays 1 (TTS scales cards the same way); the shadow darkens the back.
6. `System.GetVersion` returns that version string.

## Alternative scenarios and errors

* **4a. Card scale outside 0.1–10:** the dialog shows an error and does not save; the binding returns `apperr.ErrBadCardScale` ("Card scale must be between 0.1 and 10") and saves nothing.
* **4b. Empty or unknown lang:** the stored language is kept. The dialog has no language field yet.
* **Stored file:** keeps the older `card_size{scaleX,scaleY,scaleZ}` shape. The scale is read from `scaleX`; a missing or invalid one reads as 1, an out-of-range one is clamped (`Settings.Normalize`).

## Postconditions

* Card scale, back shadow, and language are stored in the fsentry settings file when they changed.
* The about chrome shows the process version string.
