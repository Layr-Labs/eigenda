FROM nvidia/cuda:12.0.0-devel-ubuntu22.04 as builder

# Install necessary build tools and Go
RUN apt-get update && apt-get install -y \
    make \
    gcc \
    git \
    jq \
    bash \
    curl \
    cmake \
    && rm -rf /var/lib/apt/lists/*

# Install Go
ENV GOLANG_VERSION=1.21.1
RUN curl -L https://go.dev/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz | tar -xz -C /usr/local
ENV PATH="/usr/local/go/bin:${PATH}"

# Set up the working directory
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY ./disperser /app/disperser
COPY common /app/common
COPY contracts /app/contracts
COPY core /app/core
COPY api /app/api
COPY indexer /app/indexer
COPY encoding /app/encoding

# Build Icicle libraries
# WORKDIR /app/icicle/wrappers/golang
# RUN ./build.sh -curve=bn254 -ecntt

# Set up CUDA paths
ENV CUDA_PATH=/usr/local/cuda
ENV LD_LIBRARY_PATH=${CUDA_PATH}/lib64:${LD_LIBRARY_PATH}
ENV CPATH=${CUDA_PATH}/include:${CPATH}

# Set up CGO flags for shared libraries
ENV CGO_LDFLAGS="-L/app/encoding/lib"

# Build the server with GPU support
WORKDIR /app/disperser
RUN go build -tags gpu -o ./bin/server ./cmd/encoder

# Start a new stage for the base image
FROM nvidia/cuda:12.0.0-base-ubuntu22.04

COPY --from=builder /app/disperser/bin/server /usr/local/bin/server

ENTRYPOINT ["server"]