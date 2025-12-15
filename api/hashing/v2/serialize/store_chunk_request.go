package serialize

import (
	"bytes"
	"fmt"
	"math"

	commonv2 "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	grpc "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/lunixbochs/struc"
)

// This file provides a struc-based encoder that preserves the *exact* byte layout of
// store_chunk.go's manual serializer.
//
// Wire-format invariants preserved:
// - Big-endian integers (struc defaults to big-endian)
// - For []byte and []uint32: a uint32 length prefix, followed by elements/bytes
// - Domain is written as raw bytes (no length prefix) at the start
// - Redundant QuorumNumbersLength field is preserved (it appears before the slice length prefix)

// initialBufCap is a preallocation hint to reduce allocations.
const initialBufCap = 512

// validatorStoreChunksRequestDomain is the StoreChunksRequest hash domain prefix.
// Kept here to avoid an import cycle (hashing <-> serialization).
const validatorStoreChunksRequestDomain = "validator.StoreChunksRequest"

type canonicalStoreChunksRequestBody struct {
	BatchHeader canonicalBatchHeader

	BlobCertificatesLen uint32 `struc:"uint32,sizeof=BlobCertificates"`
	BlobCertificates    []canonicalBlobCertificate

	DisperserID uint32
	Timestamp   uint32
}

type canonicalBatchHeader struct {
	RootLen uint32 `struc:"uint32,sizeof=Root"`
	Root    []byte

	ReferenceBlockNumber uint64
}

type canonicalBlobCertificate struct {
	BlobHeader canonicalBlobHeader

	SignatureLen uint32 `struc:"uint32,sizeof=Signature"`
	Signature    []byte

	RelayKeysLen uint32 `struc:"uint32,sizeof=RelayKeys"`
	RelayKeys    []uint32
}

type canonicalBlobHeader struct {
	Version uint32

	// Kept for backwards-compatible encoding: this is written first...
	QuorumNumbersLength uint32
	// ...then the real slice length prefix (same value) followed by elements.
	QuorumNumbersLen uint32 `struc:"uint32,sizeof=QuorumNumbers"`
	QuorumNumbers    []uint32

	Commitment    canonicalBlobCommitment
	PaymentHeader canonicalPaymentHeader
}

type canonicalBlobCommitment struct {
	CommitmentLen uint32 `struc:"uint32,sizeof=Commitment"`
	Commitment    []byte

	LengthCommitmentLen uint32 `struc:"uint32,sizeof=LengthCommitment"`
	LengthCommitment    []byte

	LengthProofLen uint32 `struc:"uint32,sizeof=LengthProof"`
	LengthProof    []byte

	Length uint32
}

type canonicalPaymentHeader struct {
	// store_chunk.go encodes AccountId as serializeBytes([]byte(string))
	AccountIdLen uint32 `struc:"uint32,sizeof=AccountId"`
	AccountId    []byte

	Timestamp int64

	CumulativePaymentLen uint32 `struc:"uint32,sizeof=CumulativePayment"`
	CumulativePayment    []byte
}

func SerializeStoreChunksRequest(request *grpc.StoreChunksRequest) ([]byte, error) {
	if request.GetBatch() == nil || request.GetBatch().GetHeader() == nil {
		return nil, fmt.Errorf("missing batch/header")
	}

	certs := request.GetBatch().GetBlobCertificates()
	if len(certs) > math.MaxUint32 {
		return nil, fmt.Errorf("array is too long: %d", len(certs))
	}

	body := canonicalStoreChunksRequestBody{
		BatchHeader: canonicalBatchHeader{
			Root:                 request.GetBatch().GetHeader().GetBatchRoot(),
			ReferenceBlockNumber: request.GetBatch().GetHeader().GetReferenceBlockNumber(),
		},
		BlobCertificates: make([]canonicalBlobCertificate, len(certs)),
		DisperserID:      request.GetDisperserID(),
		Timestamp:        request.GetTimestamp(),
	}

	for i, cert := range certs {
		if cert == nil || cert.GetBlobHeader() == nil ||
			cert.GetBlobHeader().GetCommitment() == nil ||
			cert.GetBlobHeader().GetPaymentHeader() == nil {
			return nil, fmt.Errorf("missing blob certificate fields at index %d", i)
		}

		bh := cert.GetBlobHeader()
		commitment := bh.GetCommitment()
		payment := bh.GetPaymentHeader()

		qnums := bh.GetQuorumNumbers()
		qnLen := uint32(len(qnums))

		body.BlobCertificates[i] = canonicalBlobCertificate{
			BlobHeader: canonicalBlobHeader{
				Version:             bh.GetVersion(),
				QuorumNumbersLength: qnLen,
				QuorumNumbers:       qnums,
				Commitment: canonicalBlobCommitment{
					Commitment:       commitment.GetCommitment(),
					LengthCommitment: commitment.GetLengthCommitment(),
					LengthProof:      commitment.GetLengthProof(),
					Length:           commitment.GetLength(),
				},
				PaymentHeader: canonicalPaymentHeader{
					AccountId:         []byte(payment.GetAccountId()),
					Timestamp:         payment.GetTimestamp(),
					CumulativePayment: payment.GetCumulativePayment(),
				},
			},
			Signature: cert.GetSignature(),
			RelayKeys: cert.GetRelayKeys(),
		}
	}

	var buf bytes.Buffer
	buf.Grow(initialBufCap)

	_, _ = buf.WriteString(validatorStoreChunksRequestDomain)

	if err := struc.Pack(&buf, &body); err != nil {
		return nil, fmt.Errorf("failed to pack canonical StoreChunksRequest: %w", err)
	}
	return buf.Bytes(), nil
}

func SerializeBlobHeader(header *commonv2.BlobHeader) ([]byte, error) {
	if header == nil || header.GetCommitment() == nil || header.GetPaymentHeader() == nil {
		return nil, fmt.Errorf("missing blob header fields")
	}

	qnums := header.GetQuorumNumbers()
	qnLen := uint32(len(qnums))

	// Preserve current SerializeBlobHeader behavior from store_chunk.go:
	// it only sets Commitment.Commitment and leaves the rest empty/zero.
	ch := canonicalBlobHeader{
		Version:             header.GetVersion(),
		QuorumNumbersLength: qnLen,
		QuorumNumbers:       qnums,
		Commitment: canonicalBlobCommitment{
			Commitment:       header.GetCommitment().GetCommitment(),
			LengthCommitment: header.GetCommitment().GetLengthCommitment(),
			LengthProof:      header.GetCommitment().GetLengthProof(),
			Length:           header.GetCommitment().GetLength(),
		},
		PaymentHeader: canonicalPaymentHeader{
			AccountId:         []byte(header.GetPaymentHeader().GetAccountId()),
			Timestamp:         header.GetPaymentHeader().GetTimestamp(),
			CumulativePayment: header.GetPaymentHeader().GetCumulativePayment(),
		},
	}

	var buf bytes.Buffer
	buf.Grow(initialBufCap)

	if err := struc.Pack(&buf, &ch); err != nil {
		return nil, fmt.Errorf("failed to pack canonical BlobHeader: %w", err)
	}
	return buf.Bytes(), nil
}
