package verify

import (
	"fmt"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/ethereum/go-ethereum/log"

	binding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"

	"github.com/Layr-Labs/eigenda/api/grpc/common"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	kzgverifier "github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/Layr-Labs/eigenda/encoding/rs"
)

type Config struct {
	Verify               bool
	RPCURL               string
	SvcManagerAddr       string
	KzgConfig            *kzg.KzgConfig
	EthConfirmationDepth uint64
}

type Verifier struct {
	verifyCert  bool
	kzgVerifier *kzgverifier.Verifier
	cv          *CertVerifier
}

func NewVerifier(cfg *Config, l log.Logger) (*Verifier, error) {
	var cv *CertVerifier
	var err error

	if cfg.Verify {
		cv, err = NewCertVerifier(cfg, l)
		if err != nil {
			return nil, err
		}
	}

	kzgVerifier, err := kzgverifier.NewVerifier(cfg.KzgConfig, false)
	if err != nil {
		return nil, err
	}

	return &Verifier{
		verifyCert:  cfg.Verify,
		kzgVerifier: kzgVerifier,
		cv:          cv,
	}, nil
}

// verifies V0 eigenda certificate type
func (v *Verifier) VerifyCert(cert *Certificate) error {
	if !v.verifyCert {
		return nil
	}

	// 1 - verify batch
	header := binding.IEigenDAServiceManagerBatchHeader{
		BlobHeadersRoot:       [32]byte(cert.Proof().GetBatchMetadata().GetBatchHeader().GetBatchRoot()),
		QuorumNumbers:         cert.Proof().GetBatchMetadata().GetBatchHeader().GetQuorumNumbers(),
		ReferenceBlockNumber:  cert.Proof().GetBatchMetadata().GetBatchHeader().GetReferenceBlockNumber(),
		SignedStakeForQuorums: cert.Proof().GetBatchMetadata().GetBatchHeader().GetQuorumSignedPercentages(),
	}

	err := v.cv.VerifyBatch(&header, cert.Proof().GetBatchId(), [32]byte(cert.Proof().BatchMetadata.GetSignatoryRecordHash()), cert.Proof().BatchMetadata.GetConfirmationBlockNumber())
	if err != nil {
		return fmt.Errorf("failed to verify batch: %w", err)
	}

	// 2 - verify merkle inclusion proof
	err = v.cv.VerifyMerkleProof(cert.Proof().GetInclusionProof(), cert.BatchHeaderRoot(), cert.Proof().GetBlobIndex(), cert.ReadBlobHeader())
	if err != nil {
		return fmt.Errorf("failed to verify merkle proof: %w", err)
	}

	// 3 - verify security parameters
	err = v.VerifySecurityParams(cert.ReadBlobHeader(), header)
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
func (v *Verifier) VerifyCommitment(expectedCommit *common.G1Commitment, blob []byte) error {
	actualCommit, err := v.Commit(blob)
	if err != nil {
		return err
	}

	expectedX := &fp.Element{}
	expectedX.Unmarshal(expectedCommit.X)
	expectedY := &fp.Element{}
	expectedY.Unmarshal(expectedCommit.Y)

	errMsg := ""
	if !actualCommit.X.Equal(expectedX) || !actualCommit.Y.Equal(expectedY) {
		errMsg += fmt.Sprintf("field elements do not match, x actual commit: %x, x expected commit: %x, ", actualCommit.X.Marshal(), expectedX.Marshal())
		errMsg += fmt.Sprintf("y actual commit: %x, y expected commit: %x", actualCommit.Y.Marshal(), expectedY.Marshal())
		return fmt.Errorf("%s", errMsg)
	}

	return nil
}

// VerifySecurityParams ensures that returned security parameters are valid
func (v *Verifier) VerifySecurityParams(blobHeader BlobHeader, batchHeader binding.IEigenDAServiceManagerBatchHeader) error {
	confirmedQuorums := make(map[uint8]bool)

	// require that the security param in each blob is met
	for i := 0; i < len(blobHeader.QuorumBlobParams); i++ {
		if batchHeader.QuorumNumbers[i] != blobHeader.QuorumBlobParams[i].QuorumNumber {
			return fmt.Errorf("quorum number mismatch, expected: %d, got: %d", batchHeader.QuorumNumbers[i], blobHeader.QuorumBlobParams[i].QuorumNumber)
		}

		if blobHeader.QuorumBlobParams[i].AdversaryThresholdPercentage > blobHeader.QuorumBlobParams[i].ConfirmationThresholdPercentage {
			return fmt.Errorf("adversary threshold percentage must be greater than or equal to confirmation threshold percentage")
		}

		quorumAdversaryThreshold, err := v.getQuorumAdversaryThreshold(blobHeader.QuorumBlobParams[i].QuorumNumber)
		if err != nil {
			log.Warn("failed to get quorum adversary threshold", "err", err)
		}

		if quorumAdversaryThreshold > 0 && blobHeader.QuorumBlobParams[i].AdversaryThresholdPercentage < quorumAdversaryThreshold {
			return fmt.Errorf("adversary threshold percentage must be greater than or equal to quorum adversary threshold percentage")
		}

		if batchHeader.SignedStakeForQuorums[i] < blobHeader.QuorumBlobParams[i].ConfirmationThresholdPercentage {
			return fmt.Errorf("signed stake for quorum must be greater than or equal to confirmation threshold percentage")
		}

		confirmedQuorums[blobHeader.QuorumBlobParams[i].QuorumNumber] = true
	}

	requiredQuorums, err := v.cv.manager.QuorumNumbersRequired(nil)
	if err != nil {
		log.Warn("failed to get required quorum numbers", "err", err)
	}

	// ensure that required quorums are present in the confirmed ones
	for _, quorum := range requiredQuorums {
		if !confirmedQuorums[quorum] {
			return fmt.Errorf("quorum %d is required but not present in confirmed quorums", quorum)
		}
	}

	return nil
}

// getQuorumAdversaryThreshold reads the adversarial threshold percentage for a given quorum number
// returns 0 if DNE
func (v *Verifier) getQuorumAdversaryThreshold(quorumNum uint8) (uint8, error) {
	percentages, err := v.cv.manager.QuorumAdversaryThresholdPercentages(nil)
	if err != nil {
		return 0, err
	}

	if len(percentages) > int(quorumNum) {
		return percentages[quorumNum], nil
	}

	return 0, nil
}
