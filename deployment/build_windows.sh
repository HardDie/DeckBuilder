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
	-a -o ../../deployment/DeckBuilder.exe -v .

cd ../../deployment
zip out/deckbuilder.windows-amd64.zip DeckBuilder.exe
rm DeckBuilder.exe

cd ../cmd/deck_builder

CGO_ENABLED=0 GOARCH=386 GOOS=windows go build -ldflags \
	"-X main.BackendCommit=${BACKEND} -X main.FrontendCommit=${FRONTEND}" \
	-a -o ../../deployment/DeckBuilder.exe -v .
rm resource.syso

cd ../../deployment
zip out/deckbuilder.windows-386.zip DeckBuilder.exe
rm DeckBuilder.exe
