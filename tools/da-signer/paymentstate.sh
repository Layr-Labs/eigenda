#set -x

# Get current timestamp in nanoseconds
TIMESTAMP=$(date +%s%N)

# Sign using cast and ensure we get a 65-byte signature
SIGNATURE=$(go run da-signer.go $ACCOUNT_ADDR $ACCOUNT_PKEY $TIMESTAMP)

# Make the gRPC request with all required fields
set -x
grpcurl -d "{
  \"account_id\": \"$ACCOUNT_ADDR\",
  \"signature\": \"$SIGNATURE\",
  \"timestamp\": $TIMESTAMP
}" disperser-v2-testnet-sepolia.eigenda.xyz:443 disperser.v2.Disperser/GetPaymentState
