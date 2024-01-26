FROM golang:1.21-alpine as build-server

RUN apk add build-base

# Initialization
RUN mkdir -p /app
WORKDIR /app

# Dependencies
COPY go.mod ./
COPY go.sum .

RUN go mod download

COPY cmd ./cmd
COPY pkg ./pkg

RUN GOOS=linux GOARCH=amd64 go build -o valetudo-telegram-bot ./cmd/valetudo-telegram-bot/main.go

FROM debian:bullseye-slim

# Options
ENV TELEGRAM_BOT_TOKEN ""
ENV TELEGRAM_CHAT_IDS ""
ENV VALETUDO_URL ""
ENV TELEGRAM_DEBUG false

# Copy build results
WORKDIR /app

COPY --from=build-server /app/valetudo-telegram-bot ./valetudo-telegram-bot

# Start the application
ENTRYPOINT [ "./valetudo-telegram-bot" ]
