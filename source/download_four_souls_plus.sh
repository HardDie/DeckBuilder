#!/bin/bash

set -u
set -e
set -o pipefail

CHARACTERS="
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-bum_bo-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-dark_judas-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-guppy-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-whore_of_babylon-768x1047.png
"

ETERNAL="
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-bag_o_trash-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-dark_arts-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-gimpy-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-infestation-768x1047.png
"

TREASURE="
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-1_up-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-20_20-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-abaddon-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-athame-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-black_candle-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-cursed_eye-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-daddy_long_legs-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-distant_admiration-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-divorce_papers-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-euthanasia-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-forget_me_now-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-game_breaking_bug-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-guppys_eye-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-head_of_krampus-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-head_of_the_keeper-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-hourglass-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-lard-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-libra-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-magnet-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-mama_haunt-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-moms_eye_shadow-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-mutant_spider-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-phd-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-polyphemus-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-rainbow_baby-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-red_candle-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-rubber_cement-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-smart_fly-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-telepathy_for_dummies-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-the_wiz-768x1047.png
"

LOOT="
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-questionmark_card-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-two_cents-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-three_cents-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-four_cents-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-nickel-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-a_penny-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-aaa_battery-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-ansuz-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-black_rune-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-bomb-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-butter_bean-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-dice_shard-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-get_out_of_jail_card-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-gold_key-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-lil_battery-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-perthro-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-pills-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-pills_2-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-pills_3-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-poker_chip-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-tape_worm-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-the_left_hand-768x1047.png
"

MONSTERS="
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-angel_room-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-blastocyst-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-bony-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-boss_rush-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-brain-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-curse_of_blood_lust-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-curse_of_impulse-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-cursed_globin-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-cursed_tumor-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-dingle-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-flaming_hopper-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-globin-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-head_trauma-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-headless_horseman-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-holy_bony-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-holy_chest-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-holy_mulligan-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-isaac-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-krampus-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-moms_heart-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-monstro_ii-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-nerve_ending-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-roundy-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-spiked_chest-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-sucker-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-swarmer-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-the_fallen-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-troll_bombs-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-tumor-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/fsp2-widow-768x1047.png
"

mkdir -pv four_souls_plus
pushd four_souls_plus
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
popd
