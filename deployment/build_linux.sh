#!/bin/bash

set -u
set -o pipefail
set -e

BACKEND=$(git --git-dir ../.git rev-parse --short HEAD)
FRONTEND=$(git --git-dir ../gui/.git rev-parse --short HEAD)

mkdir -p out

CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags \
	"-X main.BackendCommit=${BACKEND} -X main.FrontendCommit=${FRONTEND}" \
	-a -o DeckBuilder -v ../cmd/deck_builder
tar -czf out/deckbuilder.linux-amd64.tar.gz DeckBuilder
rm DeckBuilder

CGO_ENABLED=0 GOARCH=386 GOOS=linux go build -ldflags \
	"-X main.BackendCommit=${BACKEND} -X main.FrontendCommit=${FRONTEND}" \
	-a -o DeckBuilder -v ../cmd/deck_builder
tar -czf out/deckbuilder.linux-386.tar.gz DeckBuilder
rm DeckBuilder

CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -ldflags \
	"-X main.BackendCommit=${BACKEND} -X main.FrontendCommit=${FRONTEND}" \
	-a -o DeckBuilder -v ../cmd/deck_builder
tar -czf out/deckbuilder.linux-arm64.tar.gz DeckBuilder
rm DeckBuilder
