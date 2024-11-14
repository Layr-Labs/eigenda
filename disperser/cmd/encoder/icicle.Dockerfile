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
COPY icicle /app/icicle
COPY relay /app/relay

# Install Icicle
RUN cp -r /app/icicle/lib/* /usr/lib/ && \
    cp -r /app/icicle/include/icicle/ /usr/local/include/ && \
    cp -r /app/icicle /opt

# Build the server with GPU support
WORKDIR /app/disperser
RUN go build -tags=icicle -o ./bin/server ./cmd/encoder

# Start a new stage for the base image
FROM nvidia/cuda:12.2.2-base-ubuntu22.04

COPY --from=builder /app/disperser/bin/server /usr/local/bin/server
COPY --from=builder /usr/lib/libicicle* /usr/lib/
COPY --from=builder /usr/local/include/icicle /usr/local/include/icicle
COPY --from=builder /opt/icicle /opt/icicle

ENTRYPOINT ["server"]