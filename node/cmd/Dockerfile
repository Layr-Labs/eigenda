FROM golang:1.21.1-alpine3.18 as builder

ARG SEMVER="v0.x.y.z"
ARG GITCOMMIT="da48cbe1fd7ae96fe1f8fb2cce4f939459d97565"
ARG GITDATE="2024-04-24T00:00:00Z"

RUN apk add --no-cache make musl-dev linux-headers gcc git jq bash

# build node with local monorepo go modules
COPY ./node /app/node
COPY operators/churner /app/operators/churner
COPY indexer /app/indexer
COPY common /app/common
COPY core /app/core
COPY api /app/api
COPY contracts /app/contracts
COPY encoding /app/encoding
COPY go.mod /app
COPY go.sum /app

WORKDIR /app/node

RUN go build -ldflags="-X 'github.com/Layr-Labs/eigenda/node.SemVer=${SEMVER}' -X 'github.com/Layr-Labs/eigenda/node.GitCommit=${GITCOMMIT}' -X 'github.com/Layr-Labs/eigenda/node.GitDate=${GITDATE}'" -o ./bin/node ./cmd

FROM alpine:3.18

COPY --from=builder /app/node/bin/node /usr/local/bin

ENTRYPOINT ["node"]
