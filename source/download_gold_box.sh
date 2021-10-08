#!/bin/bash

set -u
set -e
set -o pipefail

CHARACTERS="\
	https://foursouls.com/wp-content/uploads/2021/10/g2-apollyon-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-azazel-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-the_keeper-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-the_lost-768x1047.png \
"

ETERNAL="\
	https://foursouls.com/wp-content/uploads/2021/10/g2-holy_mantle-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-lord_of_the_pit-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-void-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-wooden_nickel-768x1047.png \
"

TREASURE="\
	https://foursouls.com/wp-content/uploads/2021/10/g2-9_volt-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-crooked_penny-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-dads_key-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-fruitcake-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-guppys_tail-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-i_cant_believe_its_not_butter_bean-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-infamy-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-lemon_mishap-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-library_card-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-moms_knife-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-more_options-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-ouija_board-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-placenta-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-plan_c-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-skeleton_key-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-soy_milk-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-succubus-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-the_bible-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-the_butter_bean-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-the_missing_page-768x1047.png \
"

LOOT="\
	https://foursouls.com/wp-content/uploads/2021/10/g2-a_penny-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-a_sack-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-bomb-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-cancer-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-charged_penny-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-credit_card-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-four_cents-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-holy_card-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-jera-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-joker-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-pills-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-pink_eye-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-soul_heart-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-three_cents-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-two_cents-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-two_of_diamonds-768x1047.png \
"

MONSTERS="\
	https://foursouls.com/wp-content/uploads/2021/10/g2-begotten-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-boil-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-charger-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-curse_of_fatigue-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-curse_of_tiny_hands-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-deaths_head-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-fistula-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-gaper-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-gurglings-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-i_am_error-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-imp-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-knight-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-parabite-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-ragling-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-round_worm-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-the_cage-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-trap_door-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-hush-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-polycephalus-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/g2-steven-768x1047.png \
"

BONUS_SOULS="\
"


mkdir -pv gold_box
pushd gold_box
	if [[ `echo ${CHARACTERS} | wc -w` -ge 1 ]]; then
		mkdir -pv characters
		pushd characters
			for item in ${CHARACTERS}; do
				wget -nv ${item}
			done
		popd
	fi

	if [[ `echo ${ETERNAL} | wc -w` -ge 1 ]]; then
		mkdir -pv eternal
		pushd eternal
			for item in ${ETERNAL}; do
				wget -nv ${item}
			done
		popd
	fi

	if [[ `echo ${TREASURE} | wc -w` -ge 1 ]]; then
		mkdir -pv treasure
		pushd treasure
			for item in ${TREASURE}; do
				wget -nv ${item}
			done
		popd
	fi

	if [[ `echo ${LOOT} | wc -w` -ge 1 ]]; then
		mkdir -pv loot
		pushd loot
			for item in ${LOOT}; do
				wget -nv ${item}
			done
		popd
	fi

	if [[ `echo ${MONSTERS} | wc -w` -ge 1 ]]; then
		mkdir -pv monsters
		pushd monsters
			for item in ${MONSTERS}; do
				wget -nv ${item}
			done
		popd
	fi

	if [[ `echo ${BONUS_SOULS} | wc -w` -ge 1 ]]; then
		mkdir -pv bonus_souls
		pushd bonus_souls
			for item in ${BONUS_SOULS}; do
				wget -nv ${item}
			done
		popd
	fi
popd
