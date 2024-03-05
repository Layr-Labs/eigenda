FROM golang:1.21.1-alpine3.18 as builder

RUN apk add --no-cache make musl-dev linux-headers gcc git jq bash

# build encoder with local monorepo go modules
COPY ./disperser /app/disperser
COPY common /app/common
COPY contracts /app/contracts
COPY core /app/core
COPY api /app/api
COPY indexer /app/indexer
COPY encoding /app/encoding
COPY go.mod /app
COPY go.sum /app

WORKDIR /app/disperser

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \ 
    go build -o ./bin/server ./cmd/encoder

FROM alpine:3.18

COPY --from=builder /app/disperser/bin/server /usr/local/bin

ENTRYPOINT ["server"]
