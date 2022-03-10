# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.17-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /shortly

##
## Deploy
##
FROM alpine:latest

WORKDIR /app

COPY --from=build /shortly /app/shortly
COPY templates /app/templates

EXPOSE 8000

USER nobody:nobody

ENTRYPOINT ["/app/shortly"]