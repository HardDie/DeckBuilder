#!/bin/bash

set -u
set -o pipefail
set -e

VERSION=$(git --git-dir ../.git describe --tags --always)

rm -rf release || 1

goreleaser build --name 'DeckBuilder' \
	--company 'org.harddie.deckbuilder' \
	--image '512.png' \
	--license 'Licensed under GPLv3.' \
	--version "${VERSION}" \
	--ldflags "-X main.Version=${VERSION}" \
	--path '../cmd/deck_builder/main.go'
