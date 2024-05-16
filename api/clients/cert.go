package clients

import "github.com/Layr-Labs/eigenda/api/grpc/common"

type Cert struct {
	BatchHeaderHash      []byte
	BlobIndex            uint32
	ReferenceBlockNumber uint32
	QuorumIDs            []uint32

	// Used for kzg verification when reading blob data from DA
	BlobCommitment *common.G1Commitment
}
