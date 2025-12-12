# Multi-stage build: build Go binary, then small runtime image
FROM golang:1.25.5-alpine AS build

RUN apk add --no-cache git gcc musl-dev

WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /out/notora-server ./cmd/notora-server

FROM alpine:latest
RUN apk add --no-cache sqlite-libs tzdata

WORKDIR /app
COPY --from=build /out/notora-server /app/notora-server
COPY .env.example /app/.env.example

VOLUME ["/app/data"]

EXPOSE 8080

CMD ["/app/notora-server"]
