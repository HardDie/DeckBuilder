# 6. Documentation layout: README, CURSOR.md, and `docs/`

* **Status:** Accepted
* **Date:** 2026-09-18
* **Authors:** @oleg

---

## Context

Agents need HTTP contracts, layer rules, and generate limits. Users need a short install and TTS copy path. Dumping both into README or chat makes the vision drift.

This layout follows the same split as [HardDie/ytmemchat_wails](https://github.com/HardDie/ytmemchat_wails) (`CURSOR.md` + ADRs + use cases), adapted to a mature backend rather than a port.

## Considered options

1. **Only README.**
2. **README + a single ARCHITECTURE.md at repo root.**
3. **README (users) + CURSOR.md (agents) + `docs/architecture` (ADRs) + `docs/use-cases`.**

## Decision

Use option 3.

Keep README short ([Make a README](https://www.makeareadme.com/)): what it is, GUI link, TTS Saved Objects steps, build, swagger URL, license. Operator detail that is already in README (render folder, copy JSON) stays there until someone adds a wiki.

Use cases describe **behavior that exists in this repo**. Copy `_TEMPLATE.md`. Status is Implemented (or Planned only for agreed future work listed in the index).

## Consequences

### Positive

* Users start at README; agents start at CURSOR.md; decisions are numbered.
* Use cases stay honest about current code (e.g. PATCH settings only persists `lang` today).

### Negative and risks

* Product changes may need README + CURSOR.md + an ADR + a use case.

### Neutral

* No `docs/wiki` unless operator how-tos outgrow the README.
