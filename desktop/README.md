# Wails shell (scaffold)

Wails v2 project for a future desktop window. Production today is still `cmd/deck_builder` (loopback HTTP + embedded `web/dist`).

This directory is a **separate Go module** so `go test ./...` in the repo root does not compile CGO/webview.

```bash
cd desktop
wails dev    # placeholder vanilla UI
wails build
```

Do not move catalog/TTS onto Wails bindings in this step. See [ADR 010](../docs/architecture/010-wails-shell.md).
