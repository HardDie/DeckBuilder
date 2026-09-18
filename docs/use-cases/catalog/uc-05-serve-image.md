# UC-05: Serve an entity image

**Module:** `internal/servers/image`  
**Status:** Implemented  
**Actors:** GUI `<img>`, generator (via services), browsers  
**Goal:** Return stored image bytes for game/collection/deck/card  
**Preconditions:** Entity exists and has image data on disk

## Main scenario (happy path)

1. Client `GET /api/games/{game}/image` (or nested `…/collections/{id}/image`, `…/decks/{id}/image`, `…/cards/{card}/image`).
2. Server loads bytes from the corresponding service `GetImage`.
3. Response is image bytes (`png` / `jpeg` / `gif`) with HTTP 200.

## Alternative scenarios and errors

* **2a. No image:** `*ImageNotExists` (or equivalent) as JSON error envelope.
* **2b. Unknown type:** `UnknownImageType`.

## Postconditions

* `dto.*` `cachedImage` fields point at these paths so the SPA does not read the data directory itself.
