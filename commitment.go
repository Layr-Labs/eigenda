package plasma

import (
	"bytes"
	"errors"

	"github.com/ethereum/go-ethereum/crypto"
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

// is redstone?
type DaLayer byte

const (
	RedStone    DaLayer = 0x00
	NotRedstone DaLayer = 0x01
)

type ExtDAType byte

const (
	EigenDA ExtDAType = 0x00
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
	return append([]byte{byte(DaService), byte(NotRedstone), byte(EigenDA)}, c...)
}

func (c EigenDACommitment) TxData() []byte {
	return append([]byte{TxDataVersion1}, c.Encode()...)
}

// TODO - verify the commitment against the input blob by evaluating its polynomial representation at an arbitrary point
// and asserting that the generated output proof can be successfully verified against the commitment.
func (c EigenDACommitment) Verify(input []byte) error {
	// Cast to certificate type
	// var cert eigenda.Cert
	// if err := rlp.DecodeBytes(c, &cert); err != nil {
	// 	return err
	// }

	// input = eigenda.EncodeToBlob(input)
	// commit, err := eigenda.ComputeCommitmentToData(input)
	// if err != nil {
	// 	return err
	// }

	// kzgConfig := &kzg.KzgConfig{
	// 	G1Path:          "./kzg/g1.point",
	// 	G2Path:          "./kzg/g2.point",
	// 	G2PowerOf2Path:  "./kzg/g2.point.powerOf2",
	// 	CacheDir:        "./kzg/SRSTables",
	// 	SRSOrder:        3000,
	// 	SRSNumberToLoad: 3000,
	// 	NumWorker:       uint64(runtime.GOMAXPROCS(0)),
	// }

	// // One pad encoding
	// // Fr - Frame is a structured set of field elements
	// validInput := codec.ConvertByPaddingEmptyByte(c)
	// inputFr, err := rs.ToFrArray(validInput)
	// if err != nil {
	// 	return err
	// }

	// frLen := uint64(len(inputFr))
	// paddedInputFr := make([]fr.Element, encoding.NextPowerOf2(frLen))
	// // pad input Fr to power of 2 for computing FFT
	// for i := 0; i < len(paddedInputFr); i++ {
	// 	if i < len(inputFr) {
	// 		paddedInputFr[i].Set(&inputFr[i])
	// 	} else {
	// 		println("Padding zeros")
	// 		paddedInputFr[i].SetZero()
	// 	}
	// }

	// group, err := kzgProver.NewProver(kzgConfig, true)
	// if err != nil {
	// 	return err
	// }

	// xElement, yElement := fp.Element{}, fp.Element{}

	// xElement = *xElement.SetBytes(cert.BlobCommitment.X)
	// yElement = *yElement.SetBytes(cert.BlobCommitment.Y)

	// // RS encoding to get polynomial representation of data
	// // numSys := uint64(4)
	// // numPar := uint64(0)

	// // params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(validInput)))
	// // enc, err := group.GetKzgEncoder(params)
	// // if err != nil {
	// // 	return err
	// // }

	// // Lagrange basis SRS in normal order, not butterfly
	// // lagrangeG1SRS, err := enc.Fs.FFTG1(group.Srs.G1[:len(paddedInputFr)], true)
	// // if err != nil {
	// // 	return err
	// // }

	// // commit in lagrange form
	// commitLagrange, err := oc.CommitInLagrange(inputFr, group.Srs.G1[:len(inputFr)])
	// if err != nil {
	// 	return err
	// }

	// println(fmt.Sprintf("%+x", commitLagrange.X), len(commitLagrange.X))
	// println(fmt.Sprintf("%+x", commitLagrange.Y), len(commitLagrange.Y))

	// println(fmt.Sprintf("%+x", cert.BlobCommitment.X), len(cert.BlobCommitment.X))
	// println(fmt.Sprintf("%+x", cert.BlobCommitment.Y), len(cert.BlobCommitment.Y))

	// if commit.X.String() != xElement.String() {
	// 	return fmt.Errorf("x element mismatch %s : %s, %s : %s", "generated_commit", commit.X.String(), "initial_commit", xElement.String())
	// }
	// if commit.Y.String() != yElement.String() {
	// 	return fmt.Errorf("x element mismatch %s : %s, %s : %s", "generated_commit", commit.Y.String(), "initial_commit", yElement.String())
	// }

	return nil
}

func DecodeEigenDACommitment(commitment []byte) (EigenDACommitment, error) {
	if len(commitment) <= 3 {
		return nil, ErrInvalidCommitment
	}
	if commitment[0] != byte(DaService) {
		return nil, ErrInvalidCommitment
	}

	if commitment[1] != byte(NotRedstone) {
		return nil, ErrInvalidCommitment
	}

	if commitment[2] != byte(EigenDA) {
		return nil, ErrInvalidCommitment
	}

	c := commitment[3:]

	// TODO - Add a length check
	return c, nil
}
