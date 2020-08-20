FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git
WORKDIR /root/project
COPY . .
