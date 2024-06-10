# todolist-htmx-golang

A simplest tasklist web application as a show case for HTMX in Golang

## Technologies Used
- **Go (Golang)**: backend
- **echo**: just wanted to try it
- **HTMX**: Library for making AJAX requests and updating parts of a web page
- **Tailwind CSS**: I often heard about it so I decided to try it and I like how easy it is to use in small projects
- **Docker**: Containerization for PostgreSQL, Redis, Prometheus, and WebApp
- **PostgreSQL**: stores registered users and their tasks
- **Redis**: caches sessions
- **Prometheus**: was added just because
- **Make**: contain many commands that are I'm too lazy to write again

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
