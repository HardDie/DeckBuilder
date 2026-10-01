# `internal/application`

Developer reference. The map of the project is [Internals](Internals).

Composition root. Only this package constructs the graph and starts HTTP.

## Types

| Name | Meaning |
|---|---|
| `Application` | Holds config, router, TTS, and catalog services. `Listen` binds `127.0.0.1:5000`, then +1 on failure (max 20). |
| `Get(version)` | Builds config, stores, layers, middleware. Fatals on store `Init` errors. |

## Variables in `Get`

| Variable | Meaning | Who consumes it |
|---|---|---|
| `cfg` | Paths, stamped version | Almost every New below |
| `routes` | Root mux | `api.Register*` |
| `db` | `*fsentry.DB` on `cfg.Data`, then `db.Init()` | `internal/repositories` |
| `core` | `internal/repositories/core` ensures `games/` exists | `core.Init()` at startup; tests also `Drop()` |
| `settings` … `card` | `internal/repositories/{settings,game,collection,deck,card}` | services |
| `service*` / `server*` | Business + HTTP | `api.Register*` |

Order matters: `db.Init()` (root + lock) then `core.Init()` (`games/` folder).

## Methods

| Method | Role | Called from |
|---|---|---|
| `Get` | Wire everything | `cmd/deck_builder`, `main.go` |
| `Listen` | Bind loopback (retry port), set TTS port | `Run`, `main.go` |
| `Serve` | `http.Serve` the mux | `Run`, `main.go` |
| `Handler` | The mux (for Wails `AssetServer.Handler`) | `main.go` |
| `Config` / `GameService` / `CollectionService` / `DeckService` / `CardService` | Dependencies for Wails bindings | `main.go` |
| `Run` | `Listen` + `Serve` | `cmd/deck_builder` |
| `ReplaceService` | Replace rules for the Wails binding | `main.go` |
| `corsSetupHeaders` / `corsMiddleware` | `GET,OPTIONS` for image fetches | `routes.Use` |

No catalog logic lives here. Adding a new HTTP feature means a `New` + `Register*` pair in this file, not a new listen path.
