# UC-02: Manage collections under a game

**Module:** `internal/services/collection`  
**Status:** Implemented  
**Actors:** GUI  
**Goal:** CRUD collections that group decks (base set vs DLC)  
**Preconditions:** Parent `{game}` exists

## Main scenario (happy path)

1. Desktop UI calls `App.CreateCollection` with `name`, `description`, optional `image` URL and/or `imageFile` bytes.
2. Collection is stored under that game; `dto.Collection` is returned.
3. List/get/update/delete use `ListCollections`, `ReadCollection`, `UpdateCollection`, `DeleteCollection`.

## Alternative scenarios and errors

* **1a. Missing game:** collection APIs fail with not-exists on the parent.
* **2a. Duplicate collection name:** `CollectionExist`.

## Postconditions

* Decks are created under a collection id, not directly under the game (except the all-decks listing in UC-03).
