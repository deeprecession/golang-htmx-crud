PROJECT_NAME := golang-htmx-crud


## up-db: `docker-compose up` a db image
.PHONY: up-db
up-db:
	docker-compose up db -d


## up-redis: `docker-compose up` a redis image
.PHONY: up-redis
up-redis:
	docker-compose up redis -d


## run-app: start golang app with ./.env enviroment
.PHONY: run-app
run-app: up-db
	./.env
	$(MAKE) -C ./golang-htmx-crud run


## air: run `air` for ./golang-htmx-crud with ./.env enviroment
.PHONY: air
air: up-db up-redis
	@set -a; \
		. ./.env; \
		$(MAKE) -C ./golang-htmx-crud air


## audit: run `audit` target for golang project
.PHONY: audit
audit:
	cd ./golang-htmx-crud && $(MAKE) audit


## up: run docker-compose up --build
.PHONY: up
up:
	docker-compose up --build


## tailwind-build: create output.css from input.css in ./templates/styles
.PHONY: tailwind-build
tailwind-build:
	npx tailwindcss build -i ./golang-htmx-crud/assets/css/tailwind.css -o ./golang-htmx-crud/assets/css/style.css


## tailwind-build-watch: create output.css from input.css in ./templates/styles and watches for changes
.PHONY: tailwind-build-watch
tailwind-build-watch:
	npx tailwindcss build -i ./golang-htmx-crud/assets/css/tailwind.css -o ./golang-htmx-crud/assets/css/style.css --watch


## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

