FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o qa-app ./cmd/question_answer

# --- Runtime image ---
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/qa-app /app/qa-app
COPY --from=builder /app/config /app/config
COPY --from=builder /app/internal/infrastructure/storage/postgres/migrations \
    /app/internal/infrastructure/storage/postgres/migrations

EXPOSE 8082
CMD ["./qa-app"]
