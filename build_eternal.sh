#!/bin/bash

set -u
set -e
set -o pipefail

echo "-> Start building eternals"

# First line
echo "Start build first line"
convert \
	source/base_game/eternal/b2-the_d6-768x1047.png \
	source/base_game/eternal/b2-sleight_of_hand-768x1047.png \
	source/base_game/eternal/b2-yum_heart-768x1047.png \
	source/base_game/eternal/b2-book_of_belial-768x1047.png \
	source/base_game/eternal/b2-blood_lust-768x1047.png \
	source/base_game/eternal/b2-the_curse-768x1047.png \
	source/base_game/eternal/b2-incubus-768x1047.png \
	+append /tmp/mytmp/line1.png

# Second line
echo "Start build second line"
convert \
	source/base_game/eternal/b2-forever_alone-768x1047.png \
	source/base_game/eternal/b2-lazarus_rags-768x1047.png \
	source/base_game/eternal/b2-the_bone-768x1047.png \
	source/gold_box/eternal/g2-holy_mantle-768x1047.png \
	source/gold_box/eternal/g2-wooden_nickel-768x1047.png \
	source/gold_box/eternal/g2-lord_of_the_pit-768x1047.png \
	source/gold_box/eternal/g2-void-768x1047.png \
	+append /tmp/mytmp/line2.png

# Third line
echo "Start build third line"
convert \
	empty.png \
	empty.png \
	empty.png \
	empty.png \
	source/tapeworm/eternal/tw-pink_proglottid-768x1047.png \
	empty.png \
	source/back/EternalCardBack-768x1047.png \
	+append /tmp/mytmp/line3.png

# Build result image
echo "Start build result image"
convert \
	/tmp/mytmp/line1.png \
	/tmp/mytmp/line2.png \
	/tmp/mytmp/line3.png \
	-append result/eternal_v2.png
