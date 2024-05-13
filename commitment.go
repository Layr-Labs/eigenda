package plasma

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"

	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/op-plasma-eigenda/eigenda"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"

	eigen_enc "github.com/Layr-Labs/eigenda/encoding"
)

// ErrInvalidCommitment is returned when the commitment cannot be parsed into a known commitment type.
var ErrInvalidCommitment = errors.New("invalid commitment")

// ErrCommitmentMismatch is returned when the commitment does not match the given input.
var ErrCommitmentMismatch = errors.New("commitment mismatch")

// CommitmentType is the commitment type prefix.
type CommitmentType byte

const (
	// default commitment type for the DA storage.
	Keccak256CommitmentType CommitmentType = 0x00
	DaService               CommitmentType = 0x01
)

type ExtDAType byte

const (
	EigenDA ExtDAType = 0x00
)

type EigenDAVersion byte

const (
	EigenV0 EigenDAVersion = 0x00
)

// Keccak256Commitment is the default commitment type for op-plasma.
type Keccak256Commitment []byte

// Encode adds a commitment type prefix self describing the commitment.
func (c Keccak256Commitment) Encode() []byte {
	return append([]byte{byte(Keccak256CommitmentType)}, c...)
}

// TxData adds an extra version byte to signal it's a commitment.
func (c Keccak256Commitment) TxData() []byte {
	return append([]byte{TxDataVersion1}, c.Encode()...)
}

// Verify checks if the commitment matches the given input.
func (c Keccak256Commitment) Verify(input []byte) error {
	if !bytes.Equal(c, crypto.Keccak256(input)) {
		return ErrCommitmentMismatch
	}
	return nil
}

// Keccak256 creates a new commitment from the given input.
func Keccak256(input []byte) Keccak256Commitment {
	return Keccak256Commitment(crypto.Keccak256(input))
}

// DecodeKeccak256 validates and casts the commitment into a Keccak256Commitment.
func DecodeKeccak256(commitment []byte) (Keccak256Commitment, error) {
	if len(commitment) == 0 {
		return nil, ErrInvalidCommitment
	}
	if commitment[0] != byte(Keccak256CommitmentType) {
		return nil, ErrInvalidCommitment
	}
	c := commitment[1:]
	if len(c) != 32 {
		return nil, ErrInvalidCommitment
	}
	return c, nil
}

// NOTE - This logic will need to be migrated into layr-labs/op-stack directly
type EigenDACommitment []byte

func (c EigenDACommitment) Encode() []byte {
	return append([]byte{byte(DaService), byte(EigenDA), byte(EigenV0)}, c...)
}

func (c EigenDACommitment) TxData() []byte {
	return append([]byte{TxDataVersion1}, c.Encode()...)
}

// TODO - verify the commitment against the input blob by evaluating its polynomial representation at an arbitrary point
// and asserting that the generated output proof can be successfully verified against the commitment.
func (c EigenDACommitment) Verify(input []byte) error {
	var cert eigenda.Cert
	if err := rlp.DecodeBytes(c, &cert); err != nil {
		return err
	}

	blob := eigenda.EncodeToBlob(input)

	kzgConfig := &kzg.KzgConfig{
		G1Path:          "./kzg/g1.point",
		G2Path:          "./kzg/g2.point",
		G2PowerOf2Path:  "./kzg/g2.point.powerOf2",
		CacheDir:        "./kzg/SRSTables",
		SRSOrder:        3000,
		SRSNumberToLoad: 3000,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
	}

	prover, err := prover.NewProver(kzgConfig, true)
	if err != nil {
		return err
	}

	params := eigen_enc.ParamsFromSysPar(6, 69, uint64(len(blob)))
	encoder, err := prover.GetKzgEncoder(params)
	if err != nil {
		return err
	}

	commit, _, _, _, _, err := encoder.EncodeBytes(blob)
	if err != nil {
		return err
	}

	x, y := cert.BlobCommitmentFields()
	xCommit := &commit.X
	yCommit := &commit.Y

	println(fmt.Sprintf("%+b, %+b ", commit.X.Bits(), x.Bits()))
	println(fmt.Sprintf("%+b, %+b ", commit.Y.Bits(), y.Bits()))

	if x.NotEqual(xCommit) != 0 {
		return fmt.Errorf("x element mismatch %s:%s %s:%s", "gen_commit", xCommit.String(), "initial_commit", x.String())
	}

	if y.NotEqual(yCommit) != 0 {
		return fmt.Errorf("x element mismatch %s:%s %s:%s", "gen_commit", yCommit.String(), "initial_commit", y.String())
	}

	return nil
}

func DecodeEigenDACommitment(commitment []byte) (EigenDACommitment, error) {
	if len(commitment) <= 3 {
		return nil, ErrInvalidCommitment
	}
	if commitment[0] != byte(DaService) {
		return nil, ErrInvalidCommitment
	}

	if commitment[1] != byte(EigenDA) {
		return nil, ErrInvalidCommitment
	}

	// additional versions will need to be hardcoded here
	if commitment[2] != byte(EigenV0) {
		return nil, ErrInvalidCommitment
	}

	c := commitment[3:]

	// TODO - Add a length check
	return c, nil
}
