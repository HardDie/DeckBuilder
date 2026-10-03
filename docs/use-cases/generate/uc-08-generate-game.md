# UC-08: Generate TTS sheets and JSON

**Module:** `internal/render/compose`  
**Status:** Implemented  
**Actors:** GUI “Render” on a game  
**Goal:** Write sprite sheets and a TTS Saved Object JSON into `result/<gameID>/`  
**Preconditions:** Game has decks with backs and cards with faces; settings are readable

## Main scenario (happy path)

1. Client calls `generator.Game(gameID, sortOrder, scale)`. `sortOrder` is optional. Scale below 1 becomes 1.
2. The call returns after validating the game and listing cards.
3. `result/<gameID>/` is created if missing; a goroutine plans pages with `internal/render/generate` (TTS Custom Deck grid up to 10×7; last face slot is the hidden image, so 69 playable cards per page), draws them with `internal/render/sheet`, and writes a separate back file plus JSON (`internal/tts_entity` save fields `FaceURL` / `BackURL` / `NumWidth` / `NumHeight`).
4. Image URLs in JSON are **`file:///` local paths** (TTS accepts path/URL). See [ADR 005](../../architecture/005-tts-generate-local-paths.md) and [ADR 007](../../architecture/007-tts-official-api.md).
5. A one-card pile is emitted as a **Card** object, not a Deck (TTS cannot spawn a one-card deck).
6. Service may attempt TTS spawn (UC-10). Progress moves `in_progress` → `done`.

## Alternative scenarios and errors

* **1a. A generate is already running:** `ErrRenderInProgress`, "a render is already running"; `result/` and progress are not touched.
* **1b. A deck has no back image or a card has no image:** `ErrMissingImages`; the message names up to 5 of them and counts the rest. Nothing is written and progress is not touched. See [ADR 024](../../architecture/024-required-deck-and-card-images.md).
* **2a. Unknown game:** error before the goroutine.
* **3a. Unreadable back/face or draw error:** progress `error`; logs the cause.
* **3b. Too many cards for one page:** extra pages (`index` increments); CustomDeck entries multiply.
* **3c. Panic while planning or drawing:** recovered; progress `error`; logs the panic and stack.
* **3d. Page unchanged since the last render:** its file name (a hash of what is drawn) already exists, so it is reused, not redrawn. Card text, variables, and count do not count as changes.
* **3e. Draw or write fails:** the previous files stay; nothing half-written is left. See [ADR 023](../../architecture/023-per-game-results-and-page-reuse.md).

## Postconditions

* `result/<gameID>/` holds exactly this render: reused and new files; files no longer used are removed. Other games' folders are untouched. User copies JSON to TTS `Saved Objects` and keeps the images next to those paths for local play.
