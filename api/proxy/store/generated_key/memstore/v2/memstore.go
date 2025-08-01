package memstore

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/certs"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/memstore/ephemeraldb"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/memstore/memconfig"
	cert_types_binding "github.com/Layr-Labs/eigenda/contracts/bindings/IEigenDACertTypeBindings"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

const (
	BytesPerFieldElement = 32
)

// unsafeRandomBytes ... Generates random byte slice provided
// size. Errors when generating are ignored since this is only
// used for constructing dummy certificates when testing insecure integrations.
// in the worst case it doesn't work and returns empty arrays which would only
// impact memstore operation in the event that two identical payloads are provided
// since they'd resolve to the same commitment and blob key. This shouldn't matter
// given this is typically used for testing standard E2E functionality against a rollup
// stack which SHOULD never submit an identical batch more than once.
func unsafeRandomBytes(size uint) []byte {
	entropy := make([]byte, size)
	_, _ = rand.Read(entropy)
	return entropy
}

func unsafeRandInt(maxValue int64) *big.Int {
	randInt, _ := rand.Int(rand.Reader, big.NewInt(maxValue))
	return randInt
}

func unsafeRandCeilAt32() uint32 {
	// #nosec G115 - downcasting only on random value
	return uint32(unsafeRandInt(32).Uint64())
}

/*
MemStore is a simple in-memory store for blobs which uses an expiration
time to evict blobs to best emulate the ephemeral nature of blobs dispersed to
EigenDA V2 operators.
*/
type MemStore struct {
	// keccak(RLP(randomlyGeneratedCert)) -> Blob
	*ephemeraldb.DB
	log logging.Logger

	g1SRS []bn254.G1Affine
	codec codecs.BlobCodec
}

var _ common.EigenDAV2Store = (*MemStore)(nil)

// New ... constructor
func New(
	ctx context.Context, log logging.Logger, config *memconfig.SafeConfig,
	g1SRS []bn254.G1Affine,
) (*MemStore, error) {
	return &MemStore{
		ephemeraldb.New(ctx, config, log),
		log,
		g1SRS,
		codecs.NewIFFTCodec(codecs.NewDefaultBlobCodec()),
	}, nil
}

// generateRandomCert ... generates a pseudo random EigenDA V2 certificate
func (e *MemStore) generateRandomCert(blobContents []byte) (coretypes.EigenDACert, error) {
	// compute kzg data commitment. this is useful for testing
	// READPREIMAGE functionality in the arbitrum x eigenda integration since
	// preimage key is computed within the VM from hashing a recomputation of the data
	// commitment
	dataCommitment, err := verification.GenerateBlobCommitment(e.g1SRS, blobContents)
	if err != nil {
		return nil, err
	}

	x := dataCommitment.X.BigInt(&big.Int{})
	y := dataCommitment.Y.BigInt(&big.Int{})

	g1CommitPoint := cert_types_binding.BN254G1Point{
		X: x,
		Y: y,
	}

	pseudoRandomBlobInclusionInfo := cert_types_binding.EigenDATypesV2BlobInclusionInfo{
		BlobCertificate: cert_types_binding.EigenDATypesV2BlobCertificate{
			BlobHeader: cert_types_binding.EigenDATypesV2BlobHeaderV2{
				Version:       0,                            // only supported version as of now
				QuorumNumbers: []byte{byte(0x0), byte(0x1)}, // quorum 0 && quorum 1
				Commitment: cert_types_binding.EigenDATypesV2BlobCommitment{
					LengthCommitment: cert_types_binding.BN254G2Point{
						X: [2]*big.Int{unsafeRandInt(1000), unsafeRandInt(1000)},
						Y: [2]*big.Int{unsafeRandInt(1000), unsafeRandInt(1000)},
					},
					LengthProof: cert_types_binding.BN254G2Point{
						X: [2]*big.Int{unsafeRandInt(1), unsafeRandInt(1)},
						Y: [2]*big.Int{unsafeRandInt(1), unsafeRandInt(1)},
					},
					Commitment: g1CommitPoint,
					// #nosec G115 - can never overflow on 16MiB blobs
					Length: uint32(len(blobContents)) / BytesPerFieldElement,
				},
				PaymentHeaderHash: [32]byte(unsafeRandomBytes(32)),
			},
			Signature: unsafeRandomBytes(48), // 384 bits
			RelayKeys: []uint32{unsafeRandCeilAt32(), unsafeRandCeilAt32()},
		},
		// #nosec G115 - max value 1000 guaranteed to be safe for uint32
		BlobIndex:      uint32(unsafeRandInt(1_000).Uint64()),
		InclusionProof: unsafeRandomBytes(128),
	}

	randomBatchHeader := cert_types_binding.EigenDATypesV2BatchHeaderV2{
		BatchRoot: [32]byte(unsafeRandomBytes(32)),
		// increase the rbn of cert to a high enough number 4294967200 < 2^32 = 4294967296
		// where random part is chosen from 0 to 32. So there is no chance of overflow.
		// a large RBN is useful to avoid failing the recency check when testing
		// See https://github.com/Layr-Labs/eigenda/blob/master/docs/spec/src/integration/spec/6-secure-integration.md
		// where the check is often done by checking the failure condition
		// certL1InclusionBlock > RecencyWindowSize + cert.RBN
		// once we increase the RBN, the above failure condition will never trigger
		ReferenceBlockNumber: unsafeRandCeilAt32() + 4294967200,
	}

	randomNonSignerStakesAndSigs := cert_types_binding.EigenDATypesV1NonSignerStakesAndSignature{
		NonSignerQuorumBitmapIndices: []uint32{unsafeRandCeilAt32(), unsafeRandCeilAt32()},
		NonSignerPubkeys: []cert_types_binding.BN254G1Point{
			{
				X: unsafeRandInt(1000),
				Y: unsafeRandInt(1000),
			},
		},
		QuorumApks: []cert_types_binding.BN254G1Point{
			{
				X: unsafeRandInt(1000),
				Y: unsafeRandInt(1000),
			},
		},
		ApkG2: cert_types_binding.BN254G2Point{
			X: [2]*big.Int{unsafeRandInt(1000), unsafeRandInt(10000)},
			Y: [2]*big.Int{unsafeRandInt(1000), unsafeRandInt(1000)},
		},
		QuorumApkIndices:  []uint32{unsafeRandCeilAt32(), unsafeRandCeilAt32()},
		TotalStakeIndices: []uint32{unsafeRandCeilAt32(), unsafeRandCeilAt32(), unsafeRandCeilAt32()},
		NonSignerStakeIndices: [][]uint32{
			{unsafeRandCeilAt32(), unsafeRandCeilAt32()},
			{unsafeRandCeilAt32(), unsafeRandCeilAt32()},
		},
		Sigma: cert_types_binding.BN254G1Point{
			X: unsafeRandInt(1000),
			Y: unsafeRandInt(1000),
		},
	}

	return &coretypes.EigenDACertV3{
		BlobInclusionInfo:           pseudoRandomBlobInclusionInfo,
		BatchHeader:                 randomBatchHeader,
		NonSignerStakesAndSignature: randomNonSignerStakesAndSigs,
	}, nil
}

// Get fetches a value from the store.
// If returnEncodedPayload is true, it returns the encoded blob without decoding.
func (e *MemStore) Get(
	_ context.Context, versionedCert certs.VersionedCert, returnEncodedPayload bool,
) ([]byte, error) {
	encodedBlob, err := e.FetchEntry(crypto.Keccak256Hash(versionedCert.SerializedCert).Bytes())
	if err != nil {
		return nil, fmt.Errorf("fetching entry via v2 memstore: %w", err)
	}

	if returnEncodedPayload {
		return encodedBlob, nil
	}

	return e.codec.DecodeBlob(encodedBlob)
}

// Put inserts a value into the store.
// ephemeral db key = keccak256(pseudo_random_cert)
// this is done to verify that a rollup must be able to provide
// the same certificate used in dispersal for retrieval
func (e *MemStore) Put(_ context.Context, value []byte) ([]byte, error) {
	encodedVal, err := e.codec.EncodeBlob(value)
	if err != nil {
		return nil, err
	}

	artificialV2Cert, err := e.generateRandomCert(encodedVal)
	if err != nil {
		return nil, fmt.Errorf("generating random cert: %w", err)
	}

	certBytes, err := rlp.EncodeToBytes(artificialV2Cert)
	if err != nil {
		return nil, fmt.Errorf("rlp decode v2 cert: %w", err)
	}

	err = e.InsertEntry(crypto.Keccak256Hash(certBytes).Bytes(), encodedVal)
	if err != nil { // don't wrap here so api.ErrorFailover{} isn't modified
		return nil, err
	}

	return certBytes, nil
}

func (e *MemStore) VerifyCert(
	_ context.Context, _ certs.VersionedCert, _ uint64,
) error {
	return nil
}

func (e *MemStore) BackendType() common.BackendType {
	return common.MemstoreV2BackendType
}
