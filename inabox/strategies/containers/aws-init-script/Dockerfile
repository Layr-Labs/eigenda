# Start from the latest Golang base image
FROM golang:1.21.1-alpine3.18 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
# Mount the Go build cache from a volume into the container to allow caching
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY ./inabox/strategies/containers/config ./inabox/strategies/containers/config
COPY ./inabox/strategies/containers/aws-init-script ./inabox/strategies/containers/aws-init-script
COPY ./disperser ./disperser
COPY ./common ./common
COPY ./core ./core
COPY ./pkg ./pkg
COPY ./api ./api
COPY ./contracts ./contracts

WORKDIR /app/inabox/strategies/containers/aws-init-script

# Build the Go app with build cache mounted for faster builds
# The CGO_ENABLED=0 setting is for fully static compilation, adjust as necessary for your application
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o ./bin/aws-init-script ./cmd

######## Start a new stage from scratch #######
# Using scratch (empty container) for minimal size - alternatively use a small base image like alpine if you need shell access
FROM alpine:latest

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/inabox/strategies/containers/aws-init-script/bin/aws-init-script .

WORKDIR /data

# Command to run the executable
CMD ["../aws-init-script"]
