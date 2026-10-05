#!/bin/sh
# Build static libwebp (libwebp.a, libsharpyuv.a) from the pinned release (ADR 027).
# It uses libwebp's makefile.unix, so it needs only make, a C compiler, curl, and tar.
# Setting EXTRA_FLAGS replaces the makefile's defaults: no threads (DeckBuilder encodes
# one image per call), no PNG/JPEG/TIFF tool dependencies, no shared library.
#
# Usage: build-libwebp.sh <prefix> [compiler flags]
#   e.g. build-libwebp.sh /tmp/webp-x64 "-arch x86_64 -mmacosx-version-min=12.0"
# CC picks the compiler (default cc). Installs <prefix>/lib/libwebp.a,
# <prefix>/lib/libsharpyuv.a, and <prefix>/include/webp/{decode,encode,types}.h.
set -eu
ver=1.6.0
prefix=$1
flags=${2:-}
cc=${CC:-cc}

work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
curl -fsSL -o "$work/libwebp.tar.gz" \
  "https://github.com/webmproject/libwebp/archive/refs/tags/v${ver}.tar.gz"
mkdir "$work/src"
tar -xzf "$work/libwebp.tar.gz" -C "$work/src" --strip-components=1

jobs=$(getconf _NPROCESSORS_ONLN 2>/dev/null || echo 4)
make -C "$work/src" -f makefile.unix -j"$jobs" CC="$cc" \
  EXTRA_FLAGS="$flags" \
  src/libwebp.a sharpyuv/libsharpyuv.a

mkdir -p "$prefix/lib" "$prefix/include/webp"
cp "$work/src/src/libwebp.a" "$work/src/sharpyuv/libsharpyuv.a" "$prefix/lib/"
# makefile.unix archives with "ar r"; the index lets every linker read them.
ranlib "$prefix/lib/libwebp.a" "$prefix/lib/libsharpyuv.a"
cp "$work/src/src/webp/decode.h" "$work/src/src/webp/encode.h" "$work/src/src/webp/types.h" \
  "$prefix/include/webp/"
echo "libwebp ${ver}: $prefix/lib/libwebp.a $prefix/lib/libsharpyuv.a"
