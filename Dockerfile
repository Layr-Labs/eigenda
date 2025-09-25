# syntax=docker/dockerfile:1

# Declare build arguments
# NOTE: to use these args, they must be *consumed* in the child scope (see node-builder)
# https://docs.docker.com/build/building/variables/#scoping
ARG SEMVER=""
ARG GITCOMMIT=""
ARG GITDATE=""

FROM golang:1.24.4-alpine3.22 AS base-builder
RUN apk add --no-cache make musl-dev linux-headers gcc git jq bash

# Common build stage
FROM base-builder AS common-builder
WORKDIR /app
COPY . .

# Churner build stage
FROM common-builder AS churner-builder
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
WORKDIR /app/retriever
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o ./bin/retriever ./cmd

# Node build stage
FROM common-builder AS node-builder
ARG SEMVER
ARG GITCOMMIT
ARG GITDATE
WORKDIR /app/node
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags="-X 'github.com/Layr-Labs/eigenda/node.SemVer=${SEMVER}' -X 'github.com/Layr-Labs/eigenda/node.GitCommit=${GITCOMMIT}' -X 'github.com/Layr-Labs/eigenda/node.GitDate=${GITDATE}'" -o ./bin/node ./cmd

# Nodeplugin build stage
FROM common-builder AS node-plugin-builder
WORKDIR /app/node
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o ./bin/nodeplugin ./plugin/cmd

# Controller build stage
FROM common-builder AS controller-builder
COPY node/auth /app/node/auth
WORKDIR /app/disperser
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o ./bin/controller ./cmd/controller

# Relay build stage
FROM common-builder AS relay-builder
WORKDIR /app/relay
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o ./bin/relay ./cmd

# Traffic Generator V1 build stage
FROM common-builder AS generator-builder
WORKDIR /app/tools/traffic
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o ./bin/generator ./cmd

# Traffic Generator V2 build stage
FROM common-builder AS generator2-builder
WORKDIR /app/test/v2
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    make build

# BlobAPI (Combined API Server and Relay) build stage
FROM common-builder AS blobapi-builder
WORKDIR /app/disperser
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags="-X main.version=${SEMVER} \
    -X main.gitCommit=${GITCOMMIT} \
    -X main.gitDate=${GITDATE}" \
    -o ./bin/blobapi ./cmd/blobapi

# Proxy build stage
FROM common-builder AS proxy-builder
WORKDIR /app/api/proxy
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags="-X main.version=${SEMVER} \
    -X main.gitCommit=${GITCOMMIT} \
    -X main.gitDate=${GITDATE}" \
    -o ./bin/eigenda-proxy ./cmd/server

# Final stages for each component
FROM alpine:3.22 AS churner
COPY --from=churner-builder /app/operators/bin/churner /usr/local/bin
ENTRYPOINT ["churner"]

FROM alpine:3.22 AS encoder
COPY --from=encoder-builder /app/disperser/bin/encoder /usr/local/bin
ENTRYPOINT ["encoder"]

FROM alpine:3.22 AS apiserver
COPY --from=apiserver-builder /app/disperser/bin/apiserver /usr/local/bin
ENTRYPOINT ["apiserver"]

FROM alpine:3.22 AS dataapi
COPY --from=dataapi-builder /app/disperser/bin/dataapi /usr/local/bin
ENTRYPOINT ["dataapi"]

FROM alpine:3.22 AS batcher
COPY --from=batcher-builder /app/disperser/bin/batcher /usr/local/bin
ENTRYPOINT ["batcher"]

FROM alpine:3.22 AS retriever
COPY --from=retriever-builder /app/retriever/bin/retriever /usr/local/bin
ENTRYPOINT ["retriever"]

FROM alpine:3.22 AS node
COPY --from=node-builder /app/node/bin/node /usr/local/bin
ENTRYPOINT ["node"]

FROM alpine:3.22 AS node-goreleaser
COPY node /usr/local/bin
ENTRYPOINT ["node"]

FROM alpine:3.22 AS nodeplugin
COPY --from=node-plugin-builder /app/node/bin/nodeplugin /usr/local/bin
ENTRYPOINT ["nodeplugin"]

FROM alpine:3.22 AS nodeplugin-goreleaser
COPY nodeplugin /usr/local/bin
ENTRYPOINT ["nodeplugin"]

FROM alpine:3.22 AS controller
COPY --from=controller-builder /app/disperser/bin/controller /usr/local/bin
ENTRYPOINT ["controller"]

FROM alpine:3.22 AS relay
COPY --from=relay-builder /app/relay/bin/relay /usr/local/bin
ENTRYPOINT ["relay"]

FROM alpine:3.22 AS generator
COPY --from=generator-builder /app/tools/traffic/bin/generator /usr/local/bin
ENTRYPOINT ["generator"]

FROM alpine:3.22 AS generator2
COPY --from=generator2-builder /app/test/v2/bin/load /usr/local/bin
ENTRYPOINT ["load", "-", "-"]

FROM alpine:3.22 AS blobapi
COPY --from=blobapi-builder /app/disperser/bin/blobapi /usr/local/bin
ENTRYPOINT ["blobapi"]

# proxy doesn't follow the same pattern as the others, because we keep it in the same
# format as when it was a separate repo: https://github.com/Layr-Labs/eigenda-proxy/blob/main/Dockerfile
FROM alpine:3.22 AS proxy
WORKDIR /app
COPY --from=proxy-builder /app/api/proxy/bin/eigenda-proxy .
COPY --from=proxy-builder /app/api/proxy/resources/ /app/resources/
# default ports for data and metrics
EXPOSE 3100 7300
ENTRYPOINT ["./eigenda-proxy"]

FROM alpine:3.22 AS proxy-goreleaser
WORKDIR /app
COPY eigenda-proxy .
COPY api/proxy/resources/*.point resources/
# default ports for data and metrics
EXPOSE 3100 7300
ENTRYPOINT ["./eigenda-proxy"]
