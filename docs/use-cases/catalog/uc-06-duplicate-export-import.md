# UC-06: Duplicate, export, and import a game

**Module:** `internal/services/game`  
**Status:** Implemented  
**Actors:** GUI  
**Goal:** Copy a whole game on disk, or move it as a zip  
**Preconditions:** Source game exists (for duplicate/export)

## Main scenario (happy path)

1. **Duplicate:** `App.DuplicateGame(gameID, name)` creates a new game folder tree.
2. **Export:** `GET /api/games/{game}/export` returns a zip of that game.
3. **Import:** `POST /api/games/import` multipart `file` (required) and optional `name` creates a game from the archive.

## Alternative scenarios and errors

* **1a. Target name exists:** `GameExist`.
* **3a. Corrupt zip:** `BadArchive`.

## Postconditions

* Duplicate/import are independent games; later edits do not change the source. Export is read-only.
