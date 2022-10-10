#!/bin/bash

set -u
set -o pipefail
set -e

BACKEND=$(git --git-dir ../.git rev-parse --short HEAD)
FRONTEND=$(git --git-dir ../gui/.git rev-parse --short HEAD)

mkdir -p out

CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags \
	"-X main.BackendCommit=${BACKEND} -X main.FrontendCommit=${FRONTEND}" \
	-a -o DeckBuilder_amd64 -v ../cmd/deck_builder
tar -czf out/deckbuilder_linux_amd64.tar.gz DeckBuilder_amd64
rm DeckBuilder_amd64

CGO_ENABLED=0 GOARCH=386 GOOS=linux go build -ldflags \
	"-X main.BackendCommit=${BACKEND} -X main.FrontendCommit=${FRONTEND}" \
	-a -o DeckBuilder_386 -v ../cmd/deck_builder
tar -czf out/deckbuilder_linux_386.tar.gz DeckBuilder_386
rm DeckBuilder_386

CGO_ENABLED=0 GOARCH=arm GOOS=linux go build -ldflags \
	"-X main.BackendCommit=${BACKEND} -X main.FrontendCommit=${FRONTEND}" \
	-a -o DeckBuilder_arm -v ../cmd/deck_builder
tar -czf out/deckbuilder_linux_arm.tar.gz DeckBuilder_arm
rm DeckBuilder_arm
