# multi container builds ftw

FROM golang:1.21.10-alpine3.19 as builder

RUN apk add --no-cache make gcc musl-dev linux-headers jq bash git

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application binary
RUN make eigenda-proxy

# Use alpine to run app
FROM alpine:3.16

WORKDIR /app
COPY --from=builder /app/bin/eigenda-proxy .

# Copy srs values
COPY --from=builder /app/resources/ /resources/

# API & metrics servers
EXPOSE 4242 7300

# Run app
CMD ["./eigenda-proxy"]