#!/bin/bash

set -u
set -e
set -o pipefail

# Prepare
mkdir -p result
convert -size 768x1047 "xc:#000000" empty.png

./build_characters.sh
./build_bonus_souls.sh
./build_eternal.sh
