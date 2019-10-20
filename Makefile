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
PROJECT := $(CURRENT_DIR)
OWNER := duaraghav8
NAME := $(notdir $(PROJECT))
GIT_COMMIT ?= $(shell git rev-parse --short HEAD)
IMPORT_PATH := github.com/${OWNER}/${NAME}
EXTERNAL_TOOLS ?=

# Current system information
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# Default os-arch combination to build
XC_OS ?= darwin linux windows
XC_ARCH ?= amd64

# List of ldflags
LD_FLAGS ?= \
	-s \
	-w \
	-X ${IMPORT_PATH}/version.GitCommit=${GIT_COMMIT} \
	-extldflags \"-static\"

# List of tests to run
TEST ?= ./ ./command/

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

native:
	@CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -a -o="_build/deployer" -ldflags "${LD_FLAGS}" -tags "${GOTAGS}"
.PHONY: native

# bootstrap installs the necessary go tools for development or build.
bootstrap:
	@echo "==> Bootstrapping ${PROJECT}"
	@for t in ${EXTERNAL_TOOLS}; do \
		echo "--> Installing $$t" ; \
		go get -u "$$t"; \
	done
.PHONY: bootstrap

# test runs the test suite.
test: fmtcheck
	@echo "==> Testing ${NAME}"
	@go test -v -timeout=300s -tags="${GOTAGS}" ${GOFILES} ${TESTARGS}
.PHONY: test

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./main.go ./meta.go ./okta/ ./command/ ./version/

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/fmtcheck.sh'"
.PHONY: fmtcheck
