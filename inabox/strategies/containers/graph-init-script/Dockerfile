# Custom to a specific subgraph version, so must be rebuilt with any subgraph changes
FROM golang:1.21.1-alpine3.18 as builder

WORKDIR /

COPY go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
# Mount the Go build cache from a volume into the container to allow caching
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY inabox/strategies/containers/graph-init-script/ /inabox/strategies/containers/graph-init-script/
COPY inabox/strategies/containers/config ./inabox/strategies/containers/config

WORKDIR /inabox/strategies/containers/graph-init-script

# Build the Go app with build cache mounted for faster builds
# The CGO_ENABLED=0 setting is for fully static compilation, adjust as necessary for your application
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o ./bin/graph-init ./cmd

FROM ghcr.io/layr-labs/foundry:latest

RUN apk add --no-cache bash nodejs npm curl
RUN npm install --global yarn

COPY subgraphs /subgraphs
COPY inabox/strategies/containers/graph-init-script/graph-init.sh .
COPY inabox/strategies/containers/graph-init-script/server.js .

RUN for d in $(ls -d /subgraphs/*/); do \
    cd $d && \
    yarn install && \
    yarn codegen && \
    yarn build && \
    cd /; \
    done

COPY --from=builder /inabox/strategies/containers/graph-init-script/bin/graph-init .

VOLUME [ "/data" ]

WORKDIR /data

CMD ["../graph-init.sh"]
