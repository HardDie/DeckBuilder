# `internal/servers`

HTTP adapters: mux vars, multipart, map entities → DTOs, write `network.Response` / `ResponseError`. They call **services**. Wails bindings in `desktop/bindings/` call **services** directly and copy DTO mapping (HTTP servers are not the Wails API).

Each aggregate: `contract.go` (handler interface) + `server.go`.

## Catalog

Game, collection, deck, and card CRUD, plus game export and import, live on Wails bindings (`desktop/bindings/`). No game HTTP handlers remain.

`cachedImage` URLs are built in `desktop/bindings/catalog`.

## Other servers

| Package | Handlers | Downstream |
|---|---|---|
| `image` | `GameHandler`, `CollectionHandler`, `DeckHandler`, `CardHandler` | corresponding `service*.GetImage` |
| `generator` | start generate | `services/generator.GenerateGame` |
| `search` | recursive search | `services/search` |
| `replace` | prepare + replace | `services/replace` |
| `tts` | `GET /api/tts/data` | `services/tts` one-shot buffer |

Servers must not import `internal/db` or fsentry.
