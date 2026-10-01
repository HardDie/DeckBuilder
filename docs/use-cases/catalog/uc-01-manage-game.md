# UC-01: Create and edit a game

**Module:** `internal/services/game`  
**Status:** Implemented  
**Actors:** GUI (Wails desktop)  
**Goal:** Persist a named game with description and optional image  
**Preconditions:** App is running; `DeckBuilderData` is initialized (`internal/repositories/core` `Init`)

## Main scenario (happy path)

1. Desktop UI calls `App.CreateGame` with `name`, `description`, optional `image` URL and/or `imageFile` bytes.
2. Service creates an fsentry folder under `games/` and stores metadata + image bytes if provided.
3. Binding returns `data` as `dto.Game` (id, timestamps, `cachedImage` URL for the GUI).
4. Client can `ListGames` (`sort`, `search`), `ReadGame`, `UpdateGame` (same write fields), or `DeleteGame`.

## Alternative scenarios and errors

* **2a. Duplicate or illegal name:** `GameExist` or `BadName`.
* **2b. Unknown id on get/update/delete:** `GameNotExists`.
* **4a. List:** `meta.total` reflects filtered count.

## Postconditions

* The game is a parent for collections. Delete removes the subtree on disk.
