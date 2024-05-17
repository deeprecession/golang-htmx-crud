PROJECT_NAME := htmx-golang-crud
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)

.PHONY: all dep build clean test lint

all: run

dep: ## Get the dependencies
	@go mod download
	@go mod tidy

lint: dep ## Lint the files
	@golangci-lint run ./...

build: dep ## Build the binary file
	@go build -v -o $(PROJECT_NAME) ./cmd/main.go

clean: ## Remove previous build
	@rm -f $(PROJECT_NAME)

build-app-image:
	docker build --rm --no-cache -t "$(PROJECT_NAME)-app" .

run-app-image:
	docker run --rm --label "$(PROJECT_NAME)-app" "$(PROJECT_NAME)-app":latest -p 5432:5432

build-db-image:
	docker build --rm --no-cache -t "$(PROJECT_NAME)-db" -f ./docker/postgres/Dockerfile .

build-prometheus-image:
	docker build --rm --no-cache -t "$(PROJECT_NAME)-prometheus" -f ./docker/prometheus/Dockerfile .

build-all-images: build-app-image build-db-image build-prometheus-image

run:
	docker-compose up --build
