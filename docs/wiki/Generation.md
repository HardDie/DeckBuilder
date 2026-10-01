# `internal/render/deprecated/generator`

Developer reference. The map of the project is [Internals](Internals). Sheet layout for one page is [Page drawer](Page-drawer). The product steps are in the [user guide](Guide).

`GenerateGame` turns one game into files under `cfg.Results()` (`<data>/result`). It writes a PNG copy of each deck back, a JPEG face sheet for each page, and one TTS Saved Object JSON. It then asks `tts` to spawn that object.

The host list for those images later is [ADR 016](https://github.com/HardDie/DeckBuilder/blob/master/docs/architecture/016-auto-upload-generated-images.md). The URLs written today are still local paths: [ADR 005](https://github.com/HardDie/DeckBuilder/blob/master/docs/architecture/005-tts-generate-local-paths.md).

## Who starts it

The window calls `bindings/generator.Game`. Scale below 1 becomes 1 in the binding, then `GenerateGame` runs.

`GenerateGame` returns as soon as the game exists, the cards are listed, and `result/` has been recreated. Drawing continues in a goroutine. The window polls system `Status` about twice a second.

`progress` is one value for the whole process (`empty`, `in_progress`, `done`, `error`). A second render started while the first is running will clobber it. When the window reads `done` or `error`, that read clears the status back to `empty`.

## Steps

1. `GetSettings`. Failure returns to the caller. The log line is `can't get config`.
2. `game.Item`. An unknown game returns before any folder change.
3. `getListOfCards` walks collections, decks, and cards.
4. `fs.RemoveFolder` on `cfg.Results()`, then `fs.CreateFolder`. The previous render is gone.
5. Progress type becomes `Image generation` and status becomes `in_progress`.
6. The goroutine calls `generateBody`. On error, status becomes `error` and the log line starts with `Generator:`. On success, status becomes `done`.

`generateBody` sets the message `Reading a list of cards from the disk...`, then `generateImages`, then `generateJson`.

## Listing cards

`getListOfCards(gameID, sortField)` calls collection `List`, then deck `List`, then card `List`. `sortField` is the request's `SortOrder`. The search string is empty.

Each list row becomes:

| Record | Fields |
|---|---|
| `Deck` | deck `ID`, `Name`, `Image` |
| `Card` | card `ID`, `GameID`, `CollectionID`, `Count` |

A deck with no cards is omitted. The same `ID`, `Name`, and `Image` in two collections is one map key, and the cards are appended in walk order.

The returned order is deck `Name`, ascending, stable. Collection order and card order stay whatever `List` returned.

## Image pass

`generateImages` walks that deck order. Progress counts catalog cards, not `Count` copies.

For each deck it builds `internal/render/sheet/page.Page` with the deck id, `cfg.Results()`, the request scale, a `commonIndex`, and settings. `EnableBackShadow` and the request scale affect the pixels. Settings `CardSize` does not.

For each card:

1. On an empty page, `deck.GetImage` loads the back. `SetBacksideImageAndSave` writes the original bytes, then keeps a shaded copy for the sheet. Shadow on darkens by 30 brightness. Shadow off still converts the image. The file on disk stays the original bytes.
2. A full page (69 faces) is saved, then `Inherit` starts the next page with the same back and cell.
3. `card.GetImage` loads the face. `AddImage` decodes it. The first face sets the cell. Later faces are resized to that cell with Lanczos. A 10-wide or 7-tall page that would pass 10,000 pixels is scaled down. Request scale divides the cell. `0` becomes `1` inside that math.

A missing back or face returns that error. The log names the game, collection, deck, and card id.

`saveSheet` asks the page for an RGBA image, encodes it with libjpeg-turbo at quality 80 (`internal/render/sheet/draw/libjpeg`), and writes the file. Columns and rows come from the smallest grid that holds the faces plus the back, from 2×2 through 10×7.

The sheet file name:

```text
<commonIndex>_<deckId>_<page>_<faceCount>_<columns>x<rows>.jpg
```

The back file name:

```text
backside_<deckId>_<6 hex chars>.png
```

The hex is the first three MD5 bytes of the original back, printed with `%x`.

`commonIndex` on this pass starts at 0, increments once per deck, and increments again when a full page is saved. It does not increment per card.

The map returned to the JSON pass is keyed by `deckID + "_" + pageIndex` (the page index inside the deck, starting at 1). Each value is `PageInfo`: absolute sheet path, absolute back path, columns, rows.

An empty page writes nothing. The image pass does not build TTS objects.

## JSON pass

`generateJson` walks the same decks and cards again. It does not redraw. It pushes a 10×10 dummy JPEG through `internal/render/deprecated/page_drawer` so the page index and the slot stay aligned with the files already written. `New` uses an empty directory and scale `1`. It never calls `Save`.

The root file is a Saved Object: `ObjectStates` with one bag, the game. `Nickname` is the game name. Inside it, one bag per collection id that contributed a card. The collection bag's `Nickname` is the collection id. Inside a collection bag, the decks and single cards, in deck-name order.

A deck object carries `CustomDeck`. The map key is `pageIndex + deckIdOffset`. `deckIdOffset` starts at 0 and, after each deck, increases by that deck's last page index. The value is `FaceURL`, `BackURL`, `NumWidth`, `NumHeight` from `PageInfo`. Both URLs are `file:///` plus the absolute path. `BackIsHidden` and `UniqueBack` stay false. `Type` stays 0.

Each catalog card becomes a TTS card:

| Catalog | TTS JSON |
|---|---|
| Name | `Nickname`. `Name` is the string `"Card"`. |
| Description | `Description` |
| Variables | `LuaScript`, one `key="value"` line per entry. Go map order. |
| Which cell | `CardID = pageKey*100 + slot`. Slot is 0-based, `Size()-1` after that face was added. |
| How many | `Count` copies appended to the deck, from `card.Item`, not from the list row. |
| GUID | `%06d` of the JSON `commonIndex`, then that counter increments. |

Card and deck transforms use settings `CardSize`. Bags use scale 1. Positions stay 0.

A stack with one contained object is written as that card. TTS will not treat a one-card stack as a deck. Two or more contained objects are a deck object (`Name` `"Deck"`, `Nickname` the deck name). `Count` 2 on a single catalog card is two contained objects, so it stays a deck.

The JSON `commonIndex` is a different counter from the image file name. It starts at 0, increments once per deck, increments when a page fills, and increments again for every card GUID. A deck of three cards therefore consumes more numbers than the image pass, and the next deck's sheet prefix is not this counter.

`result/<gameId>.json` is the `ObjectStates` wrapper, tab-indented, for `Saves/Saved Objects`. The game bag description is `Created at: ` plus `2006-01-02 15:04:05` at the moment of the write.

`SendToTTS` stores the inner bag and dials `127.0.0.1:39999`. The script is External Editor `messageID` 3 to Global (`guid` `-1`): `WebRequest.get` of `http://127.0.0.1:<port>/api/tts/data`, then `spawnObjectJSON`. If nothing is listening, render still succeeds. The log says the dial failed. The HTTP handler returns the bag once and clears it.

## What render does not do

It does not upload the sheets. `result/<gameId>.json` keeps `file:///` paths after you host the pictures somewhere else. Replace is a separate step: it rewrites those two fields from a mapping file. Steam Cloud upload is a TTS action, described in [Steam Cloud](Steam-Cloud).
