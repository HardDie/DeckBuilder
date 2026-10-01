# `internal/render/generate`

Plans one game render. It decides which cards go on which sheet, how many copies each card has, the grid, the file names, the absolute paths, and the TTS fields.

`Prepare` returns that plan. It reads catalog bytes and image headers. Pixel work happens later, outside this package.

`internal/render/compose` is the caller that sits outside both this package and `internal/render/sheet`. It asks `Prepare` for the plan, writes the raw back files and the JSON, and hands each sheet to `internal/render/sheet/write.Draw`.

## Entry point

```go
func Prepare(
    dir string,
    gameItem *entitiesGame.Game,
    decks map[catalog.Deck][]catalog.Card,
    order []catalog.Deck,
    scale int,
    cfg *entitiesSettings.Settings,
    deckSvc servicesDeck.Deck,
    cardSvc servicesCard.Card,
) (Plan, error)
```

| Argument | Role |
|---|---|
| `dir` | Directory for the result files. Usually `cfg.Results()`. |
| `gameItem` | Game id becomes the JSON file name. Game name becomes the TTS bag nickname. |
| `decks` / `order` | Cards grouped by deck, and the deck order. `catalog.Collect` builds both. |
| `scale` | Integer divisor of the cell. `0` becomes `1` in the cell math. |
| `cfg` | `EnableBackShadow` is copied onto each sheet as a flag. `CardSize` is the TTS transform scale. |
| `deckSvc` / `cardSvc` | `GetImage` for raw bytes. `Item` for the card name, description, variables, and count. |

`Prepare` loads the back from the first card's collection, then each face in order. A missing back or face returns that error. The log names the game, collection, deck, and card id.

## What comes out

`Plan`:

| Field | Contents |
|---|---|
| `Sheets` | One page ready to paint: original face bytes, original back bytes, `CellW`, `CellH`, `Shadow`, and the absolute JPEG path. |
| `Backs` | One raw PNG per deck: absolute path and the catalog bytes, unchanged. |
| `JSONPath` | `<dir>/<game id>.json`. |
| `Root` | Saved Object wrapper, `ObjectStates` with one game bag. |
| `Bag` | That same game bag, including the `Created at:` description, for TTS. |

`Prepare` does not write these files. The caller does.

## Subpackages

### `catalog`

`Collect` walks collection `List`, then deck `List`, then card `List`. The sort field is the caller's `SortOrder`. The search string is empty.

| Record | Fields |
|---|---|
| `Deck` | `ID`, `Name`, `Image` |
| `Card` | `ID`, `GameID`, `CollectionID`, `Count` |

A deck with no cards is omitted. The same `ID`, `Name`, and `Image` in two collections is one map key. Cards are appended in walk order.

The returned order is deck `Name`, ascending, stable. Collection order and card order stay whatever `List` returned.

### `layout`

`Pages` measures sheets. It does not open pixels.

For each deck that has faces:

1. `common` starts at 0 and increments once per deck.
2. The cell comes from the first face. Width and height are read from a PNG `IHDR` or a JPEG SOF marker. A 10-wide page wider than 10,000 pixels, or a 7-tall page taller than 10,000 pixels, shrinks `inner`. Scale `0` becomes `1`. `inner` `0` becomes `1`. The cell is `trunc(side * inner) / scale`. Later pages of that deck keep this cell.
3. The back file name is `backside_<deck id>_<6 hex chars>.png`. The hex is the first three MD5 bytes of the raw back, printed with `%x`. The path is absolute. The body is the original bytes.
4. Faces fill a page until `config.MaxCount` (69). The next face starts another page. `common` and the page index both increment on that split. The page index inside the deck starts at 1.

The sheet file name:

```text
<common>_<deck id>_<page>_<face count>_<columns>x<rows>.jpg
```

Columns and rows are the smallest grid that can hold the faces plus one cell for the back, from 2×2 through 10×7. A short page stays a tight grid. 69 faces plus the back is 10×7.

`Shadow` on the page is the settings flag, copied through.

### `script`

`Build` turns the measured pages and the catalog cards into TTS objects. Page breaks use the same 69-face counter, so the page index and the slot match `layout`.

The lookup key is `deckID + "_" + pageIndex`.

| TTS field | Source |
|---|---|
| Game bag `Nickname` | Game name. Bag scale is 1. |
| Collection bag `Nickname` | Collection id. One bag per collection that contributed a card. |
| Deck `Nickname` | Deck name. Transform scale is settings `CardSize`. |
| `CustomDeck` key | `pageIndex + deckIdOffset`. After each deck, `deckIdOffset` increases by that deck's last page index. |
| `FaceURL`, `BackURL` | `file:///` plus the absolute path from `layout`. |
| `NumWidth`, `NumHeight` | Columns and rows from `layout`. `BackIsHidden`, `UniqueBack`, and `Type` stay zero. |
| Card `Name` | The string `"Card"`. `Nickname` is the catalog name. |
| `Description` | Catalog description. |
| `LuaScript` | One `key="value"` line per variable. Go map order. |
| `CardID` | `(pageIndex + deckIdOffset) * 100 + slot`. Slot is 0-based on that page. |
| `GUID` | `%06d` of this pass's `commonIndex`, then that counter increments. |
| Copies | `card.Item` `Count`, appended one by one. |

This `commonIndex` is a different counter from the number in the JPEG file name. It increments once per deck, again when a page fills, and again for every card GUID.

A stack with one contained object is placed as that card. Two or more contained objects stay a deck. `Count` 2 on a single catalog card is two contained objects, so it stays a deck.

The game bag description is `Created at: ` plus the clock, formatted `2006-01-02 15:04:05`.

`fake` is an in-memory catalog for tests. Production code does not call it.
