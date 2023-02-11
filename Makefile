# Get the release name through Git via a sub-shell command
RELEASE_NAME = $(shell git describe --exact-match --abbrev=0 2>/dev/null)
COMMIT_HASH = $(shell git rev-parse --short HEAD)

# Define directories
ROOT_DIR ?= ${CURDIR}
TOOLS_DIR ?= ${ROOT_DIR}/tools
BUILD_DIR ?= ${ROOT_DIR}/.tmpbuild

# Set a local GOBIN to run our local tooling
export GOBIN ?= ${TOOLS_DIR}/bin

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
GO_CLEAN_FLAGS ?= -i -r -x ${GO_BUILD_FLAGS}

# Compilation flags
XC_ARCHITECTURES ?= 386 amd64 arm arm64
XC_OPERATING_SYSTEMS ?= linux darwin windows freebsd netbsd openbsd
XC_OSARCHS ?= !darwin/arm !darwin/arm64
GOX_BUILD_FLAGS ?= -verbose -ldflags="${GO_LD_FLAGS}" -arch="${XC_ARCHITECTURES}" -os="${XC_OPERATING_SYSTEMS}" -osarch="${XC_OSARCHS}" -output="${BUILD_DIR}/{{.Dir}}_{{.OS}}_{{.Arch}}"

# Tool flags
GOFMT_FLAGS ?= -s
GOIMPORTS_FLAGS ?=
GOLINT_MIN_CONFIDENCE ?= 0.3

# Set the mode for code-coverage
GO_TEST_COVERAGE_MODE ?= count
GO_TEST_COVERAGE_FILE_NAME ?= coverage.out

# Validate
ifndef BUILD_DIR
$(error BUILD_DIR must be set and non-empty)
endif
ifndef GOBIN
$(error GOBIN must be set and non-empty)
endif

# Global/default target
all: install-deps build install

clean:
	go clean ${GO_CLEAN_FLAGS} ./...

clean-release: clean
	rm -rf ${BUILD_DIR}

build: install-deps
	go build ${GO_BUILD_FLAGS}

${BUILD_DIR}: install-deps install-deps-dev
	gox ${GOX_BUILD_FLAGS}
	for file in ${BUILD_DIR}/* ; do sha256sum "$${file}" > "$${file}.sha256"; done

build-release: ${BUILD_DIR}

install: install-deps
	go install ${GO_BUILD_FLAGS}

install-deps:
	go mod download

tools install-deps-dev: install-deps
	cd tools && go install \
		golang.org/x/lint/golint \
		golang.org/x/tools/cmd/goimports \
		honnef.co/go/tools/cmd/staticcheck \
		github.com/mitchellh/gox

update-deps:
	go get ./...

test:
	go test -v ./...

test-with-coverage:
	go test -cover -covermode ${GO_TEST_COVERAGE_MODE} ./...

test-with-coverage-formatted:
	go test -cover -covermode ${GO_TEST_COVERAGE_MODE} ./... | column -t | sort -r

test-with-coverage-profile:
	go test -covermode ${GO_TEST_COVERAGE_MODE} -coverprofile ${GO_TEST_COVERAGE_FILE_NAME} ./...

format-lint:
	@errors=$$(gofmt -l ${GOFMT_FLAGS} .); if [ "$${errors}" != "" ]; then echo "Format lint failed on:\n$${errors}\n"; exit 1; fi

import-lint: install-deps-dev
	@errors=$$(${GOBIN}/goimports -l ${GOIMPORTS_FLAGS} .); if [ "$${errors}" != "" ]; then echo "Import lint failed on:\n$${errors}\n"; exit 1; fi

style-lint: install-deps-dev
	${GOBIN}/golint -min_confidence=${GOLINT_MIN_CONFIDENCE} -set_exit_status ./...
	${GOBIN}/staticcheck ./...

lint: install-deps-dev format-lint import-lint style-lint

vet:
	go vet ./...

format-fix:
	gofmt -w ${GOFMT_FLAGS} .

import-fix:
	${GOBIN}/goimports -w ${GOIMPORTS_FLAGS} .

fix: install-deps-dev format-fix import-fix
	go fix ./...


.PHONY: all clean clean-release build build-release install install-deps tools install-deps-dev update-deps test test-with-coverage format-lint import-lint style-lint lint vet format-fix import-fix fix
