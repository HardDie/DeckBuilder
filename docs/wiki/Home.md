# DeckBuilder wiki

Developer notes for the Go process. The UI is the Wails window. Product intro and TTS Saved Objects path: [README](https://github.com/HardDie/DeckBuilder#readme). HTTP contract and agent rules: [CURSOR.md](https://github.com/HardDie/DeckBuilder/blob/master/CURSOR.md).

| Page | What it covers |
|---|---|
| [Overview](Overview) | Request flow, two fsentry handles, on-disk layout, IDs |
| [Application](Application) | Composition root in `internal/application` |
| [Config](Config) | Paths, generate limits, `Games()` vs `gamesPath` |
| [API](API) | Route registration |
| [Servers](Servers) | HTTP adapters |
| [Services](Services) | Rules: catalog, generate, search, replace, system, TTS |
| [Repositories](Repositories) | Images, zip, config paths |
| [DB](DB) | fsentry mapping, fields, methods, callers |
| [Entities and DTOs](Entities-and-DTOs) | Domain vs GUI JSON vs TTS JSON |
| [Helpers](Helpers) | errors, network, fs, images, progress, page_drawer |
| [Cmd and tools](Cmd-and-tools) | `cmd/deck_builder`, root Wails window, `tools/` |

## Layer rule

```text
main / cmd → application → api + servers → services → repositories → db
                                      ↘ generator uses page_drawer, tts_entity, progress, fs
```

Servers must not import fsentry. `internal/api` must not contain business rules.

Decisions: [docs/architecture](https://github.com/HardDie/DeckBuilder/tree/master/docs/architecture). Scenarios: [docs/use-cases](https://github.com/HardDie/DeckBuilder/tree/master/docs/use-cases).

## Publishing these pages

The markdown in `docs/wiki/` in the source repository is the source for the GitHub wiki (`https://github.com/HardDie/DeckBuilder/wiki`). Copy `*.md` (including `_Sidebar.md` and `_Footer.md`) into the wiki git remote:

```bash
git clone https://github.com/HardDie/DeckBuilder.wiki.git
cp /path/to/DeckBuilder/docs/wiki/*.md DeckBuilder.wiki/
cd DeckBuilder.wiki
git add .
git commit -m "Sync wiki from repo docs/wiki/"
git push
```

GitHub wiki links in these files are **page slugs without `.md`**, for example `[DB](DB)` → `https://github.com/HardDie/DeckBuilder/wiki/DB`.
