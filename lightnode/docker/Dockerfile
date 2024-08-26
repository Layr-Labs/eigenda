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
# Remove the default .bashrc, this isn't something we want to use.
RUN rm .bashrc && touch .bashrc

# Install golang
RUN wget https://go.dev/dl/go1.21.12.linux-arm64.tar.gz # TODO make this an argument
RUN tar -xf ./*.tar.gz
RUN rm ./*.tar.gz
RUN mkdir -p ~/.local/share
RUN mv go ~/.local/share/go

# Setup golang environment for lnode user
RUN echo 'export GOPATH=~/.local/share/go' >> ~/.bashrc
RUN echo 'export PATH=~/.local/share/go/bin:$PATH' >> ~/.bashrc

# Data that should persist between container restarts will be stored in this mount point.
RUN mkdir ~/data

# Clone eigenda repository
RUN git clone $GIT_URL eigenda
WORKDIR /home/lnode/eigenda
RUN git fetch # TODO necessary?
RUN git checkout $BRANCH_OR_COMMIT
RUN git pull # TODO necessary?
WORKDIR /home/lnode

# TODO make the eigenda directory read only, why not

# Build the light node
WORKDIR /home/lnode/eigenda/lightnode
RUN bash -c 'cat ~/.bashrc && source ~/.bashrc && echo $PATH && make build'
WORKDIR /home/lnode

# Run the light node when the container starts.
CMD eigenda/lightnode/bin/lnode
