# Wails shell

Wails v2 window with a copy of the Vue GUI (`gui/` → `desktop/frontend`). Catalog calls stay `fetch('/api/...')`. The process starts the same loopback HTTP server as `cmd/deck_builder` (no browser). In `wails dev`, Vite GET `/api` still targets `http://127.0.0.1:5000`.

`cmd/deck_builder` remains the production HTTP+SPA binary. This module is separate so root `go test ./...` does not compile CGO.

Wails v2.16+ CLI is required (v2.11 fails on Go 1.27 with `package "context" without types`).

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@v2.16.0
make wails-dev
```

Do not also run `cmd/deck_builder` on `:5000` at the same time.

See [ADR 010](../docs/architecture/010-wails-shell.md).
