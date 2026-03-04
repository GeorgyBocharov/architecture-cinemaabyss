FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /application cmd/*

FROM alpine:latest

COPY --from=builder /application /application

ENTRYPOINT ["/application"]
