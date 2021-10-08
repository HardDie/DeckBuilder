#!/bin/bash

set -u
set -e
set -o pipefail

BACKS="
	https://foursouls.com/wp-content/uploads/2021/10/CharacterCardBack-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/EternalCardBack-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/TreasureCardBack-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/LootCardBack-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/MonsterCardBack-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/BonusSoulCardBack-768x1047.png
"

mkdir -pv back
pushd back
	for item in ${BACKS}; do
		wget -nv ${item}
	done
popd
