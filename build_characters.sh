#!/bin/bash

set -u
set -e
set -o pipefail

echo "-> Start building characters"

# First line
echo "Start build first line"
convert \
	source/base_game/characters/b2-isaac-768x1047.png \
	source/base_game/characters/b2-cain-768x1047.png \
	source/base_game/characters/b2-maggy-768x1047.png \
	source/base_game/characters/b2-judas-768x1047.png \
	source/base_game/characters/b2-samson-768x1047.png \
	source/base_game/characters/b2-eve-768x1047.png \
	source/base_game/characters/b2-lilith-768x1047.png \
	source/base_game/characters/b2-blue_baby-768x1047.png \
	+append /tmp/mytmp/line1.png

# Second line
echo "Start build second line"
convert \
	source/base_game/characters/b2-lazarus-768x1047.png \
	source/base_game/characters/b2-the_forgotten-768x1047.png \
	source/base_game/characters/b2-eden-768x1047.png \
	source/base_game/characters/b2-the_forgotten-768x1047.png \
	source/gold_box/characters/g2-the_keeper-768x1047.png \
	source/gold_box/characters/g2-azazel-768x1047.png \
	source/gold_box/characters/g2-apollyon-768x1047.png \
	empty.png \
	+append /tmp/mytmp/line2.png

# Third line
echo "Start build third line"
convert \
	empty.png \
	empty.png \
	empty.png \
	source/shop_promos/characters/sp-eden-768x1047.png \
	source/pre_signup_promos/characters/psp-eden-768x1047.png \
	source/tapeworm/characters/tw-tapeworm-768x1047.png \
	empty.png \
	source/back/CharacterCardBack-768x1047.png \
	+append /tmp/mytmp/line3.png

# Build result image
echo "Start build result image"
convert \
	/tmp/mytmp/line1.png \
	/tmp/mytmp/line2.png \
	/tmp/mytmp/line3.png \
	-append result/characters_v2.png
