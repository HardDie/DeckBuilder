#!/bin/bash

RED="\033[0;31m"
GRN="\033[0;32m"
BLU="\033[0;36m"
NC="\033[0m"

TOTAL=0

check_bin() {
	#$1 - binary file name
	#$2 - package description
	#$3 - possible package name
	which $1 > /dev/null 2>/dev/null
	RES=$?
	if [[ $RES -eq 1 ]]; then
		echo -e "[$RED FAIL $NC]    Can't find '$1' bin. $2\n$3"
		TOTAL=1
	else
		echo -e "[$GRN  OK  $NC]    $1"
	fi
}

echo "Required build packets:"
check_bin "png2icns" \
	"Required to build the MacOS icon pack." \
	"The package may have the name '${BLU}libicns${NC}', try to install it using the package manager."
check_bin "zip" \
	"Required to compress a binary file for windows." \
	"The package may have the name '${BLU}zip${NC}', try to install it using the package manager."
check_bin "tar" \
	"Required to compress a binary file for linux and MacOS." \
	"The package may have the name '${BLU}tar${NC}', try to install it using the package manager."
check_bin "yarn" \
	"Required to build vue web." \
	"The package may have the name '${BLU}yarn${NC}', try to install it using the package manager."
check_bin "goversioninfo" \
	"Required to create a bin file containing the icon and description for a compiled windows application." \
	"How to install: ${BLU}go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest${NC}"

exit ${TOTAL}
