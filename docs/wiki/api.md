# `internal/api`

Registers mux routes and holds **go-swagger** request/response types. Handlers on `Unimplemented*Server` are empty; they exist so `make swagger` can scan comments.

## Files

| File | Registers | Server interface |
|---|---|---|
| `game.go` | `/api/games/import` | `servers/game.Game` |
| `image.go` | `…/image` for game/collection/deck/card | `servers/image` |
| `replace.go` | `/api/replace/prepare`, `/api/replace` | replace server |
| `tts_upload.go` | `/api/tts/data` | TTS server |
| `static.go` | `/`, SPA fallback, `/docs`, swagger JSON | — |

Generate and system routes live next to those servers’ register functions (same pattern).

## How it is used

`application.Get` calls `RegisterGameServer(routes, serverGame)` and the rest. The **real** handler is the implemented server, not `Unimplemented*`.

Do not put fsentry, generate, or validation here. If the GUI needs a new field, add it on a DTO and a swagger struct in the same file as the route.

`internal/api/swagger_headers` is a small helper for generated spec headers.
