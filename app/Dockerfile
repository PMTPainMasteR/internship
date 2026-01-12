FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod init wifi_sim || true

RUN go build -o wifi_visit wifi_visit_v5.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/wifi_visit .

CMD ["./wifi_visit"]
