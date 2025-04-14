# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy module files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Copy .env file
COPY .env ./

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -o 2fa-app .

# Runtime stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/2fa-app .
COPY --from=builder /app/.env ./

EXPOSE 8000
CMD ["./2fa-app"]