# 14. Wails project lives in the module root

* **Status:** Accepted
* **Date:** 2026-09-30
* **Authors:** @oleg

---

## Context

1. The window is the only UI ([ADR 013](013-http-images-and-tts.md)).
2. Wails lived in `desktop/` as its own module ([ADR 010](010-wails-shell.md)).
3. There is no separate browser GUI.

## Considered options

1. **Keep the nested `desktop/` module.**
2. **Move Wails to the module root.**
   1. One `go.mod`.
   2. Vue stays in `frontend/`.
   3. Bindings stay in `bindings/`.

## Decision

Use option 2.

1. `main.go`, `app.go`, and `wails.json` sit at the repo root.
2. Vue sources are `frontend/`.
3. Catalog bindings are `bindings/`.
4. `make build` writes `build/bin`.
5. `main.go` is the only entrypoint.
6. `frontend/dist/.gitkeep` lets `//go:embed` compile before the first UI build.

This supersedes the nested `desktop/` module in [ADR 010](010-wails-shell.md).

## Consequences

### Positive

1. One Go module.
2. `wails dev` and `wails build` run from the repo root.

### Negative and risks

1. `go test ./...` compiles the Wails packages.
2. That compile needs CGO and a WebView.
3. CI uses `-tags=nomain` ([ADR 015](015-github-actions-test-and-release.md)).

### Neutral

1. HTTP still serves images and TTS ([ADR 013](013-http-images-and-tts.md)).
