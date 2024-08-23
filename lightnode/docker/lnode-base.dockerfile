# lnode-base is the base image for light nodes.

FROM debian

# Install core dependencies
RUN apt update
RUN apt install -y wget
RUN apt install -y git
RUN apt install -y build-essential

# Install golang
RUN wget https://go.dev/dl/go1.21.12.linux-arm64.tar.gz
RUN tar -xf ./*.tar.gz
RUN rm ./*.tar.gz
RUN mv go /opt/go

# Set up lnode user
RUN useradd -m -s /bin/bash lnode
USER lnode
WORKDIR /home/lnode

# Data that should persist between container restarts will be stored in this mount point.
RUN mkdir ~/data

# Setup golang environment for lnode user
RUN echo 'export GOPATH=/opt/go' >> ~/.bashrc
RUN echo 'export PATH=/opt/go/bin:$PATH' >> ~/.bashrc

