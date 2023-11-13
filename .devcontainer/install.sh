# Install foundry
curl -L https://foundry.paradigm.xyz | bash
~/.foundry/bin/foundryup 

# Install go dependencies
go install github.com/onsi/ginkgo/v2/ginkgo@v2.2.0
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
# go install github.com/mikefarah/yq/v4@latest

# yarn global add @graphprotocol/graph-cli@0.51.0