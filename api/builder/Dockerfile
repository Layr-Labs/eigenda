FROM debian

# Install core dependencies
RUN apt update
RUN apt install -y wget unzip bash build-essential

# Set up user
RUN useradd -m -s /bin/bash user
USER user
WORKDIR /home/user
# Remove default crud
RUN rm .bashrc
RUN rm .bash_logout
RUN rm .profile

# Install golang
RUN wget https://go.dev/dl/go1.21.12.linux-arm64.tar.gz
RUN tar -xf ./*.tar.gz
RUN rm ./*.tar.gz
RUN mkdir -p ~/.local/share
RUN mv go ~/.local/share/go

# Install protoc
RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v23.4/protoc-23.4-linux-aarch_64.zip
RUN mkdir protoc
RUN cd protoc && unzip ../*.zip
RUN rm ./*.zip
RUN mv protoc ~/.local/share/protoc

# Setup PATH
RUN touch ~/.bashrc
RUN echo 'export GOPATH=~/.local/share/go' >> ~/.bashrc
RUN echo 'export PATH=~/.local/share/go/bin:$PATH' >> ~/.bashrc
RUN echo 'export PATH=~/.local/share/protoc/bin:$PATH' >> ~/.bashrc

# Install go protobuf extensions
RUN bash -c 'source ~/.bashrc && go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1'
RUN bash -c 'source ~/.bashrc && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0'
