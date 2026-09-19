# `internal/servers`

HTTP adapters: mux vars, multipart, map entities → DTOs, write `network.Response` / `ResponseError`. They call **services**. Wails bindings in `desktop/bindings/` call **services** directly and copy DTO mapping (HTTP servers are not the Wails API).

Each aggregate: `contract.go` (handler interface) + `server.go`.

## Shared constructor fields

Catalog servers typically hold:

| Field | Meaning |
|---|---|
| `cfg` | For `CachedImage` URL templates |
| `service*` | Use-case methods |

## Catalog servers (`game`, `deck`)

Game, collection, and card CRUD live on Wails bindings (`desktop/bindings/game`, `desktop/bindings/collection`, `desktop/bindings/card`). HTTP catalog servers left:

| Handler | HTTP | Service call |
|---|---|---|
| `ListHandler` | GET decks | `List(sort, search)` |
| `CreateHandler` | POST multipart | `Create(…)` |
| `ItemHandler` | GET `{id}` | `Item(id)` |
| `UpdateHandler` | PATCH multipart | `Update` |
| `DeleteHandler` | DELETE | `Delete` |
| `ExportHandler` / `ImportHandler` | game zip | repository via service |

`calculateCachedImage` fills `dto.*.CachedImage` from `cfg.*ImagePath`.

## Other servers

| Package | Handlers | Downstream |
|---|---|---|
| `image` | `GameHandler`, `CollectionHandler`, `DeckHandler`, `CardHandler` | corresponding `service*.GetImage` |
| `generator` | start generate | `services/generator.GenerateGame` |
| `search` | recursive search | `services/search` |
| `replace` | prepare + replace | `services/replace` |
| `tts` | `GET /api/tts/data` | `services/tts` one-shot buffer |
| `system` | quit, settings, status, version | `services/system`, `progress` |

Servers must not import `internal/db` or fsentry.
