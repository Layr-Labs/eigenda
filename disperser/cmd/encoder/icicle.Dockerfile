# This Dockerfile has been tested on Ubuntu 24.04
# Note: Will fail on macOS with "gcc: error: unrecognized command-line option '-m64'" during cgo compilation, which is expected because cuda is not available.
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
COPY . .

# Define Icicle versions and checksums
# If you ever change the ICICLE_VERSION, first find the new artifact links from
# https://github.com/ingonyama-zk/open-icicle/releases, and then compute the new checksums by running:
#  wget https://github.com/ingonyama-zk/open-icicle/releases/download/v3.9.2/open-icicle_3_9_2-ubuntu22.tar.gz
#  sha256sum open-icicle_3_9_2-ubuntu22.tar.gz
#  wget https://github.com/ingonyama-zk/open-icicle/releases/download/v3.9.2/open-icicle_3_9_2-ubuntu22-cuda122.tar.gz
#  sha256sum open-icicle_3_9_2-ubuntu22-cuda122.tar.gz
ENV ICICLE_VERSION=3_9_2
ENV ICICLE_BASE_SHA256=d3576c2abff38a36a6129642a5617cf97a6bebf0a18f90e76d550c5c1660b50d
ENV ICICLE_CUDA_SHA256=1f9c194bde786d178b5ac10f151d2864885a40f411ebe508bbe2a29fddfb030f

# Download Icicle tarballs from the open-icicle release
ADD https://github.com/ingonyama-zk/open-icicle/releases/download/v3.9.2/open-icicle_${ICICLE_VERSION}-ubuntu22.tar.gz /tmp/icicle.tar.gz
ADD https://github.com/ingonyama-zk/open-icicle/releases/download/v3.9.2/open-icicle_${ICICLE_VERSION}-ubuntu22-cuda122.tar.gz /tmp/icicle-cuda.tar.gz

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
