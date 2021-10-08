#!/bin/bash

set -u
set -e
set -o pipefail

echo "-> Start building bonus souls"

# First line
echo "Start build first line"
convert \
	source/base_game/bonus_souls/b2-soul_of_greed-768x1047.png \
	source/base_game/bonus_souls/b2-soul_of_gluttony-768x1047.png \
	+append /tmp/mytmp/line1.png

# Second line
echo "Start build second line"
convert \
	source/base_game/bonus_souls/b2-soul_of_guppy-768x1047.png \
	source/back/BonusSoulCardBack-768x1047.png \
	+append /tmp/mytmp/line2.png

# Build result image
echo "Start build result image"
convert \
	/tmp/mytmp/line1.png \
	/tmp/mytmp/line2.png \
	-append result/bonus_souls_v2.png
