# `internal/api`

Registers mux routes. Handlers live in `internal/servers`.

## Files

| File | Registers | Server interface |
|---|---|---|
| `image.go` | `…/image` for game/collection/deck/card | `servers/image` |
| `tts_upload.go` | `/api/tts/data` | TTS server |

## How it is used

`application.Get` calls `RegisterImageServer` and `RegisterTTSServer`.

Do not put fsentry, generate, or validation here. If the GUI needs a new field, add it on a DTO.
