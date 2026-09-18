# UC-11: Map local image paths to URLs

**Module:** `internal/services/replace`  
**Status:** Implemented  
**Actors:** User who hosted PNGs somewhere and wants a shareable TTS object  
**Goal:** Rewrite FaceURL/BackURL in generated JSON using a mapping file  
**Preconditions:** A generate JSON file exists (Saved Object / `result` JSON)

## Main scenario (happy path)

1. Client `POST /api/replace/prepare` multipart `file` = generated JSON.
2. Response lists unique image keys (`Couple`) with empty values for the user to fill.
3. Client `POST /api/replace` with `file` (original JSON) and `mapping` (filled keys → HTTP URLs).
4. Response is new JSON with paths replaced; user can save a table others can load.

## Alternative scenarios and errors

* **1a. Unexpected ObjectStates shape:** `ErrorInvalidDeckDescription`.
* **3a. Mapping missing a FaceURL/BackURL:** `ErrorInvalidDeckDescription` / `ErrorInvalidMapping`.

## Postconditions

* Catalog on disk is unchanged. Only the JSON blob is transformed.
