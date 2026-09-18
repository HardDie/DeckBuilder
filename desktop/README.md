# Wails shell

Wails v2 window with a copy of the Vue GUI (`gui/` → `desktop/frontend`). Catalog calls stay `fetch('/api/...')`. Vite proxies `/api` to `http://127.0.0.1:5000` during `wails dev`.

Production HTTP + embed is still `cmd/deck_builder`. This module is separate so root `go test ./...` does not compile CGO.

Wails v2.16+ CLI is required (v2.11 fails on Go 1.27 with `package "context" without types`).

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@v2.16.0
# terminal 1
go run ./cmd/deck_builder -debug

# terminal 2
make wails-dev
```

See [ADR 010](../docs/architecture/010-wails-shell.md).
