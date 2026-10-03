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

* **2a. Duplicate name:** `ErrDeckExists`.
* **Generate (UC-08) without a back image:** render does not start; the error names the deck (UC-08 1b).
* **2b. Image URL or file fails (download, timeout, size, format):** the deck is still saved; the image part is not applied (create stores no URL, update keeps the old image and URL); the result carries `warning`, shown as a toast. Limits: [ADR 020](../../architecture/020-image-input-limits.md).
* **2c. Saved without an image:** the deck is saved; the result carries the reminder "This deck has no image…" as a warning toast, and the deck list shows it with a "No image" placeholder and a warning badge (`hasImage` is false).
* **2d. A card in the deck has no image:** the deck tile gets a warning badge, "Some cards have no image" (`cardsMissingImage`); the cards are not named.

## Postconditions

* Cards in this deck share this back on every generated page for that deck.
