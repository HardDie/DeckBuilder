# 8. Require Go 1.27.1

* **Status:** Accepted
* **Date:** 2026-09-18
* **Authors:** @oleg

---

## Context

The module declared `go 1.19`. That language version is years behind the current stable toolchain (Go 1.27.1 as of 2026-09). Developers already run a modern `go`; the `go` line should match so `GOTOOLCHAIN` and CI use the same compiler.

## Considered options

1. **Keep `go 1.19`** — still builds on new toolchains, but blocks language/std features and leaves an outdated minimum.
2. **`go 1.27` (minor only)** — language 1.27, patch chosen by the local toolchain.
3. **`go 1.27.1`** — minimum toolchain is the current stable patch.

## Decision

Use option 3. `go.mod` has `go 1.27.1`. Bump this line when upgrading to a newer stable patch or minor.

## Consequences

### Positive

* `go test`, fuzz, and `go install` use a current compiler.
* Loop-var and other post-1.21 semantics apply.
* Go 1.27 rejects `*F` methods (`f.Fatal`) inside `f.Fuzz` callbacks; those tests must use `t`.

### Negative and risks

* Contributors must have Go 1.27.1+ (or let `GOTOOLCHAIN=auto` download it).

### Neutral

* Dependency versions were not required to jump with this change; `go mod tidy` keeps them consistent with 1.27.
