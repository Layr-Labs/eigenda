package srs

import (
	_ "embed"
	"fmt"

	"github.com/consensys/gnark-crypto/ecc/bn254"
)

//go:embed g2.point.powerOf2
var serializedG2PowerOf2Data []byte

// G2PowerOf2SRS contains 28 G2 points: [1] [tau^1] [tau^2]..[tau^{2^27}]
//
// the powerOf2 file, only [tau^exp] are stored.
// exponent    0,    1,       2,    , ..
// actual pow [tau],[tau^2],[tau^4],.. (stored in the file)
// In our convention SRSOrder contains the total number of series of g1, g2 starting with generator
// i.e. [1] [tau] [tau^2]..
// So the actual power of tau is SRSOrder - 1
// The mainnet SRS, the max power is 2^28-1, so the last power in powerOf2 file is [tau^(2^27)]
var G2PowerOf2SRS []bn254.G2Affine

func init() {
	// Note that we can't use bn254.NewDecoder(bytes.NewReader(g2PowerOf2Data)).Decode(&G2PowerOf2Data)
	// because the file was not encoded using gnark-crypto's encoder.
	// It only contains the 28 raw serialized points, each taking 64 bytes.
	// gnark-crypto's encoder/decoder adds a 4 bytes header.
	for pointIndex := 0; pointIndex < len(serializedG2PowerOf2Data); pointIndex += 64 {
		serializedPoint := serializedG2PowerOf2Data[pointIndex : pointIndex+64]
		var p bn254.G2Affine
		if _, err := p.SetBytes(serializedPoint); err != nil {
			panic(fmt.Sprintf("error deserializing G2 point at index %d: %v", pointIndex, err))
		}
		G2PowerOf2SRS = append(G2PowerOf2SRS, p)
	}
}
