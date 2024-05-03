package eigenda

// EigenDA finalized blob certificate
type Cert struct {
	BatchHeaderHash      []byte
	BlobIndex            uint32
	ReferenceBlockNumber uint32
	QuorumIDs            []uint32
}
