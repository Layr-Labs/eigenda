package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
)

// Useful for converting back and forth between 4844 padded base64 representations of
// unicode input data, for testing purposes.
//
// An example:
//
//	grpcurl \
//		-proto ./api/proto/disperser/disperser.proto \
//		-import-path ./api/proto \
//		-d '{"data": "'$(tools/kzgpad/bin/kzgpad -e hello)'"}' \
//		disperser-holesky.eigenda.xyz:443 disperser.Disperser/DisperseBlob
//
// Then poll for confirmation using GetBlobStatus, then retrieve blob:
//
//	grpcurl \
//	  -import-path ./api/proto \
//	  -proto ./api/proto/disperser/disperser.proto \
//	  -d '{"batch_header_hash": "INSERT_VALUE", "blob_index":"INSERT_VALUE"}' \
//	  disperser-holesky.eigenda.xyz:443 disperser.Disperser/RetrieveBlob | \
//	  jq -r .data | \
//	  tools/kzgpad/bin/kzgpad -d -

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: go run main.go [-e|-d] [input]")
		os.Exit(1)
	}

	mode := os.Args[1]
	input := os.Args[2]

	if input == "-" {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			processInput(mode, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading stdin:", err)
			os.Exit(1)
		}
	} else {
		processInput(mode, input)
	}
}

func processInput(mode, text string) {
	switch mode {
	case "-e":
		// Encode the input to base64
		bz := []byte(text)
		padded := codec.ConvertByPaddingEmptyByte(bz)
		encoded := base64.StdEncoding.EncodeToString(padded)
		fmt.Println(encoded)
	case "-d":
		// Decode the base64 input
		decoded, err := base64.StdEncoding.DecodeString(text)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error decoding base64:", err)
			return
		}
		unpadded := codec.RemoveEmptyByteFromPaddedBytes(decoded)
		fmt.Println(string(unpadded))
	default:
		fmt.Fprintln(os.Stderr, "Invalid mode. Use -e for encoding or -d for decoding.")
	}
}
