FROM golang:1.22.0 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/main.go

FROM alpine:3.14

WORKDIR /

COPY --from=builder /api /api

EXPOSE 8080

CMD ["/api"]
