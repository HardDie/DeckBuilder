# UC-07: Recursive search

**Module:** `internal/services/search`  
**Status:** Implemented  
**Actors:** GUI search box  
**Goal:** Find matching games, collections, decks, and cards in one response  
**Preconditions:** Catalog may be empty (then `data` lists are empty)

## Main scenario (happy path)

1. Client `GET /api/search?search=&sort=` (optional filters).
2. Service walks all games (and nested objects) applying the same sort/search rules as list endpoints.
3. Response `data` is `dto.RecursiveSearch` (id tuples, not full entities).
4. Scoped walks: `GET /api/search/games/{game}` and `…/collections/{collection}`.

## Alternative scenarios and errors

* **4a. Unknown game/collection:** not-exists error from the nested service.

## Postconditions

* Clients fetch full records with the id tuples via the catalog GET routes (UC-01–04).
