# syntax=docker/dockerfile:1

##
## Build stage.
##
FROM golang:1.21.3-alpine AS build
ENV GO111MODULE=on

WORKDIR /app

# Download dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy app
COPY . .

# build the actual binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o entry

##
## Deploy stage
##
FROM alpine:3.18

# install common deps
RUN apk add curl wget bash

# copy the prebuilt file
WORKDIR /
COPY --from=build /app/entry /usr/bin/entry

# set app as startup app
ENTRYPOINT ["/usr/bin/entry"]
