package verify

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/api/grpc/common"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/ethereum/go-ethereum/log"

	proxy_common "github.com/Layr-Labs/eigenda-proxy/common"

	binding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"

	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/rs"
)

type Config struct {
	Verify         bool
	RPCURL         string
	SvcManagerAddr string
	KzgConfig      *kzg.KzgConfig
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

func (v *Verifier) VerifyCert(cert *proxy_common.Certificate) error {
	if !v.verifyCert {
		return nil
	}

	// 1 - verify batch

	header := binding.IEigenDAServiceManagerBatchHeader{
		BlobHeadersRoot:       [32]byte(cert.GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetBatchRoot()),
		QuorumNumbers:         cert.GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetQuorumNumbers(),
		ReferenceBlockNumber:  cert.GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetReferenceBlockNumber(),
		SignedStakeForQuorums: cert.GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetQuorumSignedPercentages(),
	}

	err := v.cv.VerifyBatch(&header, cert.BlobVerificationProof.BatchId, [32]byte(cert.BlobVerificationProof.BatchMetadata.SignatoryRecordHash), cert.BlobVerificationProof.BatchMetadata.GetConfirmationBlockNumber())
	if err != nil {
		return err
	}

	// 2 - TODO: verify merkle proof

	// 3 - TODO: verify security params
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

	// convert to field elements
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
