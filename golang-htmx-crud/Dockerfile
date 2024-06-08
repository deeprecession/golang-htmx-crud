# -----------------Build Environment-----------------------------------

FROM golang:alpine AS build-stage

LABEL stage=gobuilder

WORKDIR /app

COPY . .

RUN go mod download
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /htmx-golang-crud ./cmd

# ----------------------Deployment Environment------------------------------

FROM alpine

WORKDIR /

ADD templates templates
ADD assets assets
COPY --from=build-stage /htmx-golang-crud htmx-golang-crud

EXPOSE 42069

ENTRYPOINT ["/htmx-golang-crud"]
