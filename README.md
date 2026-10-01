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
Right-click on the game and select the menu item "Render". As a result you will have a folder DeckBuilderData/result, in which you will find png images and one json file.

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
make wails-dev
```

Build the app
```
make build
```

Sheet JPEGs use libjpeg-turbo through `github.com/pixiv/go-libjpeg` (cgo).
On a Mac, install the library and point cgo at it before `make build` or `make wails-dev`:

```
brew install jpeg-turbo
export CGO_CFLAGS="-I$(brew --prefix jpeg-turbo)/include"
export CGO_LDFLAGS="-L$(brew --prefix jpeg-turbo)/lib"
```

The binary is in `build/bin`

## Documentation

- Users: this README (product, TTS export, build). Longer tour: [docs/wiki/Guide.md](docs/wiki/Guide.md).
- Agents and contributors: [CURSOR.md](CURSOR.md) (HTTP contract, packages, generate rules).
- Decisions: [docs/architecture](docs/architecture/INDEX.md). Scenarios: [docs/use-cases](docs/use-cases/INDEX.md).
- Go packages (fields, methods, callers): [docs/wiki](docs/wiki/Home.md). Sync to the GitHub wiki with `git clone https://github.com/HardDie/DeckBuilder.wiki.git` and copy `docs/wiki/*.md`.

