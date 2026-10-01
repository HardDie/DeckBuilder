# DeckBuilder wiki

Two guides. Pick the one that matches why you opened this.

| | Start here | What it is |
|---|---|---|
| Using the app | [Guide](Guide) | What the app does, how the four levels work, and why it is a better place to build a card game than the editor inside Tabletop Simulator |
| | [Steam Cloud](Steam-Cloud) | Upload a rendered table's images and swap local paths for Steam Cloud URLs |
| Developing | [Internals](Internals) | Each layer, from the Wails window down to the files on disk |
| | [Generation](Generation) | How Render draws sheets and writes the TTS JSON |

The pictures in the guide are an example catalog from [The Binding of Isaac: Four Souls](https://foursouls.com/).

Product intro and the Saved Objects copy path also live in the [README](https://github.com/HardDie/DeckBuilder#readme). HTTP contracts and agent rules: [CURSOR.md](https://github.com/HardDie/DeckBuilder/blob/master/CURSOR.md). Decisions: [docs/architecture](https://github.com/HardDie/DeckBuilder/tree/master/docs/architecture).

## Developer reference

These pages are the field notes. [Internals](Internals) is the map.

| Page | What it covers |
|---|---|
| [Overview](Overview) | Request flow, fsentry handle, on-disk layout, IDs |
| [Application](Application) | Composition root in `wire.go` |
| [Config](Config) | Paths, generate limits, `Games()` vs `gamesPath` |
| [API](API) | Image and TTS routes |
| [Servers](Servers) | HTTP adapters |
| [Services](Services) | Rules: catalog, generate, search, replace, system, TTS |
| [Repositories](Repositories) | Images, zip, config paths |
| [DB](DB) | fsentry mapping, fields, methods, callers |
| [Entities and DTOs](Entities-and-DTOs) | Domain vs GUI JSON vs TTS JSON |
| [Helpers](Helpers) | errors, network, fs, images, render/deprecated/progress, render/deprecated/page_drawer |
| [Page drawer](Page-drawer) | `PageDrawer` methods, cell size, sheet and back files |
| [Cmd and tools](Cmd-and-tools) | root Wails window, `tools/` |

## Publishing these pages

The markdown in `docs/wiki/` is the source for the GitHub wiki (`https://github.com/HardDie/DeckBuilder/wiki`). Copy the pages and the `images/` folder into the wiki git remote:

```bash
git clone https://github.com/HardDie/DeckBuilder.wiki.git
cp /path/to/DeckBuilder/docs/wiki/*.md DeckBuilder.wiki/
mkdir -p DeckBuilder.wiki/images
cp -R /path/to/DeckBuilder/docs/wiki/images/. DeckBuilder.wiki/images/
cd DeckBuilder.wiki
git add .
git commit -m "Sync wiki from repo docs/wiki/"
git push
```

GitHub wiki links in these files are **page slugs without `.md`**, for example `[Guide](Guide)`. Images use a relative path such as `images/screen-games.png`.
