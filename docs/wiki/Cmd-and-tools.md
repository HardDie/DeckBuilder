# Cmd and tools

## `cmd/deck_builder`

`main`: stamp `Version` / git commits, `application.Get`, `app.Run()` (listen for images and TTS). Swagger `meta` comments live on this package. No browser.

## `desktop/`

Wails v2 window (separate Go module). Starts the same image and TTS HTTP server. Vue UI is `desktop/frontend`. Catalog, export, import, generate, and replace are Wails bindings ([ADR 011](../architecture/011-wails-bindings.md), [ADR 013](../architecture/013-http-images-and-tts.md)). `AssetServer.Handler` serves `/api` images (Vite GET in `wails dev` still targets `:5000`). `make wails-dev`.

## `web/`

`swagger.json` from `make swagger`. Not served over HTTP.

## `tools/copy_cards_variables`

Offline helper that copies `variables` between two `cards/.info.json` files. Uses `fsentry.QuotedString`. Not on the HTTP path.

## `tools/join`

Small CLI around HTTP requests (fill/join helpers). Not part of the server binary.
