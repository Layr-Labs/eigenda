FROM golang:1.21.1-alpine3.18 as builder

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

RUN go build -o ./bin/nodeplugin ./plugin/cmd

FROM alpine:3.18

COPY --from=builder /app/node/bin/nodeplugin /usr/local/bin

ENTRYPOINT ["nodeplugin"]
