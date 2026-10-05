# 28. C libraries built from source in third_party

* **Status:** Accepted
* **Date:** 2026-10-05
* **Authors:** @oleg

---

## Context

1. The app links two C libraries statically: libjpeg-turbo and libwebp (ADR 027).
2. Each came from a different place on each machine.
   1. libjpeg: Homebrew, the distro package, MSYS2, or a CI source build.
   2. libwebp: Homebrew or the distro package locally, a source build in CI.
3. Dev, CI tests, and releases could link different versions.
4. The Makefile carried an `if` chain per library and per source.
5. libjpeg needed a link wrapper (`scripts/cc-static-jpeg`, plus a C copy for Windows).
   1. go-libjpeg asks for `-ljpeg`; the wrapper dropped it.
   2. A symlink folder and `LIBRARY_PATH` covered `wails generate`, which drops `CC`.
6. CI downloaded source tarballs without checking them.

## Considered options

1. **Keep package managers, tidy the Makefile**
   1. Lost. Versions still differ between machines.
2. **Source builds in CI only**
   1. Lost. Dev and release still link different code.
3. **Source builds everywhere, pinned and checksummed**
   1. Won. One version and one path on every machine.

## Decision

1. Layout
   1. `third_party/<library>/VERSION`: version, sha256, and URL of the tarball.
   2. `third_party/<library>/build.sh`: builds static archives into a prefix.
   3. `third_party/lib.sh`: download, checksum, unpack, and the macOS slices.
   4. `third_party/deps.mk`: included by the root Makefile.
2. Libraries
   1. libjpeg-turbo 3.2.0 with cmake: static only, no TurboJPEG API, tools, or tests.
   2. `REQUIRE_SIMD` is on, so a missing nasm fails the build.
   3. libwebp 1.6.0 with its `makefile.unix`: static, no threads.
   4. `stb_image_resize2.h` and its cgo wrapper are a Go package: `third_party/stbresize`.
3. Build
   1. `make deps` builds into `build/third_party/<os>_<arch>/`.
   2. Every cgo target depends on it.
   3. macOS always builds arm64 and x86_64 and joins them with `lipo`.
   4. macOS objects target 13.0, the app's minimum, including nasm's (a pre-included pragma).
   5. Tarballs are kept in `build/third_party/downloads` and checked on every use.
   6. A checksum mismatch stops the build.
4. Rebuild
   1. A library rebuilds when the content of its `build.sh`, `VERSION`, or `lib.sh` changes.
   2. Its marker file names the checksum of those files, not their time.
   3. A CI cache restores old file times, so times would rebuild every run.
5. Linking
   1. `deps.mk` exports `CGO_ENABLED`, `CGO_CFLAGS`, and `CGO_LDFLAGS`.
   2. `-L build/third_party/…/lib` comes first and the folder holds only static archives.
   3. go-libjpeg's `-ljpeg` therefore finds `libjpeg.a`.
   4. libwebp is linked by path.
   5. The `CC` wrapper, its C copy, the symlink folder, and `LIBRARY_PATH` are removed.
6. Tools
   1. cmake and nasm (x86 SIMD), make, a C compiler, curl, tar.
   2. macOS: `brew install cmake nasm`; Ubuntu: `apt install cmake nasm`.
   3. Windows: MSYS2 MINGW64 `make`, `gcc`, `cmake`, `nasm`.
7. CI
   1. Release and test workflows install the tools and cache `build/third_party`.
   2. The cache key covers every file in `third_party/`.
   3. Windows runs `make deps` in the MSYS2 shell.
   4. The shared-library check after `make build` stays.
8. One `go test` outside make: `eval "$(make -s cgo-env)"`.

## Consequences

### Positive

1. The same library versions on every machine and in every release.
2. Downloads are checked against a pinned sha256.
3. The Makefile has no per-library or per-OS search logic.
4. No link wrapper; `wails generate` needs no special path.
5. A version bump is one edit to `VERSION`.

### Negative and risks

1. Developers need cmake and nasm.
2. The first `make` after a clone or a bump builds the libraries (about 30 s on an M4).
3. Security fixes in the libraries arrive only when `VERSION` is bumped.

### Neutral

1. The universal macOS libraries are built on every Mac, even for a native build.
