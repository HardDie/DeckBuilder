# Cmd and tools

## `cmd/deck_builder`

`main`: parse `-debug`, stamp `Version` / git commits, `application.Get`, open browser unless debug, `app.Run()`. Swagger `meta` comments live on this package.

## `web/`

`//go:embed` of GUI `dist` + `swagger.json`. Served by `api/static.go`. Rebuild GUI with `make web-build`.

## `tools/copy_cards_variables`

Offline helper that reads card JSON using **legacy** `QuotedString` types. Not on the HTTP path. Update it when card payload types move to v0.1 `fsentry.QuotedString`.

## `tools/join`

Small CLI around HTTP requests (fill/join helpers). Not part of the server binary.
