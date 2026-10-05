#!/bin/sh
# libjpeg-turbo: the JPEG encoder for render sheets (github.com/pixiv/go-libjpeg).
# Static only, no TurboJPEG API, no tools or tests. REQUIRE_SIMD fails the build
# instead of quietly making a slow library when nasm is missing.
# Usage: build.sh <prefix>
dir=$(cd "$(dirname "$0")" && pwd)
NAME=libjpeg-turbo
. "$dir/VERSION"
. "$dir/../lib.sh"
need cmake
need curl

build_one() {
  src=$1 out=$2 arch=$3
  set -- -DENABLE_SHARED=OFF -DENABLE_STATIC=ON \
    -DWITH_TURBOJPEG=OFF -DWITH_TOOLS=OFF -DWITH_TESTS=OFF \
    -DREQUIRE_SIMD=ON -DCMAKE_POSITION_INDEPENDENT_CODE=ON \
    -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX="$out" -DCMAKE_INSTALL_LIBDIR=lib
  case "$(uname -s)" in
    MINGW* | MSYS* | CYGWIN*) set -- "$@" -G "MSYS Makefiles" ;;
    *) set -- "$@" -G "Unix Makefiles" ;;
  esac
  if [ -n "$arch" ]; then
    set -- "$@" -DCMAKE_OSX_ARCHITECTURES="$arch" -DCMAKE_OSX_DEPLOYMENT_TARGET="$MACOS_MIN"
  fi
  if [ "$arch" = "x86_64" ]; then
    # nasm writes no macOS version into its objects unless told, and the linker
    # then warns "no platform load command". -P pre-includes this pragma in every .asm.
    inc="$src/macos-version.inc"
    printf '%%pragma macho build_version macos, %s, %s\n' "${MACOS_MIN%%.*}" "${MACOS_MIN#*.}" >"$inc"
    set -- "$@" -DCMAKE_ASM_NASM_FLAGS="-P$inc"
  fi
  if [ "$arch" != "arm64" ] && [ "$(uname -m)" != "aarch64" ] && [ "$(uname -m)" != "arm64" ]; then
    need nasm # x86 SIMD is assembly; arm64 uses intrinsics
  fi
  cmake -S "$src" -B "$src/build" "$@" >/dev/null
  cmake --build "$src/build" --target install --parallel "$(jobs)" >/dev/null
}

build_all "$1"
