FROM golang:1.26.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/migrate ./cmd/migrate

FROM alpine:3.22

RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /bin/server /app/server
COPY --from=builder /bin/migrate /app/migrate
COPY migration /app/migration

# ENV APP_PORT=3000

EXPOSE 3000

USER appuser

CMD ["/app/server"]
