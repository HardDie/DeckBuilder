# Architecture decision records

ADRs live in this folder. Number them in order (`001`, `002`, …). Status is one of: Proposed, Accepted, Deprecated, Superseded.

These records describe **decisions already embodied in this repository**, so agents and humans share the same product vision. They are not a port plan.

| ID | Title | Status |
|---|---|---|
| [001](001-loopback-http-embedded-spa.md) | Loopback HTTP process; GUI is a separate SPA | Superseded |
| [002](002-four-level-catalog.md) | Game → collection → deck → card | Accepted |
| [003](003-fsentry-on-disk.md) | Catalog stored with fsentry under DeckBuilderData | Accepted |
| [004](004-layered-go-packages.md) | api → servers → services → repositories | Accepted |
| [005](005-tts-generate-local-paths.md) | Generate sprite sheets + TTS JSON with local image paths | Accepted |
| [006](006-docs-layout.md) | README, CURSOR.md, `docs/` (ADRs, use cases) | Superseded |
| [007](007-tts-official-api.md) | TTS integration follows api.tabletopsimulator.com | Accepted |
| [008](008-require-go-1.27.md) | Require Go 1.27.1 | Accepted |
| [009](009-catalog-timestamps.md) | Normalize catalog createdAt/updatedAt in memory | Accepted |
| [010](010-wails-shell.md) | Wails desktop shell beside loopback HTTP | Accepted |
| [011](011-wails-bindings.md) | Incremental Wails catalog bindings beside HTTP | Accepted |
| [012](012-card-image-urls.md) | Card faces load from an image URL | Accepted |
| [013](013-http-images-and-tts.md) | HTTP serves images and TTS; Wails hosts the UI | Accepted |
| [014](014-wails-at-module-root.md) | Wails project lives in the module root | Accepted |
| [015](015-github-actions-test-and-release.md) | GitHub Actions tests on push; binaries on tag | Accepted |
| [016](016-auto-upload-generated-images.md) | Auto-upload generated images to a web host | Proposed |
| [017](017-sheet-draw-speed.md) | Sheet draw speed trials | Accepted |
| [018](018-generation-progress.md) | Generation progress in `internal/render/progress` | Accepted |
| [019](019-claude-md.md) | CLAUDE.md replaces CURSOR.md | Accepted |
| [020](020-image-input-limits.md) | Limits on incoming images | Accepted |
| [021](021-image-download-progress.md) | Image download progress | Accepted |
| [022](022-stb-image-resize.md) | Resize faces with stb_image_resize2 | Accepted |
| [023](023-per-game-results-and-page-reuse.md) | Per-game result folders and page reuse | Accepted |
