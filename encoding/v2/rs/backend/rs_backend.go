package backend

import (
	"context"
	"github.com/Layr-Labs/eigenda/encoding/v2/rs/backend/gnark"
	"github.com/Layr-Labs/eigenda/encoding/v2/rs/backend/icicle"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// Proof device represents a device capable of computing reed-solomon operations.
type RSEncoderBackend interface {
	ExtendPolyEvalV2(ctx context.Context, coeffs []fr.Element) ([]fr.Element, error)
}

// We implement two backends: gnark and icicle.
//   - Gnark uses the gnark library and is the default CPU-based backend, and is always available.
//   - Icicle uses the icicle library and can leverage GPU acceleration, but requires building with the icicle tag.
//     Building with the icicle tag will inject the dynamic libraries required to use icicle.
var _ RSEncoderBackend = &gnark.RSBackend{}
var _ RSEncoderBackend = &icicle.RSBackend{}
