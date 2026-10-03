# UC-10: Spawn last generate in TTS

**Module:** `internal/services/tts`  
**Status:** Implemented  
**Actors:** Generator after writing JSON; Replace after rewriting URLs ([UC-11](../replace/uc-11-replace-urls.md)); TTS Lua; `GET /api/tts/data`  
**Goal:** If Tabletop Simulator is listening, spawn the generated bag on the table  
**Preconditions:** Generate produced a bag; [External Editor](https://api.tabletopsimulator.com/externaleditorapi/) TCP on the TTS host may or may not be up

## Main scenario (happy path)

1. `SendToTTS` stores the marshaled **Bag** (not `ObjectStates`) and dials TTS’s server port **`127.0.0.1:39999`**.
2. It writes [Execute Lua Code](https://api.tabletopsimulator.com/externaleditorapi/): `messageID` **3**, `guid` **`-1`** (Global), `script` calling [`WebRequest.get`](https://api.tabletopsimulator.com/webrequest/manager/) on `http://127.0.0.1:5000/api/tts/data`.
3. Host-side TTS GETs that URL; `DataHandler` returns the bag JSON; the in-memory buffer is **cleared**.
4. Lua checks `request.is_error`, then [`spawnObjectJSON`](https://api.tabletopsimulator.com/base/) with `json = request.text` and a `callback_function`.

## Alternative scenarios and errors

* **1a. TTS not running / nothing listening on 39999:** dial fails; generate still succeeded; log only.
* **2a. Execute Lua on a GUID that has no script:** TTS errors (`Object reference not set…`). We only use Global `-1`.
* **3a. Second GET `/api/tts/data`:** "nothing to send to Tabletop Simulator" (HTTP 404, `ErrNothingForTTS`).
* **4a. `request.is_error`:** Lua prints the error and does not spawn.

## Postconditions

* Spawn is best-effort. `result/<gameId>.json` (`ObjectStates`) remains what the user copies into TTS `Saves/Saved Objects` (README). This app does not listen on editor port **39998**.
