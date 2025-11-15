FROM golang:1.22-alpine AS builder

WORKDIR /src

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app-go .