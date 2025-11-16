FROM golang:1.23.3 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/main ./cmd/app/main.go

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/migrate ./cmd/migrations/main.go

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/bin/main .
COPY --from=builder /app/bin/migrate .

COPY --from=builder /app/db/migrations ./db/migrations

EXPOSE 8080