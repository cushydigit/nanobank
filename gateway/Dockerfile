# gateway/Dockerfile => require to build from the root
FROM docker.arvancloud.ir/golang:1.24-alpine AS builder

WORKDIR /app

COPY shared ./shared
COPY gateway ./gateway

# Download dependencies for shared module
WORKDIR /app/shared
RUN go mod download

# Download and build service
WORKDIR /app/gateway
RUN go mod download
RUN GOOS=linux CGO_ENABLED=0 go build -o /bin/gateway .

# --- Final Image ---
FROM docker.arvancloud.ir/alpine:latest

WORKDIR /app

COPY --from=builder /bin/gateway .

CMD ["/app/gateway"]

