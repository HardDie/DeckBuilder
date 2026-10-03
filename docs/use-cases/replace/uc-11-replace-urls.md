# UC-11: Map local image paths to URLs

**Module:** `internal/services/replace`  
**Status:** Implemented  
**Actors:** User who hosted PNGs somewhere and wants a shareable TTS object  
**Goal:** Rewrite FaceURL/BackURL in generated JSON using a mapping file  
**Preconditions:** A generate JSON file exists (Saved Object / `result` JSON)

## Main scenario (happy path)

1. The replace binding `Prepare` takes the generated JSON bytes.
2. The result lists unique image keys (`Couple`) with empty values for the user to fill.
3. `Replace` takes the original JSON and the mapping (filled keys → HTTP URLs).
4. The result is new JSON with paths replaced; the user can save a table others can load.
5. `Replace` also sends the new root bag to TTS, like generate does ([UC-10](../generate/uc-10-tts-spawn.md)).
   1. The bag is `ObjectStates[0]`, not the `ObjectStates` wrapper.
   2. It replaces whatever `/api/tts/data` still holds.
   3. If TTS spawns it, the user can check the hosted URLs load.

## Alternative scenarios and errors

* **1a. Unexpected ObjectStates shape:** `ErrBadRenderFile`.
* **3a. Mapping missing a FaceURL/BackURL:** `ErrBadRenderFile` / `ErrBadMappingFile`.
* **5a. TTS not running:** dial fails, log only; `Replace` still returns the new JSON and nothing is stored for `/api/tts/data`.

## Postconditions

* Catalog on disk is unchanged. `result/` is unchanged too: the new JSON is only returned to the caller.
* If TTS is listening, the replaced bag is spawned on the table (best effort, as in UC-10).
