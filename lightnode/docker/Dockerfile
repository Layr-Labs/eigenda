# syntax=docker/dockerfile:1

# lnode-git builds on top of lnode-base by cloning the eigenda repository. A useful caching step to avoid
# re-downloading the entire repository every time the image is built.

FROM debian

# The url of the git repository to clone.
ARG GIT_URL
# The branch or commit to check out.
ARG BRANCH_OR_COMMIT

# Install core dependencies
RUN apt update
RUN apt install -y wget git build-essential bash

# Set up lnode user
RUN useradd -m -s /bin/bash lnode
USER lnode
WORKDIR /home/lnode
# Remove default crud
RUN rm .bashrc
RUN rm .bash_logout
RUN rm .profile

# Install golang
RUN wget https://go.dev/dl/go1.21.12.linux-arm64.tar.gz # TODO make this an argument
RUN tar -xf ./*.tar.gz
RUN rm ./*.tar.gz
RUN mkdir -p ~/.local/share
RUN mv go ~/.local/share/go
RUN rm .wget-hsts

# Setup golang environment for lnode user
RUN touch ~/.bashrc
RUN echo 'export GOPATH=~/.local/share/go' >> ~/.bashrc
RUN echo 'export PATH=~/.local/share/go/bin:$PATH' >> ~/.bashrc

# Clone eigenda repository
RUN git clone $GIT_URL eigenda
WORKDIR /home/lnode/eigenda
RUN git checkout $BRANCH_OR_COMMIT
WORKDIR /home/lnode

# Download all go dependencies.
# This is done as a separate step to avoid repeat work every time the latest commit in the branch is updated.
WORKDIR /home/lnode/eigenda/lightnode
RUN bash -c 'source ~/.bashrc && go mod download'
WORKDIR /home/lnode

# Copy a file containing build information. Useful for detective work on an otherwise unlabelled image.
# This is also useful for forcing docker to invalidate caches when the build target is the latest commit in a branch.
WORKDIR /home/lnode/eigenda/lightnode
COPY --chown=lnode docker/build-info.txt /home/lnode
WORKDIR /home/lnode

# Just in case we are tracking the latest commit in a branch, pull again. This is a no-op if the target is a commit sha.
WORKDIR /home/lnode/eigenda
RUN git pull
WORKDIR /home/lnode

# Build the light node
WORKDIR /home/lnode/eigenda/lightnode
RUN bash -c 'source ~/.bashrc && make build'
RUN ln -s /home/lnode/eigenda/lightnode/bin/lnode ~/lnode
WORKDIR /home/lnode

# Data that should persist between container restarts will be stored in this mount point.
RUN mkdir ~/data

# Data that shouldn't persist between container restarts will be stored in this location.
RUN mkdir ~/tmp

# Make everything except ~/tmp and ~/data read only to enforce good file system hygiene.
RUN chmod -R -w ~
RUN chmod -R +w ~/tmp
RUN chmod -R +w ~/data

# Run the light node when the container starts.
CMD ["/home/lnode/lnode"]
