# Use the latest foundry image
FROM --platform=linux/amd64 ghcr.io/foundry-rs/foundry:latest

# Copy our source code into the container
WORKDIR /app

# Build and test the source code
COPY . .
RUN forge build
RUN forge test

# Set the entrypoint to the forge command
ENTRYPOINT ["forge"]