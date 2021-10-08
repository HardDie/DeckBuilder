#!/bin/bash

set -u
set -e
set -o pipefail

TREASURE="\
	https://foursouls.com/wp-content/uploads/2021/10/t-dead_eye-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/t-epic_fetus-768x1047.png \
	https://foursouls.com/wp-content/uploads/2021/10/t-marked-768x1047.png \
"

mkdir -pv target
pushd target
	if [[ `echo ${TREASURE} | wc -w` -ge 1 ]]; then
		mkdir -pv treasure
		pushd treasure
			for item in ${TREASURE}; do
				wget -nv ${item}
			done
		popd
	fi
popd
