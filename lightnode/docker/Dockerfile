FROM golang:1.21.12-bookworm AS base

# Install core dependencies
RUN apt update
RUN apt install -y build-essential

# Set up lnode user
RUN useradd -m -s /bin/bash lnode
USER lnode
WORKDIR /home/lnode
# Remove default crud
RUN rm .bashrc
RUN rm .bash_logout
RUN rm .profile

# Copy the eigenda repository
COPY --chown=lnode . /home/lnode/eigenda

# Download all go dependencies and build the binary.
WORKDIR /home/lnode/eigenda/lightnode
RUN go mod download
RUN make build

# In order to equip this image with a shell for debugging,
# swap out the "FROM scratch" below with "FROM golang:1.21.12-bookworm"
FROM scratch AS final

# Copy over files needed for lnode user.
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group

USER lnode
WORKDIR /home/lnode

# Copy the executable binary.
COPY --from=base /home/lnode/eigenda/lightnode/bin/lnode /home/lnode/lnode

# Run the light node when the container starts.
CMD ["/home/lnode/lnode"]