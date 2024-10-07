FROM golang:1.21.12-bookworm

# The URL where the protoc binary can be downloaded. Is different depending on architecture.
ARG PROTOC_URL

# The UID of the user to create
ARG UID

# Install core dependencies
RUN apt update
RUN apt install -y wget unzip bash

# Set up user
RUN useradd -u $UID -m -s /bin/bash user
USER user
WORKDIR /home/user
# Remove default crud
RUN rm .bashrc
RUN rm .bash_logout
RUN rm .profile

# Install protoc
RUN wget $PROTOC_URL
RUN mkdir protoc
RUN cd protoc && unzip ../*.zip
RUN rm ./*.zip

# Setup PATH
RUN touch ~/.bashrc
RUN echo 'export PATH=~/protoc/bin:$PATH' >> ~/.bashrc
RUN echo 'export GOPATH=/go' >> ~/.bashrc
RUN echo 'export PATH=/usr/local/go/bin:$PATH' >> ~/.bashrc

# Install go protobuf extensions
RUN bash -c 'source ~/.bashrc && go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1'
RUN bash -c 'source ~/.bashrc && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0'
