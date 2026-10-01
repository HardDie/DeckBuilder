# 13. HTTP serves images and TTS; Wails hosts the UI

* **Status:** Accepted
* **Date:** 2026-09-29
* **Authors:** @oleg

---

## Context

1. The window is Wails ([ADR 010](010-wails-shell.md)).
2. Catalog verbs are bindings ([ADR 011](011-wails-bindings.md)).
3. The mux still served `web/dist`, SPA fallbacks, and `/docs`.
4. Startup opened a browser.

## Considered options

1. **Keep the embedded SPA** next to the Wails window.
2. **Drop Vue hosting on HTTP.**
   1. Keep image GETs.
   2. Keep `GET /api/tts/data`.

## Decision

Use option 2.

1. No embedded GUI.
2. No browser launch.
3. No `/docs` or `/swagger.json` route.
4. Loopback HTTP stays on `127.0.0.1:5000` (next port, 20 tries).
5. Routes that stay:
   1. `GET /api/.../image`
   2. `GET /api/tts/data`
6. Replace is `bindings/replace`.
7. `AssetServer.Handler` still forwards `/api` so `<img>` can load faces.
8. `main.go` starts that HTTP server.

This supersedes the SPA embed and browser in [ADR 001](001-loopback-http-embedded-spa.md).
Loopback listen stays.
This supersedes the “production HTTP+SPA” line in [ADR 010](010-wails-shell.md).

## Consequences

### Positive

1. One UI: the Wails window.
2. TTS and image URLs stay HTTP.

### Negative and risks

1. `make web-build` no longer feeds a served GUI.

### Neutral

1. Image reads stay `cachedImage` ([ADR 012](012-card-image-urls.md)).
2. Vite in `wails dev` still proxies GET `/api` to `:5000`.
