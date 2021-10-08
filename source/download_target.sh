#!/bin/bash

set -u
set -e
set -o pipefail

CHARACTERS="\
"

ETERNAL="\
"

TREASURE="\
	https://foursouls.com/wp-content/uploads/2021/10/t-dead_eye-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/t-epic_fetus-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/t-marked-768x1047.png \
"

LOOT="\
"

MONSTERS="\
"

BONUS_SOULS="\
"


mkdir -pv target
pushd target
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
