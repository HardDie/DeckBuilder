# Entities, DTOs, TTS JSON

Developer reference. The map of the project is [Internals](Internals).

Three different structs for “a card”. Do not reuse TTS types as API JSON.

## `internal/entities`

In-process domain. Used by repositories → services → servers (then copied into DTOs).

| Package | Extra fields |
|---|---|
| `game` | `ID`, `Name`, `Description`, `Image` (source URL string), timestamps |
| `collection` | + `GameID` |
| `deck` | + `GameID`, `CollectionID` |
| `card` | + parent ids, `Variables map[string]string`, `Count`, `ID int64` |
| `settings` | `Default()`, lang, shadow, card scale |
| `status` | generate progress snapshot |

`GetName` / `GetCreatedAt` exist for `utils.Sort`.

## `internal/dto`

JSON the GUI sees. Same catalog fields plus `CachedImage` (API URL for `<img src>`).
Catalog `createdAt` / `updatedAt` are RFC3339 strings.
Binding list results carry `network.Meta` with `total` (and `cardsTotal` where relevant).

Recursive search DTO nests filtered games/collections/decks/cards.

Do not add fsentry envelopes (`id`/`data`/`createdAt` at the store layer) to DTOs.

## `internal/tts_entity`

Saved Object / `spawnObjectJSON` shapes: `RootObjects` (`ObjectStates`), `Bag`, `Deck`, `Card`, `DeckDescription` (`FaceURL`, `BackURL`, `NumWidth`, `NumHeight`, `BackIsHidden`).

Used by **generator** (write `result/*.json`) and **replace** (rewrite URLs). TTS card `CardID` is not the catalog `int64`.
