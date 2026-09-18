# `internal/application`

Composition root. Only this package constructs the graph and starts HTTP.

## Types

| Name | Meaning |
|---|---|
| `Application` | Holds `*mux.Router`. `Run` binds `127.0.0.1:5000`. |
| `Get(debugFlag, version)` | Builds config, stores, layers, middleware. Fatals on store `Init` errors. |

## Variables in `Get`

| Variable | Meaning | Who consumes it |
|---|---|---|
| `cfg` | Paths, debug, stamped version | Almost every New below |
| `routes` | Root mux | `api.Register*` |
| `db` | `*fsentry.DB` on `cfg.Data`, then `db.Init()` | every `db/*.New` |
| `core` | Ensures `games/` exists | `core.Init()` at startup; tests also `Drop()` |
| `settings` … `card` | db aggregates | repositories / system service |
| `service*` / `server*` | Business + HTTP | `api.Register*` |

Order matters: `db.Init()` (root + lock) then `core.Init()` (`games/` folder).

## Methods

| Method | Role | Called from |
|---|---|---|
| `Get` | Wire everything | `cmd/deck_builder` |
| `Run` | `http.Handle("/", router)` + listen | `cmd/deck_builder` |
| `corsSetupHeaders` / `corsMiddleware` | Allow GUI-dev on another origin | `routes.Use` |

No catalog logic lives here. Adding a new HTTP feature means a `New` + `Register*` pair in this file, not a new listen path.
