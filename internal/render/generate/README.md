# `internal/render/generate`

Plans one render of a game. It decides which card faces go on which sheet, how big a cell is, what every file is called, and what the TTS JSON says. It draws nothing and writes nothing.

## Where this fits

`render/compose` calls `Prepare`, gets a `Plan`, and then does the work: it hands each sheet to `render/sheet` and writes the files. This package never imports `render/sheet`; see [compose](../compose/README.md) for the whole layer.

## The logic in short

1. **Walk the catalog** (`catalog.Collect`): every collection, deck, and card, grouped by deck, decks ordered by name.
2. **Load the raw images** (`Prepare`): one back per deck, one face per card. No decoding; only the image header is read later, for the size.
3. **Lay out pages** (`layout.Pages`): split each deck into sheets of at most 69 faces, pick the cell size from the deck's first face, and name every file with a hash of its content.
4. **Build the TTS object** (`script.Build`): one game bag, one bag per collection, a deck (or a single card) per deck, one card object per copy.
5. **Return the `Plan`.** The caller writes it.

A card's `count` only adds copies in the JSON. On the sheet each card's face appears once.

## `Prepare`

```go
func Prepare(dir string, gameItem *entitiesGame.Game,
    decks map[catalog.Deck][]catalog.Card, order []catalog.Deck,
    scale int, cfg *entitiesSettings.Settings,
    deckSvc servicesDeck.Deck, cardSvc servicesCard.Card) (Plan, error)
```

| Argument | Role |
|---|---|
| `dir` | The game's output folder, `result/<gameID>`. Paths in the plan are absolute and inside it. |
| `gameItem` | Its id names the JSON file; its name is the game bag's nickname. |
| `decks`, `order` | From `catalog.Collect`. |
| `scale` | Integer divisor of the cell; 0 counts as 1. |
| `cfg` | `EnableBackShadow` goes onto every sheet; `CardSize` is the TTS scale of decks and cards. |
| `deckSvc`, `cardSvc` | `GetImage` for raw bytes; card `Item` for name, description, variables, and count. |

The back comes from the deck in the first card's collection. A missing back or face stops the plan with that error, and the log names game, collection, deck, and card.

### The `Plan`

| Field | Contents |
|---|---|
| `Sheets` | One per page: original face bytes in order, original back bytes, `CellW`, `CellH`, `Shadow`, and the JPEG path. Everything `sheet/write.Draw` needs. |
| `Backs` | One per deck: path and the original back bytes. |
| `JSONPath` | `<dir>/<gameID>.json`. |
| `Root` | The Saved Object file: `{"ObjectStates": [game bag]}`. |
| `Bag` | The same game bag alone, which is what TTS spawns. |

## `catalog`: walking the game

`Collect` calls collection `List`, then deck `List` for each collection, then card `List` for each deck, all with the caller's sort order and an empty search.

| Record | Fields |
|---|---|
| `Deck` | `ID`, `Name`, `Image`: used as the map key |
| `Card` | `ID`, `GameID`, `CollectionID`, `Count` |

- A deck with no cards is left out.
- Decks with the same id, name, and image in two collections share one key. Their cards are concatenated, and become one set of sheets.
- The returned `order` is by deck name, ascending, stable. Collection and card order stay as `List` returned them.

## `layout`: pages, cells, and names

`Pages(dir, decks, scale, shadow)` measures. It reads only the width and height from the image header (`image.DecodeConfig`), for every accepted format: PNG, JPEG, and GIF. A first face whose header cannot be read returns `*layout.ImageError` with the deck id; `Prepare` turns it into `UnreadableImage(deckName)`: "an image in deck … could not be read".

**Splitting.** Faces fill a page in card order until 69 (`config.MaxCount`); the next face starts a new page. Page numbers start at 1 for each deck. The 70th cell of a full 10×7 page holds the back.

**Grid.** The smallest columns × rows, from 2×2 to 10×7, that fits the faces plus one cell for the back. A short page stays a small grid. `sheet/grid.Size` has the same rule; the two copies must stay identical.

**Cell size.** It comes from the deck's first face, and every page of that deck keeps it. TTS limits a sheet to 10,000 pixels per side, so:

1. Start at the face's own size (factor 1).
2. If 10 faces side by side would exceed 10,000 px, shrink the factor so they fit.
3. If 7 faces stacked would exceed 10,000 px, shrink it so they fit.
4. Cell = `trunc(side × factor) / scale`.

**Names.**

| File | Name |
|---|---|
| Sheet | `<deck>_<page>_<hash>.jpg` |
| Back | `backside_<deck>_<hash>.png` (hash of the back bytes) |

The sheet hash is xxHash64 over, in order: `sheetVersion`, cell width, cell height, columns, rows, shadow flag, face count, the back bytes, then each face's bytes. Every image is prefixed with its length. So the name changes exactly when the page would look different, and a swapped card order gives a new name. Card text is not drawn and not hashed.

`sheetVersion` is a constant in `layout.go`. **Bump it whenever sheet drawing changes pixels** (resize filter, paint, shadow, encoder, quality); otherwise pages drawn by older code are reused. See [ADR 023](../../../docs/architecture/023-per-game-results-and-page-reuse.md).

## `script`: the TTS object

`Build` walks the decks again with the same 69-face split, so page numbers and slots match `layout`. It finds each page by `deckID + "_" + page`.

**Object tree.**

```text
Bag  (game name)
└── Bag  (collection id), one per collection that has cards, by collection id
    ├── Deck  (deck name), when the deck holds two or more card objects
    │   └── Card × count, for each card
    └── Card  (card name), when the deck holds exactly one card object
```

TTS cannot spawn a deck of one, so a one-object deck is stored as that card. A single card with count 2 is two objects, so it stays a deck.

**Fields.**

| TTS field | Source |
|---|---|
| `CustomDeck` key | Page number + `deckIDOffset`. The offset grows by each deck's page count, so keys are unique across the game. |
| `FaceURL`, `BackURL` | `file:///` + the absolute paths from `layout`. |
| `NumWidth`, `NumHeight` | Columns and rows from `layout`. |
| Card `Nickname`, `Description` | Catalog name and description. `Name` is always `"Card"`. |
| `LuaScript` | One `key="value"` line per card variable, sorted by key. |
| `CardID` | `CustomDeck key × 100 + slot`, slot counted from 0 on that page. |
| `GUID` | A 6-digit counter: it steps once per deck, once per page break, and once per card. |
| Deck and card transform | Settings `CardSize`. Bags use scale 1. |
| Game bag `Description` | `Created at: <time>`. |

**Order is stable.** Collection bags follow collection id, decks follow deck name, cards follow the catalog list, and `LuaScript` lines follow variable name. An unchanged game gives the same JSON on every render, except the `Created at:` time.

## `fake`

An in-memory game, collections, decks, cards, settings, and TTS for tests. Production code does not use it.
