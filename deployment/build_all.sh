#!/bin/bash

set -u
set -o pipefail
set -e

VERSION=$(go run ../pkg/version/cmd/version)
RELEASE=$(git --git-dir ../.git describe --tags --always)

rm -rf release || 1

goreleaser build --name 'DeckBuilder' \
	--company 'org.harddie.deckbuilder' \
	--image '512.png' \
	--license 'Licensed under GPLv3.' \
	--version "${RELEASE}" \
	--ldflags "-X github.com/HardDie/DeckBuilder/pkg/version.Build=${VERSION}" \
	--path '../cmd/deck_builder/main.go'
