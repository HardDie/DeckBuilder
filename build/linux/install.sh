#!/bin/sh
# Install DeckBuilder for the current user. No root.
# Run this from the unpacked release folder.
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
NAME=DeckBuilder

BIN_SRC="$ROOT/$NAME"
ICON_SRC="$ROOT/$NAME.png"
DESKTOP_SRC="$ROOT/$NAME.desktop"

if [ ! -f "$BIN_SRC" ] || [ ! -f "$ICON_SRC" ] || [ ! -f "$DESKTOP_SRC" ]; then
  echo "Run install.sh from the unpacked release folder." >&2
  echo "Expected $NAME, $NAME.png, and $NAME.desktop next to this script." >&2
  exit 1
fi

BIN_DEST="$HOME/.local/bin/$NAME"
ICON_DEST="$HOME/.local/share/icons/hicolor/512x512/apps/$NAME.png"
DESKTOP_DEST="$HOME/.local/share/applications/$NAME.desktop"

mkdir -p \
  "$(dirname "$BIN_DEST")" \
  "$(dirname "$ICON_DEST")" \
  "$(dirname "$DESKTOP_DEST")"

cp "$BIN_SRC" "$BIN_DEST"
chmod 755 "$BIN_DEST"

cp "$ICON_SRC" "$ICON_DEST"
chmod 644 "$ICON_DEST"

# Line match avoids sed replacement metacharacters in $HOME.
while IFS= read -r line || [ -n "$line" ]; do
  case "$line" in
    Exec=__EXEC__) printf 'Exec=%s\n' "$BIN_DEST" ;;
    Icon=__ICON__) printf 'Icon=%s\n' "$ICON_DEST" ;;
    *) printf '%s\n' "$line" ;;
  esac
done < "$DESKTOP_SRC" > "$DESKTOP_DEST"
chmod 644 "$DESKTOP_DEST"

if command -v update-desktop-database >/dev/null 2>&1; then
  update-desktop-database "$HOME/.local/share/applications" || true
fi

echo "Installed $NAME for this user."
echo "Binary:   $BIN_DEST"
echo "Launcher: $DESKTOP_DEST"
echo "Icon:     $ICON_DEST"
echo "Open $NAME from the app list."
