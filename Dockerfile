FROM golang:1.22.0-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o myapp .

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/myapp .

CMD ["./myapp"]
