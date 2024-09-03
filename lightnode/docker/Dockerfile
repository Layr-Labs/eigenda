FROM golang:1.21.12-bookworm AS base

# The url of the git repository to clone.
ARG GIT_URL
# The branch or commit to check out.
ARG BRANCH_OR_COMMIT

# Install core dependencies
RUN apt update
RUN apt install -y git build-essential bash

# Set up lnode user
RUN useradd -m -s /bin/bash lnode
USER lnode
WORKDIR /home/lnode
# Remove default crud
RUN rm .bashrc
RUN rm .bash_logout
RUN rm .profile

# Clone eigenda repository
RUN git clone $GIT_URL eigenda
WORKDIR /home/lnode/eigenda
RUN git checkout $BRANCH_OR_COMMIT
WORKDIR /home/lnode

# Download all go dependencies.
# This is done as a separate step to avoid repeat work every time the latest commit in the branch is updated.
WORKDIR /home/lnode/eigenda/lightnode
RUN go mod download
WORKDIR /home/lnode

# Copy a file containing build information. Useful for detective work on an otherwise unlabelled image.
# This is also useful for forcing docker to invalidate caches when a target branch is updated.
WORKDIR /home/lnode/eigenda/lightnode
COPY --chown=lnode build-info.txt /home/lnode
WORKDIR /home/lnode

# Just in case we are tracking the latest commit in a branch, pull again. This is a no-op if the target is a commit sha.
WORKDIR /home/lnode/eigenda
RUN git pull
WORKDIR /home/lnode

# Build the light node binary.
WORKDIR /home/lnode/eigenda/lightnode
RUN make build
WORKDIR /home/lnode

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

# Copy the build info file.
COPY --from=base /home/lnode/build-info.txt /home/lnode/build-info.txt

# Run the light node when the container starts.
CMD ["/home/lnode/lnode"]