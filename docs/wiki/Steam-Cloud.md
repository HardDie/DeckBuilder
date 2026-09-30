# Upload sheets to Steam Cloud

Generate writes `file:///` paths into `FaceURL` and `BackURL`. Tabletop Simulator can upload those loaded files to your Steam Cloud and swap the live objects onto `http://cloud-3.steamusercontent.com/ugc/…/` URLs. DeckBuilder does not perform that upload. The button is inside TTS.

Why a Lua poke cannot do this: [ADR 016](https://github.com/HardDie/DeckBuilder/blob/master/docs/architecture/016-auto-upload-generated-images.md). Official UI: [Cloud Manager](https://kb.tabletopsimulator.com/custom-content/cloud-manager/).

## What you need

1. A rendered game. `result/` holds the PNG sheets and one JSON file.
2. Those PNGs still at the paths stored in the JSON. Do not move `result/` before this procedure.
3. TTS running, logged into Steam, on an empty table.

`result/` is `DeckBuilderData/result` next to the working directory, or `~/DeckBuilderData/result` on macOS.

Upload All sends every loaded file on the current table, local and cached. An empty table keeps other mods' images off your cloud.

## 1. Put the objects on the table

Either path is enough.

**Spawn.** Render while TTS is already open. DeckBuilder dials `127.0.0.1:39999` and the bag appears on the table. See [UC-10](https://github.com/HardDie/DeckBuilder/blob/master/docs/use-cases/generate/uc-10-tts-spawn.md).

**Saved Object.** Copy `result/<game>.json` into `Documents/My Games/Tabletop Simulator/Saves/Saved Objects`. In TTS, open Saved Objects and spawn that file.

The decks still store local paths, for example `file:///…/result/page.png`.

## 2. Save that table

**Games → Save.** Name it so the local paths stay obvious, for example `MyGame_local`.

Upload All reloads the game. The log line is `Reloading game to replace all files with Steam Cloud version`. Objects that exist only as an unsaved spawn are not part of that reload.

## 3. Upload all loaded files

1. Open **Modding → Cloud Manager**. Older builds use **Upload → Cloud Manager**.
2. Optional: type a folder name so this game's files stay together. The quota is 100 GB on your Steam account.
3. Click the up-arrow labeled **Upload all loaded files**.

TTS uploads each loaded file through Steam and reloads the save. On that reload it replaces each local address on the live objects with the new cloud URL. A custom deck's `FaceURL` and `BackURL` both change. Cards that share a sheet share the new face URL.

Cloud Manager does not edit `result/<game>.json`. That file still has `file:///` paths. Rendering again writes local paths again.

## 4. Save the cloud copy

**Games → Save** under a second name, for example `MyGame_cloud`.

Leave `MyGame_local` as it is. The cloud save is the one you load with friends or send to the Workshop. Without this second save, the file on disk still points at the local paths.

## Check

1. In Cloud Manager, click a sheet. TTS copies a URL like `http://cloud-3.steamusercontent.com/ugc/<id>/<hash>/`.
2. Load `MyGame_cloud` and inspect a deck. `FaceURL` and `BackURL` are that host.

The first time someone else loads that save, TTS downloads the sheets from Steam Cloud and caches them under `Documents/My Games/Tabletop Simulator/Mods`.

## Limits

1. Delete in Cloud Manager is permanent.
2. Cached assets already on the table are uploaded too.
3. This is separate from the replace binding. Replace rewrites a JSON file on disk from a mapping you supply. This procedure rewrites the TTS save only.
