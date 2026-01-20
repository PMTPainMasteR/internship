FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY wifi_visit_v5.go .

RUN CGO_ENABLED=0 go build -o wifi_visit wifi_visit_v5.go

FROM gcr.io/distroless/base-debian12
WORKDIR /
# Copy binary from builder
COPY --from=builder /app/wifi_visit /
# Command to run
CMD ["/wifi_visit"]
