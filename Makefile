PROJECT_NAME := htmx-golang-crud

build-app-image:
	docker build --rm --no-cache -t "$(PROJECT_NAME)-app" ./golang-htmx-crud

run-app-image:
	docker run --rm --label "$(PROJECT_NAME)-app" "$(PROJECT_NAME)-app":latest -p 5432:5432

build-db-image:
	docker build --rm --no-cache -t "$(PROJECT_NAME)-db" -f ./docker/postgres/Dockerfile .

build-prometheus-image:
	docker build --rm --no-cache -t "$(PROJECT_NAME)-prometheus" -f ./docker/prometheus/Dockerfile .

build-all-images: build-app-image build-db-image build-prometheus-image

run:
	docker-compose up --build
