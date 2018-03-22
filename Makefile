# Get the release name through Git via a sub-shell command
RELEASE_NAME = $(shell git describe --exact-match --abbrev=0 2>/dev/null)
COMMIT_HASH = $(shell git rev-parse --short HEAD)

# Define directories
ROOT_DIR ?= ${CURDIR}
BUILD_DIR ?= ${ROOT_DIR}/.tmpbuild
VENDOR_DIR ?= ${ROOT_DIR}/vendor
GOPATH_FIRST ?= $(word 1, $(subst :, , ${GOPATH}))
GOBIN ?= ${GOPATH_FIRST}/bin

# Define program imports and variables for the compiler
APP_VERSION_IMPORT_PATH ?= github.com/Rican7/define/internal/version
APP_VERSION_ID_VAR ?= ${APP_VERSION_IMPORT_PATH}.identifier
APP_VERSION_COMMIT_HASH_VAR ?= ${APP_VERSION_IMPORT_PATH}.commitHash

# Linker flags
GO_LD_FLAGS += -X ${APP_VERSION_COMMIT_HASH_VAR}=${COMMIT_HASH}
ifneq (${RELEASE_NAME},)
GO_LD_FLAGS += -X ${APP_VERSION_ID_VAR}=${RELEASE_NAME}
endif

# Build flags
GO_BUILD_FLAGS ?= -ldflags "${GO_LD_FLAGS}" -v
GO_CLEAN_FLAGS ?= -i -x ${GO_BUILD_FLAGS}
GOX_BUILD_FLAGS ?= -verbose -ldflags="${GO_LD_FLAGS}" -output="${BUILD_DIR}/{{.Dir}}_{{.OS}}_{{.Arch}}"

# Tool flags
GOFMT_FLAGS ?= -s
GOIMPORTS_FLAGS ?=
GOLINT_MIN_CONFIDENCE ?= 0.3

# Validate
ifndef BUILD_DIR
$(error BUILD_DIR must be set and non-empty)
endif
ifndef GOBIN
$(error GOBIN must be set and non-empty)
endif

# Global/default target
all: install-deps build install

check-dep:
	@command -v dep &> /dev/null || (echo 'The `dep` command is not available. Download it at https://github.com/golang/dep' && exit 1)

clean-deps:
	rm -rf ${VENDOR_DIR}

clean: clean-deps
	go clean ${GO_CLEAN_FLAGS} ./...

clean-release: clean
	rm -rf ${BUILD_DIR}

build: install-deps
	go build ${GO_BUILD_FLAGS}

${BUILD_DIR}: install-deps install-deps-dev
	gox ${GOX_BUILD_FLAGS}

build-release: ${BUILD_DIR}

install: install-deps
	go install ${GO_BUILD_FLAGS}

${VENDOR_DIR}: | check-dep
	dep ensure

install-deps: check-dep | ${VENDOR_DIR}

install-deps-dev: install-deps
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/goimports
	go get github.com/mitchellh/gox

update-deps: check-dep
	dep ensure -update

update-deps-dev: update-deps
	go get -u github.com/golang/lint/golint
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/mitchellh/gox

test: install-deps
	go test -v ./...

test-with-coverage: install-deps
	go test -cover ./...

format-lint:
	@errors=$$(gofmt -l ${GOFMT_FLAGS} .); if [ "$${errors}" != "" ]; then echo "Format lint failed on:\n$${errors}\n"; exit 1; fi

import-lint:
	@errors=$$(goimports -l ${GOIMPORTS_FLAGS} .); if [ "$${errors}" != "" ]; then echo "Import lint failed on:\n$${errors}\n"; exit 1; fi

style-lint:
	@errors=$$(golint -min_confidence=${GOLINT_MIN_CONFIDENCE} $$(go list ./... | grep -v /vendor/)); if [ "$${errors}" != "" ]; then echo "Style lint failed on:\n$${errors}\n"; exit 1; fi

lint: install-deps-dev format-lint import-lint style-lint

format-fix:
	gofmt -w ${GOFMT_FLAGS} .

import-fix:
	goimports -w ${GOIMPORTS_FLAGS} .

fix: install-deps-dev format-fix import-fix
	go fix ./...

vet:
	go vet ./...


.PHONY: all check-dep clean-deps clean clean-release build build-release install install-deps install-deps-dev update-deps update-deps-dev test test-with-coverage format-lint import-lint style-lint lint format-fix import-fix fix vet
