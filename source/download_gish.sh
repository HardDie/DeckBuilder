#!/bin/bash

set -u
set -e
set -o pipefail

TREASURE="
	https://foursouls.com/wp-content/uploads/2021/07/gi-lil_gish-768x1047.png
"

MONSTERS="
	https://foursouls.com/wp-content/uploads/2021/10/gi-gish-768x1047.png
"


mkdir -pv gish
pushd gish
	if [[ `echo ${TREASURE} | wc -w` -ge 1 ]]; then
		mkdir -pv treasure
		pushd treasure
			for item in ${TREASURE}; do
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
