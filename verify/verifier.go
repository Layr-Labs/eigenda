package verify

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/ethereum/go-ethereum/log"

	binding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"

	"github.com/Layr-Labs/eigenda/api/grpc/common"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
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
	verifyCert bool
	prover     *prover.Prover
	cv         *CertVerifier
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

	prover, err := prover.NewProver(cfg.KzgConfig, false) // don't load G2 points
	if err != nil {
		return nil, err
	}

	return &Verifier{
		verifyCert: cfg.Verify,
		prover:     prover,
		cv:         cv,
	}, nil
}

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
		return err
	}

	// 2 - verify merkle inclusion proof
	err = v.cv.VerifyMerkleProof(cert.Proof().GetInclusionProof(), cert.BatchHeaderRoot(), cert.Proof().GetBlobIndex(), cert.ReadBlobHeader())
	if err != nil {
		return err
	}

	// 3 - verify security parameters
	err = v.VerifySecurityParams(cert.ReadBlobHeader(), header)
	if err != nil {
		return err
	}

	return nil
}

func (v *Verifier) Commit(blob []byte) (*bn254.G1Affine, error) {
	// ChunkLength and TotalChunks aren't relevant for computing data
	// commitment which is why they're currently set arbitrarily
	encoder, err := v.prover.GetKzgEncoder(
		encoding.ParamsFromSysPar(420, 69, uint64(len(blob))),
	)
	if err != nil {
		return nil, err
	}

	inputFr, err := rs.ToFrArray(blob)
	if err != nil {
		return nil, fmt.Errorf("cannot convert bytes to field elements, %w", err)
	}

	commit, err := encoder.Commit(inputFr)
	if err != nil {
		return nil, err
	}

	return &commit, nil
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
		errMsg += fmt.Sprintf("field elements do not match, x actual commit: %x, x expected commit: %x, ", actualCommit.X.Marshal(), (*expectedX).Marshal())
		errMsg += fmt.Sprintf("y actual commit: %x, y expected commit: %x", actualCommit.Y.Marshal(), (*expectedY).Marshal())
		return fmt.Errorf(errMsg)
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
