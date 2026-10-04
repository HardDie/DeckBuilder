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
# Static libjpeg.a for github.com/pixiv/go-libjpeg. make sets the flags.
# Wails drops CC while generating bindings, so -ljpeg still reaches the linker.
# JPEG_LINK_DIR holds only libjpeg.a. LIBRARY_PATH makes -ljpeg use that archive.
# JPEG_LIBJPEG_A / JPEG_INCLUDE point at an archive built in CI.
# WAILS_BUILD_ARGS is for release CI (-skipbindings -platform …).
CGO_ENABLED := 1
REAL_CC := $(shell $(GO) env CC)
UNAME_S := $(shell uname -s)
JPEG_CC := $(CURDIR)/scripts/cc-static-jpeg
ifneq ($(filter MINGW% MSYS% CYGWIN%,$(UNAME_S)),)
JPEG_CC := $(shell cygpath -m "$(CURDIR)/scripts/cc-static-jpeg.exe" 2>/dev/null || echo "$(CURDIR)/scripts/cc-static-jpeg.exe")
endif
JPEG_EXTRA_LDFLAGS :=
ifeq ($(UNAME_S),Linux)
JPEG_EXTRA_LDFLAGS := -lm
endif
WAILS_BUILD_ARGS ?=
JPEG_LINK_DIR := $(CURDIR)/build/libjpeg-link
JPEG_LIBRARY_PATH := $(JPEG_LINK_DIR)
ifneq ($(filter MINGW% MSYS% CYGWIN%,$(UNAME_S)),)
JPEG_LIBRARY_PATH := $(shell cygpath -m "$(JPEG_LINK_DIR)" 2>/dev/null || echo "$(JPEG_LINK_DIR)")
endif
ifneq ($(JPEG_LIBJPEG_A),)
CGO_LDFLAGS := $(JPEG_LIBJPEG_A)
ifneq ($(JPEG_INCLUDE),)
CGO_CFLAGS := -I$(JPEG_INCLUDE)
endif
else
JPEG_BREW_PREFIX := $(shell brew --prefix jpeg-turbo 2>/dev/null)
ifneq ($(JPEG_BREW_PREFIX),)
CGO_CFLAGS := -I$(JPEG_BREW_PREFIX)/include
CGO_LDFLAGS := $(JPEG_BREW_PREFIX)/lib/libjpeg.a
else
ifeq ($(UNAME_S),Linux)
CGO_CFLAGS := $(shell pkg-config --cflags libjpeg 2>/dev/null)
CGO_LDFLAGS := $(shell ./scripts/find-libjpeg-a.sh)
else
ifneq ($(filter MINGW% MSYS% CYGWIN%,$(UNAME_S)),)
CGO_LDFLAGS := $(shell ./scripts/find-libjpeg-a.sh)
ifneq ($(CGO_LDFLAGS),)
CGO_CFLAGS := -I$(patsubst %/lib/libjpeg.a,%/include,$(CGO_LDFLAGS))
endif
endif
endif
endif
endif

.DEFAULT_GOAL := help

.PHONY: help version dev build jpeg-link require-jpeg generate test test-integration test-all \
	vet fmt tidy doc doc-all docs-site frontend-install screenshots clean ci \
	linter-install linter-run fuzz_game fuzz_collection fuzz_deck fuzz_card

ifneq ($(filter MINGW% MSYS% CYGWIN%,$(UNAME_S)),)
dev build: $(JPEG_CC)
$(JPEG_CC): scripts/cc-static-jpeg.c
	$(REAL_CC) -O2 -o $@ scripts/cc-static-jpeg.c
endif

dev build: jpeg-link

jpeg-link: require-jpeg
	mkdir -p "$(JPEG_LINK_DIR)"
	ln -sfn "$(CGO_LDFLAGS)" "$(JPEG_LINK_DIR)/libjpeg.a"

## help: Show this list
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## //' | awk -F': ' '{printf "  %-22s %s\n", $$1, $$2}'

## version: Set wails.json productVersion from the last git tag
version:
	PRODUCT_VERSION="$(PRODUCT_VERSION)" ./scripts/sync-product-version.sh

## dev: Run the Wails app with frontend hot reload
dev: require-wails require-jpeg version
	LIBRARY_PATH="$(JPEG_LIBRARY_PATH)$${LIBRARY_PATH:+:$$LIBRARY_PATH}" \
	CGO_ENABLED="$(CGO_ENABLED)" CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS) $(JPEG_EXTRA_LDFLAGS)" \
		CC="$(JPEG_CC)" REAL_CC="$(REAL_CC)" \
		$(WAILS) dev $(WAILS_TAGS) -ldflags "$(VERSION_LDFLAGS)"

## build: Production binary for this machine (build/bin)
build: require-wails require-jpeg version
	LIBRARY_PATH="$(JPEG_LIBRARY_PATH)$${LIBRARY_PATH:+:$$LIBRARY_PATH}" \
	CGO_ENABLED="$(CGO_ENABLED)" CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS) $(JPEG_EXTRA_LDFLAGS)" \
		CC="$(JPEG_CC)" REAL_CC="$(REAL_CC)" \
		$(WAILS) build $(WAILS_TAGS) $(WAILS_BUILD_ARGS) -clean -trimpath -ldflags "$(VERSION_LDFLAGS)"
	@if [ "$$(uname)" = "Linux" ]; then \
		cp build/linux/install.sh build/linux/$(APP_NAME).desktop build/bin/ && \
		cp build/appicon.png build/bin/$(APP_NAME).png && \
		chmod +x build/bin/install.sh; \
	fi

## generate: Regenerate frontend/wailsjs bindings from Go
generate: require-wails jpeg-link
	LIBRARY_PATH="$(JPEG_LIBRARY_PATH)$${LIBRARY_PATH:+:$$LIBRARY_PATH}" \
	CGO_ENABLED="$(CGO_ENABLED)" CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS) $(JPEG_EXTRA_LDFLAGS)" \
		CC="$(JPEG_CC)" REAL_CC="$(REAL_CC)" \
		$(WAILS) generate module $(WAILS_TAGS)

## test: Unit tests for package main, bindings, internal, and pkg (same as CI, with race)
test: jpeg-link
	LIBRARY_PATH="$(JPEG_LIBRARY_PATH)$${LIBRARY_PATH:+:$$LIBRARY_PATH}"; \
	export LIBRARY_PATH CGO_ENABLED="$(CGO_ENABLED)" CGO_CFLAGS="$(CGO_CFLAGS)" \
		CGO_LDFLAGS="$(CGO_LDFLAGS) $(JPEG_EXTRA_LDFLAGS)" \
		CC="$(JPEG_CC)" REAL_CC="$(REAL_CC)"; \
	$(GO) test -race -count=1 -tags=$(TEST_MAIN_TAGS) . && \
	$(GO) test -race -count=1 -tags=$(TEST_MAIN_TAGS) ./bindings/... && \
	$(GO) test -race -count=1 -tags=$(TEST_MAIN_TAGS) $(INTERNAL) && \
	$(GO) test -race -count=1 -tags=$(TEST_MAIN_TAGS) ./pkg/...

## test-integration: Integration tests under internal (build tag integration)
test-integration: jpeg-link
	LIBRARY_PATH="$(JPEG_LIBRARY_PATH)$${LIBRARY_PATH:+:$$LIBRARY_PATH}"; \
	export LIBRARY_PATH CGO_ENABLED="$(CGO_ENABLED)" CGO_CFLAGS="$(CGO_CFLAGS)" \
		CGO_LDFLAGS="$(CGO_LDFLAGS) $(JPEG_EXTRA_LDFLAGS)" \
		CC="$(JPEG_CC)" REAL_CC="$(REAL_CC)"; \
	$(GO) test -tags=integration -count=1 $(INTERNAL)

## test-all: Unit then integration tests
test-all: test test-integration

## ci: What GitHub Actions test.yml runs
ci: test-all

## vet: Go vet on package main (no Wails CGO), bindings, internal, and pkg
vet: jpeg-link
	LIBRARY_PATH="$(JPEG_LIBRARY_PATH)$${LIBRARY_PATH:+:$$LIBRARY_PATH}"; \
	export LIBRARY_PATH CGO_ENABLED="$(CGO_ENABLED)" CGO_CFLAGS="$(CGO_CFLAGS)" \
		CGO_LDFLAGS="$(CGO_LDFLAGS) $(JPEG_EXTRA_LDFLAGS)" \
		CC="$(JPEG_CC)" REAL_CC="$(REAL_CC)"; \
	$(GO) vet -tags=$(TEST_MAIN_TAGS) . && \
	$(GO) vet -tags=$(TEST_MAIN_TAGS) ./bindings/... && \
	$(GO) vet -tags=$(TEST_MAIN_TAGS) $(INTERNAL) && \
	$(GO) vet -tags=$(TEST_MAIN_TAGS) ./pkg/...

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

require-jpeg:
	@if [ ! -f "$(CGO_LDFLAGS)" ]; then \
		echo "install the static libjpeg package"; \
		exit 1; \
	fi
