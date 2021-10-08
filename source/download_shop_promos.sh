#!/bin/bash

set -u
set -e
set -o pipefail

CHARACTERS="\
	https://foursouls.com/wp-content/uploads/2021/10/sp-eden-768x1047.png \
"

ETERNAL="\
"

TREASURE="\
"

LOOT="\
"

MONSTERS="\
	https://foursouls.com/wp-content/uploads/2021/10/sp-it_lives-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/so-the_bloat-768x1047.png \
"

BONUS_SOULS="\
"


mkdir -pv shop_promos
pushd shop_promos
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
