package v2

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	commonpb "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	disperserpb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

type BlobVersion = uint16

// Assignment contains information about the set of chunks that a specific node will receive
type Assignment struct {
	StartIndex uint32
	NumChunks  uint32
}

// GetIndices generates the list of ChunkIndices associated with a given assignment
func (c *Assignment) GetIndices() []uint32 {
	indices := make([]uint32, c.NumChunks)
	for ind := range indices {
		indices[ind] = c.StartIndex + uint32(ind)
	}
	return indices
}

type BlobKey [32]byte

func (b BlobKey) Hex() string {
	return hex.EncodeToString(b[:])
}

func HexToBlobKey(h string) (BlobKey, error) {
	s := strings.TrimPrefix(h, "0x")
	s = strings.TrimPrefix(s, "0X")
	b, err := hex.DecodeString(s)
	if err != nil {
		return BlobKey{}, err
	}
	return BlobKey(b), nil
}

func BytesToBlobKey(bytes []byte) (BlobKey, error) {
	// Validate length
	if len(bytes) != 32 {
		return BlobKey{}, fmt.Errorf("invalid blob key length: expected 32 bytes, got %d", len(bytes))
	}

	var blobKey BlobKey
	copy(blobKey[:], bytes)
	return blobKey, nil
}

// BlobHeader contains all metadata related to a blob including commitments and parameters for encoding
type BlobHeader struct {
	BlobVersion BlobVersion

	BlobCommitments encoding.BlobCommitments

	// QuorumNumbers contains the quorums the blob is dispersed to
	QuorumNumbers []core.QuorumID

	// PaymentMetadata contains the payment information for the blob
	PaymentMetadata core.PaymentMetadata

	// Signature is the signature of the blob header by the account ID
	Signature []byte
}

func BlobHeaderFromProtobuf(proto *commonpb.BlobHeader) (*BlobHeader, error) {
	commitment, err := new(encoding.G1Commitment).Deserialize(proto.GetCommitment().GetCommitment())
	if err != nil {
		return nil, err
	}
	lengthCommitment, err := new(encoding.G2Commitment).Deserialize(proto.GetCommitment().GetLengthCommitment())
	if err != nil {
		return nil, err
	}
	lengthProof, err := new(encoding.LengthProof).Deserialize(proto.GetCommitment().GetLengthProof())
	if err != nil {
		return nil, err
	}

	if !(*bn254.G1Affine)(commitment).IsInSubGroup() {
		return nil, errors.New("commitment is not in the subgroup")
	}

	if !(*bn254.G2Affine)(lengthCommitment).IsInSubGroup() {
		return nil, errors.New("lengthCommitment is not in the subgroup")
	}

	if !(*bn254.G2Affine)(lengthProof).IsInSubGroup() {
		return nil, errors.New("lengthProof is not in the subgroup")
	}

	quorumNumbers := make([]core.QuorumID, len(proto.QuorumNumbers))
	for i, q := range proto.GetQuorumNumbers() {
		if q > MaxQuorumID {
			return nil, errors.New("quorum number exceeds maximum allowed")
		}
		quorumNumbers[i] = core.QuorumID(q)
	}

	paymentMetadata := core.ConvertToPaymentMetadata(proto.GetPaymentHeader())
	if paymentMetadata == nil {
		return nil, errors.New("payment metadata is nil")
	}

	return &BlobHeader{
		BlobVersion: BlobVersion(proto.GetVersion()),
		BlobCommitments: encoding.BlobCommitments{
			Commitment:       commitment,
			LengthCommitment: lengthCommitment,
			LengthProof:      lengthProof,
			Length:           uint(proto.GetCommitment().GetLength()),
		},
		QuorumNumbers:   quorumNumbers,
		PaymentMetadata: *paymentMetadata,
		Signature:       proto.GetSignature(),
	}, nil
}

func (b *BlobHeader) ToProtobuf() (*commonpb.BlobHeader, error) {
	quorums := make([]uint32, len(b.QuorumNumbers))
	for i, q := range b.QuorumNumbers {
		quorums[i] = uint32(q)
	}

	commitments, err := b.BlobCommitments.ToProtobuf()
	if err != nil {
		return nil, fmt.Errorf("failed to convert blob commitments to protobuf: %v", err)
	}

	return &commonpb.BlobHeader{
		Version:       uint32(b.BlobVersion),
		QuorumNumbers: quorums,
		Commitment:    commitments,
		PaymentHeader: b.PaymentMetadata.ToProtobuf(),
		Signature:     b.Signature,
	}, nil
}

func (b *BlobHeader) GetEncodingParams(blobParams *core.BlobVersionParameters) (encoding.EncodingParams, error) {
	length, err := GetChunkLength(uint32(b.BlobCommitments.Length), blobParams)
	if err != nil {
		return encoding.EncodingParams{}, err
	}

	return encoding.EncodingParams{
		NumChunks:   uint64(blobParams.NumChunks),
		ChunkLength: uint64(length),
	}, nil
}

type RelayKey = uint32

type BlobCertificate struct {
	BlobHeader *BlobHeader

	// RelayKeys
	RelayKeys []RelayKey
}

func (c *BlobCertificate) ToProtobuf() (*commonpb.BlobCertificate, error) {
	if c.BlobHeader == nil {
		return nil, fmt.Errorf("blob header is nil")
	}

	blobHeader, err := c.BlobHeader.ToProtobuf()
	if err != nil {
		return nil, fmt.Errorf("failed to convert blob header to protobuf: %v", err)
	}

	relays := make([]uint32, len(c.RelayKeys))
	for i, r := range c.RelayKeys {
		relays[i] = uint32(r)
	}

	return &commonpb.BlobCertificate{
		BlobHeader: blobHeader,
		Relays:     relays,
	}, nil
}

func BlobCertificateFromProtobuf(proto *commonpb.BlobCertificate) (*BlobCertificate, error) {
	if proto.GetBlobHeader() == nil {
		return nil, errors.New("missing blob header in blob certificate")
	}

	blobHeader, err := BlobHeaderFromProtobuf(proto.GetBlobHeader())
	if err != nil {
		return nil, fmt.Errorf("failed to create blob header: %v", err)
	}

	relayKeys := make([]RelayKey, len(proto.GetRelays()))
	for i, r := range proto.GetRelays() {
		relayKeys[i] = RelayKey(r)
	}

	return &BlobCertificate{
		BlobHeader: blobHeader,
		RelayKeys:  relayKeys,
	}, nil
}

type BatchHeader struct {
	// BatchRoot is the root of a Merkle tree whose leaves are the keys of the blobs in the batch
	BatchRoot [32]byte
	// ReferenceBlockNumber is the block number at which all operator information (stakes, indexes, etc.) is taken from
	ReferenceBlockNumber uint64
}

func (h *BatchHeader) ToProtobuf() *commonpb.BatchHeader {
	return &commonpb.BatchHeader{
		BatchRoot:            h.BatchRoot[:],
		ReferenceBlockNumber: h.ReferenceBlockNumber,
	}
}

type Batch struct {
	BatchHeader      *BatchHeader
	BlobCertificates []*BlobCertificate
}

func (b *Batch) ToProtobuf() (*commonpb.Batch, error) {
	if b.BatchHeader == nil {
		return nil, errors.New("batch header is nil")
	}

	if b.BatchHeader.BatchRoot == [32]byte{} {
		return nil, errors.New("batch root is empty")
	}

	if b.BatchHeader.ReferenceBlockNumber == 0 {
		return nil, errors.New("reference block number is 0")
	}

	blobCerts := make([]*commonpb.BlobCertificate, len(b.BlobCertificates))
	for i, cert := range b.BlobCertificates {
		blobCert, err := cert.ToProtobuf()
		if err != nil {
			return nil, fmt.Errorf("failed to convert blob certificate to protobuf: %v", err)
		}
		blobCerts[i] = blobCert
	}

	return &commonpb.Batch{
		Header: &commonpb.BatchHeader{
			BatchRoot:            b.BatchHeader.BatchRoot[:],
			ReferenceBlockNumber: b.BatchHeader.ReferenceBlockNumber,
		},
		BlobCertificates: blobCerts,
	}, nil
}

func BatchFromProtobuf(proto *commonpb.Batch) (*Batch, error) {
	if len(proto.GetBlobCertificates()) == 0 {
		return nil, errors.New("missing blob certificates in batch")
	}

	if proto.GetHeader() == nil {
		return nil, errors.New("missing header in batch")
	}

	if len(proto.GetHeader().GetBatchRoot()) != 32 {
		return nil, errors.New("batch root must be 32 bytes")
	}

	batchHeader := &BatchHeader{
		BatchRoot:            [32]byte(proto.GetHeader().GetBatchRoot()),
		ReferenceBlockNumber: proto.GetHeader().GetReferenceBlockNumber(),
	}

	blobCerts := make([]*BlobCertificate, len(proto.GetBlobCertificates()))
	for i, cert := range proto.GetBlobCertificates() {
		blobHeader, err := BlobHeaderFromProtobuf(cert.GetBlobHeader())
		if err != nil {
			return nil, fmt.Errorf("failed to create blob header: %v", err)
		}

		blobCerts[i] = &BlobCertificate{
			BlobHeader: blobHeader,
			RelayKeys:  make([]RelayKey, len(cert.GetRelays())),
		}
		for j, r := range cert.GetRelays() {
			blobCerts[i].RelayKeys[j] = RelayKey(r)
		}
	}

	return &Batch{
		BatchHeader:      batchHeader,
		BlobCertificates: blobCerts,
	}, nil
}

type Attestation struct {
	*BatchHeader

	// AttestedAt is the time the attestation was made
	AttestedAt uint64
	// NonSignerPubKeys are the public keys of the operators that did not sign the blob
	NonSignerPubKeys []*core.G1Point
	// APKG2 is the aggregate public key of all signers
	APKG2 *core.G2Point
	// QuorumAPKs is the aggregate public keys of all operators in each quorum
	QuorumAPKs map[core.QuorumID]*core.G1Point
	// Sigma is the aggregate signature of all signers
	Sigma *core.Signature
	// QuorumNumbers contains the quorums relevant for the attestation
	QuorumNumbers []core.QuorumID
	// QuorumResults contains the results of the quorum verification
	QuorumResults map[core.QuorumID]uint8
}

func (a *Attestation) ToProtobuf() (*disperserpb.Attestation, error) {
	nonSignerPubKeys := make([][]byte, len(a.NonSignerPubKeys))
	for i, p := range a.NonSignerPubKeys {
		pubkeyBytes := p.Bytes()
		nonSignerPubKeys[i] = pubkeyBytes[:]
	}

	quorumAPKs := make([][]byte, len(a.QuorumAPKs))
	quorumNumbers := make([]uint32, len(a.QuorumNumbers))
	quorumResults := make([]uint8, len(a.QuorumResults))
	for i, q := range a.QuorumNumbers {
		quorumNumbers[i] = uint32(q)

		apk, ok := a.QuorumAPKs[q]
		if !ok {
			return nil, fmt.Errorf("missing quorum APK for quorum %d", q)
		}
		apkBytes := apk.Bytes()
		quorumAPKs[i] = apkBytes[:]
		quorumResults[i] = a.QuorumResults[q]
	}

	apkG2Bytes := a.APKG2.Bytes()
	sigmaBytes := a.Sigma.Bytes()

	return &disperserpb.Attestation{
		NonSignerPubkeys:        nonSignerPubKeys,
		ApkG2:                   apkG2Bytes[:],
		QuorumApks:              quorumAPKs,
		Sigma:                   sigmaBytes[:],
		QuorumNumbers:           quorumNumbers,
		QuorumSignedPercentages: quorumResults,
	}, nil
}

type BlobVerificationInfo struct {
	*BatchHeader

	BlobKey
	BlobIndex      uint32
	InclusionProof []byte
}

func (v *BlobVerificationInfo) ToProtobuf(blobCert *BlobCertificate) (*disperserpb.BlobVerificationInfo, error) {
	blobCertProto, err := blobCert.ToProtobuf()
	if err != nil {
		return nil, err
	}
	return &disperserpb.BlobVerificationInfo{
		BlobCertificate: blobCertProto,
		BlobIndex:       v.BlobIndex,
		InclusionProof:  v.InclusionProof,
	}, nil
}

// DispersalRequest is a request to disperse a batch to a specific operator
type DispersalRequest struct {
	core.OperatorID `dynamodbav:"-"`
	OperatorAddress gethcommon.Address
	Socket          string
	DispersedAt     uint64

	BatchHeader
}

// DispersalResponse is a response to a dispersal request
type DispersalResponse struct {
	*DispersalRequest

	RespondedAt uint64
	// Signature is the signature of the response by the operator
	Signature [32]byte
	// Error is the error message if the dispersal failed
	Error string
}

const (
	// We use uint8 to count the number of quorums, so we can have at most 255 quorums,
	// which means the max ID can not be larger than 254 (from 0 to 254, there are 255
	// different IDs).
	MaxQuorumID = 254
)
