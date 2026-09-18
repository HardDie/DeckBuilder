# Cmd and tools

## `cmd/deck_builder`

`main`: parse `-debug`, stamp `Version` / git commits, `application.Get`, open browser unless debug, `app.Run()`. Swagger `meta` comments live on this package.

## `web/`

`//go:embed` of GUI `dist` + `swagger.json`. Served by `api/static.go`. Rebuild GUI with `make web-build`.

## `tools/copy_cards_variables`

Offline helper that copies `variables` between two `cards/.info.json` files. Uses `fsentry.QuotedString`. Not on the HTTP path.

## `tools/join`

Small CLI around HTTP requests (fill/join helpers). Not part of the server binary.
