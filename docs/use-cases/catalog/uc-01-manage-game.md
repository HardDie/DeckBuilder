# UC-01: Create and edit a game

**Module:** `internal/services/game`  
**Status:** Implemented  
**Actors:** GUI (or any HTTP client)  
**Goal:** Persist a named game with description and optional image  
**Preconditions:** Server is running; `DeckBuilderData` is initialized (`db/core.Init`)

## Main scenario (happy path)

1. Client `POST /api/games` as multipart: `name`, `description`, optional `image` URL and/or `imageFile`.
2. Service creates an fsentry folder under `games/` and stores metadata + image bytes if provided.
3. Response `200` with `data` as `dto.Game` (id, timestamps, `cachedImage` URL for the GUI).
4. Client can `GET /api/games` (`sort`, `search`), `GET /api/games/{game}`, `PATCH` the same multipart fields, or `DELETE /api/games/{game}`.

## Alternative scenarios and errors

* **2a. Duplicate or illegal name:** `GameExist` or `BadName` (HTTP 400).
* **2b. Unknown id on get/update/delete:** `GameNotExists`.
* **4a. List:** `meta.total` reflects filtered count.

## Postconditions

* The game is a parent for collections. Delete removes the subtree on disk.
