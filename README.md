# sna-project

## Setup Instructions

1. Change default database credentials variables in a `.env` file in the project root directory.

```
DB_USER="default"
DB_PASSWORD="default"
```

2. To build and run the project use `make` or `docker-compose up --build`

## Services

`::8080` balancer for webapp

`app_[1-3]:42069` webapp server

- GET `/`
- GET `/metrics` statistics for Prometheus
- POST `/tasks`
- PUT `/tasks/:id`
- DELETE `/tasks/:id`

`::9090` Prometheus client

`::5432` Postgres database
