#!/bin/bash

set -u
set -e
set -o pipefail

./download_back.sh
./download_base_game.sh
./download_gold_box.sh
./download_four_souls_plus.sh
# requiem
# warp
# big goi
./download_gish.sh
./download_target.sh
./download_shop_promos.sh
./download_pre_signup_promos.sh
./download_tapeworm.sh
