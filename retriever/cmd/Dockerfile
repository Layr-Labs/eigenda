FROM golang:1.21.1-alpine3.18 as builder

RUN apk add --no-cache make musl-dev linux-headers gcc git jq bash

# build node with local monorepo go modules
COPY ./retriever /app/retriever
COPY common /app/common
COPY core /app/core
COPY api /app/api
COPY contracts /app/contracts
COPY disperser /app/disperser
COPY node /app/node
COPY operators/churner /app/operators/churner
COPY indexer /app/indexer
COPY encoding /app/encoding
COPY go.mod /app
COPY go.sum /app

WORKDIR /app/retriever

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \ 
    go build -o ./bin/retriever ./cmd

FROM alpine:3.18

COPY --from=builder /app/retriever/bin/retriever /usr/local/bin

ENTRYPOINT ["retriever"]
