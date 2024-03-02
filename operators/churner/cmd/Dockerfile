FROM golang:1.21.1-alpine3.18 as builder

RUN apk add --no-cache make musl-dev linux-headers gcc git jq bash

# build node with local monorepo go modules
COPY ./operators/churner /app/operators/churner
COPY common /app/common
COPY core /app/core
COPY api /app/api
COPY indexer /app/indexer
COPY contracts /app/contracts
COPY encoding /app/encoding
COPY go.mod /app
COPY go.sum /app

WORKDIR /app/operators/churner

RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go build -o ./bin/churner ./cmd

FROM alpine:3.18

COPY --from=builder /app/operators/churner/bin/churner /usr/local/bin

ENTRYPOINT ["churner"]
