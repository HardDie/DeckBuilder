# Generation

What **Render** does to a game, from the binding to the files in `result/`. The product steps are in the [user guide](Guide). The host list for those images later is [ADR 016](https://github.com/HardDie/DeckBuilder/blob/master/docs/architecture/016-auto-upload-generated-images.md). The current rule is still local paths: [ADR 005](https://github.com/HardDie/DeckBuilder/blob/master/docs/architecture/005-tts-generate-local-paths.md).

## Who starts it

The window calls the generator binding. `GenerateGame` returns immediately and runs the work in a goroutine. The window polls system `Status` about twice a second and shows the progress circle.

`progress` is one value for the whole process (`empty`, `in_progress`, `done`, `error`). A second render started while the first is running will clobber it. When the window reads `done` or `error`, that read clears the status back to `empty`.

The job loads settings (card scale, back shadow), walks every collection and deck, and fails if a deck has no back image or a card has no face.

## Two passes

`generateBody` does images first, then JSON.

1. `generateImages` reads face and back bytes and writes the sheet files. It returns a map keyed by deck id plus page index. Each entry has the sheet path, the back path, and the grid size.
2. `generateJson` walks the same decks again and builds TTS objects. It does not redraw. For the grid math it pushes a dummy image per card so the page index and the card slot stay aligned with the files already written.

## How a sheet is packed

`internal/page_drawer` owns one page. TTS custom decks default to 10 columns by 7 rows. With `BackIsHidden` left false, TTS treats the **last cell of the face sheet** as the hidden card. That cell is not a playable face. `MaxCount` is `10*7 - 1`, so 69 faces per page. A 70th card starts another page.

For each deck:

1. The first card loads the deck's image and saves it as `backside_<deck>_<hash>.png`. If back shadow is on, the copy drawn into the sheet is darkened by 30 brightness. The file on disk stays the original bytes.
2. Each face is decoded. The first face sets the cell size. Later faces are resized to that cell with Lanczos. If a 10-wide or 7-tall page would pass 10,000 pixels, the cell is scaled down so the page stays inside that limit.
3. When the page holds 69 faces, `Save` runs and a new page inherits the back and the cell size.
4. After the last card, the open page is saved if it has any faces.

`Save` asks `CalculateGridSize` for the smallest grid that can hold the faces **plus one** (the back), inside 2×2 up to 10×7. It prefers a tighter grid over a full 10×7 when the page is not full. Cells fill left to right, top to bottom. The back is drawn in the last cell (`columns-1`, `rows-1`).

The sheet file name looks like:

```text
<commonIndex>_<deckId>_<page>_<faceCount>_<columns>x<rows>.jpg
```

`commonIndex` counts pages across the whole game, so two decks do not write the same name. The sheet is JPEG. The standalone back is PNG.

## How the JSON is built

The root file is a Saved Object: `ObjectStates` with one bag, the game. Inside it, one bag per collection. Inside a collection bag, the decks.

A deck object carries `CustomDeck`. The map key is the page index. The value is `FaceURL`, `BackURL`, `NumWidth`, `NumHeight`. Both URLs are `file:///` plus the absolute path `generateImages` returned.

Each catalog card becomes a TTS card:

| Catalog | TTS JSON |
|---|---|
| Name | `Nickname`. `Name` is the string `"Card"`. |
| Description | `Description` |
| Variables | `LuaScript`, one `key="value"` line per entry |
| Which cell | `CardID = pageId*100 + indexOnThatPage` |
| How many | `count` copies appended to the deck |

A deck with a single card is written as that card object. TTS will not treat a one-card stack as a deck. Two or more cards are a deck object.

Card scale from settings is the TTS transform scale on the bag, the decks, and the cards.

The file written to `result/<gameId>.json` is the `ObjectStates` wrapper, for `Saves/Saved Objects`. `SendToTTS` sends the inner bag only. `spawnObjectJSON` wants one object, not the wrapper.

## Spawn

`SendToTTS` stores the bag, dials `127.0.0.1:39999`, and sends External Editor message `messageID` 3 to Global (`guid` `-1`). The script is `WebRequest.get` of `http://127.0.0.1:<port>/api/tts/data`, then `spawnObjectJSON`. If nothing is listening, render still succeeds. The log says the dial failed.

The HTTP handler returns the bag once and clears it.

## What render does not do

It does not upload the sheets. `result/<gameId>.json` keeps `file:///` paths after you host the pictures somewhere else. Replace is a separate step: it rewrites those two fields from a mapping file. Steam Cloud upload is a TTS action, described in [Steam Cloud](Steam-Cloud).
