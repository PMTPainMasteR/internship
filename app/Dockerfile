FROM golang:1.2 AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY wifi_visit_v5.go .

RUN CGO_ENABLED=0 go build -o wifi_visit wifi_visit_v5.go

FROM UBUNTU:24.04
WORKDIR /root/
COPY --from=builder /app/wifi_visit .
CMD ["./wifi_visit"]
