# Builder
FROM golang:1.18-alpine as builder

RUN         apk add gcc g++ make

WORKDIR     /app

COPY        go.mod ./

COPY        go.sum ./

RUN         go mod download

COPY        . ./

# RUN         CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -v -o server cmd/main.go

# 'CGO_ENABLED=0', go-sqlite3 requires cgo to work.
# RUN         go build -ldflags="-w -s" -v -o server cmd/main.go

RUN         CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -v -o server main.go

# Server
# FROM  scratch # x509: certificate signed by unknown authority
FROM alpine:latest

WORKDIR /app

LABEL       MAINTAINER="Zero<tobewhatwewant@gmail.com>"

ARG         VERSION=v1

COPY        --from=builder /app/server /app/server

EXPOSE      8080

ENV         GIN_MODE=release

ENV         VERSION=${VERSION}

CMD  ["/app/server", "-c", "/conf/config.yml"]