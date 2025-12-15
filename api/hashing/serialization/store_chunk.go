package serialization

import (
	"encoding/binary"
	"fmt"
	"math"

	commonv2 "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	grpc "github.com/Layr-Labs/eigenda/api/grpc/validator"
)

// initialStoreChunksRequestCap is just a preallocation hint to reduce allocations.
// It's not a limit: `append` will grow the slice as needed.
const initialStoreChunksRequestCap = 512

const initialBlobHeaderCap = 512

// validatorStoreChunksRequestDomain is the StoreChunksRequest hash domain prefix.
// Kept here to avoid an import cycle (hashing <-> serialization).
const validatorStoreChunksRequestDomain = "validator.StoreChunksRequest"

type canonicalStoreChunksRequest struct {
	Domain           string
	BatchHeader      canonicalBatchHeader
	BlobCertificates []canonicalBlobCertificate
	DisperserID      uint32
	Timestamp        uint32
}

type canonicalBatchHeader struct {
	Root                 []byte
	ReferenceBlockNumber uint64
}

type canonicalBlobCertificate struct {
	BlobHeader canonicalBlobHeader
	Signature  []byte
	RelayKeys  []uint32
}

type canonicalBlobHeader struct {
	Version uint32
	// TODO(taras): QuorumNumbersLength is redundant. As QuorumNumbers is a list and length will
	// the first uint32 in the list
	QuorumNumbersLength uint32
	QuorumNumbers       []uint32
	Commitment          canonicalBlobCommitment
	PaymentHeader       canonicalPaymentHeader
}

type canonicalBlobCommitment struct {
	Commitment       []byte
	LengthCommitment []byte
	LengthProof      []byte
	Length           uint32
}

type canonicalPaymentHeader struct {
	AccountId         string
	Timestamp         int64
	CumulativePayment []byte
}

func (h canonicalBatchHeader) serialize(dst []byte) ([]byte, error) {
	var err error
	dst, err = serializeBytes(dst, h.Root)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize batch header root: %w", err)
	}
	dst = serializeU64(dst, h.ReferenceBlockNumber)
	return dst, nil
}

func (c canonicalBlobCommitment) serialize(dst []byte) ([]byte, error) {
	var err error
	dst, err = serializeBytes(dst, c.Commitment)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize commitment: %w", err)
	}
	dst, err = serializeBytes(dst, c.LengthCommitment)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize length commitment: %w", err)
	}
	dst, err = serializeBytes(dst, c.LengthProof)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize length proof: %w", err)
	}
	dst = serializeU32(dst, c.Length)
	return dst, nil
}

func (h canonicalPaymentHeader) serialize(dst []byte) ([]byte, error) {
	var err error
	dst, err = serializeBytes(dst, []byte(h.AccountId))
	if err != nil {
		return nil, fmt.Errorf("failed to serialize account id: %w", err)
	}
	dst = serializeI64(dst, h.Timestamp)
	dst, err = serializeBytes(dst, h.CumulativePayment)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize cumulative payment: %w", err)
	}
	return dst, nil
}

func (h canonicalBlobHeader) serialize(dst []byte) ([]byte, error) {
	dst = serializeU32(dst, h.Version)
	dst = serializeU32(dst, h.QuorumNumbersLength)
	var err error
	dst, err = serializeU32Slice(dst, h.QuorumNumbers)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize quorum numbers: %w", err)
	}
	dst, err = h.Commitment.serialize(dst)
	if err != nil {
		return nil, err
	}
	dst, err = h.PaymentHeader.serialize(dst)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func (c canonicalBlobCertificate) serialize(dst []byte) ([]byte, error) {
	var err error
	dst, err = c.BlobHeader.serialize(dst)
	if err != nil {
		return nil, err
	}
	dst, err = serializeBytes(dst, c.Signature)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize signature: %w", err)
	}
	dst, err = serializeU32Slice(dst, c.RelayKeys)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize relay keys: %w", err)
	}
	return dst, nil
}

func (r canonicalStoreChunksRequest) serialize(dst []byte) ([]byte, error) {
	dst = append(dst, []byte(r.Domain)...)

	var err error
	dst, err = r.BatchHeader.serialize(dst)
	if err != nil {
		return nil, err
	}

	if len(r.BlobCertificates) > math.MaxUint32 {
		return nil, fmt.Errorf("array is too long: %d", len(r.BlobCertificates))
	}
	dst = serializeU32(dst, uint32(len(r.BlobCertificates)))
	for i, cert := range r.BlobCertificates {
		dst, err = cert.serialize(dst)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize blob certificate at index %d: %w", i, err)
		}
	}

	dst = serializeU32(dst, r.DisperserID)
	dst = serializeU32(dst, r.Timestamp)
	return dst, nil
}

func serializeU32(dst []byte, v uint32) []byte {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], v)
	return append(dst, b[:]...)
}

func serializeU64(dst []byte, v uint64) []byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], v)
	return append(dst, b[:]...)
}

func serializeI64(dst []byte, v int64) []byte {
	return serializeU64(dst, uint64(v))
}

func serializeBytes(dst []byte, b []byte) ([]byte, error) {
	if len(b) > math.MaxUint32 {
		return nil, fmt.Errorf("byte array is too long: %d", len(b))
	}
	dst = serializeU32(dst, uint32(len(b)))
	dst = append(dst, b...)
	return dst, nil
}

func serializeU32Slice(dst []byte, s []uint32) ([]byte, error) {
	if len(s) > math.MaxUint32 {
		return nil, fmt.Errorf("uint32 array is too long: %d", len(s))
	}
	dst = serializeU32(dst, uint32(len(s)))
	for _, v := range s {
		dst = serializeU32(dst, v)
	}
	return dst, nil
}

func SerializeStoreChunksRequest(request *grpc.StoreChunksRequest) ([]byte, error) {
	if request.GetBatch() == nil || request.GetBatch().GetHeader() == nil {
		return nil, fmt.Errorf("missing batch/header")
	}

	canonicalRequest := canonicalStoreChunksRequest{
		Domain: validatorStoreChunksRequestDomain,
		BatchHeader: canonicalBatchHeader{
			Root:                 request.GetBatch().GetHeader().GetBatchRoot(),
			ReferenceBlockNumber: request.GetBatch().GetHeader().GetReferenceBlockNumber(),
		},
		BlobCertificates: make([]canonicalBlobCertificate, len(request.GetBatch().GetBlobCertificates())),
		DisperserID:      request.GetDisperserID(),
		Timestamp:        request.GetTimestamp(),
	}
	for i, cert := range request.GetBatch().GetBlobCertificates() {
		if cert == nil || cert.GetBlobHeader() == nil ||
			cert.GetBlobHeader().GetCommitment() == nil ||
			cert.GetBlobHeader().GetPaymentHeader() == nil ||
			cert.GetSignature() == nil ||
			cert.GetRelayKeys() == nil {
			return nil, fmt.Errorf("missing blob certificate fields at index %d", i)
		}
		canonicalRequest.BlobCertificates[i] = canonicalBlobCertificate{
			BlobHeader: canonicalBlobHeader{
				Version:             cert.GetBlobHeader().GetVersion(),
				QuorumNumbersLength: uint32(len(cert.GetBlobHeader().GetQuorumNumbers())),
				QuorumNumbers:       cert.GetBlobHeader().GetQuorumNumbers(),
				Commitment: canonicalBlobCommitment{
					Commitment:       cert.GetBlobHeader().GetCommitment().GetCommitment(),
					LengthCommitment: cert.GetBlobHeader().GetCommitment().GetLengthCommitment(),
					LengthProof:      cert.GetBlobHeader().GetCommitment().GetLengthProof(),
					Length:           cert.GetBlobHeader().GetCommitment().GetLength(),
				},
				PaymentHeader: canonicalPaymentHeader{
					AccountId:         cert.GetBlobHeader().GetPaymentHeader().GetAccountId(),
					Timestamp:         cert.GetBlobHeader().GetPaymentHeader().GetTimestamp(),
					CumulativePayment: cert.GetBlobHeader().GetPaymentHeader().GetCumulativePayment(),
				},
			},
			Signature: cert.GetSignature(),
			RelayKeys: cert.GetRelayKeys(),
		}
	}

	out := make([]byte, 0, initialStoreChunksRequestCap)
	return canonicalRequest.serialize(out)
}

func SerializeBlobHeader(header *commonv2.BlobHeader) ([]byte, error) {
	if header == nil || header.GetCommitment() == nil || header.GetPaymentHeader() == nil {
		return nil, fmt.Errorf("missing blob header fields")
	}
	canonicalHeader := canonicalBlobHeader{
		Version:             header.GetVersion(),
		QuorumNumbersLength: uint32(len(header.GetQuorumNumbers())),
		QuorumNumbers:       header.GetQuorumNumbers(),
		Commitment: canonicalBlobCommitment{
			Commitment: header.GetCommitment().GetCommitment(),
		},
		PaymentHeader: canonicalPaymentHeader{
			AccountId:         header.GetPaymentHeader().GetAccountId(),
			Timestamp:         header.GetPaymentHeader().GetTimestamp(),
			CumulativePayment: header.GetPaymentHeader().GetCumulativePayment(),
		},
	}
	out := make([]byte, 0, initialBlobHeaderCap)
	return canonicalHeader.serialize(out)
}
