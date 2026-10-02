# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

DeckBuilder is a **local authoring tool** for card games.
It exports to [Tabletop Simulator](https://tabletopsimulator.com/) (TTS).

1. UI
   1. Wails window at the repo root.
   2. Vue sources live in `frontend/`.
   3. Catalog verbs are Wails bindings ([ADR 011](docs/architecture/011-wails-bindings.md)).
2. HTTP ([ADR 013](docs/architecture/013-http-images-and-tts.md))
   1. Loopback `127.0.0.1:5000`.
   2. Entity image bytes.
   3. `GET /api/tts/data` for TTS Lua.

Not this product:
1. A multi-user cloud service.
2. A TTS plugin inside the game.
3. More than a one-shot External Editor TCP poke.

Keep this file updated when routes, disk layout, generate rules, or layers change.

## Commands

All Go builds and tests need cgo and a **static libjpeg**.
On macOS: `brew install jpeg-turbo`.
The Makefile sets `CGO_*`, `CC`, and `LIBRARY_PATH`.
Prefer `make` targets over bare `go` commands.

```bash
make dev               # Wails window with frontend hot reload
make build             # production binary -> build/bin
make test              # unit tests, -race -tags=nomain (root, bindings, internal, pkg)
make test-integration  # go test -tags=integration ./internal/...
make test-all          # what CI runs (.github/workflows/test.yml)
make vet
make fmt               # gofmt check; fails and lists files that need formatting
make linter-install && make linter-run   # gosec + golangci-lint (.golangci.yaml)
make generate          # regenerate frontend/wailsjs after changing bindings/ signatures
make fuzz_game         # also fuzz_collection, fuzz_deck, fuzz_card
make help
```

Running a single test:
1. Always pass `-tags=nomain`.
   1. `main.go` and `wire.go` are `//go:build !nomain`.
   2. They need Wails cgo.
2. Packages that reach `internal/render/sheet/draw/libjpeg` also need the cgo flags.

On macOS:

```bash
CGO_ENABLED=1 CGO_CFLAGS="-I$(brew --prefix jpeg-turbo)/include" \
CGO_LDFLAGS="$(brew --prefix jpeg-turbo)/lib/libjpeg.a" \
go test -race -count=1 -tags=nomain -run TestName ./internal/services/game/
```

Tests:
1. Unit tests sit next to the package.
2. Prefer table-driven tests.
3. Integration tests use `-tags=integration`.
4. Render goldens live in `internal/render/compose/testdata/golden`.
5. Regenerate goldens with `WRITE_GOLDEN=1`.
6. Tests use `config.SetDataPath` with a temp dir.
7. Never write into the developer's `DeckBuilderData`.

Frontend:
1. `frontend/` is Vue 3 + Vite + Pinia + naive-ui, managed with yarn.
2. Install with `make frontend-install`.
3. Lint with `yarn --cwd frontend lint`.

## Product scope

Must have:
1. Catalog tree: **game → collection → deck → card**.
2. Images per entity, from a file upload or a URL.
   1. GUI and TTS load them via `/api/…/image`.
3. Cards carry Lua variables (`map[string]string`).
4. Cards carry a `count`: copies in the rendered deck.
5. Render a game.
   1. Output: sprite-sheet images plus one TTS Saved Object JSON.
   2. Output dir: `DeckBuilderData/result`.
6. Export and import a game as a zip.
7. Duplicate a game.
8. Recursive search across the tree.
9. Replace local `file://` image paths in generated JSON with hosted URLs.
10. Optional TTS spawn after generate, via the External Editor port.
11. Settings persist in the fsentry store.
    1. Language `en` or `ru`.
    2. Card scale.
    3. Back shadow.

Known limitation. Do not "fix" it by inventing a host.
1. Generated JSON points at local image paths.
2. Users keep files in `result/`, or run replace after uploading.
3. Auto-upload is out of scope unless explicitly requested.
4. Future hosts: [ADR 016](docs/architecture/016-auto-upload-generated-images.md).
   1. Not implemented.
   2. Steam Cloud via the Lua poke is rejected there.

Out of scope unless explicitly requested:
1. Binding on non-loopback interfaces.
2. Adding auth.
3. Changing the domain tree.
4. Replacing fsentry with SQL.
5. Serving the Vue UI over HTTP.

## Domain model

| Level | ID type | Role (example) |
|---|---|---|
| Game | string (folder id from name) | "Munchkin" |
| Collection | string | Base game or a DLC |
| Deck | string | "Monster"; **deck image is the card back** |
| Card | `int64` | Face image, description, variables (`HP: 2`), `count` |

1. List calls accept `sort` and `search`.
   1. Sort values: `name`, `name_desc`, `created`, `created_desc`.
2. Create and update take `name`, `description`, optional `image` URL, optional `imageFile`.
3. Card writes also take `variables` (JSON string) and `count`.
4. JSON envelope: `{ "data": …, "meta": { "total", "cardsTotal?" }, "error": … }`.
5. String ids come from fsentry and derive from the name.
6. Card ids are integers.
   1. Path params named `{card}` are parsed as that integer.
7. A missing deck back fails generate.

## Architecture

Process:
1. `main.go` is the only entrypoint.
2. `wire.go` builds the config, the fsentry store, repositories, and services.
3. `main.go` starts the loopback HTTP server (`internal/servers`).
4. `main.go` then opens the Wails window.
5. `api_proxy.go` sends the Wails AssetServer's `/api` requests to the same mux.

UI calls:
1. Catalog verbs are Wails bindings in `bindings/<area>`, not REST.
   1. Game (with export and import), collection, deck, card.
   2. System, search, generator, replace.
2. One binding struct per pane.
3. Bindings map DTOs from `internal/dto` to services.
4. DTO timestamps are RFC3339 strings.
   1. Wails bindings cannot resolve `time.Time`.
   2. Entities stay `time.Time`.
5. Generated JS lives in `frontend/wailsjs` (`make generate`).
6. If the GUI needs a field, add it to the existing DTO with `json` tags.

Layers:
1. Imports go `main/bindings/servers` → `services` → `repositories` → fsentry.
2. `servers` register routes and parse HTTP.
   1. They map to services and `dto`.
   2. They never touch fsentry.
   3. They hold no business rules.
3. `services` own rules: validation, generate orchestration, replace.
   1. They do not import `net/http`.
4. `repositories` own disk: fsentry folders, image files, config paths.
5. Keep three struct families apart:
   1. Domain structs: `internal/entities`.
   2. GUI JSON: `internal/dto`.
   3. TTS JSON: `internal/tts_entity`.
   4. TTS `CardID` differs from the catalog `int64` id.
6. Each package declares its interface in `contract.go`.
7. Implementation structs stay unexported.
8. Errors come from `internal/errors`.
9. Do not add a top-level `internal/` package for a single handler.
10. Do not collapse `servers` / `services` / `repositories` without an ADR.

Other locations:
1. Paths and sheet limits: `internal/config`.
2. JSON envelope: `internal/network`.
3. Version: `pkg/version`.
4. One-off CLIs: `tools/`.
5. `web/` is a dist stub. Do not put the Vue UI there.

## HTTP contract

1. Listen on `127.0.0.1:5000`.
   1. If bind fails, try `:5001`, `:5002`, and so on.
   2. Stop after 20 attempts and fail startup.
   3. Never bind on all interfaces.
2. CORS is `*` with `GET,OPTIONS`.
   1. Preflight `OPTIONS` returns 204.
3. Do not invent REST aliases.

| Method | Route | Role |
|---|---|---|
| GET | `/api/games/{game}/image` (and nested `…/image`) | Entity image bytes |
| GET | `/api/tts/data` | One-shot last generated JSON for TTS Lua |

## Render pipeline

1. The window calls `bindings/generator.Game`.
   1. A scale below 1 becomes 1.
2. `compose` (`internal/render/compose`) is the live generator.
   1. `wire.go` builds it as `compose.Generator`.
3. `GenerateGame` lists cards, clears `result/`, and returns.
   1. Never block the binding on image drawing.
4. A goroutine then runs three steps.
   1. `generate.Prepare` plans the catalog walk, layout, and TTS JSON.
   2. `sheet/write.Draw` draws JPEG sheets with libjpeg-turbo.
   3. `services/tts.SendToTTS` pokes TTS, best effort.
5. Sheets follow TTS Custom Deck defaults.
   1. Grid is at most 10 wide × 7 high.
   2. One page holds 69 faces (`MaxCount`).
   3. `BackIsHidden` is false, so the last cell is the hidden image.
   4. That cell holds the back.
   5. A dedicated back file is still written for `BackURL`.
   6. Overflow starts a new page.
6. Card scale and back shadow come from settings.
7. Progress (`internal/render/progress`) is a process-wide singleton.
   1. The GUI polls system `Status`.
   2. Values: `empty`, `in_progress`, `done`, `error`.
   3. Reading `done` or `error` resets it to `empty`.
   4. Pollers treat that `empty` as "already observed".
   5. A second generate clobbers it.
   6. Do not allow overlapping generates without a product decision.
8. `internal/render/deprecated/*` is not called by compose.
   1. It holds the old `generator`, `page_drawer`, and `progress`.
   2. `docs/wiki/Generation.md` still describes it.
   3. Trust the code and `internal/render/compose/README.md`.

## Tabletop Simulator

Source of truth: [api.tabletopsimulator.com](https://api.tabletopsimulator.com/).
Do not invent Lua globals or TCP message ids.

Custom deck fields:
1. Lua `setCustomObject` uses `face`, `back`, `width`, `height`, `number`.
   1. Also `unique_back`, `sideways`, `back_is_hidden`.
2. Saved Object JSON uses `FaceURL`, `BackURL`, `NumWidth`, `NumHeight`.
   1. Also `BackIsHidden`, `UniqueBack`.
   2. See `internal/tts_entity.DeckDescription`.
3. Our saved objects use `Name: "Deck"`, `"Card"`, and `"Bag"`.
4. A deck with one card must be a card object, not a deck.
   1. Generate already splits that case.

Spawn flow:
1. TTS External Editor listens on TCP `39999`.
   1. We connect there and send once.
   2. We do not listen on `39998`.
2. Message: Execute Lua Code.
   1. `messageID` is 3.
   2. `guid` is `-1` (Global script).
3. The Lua calls `WebRequest.get("http://127.0.0.1:5000/api/tts/data")`.
   1. Check `request.is_error` before `request.text`.
4. Then it calls `spawnObjectJSON({ json = request.text, … })`.
5. `/api/tts/data` serves the root bag once, then clears the buffer.
6. If TTS is not running, generate still writes `result/`.

File formats:
1. `result/*.json` is a Saved Object with an `ObjectStates` array.
   1. Users copy it to `Tabletop Simulator/Saves/Saved Objects`.
2. `spawnObjectJSON` wants one object, not that wrapper.
3. `file:///` URLs work on one machine only.
   1. Sharing a table needs hosted HTTP(S) URLs.
   2. See [ADR 005](docs/architecture/005-tts-generate-local-paths.md).

## Data on disk

1. Storage is [fsentry](https://github.com/HardDie/fsentry): folders plus JSON sidecars.
   1. No SQL, no migrations, no other config file.
2. Data root (`cfg.Data`):
   1. Linux and Windows: `./DeckBuilderData` in the working directory.
   2. macOS: `~/DeckBuilderData`, since a `.app` dir is not writable.
3. Layout: `Data/games/<game>/…` for the catalog and images.
4. `Data/result/` holds the last generate.
5. Settings live in the same fsentry tree.
6. Catalog `createdAt` and `updatedAt` are written on every create and update.
   1. Empty values in old data are filled in memory.
   2. See [ADR 009](docs/architecture/009-catalog-timestamps.md).

## Version

1. Package `pkg/version`.
2. Stamp `-X github.com/HardDie/DeckBuilder/pkg/version.Build`.
3. Value is an exact git tag, or the 12-character commit hash.
4. A blank stamp is `dev`.

## Documentation

1. [README.md](README.md) is for users.
   1. What the app is, build, TTS copy path, honest limitations.
   2. Keep it short.
2. This file is for agents and contributors.
3. Decisions: [docs/architecture](docs/architecture/INDEX.md).
   1. New core decisions get an ADR with the next number.
   2. Add a row to `INDEX.md`.
4. Scenarios: [docs/use-cases](docs/use-cases/INDEX.md).
   1. New user-visible behavior gets a use case from `_TEMPLATE.md`.
5. Package notes: [docs/wiki](docs/wiki/Home.md).
   1. GitHub-wiki markdown, slug links without `.md`.
   2. Publish by copying into `DeckBuilder.wiki.git`.

## Working agreements

1. Prefer the smallest change that keeps the GUI contract.
2. Avoid layer rewrites.
