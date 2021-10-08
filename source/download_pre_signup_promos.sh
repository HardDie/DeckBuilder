#!/bin/bash

set -u
set -e
set -o pipefail

CHARACTERS="
	https://foursouls.com/wp-content/uploads/2021/10/psp-eden-768x1047.png
"

MONSTERS="
	https://foursouls.com/wp-content/uploads/2021/10/psp-corrupted_data-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/psp-the_bloat-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/psp-ultra_pride-768x1047.png
"

mkdir -pv pre_signup_promos
pushd pre_signup_promos
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
