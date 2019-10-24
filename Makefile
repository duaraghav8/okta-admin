# Metadata about this makefile and position
MKFILE_PATH := $(lastword $(MAKEFILE_LIST))
CURRENT_DIR := $(patsubst %/,%,$(dir $(realpath $(MKFILE_PATH))))

# List all our actual files, excluding vendor
GOFILES ?= $(shell go list $(TEST) | grep -v /vendor/)

# Tags specific for building
GOTAGS ?=

TESTARGS ?=

# Number of procs to use
GOMAXPROCS ?= 4

# Get the project metadata
GOVERSION := 1.13.1
PROJECT := github.com/duaraghav8/okta-admin
NAME := $(notdir $(PROJECT))
GIT_COMMIT ?= $(shell git rev-parse --short HEAD)
VERSION := $(shell awk -F\" '/Version/ { print $$2; exit }' "${CURRENT_DIR}/version/version.go")

# Current system information
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# Default os-arch combination to build
XC_OS ?= darwin linux windows
XC_ARCH ?= amd64

# GPG Signing key (blank by default, means no GPG signing)
GPG_KEY ?=

# List of ldflags
LD_FLAGS ?= \
	-s \
	-w \
	-X ${PROJECT}/version.GitCommit=${GIT_COMMIT} \
	-extldflags \"-static\"

# List of tests to run
TEST ?= ./...

# Create a cross-compile target for every os-arch pairing. This will generate
# a make target for each os/arch like "make linux/amd64" as well as generate a
# meta target (build) for compiling everything.
define make-xc-target
  $1/$2:
	@printf "%s%20s %s\n" "-->" "${1}/${2}:" "${PROJECT}"
	@docker run \
		--interactive \
		--rm \
		--dns="8.8.8.8" \
		--volume="${CURRENT_DIR}:/go/src/${PROJECT}" \
		--workdir="/go/src/${PROJECT}" \
		"golang:${GOVERSION}" \
		env \
			CGO_ENABLED="0" \
			GOOS="${1}" \
			GOARCH="${2}" \
			go build \
			  -a \
				-o="_build/${NAME}${3}_${1}_${2}" \
				-ldflags "${LD_FLAGS}" \
				-tags "${GOTAGS}"
  .PHONY: $1/$2

  $1:: $1/$2
  .PHONY: $1

  build:: $1/$2
  .PHONY: build
endef
$(foreach goarch,$(XC_ARCH),$(foreach goos,$(XC_OS),$(eval $(call make-xc-target,$(goos),$(goarch),$(if $(findstring windows,$(goos)),.exe,)))))

# dist builds the binaries and then signs and packages them for distribution
dist:
ifndef GPG_KEY
	@echo "==> ERROR: No GPG key specified! Without a GPG key, this release cannot"
	@echo "           be signed. Set the environment variable GPG_KEY to the ID of"
	@echo "           the GPG key to continue."
	@exit 127
else
	@$(MAKE) -f "${MKFILE_PATH}" _cleanup
	@$(MAKE) -f "${MKFILE_PATH}" -j4 build
	@$(MAKE) -f "${MKFILE_PATH}" _checksum _sign
endif
.PHONY: dist

# _cleanup removes any previous binaries
_cleanup:
	@rm -rf "${CURRENT_DIR}/_build/"

# _checksum produces the checksums for the binaries in _build
_checksum:
	@cd "${CURRENT_DIR}/_build" && \
		shasum --algorithm 256 * > ${CURRENT_DIR}/_build/${NAME}_${VERSION}_SHA256SUMS && \
		cd - &>/dev/null
.PHONY: _checksum

# _sign signs the binaries using the given GPG_KEY. This should not be called
# as a separate function.
_sign:
	@echo "==> Signing ${PROJECT} at v${VERSION}"
	@gpg \
		--default-key "${GPG_KEY}" \
		--detach-sig "${CURRENT_DIR}/_build/${NAME}_${VERSION}_SHA256SUMS"
	@git commit \
		--allow-empty \
		--gpg-sign="${GPG_KEY}" \
		--message "Release v${VERSION}" \
		--quiet \
		--signoff
	@git tag \
		--annotate \
		--create-reflog \
		--local-user "${GPG_KEY}" \
		--message "Version ${VERSION}" \
		--sign \
		"v${VERSION}" master
	@echo "--> Do not forget to run:"
	@echo ""
	@echo "    git push && git push --tags"
	@echo ""
	@echo "And then upload the binaries in dist/!"
.PHONY: _sign

# test runs the test suite.
test: fmtcheck
	@echo "==> Testing ${NAME}"
	@go test -v -timeout=300s -tags="${GOTAGS}" ${GOFILES} ${TESTARGS}
.PHONY: test

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./main.go ./meta.go ./okta ./command ./version

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/fmtcheck.sh'"
.PHONY: fmtcheck
