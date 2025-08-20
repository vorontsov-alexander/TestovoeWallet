# 1) build stage
FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o wallet-service .

# 2) runtime
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/wallet-service /app/wallet-service
COPY config.env /app/config.env
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/wallet-service"]
