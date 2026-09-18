# 7. Tabletop Simulator integration follows the official Lua / External Editor API

* **Status:** Accepted
* **Date:** 2026-09-18
* **Authors:** @oleg

---

## Context

DeckBuilder’s generate step must produce objects TTS can load (Saved Objects folder) and, when TTS is running, spawn on the table. Lua names, TCP ports, and cardsheet rules are easy to guess wrong. The canonical docs are [api.tabletopsimulator.com](https://api.tabletopsimulator.com/).

## Considered options

1. **Reverse-engineer** save files and Atom plugin traffic only.
2. **Treat official API pages as the contract** and map our Go types onto them: External Editor, `spawnObjectJSON`, Custom Deck, `WebRequest`.

## Decision

Use option 2. When TTS behavior is unclear, read those pages before changing `internal/tts_entity`, `internal/services/tts`, or sheet layout in `internal/config` / `page_drawer`.

Concrete mapping:

| Official API | This app |
|---|---|
| Custom Deck defaults width 10, height 7; last face slot = hidden image unless `back_is_hidden` | `MaxWidth`/`MaxHeight`/`MaxCount`; `BackIsHidden` left false; last cell drawn with the back |
| Save JSON `FaceURL` / `BackURL` / `NumWidth` / `NumHeight` | `tts_entity.DeckDescription`; generate uses `file:///` + absolute paths |
| `spawnObjectJSON({ json = … })` | Lua injected after generate; JSON is one **Bag**, not `ObjectStates` |
| Saved Object file `ObjectStates` | `tts_entity.RootObjects` written to `result/<gameId>.json` |
| TTS listens on **localhost 39999** | `net.Dial("tcp", "127.0.0.1:39999")` |
| Execute Lua Code `messageID` 3, `guid` `-1` = Global | `internal/services/tts.Message` |
| Editor listen port **39998** | Unused (we do not implement an editor server) |
| `WebRequest.get` on the **host** | Lua downloads `http://127.0.0.1:5000/api/tts/data` |
| One-card “decks” are cards (`Card` / `CardCustom`), not `Deck` | Generator emits a single card object when `len == 1` |

Card Lua variables are written to each card’s `LuaScript` as `key="value"` lines (object script in the save), not Custom Deck fields.

## Consequences

### Positive

* Sheet size and hidden-slot reservation match TTS, not an arbitrary grid.
* Spawn uses documented Execute Lua + `spawnObjectJSON` instead of a private plugin protocol.
* Saved Object copy-paste (README) stays the offline path.

### Negative and risks

* Official Lua `setCustomObject` field names differ from save JSON; mix them up and TTS ignores or breaks decks.
* Execute Lua on a non-scripted object GUID fails (`Object reference not set…`); we only target Global `-1`.
* `WebRequest` only runs on the host; spawn will not fetch DeckBuilder’s loopback URL from a client machine.

### Neutral

* Do not listen on 39998 unless we need TTS→app callbacks (`print`, errors, `messageID` 5 return values).
