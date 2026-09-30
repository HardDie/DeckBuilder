# DeckBuilder Development Guidelines

DeckBuilder is a **local authoring tool** for card games that export to [Tabletop Simulator](https://tabletopsimulator.com/) (TTS).

1. **UI**
   1. Wails window in `desktop/`.
   2. Vue sources live in `desktop/frontend`.
2. **HTTP** ([ADR 013](docs/architecture/013-http-images-and-tts.md))
   1. Loopback `127.0.0.1:5000`.
   2. Entity image bytes.
   3. `GET /api/tts/data` for TTS Lua.

**Not this product:**

1. A multi-user cloud service.
2. A TTS plugin inside the game.
3. More than a one-shot External Editor TCP poke.

Catalog verbs are Wails bindings ([ADR 011](docs/architecture/011-wails-bindings.md)).

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
- Settings (language `en`/`ru`, card scale, back shadow) persisted in the fsentry store.

**Known limitation (do not “fix” by inventing a host)**

- Generated JSON currently points at **local image paths**. Users must keep files in `result/` (or run **replace** after uploading images). Automatic upload to image hosting is **out of scope unless explicitly requested**.

**Out of scope unless explicitly requested**

- Binding on non-loopback interfaces or adding auth.
- Changing the domain tree (game/collection/deck/card).
- Replacing fsentry with SQL.
- Serving the Vue UI over HTTP. Images and TTS stay on loopback HTTP.

---

## Domain model

| Level | ID type | Role (example) |
|---|---|---|
| Game | string (folder id from name) | “Munchkin” |
| Collection | string | Base game or a DLC |
| Deck | string | “Monster”; **deck image is the card back** |
| Card | `int64` | Face image, description, variables (`HP: 2`), `count` |

List endpoints accept `sort` (`name`, `name_desc`, `created`, `created_desc`) and `search`. Create/update of named entities use **multipart/form-data** (`name`, `description`, optional `image` URL, optional `imageFile`). Card writes also send `variables` (JSON string) and `count` (desktop bindings map those fields).

JSON envelope: `{ "data": …, "meta": { "total", "cardsTotal?" }, "error": … }`.

---

## HTTP contract

Listen: **`127.0.0.1:5000`**, then `:5001` … if bind fails, up to **20** ports. Fail startup if all 20 fail. CORS is `*` with `GET,OPTIONS`; preflight `OPTIONS` returns 204.

| Method | Route | Role |
|---|---|---|
| GET | `/api/games/{game}/image` (and nested `…/image`) | Entity image bytes |
| GET | `/api/tts/data` | One-shot last generated JSON for TTS Lua |

Do not invent REST aliases. If the GUI needs a field, add it on the existing DTO with `json` tags and keep `error` / `data` wrapping.

**Generate** (`generator.Game`) returns immediately; work runs in a goroutine. The GUI polls system `Status`. Status values: `empty`, `in_progress`, `done`, `error`.

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
| GUI | Vue in `desktop/frontend`; Wails window ([ADR 013](docs/architecture/013-http-images-and-tts.md)) |
| HTTP | Go 1.27.1 (`go.mod`), `gorilla/mux`, listen `127.0.0.1:5000` (try next port up to 20 times if busy) |
| Persistence | `github.com/HardDie/fsentry` |
| Images | `disintegration/imaging`, `internal/page_drawer`, `internal/images` |
| Tests | `go test ./... -race`; fuzz targets under services |

---

## Build & Run

```bash
git clone https://github.com/HardDie/DeckBuilder
make build       # Wails app for this machine (desktop/build/bin)
make wails-dev   # Wails window
make test
make linter-run  # after make linter-install
```

`cmd/deck_builder` listens for images and TTS. It does not open a browser. The UI is `make wails-dev`.

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
| HTTP routes | `internal/api` |
| Process wiring | `internal/application/application.go` |
| Paths, sheet limits | `internal/config` |
| Domain structs | `internal/entities` |
| JSON DTOs | `internal/dto` |
| TTS Saved Object shapes | `internal/tts_entity` |
| HTTP handlers | `internal/servers` |
| Business logic | `internal/services` |
| Image + fs helpers | `internal/repositories`, `internal/db`, `internal/fs` |
| Sprite composition | `internal/page_drawer` |
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
├── cmd/deck_builder/       # image + TTS HTTP server
├── desktop/                # Wails v2 window (own module; starts the same HTTP server)
├── internal/
│   ├── application/        # mux, DI, ListenAndServe
│   ├── api/                # route registration
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
│   ├── network/            # JSON envelope
│   ├── errors/
│   └── …
├── pkg/version/            # exact git tag or 12-character commit hash
├── web/                    # dist stub, not served
└── tools/                  # join, copy_cards_variables
```

**Layer rules**

- **`api` registers routes.** Handlers live in `servers`.
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

- **The UI is the Wails window.** Game (including export and import), collection, deck, card, system, search, generator, and replace use bindings under `desktop/bindings/` ([ADR 011](docs/architecture/011-wails-bindings.md), [ADR 013](docs/architecture/013-http-images-and-tts.md)). HTTP serves images and TTS.
- **Loopback only.** Bind `127.0.0.1` starting at port **5000**; if that fails, try the next port, at most **20** attempts. Do not bind on all interfaces.
- **Generate is async.** Never block the binding on image drawing. Progress is a **process-wide singleton**; overlapping generates will clobber it — do not start a second generate without an explicit product decision.
- **System `Status` consumes terminal states** (`done` / `error` → flush). Pollers must treat a following `empty` as “already observed.”
- **TTS TCP is best-effort Execute Lua on Global (`39999`, `messageID` 3, `guid` `-1`).** If TTS is not running, generate still writes `result/`. `/api/tts/data` is one-shot (clears the buffer). Official protocol: [External Editor API](https://api.tabletopsimulator.com/externaleditorapi/).
- **IDs:** string ids come from fsentry (name-derived). Card ids are integers. Path params named `{card}` are still parsed as that integer id.
- **Deck image = card back** for that deck’s sheets. Missing backside fails generate.
- **macOS data path is home**, not next to the binary.
- **Version** is `pkg/version`.
  - Stamp `-X github.com/HardDie/DeckBuilder/pkg/version.Build`.
  - Value is an exact git tag, or the 12-character commit hash.
  - A blank stamp is `dev`.

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
- The Vue UI that ships is `desktop/frontend`. Do not embed it under `web/`.
}
