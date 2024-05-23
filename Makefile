PROJECT_NAME := golang-htmx-crud


## up-db: `docker-compose up` a db image
.PHONY: up-db
up-db:
	docker-compose up db -d


## run-app: start golang app with ./.env enviroment
.PHONY: run-app
run-app: up-db
	./.env
	$(MAKE) -C ./golang-htmx-crud run


## air: run `air` for ./golang-htmx-crud with ./.env enviroment
.PHONY: air
air: up-db
	@set -a; \
		. ./.env; \
		$(MAKE) -C ./golang-htmx-crud air


## up: run docker-compose up --build
.PHONY: up
up:
	docker-compose up --build


## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

