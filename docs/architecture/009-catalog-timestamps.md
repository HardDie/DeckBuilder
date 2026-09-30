# 9. Normalize catalog createdAt and updatedAt in memory

* **Status:** Accepted
* **Date:** 2026-09-18
* **Authors:** @oleg

---

## Context

Preview builds stored `createdAt` / `updatedAt` inconsistently. Folder envelopes (`.info.json`) and card payloads could omit the fields, write JSON `null`, or leave an empty value. Each catalog repository had its own `convertCreateUpdate`: fill `createdAt` with now, copy it into `updatedAt` when that is also missing.

Current writes always store both timestamps. Older projects and imported zips can still contain the preview shape. Services expect non-zero `time.Time`.

## Considered options

1. **Keep four copies of `convertCreateUpdate`.**
2. **One in-memory helper** — fill missing times when mapping to entities; do not rewrite files.
3. **Rewrite `.info.json` on import or read** — persist the new format as a migration.

## Decision

Use option 2.

`utils.NormalizeTimestamps` is the only fill rule. Repositories apply it in `toEntity` before returning to services. Preview files stay as they are until a later create/update writes new JSON. There is no catalog timestamp migration.

New cards write `createdAt` and `updatedAt` together. Missing `updatedAt` still copies `createdAt` (not “now”).

## Consequences

### Positive

* Preview projects still list and sort; the GUI always sees timestamps.
* Timestamp handling is not copied per repository.
* Reads do not rewrite the catalog.

### Negative and risks

* On-disk preview JSON can remain null/empty until the user edits that object.

### Neutral

* Entities stay non-pointer `time.Time`.
* DTO timestamps are RFC3339 strings.
