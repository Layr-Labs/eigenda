FROM golang:1.21.1-alpine3.18 as builder

RUN apk add --no-cache make musl-dev linux-headers gcc git jq bash

WORKDIR /app

# Copy Entire Repo here in order to not copy individual dependencies
COPY . .

WORKDIR /app/tools/traffic

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o ./bin/generator ./cmd2

FROM alpine:3.18 as generator2

COPY --from=builder /app/tools/traffic/bin/generator /usr/local/bin

ENTRYPOINT ["generator"]
