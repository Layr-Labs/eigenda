package verify

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/ethereum/go-ethereum/log"

	"github.com/Layr-Labs/eigenda/api/grpc/common"
	"github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	kzgverifier "github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/Layr-Labs/eigenda/encoding/rs"
)

const (
	HoleskySVCManagerV1Address = "0xD4A7E1Bd8015057293f0D0A557088c286942e84b"
)

type Config struct {
	KzgConfig   *kzg.KzgConfig
	VerifyCerts bool
	// below fields are only required if VerifyCerts is true
	RPCURL               string
	SvcManagerAddr       string
	EthConfirmationDepth uint64
	WaitForFinalization  bool
}

// Custom MarshalJSON function to control what gets included in the JSON output
func (c Config) MarshalJSON() ([]byte, error) {
	type Alias Config // Use an alias to avoid recursion with MarshalJSON
	aux := (Alias)(c)
	// Conditionally include a masked password if it is set
	if aux.RPCURL != "" {
		aux.RPCURL = "*****"
	}
	return json.Marshal(aux)
}

// TODO: right now verification and confirmation depth are tightly coupled. we should decouple them
type Verifier struct {
	// kzgVerifier is needed to commit blobs to the memstore
	kzgVerifier *kzgverifier.Verifier
	// cert verification is optional, and verifies certs retrieved from eigenDA when turned on
	verifyCerts bool
	cv          *CertVerifier
	// holesky is a flag to enable/disable holesky specific checks
	holesky bool
}

func NewVerifier(cfg *Config, l log.Logger) (*Verifier, error) {
	var cv *CertVerifier
	var err error

	if cfg.VerifyCerts {
		cv, err = NewCertVerifier(cfg, l)
		if err != nil {
			return nil, fmt.Errorf("failed to create cert verifier: %w", err)
		}
	}

	kzgVerifier, err := kzgverifier.NewVerifier(cfg.KzgConfig, false)
	if err != nil {
		return nil, fmt.Errorf("failed to create kzg verifier: %w", err)
	}

	return &Verifier{
		kzgVerifier: kzgVerifier,
		verifyCerts: cfg.VerifyCerts,
		cv:          cv,
		holesky:     isHolesky(cfg.SvcManagerAddr),
	}, nil
}

// verifies V0 eigenda certificate type
func (v *Verifier) VerifyCert(ctx context.Context, cert *Certificate) error {
	if !v.verifyCerts {
		return nil
	}

	// 1 - verify batch in the cert is confirmed onchain
	err := v.cv.verifyBatchConfirmedOnChain(ctx, cert.Proof().GetBatchId(), cert.Proof().GetBatchMetadata())
	if err != nil {
		return fmt.Errorf("failed to verify batch: %w", err)
	}

	// 2 - verify merkle inclusion proof
	err = v.cv.verifyMerkleProof(cert.Proof().GetInclusionProof(), cert.BatchHeaderRoot(), cert.Proof().GetBlobIndex(), cert.ReadBlobHeader())
	if err != nil {
		return fmt.Errorf("failed to verify merkle proof: %w", err)
	}

	// 3 - verify security parameters
	batchHeader := cert.Proof().GetBatchMetadata().GetBatchHeader()
	err = v.verifySecurityParams(cert.ReadBlobHeader(), batchHeader)
	if err != nil {
		return fmt.Errorf("failed to verify security parameters: %w", err)
	}

	return nil
}

// compute kzg-bn254 commitment of raw blob data using SRS
func (v *Verifier) Commit(blob []byte) (*bn254.G1Affine, error) {
	inputFr, err := rs.ToFrArray(blob)
	if err != nil {
		return nil, fmt.Errorf("cannot convert bytes to field elements, %w", err)
	}

	if len(v.kzgVerifier.Srs.G1) < len(inputFr) {
		return nil, fmt.Errorf("cannot verify commitment because the number of stored srs in the memory is insufficient, have %v need %v", len(v.kzgVerifier.Srs.G1), len(inputFr))
	}

	config := ecc.MultiExpConfig{}
	var commitment bn254.G1Affine
	_, err = commitment.MultiExp(v.kzgVerifier.Srs.G1[:len(inputFr)], inputFr, config)
	if err != nil {
		return nil, err
	}

	return &commitment, nil
}

// Verify regenerates a commitment from the blob and asserts equivalence
// to the commitment in the certificate
// TODO: Optimize implementation by opening a point on the commitment instead
func (v *Verifier) VerifyCommitment(certCommitment *common.G1Commitment, blob []byte) error {
	actualCommit, err := v.Commit(blob)
	if err != nil {
		return err
	}

	certCommitmentX := &fp.Element{}
	certCommitmentX.Unmarshal(certCommitment.X)
	certCommitmentY := &fp.Element{}
	certCommitmentY.Unmarshal(certCommitment.Y)

	certCommitmentAffine := bn254.G1Affine{
		X: *certCommitmentX,
		Y: *certCommitmentY,
	}

	if !certCommitmentAffine.IsOnCurve() {
		return fmt.Errorf("commitment (x,y) field elements are not on the BN254 curve")
	}

	errMsg := ""
	if !actualCommit.X.Equal(certCommitmentX) || !actualCommit.Y.Equal(certCommitmentY) {
		errMsg += fmt.Sprintf("field elements do not match, x actual commit: %x, x expected commit: %x, ", actualCommit.X.Marshal(), certCommitmentX.Marshal())
		errMsg += fmt.Sprintf("y actual commit: %x, y expected commit: %x", actualCommit.Y.Marshal(), certCommitmentY.Marshal())
		return fmt.Errorf("%s", errMsg)
	}

	return nil
}

// verifySecurityParams ensures that returned security parameters are valid
func (v *Verifier) verifySecurityParams(blobHeader BlobHeader, batchHeader *disperser.BatchHeader) error {
	confirmedQuorums := make(map[uint8]bool)

	// require that the security param in each blob is met
	for i := 0; i < len(blobHeader.QuorumBlobParams); i++ {
		if batchHeader.QuorumNumbers[i] != blobHeader.QuorumBlobParams[i].QuorumNumber {
			return fmt.Errorf("quorum number mismatch, expected: %d, got: %d", batchHeader.QuorumNumbers[i], blobHeader.QuorumBlobParams[i].QuorumNumber)
		}

		if blobHeader.QuorumBlobParams[i].AdversaryThresholdPercentage > blobHeader.QuorumBlobParams[i].ConfirmationThresholdPercentage {
			return fmt.Errorf("adversary threshold percentage must be greater than or equal to confirmation threshold percentage")
		}
		// we get the quorum adversary threshold at the batch's reference block number. This is not strictly needed right now
		// since this threshold is hardcoded into the contract: https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/core/EigenDAServiceManagerStorage.sol
		// but it is good practice in case the contract changes in the future
		quorumAdversaryThreshold, ok := v.cv.quorumAdversaryThresholds[blobHeader.QuorumBlobParams[i].QuorumNumber]
		if !ok {
			log.Warn("CertVerifier.quorumAdversaryThresholds map does not contain quorum number", "quorumNumber", blobHeader.QuorumBlobParams[i].QuorumNumber)
		} else if blobHeader.QuorumBlobParams[i].AdversaryThresholdPercentage < quorumAdversaryThreshold {
			return fmt.Errorf("adversary threshold percentage must be greater than or equal to quorum adversary threshold percentage")
		}

		if batchHeader.QuorumSignedPercentages[i] < blobHeader.QuorumBlobParams[i].ConfirmationThresholdPercentage {
			return fmt.Errorf("signed stake for quorum must be greater than or equal to confirmation threshold percentage")
		}

		confirmedQuorums[blobHeader.QuorumBlobParams[i].QuorumNumber] = true
	}

	// ensure that required quorums are present in the confirmed ones
	for _, quorum := range requiredQuorum(batchHeader.ReferenceBlockNumber, v) {
		if !confirmedQuorums[quorum] {
			return fmt.Errorf("quorum %d is required but not present in confirmed quorums", quorum)
		}
	}

	return nil
}

func requiredQuorum(referenceBlockNumber uint32, v *Verifier) []uint8 {
	// This check is required due to a bug we had when we updated the EigenDAServiceManager in Holesky. For a brief period of time, the quorum 1 was not
	// required for the commitment to be confirmed, so the disperser created batches with only quorum 0 signatures.
	// Archive nodes trying to sync from these stored batches would thus fail validation here since
	// quorumsRequired is read from the latestBlock, where the bug has been fixed and both quorums are required.
	// This check is only for testnet and for a specific block range.
	if v.holesky && referenceBlockNumber >= 2950000 && referenceBlockNumber < 2960000 {
		return []uint8{0}
	}
	return v.cv.quorumsRequired
}

func isHolesky(svcAddress string) bool {
	return strings.EqualFold(strings.TrimPrefix(svcAddress, "0x"), strings.TrimPrefix(HoleskySVCManagerV1Address, "0x"))
}
