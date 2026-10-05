# Common developer commands. Run `make help` for the list.

APP_NAME     := DeckBuilder
INTERNAL     := ./internal/...
# Exclude the Wails entrypoint (CGO) so package main tests run on CI.
TEST_MAIN_TAGS := nomain
GO           := go
WAILS        := wails
PKG          ?= ./pkg/version
PKGSITE_ADDR ?= localhost:8081
# Exact git tag when HEAD is tagged, otherwise the short commit (pkg/version).
# Override with BUILD_VERSION=… on the make command line or in the environment.
ifndef BUILD_VERSION
BUILD_VERSION = $(shell $(GO) run ./pkg/version/cmd/version 2>/dev/null || echo dev)
endif
VERSION_LDFLAGS = -X github.com/HardDie/DeckBuilder/pkg/version.Build=$(BUILD_VERSION)
# Last git tag without a leading v (macOS CFBundle / Windows ProductVersion). Override with PRODUCT_VERSION=…
PRODUCT_VERSION ?= $(patsubst v%,%,$(shell git describe --tags --abbrev=0 2>/dev/null))
# Ubuntu 24.04+ ships webkit2gtk-4.1; Wails needs this tag instead of 4.0.
WAILS_TAGS := $(shell pkg-config --exists webkit2gtk-4.1 2>/dev/null && echo -tags webkit2_41)
# WAILS_BUILD_ARGS is for release CI (-skipbindings -platform …).
WAILS_BUILD_ARGS ?=
UNAME_S := $(shell uname -s)
# make dev builds through GO_OPTIMIZED: it drops the -gcflags "all=-N -l" that wails dev adds.
GO_OPTIMIZED := $(CURDIR)/build/go-optimized/go-optimized
ifneq ($(filter MINGW% MSYS% CYGWIN%,$(UNAME_S)),)
GO_OPTIMIZED := $(shell cygpath -m "$(CURDIR)/build/go-optimized/go-optimized.exe" 2>/dev/null || echo "$(CURDIR)/build/go-optimized/go-optimized.exe")
endif

# Static C libraries (libjpeg-turbo, libwebp) and the cgo flags: see third_party/.
include third_party/deps.mk

.DEFAULT_GOAL := help

.PHONY: help version dev dev-debug build generate test test-integration test-all \
	vet fmt tidy doc doc-all docs-site frontend-install screenshots clean ci \
	linter-install linter-run fuzz_game fuzz_collection fuzz_deck fuzz_card

$(GO_OPTIMIZED): scripts/go-optimized/main.go
	$(GO) build -o "$@" ./scripts/go-optimized

## help: Show this list
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -h -E '^## ' $(MAKEFILE_LIST) | sed 's/^## //' | awk -F': ' '{printf "  %-22s %s\n", $$1, $$2}'

## version: Set wails.json productVersion from the last git tag
version:
	PRODUCT_VERSION="$(PRODUCT_VERSION)" ./scripts/sync-product-version.sh

## dev: Run the Wails app with frontend hot reload; Go is optimized like make build
dev: require-wails deps version $(GO_OPTIMIZED)
	REAL_GO="$(GO)" $(WAILS) dev $(WAILS_TAGS) -compiler "$(GO_OPTIMIZED)" -ldflags "$(VERSION_LDFLAGS)"

## dev-debug: Like dev, but Go without optimization or inlining (-N -l), for a debugger
dev-debug: require-wails deps version
	$(WAILS) dev $(WAILS_TAGS) -ldflags "$(VERSION_LDFLAGS)"

## build: Production binary for this machine (build/bin)
build: require-wails deps version
	$(WAILS) build $(WAILS_TAGS) $(WAILS_BUILD_ARGS) -clean -trimpath -ldflags "$(VERSION_LDFLAGS)"
	@if [ "$$(uname)" = "Linux" ]; then \
		cp build/linux/install.sh build/linux/$(APP_NAME).desktop build/bin/ && \
		cp build/appicon.png build/bin/$(APP_NAME).png && \
		chmod +x build/bin/install.sh; \
	fi

## generate: Regenerate frontend/wailsjs bindings from Go
generate: require-wails deps
	$(WAILS) generate module $(WAILS_TAGS)

## test: Unit tests for package main, bindings, internal, pkg, and third_party (same as CI, with race)
test: deps
	$(GO) test -race -count=1 -tags=$(TEST_MAIN_TAGS) . && \
	$(GO) test -race -count=1 -tags=$(TEST_MAIN_TAGS) ./bindings/... && \
	$(GO) test -race -count=1 -tags=$(TEST_MAIN_TAGS) $(INTERNAL) && \
	$(GO) test -race -count=1 -tags=$(TEST_MAIN_TAGS) ./pkg/... && \
	$(GO) test -race -count=1 -tags=$(TEST_MAIN_TAGS) ./third_party/...

## test-integration: Integration tests under internal (build tag integration)
test-integration: deps
	$(GO) test -tags=integration -count=1 $(INTERNAL)

## test-all: Unit then integration tests
test-all: test test-integration

## ci: What GitHub Actions test.yml runs
ci: test-all

## vet: Go vet on package main (no Wails CGO), bindings, internal, pkg, and third_party
vet: deps
	$(GO) vet -tags=$(TEST_MAIN_TAGS) . && \
	$(GO) vet -tags=$(TEST_MAIN_TAGS) ./bindings/... && \
	$(GO) vet -tags=$(TEST_MAIN_TAGS) $(INTERNAL) && \
	$(GO) vet -tags=$(TEST_MAIN_TAGS) ./pkg/... && \
	$(GO) vet -tags=$(TEST_MAIN_TAGS) ./third_party/...

## fmt: Format Go files; fail if any file needed formatting
fmt:
	@files="$$(gofmt -l $$(find . -name '*.go' -not -path './frontend/*'))"; \
	if [ -n "$$files" ]; then echo "$$files"; echo "run: gofmt -w on the files above"; exit 1; fi

## tidy: go mod download and tidy
tidy:
	$(GO) mod download
	$(GO) mod tidy

## doc: Package summary (PKG=./pkg/version)
doc:
	$(GO) doc $(PKG)

## doc-all: Package + all exports (PKG=./pkg/version)
doc-all:
	$(GO) doc -all $(PKG)

## docs-site: HTML godoc at PKGSITE_ADDR (default localhost:8081)
docs-site:
	$(GO) run golang.org/x/pkgsite/cmd/pkgsite@latest -http $(PKGSITE_ADDR)

## frontend-install: yarn install in frontend/
frontend-install: require-wails
	yarn --cwd frontend install

## screenshots: Regenerate wiki guide images from the Vue window
screenshots:
	npm --prefix scripts/screenshots install
	npm --prefix scripts/screenshots exec -- playwright install chromium
	node scripts/screenshots/capture.mjs

## clean: Remove Wails/Go build artifacts
clean:
	$(GO) clean
	rm -rf build/bin frontend/dist

## linter-install: Install gosec and golangci-lint into ./bin
linter-install:
	curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b ./bin
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin

## linter-run: Run gosec and golangci-lint
linter-run:
	./bin/gosec -fmt=sonarqube ./... || echo "gosec found issues"
	./bin/golangci-lint run

## fuzz_game: Run fuzz test for the game API
fuzz_game:
	cd internal/services/game && go test -fuzz=FuzzGame -v

## fuzz_collection: Run fuzz test for the collection API
fuzz_collection:
	cd internal/services/collection && go test -fuzz=FuzzCollection -v

## fuzz_deck: Run fuzz test for the deck API
fuzz_deck:
	cd internal/services/deck && go test -fuzz=FuzzDeck -v

## fuzz_card: Run fuzz test for the card API
fuzz_card:
	cd internal/services/card && go test -fuzz=FuzzCard -v

require-wails:
	@test -f wails.json || (echo "wails.json missing: scaffold with wails init before this target"; exit 1)

