# 4. Layered packages: api, servers, services, repositories

* **Status:** Accepted
* **Date:** 2026-09-18
* **Authors:** @oleg

---

## Context

The backend grew CRUD for four aggregates plus generate, search, replace, images, and system. Mixing mux, swagger, and disk in one package makes Swagger generation and tests harder.

## Considered options

1. **Handlers in `package main` / one `internal/http`.**
2. **Hexagonal ports with many small packages per verb.**
3. **Fixed layers used everywhere:** `internal/servers` (image and TTS GET), `internal/services` (rules), `internal/repositories` (disk).

## Decision

Use option 3. Wiring is only in `wire.go`. Each aggregate has `contract.go`, `service.go`, and `repository.go`. Image and TTS GET live in `internal/servers`.

`internal/entities` are in-process structs. `internal/dto` is the JSON the GUI sees. `internal/tts_entity` is the TTS Saved Object schema — not the same as catalog cards.

## Consequences

### Positive

* Services are testable without `net/http` (see fuzz tests under `internal/services/*`).
* New endpoints follow an existing file to copy.

### Negative and risks

* A one-line change can touch api + server + service + dto.
* `wire.go` is a long composition root by design.

### Neutral

* Helpers (`fs`, `images`, `page_drawer`, `progress`, `network`) are not a fifth “business” layer; generator depends on them from services.
