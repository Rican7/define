# Get the release name through Git via a sub-shell command
RELEASE_NAME = $(shell git describe --exact-match --abbrev=0 2>/dev/null)
COMMIT_HASH = $(shell git rev-parse --short HEAD)

# Define directories
ROOT_DIR ?= ${CURDIR}
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

# Validate
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

build: install-deps
	go build ${GO_BUILD_FLAGS}

install: install-deps
	go install ${GO_BUILD_FLAGS}

${VENDOR_DIR}: | check-dep
	dep ensure

install-deps: check-dep | ${VENDOR_DIR}

update-deps: check-dep
	dep ensure -update

test: install-deps
	go test -v ./...

test-with-coverage: install-deps
	go test -cover ./...



.PHONY: all check-dep clean-deps clean build install install-deps update-deps test test-with-coverage
