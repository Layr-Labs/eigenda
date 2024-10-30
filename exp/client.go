package main

import (
	"context"

	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
	"github.com/Layr-Labs/eigenda/encoding"
)

func client() {
	// List of blob keys and their expected sizes
	blobInfos := []struct {
		key      string
		size     int
		chunkLen int
	}{
		// {"886e35653b423c0ab4b6d7b78d24c8d85bd7a4c159c1983935f607aebfe0cf6b", 2048, 1},
		// {"a934a37212568f1c5d01b784f849eff8db639373df80650b4ffc98015ebb8abf", 4096, 1},
		// {"2ecd9f89c210d078ee33f36e96f76c6dadb1fef03ad4f6a1339243eb6ad00804", 8192, 1},
		// {"f8d9ac616cca03dc43dac510bf3c7aab7294af6acb3039da43758962ef4f0670", 16384, 1},
		// {"58af7b30ad77e85138073ac1aa3d38c46b5a8d6d7192c841c426af74771c7591", 32768, 1},
		// {"550f63e5d246a3c11d708d537a49a9dc7cb4c1cd7d8252ea0afaee2714769200", 65536, 2},
		// {"9047cb63013ba1068190ee23afc20e1ba4c26019dc8bd98b234a3d909861d3a1", 131072, 4},
		// {"a73df72f43ddff57cfc6fd00b424a8d0a8c257a0f0915d8cfcba70c6cce69bbc", 262144, 8},
		// {"3d451abde44f0cc4ef26af6a5cc1981dee4ceeeebbeff62ea413459ab7a071d5", 524288, 16},
		// {"bbedd2c78b43aafe134c8a0a48afa9d54664a4eef817f71bacc377842bd3812b", 1048576, 32},
		// {"cd1d7e1a2b94bfe7b12712279dd53beeabe8d44522e55273ad73c61f02c98dc0", 2097152, 64},
		// {"0e43067dcd6425acc76f24753ca6cd9acaca2806f10c855763956009e3fd277b", 4194304, 128},
		// {"4e56cf9801f0ebfe1e59f9b412f2892abbd60b1810675179bc8df48451e0b8b6", 8388608, 256},
		// {"d7897a41b4581572646b1b717dbb38d5a2f7915a5d185a43b1f4efaa16d23559", 16777216, 512},
		{"7e1e9f91bdce51901df4c3196c239e50bf6d09fb4d46b603372d740999952cf4", 33554432, 1024},
	}

	// Create encoder client
	encoderClient, err := encoder.NewEncoderClientV2("localhost:34000")
	if err != nil {
		panic(err)
	}

	// Encode blob
	for _, info := range blobInfos {
		blobKey, err := corev2.HexToBlobKey(info.key)
		if err != nil {
			panic(err)
		}

		// Set encoding params
		encoderParams := encoding.EncodingParams{
			ChunkLength: uint64(info.chunkLen),
			NumChunks:   8192,
		}

		_, err = encoderClient.EncodeBlob(context.Background(), blobKey, encoderParams)
		if err != nil {
			panic(err)
		}
	}
}
