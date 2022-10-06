#!/bin/bash

set -u
set -o pipefail
set -e

mkdir -p out

cd ../cmd/deck_builder
go generate

CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -a -o ../../deployment/DeckBuilder_amd64.exe -v .
CGO_ENABLED=0 GOARCH=386 GOOS=windows go build -a -o ../../deployment/DeckBuilder_386.exe -v .
rm resource.syso

cd ../../deployment

zip out/deckbuilder_windows_amd64.zip DeckBuilder_amd64.exe
zip out/deckbuilder_windows_386.zip DeckBuilder_386.exe

rm DeckBuilder_amd64.exe
rm DeckBuilder_386.exe
