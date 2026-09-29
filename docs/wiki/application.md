# `internal/application`

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
| `db` | `*fsentry.DB` on `cfg.Data`, then `db.Init()` | every `db/*.New` |
| `core` | Ensures `games/` exists | `core.Init()` at startup; tests also `Drop()` |
| `settings` … `card` | db aggregates | repositories / system service |
| `service*` / `server*` | Business + HTTP | `api.Register*` |

Order matters: `db.Init()` (root + lock) then `core.Init()` (`games/` folder).

## Methods

| Method | Role | Called from |
|---|---|---|
| `Get` | Wire everything | `cmd/deck_builder`, `desktop/` |
| `Listen` | Bind loopback (retry port), set TTS port | `Run`, `desktop/` |
| `Serve` | `http.Serve` the mux | `Run`, `desktop/` |
| `Handler` | The mux (for Wails `AssetServer.Handler`) | `desktop/` |
| `Config` / `GameService` / `CollectionService` / `DeckService` / `CardService` | Dependencies for Wails bindings | `desktop/` |
| `Run` | `Listen` + `Serve` | `cmd/deck_builder` |
| `ReplaceService` | Replace rules for the Wails binding | `desktop/` |
| `corsSetupHeaders` / `corsMiddleware` | `GET,OPTIONS` for image fetches | `routes.Use` |

No catalog logic lives here. Adding a new HTTP feature means a `New` + `Register*` pair in this file, not a new listen path.
