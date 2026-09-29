# Wails shell

Wails v2 window. Vue UI is `desktop/frontend`. Catalog verbs, export, import, generate, and replace are Wails bindings. The process also starts the loopback HTTP server for images and TTS. In `wails dev`, Vite GET `/api` still targets `http://127.0.0.1:5000`.

`cmd/deck_builder` is that HTTP server without a window. This module is separate so root `go test ./...` does not compile CGO.

Wails v2.16+ CLI is required (v2.11 fails on Go 1.27 with `package "context" without types`).

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@v2.16.0
make wails-dev
```

Do not also run `cmd/deck_builder` on `:5000` at the same time.

See [ADR 010](../docs/architecture/010-wails-shell.md), [ADR 011](../docs/architecture/011-wails-bindings.md), and [ADR 013](../docs/architecture/013-http-images-and-tts.md).
