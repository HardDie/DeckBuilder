# UC-03: Manage decks (including card back)

**Module:** `internal/services/deck`  
**Status:** Implemented  
**Actors:** GUI  
**Goal:** CRUD a deck whose **image is the card back** used at generate time  
**Preconditions:** Parent game and collection exist

## Main scenario (happy path)

1. Client `POST …/collections/{collection}/decks` with multipart fields and a back image (file or URL).
2. Deck id is stored; `GET …/decks/{deck}/image` later returns the back.
3. `GET /api/games/{game}/decks` lists decks from **all** collections in that game.
4. Nested list remains `GET …/collections/{collection}/decks`.

## Alternative scenarios and errors

* **2a. Duplicate name:** `DeckExist`.
* **Generate (UC-08) without a back image:** generate fails when `GetImage` for the deck errors.

## Postconditions

* Cards in this deck share this back on every generated page for that deck.
