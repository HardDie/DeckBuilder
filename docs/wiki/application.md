# `wire.go`

Developer reference. The map of the project is [Internals](Internals).

Composition root in package `main`. `wire` loads config, opens fsentry, creates the `games` folder, and constructs the services. `main.go` calls `wire`, then `servers.New` and the Wails bindings.

## `wire(version)`

Fatals if the store `Init` fails. Fatals if the `games` folder `Init` fails.

| Step | Meaning |
|---|---|
| `config.Get` | Paths and stamped version |
| `fsentry.New` + `Init` | Store on `cfg.Data` |
| `core.Init` | Creates `games/` |
| services | settings, game, collection, deck, card, TTS, generator, replace, search |

Order matters: `db.Init()` (root + lock) then `core.Init()` (`games/` folder).

`main.go` passes game, collection, deck, card, and TTS into `servers.New`. Bindings take the same services.
