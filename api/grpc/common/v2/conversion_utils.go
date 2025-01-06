package v2

import (
	"fmt"
	verifierBindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDABlobVerifier"
	"github.com/Layr-Labs/eigenda/core"
	"math"
)

// ToBinding converts a BatchHeader into a contractEigenDABlobVerifier.BatchHeaderV2
func (h *BatchHeader) ToBinding() (*verifierBindings.BatchHeaderV2, error) {
	var outputBatchRoot [32]byte

	inputBatchRoot := h.GetBatchRoot()
	if len(inputBatchRoot) != 32 {
		return nil, fmt.Errorf("BatchRoot must be 32 bytes (length was %d)", len(inputBatchRoot))
	}
	copy(outputBatchRoot[:], inputBatchRoot[:])

	inputReferenceBlockNumber := h.GetReferenceBlockNumber()
	if inputReferenceBlockNumber > math.MaxUint32 {
		return nil, fmt.Errorf(
			"ReferenceBlockNumber overflow: value was %d, but max allowable value is %d",
			inputReferenceBlockNumber,
			math.MaxUint32)
	}

	convertedHeader := &verifierBindings.BatchHeaderV2{
		BatchRoot:            outputBatchRoot,
		ReferenceBlockNumber: uint32(inputReferenceBlockNumber),
	}

	return convertedHeader, nil
}

// ToBinding converts a BlobCertificate into a contractEigenDABlobVerifier.BlobCertificate
func (c *BlobCertificate) ToBinding() (*verifierBindings.BlobCertificate, error) {
	convertedBlobHeader, err := c.GetBlobHeader().toBinding()
	if err != nil {
		return nil, fmt.Errorf("convert blob header: %s", err)
	}

	return &verifierBindings.BlobCertificate{
		BlobHeader: *convertedBlobHeader,
		RelayKeys:  c.GetRelays(),
	}, nil
}

// toBinding converts a BlobHeader into a contractEigenDABlobVerifier.BlobHeaderV2
func (h *BlobHeader) toBinding() (*verifierBindings.BlobHeaderV2, error) {
	inputVersion := h.GetVersion()
	if inputVersion > math.MaxUint16 {
		return nil, fmt.Errorf(
			"version overflow: value was %d, but max allowable value is %d",
			inputVersion,
			math.MaxUint16)
	}

	var quorumNumbers []byte
	for _, quorumNumber := range h.GetQuorumNumbers() {
		if quorumNumber > math.MaxUint8 {
			return nil, fmt.Errorf("quorum number overflow: value was %d, but max allowable value is %d", quorumNumber, uint8(math.MaxUint8))
		}

		quorumNumbers = append(quorumNumbers, byte(quorumNumber))
	}

	convertedBlobCommitment, err := h.GetCommitment().ToBinding()
	if err != nil {
		return nil, fmt.Errorf("convert blob commitment: %s", err)
	}

	paymentHeaderHash, err := core.ConvertToPaymentMetadata(h.GetPaymentHeader()).Hash()
	if err != nil {
		return nil, fmt.Errorf("hash payment header: %s", err)
	}

	return &verifierBindings.BlobHeaderV2{
		Version:           uint16(inputVersion),
		QuorumNumbers:     quorumNumbers,
		Commitment:        *convertedBlobCommitment,
		PaymentHeaderHash: paymentHeaderHash,
	}, nil
}
