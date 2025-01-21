FROM golang:1.21.13-alpine3.20 as builder

RUN apk add --no-cache make musl-dev linux-headers gcc git jq bash

WORKDIR /app

# Copy Entire Repo here in order to not copy individual dependencies
COPY . .

WORKDIR /app/tools/traffic

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \ 
    go build -o ./bin/generator ./cmd

FROM alpine:3.18 AS generator

COPY --from=builder /app/tools/traffic/bin/generator /usr/local/bin

ENTRYPOINT ["generator"]
