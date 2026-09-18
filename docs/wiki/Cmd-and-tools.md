# Cmd and tools

## `cmd/deck_builder`

`main`: parse `-debug`, stamp `Version` / git commits, `application.Get`, `app.Run()` (listen + optional browser). Swagger `meta` comments live on this package.

## `desktop/`

Wails v2 scaffold (separate Go module). Vue UI copied from `gui/` into `desktop/frontend`; still talks HTTP (`/api` → `127.0.0.1:5000`: Vite for GET in `wails dev`, `AssetServer.Handler` for POST/PATCH/DELETE). Not the production binary ([ADR 010](../architecture/010-wails-shell.md)). `cd desktop && wails dev`.

## `web/`

`//go:embed` of GUI `dist` + `swagger.json`. Served by `api/static.go`. Rebuild GUI with `make web-build`.

## `tools/copy_cards_variables`

Offline helper that copies `variables` between two `cards/.info.json` files. Uses `fsentry.QuotedString`. Not on the HTTP path.

## `tools/join`

Small CLI around HTTP requests (fill/join helpers). Not part of the server binary.
