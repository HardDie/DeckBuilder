#!/bin/bash

set -u
set -o pipefail
set -e

BACKEND=$(git --git-dir ../.git rev-parse --short HEAD)
FRONTEND=$(git --git-dir ../gui/.git rev-parse --short HEAD)

mkdir -p out

cp versioninfo.json ../cmd/deck_builder
cd ../cmd/deck_builder
goversioninfo -icon=../../deployment/win_icon.ico -64

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
rm versioninfo.json

cd ../../deployment
zip out/deckbuilder.windows-386.zip DeckBuilder.exe
rm DeckBuilder.exe
