#!/bin/bash

set -u
set -e
set -o pipefail

CHARACTERS="
	https://foursouls.com/wp-content/uploads/2021/10/sp-eden-768x1047.png
"

MONSTERS="
	https://foursouls.com/wp-content/uploads/2021/10/sp-it_lives-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/so-the_bloat-768x1047.png
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

	if [[ `echo ${MONSTERS} | wc -w` -ge 1 ]]; then
		mkdir -pv monsters
		pushd monsters
			for item in ${MONSTERS}; do
				wget -nv ${item}
			done
		popd
	fi
popd
