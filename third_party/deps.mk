# C libraries in third_party/: pinned, checksummed, built static from source on
# every machine (ADR 028). The root Makefile includes this file.
#
# `make deps` builds what is missing into build/third_party/<os>_<arch>/;
# every cgo target depends on it. A library is rebuilt when the content of its
# build.sh, VERSION, or lib.sh changes. macOS always builds arm64 + x86_64.
#
# Linking stays static without wrapper scripts:
#   - pixiv/go-libjpeg asks for -ljpeg. -L points at a folder that holds only
#     libjpeg.a, and that folder is searched before the system ones.
#   - libwebp is named by path, so nothing can pick a shared libwebp.

DEPS_OS := $(shell $(GO) env GOOS)
DEPS_ARCH := $(if $(filter darwin,$(DEPS_OS)),universal,$(shell $(GO) env GOARCH))
# DEPS_REL names make targets (no drive letter on Windows); DEPS_DIR goes to compilers.
DEPS_REL := build/third_party/$(DEPS_OS)_$(DEPS_ARCH)
DEPS_DIR := $(CURDIR)/$(DEPS_REL)
ifneq ($(filter MINGW% MSYS% CYGWIN%,$(shell uname -s)),)
DEPS_DIR := $(shell cygpath -m "$(DEPS_DIR)")
endif
DEPS_LIBS := libjpeg-turbo libwebp

# The oldest macOS the app supports; lib.sh holds it, so the libraries match.
# Wails adds -mmacosx-version-min=10.13 only when the flags carry none.
DEPS_MACOS_MIN := $(shell sed -n 's/^MACOS_MIN=//p' third_party/lib.sh)
DEPS_OS_FLAGS := $(if $(filter darwin,$(DEPS_OS)),-mmacosx-version-min=$(DEPS_MACOS_MIN))

export CGO_ENABLED := 1
export CGO_CFLAGS := -I$(DEPS_DIR)/include $(DEPS_OS_FLAGS)
export CGO_LDFLAGS := -L$(DEPS_DIR)/lib $(DEPS_DIR)/lib/libwebp.a $(DEPS_DIR)/lib/libsharpyuv.a $(DEPS_OS_FLAGS) $(if $(filter linux,$(DEPS_OS)),-lm)

# A library is current when its marker names the checksum of its build inputs.
# Content, not file times: a CI cache restores old times next to a fresh checkout.
deps_sum = $(shell cat third_party/lib.sh third_party/$(1)/build.sh third_party/$(1)/VERSION | cksum | cut -d' ' -f1)
deps_stamp = $(DEPS_REL)/.built-$(1)-$(call deps_sum,$(1))

define deps_rule
$(call deps_stamp,$(1)):
	rm -f $(DEPS_REL)/.built-$(1)-*
	third_party/$(1)/build.sh "$(DEPS_DIR)" >&2
	touch "$$@"
endef
$(foreach lib,$(DEPS_LIBS),$(eval $(call deps_rule,$(lib))))

.PHONY: deps cgo-env

## deps: Build the static C libraries (third_party/) for this machine
deps: $(foreach lib,$(DEPS_LIBS),$(call deps_stamp,$(lib)))

# For one go test outside make: eval "$(make -s cgo-env)"
## cgo-env: Print the cgo exports for a single go test (see third_party/README.md)
cgo-env: deps
	@echo 'export CGO_ENABLED=1'
	@echo "export CGO_CFLAGS='$(CGO_CFLAGS)'"
	@echo "export CGO_LDFLAGS='$(CGO_LDFLAGS)'"
