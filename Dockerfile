# syntax=docker/dockerfile:1

# Declare build arguments
# NOTE: to use these args, they must be *consumed* in the child scope (see node-builder)
# https://docs.docker.com/build/building/variables/#scoping
ARG SEMVER=""
ARG GITCOMMIT=""
ARG GITDATE=""

FROM golang:1.21.1-alpine3.18 AS base-builder
RUN apk add --no-cache make musl-dev linux-headers gcc git jq bash

# Common build stage
FROM base-builder AS common-builder
WORKDIR /app
COPY go.mod go.sum ./
COPY disperser /app/disperser
COPY common /app/common
COPY core /app/core
COPY api /app/api
COPY contracts /app/contracts
COPY indexer /app/indexer
COPY encoding /app/encoding
COPY relay /app/relay

# Churner build stage
FROM common-builder AS churner-builder
COPY operators ./operators
WORKDIR /app/operators
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o ./bin/churner ./churner/cmd

# Encoder build stage
FROM common-builder AS encoder-builder
WORKDIR /app/disperser
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o ./bin/encoder ./cmd/encoder

# API Server build stage
FROM common-builder AS apiserver-builder
WORKDIR /app/disperser
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o ./bin/apiserver ./cmd/apiserver

# DataAPI build stage
FROM common-builder AS dataapi-builder
COPY operators ./operators
WORKDIR /app/disperser
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o ./bin/dataapi ./cmd/dataapi

# Batcher build stage
FROM common-builder AS batcher-builder
WORKDIR /app/disperser
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o ./bin/batcher ./cmd/batcher

# Retriever build stage
FROM common-builder AS retriever-builder
COPY retriever /app/retriever
COPY node /app/node
COPY operators ./operators
WORKDIR /app/retriever
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o ./bin/retriever ./cmd

# Node build stage
FROM common-builder AS node-builder
ARG SEMVER
ARG GITCOMMIT
ARG GITDATE
COPY node /app/node
COPY operators ./operators
WORKDIR /app/node
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags="-X 'github.com/Layr-Labs/eigenda/node.SemVer=${SEMVER}' -X 'github.com/Layr-Labs/eigenda/node.GitCommit=${GITCOMMIT}' -X 'github.com/Layr-Labs/eigenda/node.GitDate=${GITDATE}'" -o ./bin/node ./cmd

# Nodeplugin build stage
FROM common-builder AS node-plugin-builder
COPY ./node /app/node
COPY operators ./operators
WORKDIR /app/node
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o ./bin/nodeplugin ./plugin/cmd


# Final stages for each component
FROM alpine:3.18 AS churner
COPY --from=churner-builder /app/operators/bin/churner /usr/local/bin
ENTRYPOINT ["churner"]

FROM alpine:3.18 AS encoder
COPY --from=encoder-builder /app/disperser/bin/encoder /usr/local/bin
ENTRYPOINT ["encoder"]

FROM alpine:3.18 AS apiserver
COPY --from=apiserver-builder /app/disperser/bin/apiserver /usr/local/bin
ENTRYPOINT ["apiserver"]

FROM alpine:3.18 AS dataapi
COPY --from=dataapi-builder /app/disperser/bin/dataapi /usr/local/bin
ENTRYPOINT ["dataapi"]

FROM alpine:3.18 AS batcher
COPY --from=batcher-builder /app/disperser/bin/batcher /usr/local/bin
ENTRYPOINT ["batcher"]

FROM alpine:3.18 AS retriever
COPY --from=retriever-builder /app/retriever/bin/retriever /usr/local/bin
ENTRYPOINT ["retriever"]

FROM alpine:3.18 AS node
COPY --from=node-builder /app/node/bin/node /usr/local/bin
ENTRYPOINT ["node"]

FROM alpine:3.18 AS nodeplugin
COPY --from=node-plugin-builder /app/node/bin/nodeplugin /usr/local/bin
ENTRYPOINT ["nodeplugin"]
