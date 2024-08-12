#!/bin/bash

set -u
set -o pipefail
set -e

BACKEND=$(git --git-dir ../.git rev-parse --short HEAD)
FRONTEND=$(git --git-dir ../gui/.git rev-parse --short HEAD)
TAG=$(git --git-dir ../.git describe --tags)

BACKEND=
FRONTEND=
VERSION="v1.0.2"

rm -rf release || 1

goreleaser build --name 'DeckBuilder' \
	--company 'org.harddie.deckbuilder' \
	--image '512.png' \
	--license 'Licensed under GPLv3.' \
	--version "${TAG}" \
	--ldflags "-X main.BackendCommit=${BACKEND} -X main.FrontendCommit=${FRONTEND} -X main.Version=${VERSION}" \
	--path '../cmd/deck_builder/main.go'
