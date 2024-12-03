FROM nvidia/cuda:12.2.2-devel-ubuntu22.04 AS builder

# Install Go
ENV GOLANG_VERSION=1.21.1
ENV GOLANG_SHA256=b3075ae1ce5dab85f89bc7905d1632de23ca196bd8336afd93fa97434cfa55ae

ADD https://go.dev/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz /tmp/go.tar.gz
RUN echo "${GOLANG_SHA256} /tmp/go.tar.gz" | sha256sum -c - && \
    tar -C /usr/local -xzf /tmp/go.tar.gz && \
    rm /tmp/go.tar.gz
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
COPY relay /app/relay

# Define Icicle versions and checksums
ENV ICICLE_VERSION=3.1.0
ENV ICICLE_BASE_SHA256=2e4e33b8bc3e335b2dd33dcfb10a9aaa18717885509614a24f492f47a2e4f4b1
ENV ICICLE_CUDA_SHA256=cdba907eac6297445a6c128081ebba5c711d352003f69310145406a8fd781647

# Download Icicle tarballs
ADD https://github.com/ingonyama-zk/icicle/releases/download/v${ICICLE_VERSION}/icicle_${ICICLE_VERSION//./_}-ubuntu22.tar.gz /tmp/icicle.tar.gz
ADD https://github.com/ingonyama-zk/icicle/releases/download/v${ICICLE_VERSION}/icicle_${ICICLE_VERSION//./_}-ubuntu22-cuda122.tar.gz /tmp/icicle-cuda.tar.gz

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
