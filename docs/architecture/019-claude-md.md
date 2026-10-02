# 19. CLAUDE.md replaces CURSOR.md

* **Status:** Accepted
* **Date:** 2026-10-02
* **Authors:** @oleg

---

## Context

1. [ADR 006](006-docs-layout.md) put the agent spec in `CURSOR.md`.
2. `.cursor/rules/short-docs.mdc` set the line style for it.
3. Agents now run in Claude Code.
4. Claude Code loads `CLAUDE.md`, not `CURSOR.md`.
5. Two agent files would drift.

## Considered options

1. **Keep CURSOR.md, point CLAUDE.md at it**
   1. Lost. Agents read two files for one spec.
2. **Merge CURSOR.md into CLAUDE.md**
   1. Won. One agent file, loaded automatically.

## Decision

Use option 2.

1. `CLAUDE.md` holds the agent and contributor spec.
2. `CURSOR.md` is removed.
3. `.cursor/rules/` is removed.
4. The short-line rule now lives in the global `~/.claude/CLAUDE.md`.
   1. It governs `CLAUDE.md` and ADRs.
5. The rest of ADR 006 stands.
   1. README for users.
   2. ADRs, use cases, and wiki under `docs/`.

## Consequences

### Positive

1. One agent file.
2. Claude Code loads it without a pointer.

### Negative and risks

1. Cursor no longer gets project rules.

### Neutral

1. Older ADRs still name `CURSOR.md`.
2. Read that name as `CLAUDE.md`.
