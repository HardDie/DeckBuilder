# HTTP routes

The routes live in `internal/servers`.

`application.Get` calls `servers.Register`.

| Route | Response |
|---|---|
| `GET /api/games/{game}/image` | game image bytes |
| `GET /api/games/{game}/collections/{collection}/image` | collection image bytes |
| `GET /api/games/{game}/collections/{collection}/decks/{deck}/image` | deck image bytes |
| `GET /api/games/{game}/collections/{collection}/decks/{deck}/cards/{card}/image` | card image bytes |
| `GET /api/tts/data` | one-shot TTS JSON |

Details: [Servers](Servers).
