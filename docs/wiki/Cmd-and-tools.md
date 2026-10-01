# Cmd and tools

Developer reference. The map of the project is [Internals](Internals).

## Wails window

Root module (`main.go`, `wails.json`). The only entrypoint. Starts the image and TTS HTTP server. Vue UI is `frontend/`. Catalog, export, import, generate, and replace are Wails bindings ([ADR 011](../architecture/011-wails-bindings.md), [ADR 013](../architecture/013-http-images-and-tts.md), [ADR 014](../architecture/014-wails-at-module-root.md)). `AssetServer.Handler` serves `/api` images (Vite GET in `wails dev` still targets `:5000`). `make dev`. `make build` writes the binary to `build/bin`.

## `scripts/screenshots`

`make screenshots` installs Playwright under `scripts/screenshots/` and writes `docs/wiki/images/screen-*.png`. It starts Vite with `?screenshot=1`. That flag stubs the Wails bindings with the Four Souls fixture in `frontend/src/api/screenshot.js`. No Go process is required.

## `tools/copy_cards_variables`

Offline helper that copies `variables` between two `cards/.info.json` files. Uses `fsentry.QuotedString`. Not on the HTTP path.

## `tools/join`

Small CLI around HTTP requests (fill/join helpers). Not part of the server binary.
