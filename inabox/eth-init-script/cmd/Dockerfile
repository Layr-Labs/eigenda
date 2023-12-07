# Start from the latest Golang base image
FROM golang:1.21.1-alpine3.18 as builder

RUN apk add --no-cache make musl-dev linux-headers gcc git jq bash

# Set the Current Working Directory inside the container
WORKDIR /app/eth-init-script

# Copy go mod and sum files
COPY go.mod go.sum ./
COPY ./api ./api

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
# Mount the Go build cache from a volume into the container to allow caching
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY ./inabox/config ./inabox/config
COPY ./inabox/eth-init-script ./inabox/eth-init-script
COPY inabox/utils ./inabox/utils

COPY ./disperser ./disperser
COPY ./common ./common
COPY ./core ./core
COPY ./pkg ./pkg
COPY ./contracts ./contracts

WORKDIR /app/eth-init-script/inabox/eth-init-script

# Build the Go app with build cache mounted for faster builds
# The CGO_ENABLED=0 setting is for fully static compilation, adjust as necessary for your application
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o ./bin/init ./cmd

######## Start a new stage from scratch #######
# Using scratch (empty container) for minimal size - alternatively use a small base image like alpine if you need shell access
FROM ghcr.io/foundry-rs/foundry:nightly-f5c91995f80b5bf3b4c29c934d414cc198c9e7a8

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/eth-init-script/inabox/eth-init-script/bin/init .

WORKDIR /data

# Command to run the executable
CMD ["../init"]
