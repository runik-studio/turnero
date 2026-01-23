# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
# Copy .env if it exists (usually better to pass env vars in docker-compose, but useful for standalone)
COPY --from=builder /app/.env .

COPY --from=builder /app/firebaseCredentials.json .


EXPOSE 8080
CMD ["./main"]
