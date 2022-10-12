#!/bin/bash

set -u
set -o pipefail
set -e

BACKEND=$(git --git-dir ../.git rev-parse --short HEAD)
FRONTEND=$(git --git-dir ../gui/.git rev-parse --short HEAD)

mkdir -p out

cd ../cmd/deck_builder
go generate

CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -ldflags \
	"-X main.BackendCommit=${BACKEND} -X main.FrontendCommit=${FRONTEND}" \
	-a -o ../../deployment/DeckBuilder_amd64.exe -v .

CGO_ENABLED=0 GOARCH=386 GOOS=windows go build -ldflags \
	"-X main.BackendCommit=${BACKEND} -X main.FrontendCommit=${FRONTEND}" \
	-a -o ../../deployment/DeckBuilder_386.exe -v .
rm resource.syso

cd ../../deployment

zip out/deckbuilder_windows_amd64.zip DeckBuilder_amd64.exe
zip out/deckbuilder_windows_386.zip DeckBuilder_386.exe

rm DeckBuilder_amd64.exe
rm DeckBuilder_386.exe
