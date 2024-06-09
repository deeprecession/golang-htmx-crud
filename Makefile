PROJECT_NAME := todolist-htmx-golang


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
	$(MAKE) -C ./$(PROJECT_NAME) run


## air: run `air` for project with ./.env enviroment
.PHONY: air
air: up-db up-redis
	@set -a; \
		. ./.env; \
		$(MAKE) -C ./$(PROJECT_NAME) air


## audit: run `audit` target for golang project
.PHONY: audit
audit:
	cd ./$(PROJECT_NAME) && $(MAKE) audit


## up: run docker-compose up --build
.PHONY: up
up:
	docker-compose up --build


## tailwind-build: create output.css from input.css in ./templates/styles
.PHONY: tailwind-build
tailwind-build:
	npx tailwindcss build -i ./$(PROJECT_NAME)/assets/css/tailwind.css -o ./$(PROJECT_NAME)/assets/css/style.css


## tailwind-build-watch: create output.css from input.css in ./templates/styles and watches for changes
.PHONY: tailwind-build-watch
tailwind-build-watch:
	npx tailwindcss build -i ./$(PROJECT_NAME)/assets/css/tailwind.css -o ./$(PROJECT_NAME)/assets/css/style.css --watch


## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

