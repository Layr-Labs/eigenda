FROM golang:1.21.1-alpine3.18 as builder

RUN apk add --no-cache make musl-dev linux-headers gcc git jq bash

# build batcher with local monorepo go modules
COPY ./disperser /app/disperser
COPY common /app/common
COPY core /app/core
COPY api /app/api
COPY contracts /app/contracts
COPY indexer /app/indexer
COPY pkg /app/pkg
COPY go.mod /app
COPY go.sum /app

WORKDIR /app/disperser

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \ 
    go build -o ./bin/server ./cmd/batcher

FROM alpine:3.18

COPY --from=builder /app/disperser/bin/server /usr/local/bin

ENTRYPOINT ["server"]
