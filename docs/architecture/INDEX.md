# Architecture decision records

ADRs live in this folder. Number them in order (`001`, `002`, …). Status is one of: Proposed, Accepted, Deprecated, Superseded.

These records describe **decisions already embodied in this repository**, so agents and humans share the same product vision. They are not a port plan.

| ID | Title | Status |
|---|---|---|
| [001](001-loopback-http-embedded-spa.md) | Loopback HTTP process; GUI is a separate SPA | Accepted |
| [002](002-four-level-catalog.md) | Game → collection → deck → card | Accepted |
| [003](003-fsentry-on-disk.md) | Catalog stored with fsentry under DeckBuilderData | Accepted |
| [004](004-layered-go-packages.md) | api → servers → services → repositories → db | Accepted |
| [005](005-tts-generate-local-paths.md) | Generate sprite sheets + TTS JSON with local image paths | Accepted |
| [006](006-docs-layout.md) | README, CURSOR.md, `docs/` (ADRs, use cases) | Accepted |
| [007](007-tts-official-api.md) | TTS integration follows api.tabletopsimulator.com | Accepted |
| [008](008-require-go-1.27.md) | Require Go 1.27.1 | Accepted |
| [009](009-catalog-timestamps.md) | Normalize catalog createdAt/updatedAt in memory | Accepted |
| [010](010-wails-shell.md) | Wails desktop shell beside loopback HTTP | Accepted |
| [011](011-wails-bindings.md) | Incremental Wails catalog bindings beside HTTP | Accepted |
| [012](012-card-image-urls.md) | Card faces load from an image URL | Accepted |
