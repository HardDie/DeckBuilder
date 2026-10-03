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
* **1c. Image URL or file fails (download, timeout, size, format):** the card is still saved; the image part is not applied (create stores no URL, update keeps the old image and URL); the result carries `warning`, shown as a toast. Limits: [ADR 020](../../architecture/020-image-input-limits.md).
* **1d. Count below 1:** create and update store 1; old or imported data with a lower count reads as 1.

## Postconditions

* Variables are catalog data for TTS Lua, not HTTP query parameters. They must round-trip on GET after PATCH.
