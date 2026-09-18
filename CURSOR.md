# DeckBuilder Development Guidelines

DeckBuilder is a **local authoring tool** for card games that export to [Tabletop Simulator](https://tabletopsimulator.com/) (TTS). This repository is the **Go HTTP backend**. The GUI lives in a separate repo, cloned as the `gui` git submodule ([HardDie/DeckBuilderGUI](https://github.com/HardDie/DeckBuilderGUI); README also points at [lmm1ng/DeckBuilderGUI](https://github.com/lmm1ng/DeckBuilderGUI)). Built frontend assets are copied into `web/dist` and **embedded** into the binary.

**Two surfaces**

1. **SPA (embedded `web/dist`)** — the authoring UI. It talks only to this process over HTTP. SPA routes such as `/games/{id}/…` are SPA paths; Go forwards them to `index.html` so reloads work.
2. **REST API (`/api/…`)** — JSON (and multipart) for CRUD, images, generate, search, replace, system. TTS Lua downloads the last generated object from `GET /api/tts/data`.

**Not this product:** a multi-user cloud service, a TTS plugin inside the game (beyond a one-shot External Editor TCP poke), or a second in-process desktop toolkit such as Wails. The process is a **loopback HTTP server** on `127.0.0.1:5000`.

---

## Product scope

**Must have (current product)**

- Hierarchical catalog: **game → collection → deck → card**.
- Images per entity (file upload or URL); GUI and TTS consume them via `/api/…/image`.
- Cards carry **Lua variables** (`map[string]string`) and a **count** (how many copies go into the rendered deck).
- **Render/generate** a game: sprite-sheet PNGs + one TTS Saved Object JSON under `DeckBuilderData/result` (on macOS, that data dir is under the home directory).
- **Export/import** a game as a zip; **duplicate** a game.
- Recursive **search** across the tree.
- **Replace** local `file://` (or disk) image paths in generated JSON with hosted URLs so a table can be shared.
- Optional **TTS spawn**: after generate, try TCP `127.0.0.1:39999` (TTS External Editor) and tell the game to `WebRequest.get` `/api/tts/data`.
- Settings (language `en`/`ru`, card scale, back shadow) persisted in the fsentry store; **quit** from the UI with a delay (skipped in `-debug`).

**Known limitation (do not “fix” by inventing a host)**

- Generated JSON currently points at **local image paths**. Users must keep files in `result/` (or run **replace** after uploading images). Automatic upload to image hosting is **out of scope unless explicitly requested**.

**Out of scope unless explicitly requested**

- Binding on non-loopback interfaces or adding auth.
- Changing the domain tree (game/collection/deck/card).
- Replacing fsentry with SQL.
- Merging the GUI into this Go module as a Wails/webview app.

---

## Domain model

| Level | ID type | Role (example) |
|---|---|---|
| Game | string (folder id from name) | “Munchkin” |
| Collection | string | Base game or a DLC |
| Deck | string | “Monster”; **deck image is the card back** |
| Card | `int64` | Face image, description, variables (`HP: 2`), `count` |

List endpoints accept `sort` (`name`, `name_desc`, `created`, `created_desc`) and `search`. Create/update of named entities use **multipart/form-data** (`name`, `description`, optional `image` URL, optional `imageFile`). Cards also send `variables` (JSON string) and `count`.

JSON envelope: `{ "data": …, "meta": { "total", "cardsTotal?" }, "error": … }`.

---

## HTTP contract

Listen: **`127.0.0.1:5000`**, then `:5001` … if bind fails, up to **20** ports. Fail startup if all 20 fail. CORS is `*` with `GET,POST,PATCH,DELETE` (the SPA may be served from the same origin when embedded; CORS exists for local GUI-dev against the API).

| Method | Route | Role |
|---|---|---|
| GET/POST | `/api/games` | List / create games |
| GET/PATCH/DELETE | `/api/games/{game}` | Item / update / delete |
| POST | `/api/games/{game}/duplicate` | Copy game (JSON body `{ "name" }`) |
| GET | `/api/games/{game}/export` | Zip archive |
| POST | `/api/games/import` | Import zip (`file`, optional `name`) |
| GET/POST | `/api/games/{game}/collections` | List / create |
| GET/PATCH/DELETE | `/api/games/{game}/collections/{collection}` | Item / update / delete |
| GET | `/api/games/{game}/decks` | All decks in a game |
| GET/POST | `/api/games/{game}/collections/{collection}/decks` | List / create |
| GET/PATCH/DELETE | `…/decks/{deck}` | Item / update / delete |
| GET/POST | `…/decks/{deck}/cards` | List / create |
| GET/PATCH/DELETE | `…/cards/{card}` | Item / update / delete |
| GET | `/api/games/{game}/image` (and nested `…/image`) | Entity image bytes |
| POST | `/api/games/{game}/generate` | Start background render (`sortOrder`, `scale`) |
| GET | `/api/search`, `/api/search/games/{game}`, `…/collections/{collection}` | Recursive search |
| POST | `/api/replace/prepare` | Extract unique FaceURL/BackURL keys |
| POST | `/api/replace` | Rewrite JSON with mapping file |
| GET | `/api/tts/data` | One-shot last generated JSON for TTS Lua |
| DELETE | `/api/system/quit` | Start 60s quit timer (no-op in `-debug`) |
| GET/PATCH | `/api/system/settings` | Settings (`lang` on PATCH today) |
| GET | `/api/system/status` | Generator progress; **clears** `done`/`error` after read |
| GET | `/api/system/version` | Stamped version string |
| GET | `/docs`, `/swagger.json` | Redoc + spec |
| GET | `/` and SPA asset paths | Embedded GUI |

Do not invent REST aliases. If the GUI needs a field, add it on the existing DTO with `json` tags and keep `error` / `data` wrapping.

**Generate** returns immediately; work runs in a goroutine. The GUI polls `/api/system/status`. Status values: `empty`, `in_progress`, `done`, `error`.

**Sprite sheets:** pages follow TTS [Custom Deck](https://api.tabletopsimulator.com/custom-game-objects/) defaults: **width 10 × height 7**. One cell is reserved (`MaxCount = 69`) because, with `BackIsHidden` false (our generate default), TTS uses **the last slot of the face sheet as the hidden image**, not a separate back. `page_drawer.Save` draws the back into that last cell and still writes a dedicated back file for `BackURL`. Overflow starts a new page. Card scale / back shadow come from settings. Output: images + JSON in `cfg.Results()`.

---

## Tabletop Simulator (official API)

Source of truth: [api.tabletopsimulator.com](https://api.tabletopsimulator.com/). Do not invent Lua globals or TCP message IDs.

**Custom objects (Lua `setCustomObject` vs save JSON).** Lua custom-deck fields are `face`, `back`, `width`, `height`, `number`, `unique_back`, `sideways`, `back_is_hidden`. Saved-object / `spawnObjectJSON` JSON uses `FaceURL`, `BackURL`, `NumWidth`, `NumHeight`, `BackIsHidden`, `UniqueBack` (`internal/tts_entity.DeckDescription`). Spawn type name for a custom deck is `DeckCustom`; our saved decks use `Name: "Deck"` / cards `Name: "Card"` / bags `Name: "Bag"` as in TTS save data. A deck with **one card must be a card object**, not a deck (generator already splits that case).

**`spawnObjectJSON`** ([base](https://api.tabletopsimulator.com/base/)): mandatory `json` string (one object); optional `position` / `rotation` / `scale` / `callback_function`. Prefer this over hand-building Custom Deck in Lua.

**`WebRequest.get(url, callback)`** ([manager](https://api.tabletopsimulator.com/webrequest/manager/)): HTTP from the **host computer only**. Callback receives a request instance; check `request.is_error` / `request.error` before `request.text`.

**External Editor** ([externaleditorapi](https://api.tabletopsimulator.com/externaleditorapi/)): two localhost TCP sockets. **TTS listens on `39999`** (we connect here). The editor would listen on `39998` for replies; this app does **not** listen on 39998 (fire-and-forget). Messages are JSON with `messageID`. We send **Execute Lua Code** (`messageID` **3**) with `guid` **`-1`** (Global script) and a `script` string. Objects only accept Execute Lua if they already have a script in TTS; Global `-1` is the right target for spawn-from-HTTP. After generate, `SendToTTS` marshals the **root bag** (not the `ObjectStates` wrapper). Lua then `WebRequest.get("http://127.0.0.1:5000/api/tts/data")` and `spawnObjectJSON({ json = request.text, … })`. The file in `result/*.json` is a **Saved Object** (`ObjectStates` array) for copying into TTS `Saves/Saved Objects`; that wrapper is not what `spawnObjectJSON` wants.

**`file:///` FaceURL/BackURL:** TTS custom images are “path/URL”. Local paths work for one machine; sharing a table still needs hosted HTTP(S) URLs ([ADR 005](docs/architecture/005-tts-generate-local-paths.md)).

---

## Data on disk

Storage is **[fsentry](https://github.com/HardDie/fsentry)** (folders + JSON sidecar), not a SQL DB.

| OS | Data root (`cfg.Data`) |
|---|---|
| Linux / Windows | `./DeckBuilderData` next to the working directory |
| macOS | `~/DeckBuilderData` (cannot write next to a `.app`) |

Layout (conceptual): `Data/games/<game>/…` collections, decks, cards, images; `Data/result/` last generate; settings in the same fsentry tree. **Never assume SQL, migrations, or a user-config JSON besides this store.** Catalog `createdAt`/`updatedAt` are always written on create/update; preview null/empty values are filled in memory when mapping to entities ([ADR 009](docs/architecture/009-catalog-timestamps.md)).

---

## Stack

| Layer | Choice |
|---|---|
| GUI | Vue SPA in `gui/` submodule; production assets `web/dist` (`//go:embed`) |
| HTTP | Go 1.27.1 (`go.mod`), `gorilla/mux`, listen `127.0.0.1:5000` (try next port up to 20 times if busy) |
| Persistence | `github.com/HardDie/fsentry` |
| Images | `disintegration/imaging`, `internal/page_drawer`, `internal/images` |
| API docs | go-swagger comments in `internal/api`; `make swagger` → `web/swagger.json` |
| Tests | `go test ./... -race`; fuzz targets under services |

---

## Build & Run

```bash
git clone https://github.com/HardDie/DeckBuilder --recursive
./deployment/check_binary.sh
make web-build   # gui: yarn install && yarn build → web/dist
make build       # deployment/build_all.sh → deployment/out
make test
make linter-run  # after make linter-install
```

Local API + embedded UI: build `cmd/deck_builder` and run it. `-debug` skips opening the browser and ignores `/system/quit`. Swagger UI: `http://localhost:5000/docs`.

GUI-only development typically runs the Vue app separately against this API (CORS is open). Keep API contracts stable.

---

## Where things live

This file stays lean. **[README.md](README.md)** is the short user entry (what the app is, build, TTS copy path). This file is the agent/implementation spec. Decisions: **[docs/architecture](docs/architecture/INDEX.md)**. Scenarios: **[docs/use-cases](docs/use-cases/INDEX.md)**. Package walkthrough: **[docs/wiki](docs/wiki/Home.md)** (GitHub wiki source; copy into `DeckBuilder.wiki.git`).

| If you need… | Read |
|---|---|
| Product, clone, build, TTS Saved Objects path | **[README.md](README.md)** |
| TTS Lua, External Editor, Custom Deck | **[api.tabletopsimulator.com](https://api.tabletopsimulator.com/)** — [External Editor](https://api.tabletopsimulator.com/externaleditorapi/), [spawnObjectJSON](https://api.tabletopsimulator.com/base/), [Custom Deck](https://api.tabletopsimulator.com/custom-game-objects/), [WebRequest](https://api.tabletopsimulator.com/webrequest/manager/) |
| Architecture decisions | **[docs/architecture](docs/architecture/INDEX.md)** |
| Use cases | **[docs/use-cases](docs/use-cases/INDEX.md)** |
| Go modules, fields, who calls what | **[docs/wiki](docs/wiki/Home.md)** |
| HTTP routes + swagger comments | `internal/api` |
| Process wiring | `internal/application/application.go` |
| Paths, sheet limits | `internal/config` |
| Domain structs | `internal/entities` |
| JSON DTOs | `internal/dto` |
| TTS Saved Object shapes | `internal/tts_entity` |
| HTTP handlers | `internal/servers` |
| Business logic | `internal/services` |
| Image + zip + fs helpers | `internal/repositories`, `internal/db`, `internal/fs` |
| Sprite composition | `internal/page_drawer` |
| Embedded UI + swagger | `web/` |
| Wails window (scaffold) | `desktop/` |
| One-off CLIs | `tools/` |

### Tree

```text
.
├── Makefile
├── README.md
├── CURSOR.md
├── docs/
│   ├── architecture/       # ADRs
│   ├── use-cases/
│   └── wiki/               # GitHub wiki source (Home.md, _Sidebar.md)
├── cmd/deck_builder/       # production main, swagger:meta, -debug, version ldflags
├── desktop/                # Wails v2 scaffold (own module; not the HTTP server)
├── internal/
│   ├── application/        # mux, DI, ListenAndServe
│   ├── api/                # route registration + swagger types
│   ├── servers/            # HTTP adapters
│   ├── services/           # use-case logic
│   ├── repositories/       # images + persistence glue
│   ├── db/                 # fsentry folders per aggregate
│   ├── entities/           # in-process domain
│   ├── dto/                # JSON to the GUI
│   ├── tts_entity/         # TTS JSON schema
│   ├── page_drawer/        # sheet layout
│   ├── progress/           # singleton generate status
│   ├── config/
│   ├── network/            # JSON envelope, open browser
│   ├── errors/
│   └── …
├── web/                    # embed dist + swagger.json
├── gui/                    # submodule: DeckBuilderGUI
├── deployment/
└── tools/                  # join, copy_cards_variables
```

**Layer rules**

- **`api` registers routes** and holds swagger request/response types. Unimplemented structs exist so `go-swagger` can scan comments.
- **`servers` parse HTTP** (mux vars, multipart) and map to services + `dto`. They must not talk to fsentry directly.
- **`services` own rules** (validation, generate orchestration, replace). They depend on repositories / other services, not `net/http`.
- **`repositories` + `db`** own disk. `db` is the fsentry mapping; repositories add image files and config paths.
- **`entities` vs `dto` vs `tts_entity`:** do not reuse TTS JSON structs as API DTOs. Cards in TTS have different IDs (`CardID` vs catalog `int64`).
- **Import direction:** `application` → servers → services → repositories → db. `page_drawer` and `tts_entity` are used by generator, not by `api`.
- **Do not add a new top-level `internal/` package** for a single handler; put it on an existing server.

### Tests

- Unit tests sit next to the package (`internal/services/game/game_test.go`).
- Prefer table-driven tests. Fuzz: `make fuzz_game` (and collection/deck/card).
- Tests may use `config.SetDataPath` for a temp dir. Do not write into the developer’s `DeckBuilderData`.

---

## Key facts to keep in mind

- **This repo is the backend.** GUI changes belong in DeckBuilderGUI unless you are only embedding a new `web/dist`. Wails lives in `desktop/` as a future window ([ADR 010](docs/architecture/010-wails-shell.md)); do not replace `/api` with bindings yet.
- **Loopback only.** Bind `127.0.0.1` starting at port **5000**; if that fails, try the next port, at most **20** attempts. Do not bind on all interfaces.
- **Generate is async.** Never block the HTTP handler on image drawing. Progress is a **process-wide singleton**; overlapping generates will clobber it — do not start a second generate without an explicit product decision.
- **`GET /api/system/status` consumes terminal states** (`done` / `error` → flush). Pollers must treat a following `empty` as “already observed.”
- **Quit is delayed 60s** and cancelled when other handlers call `StopQuit` (so navigating the SPA does not kill the process immediately). `-debug` never quits.
- **TTS TCP is best-effort Execute Lua on Global (`39999`, `messageID` 3, `guid` `-1`).** If TTS is not running, generate still writes `result/`. `/api/tts/data` is one-shot (clears the buffer). Official protocol: [External Editor API](https://api.tabletopsimulator.com/externaleditorapi/).
- **IDs:** string ids come from fsentry (name-derived). Card ids are integers. Path params named `{card}` are still parsed as that integer id.
- **Deck image = card back** for that deck’s sheets. Missing backside fails generate.
- **Keep swagger comments in `internal/api` in sync** when routes or bodies change (`make swagger`).
- **macOS data path is home**, not next to the binary.

---

## Agent working agreements

- Match existing Go style: interfaces in `contract.go`, unexported impl structs, errors from `internal/errors`.
- Keep this file updated when routes, disk layout, generate rules, or layers change.
- **Keep [README.md](README.md) short.** User-facing install, TTS copy path, and honest limitations stay there. Deep contracts stay here.
- **New core decisions get an ADR** in `docs/architecture/` (next number, update `INDEX.md`).
- **New user-visible behavior gets a use case** from `docs/use-cases/_TEMPLATE.md` and a row in `INDEX.md`.
- **Wiki pages** live in `docs/wiki/` as GitHub wiki markdown (`Home.md`, `_Sidebar.md`, slug links without `.md`). Copy into `DeckBuilder.wiki.git` to publish.
- Prefer the smallest change that preserves the GUI contract over a layer rewrite.
- Do not collapse `servers` / `services` / `db` “to simplify” without an ADR.
- Do not put business rules in `internal/api` (comments + registration only).
- Frontend is a submodule: `git clone --recursive`. Do not vendor a second copy of the GUI into `web/` by hand except via `make web-build`.
}
