#!/usr/bin/env bash
# Write the last git tag (no leading v) into wails.json info.productVersion.
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$root"

tag="${PRODUCT_VERSION:-}"
if [ -z "$tag" ]; then
  tag="$(git describe --tags --abbrev=0 2>/dev/null || true)"
fi
tag="${tag#v}"
if [ -z "$tag" ]; then
  tag="0.0.0"
fi

node -e '
const fs = require("fs");
const path = "wails.json";
const p = JSON.parse(fs.readFileSync(path, "utf8"));
const v = process.argv[1];
p.info = Object.assign({}, p.info, { productVersion: v });
const next = JSON.stringify(p, null, 2) + "\n";
if (fs.readFileSync(path, "utf8") !== next) {
  fs.writeFileSync(path, next);
}
' "$tag"
