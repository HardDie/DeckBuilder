<p align="center">
  <img alt="logo" src="deployment/512.png" height="150" />
  <h3 align="center">DeckBuilder</h3>
  <p align="center">Create your own deck of cards for the Tabletop Simulator</p>
</p>

---

The frontend for this application: https://github.com/lmm1ng/DeckBuilderGUI

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
git clone https://github.com/HardDie/DeckBuilder --recursive
```

Check that all necessary packages are installed
```
./deployment/check_binary.sh
```

Build web
```
make web-build
```

Build binary
```
make build
```

The resulting files can be found in the deployment/out folder
