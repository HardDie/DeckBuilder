# UC-04: Manage cards, variables, and count

**Module:** `internal/services/card`  
**Status:** Implemented  
**Actors:** GUI  
**Goal:** Add face cards with Lua variables and copy count  
**Preconditions:** Parent deck exists

## Main scenario (happy path)

1. Client `POST …/decks/{deck}/cards` multipart: `name`, `description`, `count`, optional `image` / `imageFile`, optional `variables` (JSON object as string).
2. Card is stored with an **integer** id; response is `dto.Card` including `variables` map and `cachedImage`.
3. Get/list/patch/delete use `{card}` as that id.
4. Generate (UC-08) repeats the face according to `count`.

## Alternative scenarios and errors

* **1a. Bad id on later requests:** `CardNotExists` / `BadId`.
* **1b. Missing image:** card can exist; generate fails later if the face cannot be read.

## Postconditions

* Variables are catalog data for TTS Lua, not HTTP query parameters. They must round-trip on GET after PATCH.
