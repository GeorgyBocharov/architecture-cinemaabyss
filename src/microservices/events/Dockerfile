FROM golang:1.25-alpine AS builder

RUN apk add --no-cache gcc musl-dev git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -tags musl -o /application cmd/*

FROM alpine:latest

COPY --from=builder /application /application

ENTRYPOINT ["/application"]
