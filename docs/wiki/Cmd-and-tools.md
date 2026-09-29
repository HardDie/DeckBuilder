# Cmd and tools

## `cmd/deck_builder`

`main`: parse `-debug`, stamp `Version` / git commits, `application.Get`, `app.Run()` (listen + optional browser). Swagger `meta` comments live on this package.

## `desktop/`

Wails v2 window (separate Go module). Starts the same `internal/application` HTTP server (no browser). Vue UI copied from `gui/` into `desktop/frontend`. Game, collection, deck, and card catalog verbs use Wails bindings. Game export is a Wails binding (native save dialog). Game import is a Wails binding. Generate still `fetch('/api/...')` ([ADR 011](../architecture/011-wails-bindings.md)). `AssetServer.Handler` serves remaining `/api` (Vite GET in `wails dev` still targets `:5000`). Not the production `cmd/deck_builder` binary ([ADR 010](../architecture/010-wails-shell.md)). `make wails-dev`.

## `web/`

`//go:embed` of GUI `dist` + `swagger.json`. Served by `api/static.go`. Rebuild GUI with `make web-build`.

## `tools/copy_cards_variables`

Offline helper that copies `variables` between two `cards/.info.json` files. Uses `fsentry.QuotedString`. Not on the HTTP path.

## `tools/join`

Small CLI around HTTP requests (fill/join helpers). Not part of the server binary.
