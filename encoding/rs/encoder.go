package rs

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	_ "go.uber.org/automaxprocs"
)

type Config struct {
	NumWorker int
}

type Encoder struct {
	*Config
	NumRSWorker         int
	mu                  sync.Mutex
	ParametrizedEncoder map[encoding.EncodingParams]*ParametrizedEncoder
	verbose             bool
}

// Proof device represents a device capable of computing reed-solomon operations.
type EncoderDevice interface {
	ExtendPolyEval(coeffs []fr.Element) ([]fr.Element, error)
}

// // RsComputeDevice represents a device capable of performing Reed-Solomon encoding computations.
// // Implementations of this interface are expected to handle polynomial evaluation extensions.
// type RsComputeDevice interface {
// 	// ExtendPolyEval extends the evaluation of a polynomial given its coefficients.
// 	// It takes a slice of polynomial coefficients and returns an extended evaluation.
// 	//
// 	// Parameters:
// 	//   - coeffs: A slice of fr.Element representing the polynomial coefficients.
// 	//
// 	// Returns:
// 	//   - A slice of fr.Element representing the extended polynomial evaluation.
// 	//   - An error if the extension process fails.
// 	ExtendPolyEval(coeffs []fr.Element) ([]fr.Element, error)
// }

func NewEncoder() (*Encoder, error) {
	fmt.Println("rs numthread", runtime.GOMAXPROCS(0))
	return &Encoder{
		Config: &Config{
			NumWorker: runtime.GOMAXPROCS(0),
		},
		ParametrizedEncoder: make(map[encoding.EncodingParams]*ParametrizedEncoder),
	}, nil
}

// The function creates a high level struct that determines the encoding the a data of a
// specific length under (num systematic node, num parity node) setup. A systematic node
// stores a systematic data chunk that contains part of the original data. A parity node
// stores a parity data chunk which is an encoding of the original data. A receiver that
// collects all systematic chunks can simply stitch data together to reconstruct the
// original data. When some systematic chunks are missing but identical parity chunk are
// available, the receive can go through a Reed Solomon decoding to reconstruct the
// original data.
func (g *Encoder) GetRsEncoder(params encoding.EncodingParams) (*ParametrizedEncoder, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	enc, ok := g.ParametrizedEncoder[params]
	if ok {
		return enc, nil
	}

	enc, err := g.newEncoder(params)
	if err == nil {
		g.ParametrizedEncoder[params] = enc
	}

	return enc, err
}
