# This Dockerfile has been tested on Ubuntu 24.04
# Note: Will fail on macOS with "gcc: error: unrecognized command-line option '-m64'" during cgo compilation
FROM nvidia/cuda:12.2.2-devel-ubuntu22.04 AS builder

# Install Go 1.24.4 to match go.mod requirements
ENV GOLANG_VERSION=1.24.4
ENV GOLANG_SHA256=77e5da33bb72aeaef1ba4418b6fe511bc4d041873cbf82e5aa6318740df98717

ADD https://go.dev/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz /tmp/go.tar.gz
RUN echo "${GOLANG_SHA256} /tmp/go.tar.gz" | sha256sum -c - && \
    tar -C /usr/local -xzf /tmp/go.tar.gz && \
    rm /tmp/go.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"

# Set up the working directory
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./

# Copy api/proxy/clients for the replace directive
COPY api/proxy/clients ./api/proxy/clients

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY ./disperser /app/disperser
COPY common /app/common
COPY contracts /app/contracts
COPY core /app/core
COPY api /app/api
COPY indexer /app/indexer
COPY encoding /app/encoding
COPY relay /app/relay

# Define Icicle versions and checksums
ENV ICICLE_VERSION=3.4
ENV ICICLE_BASE_SHA256=e3eec1d5fca0e4ba52e09630dc360eb5f1c1d54b3bb1834eeef1624c2f7f4c48
ENV ICICLE_CUDA_SHA256=090bdf255b1beab05efa7eb2e67d50481c0c6f57beda969356e71bd11bb39f17

# Download Icicle tarballs
# https://github.com/ingonyama-zk/icicle/releases/download/v3.4.0/icicle_3_4-ubuntu22-cuda122.tar.gz
ADD https://github.com/ingonyama-zk/icicle/releases/download/v${ICICLE_VERSION}.0/icicle_${ICICLE_VERSION//./_}-ubuntu22.tar.gz /tmp/icicle.tar.gz
ADD https://github.com/ingonyama-zk/icicle/releases/download/v${ICICLE_VERSION}.0/icicle_${ICICLE_VERSION//./_}-ubuntu22-cuda122.tar.gz /tmp/icicle-cuda.tar.gz

# Verify checksums and install Icicle
RUN echo "${ICICLE_BASE_SHA256} /tmp/icicle.tar.gz" | sha256sum -c - && \
    echo "${ICICLE_CUDA_SHA256} /tmp/icicle-cuda.tar.gz" | sha256sum -c - && \
    tar xzf /tmp/icicle.tar.gz && \
    cp -r ./icicle/lib/* /usr/lib/ && \
    cp -r ./icicle/include/icicle/ /usr/local/include/ && \
    tar xzf /tmp/icicle-cuda.tar.gz -C /opt && \
    rm /tmp/icicle.tar.gz /tmp/icicle-cuda.tar.gz

# Build the server with icicle backend
WORKDIR /app/disperser
RUN go build -tags=icicle -o ./bin/server ./cmd/encoder

# Start a new stage for the base image
FROM nvidia/cuda:12.2.2-base-ubuntu22.04

COPY --from=builder /app/disperser/bin/server /usr/local/bin/server
COPY --from=builder /usr/lib/libicicle* /usr/lib/
COPY --from=builder /usr/local/include/icicle /usr/local/include/icicle
COPY --from=builder /opt/icicle /opt/icicle

ENTRYPOINT ["server"]
