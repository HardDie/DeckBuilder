#!/bin/sh
# libwebp: the stored image format (ADR 027). Static libwebp.a and libsharpyuv.a.
# Built with libwebp's makefile.unix, so it needs only make and a C compiler.
# Setting EXTRA_FLAGS replaces the makefile's defaults: no threads (one image per
# call, and no pthread or winpthread to link), no PNG/JPEG/TIFF tool dependencies.
# Usage: build.sh <prefix>
dir=$(cd "$(dirname "$0")" && pwd)
NAME=libwebp
. "$dir/VERSION"
. "$dir/../lib.sh"
need make
need curl

build_one() {
  src=$1 out=$2 arch=$3
  flags=""
  case "$(uname -s)" in
    Linux) flags="-fPIC" ;;
  esac
  if [ -n "$arch" ]; then
    flags="-arch $arch -mmacosx-version-min=$MACOS_MIN"
  fi
  case "$(uname -s)" in
    MINGW* | MSYS* | CYGWIN*) cc=gcc ;;
    *) cc=${CC:-cc} ;;
  esac
  make -C "$src" -f makefile.unix -j"$(jobs)" CC="$cc" EXTRA_FLAGS="$flags" \
    src/libwebp.a sharpyuv/libsharpyuv.a >/dev/null
  mkdir -p "$out/lib" "$out/include/webp"
  cp "$src/src/libwebp.a" "$src/sharpyuv/libsharpyuv.a" "$out/lib/"
  # makefile.unix archives with "ar r"; the index lets every linker read them.
  ranlib "$out/lib/libwebp.a" "$out/lib/libsharpyuv.a"
  cp "$src/src/webp/decode.h" "$src/src/webp/encode.h" "$src/src/webp/types.h" "$out/include/webp/"
}

build_all "$1"
