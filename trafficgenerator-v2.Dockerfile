FROM golang:1.21.13-alpine3.20 as builder

RUN apk add --no-cache make musl-dev linux-headers gcc git jq bash

WORKDIR /app

RUN apk add --no-cache make

# Copy Entire Repo here in order to not copy individual dependencies
COPY . .

WORKDIR app/test/v2

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    make build

FROM alpine:3.18 AS generator2

COPY --from=builder /app/test/v2/bin/load /usr/local/bin/load

ENTRYPOINT ["generator"]
