# service/Dockerfile => require to build from the root
FROM docker.arvancloud.ir/golang:1.24-alpine AS builder

ARG SERVICE_NAME=account-service
ENV SERVICE_NAME=${SERVICE_NAME}

WORKDIR /app

COPY shared ./shared
COPY ${SERVICE_NAME} ./${SERVICE_NAME}

# Download dependencies for shared module
WORKDIR /app/shared
RUN go mod download

# Download and build service
WORKDIR /app/${SERVICE_NAME}
RUN go mod download
RUN GOOS=linux CGO_ENABLED=0 go build -o /bin/${SERVICE_NAME} ./cmd/

# --- Final Image ---
FROM docker.arvancloud.ir/alpine:latest

ARG SERVICE_NAME=account-service
ENV SERVICE_NAME=${SERVICE_NAME}

WORKDIR /app

COPY --from=builder /bin/${SERVICE_NAME} .

CMD sh -c  "/app/$SERVICE_NAME"

