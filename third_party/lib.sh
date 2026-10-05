# Shared by third_party/*/build.sh. Not run on its own.
#
# A build.sh sets NAME, VERSION, SHA256, and URL (from its VERSION file), defines
# build_one <src> <prefix> <arch>, and calls build_all <prefix>.
# build_one gets an empty arch for a native build, or arm64 / x86_64 on macOS.
set -eu

die() {
  echo "third_party/$NAME: $*" >&2
  exit 1
}

need() {
  command -v "$1" >/dev/null 2>&1 || die "$1 is missing. macOS: brew install cmake nasm; Ubuntu: apt install cmake nasm; MSYS2: pacman -S mingw-w64-x86_64-cmake mingw-w64-x86_64-nasm"
}

sha256() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | cut -d' ' -f1
  else
    shasum -a 256 "$1" | cut -d' ' -f1
  fi
}

jobs() {
  getconf _NPROCESSORS_ONLN 2>/dev/null || echo 4
}

# fetch <dir>: download URL once into the shared download folder, check SHA256,
# and unpack a fresh copy into <dir>.
fetch() {
  tarball="$DOWNLOADS/$NAME-$VERSION.tar.gz"
  if [ ! -f "$tarball" ] || [ "$(sha256 "$tarball")" != "$SHA256" ]; then
    mkdir -p "$DOWNLOADS"
    curl -fsSL -o "$tarball.part" "$URL"
    got=$(sha256 "$tarball.part")
    if [ "$got" != "$SHA256" ]; then
      rm -f "$tarball.part"
      die "checksum mismatch for $URL: got $got, want $SHA256"
    fi
    mv "$tarball.part" "$tarball"
  fi
  rm -rf "$1"
  mkdir -p "$1"
  # Read the archive from stdin: GNU tar (MSYS2) takes "D:/…" as host "D" and a remote path.
  tar -xzf - -C "$1" --strip-components=1 <"$tarball"
}

# build_all <prefix>: native build, or on macOS an arm64 and an x86_64 build
# joined with lipo, so one prefix serves every macOS target (universal included).
build_all() {
  prefix=$1
  DOWNLOADS="$(dirname "$prefix")/downloads"
  work="$(dirname "$prefix")/work/$NAME"
  rm -rf "$work"
  if [ "$(uname -s)" = "Darwin" ]; then
    for arch in arm64 x86_64; do
      fetch "$work/src-$arch"
      build_one "$work/src-$arch" "$work/out-$arch" "$arch"
    done
    mkdir -p "$prefix/lib" "$prefix/include"
    cp -R "$work/out-arm64/include/." "$prefix/include/"
    for lib in "$work"/out-arm64/lib/*.a; do
      name=$(basename "$lib")
      lipo -create "$lib" "$work/out-x86_64/lib/$name" -output "$prefix/lib/$name"
    done
  else
    fetch "$work/src"
    build_one "$work/src" "$work/out" ""
    mkdir -p "$prefix/lib" "$prefix/include"
    cp -R "$work/out/include/." "$prefix/include/"
    cp "$work"/out/lib/*.a "$prefix/lib/"
  fi
  rm -rf "$work"
  echo "third_party/$NAME $VERSION: $(cd "$prefix/lib" && ls)"
}

# The oldest macOS the app supports: Go 1.27's linker supports macOS 13 and later.
# deps.mk reads this line for the app's own cgo flags; keep Info.plist in step.
MACOS_MIN=13.0
