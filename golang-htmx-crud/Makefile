
.PHONY: test all clean

MAIN_PACKAGE_PATH := ./cmd/main.go
VERSION := v1.0
BINARY_PATH := ./golang-htmx-crud-${VERSION}

## tidy: fmt and tidy
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v


## audit: mod verify, vet, golangci-lint, staticcheck, govulncheck, test -race -buildvcs
.PHONY: audit
audit:
	go mod verify
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run


## build: builds a package
.PHONY: build
build: tidy
	go build -o=${BINARY_PATH} ${MAIN_PACKAGE_PATH}


## run: builds and runs a package
.PHONY: run
run: build
	${BINARY_PATH}


## air: run `air` for current project
.PHONY: air
air: tidy
	air


## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
