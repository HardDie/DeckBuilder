#!/bin/bash

set -u
set -e
set -o pipefail

CHARACTERS="
	https://foursouls.com/wp-content/uploads/2021/10/tw-tapeworm-768x1047.png
"

ETERNAL="
	https://foursouls.com/wp-content/uploads/2021/10/tw-pink_proglottid-768x1047.png
"

TREASURE="
	https://foursouls.com/wp-content/uploads/2021/10/tw-black_proglottid-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/tw-red_proglottid-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/tw-white_proglottid-768x1047.png
"

LOOT="
	https://foursouls.com/wp-content/uploads/2021/10/tw-rainbow_tapeworm-768x1047.png
"

MONSTERS="
	https://foursouls.com/wp-content/uploads/2021/10/tw-tapeworm_monster-768x1047.png
"


mkdir -pv tapeworm
pushd tapeworm
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
