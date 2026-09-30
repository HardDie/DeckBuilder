# Cmd and tools

## `cmd/deck_builder`

`main`: `version.String()` into `application.Get`, then `app.Run()` (listen for images and TTS). Stamp `-X github.com/HardDie/DeckBuilder/pkg/version.Build`. No browser.

## Wails window

Root module (`main.go`, `wails.json`). Starts the same image and TTS HTTP server. Vue UI is `frontend/`. Catalog, export, import, generate, and replace are Wails bindings ([ADR 011](../architecture/011-wails-bindings.md), [ADR 013](../architecture/013-http-images-and-tts.md), [ADR 014](../architecture/014-wails-at-module-root.md)). `AssetServer.Handler` serves `/api` images (Vite GET in `wails dev` still targets `:5000`). `make wails-dev`. `make build` writes the binary to `build/bin`.

## `tools/copy_cards_variables`

Offline helper that copies `variables` between two `cards/.info.json` files. Uses `fsentry.QuotedString`. Not on the HTTP path.

## `tools/join`

Small CLI around HTTP requests (fill/join helpers). Not part of the server binary.
