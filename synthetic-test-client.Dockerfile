# Use the official Go image as the base image
FROM golang:1.21.1-alpine3.18 as builder

# Copy only the test file and necessary files to the container
COPY ./disperser /app/disperser
COPY ./test/synthetic-test /app
COPY go.mod /app
COPY go.sum /app
COPY api /app/api
COPY clients /app/clients
COPY node /app/node
COPY common /app/common
COPY churner /app/churner
COPY core /app/core
COPY indexer /app/indexer
COPY contracts /app/contracts
COPY pkg /app/pkg
# Set the working directory inside the container
WORKDIR /app

# TODO eventually this will be replaced with an executable
# Run the Go test command for the specific test file
CMD ["go", "test", "-v", "synthetic_client_test.go"]
