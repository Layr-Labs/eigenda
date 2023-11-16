# Start from the latest Golang base image
FROM golang:1.21.1-alpine3.18 as builder

RUN apk add --no-cache make musl-dev linux-headers gcc git jq bash

# Set the Current Working Directory inside the container
WORKDIR /app/deploy

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
# Mount the Go build cache from a volume into the container to allow caching
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY ./inabox-docker/deploy ./inabox-docker/deploy
COPY ./inabox-docker/local-init ./inabox-docker/local-init
COPY ./disperser ./disperser
COPY ./common ./common
COPY ./core ./core
COPY ./api ./api
COPY ./pkg ./pkg
COPY ./contracts ./contracts

WORKDIR /app/deploy/inabox-docker/local-init

# Build the Go app with build cache mounted for faster builds
# The CGO_ENABLED=0 setting is for fully static compilation, adjust as necessary for your application
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/init ./cmd

######## Start a new stage from scratch #######
# Using scratch (empty container) for minimal size - alternatively use a small base image like alpine if you need shell access
FROM ghcr.io/foundry-rs/foundry:latest
# FROM foundry:latest

RUN apk add --no-cache bash nodejs npm file
RUN npm install --global yarn

# this is really slow in the container for some reason so we do it here.
COPY ./subgraphs /subgraphs
RUN for d in $(ls -d /subgraphs/*/); do \
    cd $d && \
    yarn install && \
    # yarn codegen && \
    cd /; \
    done

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/deploy/inabox-docker/local-init/bin/init .

# Command to run the executable
CMD ["./init"]