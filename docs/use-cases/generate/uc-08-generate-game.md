# UC-08: Generate TTS sheets and JSON

**Module:** `internal/services/generator`  
**Status:** Implemented  
**Actors:** GUI “Render” on a game  
**Goal:** Write sprite-sheet PNGs and a TTS Saved Object JSON into `result/`  
**Preconditions:** Game has decks with backs and cards with faces; settings are readable

## Main scenario (happy path)

1. Client `POST /api/games/{game}/generate` with optional JSON `{ "sortOrder", "scale" }`.
2. Handler returns 200 immediately after validating the game and listing cards.
3. A goroutine clears `cfg.Results()`, packs faces with `page_drawer` (TTS Custom Deck grid up to 10×7; last face slot is the hidden image, so 69 playable cards per page), writes a separate back file plus JSON (`internal/tts_entity` save fields `FaceURL` / `BackURL` / `NumWidth` / `NumHeight`).
4. Image URLs in JSON are **`file:///` local paths** (TTS accepts path/URL). See [ADR 005](../../architecture/005-tts-generate-local-paths.md) and [ADR 007](../../architecture/007-tts-official-api.md).
5. A one-card pile is emitted as a **Card** object, not a Deck (TTS cannot spawn a one-card deck).
6. Service may attempt TTS spawn (UC-10). Progress moves `in_progress` → `done`.

## Alternative scenarios and errors

* **2a. Unknown game:** error before the goroutine.
* **3a. Missing back/face or draw error:** progress `error`; logs the cause.
* **3b. Too many cards for one page:** extra pages (`index` increments); CustomDeck entries multiply.

## Postconditions

* Previous `result/` contents are gone. User copies JSON to TTS `Saved Objects` and keeps PNGs next to those paths for local play.
