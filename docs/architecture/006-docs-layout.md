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
3. **README (users) + CURSOR.md (agents) + `docs/architecture` (ADRs) + `docs/use-cases` + `docs/wiki` (package notes).**

## Decision

Use option 3.

Keep README short ([Make a README](https://www.makeareadme.com/)): what it is, GUI link, TTS Saved Objects steps, build, swagger URL, license. Operator how-tos that outgrow the README can go in the wiki; HTTP contracts stay in CURSOR.md.

Use cases describe **behavior that exists in this repo**. Copy `_TEMPLATE.md`. Status is Implemented (or Planned only for agreed future work listed in the index).

Developer module notes (packages, fields, who calls what) live in **[docs/wiki](../wiki/Home.md)**. They explain current code; they do not replace ADRs. Copy `docs/wiki/*.md` (including `_Sidebar.md` and `_Footer.md`) into the GitHub wiki remote, same as [fsentry](https://github.com/HardDie/fsentry/tree/master/docs/wiki).

## Consequences

### Positive

* Users start at README; agents start at CURSOR.md; decisions are numbered.
* Use cases stay honest about current code (e.g. PATCH settings only persists `lang` today).

### Negative and risks

* Product changes may need README + CURSOR.md + an ADR + a use case.

### Neutral

* Wiki markdown uses GitHub wiki slugs (`[DB](DB)`, `Home.md`, `_Sidebar.md`), not Jekyll.
