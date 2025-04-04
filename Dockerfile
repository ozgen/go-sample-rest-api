# Step 1: Build Stage
FROM golang:1.23.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /api ./cmd/main.go
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /migrate ./cmd/migrate/main.go

COPY cmd/migrate/migrations /app/migrations

# Step 2: Runtime Stage
FROM alpine:3.14

WORKDIR /

# ✅ Install required runtime dependencies
RUN apk --no-cache add ca-certificates

# Copy binaries from build stage
COPY --from=builder /api /api
COPY --from=builder /migrate /migrate
RUN mkdir -p /cmd/migrate/migrations
COPY --from=builder /app/cmd/migrate/migrations /cmd/migrate/migrations

# Copy entrypoint script and make executable
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]
