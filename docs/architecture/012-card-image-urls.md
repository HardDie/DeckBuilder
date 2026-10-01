# 12. Card faces load from an image URL

* **Status:** Accepted
* **Date:** 2026-09-27
* **Authors:** @oleg

---

## Context

1. A card row needs a face plus name, description, variables, and count.
2. Catalog verbs are Wails bindings ([ADR 011](011-wails-bindings.md)).
3. Those bindings already return `cachedImage`.
4. The value is `/api/.../image?<hash of updatedAt>`.
5. `internal/servers` writes the file bytes and `Content-Type`.
6. `main.go` mounts that mux as `AssetServer.Handler`.
7. The Vue grid binds `:img="item.cachedImage"`.
8. TTS still needs an HTTP image URL ([ADR 005](005-tts-generate-local-paths.md)).

## What Wails shows

1. [Dogs API tutorial](https://wails.io/docs/tutorials/dogsapi/)
   1. Bindings return URL strings.
   2. `<img src>` loads the pixels.
   3. The binding does not carry the file.
2. [Dynamic assets](https://wails.io/docs/guides/dynamic-assets/)
   1. Files outside the embed use `AssetServer.Handler`.
   2. That handler is `http.Handler`.
   3. The page still uses a normal `src`.
   4. Keep the handler on a known path.
3. The webview blocks `file://` ([discussion 1538](https://github.com/wailsapp/wails/discussions/1538)).
4. Showcases
   1. GameStacker and YTD pages do not document image transport.
   2. Clustta stores preview blobs in sqlite.
   3. Clustta syncs those blobs with HTTP `/previews`.
5. A bound `[]byte` return
   1. The dispatcher `json.Marshal`s the result (`calls.go`).
   2. Go encodes `[]byte` as a base64 string.
   3. The generated JS type for that return is `string`.

## Considered options

1. **A read binding returns `[]byte`.**
   1. JS decodes base64.
   2. JS builds a blob URL.
   3. One bridge call per card.
   4. The webview image cache does not apply.
   5. `loading="lazy"` starts after that call.
   6. A list response must not embed the bytes.
2. **A binding returns a disk path.**
   1. The webview rejects `file://`.
   2. Display still needs an HTTP handler.
3. **A binding returns `cachedImage`. HTTP serves the bytes.**
   1. This matches the current code.
   2. The GUI and TTS share one URL.
   3. The handler can stay on the asset server.
4. **A remote URL only, as in the dogs tutorial.**
   1. Faces are files under the catalog.
   2. The `image` field is the source link.
   3. `cachedImage` is the local file.

## Decision

Use option 3.

1. Bindings return card fields and `cachedImage`.
2. Image reads stay `GET /api/.../image` on the loopback mux.
3. `AssetServer.Handler` already forwards `/api`.
4. A write may still send `imageFile []byte`.
   1. That is one file.
   2. That is one user action.
5. Cache bust stays the `updatedAt` hash on the query string.

## Consequences

### Positive

1. The grid uses the webview image loader.
2. Many faces load in parallel.
3. The HTTP cache keys off `cachedImage`.
4. Later REST removals can leave this handler in place.

### Negative and risks

1. Image GET stays HTTP after catalog CRUD is bindings-only.
2. Wails notes that the dynamic-asset example breaks on Vite 5.
   1. Our image URLs use the `/api` prefix.
   2. That prefix already hits `AssetServer.Handler`.

### Neutral

1. ADR 011 still stands for catalog verbs.
2. This record covers the image read path only.
