FROM golang:1.21.13-alpine3.20 AS builder

RUN apk add --no-cache make musl-dev linux-headers gcc git jq bash

WORKDIR /app

RUN apk add --no-cache make

# Copy Entire Repo here in order to not copy individual dependencies
COPY . .

WORKDIR app/test/v2

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o test/v2/bin/load test/v2/load/main/load_main.go

FROM alpine:3.18 AS generator2

COPY --from=builder /app/test/v2/bin/load /usr/local/bin/load

ENTRYPOINT ["generator"]
