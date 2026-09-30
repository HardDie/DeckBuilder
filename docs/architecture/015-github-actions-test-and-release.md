# 15. GitHub Actions: test on push, binaries on tag

* **Status:** Accepted
* **Date:** 2026-09-30
* **Authors:** @oleg

---

## Context

1. Catalog tests should run on every push.
2. Players need downloadable binaries.
3. Wails uses CGO and a WebView.
4. `main.go` must stay out of CI tests.
5. Linux releases need a user-level install, not only a raw binary.

## Considered options

1. **No CI** — local `go test` only.
2. **One Linux cross-compile** — misses WebView and the other OS binaries.
3. **Repo workflows**
   1. `test.yml` on every push and pull request.
   2. `release.yml` on `v*` tags.

## Decision

Use option 3.

1. **Tests**
   1. `//go:build !nomain` is on `main.go`.
   2. `make test` and CI use `-race -count=1 -tags=nomain`.
   3. Packages are `.`, `./bindings/...`, `./internal/...`, and `./pkg/...`.
   4. Then `go test -tags=integration -count=1 ./internal/...`.
   5. Go is 1.27.1.
   6. `make ci` is that sequence.
2. **Releases**
   1. Tag `v1.2.3`.
   2. Checkout uses `fetch-depth: 0` so tags exist.
   3. Wails CLI matches `go.mod`.
   4. `go mod download` retries with `GODEBUG=http2client=0`.
   5. Build is `wails build -skipbindings -platform …`.
   6. Stamp `-X github.com/HardDie/DeckBuilder/pkg/version.Build=…`.
   7. `frontend/wailsjs` stays committed, so generate it locally after a binding change.
   8. Linux passes `-tags webkit2_41` (`libwebkit2gtk-4.1-dev`).
   9. Runners: `ubuntu-latest`, `ubuntu-24.04-arm`, `windows-latest`, `macos-latest` (`darwin/universal`).
   10. Linux installs GTK, WebKit 4.1, `libx11-dev`, and runs under `xvfb-run`.
   11. Frontend install is yarn 1 (`frontend/yarn.lock`).
   12. `scripts/sync-product-version.sh` writes `info.productVersion` (tag without `v`).
   13. That fills macOS bundle versions and Windows `build/windows/info.json`.
   14. The publish job uploads archives and `SHA256SUMS.txt`.
   15. Names look like `DeckBuilder-v1.2.3-linux-amd64.tar.gz`.
3. **Linux archive**
   1. Also contains `install.sh`, `DeckBuilder.desktop`, and `DeckBuilder.png`.
   2. `install.sh` copies them into the user XDG dirs.
   3. No root.
4. A tag without `wails.json` fails the release job.
5. Code signing is out of scope.

## Consequences

### Positive

1. Tests run on every push.
2. Four OS/arch archives from one tag.
3. `make ci` matches the test workflow.

### Negative and risks

1. macOS and Windows builds use GitHub-hosted minutes.
   1. Public repo: standard runners are free.
   2. Private repo: those minutes are billed.
2. If `ubuntu-24.04-arm` disappears, switch arm64 to a cross-compile.
3. `go test ./...` still compiles Wails CGO.

### Neutral

1. Workflows: `.github/workflows/test.yml`, `.github/workflows/release.yml`.
2. Linux files: `build/linux/install.sh`, `build/linux/DeckBuilder.desktop`.
