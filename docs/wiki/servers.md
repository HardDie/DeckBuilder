# `internal/servers`

Developer reference. The map of the project is [Internals](Internals).

One package serves the loopback GETs. It calls **services** and writes image bytes or `network.ResponseError`. Wails bindings in `bindings/` call **services** directly. Catalog writes are not HTTP.

`cachedImage` URLs are built in `bindings/catalog`.

| Route | Downstream |
|---|---|
| `GET /api/games/{game}/image` | `game.GetImage` |
| `GET /api/games/{game}/collections/{collection}/image` | `collection.GetImage` |
| `GET /api/games/{game}/collections/{collection}/decks/{deck}/image` | `deck.GetImage` |
| `GET /api/games/{game}/collections/{collection}/decks/{deck}/cards/{card}/image` | `card.GetImage` |
| `GET /api/tts/data` | `services/tts` one-shot buffer |

`internal/servers` must not import `internal/repositories` or fsentry.
