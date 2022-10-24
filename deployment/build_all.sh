#!/bin/bash

set -u
set -o pipefail
set -e

./check_binary.sh

rm -rf out || 1
./build_linux.sh
./build_darwin.sh
./build_windows.sh
