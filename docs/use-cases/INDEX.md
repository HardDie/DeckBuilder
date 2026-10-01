# Use cases

Scenarios for **behavior that exists in this repository**. Copy [_TEMPLATE.md](_TEMPLATE.md) for a new file. Keep numbering stable.

**Status:** Planned · In progress · Implemented

| ID | Name | Module | Status | File |
|---|---|---|---|---|
| UC-01 | Create and edit a game | `internal/services/game` | Implemented | [catalog/uc-01-manage-game.md](catalog/uc-01-manage-game.md) |
| UC-02 | Manage collections under a game | `internal/services/collection` | Implemented | [catalog/uc-02-manage-collection.md](catalog/uc-02-manage-collection.md) |
| UC-03 | Manage decks (including card back) | `internal/services/deck` | Implemented | [catalog/uc-03-manage-deck.md](catalog/uc-03-manage-deck.md) |
| UC-04 | Manage cards, variables, and count | `internal/services/card` | Implemented | [catalog/uc-04-manage-card.md](catalog/uc-04-manage-card.md) |
| UC-05 | Serve an entity image | `internal/servers` | Implemented | [catalog/uc-05-serve-image.md](catalog/uc-05-serve-image.md) |
| UC-06 | Duplicate, export, and import a game | `internal/services/game` | Implemented | [catalog/uc-06-duplicate-export-import.md](catalog/uc-06-duplicate-export-import.md) |
| UC-07 | Recursive search | `internal/services/search` | Implemented | [search/uc-07-recursive-search.md](search/uc-07-recursive-search.md) |
| UC-08 | Generate TTS sheets and JSON | `internal/services/generator` | Implemented | [generate/uc-08-generate-game.md](generate/uc-08-generate-game.md) |
| UC-09 | Poll generate progress | `internal/progress` | Implemented | [generate/uc-09-poll-status.md](generate/uc-09-poll-status.md) |
| UC-10 | Spawn last generate in TTS | `internal/services/tts` | Implemented | [generate/uc-10-tts-spawn.md](generate/uc-10-tts-spawn.md) |
| UC-11 | Map local image paths to URLs | `internal/services/replace` | Implemented | [replace/uc-11-replace-urls.md](replace/uc-11-replace-urls.md) |
| UC-12 | Settings and version | `internal/services/system` | Implemented | [system/uc-12-settings-version.md](system/uc-12-settings-version.md) |

## Layout

```text
docs/use-cases/
├── INDEX.md
├── _TEMPLATE.md
├── catalog/
├── search/
├── generate/
├── replace/
└── system/
```
