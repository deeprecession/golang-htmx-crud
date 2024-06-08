# golang-htmx-crud

## Setup Instructions

1. Change default database credentials variables in a `.env` file in the project root directory.

```
DB_USER="user"
DB_PASSWORD="pswd"

REDIS_PASSWORD="pswd"

```

2. To build and run the project use `make` or `docker-compose up --build`

## Services

`:8080` webapp server

- GET `/`
- GET `/register`
- GET `/login`
- GET `/metrics` statistics for Prometheus
- POST `/tasks`
- PUT `/tasks/:id`
- DELETE `/tasks/:id`

`::9090` Prometheus client

`::5432` Postgres database

`::6379` Redis database
