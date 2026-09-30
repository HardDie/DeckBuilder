# Exact git tag when HEAD is tagged, otherwise the 12-character commit hash.
# Override with BUILD_VERSION=… on the make command line or in the environment.
ifndef BUILD_VERSION
BUILD_VERSION = $(shell go run ./pkg/version/cmd/version 2>/dev/null || echo dev)
endif
VERSION_LDFLAGS = -X github.com/HardDie/DeckBuilder/pkg/version.Build=$(BUILD_VERSION)
# Ubuntu 24.04+ ships webkit2gtk-4.1; Wails needs this tag instead of 4.0.
WAILS_TAGS := $(shell pkg-config --exists webkit2gtk-4.1 2>/dev/null && echo -tags webkit2_41)

.PHONY: default
default: help

.PHONY: help
help: ## help information about make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## build the Wails app for this machine
	cd desktop && wails build $(WAILS_TAGS) -clean -trimpath -ldflags "$(VERSION_LDFLAGS)"

.PHONY: wails-dev
wails-dev: ## run Wails desktop app (window + image/TTS HTTP)
	cd desktop && wails dev $(WAILS_TAGS) -ldflags "$(VERSION_LDFLAGS)"

.PHONY: swagger
swagger: ## generate swagger json file
	./swagger generate spec -m -o web/swagger.json

.PHONY: swagger-install-linux
swagger-install-linux: ## install swagger for linux
	curl -o swagger -L https://github.com/go-swagger/go-swagger/releases/download/v0.29.0/swagger_linux_amd64 && chmod +x swagger

.PHONY: linter-install
linter-install: ## install linters
	curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b ./bin
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin

.PHONY: linter-run
linter-run: ## run linters
	./bin/gosec -fmt=sonarqube ./... || echo "gosec found issues"
	./bin/golangci-lint run

.PHONY: test
test: ## run unit tests
	go test ./... -v -race

.PHONY: fuzz_game
fuzz_game: ## run fuzzy test for game API's
	cd internal/services/game && go test -fuzz=FuzzGame -v

.PHONY: fuzz_collection
fuzz_collection: ## run fuzzy test for collection API's
	cd internal/services/collection && go test -fuzz=FuzzCollection -v

.PHONY: fuzz_deck
fuzz_deck: ## run fuzzy test for deck API's
	cd internal/services/deck && go test -fuzz=FuzzDeck -v

.PHONY: fuzz_card
fuzz_card: ## run fuzzy test for card API's
	cd internal/services/card && go test -fuzz=FuzzCard -v
