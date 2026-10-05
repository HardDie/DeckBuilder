# third_party

Every C library the app links, and how it is built.
See [ADR 028](../docs/architecture/028-third-party-c-libraries.md).

## Libraries

| Library | Version | Used by | Built with |
|---|---|---|---|
| libjpeg-turbo | 3.2.0 | sheet JPEG encode (`github.com/pixiv/go-libjpeg`) | cmake (+ nasm on x86) |
| libwebp | 1.6.0 | stored images (`internal/images/webp`, ADR 027) | its `makefile.unix` |
| stb_image_resize2 | 2.18 | sheet resize (`third_party/stbresize`) | none: a Go package; cgo compiles the vendored header |

## How it works

1. `make deps` (every cgo target depends on it) builds into `build/third_party/<os>_<arch>/`.
   1. macOS is always `darwin_universal`: an arm64 and an x86_64 build joined with `lipo`.
   2. macOS objects target 13.0, the app's minimum (`MACOS_MIN` in `lib.sh`; keep `Info.plist` in step).
2. Each tarball is downloaded once into `build/third_party/downloads` and checked against `SHA256`.
3. Only static archives are built: no shared libraries, tools, or tests.
4. A library is rebuilt when its `build.sh`, `VERSION`, or `lib.sh` changes in content.
5. `deps.mk` exports `CGO_ENABLED`, `CGO_CFLAGS`, and `CGO_LDFLAGS` for every recipe.
   1. libjpeg is linked as `-ljpeg` (go-libjpeg asks for it) from a folder that holds only `libjpeg.a`.
   2. libwebp is linked by path.
6. For one `go test` outside make: `eval "$(make -s cgo-env)"`.

## Bump a version

1. Edit `VERSION`: the new version and the tarball's sha256 (`shasum -a 256` or `sha256sum`).
2. Run `make deps`; the library rebuilds because its `VERSION` changed.
3. Run `make test`, and for libjpeg the render goldens.

## Add a library

1. A folder with `VERSION` (`VERSION`, `SHA256`, `URL`) and `build.sh`.
2. `build.sh` sources `lib.sh`, defines `build_one <src> <out> <arch>`, and calls `build_all "$1"`.
3. Add its name to `DEPS_LIBS` in `deps.mk` and its archives to `CGO_LDFLAGS`.
