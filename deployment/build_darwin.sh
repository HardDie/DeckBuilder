#!/bin/bash

set -u
set -o pipefail
set -e

BACKEND=$(git --git-dir ../.git rev-parse --short HEAD)
FRONTEND=$(git --git-dir ../gui/.git rev-parse --short HEAD)

mkdir -p out

CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -ldflags \
	"-X main.BackendCommit=${BACKEND} -X main.FrontendCommit=${FRONTEND}" \
	-a -o DeckBuilder_amd64 -v ../cmd/deck_builder
mkdir tmp
cp -r DeckBuilder.app tmp/
rm tmp/DeckBuilder.app/Contents/MacOS/put_binary_here
mv DeckBuilder_amd64 tmp/DeckBuilder.app/Contents/MacOS/DeckBuilder
cd tmp
tar -czf ../out/deckbuilder_darwin_amd64.tar.gz DeckBuilder.app
cd ..
rm -rf tmp

CGO_ENABLED=0 GOARCH=arm64 GOOS=darwin go build -ldflags \
	"-X main.BackendCommit=${BACKEND} -X main.FrontendCommit=${FRONTEND}" \
	-a -o DeckBuilder_arm64 -v ../cmd/deck_builder
mkdir tmp
cp -r DeckBuilder.app tmp/
rm tmp/DeckBuilder.app/Contents/MacOS/put_binary_here
mv DeckBuilder_arm64 tmp/DeckBuilder.app/Contents/MacOS/DeckBuilder
cd tmp
tar -czf ../out/deckbuilder_darwin_arm64.tar.gz DeckBuilder.app
cd ..
rm -rf tmp
