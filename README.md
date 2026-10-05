<p align="center">
  <img alt="logo" src="build/appicon.png" height="150" />
  <h3 align="center">DeckBuilder</h3>
  <p align="center">Create your own deck of cards for the Tabletop Simulator</p>
</p>

---

The app window is the Wails project at the repo root.

[Video guide [ENG]](https://www.youtube.com/watch?v=jty_nEsGGJg)

[Video guide [RUS]](https://www.youtube.com/watch?v=r0-4mW8gX1w)

## Description
This utility helps you easily create, modify, and export card games for Tabletop Simulator.
The utility has four logical objects.
- on the main screen, you can create a game, for example, Munchkin
- on the next level you can create a collection, for example, the base game or the first DLC
- in the third level, you can create a deck, for example Monster
- on the last level you can add cards, and you can set variables for the internal lua, such as HP: 2 or AT: 1

## How to create files for TTS
Right-click on the game and select the menu item "Render". As a result you will have a folder DeckBuilderData/result/<game>, in which you will find the sheet images and one json file. Rendering again only redraws the pages that changed.

At the moment the json path to the pictures is set as they are located on the HDD, and you have to keep them in the result folder. This means that for now you cannot save a deck of cards to the table.
But in the future there will be support for automatically uploading images to some image storage.

You must copy the json file to Saved Object, which is in the Tabletop Simulator save files. The path should look like this: "Tabletop Simulator/Saves/Saved Objects".
Then you can start the game, open Saved Objects and find this object.

## How to build
Clone repository:
```
git clone https://github.com/HardDie/DeckBuilder
```

Run the window
```
make dev
```

Build the app
```
make build
```

The app links two C libraries statically: libjpeg-turbo and libwebp.
The first `make` builds them from source into `build/third_party` (see [third_party](third_party/README.md)).
That needs cmake and nasm:

```
brew install cmake nasm            # macOS
sudo apt install cmake nasm        # Ubuntu
```

On Windows (MSYS2 MINGW64): `pacman -S make mingw-w64-x86_64-gcc mingw-w64-x86_64-cmake mingw-w64-x86_64-nasm`.

The binary is in `build/bin`

## Documentation

- Users: this README (product, TTS export, build). Longer tour: [docs/wiki/Guide.md](docs/wiki/Guide.md).
- Agents and contributors: [CLAUDE.md](CLAUDE.md) (HTTP contract, packages, generate rules).
- Decisions: [docs/architecture](docs/architecture/INDEX.md). Scenarios: [docs/use-cases](docs/use-cases/INDEX.md).
- Go packages (fields, methods, callers): [docs/wiki](docs/wiki/Home.md). Sync to the GitHub wiki with `git clone https://github.com/HardDie/DeckBuilder.wiki.git` and copy `docs/wiki/*.md`.

