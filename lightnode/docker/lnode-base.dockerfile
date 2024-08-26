# syntax=docker/dockerfile:1

# The base libraries required to build and run the light node.

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
