#!/bin/sh
# Print the path of a static archive (for example libwebp.a) when Homebrew is not used.
# pkg-config --static for the library, then the usual lib dirs, then MinGW prefixes.
# Usage: find-static-lib.sh libwebp.a [pkg-config name]
lib=$1
pc=${2:-}
emit() {
  if command -v cygpath >/dev/null 2>&1; then
    cygpath -m "$1"
  else
    printf '%s\n' "$1"
  fi
}
if [ -n "$pc" ] && pkg-config --exists "$pc" 2>/dev/null; then
  libs=$(pkg-config --static --libs-only-L "$pc" 2>/dev/null || true)
  for tok in $libs; do
    dir=${tok#-L}
    if [ -f "$dir/$lib" ]; then
      emit "$dir/$lib"
      exit 0
    fi
  done
fi
arch=$(uname -m)
for dir in \
  /usr/lib \
  /usr/lib64 \
  /usr/local/lib \
  "/usr/lib/${arch}-linux-gnu" \
  /lib \
  /lib64 \
  /mingw64/lib \
  /ucrt64/lib \
  /clang64/lib \
  /c/msys64/mingw64/lib \
  /c/msys64/ucrt64/lib \
  /c/msys64/clang64/lib
do
  if [ -f "$dir/$lib" ]; then
    emit "$dir/$lib"
    exit 0
  fi
done
exit 0
