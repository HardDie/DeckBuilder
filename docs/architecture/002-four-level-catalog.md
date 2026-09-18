# 2. Catalog hierarchy: game → collection → deck → card

* **Status:** Accepted
* **Date:** 2026-09-18
* **Authors:** @oleg

---

## Context

Card games for TTS are not a flat list of images. Users think in expansions and decks (example: Munchkin → base set → Monster deck → individual monster cards with HP/AT variables).

## Considered options

1. **Flat list of cards** with tags.
2. **Three levels** (game / deck / card).
3. **Four levels:** game, collection (base/DLC), deck, card.

## Decision

Use option 3. HTTP paths nest the same way:

`/api/games/{game}/collections/{collection}/decks/{deck}/cards/{card}`

Games, collections, and decks use **string ids** (fsentry folder ids derived from names). Cards use **`int64` ids**. A deck’s **image is the card back** used when generating sheets. Cards have `variables` (string map for TTS Lua) and `count` (copies in the generated deck).

`GET /api/games/{game}/decks` lists every deck in a game (all collections) for UI shortcuts.

## Consequences

### Positive

* Matches the user mental model in the README and GUI.
* Generate can walk collections then decks then cards without extra tagging.

### Negative and risks

* Nested URLs are verbose; clients must keep four ids on card screens.
* Renaming a game/collection/deck is an fsentry folder rename, not a surrogate key.

### Neutral

* Recursive search (`/api/search…`) returns typed id tuples (`dto.RecursiveSearch`) rather than flattening the tree into one entity type.
