# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies for sqlite CGO
RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the binary
RUN CGO_ENABLED=1 go build -o bin/api ./cmd/api

# Runtime stage
FROM alpine:3.21

WORKDIR /app

# Runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/bin/api .
COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./api"]
