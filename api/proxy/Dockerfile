# multi container builds ftw

FROM golang:1.24.4-alpine3.22 AS builder

RUN apk add --no-cache make gcc musl-dev linux-headers jq bash git

# Set the working directory inside the container
WORKDIR /app

# Copy the entire repository (build context should be repo root)
COPY . .

# Build the application binary from api/proxy directory
WORKDIR /app/api/proxy
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    make build

# Use alpine to run app
FROM alpine:3.22

WORKDIR /app
COPY --from=builder /app/api/proxy/bin/eigenda-proxy .

# Copy srs values
COPY --from=builder /app/api/proxy/resources/ /app/resources/

# API & metrics servers
EXPOSE 4242 7300

# Run app
ENTRYPOINT ["./eigenda-proxy"]
