FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY src/go.mod src/go.sum ./

RUN go mod download
RUN apk add --no-cache gcc musl-dev
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0

COPY src .

RUN sqlc generate
RUN CGO_ENABLED=1 go build -ldflags="-w -s" -o main .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

CMD ["./main", "--ip", "0.0.0.0", "--port", "80"]