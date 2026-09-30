# 16. Auto-upload generated images to a web host

* **Status:** Proposed
* **Date:** 2026-09-30
* **Authors:** @oleg

---

## Context

1. Generate writes local `file:///` paths ([ADR 005](005-tts-generate-local-paths.md)).
2. A shared table needs http(s) `FaceURL` / `BackURL`.
3. Replace swaps those paths after the user hosts the files.
4. Auto-upload is future work.
5. This record starts the host list.
6. Append options here when implementation starts.

## Considered options

1. **Stay on local paths.**
   1. The user keeps `result/` or runs replace.
   2. This is the current product ([ADR 005](005-tts-generate-local-paths.md)).
2. **Steam Cloud from the existing Lua poke.**
   1. Execute Lua on `127.0.0.1:39999`.
   2. `messageID` 3, Global `guid` `-1`.
   3. The script would upload sheets and return cloud URLs.
3. **More hosts.**
   1. Not chosen yet.
   2. Append them when we implement upload.

## Decision

1. Do not build auto-upload yet.
2. Reject option 2.
3. Keep option 1 until a host is chosen.

## Steam Cloud

1. TTS hosts files on the user's Steam account.
2. Quota is 100 GB.
3. A hit looks like `http://cloud-3.steamusercontent.com/ugc/…/`.
4. Upload is the in-game UI.
   1. **Upload → Cloud Manager**, folder icon or **Upload All**.
   2. The **Cloud** choice when browsing a local file.
5. Lua has no upload call.
6. `WebRequest` is HTTP on the host only.
   1. It can GET, POST, PUT, or `custom`.
   2. It cannot open Steamworks.
   3. `put` UTF-8-encodes the body.
   4. That corrupts PNG and JPEG.
7. Docs:
   1. [Cloud Manager](https://kb.tabletopsimulator.com/custom-content/cloud-manager/).
   2. [Asset Importing](https://kb.tabletopsimulator.com/custom-content/asset-importing/).
   3. [WebRequest](https://api.tabletopsimulator.com/webrequest/manager/).

## Manual Cloud path

1. Spawn still uses local paths.
2. The user clicks **Upload All** in Cloud Manager.
3. TTS rewrites loaded local URLs to cloud URLs.
4. The user must save.
5. The Lua poke cannot click that control.

## Consequences

### Positive

1. The rejected host is written down.
2. Later hosts append under Considered options.

### Negative and risks

1. Shared tables still need replace or a manual host.
2. Upload All is a user click, not generate.

### Neutral

1. ADR 005 stays the generate rule until this record is Accepted.
