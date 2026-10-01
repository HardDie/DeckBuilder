# 5. Generate TTS sheets locally; share via replace, not auto-upload

* **Status:** Accepted
* **Date:** 2026-09-18
* **Authors:** @oleg

---

## Context

Tabletop Simulator [Custom Deck](https://api.tabletopsimulator.com/custom-game-objects/) objects use a face **cardsheet**, a back image or sheet, and grid size (`width`/`height` in Lua; `NumWidth`/`NumHeight` in save JSON). Official defaults are **10×7**. With `back_is_hidden` / `BackIsHidden` false, the **last slot of the face sheet is the hidden image**. Saved Objects wrap objects in `ObjectStates`; `spawnObjectJSON` takes a **single object** JSON string ([base](https://api.tabletopsimulator.com/base/)). TTS wants HTTP(S) URLs for images if the table is saved and shared. Hosting those images is a product problem (accounts, ToS, bandwidth).

## Considered options

1. **Always upload** faces/backs to a built-in image host during generate.
2. **Write local paths** into JSON; user copies JSON to `Saves/Saved Objects` and keeps `result/` on disk for solo use; **replace** API later swaps paths for URLs the user uploaded.
3. **Only PNG sheets**, no JSON — user builds TTS objects by hand.

## Decision

Use option 2.

`generator.Game` starts a background job: `internal/render/page_drawer` packs faces (max 10×7, last cell reserved → 69 playable cards per page; that last cell is drawn as the hidden/back slot), writes images + JSON under `result/`, updates `internal/progress`. Optional [External Editor](https://api.tabletopsimulator.com/externaleditorapi/) TCP to TTS port **39999** sends Execute Lua (`messageID` 3, Global `guid` `-1`) so the host runs `WebRequest.get` on `/api/tts/data` then `spawnObjectJSON` (buffer is one-shot). The `result/*.json` file is Saved Object shaped (`ObjectStates`); the HTTP payload is the inner **Bag** only.

`POST /api/replace/prepare` lists unique FaceURL/BackURL keys. `POST /api/replace` applies a mapping file. Automatic hosting is explicitly **future / out of scope**.

Sheet limits live in `internal/config` (`MaxWidth`, `MaxHeight`, `MaxCount`).

## Consequences

### Positive

* Generate works offline with no third-party account.
* Shared tables are still possible once the user hosts images and runs replace.
* TTS spawn is optional; missing TTS is not a generate failure.

### Negative and risks

* README must stay honest: local paths do not survive sending a save to a friend without replace.
* Progress is a singleton; overlapping generates are unsafe.

### Neutral

* JSON shapes stay in `internal/tts_entity`. Do not generate a different TTS schema without an ADR.
